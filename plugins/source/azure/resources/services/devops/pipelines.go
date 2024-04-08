package devops

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devops/armdevops"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func Pipelines() *schema.Table {
	return &schema.Table{
		Name:        "azure_devops_pipelines",
		Resolver:    fetchPipelines,
		Description: "https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devops/armdevops@v0.5.0#Pipeline",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_devops_pipelines", client.Namespacemicrosoft_devops),
		Transform:   transformers.TransformWithStruct(&armdevops.Pipeline{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{client.SubscriptionID},
	}
}

func fetchPipelines(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armdevops.NewPipelinesClient(cl.SubscriptionId, cl.Creds, cl.Options)
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
