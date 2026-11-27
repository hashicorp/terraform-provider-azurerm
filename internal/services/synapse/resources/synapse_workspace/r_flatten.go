// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func flattenArmWorkspaceDataLakeStorageAccountDetails(input *synapse.DataLakeStorageAccountDetails) string {
	if input != nil && input.AccountURL != nil && input.Filesystem != nil {
		return fmt.Sprintf("%s/%s", *input.AccountURL, *input.Filesystem)
	}
	return ""
}

func flattenWorkspaceRepositoryConfiguration(config *synapse.WorkspaceRepositoryConfiguration) (repoTypeResult string, result []any) {
	if config == nil {
		return "", make([]any, 0)
	}

	if repoType := config.Type; repoType != nil {
		repo := map[string]any{}

		switch *repoType {
		case workspaceVSTSConfiguration:
			if config.ProjectName != nil {
				repo["project_name"] = *config.ProjectName
			}
			if config.TenantID != nil {
				repo["tenant_id"] = config.TenantID.String()
			}
		case workspaceGitHubConfiguration:
			if config.HostName != nil {
				repo["git_url"] = *config.HostName
			}
		}

		if config.AccountName != nil {
			repo["account_name"] = *config.AccountName
		}
		if config.CollaborationBranch != nil {
			repo["branch_name"] = *config.CollaborationBranch
		}
		if config.RepositoryName != nil {
			repo["repository_name"] = *config.RepositoryName
		}
		if config.RootFolder != nil {
			repo["root_folder"] = *config.RootFolder
		}
		if config.LastCommitID != nil {
			repo["last_commit_id"] = *config.LastCommitID
		}

		return *repoType, []any{repo}
	}

	return "", make([]any, 0)
}

func flattenIdentityControlSQLSettings(settings synapse.ManagedIdentitySQLControlSettingsModel) bool {
	if prop := settings.ManagedIdentitySQLControlSettingsModelProperties; prop != nil {
		if sqlControl := prop.GrantSQLControlToManagedIdentity; sqlControl != nil {
			if sqlControl.DesiredState == synapse.DesiredStateEnabled {
				return true
			}
		}
	}

	return false
}

func flattenEncryptionDetails(encryption *synapse.EncryptionDetails) ([]any, error) {
	output := make([]any, 0)
	if encryption == nil || encryption.Cmk == nil || encryption.Cmk.Key == nil {
		return output, nil
	}
	resultMap := map[string]any{}

	resultMap["key_name"] = pointer.From(encryption.Cmk.Key.Name)
	resultMap["key_versionless_id"] = pointer.From(encryption.Cmk.Key.KeyVaultURL)

	if encryption.Cmk.KekIdentity != nil && encryption.Cmk.KekIdentity.UserAssignedIdentity != nil {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(encryption.Cmk.KekIdentity.UserAssignedIdentity))
		if err != nil {
			return nil, err
		}
		resultMap["user_assigned_identity_id"] = parsed.ID()
	}

	return append(output, resultMap), nil
}

func flattenIdentity(input *synapse.ManagedIdentity) (any, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: nil,
		}

		if input.PrincipalID != nil {
			config.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			config.TenantId = input.TenantID.String()
		}
		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		if input.UserAssignedIdentities != nil {
			for k, v := range input.UserAssignedIdentities {
				identityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    pointer.To(v.ClientID.String()),
					PrincipalId: pointer.To(v.PrincipalID.String()),
				}
			}
		}

		config.IdentityIds = identityIds
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}
