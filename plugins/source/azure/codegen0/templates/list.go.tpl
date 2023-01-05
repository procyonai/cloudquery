// Code generated by codegen; DO NOT EDIT.
package {{.PackageName}}

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"context"
	"{{.ImportPath}}"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
)

func {{.Name | ToCamel}}() *schema.Table {
    return &schema.Table{
			Name: "azure_{{.PackageName}}_{{.Name}}",
			Resolver: fetch{{.Name | ToCamel}},
			Multiplex: client.SubscriptionMultiplexRegisteredNamespace(client.Namespace{{.NamespaceConst}}),
			Transform: transformers.TransformWithStruct(&{{.BaseImportPath}}.{{.ResponseValueStruct}}{}),
			Columns: []schema.Column{
				{
					Name:     "subscription_id",
					Type:     schema.TypeString,
					Resolver: client.ResolveAzureSubscription,
				},
				{
					Name:     "id",
					Type:     schema.TypeString,
					Resolver: schema.PathResolver("ID"),
					CreationOptions: schema.ColumnCreationOptions{
						PrimaryKey: true,
					},
				},
			},
		}
}

func fetch{{.Name | ToCamel}}(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	{{- if .NewFuncHasSubscriptionId}}
  svc, err := {{.BaseImportPath}}.{{.NewFuncName}}(cl.SubscriptionId, cl.Creds, cl.Options)
  {{- else}}
  svc, err := {{.BaseImportPath}}.{{.NewFuncName}}(cl.Creds, cl.Options)
  {{- end}}
	if err != nil {
    return err
  }
	pager := svc.{{.Pager}}(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}