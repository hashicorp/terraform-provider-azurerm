// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSynapseWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSynapseWorkspaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"azure_devops_repo": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"last_commit_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"root_folder": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"azuread_authentication_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_versionless_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"user_assigned_identity_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"github_repo": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"git_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"last_commit_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"root_folder": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"purview_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sql_identity_control_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"linking_allowed_for_aad_tenant_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"compute_subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"data_exfiltration_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"managed_virtual_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"connectivity_endpoints": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"managed_resource_group_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_data_lake_gen2_filesystem_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sql_administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceSynapseWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	sqlControlSettings, err := identitySQLControlClient.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("retrieving Sql Identity Control for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.WorkspaceProperties; props != nil {
		managedVirtualNetworkEnabled := false
		if strings.EqualFold(pointer.From(props.ManagedVirtualNetwork), "default") {
			managedVirtualNetworkEnabled = true
			if props.ManagedVirtualNetworkSettings != nil {
				d.Set("data_exfiltration_protection_enabled", pointer.From(props.ManagedVirtualNetworkSettings.PreventDataExfiltration))
				d.Set("linking_allowed_for_aad_tenant_ids", helpers.FlattenStringSlice(props.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking))
			}
		}
		d.Set("managed_virtual_network_enabled", managedVirtualNetworkEnabled)
		d.Set("storage_data_lake_gen2_filesystem_id", flattenArmWorkspaceDataLakeStorageAccountDetails(props.DefaultDataLakeStorage))
		d.Set("sql_administrator_login", pointer.From(props.SQLAdministratorLogin))
		d.Set("managed_resource_group_name", pointer.From(props.ManagedResourceGroupName))
		d.Set("connectivity_endpoints", helpers.FlattenMapStringPtrString(props.ConnectivityEndpoints))
		d.Set("public_network_access_enabled", resp.PublicNetworkAccess == synapse.WorkspacePublicNetworkAccessEnabled)
		d.Set("azuread_authentication_only", pointer.From(props.AzureADOnlyAuthentication))
		cmk, err := flattenEncryptionDetails(props.Encryption)
		if err != nil {
			return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
		}
		if err := d.Set("customer_managed_key", cmk); err != nil {
			return fmt.Errorf("setting `customer_managed_key`: %+v", err)
		}

		repoType, repo := flattenWorkspaceRepositoryConfiguration(props.WorkspaceRepositoryConfiguration)
		switch repoType {
		case workspaceVSTSConfiguration:
			if err := d.Set("azure_devops_repo", repo); err != nil {
				return fmt.Errorf("setting `azure_devops_repo`: %+v", err)
			}
		case workspaceGitHubConfiguration:
			if err := d.Set("github_repo", repo); err != nil {
				return fmt.Errorf("setting `github_repo`: %+v", err)
			}
		}

		if props.VirtualNetworkProfile != nil {
			d.Set("compute_subnet_id", pointer.From(props.VirtualNetworkProfile.ComputeSubnetID))
		}
		if props.PurviewConfiguration != nil {
			d.Set("purview_id", pointer.From(props.PurviewConfiguration.PurviewResourceID))
		}
	}
	if err := d.Set("sql_identity_control_enabled", flattenIdentityControlSQLSettings(sqlControlSettings)); err != nil {
		return fmt.Errorf("setting `sql_identity_control_enabled`: %+v", err)
	}
	flattenIdentities, err := flattenIdentity(resp.Identity)
	if err != nil {
		return err
	}
	if err := d.Set("identity", flattenIdentities); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
