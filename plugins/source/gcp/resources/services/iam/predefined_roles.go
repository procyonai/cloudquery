package iam

import (
	pb "cloud.google.com/go/iam/admin/apiv1/adminpb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func PredefinedRoles() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_predefined_roles",
		Description: `https://cloud.google.com/iam/docs/reference/rest/v1/roles#Role`,
		Resolver:    fetchPredefinedRoles,
		Multiplex:   client.ProjectMultiplexEnabledServices("iam.googleapis.com"),
		Transform:   transformers.TransformWithStruct(&pb.Role{}, append(client.Options(), transformers.WithPrimaryKeys("Name"))...),
		Columns: []schema.Column{
			{
				Name:     "organization_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveOrganization,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
			{
				Name:     "project_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveProject,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
		Relations: []*schema.Table{
			RolePolicies(),
		},
	}
}
