package storage

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func File_shares() *schema.Table {
	return &schema.Table{
		Name:        "azure_storage_file_shares",
		Resolver:    fetchFileShares,
		Description: "https://learn.microsoft.com/en-us/rest/api/storagerp/file-shares/list?tabs=HTTP#fileshareproperties",
		Transform:   transformers.TransformWithStruct(&armstorage.FileShareItem{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{},
	}
}
