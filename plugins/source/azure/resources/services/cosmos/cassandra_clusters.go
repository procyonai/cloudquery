package cosmos

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v2"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func CassandraClusters() *schema.Table {
	return &schema.Table{
		Name:        "azure_cosmos_cassandra_clusters",
		Resolver:    fetchCassandraClusters,
		Description: "https://learn.microsoft.com/en-us/rest/api/cosmos-db-resource-provider/cassandra-clusters/list-by-subscription?view=rest-cosmos-db-resource-provider-2023-11-15&tabs=HTTP#clusterresource",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_cosmos_cassandra_clusters", client.Namespacemicrosoft_documentdb),
		Transform:   transformers.TransformWithStruct(&armcosmos.ClusterResource{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{client.SubscriptionID},
	}
}

func fetchCassandraClusters(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armcosmos.NewCassandraClustersClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	pager := svc.NewListBySubscriptionPager(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
