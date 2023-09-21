package iam

import (
	pb "cloud.google.com/go/iam/admin/apiv1/adminpb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func OrganizationRoles() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_organization_roles",
		Description: `https://cloud.google.com/iam/docs/reference/rest/v1/roles#Role`,
		Resolver:    fetchOrganizationRoles,
		Multiplex:   client.OrganizationMultiplexEnabledServices("iam.googleapis.com"),
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
		},
		Relations: []*schema.Table{
			RolePolicies(),
		},
	}
}
