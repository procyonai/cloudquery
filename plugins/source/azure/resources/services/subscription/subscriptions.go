package subscription

import (
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
)

func Subscriptions() *schema.Table {
	return &schema.Table{
		Name:        "azure_subscription_subscriptions",
		Resolver:    fetchSubscriptions,
		Description: "https://learn.microsoft.com/en-us/rest/api/resources/subscriptions/list?tabs=HTTP#subscription",
		Multiplex:   client.SingleSubscriptionMultiplex,
		Transform:   client.WithColumnValueTransformer("type", "Microsoft.Subscriptions"),
		Columns:     schema.ColumnList{},
		Relations: []*schema.Table{
			locations(),
		},
	}
}
