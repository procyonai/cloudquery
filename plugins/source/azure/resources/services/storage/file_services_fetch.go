package storage

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
)

func fetchFileServices(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armstorage.NewFileServicesClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	item := parent.Item.(*armstorage.Account)
	group, err := client.ParseResourceGroup(*item.ID)
	if err != nil {
		return err
	}
	fileSvc, err := svc.List(ctx, group, *item.Name, nil)
	if err != nil {
		return err
	}
	res <- fileSvc.Value
	return nil
}
