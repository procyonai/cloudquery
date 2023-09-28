package compute

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
	iamadmin "cloud.google.com/go/iam/admin/apiv1"
	iampb "cloud.google.com/go/iam/admin/apiv1/adminpb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func InstancesGrantableRoles() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_compute_instance_grantable_roles",
		Description: `https://cloud.google.com/compute/docs/reference/rest/v1/instances/getIamPolicy`,
		Resolver:    fetchInstanceGrantableRoles,
		Multiplex:   client.ProjectMultiplexEnabledServices("iam.googleapis.com"),
		Transform:   transformers.TransformWithStruct(&iampb.Role{}),
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

func fetchInstanceGrantableRoles(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	p := parent.Item.(*computepb.Instance)
	nextPageToken := ""
	iamClient, err := iamadmin.NewIamClient(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	// Get the instance
	zoneSlice := strings.Split(*p.Zone, "/")
	zone := zoneSlice[len(zoneSlice)-1]

	fullResourceName := fmt.Sprintf("//compute.googleapis.com/projects/%s/zones/%s/instances/%s", c.ProjectId, zone, *p.Name)
	for {
		req := &iampb.QueryGrantableRolesRequest{
			FullResourceName: fullResourceName,
			PageSize:         1000,
			PageToken:        nextPageToken,
		}
		resp, err := iamClient.QueryGrantableRoles(ctx, req)
		if err != nil {
			return err
		}

		res <- resp.Roles

		if resp.NextPageToken == "" {
			break
		}
		nextPageToken = resp.NextPageToken
	}

	return nil
}
