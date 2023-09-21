package iam

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func OrgRoleBinding() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_org_role_binding",
		Description: `https://cloud.google.com/iam/docs/reference/rest/v1/roles#Role`,
		Resolver:    fetchOrgRoleBinding,
		Multiplex:   client.OrganizationMultiplexEnabledServices("iam.googleapis.com"),
		Transform:   transformers.TransformWithStruct(&gcpBinding{}),
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
	}
}
