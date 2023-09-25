package compute

import (
	"context"
	"fmt"
	"strings"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func InstancesIamPolicy() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_instance_role_binding",
		Description: `https://cloud.google.com/compute/docs/reference/rest/v1/instances/getIamPolicy`,
		Resolver:    fetchInstanceIamPolicy,
		Multiplex:   client.ProjectMultiplexEnabledServices("iam.googleapis.com"),
		Transform:   transformers.TransformWithStruct(&assetpb.IamPolicyAnalysisResult{}),
		Columns: []schema.Column{
			{
				Name:     "project_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveProject,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
			{
				Name:     "instance_id",
				Type:     schema.TypeInt,
				Resolver: schema.ParentColumnResolver("id"),
			},
		},
	}
}

// It will return duplicate rows
func fetchInstanceIamPolicy(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	p := parent.Item.(*computepb.Instance)
	client, err := asset.NewClient(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	// Get the instance
	zoneSlice := strings.Split(*p.Zone, "/")
	zone := zoneSlice[len(zoneSlice)-1]
	scope := fmt.Sprintf("projects/%s", c.ProjectId)
	fullResourceName := fmt.Sprintf("//compute.googleapis.com/projects/%s/zones/%s/instances/%s", c.ProjectId, zone, *p.Name)
	if c.FolderId == "" && c.OrgId == "" {
		return analyze(client, scope, fullResourceName, ctx, res)
	}

	if c.FolderId != "" && c.OrgId == "" {
		scope = c.FolderId
		return analyze(client, scope, fullResourceName, ctx, res)
	}
	if c.OrgId != "" && c.FolderId == "" {
		scope = c.OrgId
		return analyze(client, scope, fullResourceName, ctx, res)
	}

	return nil
}

func analyze(client *asset.Client, scope string, fullResourceName string, ctx context.Context, res chan<- any) error {
	req := &assetpb.AnalyzeIamPolicyRequest{
		AnalysisQuery: &assetpb.IamPolicyAnalysisQuery{
			Scope: scope,
			ResourceSelector: &assetpb.IamPolicyAnalysisQuery_ResourceSelector{
				FullResourceName: fullResourceName,
			},
			Options: &assetpb.IamPolicyAnalysisQuery_Options{
				ExpandGroups:     true,
				OutputGroupEdges: true,
			},
		},
	}
	op, err := client.AnalyzeIamPolicy(ctx, req)
	if err != nil {
		return err
	}
	for _, analysis := range op.MainAnalysis.AnalysisResults {
		res <- analysis
	}
	return nil
}
