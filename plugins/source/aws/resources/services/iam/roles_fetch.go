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

func fetchIamRoles(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	var config iam.ListRolesInput
	svc := meta.(*client.Client).Services().Iam
	paginator := iam.NewListRolesPaginator(svc, &config)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- response.Roles
	}
	return nil
}

func getRole(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	role := resource.Item.(types.Role)
	svc := meta.(*client.Client).Services().Iam
	roleDetails, err := svc.GetRole(ctx, &iam.GetRoleInput{
		RoleName: role.RoleName,
	})
	if err != nil {
		return err
	}
	resource.Item = roleDetails.Role
	return nil
}

func resolveIamRolePolicies(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(*types.Role)
	cl := meta.(*client.Client)
	svc := cl.Services().Iam
	input := iam.ListAttachedRolePoliciesInput{
		RoleName: r.RoleName,
	}
	policies := map[string]*string{}
	paginator := iam.NewListAttachedRolePoliciesPaginator(svc, &input)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			if cl.IsNotFoundError(err) {
				return nil
			}
			return err
		}
		for _, p := range response.AttachedPolicies {
			policies[*p.PolicyArn] = p.PolicyName
		}
	}
	return resource.Set("policies", policies)
}

func resolveIamRolesPoliciesDocuments(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(*types.Role)
	svc := meta.(*client.Client).Services().Iam
	config := iam.ListAttachedRolePoliciesInput{
		RoleName: r.RoleName,
	}
	response, err := svc.ListAttachedRolePolicies(ctx, &config)
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

func resolveRolesAssumeRolePolicyDocument(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(*types.Role)
	if r.AssumeRolePolicyDocument == nil {
		return nil
	}
	decodedDocument, err := url.QueryUnescape(*r.AssumeRolePolicyDocument)
	if err != nil {
		return err
	}
	var d map[string]any
	err = json.Unmarshal([]byte(decodedDocument), &d)
	if err != nil {
		return err
	}
	return resource.Set("assume_role_policy_document", d)
}

func fetchIamRolePolicies(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	svc := c.Services().Iam
	role := parent.Item.(*types.Role)
	paginator := iam.NewListRolePoliciesPaginator(svc, &iam.ListRolePoliciesInput{
		RoleName: role.RoleName,
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if c.IsNotFoundError(err) {
				return nil
			}
			return err
		}
		res <- output.PolicyNames
	}
	return nil
}

func getRolePolicy(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	c := meta.(*client.Client)
	svc := c.Services().Iam
	p := resource.Item.(string)
	role := resource.Parent.Item.(*types.Role)

	policyResult, err := svc.GetRolePolicy(ctx, &iam.GetRolePolicyInput{PolicyName: &p, RoleName: role.RoleName})
	if err != nil {
		return err
	}
	resource.Item = policyResult
	return nil
}

func resolveRolePoliciesPolicyDocument(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
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
