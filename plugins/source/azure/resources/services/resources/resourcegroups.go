package resources

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func ResourceGroups() *schema.Table {
	return &schema.Table{
		Name:        "azure_resources_resourcegroups",
		Resolver:    fetchResourceGroups,
		Description: "https://learn.microsoft.com/en-us/rest/api/resources/resource-groups/list#resourcegroup",
		Multiplex:   client.SubscriptionMultiplex,
		Transform:   transformers.TransformWithStruct(&armresources.GenericResourceExpanded{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{},
	}
}

func fetchResourceGroups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	resourceGroupObjects := cl.ResourceGroupObjs()
	for _, resourceGroupObjs := range resourceGroupObjects {
		for _, resourceGroupObj := range resourceGroupObjs {
			res <- resourceGroupObj
		}
	}
	return nil
}
