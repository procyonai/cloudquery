package client

import (
	"context"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

var _ transformers.NameTransformer = ETagNameTransformer

func ETagNameTransformer(fld reflect.StructField) (string, error) {
	if fld.Name == "ETag" {
		return "etag", nil
	}
	return transformers.DefaultNameTransformer(fld)
}

func WithColumnValueTransformer(column, value string) schema.Transform {
	return func(table *schema.Table) error {
		err := transformers.TransformWithStruct(&armsubscription.Subscription{}, transformers.WithPrimaryKeys("ID"))(table)
		if err != nil {
			return err
		}
		column := schema.Column{
			Name: column,
			Type: schema.TypeString,
			Resolver: func(_ context.Context, _ schema.ClientMeta, r *schema.Resource, c schema.Column) error {
				return r.Set(c.Name, value)
			},
		}
		table.Columns = append(table.Columns, column)
		return nil
	}
}
