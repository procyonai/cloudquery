package iam

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugins/source/gcp/client"

	"google.golang.org/api/cloudresourcemanager/v3"
)

type gcpBinding struct {
	Member string
	Role   string
}

func fetchRoleBinding(ctx context.Context, meta schema.ClientMeta, r *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	iamClient, err := cloudresourcemanager.NewService(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	// List the users at the project level
	policy, err := iamClient.Projects.GetIamPolicy(fmt.Sprintf("projects/%s", c.ProjectId),
		&cloudresourcemanager.GetIamPolicyRequest{}).Do()
	if err != nil {
		return err
	}
	var bindings []gcpBinding
	for _, binding := range policy.Bindings {
		for _, member := range binding.Members {
			bindings = append(bindings, gcpBinding{
				Member: member,
				Role:   binding.Role,
			})
		}
	}
	res <- bindings

	return nil
}

func fetchOrgRoleBinding(ctx context.Context, meta schema.ClientMeta, r *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	iamClient, err := cloudresourcemanager.NewService(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	if c.OrgId != "" {
		// List the users at the project level
		policy, err := iamClient.Organizations.GetIamPolicy(c.OrgId,
			&cloudresourcemanager.GetIamPolicyRequest{}).Do()
		if err != nil {
			return err
		}
		var bindings []gcpBinding
		for _, binding := range policy.Bindings {
			for _, member := range binding.Members {
				bindings = append(bindings, gcpBinding{
					Member: member,
					Role:   binding.Role,
				})
			}
		}
		res <- bindings
	}

	return nil
}

func fetchFolderRoleBinding(ctx context.Context, meta schema.ClientMeta, r *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	iamClient, err := cloudresourcemanager.NewService(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	if c.FolderId != "" {
		// List the users at the project level
		policy, err := iamClient.Folders.GetIamPolicy(c.FolderId,
			&cloudresourcemanager.GetIamPolicyRequest{}).Do()
		if err != nil {
			return err
		}
		var bindings []gcpBinding
		for _, binding := range policy.Bindings {
			for _, member := range binding.Members {
				bindings = append(bindings, gcpBinding{
					Member: member,
					Role:   binding.Role,
				})
			}
		}
		res <- bindings
	}

	return nil
}
