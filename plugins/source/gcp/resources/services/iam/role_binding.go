package iam

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func RoleBinding() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_role_binding",
		Description: `https://cloud.google.com/iam/docs/reference/rest/v1/roles#Role`,
		Resolver:    fetchRoleBinding,
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
		},
	}
}
