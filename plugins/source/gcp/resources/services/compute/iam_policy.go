package compute

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
	"google.golang.org/api/compute/v1"
)

type gcpBinding struct {
	Members []string
	Role    string
}

func InstancesIamPolicy() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_instance_role_binding",
		Description: `https://cloud.google.com/compute/docs/reference/rest/v1/instances/getIamPolicy`,
		Resolver:    fetchInstanceIamPolicy,
		Multiplex:   client.ProjectMultiplexEnabledServices("iam.googleapis.com"),
		Transform:   transformers.TransformWithStruct(&gcpBinding{}),
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

func fetchInstanceIamPolicy(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	p := parent.Item.(*computepb.Instance)
	computeClient, err := compute.NewService(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	// Get the instance
	zoneSlice := strings.Split(*p.Zone, "/")
	zone := zoneSlice[len(zoneSlice)-1]
	policy, err := computeClient.Instances.GetIamPolicy(c.ProjectId, zone, *p.Name).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve instance iam policy: %v", err)
		return err
	}
	var bindings []gcpBinding
	for _, binding := range policy.Bindings {
		bindings = append(bindings, gcpBinding{
			Members: binding.Members,
			Role:    binding.Role,
		})
	}
	res <- bindings

	return nil
}
