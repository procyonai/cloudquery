package iam

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"
)

func FolderRoleBinding() *schema.Table {
	return &schema.Table{
		Name:        "gcp_iam_folder_role_binding",
		Description: `https://cloud.google.com/iam/docs/reference/rest/v1/roles#Role`,
		Resolver:    fetchFolderRoleBinding,
		Multiplex:   client.FolderMultiplexEnabledServices("iam.googleapis.com"),
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
			{
				Name:     "folder_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveFolder,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
	}
}
