package subscription

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func Tenants() *schema.Table {
	return &schema.Table{
		Name:        "azure_subscription_tenants",
		Resolver:    fetchTenants,
		Description: "https://learn.microsoft.com/en-us/rest/api/resources/tenants/list?tabs=HTTP#tenantiddescription",
		Multiplex:   client.SingleSubscriptionMultiplex,
		Transform: client.CombineTranformers(transformers.TransformWithStruct(&armsubscriptions.TenantIDDescription{}, transformers.WithPrimaryKeys("ID")),
			client.WithColumnValueTransformer("type", "Microsoft.Tenants")),
		Columns: schema.ColumnList{},
	}
}

func fetchTenants(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armsubscriptions.NewTenantsClient(cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	pager := svc.NewListPager(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
