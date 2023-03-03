package iam

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
)

func fetchIamGroups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	var config iam.ListGroupsInput
	svc := meta.(*client.Client).Services().Iam
	for {
		response, err := svc.ListGroups(ctx, &config)
		if err != nil {
			return err
		}
		res <- response.Groups
		if aws.ToString(response.Marker) == "" {
			break
		}
		config.Marker = response.Marker
	}
	return nil
}
func resolveIamGroupPolicies(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(types.Group)
	svc := meta.(*client.Client).Services().Iam
	config := iam.ListAttachedGroupPoliciesInput{
		GroupName: r.GroupName,
	}
	response, err := svc.ListAttachedGroupPolicies(ctx, &config)
	if err != nil {
		return err
	}
	policyMap := map[string]*string{}
	for _, p := range response.AttachedPolicies {
		policyMap[*p.PolicyArn] = p.PolicyName
	}
	return resource.Set(c.Name, policyMap)
}

func resolveIamGroupPoliciesDocuments(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(types.Group)
	svc := meta.(*client.Client).Services().Iam
	config := iam.ListAttachedGroupPoliciesInput{
		GroupName: r.GroupName,
	}
	response, err := svc.ListAttachedGroupPolicies(ctx, &config)
	if err != nil {
		return err
	}
	policyMap := make(map[string]*map[string]any)
	for _, p := range response.AttachedPolicies {
		resp, err := svc.GetPolicy(ctx, &iam.GetPolicyInput{PolicyArn: p.PolicyArn})
		if err != nil {
			continue
		}
		policyResult, err := svc.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{PolicyArn: p.PolicyArn, VersionId: resp.Policy.DefaultVersionId})
		if err != nil {
			continue
		}

		decodedDocument, err := url.QueryUnescape(*policyResult.PolicyVersion.Document)
		if err != nil {
			continue
		}

		var document map[string]any
		err = json.Unmarshal([]byte(decodedDocument), &document)
		if err != nil {
			continue
		}
		policyMap[*p.PolicyArn] = &document
	}

	return resource.Set(c.Name, policyMap)
}
