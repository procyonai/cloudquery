package cosmosforpostgresql

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmosforpostgresql/armcosmosforpostgresql"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func Clusters() *schema.Table {
	return &schema.Table{
		Name:        "azure_cosmos_for_postgresql_clusters",
		Resolver:    fetchClusters,
		Description: "https://learn.microsoft.com/en-us/rest/api/postgresqlhsc/clusters/list?view=rest-postgresqlhsc-2022-11-08&tabs=HTTP#cluster",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_cosmos_for_postgresql_clusters", client.Namespacemicrosoft_documentdb),
		Transform:   transformers.TransformWithStruct(&armcosmosforpostgresql.Cluster{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{client.SubscriptionID},
	}
}

func fetchClusters(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armcosmosforpostgresql.NewClustersClient(cl.SubscriptionId, cl.Creds, cl.Options)
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
