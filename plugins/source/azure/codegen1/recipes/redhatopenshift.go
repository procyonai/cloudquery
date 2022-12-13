// Code generated by codegen0; DO NOT EDIT.
package recipes

import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redhatopenshift/armredhatopenshift"

func Armredhatopenshift() []*Table {
	tables := []*Table{
		{
			NewFunc:        armredhatopenshift.NewOpenShiftClustersClient,
			PkgPath:        "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redhatopenshift/armredhatopenshift",
			URL:            "/subscriptions/{subscriptionId}/providers/Microsoft.RedHatOpenShift/openShiftClusters",
			Namespace:      "Microsoft.RedHatOpenShift",
			Multiplex:      `client.SubscriptionMultiplexRegisteredNamespace(client.NamespaceMicrosoft_RedHatOpenShift)`,
			Pager:          `NewListPager`,
			ResponseStruct: "OpenShiftClustersClientListResponse",
		},
	}
	return tables
}

func init() {
	Tables = append(Tables, Armredhatopenshift())
}