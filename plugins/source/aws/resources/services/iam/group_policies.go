package iam

import (
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func groupPolicies() *schema.Table {
	tableName := "aws_iam_group_policies"
	return &schema.Table{
		Name:                tableName,
		Description:         `https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroupPolicy.html`,
		Resolver:            fetchIamGroupPolicies,
		PreResourceResolver: getGroupPolicy,
		Transform:           transformers.TransformWithStruct(&iam.GetGroupPolicyOutput{}, transformers.WithPrimaryKeys("PolicyName")),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(true),
			{
				Name:     "group_arn",
				Type:     schema.TypeString,
				Resolver: schema.ParentColumnResolver("arn"),
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
			{
				Name:     "policy_document",
				Type:     schema.TypeJSON,
				Resolver: resolveIamGroupInlinePolicyPolicyDocument,
			},
		},
	}
}
