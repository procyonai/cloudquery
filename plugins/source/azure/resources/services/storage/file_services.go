package storage

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func File_services() *schema.Table {
	return &schema.Table{
		Name:        "azure_storage_file_services",
		Resolver:    fetchFileServices,
		Description: "https://learn.microsoft.com/en-us/rest/api/storagerp/file-services/list?tabs=HTTP#fileserviceproperties",
		Transform:   transformers.TransformWithStruct(&armstorage.FileServiceProperties{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{},
	}
}
