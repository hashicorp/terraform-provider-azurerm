// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandArmWorkspaceDataLakeStorageAccountDetails(storageDataLakeGen2FilesystemId string) *synapse.DataLakeStorageAccountDetails {
	uri, _ := url.Parse(storageDataLakeGen2FilesystemId)
	return &synapse.DataLakeStorageAccountDetails{
		AccountURL: pointer.To(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)), // https://storageaccountname.dfs.core.windows.net/filesystemname -> https://storageaccountname.dfs.core.windows.net
		Filesystem: pointer.To(strings.TrimPrefix(uri.Path, "/")),            // https://storageaccountname.dfs.core.windows.net/filesystemname -> filesystemname
	}
}

func expandWorkspaceRepositoryConfiguration(d *pluginsdk.ResourceData) *synapse.WorkspaceRepositoryConfiguration {
	if azdoList, ok := d.GetOk("azure_devops_repo"); ok {
		azdo := azdoList.([]any)[0].(map[string]any)
		config := synapse.WorkspaceRepositoryConfiguration{
			Type:                pointer.To(workspaceVSTSConfiguration),
			AccountName:         pointer.To(azdo["account_name"].(string)),
			CollaborationBranch: pointer.To(azdo["branch_name"].(string)),
			LastCommitID:        pointer.To(azdo["last_commit_id"].(string)),
			ProjectName:         pointer.To(azdo["project_name"].(string)),
			RepositoryName:      pointer.To(azdo["repository_name"].(string)),
			RootFolder:          pointer.To(azdo["root_folder"].(string)),
		}
		if azdoTenantId := uuid.FromStringOrNil(azdo["tenant_id"].(string)); azdoTenantId != uuid.Nil {
			config.TenantID = &azdoTenantId
		}
		return &config
	}

	if githubList, ok := d.GetOk("github_repo"); ok {
		github := githubList.([]any)[0].(map[string]any)
		return &synapse.WorkspaceRepositoryConfiguration{
			Type:                pointer.To(workspaceGitHubConfiguration),
			AccountName:         pointer.To(github["account_name"].(string)),
			CollaborationBranch: pointer.To(github["branch_name"].(string)),
			HostName:            pointer.To(github["git_url"].(string)),
			LastCommitID:        pointer.To(github["last_commit_id"].(string)),
			RepositoryName:      pointer.To(github["repository_name"].(string)),
			RootFolder:          pointer.To(github["root_folder"].(string)),
		}
	}

	// API won't clear an existing repository config with nil
	return &synapse.WorkspaceRepositoryConfiguration{}
}

func expandIdentityControlSQLSettings(enabled bool) *synapse.ManagedIdentitySQLControlSettingsModel {
	var desiredState synapse.DesiredState
	if enabled {
		desiredState = synapse.DesiredStateEnabled
	} else {
		desiredState = synapse.DesiredStateDisabled
	}

	return &synapse.ManagedIdentitySQLControlSettingsModel{
		ManagedIdentitySQLControlSettingsModelProperties: &synapse.ManagedIdentitySQLControlSettingsModelProperties{
			GrantSQLControlToManagedIdentity: &synapse.ManagedIdentitySQLControlSettingsModelPropertiesGrantSQLControlToManagedIdentity{
				DesiredState: desiredState,
			},
		},
	}
}

func expandEncryptionDetails(d *pluginsdk.ResourceData) *synapse.EncryptionDetails {
	if cmkList, ok := d.GetOk("customer_managed_key"); ok {
		cmk := cmkList.([]any)[0].(map[string]any)

		encryptionDetails := &synapse.EncryptionDetails{
			Cmk: &synapse.CustomerManagedKeyDetails{
				Key: &synapse.WorkspaceKeyDetails{
					Name:        pointer.To(cmk["key_name"].(string)),
					KeyVaultURL: pointer.To(cmk["key_versionless_id"].(string)),
				},
			},
		}

		if v, ok := cmk["user_assigned_identity_id"]; ok && v.(string) != "" {
			encryptionDetails.Cmk.KekIdentity = &synapse.KekIdentityProperties{
				UserAssignedIdentity:      pointer.To(v.(string)),
				UseSystemAssignedIdentity: false,
			}
		}

		return encryptionDetails
	}

	return nil
}

func expandIdentity(input []any) (*synapse.ManagedIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := synapse.ManagedIdentity{
		Type: synapse.ResourceIdentityType(string(expanded.Type)),
	}

	if len(expanded.IdentityIds) > 0 {
		userAssignedIdentities := make(map[string]*synapse.UserAssignedManagedIdentity)
		for id := range expanded.IdentityIds {
			userAssignedIdentities[id] = &synapse.UserAssignedManagedIdentity{}
		}
		out.UserAssignedIdentities = userAssignedIdentities
	}

	return &out, nil
}
