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
	if c.OrgId != "" {
		scope = c.OrgId
	}
	if c.FolderId != "" {
		scope = c.FolderId
	}

	fullResourceName := fmt.Sprintf("//compute.googleapis.com/projects/%s/zones/%s/instances/%s", c.ProjectId, zone, *p.Name)
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

// func fetchInstanceIamPolicy(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
// 	c := meta.(*client.Client)
// 	p := parent.Item.(*computepb.Instance)
// 	computeClient, err := compute.NewService(ctx, c.ClientOptions...)
// 	if err != nil {
// 		return err
// 	}
// 	// Get the instance
// 	zoneSlice := strings.Split(*p.Zone, "/")
// 	zone := zoneSlice[len(zoneSlice)-1]
// 	policy, err := computeClient.Instances.GetIamPolicy(c.ProjectId, zone, *p.Name).Do()
// 	if err != nil {
// 		log.Fatalf("Failed to retrieve instance iam policy: %v", err)
// 		return err
// 	}
// 	var bindings []gcpBinding
// 	for _, binding := range policy.Bindings {
// 		bindings = append(bindings, gcpBinding{
// 			Members: binding.Members,
// 			Role:    binding.Role,
// 		})
// 	}
// 	res <- bindings

// 	return nil
// }
