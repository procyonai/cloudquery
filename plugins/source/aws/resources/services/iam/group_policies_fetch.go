package iam

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
)

func fetchIamGroupPolicies(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	svc := c.Services().Iam
	group := parent.Item.(types.Group)
	config := iam.ListGroupPoliciesInput{
		GroupName: group.GroupName,
	}
	paginator := iam.NewListGroupPoliciesPaginator(svc, &config)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if c.IsNotFoundError(err) {
				return nil
			}
			return err
		}
		res <- page.PolicyNames
	}
	return nil
}

func getGroupPolicy(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	c := meta.(*client.Client)
	svc := c.Services().Iam
	p := resource.Item.(string)
	group := resource.Parent.Item.(types.Group)

	policyResult, err := svc.GetGroupPolicy(ctx, &iam.GetGroupPolicyInput{PolicyName: &p, GroupName: group.GroupName}, func(options *iam.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}
	resource.Item = policyResult
	return nil
}

func resolveIamGroupPolicyPolicyDocument(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	svc := meta.(*client.Client).Services().Iam
	resourceMap := resource.Item.(types.AttachedPolicy)
	policyArn := *resourceMap.PolicyArn

	resp, err := svc.GetPolicy(ctx, &iam.GetPolicyInput{PolicyArn: &policyArn})
	if err != nil {
		return err
	}
	versionId := resp.Policy.DefaultVersionId

	policyResult, err := svc.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{PolicyArn: &policyArn, VersionId: versionId})

	if err != nil {
		return err
	}

	decodedDocument, err := url.QueryUnescape(*policyResult.PolicyVersion.Document)
	if err != nil {
		return err
	}

	var document map[string]any
	err = json.Unmarshal([]byte(decodedDocument), &document)
	if err != nil {
		return err
	}
	return resource.Set(c.Name, document)
}

func resolveIamGroupInlinePolicyPolicyDocument(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(*iam.GetGroupPolicyOutput)

	decodedDocument, err := url.QueryUnescape(*r.PolicyDocument)
	if err != nil {
		return err
	}

	var document map[string]any
	err = json.Unmarshal([]byte(decodedDocument), &document)
	if err != nil {
		return err
	}
	return resource.Set(c.Name, document)
}
