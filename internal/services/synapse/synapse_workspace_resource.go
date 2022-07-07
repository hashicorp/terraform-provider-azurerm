package synapse

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2021-03-01/synapse"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	workspaceVSTSConfiguration   = "WorkspaceVSTSConfiguration"
	workspaceGitHubConfiguration = "WorkspaceGitHubConfiguration"
)

func resourceSynapseWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseWorkspaceCreate,
		Read:   resourceSynapseWorkspaceRead,
		Update: resourceSynapseWorkspaceUpdate,
		Delete: resourceSynapseWorkspaceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"storage_data_lake_gen2_filesystem_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sql_administrator_login": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlAdministratorLoginName,
			},

			"sql_administrator_login_password": {
				Type:      pluginsdk.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"linking_allowed_for_aad_tenant_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"compute_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"data_exfiltration_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"managed_virtual_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"aad_admin": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConfigMode:    pluginsdk.SchemaConfigModeAttr,
				ConflictsWith: []string{"customer_managed_key"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"login": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"object_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"sql_aad_admin": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConfigMode:    pluginsdk.SchemaConfigModeAttr,
				ConflictsWith: []string{"customer_managed_key"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"login": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"object_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"connectivity_endpoints": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"identity": commonschema.SystemAssignedIdentityRequired(),

			"managed_resource_group_name": commonschema.ResourceGroupNameOptionalComputed(),

			"azure_devops_repo": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"github_repo"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"branch_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"last_commit_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"project_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"repository_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"root_folder": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.RepoRootFolder(),
						},
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"github_repo": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"azure_devops_repo"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"branch_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"git_url": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
						"last_commit_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"repository_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"root_folder": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.RepoRootFolder(),
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"purview_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: account.ValidateAccountID,
			},

			"sql_identity_control_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"customer_managed_key": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"aad_admin", "sql_aad_admin"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_versionless_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.VersionlessNestedItemId,
						},

						// Default to cmk to ensure backwards compatibility with previous version that hardcoded the key name to cmk
						"key_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "cmk",
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSynapseWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	aadAdminClient := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	sqlAdminClient := meta.(*clients.Client).Synapse.WorkspaceSQLAadAdminsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_synapse_workspace", id.ID())
	}

	managedVirtualNetwork := ""
	if d.Get("managed_virtual_network_enabled").(bool) {
		managedVirtualNetwork = "default"
	}

	publicNetworkAccess := synapse.WorkspacePublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = synapse.WorkspacePublicNetworkAccessDisabled
	}

	workspaceInfo := synapse.Workspace{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		WorkspaceProperties: &synapse.WorkspaceProperties{
			DefaultDataLakeStorage:           expandArmWorkspaceDataLakeStorageAccountDetails(d.Get("storage_data_lake_gen2_filesystem_id").(string)),
			ManagedVirtualNetwork:            utils.String(managedVirtualNetwork),
			PublicNetworkAccess:              publicNetworkAccess,
			SQLAdministratorLogin:            utils.String(d.Get("sql_administrator_login").(string)),
			SQLAdministratorLoginPassword:    utils.String(d.Get("sql_administrator_login_password").(string)),
			ManagedResourceGroupName:         utils.String(d.Get("managed_resource_group_name").(string)),
			WorkspaceRepositoryConfiguration: expandWorkspaceRepositoryConfiguration(d),
			Encryption:                       expandEncryptionDetails(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	expandedIdentity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	workspaceInfo.Identity = expandedIdentity

	if purviewId, ok := d.GetOk("purview_id"); ok {
		workspaceInfo.WorkspaceProperties.PurviewConfiguration = &synapse.PurviewConfiguration{
			PurviewResourceID: utils.String(purviewId.(string)),
		}
	}

	if computeSubnetId, ok := d.GetOk("compute_subnet_id"); ok {
		workspaceInfo.WorkspaceProperties.VirtualNetworkProfile = &synapse.VirtualNetworkProfile{
			ComputeSubnetID: utils.String(computeSubnetId.(string)),
		}
	}

	dataExfiltrationProtectionEnabled := d.Get("data_exfiltration_protection_enabled").(bool)

	if dataExfiltrationProtectionEnabled {
		workspaceInfo.ManagedVirtualNetworkSettings = &synapse.ManagedVirtualNetworkSettings{
			PreventDataExfiltration: utils.Bool(dataExfiltrationProtectionEnabled),
		}
	}

	if allowedLinkingTenantIds, ok := d.GetOk("linking_allowed_for_aad_tenant_ids"); ok {
		if workspaceInfo.ManagedVirtualNetworkSettings == nil {
			workspaceInfo.ManagedVirtualNetworkSettings = &synapse.ManagedVirtualNetworkSettings{}
		}
		workspaceInfo.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking = utils.ExpandStringSlice(allowedLinkingTenantIds.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, workspaceInfo)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	if err := waitSynapseWorkspaceCMKState(ctx, client, &id); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}

	aadAdmin := expandArmWorkspaceAadAdminInfo(d.Get("aad_admin").([]interface{}))
	if aadAdmin != nil {
		future, err := aadAdminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *aadAdmin)
		if err != nil {
			return fmt.Errorf("configuring AzureAD Admin for %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for configuration of AzureAD Admin for %s: %+v", id, err)
		}
	}

	sqlAdmin := expandArmWorkspaceAadAdminInfo(d.Get("sql_aad_admin").([]interface{}))
	if sqlAdmin != nil {
		future, err := sqlAdminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *sqlAdmin)
		if err != nil {
			return fmt.Errorf("configuring Sql Admin for %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for configuration of Sql Admin for %s: %+v", id, err)
		}
	}

	sqlControlSettings := expandIdentityControlSQLSettings(d.Get("sql_identity_control_enabled").(bool))
	future2, err := identitySQLControlClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *sqlControlSettings)
	if err != nil {
		return fmt.Errorf("configuring Sql Identity Control for %s: %+v", id, err)
	}
	if err = future2.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for configuration of Sql Identity Control for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSynapseWorkspaceRead(d, meta)
}

func resourceSynapseWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	aadAdminClient := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	sqlAdminClient := meta.(*clients.Client).Synapse.WorkspaceSQLAadAdminsClient
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	aadAdmin, err := aadAdminClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// NOTE: AAD Admin isn't supported for a Workspace created from a Dedicated SQL Pool / SQL DataWarehouse and returns a Conflict
		if !utils.ResponseWasNotFound(aadAdmin.Response) && !utils.ResponseWasConflict(aadAdmin.Response) {
			return fmt.Errorf("retrieving AzureAD Admin for %s: %+v", *id, err)
		}
	}
	sqlAdmin, err := sqlAdminClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// NOTE: SQL Admin isn't supported for a Workspace created from a Dedicated SQL Pool / SQL DataWarehouse and returns a Conflict
		if !utils.ResponseWasNotFound(sqlAdmin.Response) && !utils.ResponseWasConflict(sqlAdmin.Response) {
			return fmt.Errorf("retrieving Sql Admin for %s: %+v", *id, err)
		}
	}

	sqlControlSettings, err := identitySQLControlClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Sql Identity Control for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("identity", flattenIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.WorkspaceProperties; props != nil {
		managedVirtualNetworkEnabled := false
		if props.ManagedVirtualNetwork != nil && strings.EqualFold(*props.ManagedVirtualNetwork, "default") {
			managedVirtualNetworkEnabled = true
			if props.ManagedVirtualNetworkSettings != nil {
				d.Set("data_exfiltration_protection_enabled", props.ManagedVirtualNetworkSettings.PreventDataExfiltration)
				d.Set("linking_allowed_for_aad_tenant_ids", utils.FlattenStringSlice(props.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking))
			}
		}
		d.Set("managed_virtual_network_enabled", managedVirtualNetworkEnabled)
		d.Set("storage_data_lake_gen2_filesystem_id", flattenArmWorkspaceDataLakeStorageAccountDetails(props.DefaultDataLakeStorage))
		d.Set("sql_administrator_login", props.SQLAdministratorLogin)
		d.Set("managed_resource_group_name", props.ManagedResourceGroupName)
		d.Set("connectivity_endpoints", utils.FlattenMapStringPtrString(props.ConnectivityEndpoints))
		d.Set("public_network_access_enabled", resp.PublicNetworkAccess == synapse.WorkspacePublicNetworkAccessEnabled)
		cmk := flattenEncryptionDetails(props.Encryption)
		if err := d.Set("customer_managed_key", cmk); err != nil {
			return fmt.Errorf("setting `customer_managed_key`: %+v", err)
		}

		repoType, repo := flattenWorkspaceRepositoryConfiguration(props.WorkspaceRepositoryConfiguration)
		if repoType == workspaceVSTSConfiguration {
			if err := d.Set("azure_devops_repo", repo); err != nil {
				return fmt.Errorf("setting `azure_devops_repo`: %+v", err)
			}
		} else if repoType == workspaceGitHubConfiguration {
			if err := d.Set("github_repo", repo); err != nil {
				return fmt.Errorf("setting `github_repo`: %+v", err)
			}
		}

		if props.VirtualNetworkProfile != nil {
			d.Set("compute_subnet_id", props.VirtualNetworkProfile.ComputeSubnetID)
		}
		if props.PurviewConfiguration != nil {
			d.Set("purview_id", props.PurviewConfiguration.PurviewResourceID)
		}
	}
	if err := d.Set("aad_admin", flattenArmWorkspaceAadAdmin(aadAdmin.AadAdminProperties)); err != nil {
		return fmt.Errorf("setting `aad_admin`: %+v", err)
	}
	if err := d.Set("sql_aad_admin", flattenArmWorkspaceAadAdmin(sqlAdmin.AadAdminProperties)); err != nil {
		return fmt.Errorf("setting `sql_aad_admin`: %+v", err)
	}
	if err := d.Set("sql_identity_control_enabled", flattenIdentityControlSQLSettings(sqlControlSettings)); err != nil {
		return fmt.Errorf("setting `sql_identity_control_enabled`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSynapseWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	aadAdminClient := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	sqlAdminClient := meta.(*clients.Client).Synapse.WorkspaceSQLAadAdminsClient
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChanges("tags", "sql_administrator_login_password", "github_repo", "azure_devops_repo", "customer_managed_key", "public_network_access_enabled") {
		publicNetworkAccess := synapse.WorkspacePublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = synapse.WorkspacePublicNetworkAccessDisabled
		}
		workspacePatchInfo := synapse.WorkspacePatchInfo{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
			WorkspacePatchProperties: &synapse.WorkspacePatchProperties{
				SQLAdministratorLoginPassword:    utils.String(d.Get("sql_administrator_login_password").(string)),
				WorkspaceRepositoryConfiguration: expandWorkspaceRepositoryConfiguration(d),
				Encryption:                       expandEncryptionDetails(d),
				PublicNetworkAccess:              publicNetworkAccess,
			},
		}

		if allowedLinkingTenantIds, ok := d.GetOk("linking_allowed_for_aad_tenant_ids"); ok {
			if workspacePatchInfo.ManagedVirtualNetworkSettings == nil {
				workspacePatchInfo.ManagedVirtualNetworkSettings = &synapse.ManagedVirtualNetworkSettings{}
			}
			workspacePatchInfo.ManagedVirtualNetworkSettings.AllowedAadTenantIdsForLinking = utils.ExpandStringSlice(allowedLinkingTenantIds.([]interface{}))
		}

		if purviewId, ok := d.GetOk("purview_id"); ok {
			workspacePatchInfo.PurviewConfiguration = &synapse.PurviewConfiguration{
				PurviewResourceID: utils.String(purviewId.(string)),
			}
		}

		future, err := client.Update(ctx, id.ResourceGroup, id.Name, workspacePatchInfo)
		if err != nil {
			return fmt.Errorf("updating Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on updating future for Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := waitSynapseWorkspaceCMKState(ctx, client, id); err != nil {
			return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
		}
	}

	if d.HasChange("aad_admin") {
		aadAdmin := expandArmWorkspaceAadAdminInfo(d.Get("aad_admin").([]interface{}))
		if aadAdmin != nil {
			workspaceAadAdminsCreateOrUpdateFuture, err := aadAdminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *aadAdmin)
			if err != nil {
				return fmt.Errorf("updating Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err = workspaceAadAdminsCreateOrUpdateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting on updating for Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		} else {
			workspaceAadAdminsDeleteFuture, err := aadAdminClient.Delete(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("setting empty Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err = workspaceAadAdminsDeleteFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting on setting empty Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
	}

	if d.HasChange("sql_aad_admin") {
		sqlAdmin := expandArmWorkspaceAadAdminInfo(d.Get("sql_aad_admin").([]interface{}))
		if sqlAdmin != nil {
			workspaceSqlAdminsCreateOrUpdateFuture, err := sqlAdminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *sqlAdmin)
			if err != nil {
				return fmt.Errorf("updating Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err = workspaceSqlAdminsCreateOrUpdateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting on updating for Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		} else {
			workspaceSqlAdminsDeleteFuture, err := sqlAdminClient.Delete(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("setting empty Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err = workspaceSqlAdminsDeleteFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting on setting empty Synapse Workspace %q Sql Admin (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
	}

	if d.HasChange("sql_identity_control_enabled") {
		sqlControlSettings := expandIdentityControlSQLSettings(d.Get("sql_identity_control_enabled").(bool))
		future, err := identitySQLControlClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *sqlControlSettings)
		if err != nil {
			return fmt.Errorf("Updating workspace identity control for SQL pool: %+v", err)
		}
		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update workspace identity control for SQL pool of %q: %+v", id, err)
		}
	}

	return resourceSynapseWorkspaceRead(d, meta)
}

func resourceSynapseWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// sometimes the waitForCompletion rest api will return 404
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for Synapse Workspace %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func waitSynapseWorkspaceCMKState(ctx context.Context, client *synapse.WorkspacesClient, id *parse.WorkspaceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			"Updating",
			"ActivatingWorkspace",
		},
		Target: []string{
			"Succeeded",
			"Consistent",
			"AwaitingUserAction",
		},
		Refresh:                   synapseWorkspaceCMKUpdateStateRefreshFunc(ctx, client, id),
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}
	return nil
}

func synapseWorkspaceCMKUpdateStateRefreshFunc(ctx context.Context, client *synapse.WorkspacesClient, id *parse.WorkspaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Encryption != nil && res.Encryption.Cmk != nil {
			return res, *res.Encryption.Cmk.Status, nil
		}
		return res, "Succeeded", nil
	}
}

func expandArmWorkspaceDataLakeStorageAccountDetails(storageDataLakeGen2FilesystemId string) *synapse.DataLakeStorageAccountDetails {
	uri, _ := url.Parse(storageDataLakeGen2FilesystemId)
	return &synapse.DataLakeStorageAccountDetails{
		AccountURL: utils.String(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)), // https://storageaccountname.dfs.core.windows.net/filesystemname -> https://storageaccountname.dfs.core.windows.net
		Filesystem: utils.String(uri.Path[1:]),                                 // https://storageaccountname.dfs.core.windows.net/filesystemname -> filesystemname
	}
}

func expandArmWorkspaceAadAdminInfo(input []interface{}) *synapse.WorkspaceAadAdminInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &synapse.WorkspaceAadAdminInfo{
		AadAdminProperties: &synapse.AadAdminProperties{
			TenantID:          utils.String(v["tenant_id"].(string)),
			Login:             utils.String(v["login"].(string)),
			AdministratorType: utils.String("ActiveDirectory"),
			Sid:               utils.String(v["object_id"].(string)),
		},
	}
}

func expandWorkspaceRepositoryConfiguration(d *pluginsdk.ResourceData) *synapse.WorkspaceRepositoryConfiguration {
	if azdoList, ok := d.GetOk("azure_devops_repo"); ok {
		azdo := azdoList.([]interface{})[0].(map[string]interface{})
		config := synapse.WorkspaceRepositoryConfiguration{
			Type:                utils.String(workspaceVSTSConfiguration),
			AccountName:         utils.String(azdo["account_name"].(string)),
			CollaborationBranch: utils.String(azdo["branch_name"].(string)),
			LastCommitID:        utils.String(azdo["last_commit_id"].(string)),
			ProjectName:         utils.String(azdo["project_name"].(string)),
			RepositoryName:      utils.String(azdo["repository_name"].(string)),
			RootFolder:          utils.String(azdo["root_folder"].(string)),
		}
		if azdoTenantId := uuid.FromStringOrNil(azdo["tenant_id"].(string)); azdoTenantId != uuid.Nil {
			config.TenantID = &azdoTenantId
		}
		return &config
	}

	if githubList, ok := d.GetOk("github_repo"); ok {
		github := githubList.([]interface{})[0].(map[string]interface{})
		return &synapse.WorkspaceRepositoryConfiguration{
			Type:                utils.String(workspaceGitHubConfiguration),
			AccountName:         utils.String(github["account_name"].(string)),
			CollaborationBranch: utils.String(github["branch_name"].(string)),
			HostName:            utils.String(github["git_url"].(string)),
			LastCommitID:        utils.String(github["last_commit_id"].(string)),
			RepositoryName:      utils.String(github["repository_name"].(string)),
			RootFolder:          utils.String(github["root_folder"].(string)),
		}
	}

	return nil
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
		cmk := cmkList.([]interface{})[0].(map[string]interface{})
		return &synapse.EncryptionDetails{
			Cmk: &synapse.CustomerManagedKeyDetails{
				Key: &synapse.WorkspaceKeyDetails{
					Name:        utils.String(cmk["key_name"].(string)),
					KeyVaultURL: utils.String(cmk["key_versionless_id"].(string)),
				},
			},
		}
	}

	return nil
}

func flattenArmWorkspaceDataLakeStorageAccountDetails(input *synapse.DataLakeStorageAccountDetails) string {
	if input != nil && input.AccountURL != nil && input.Filesystem != nil {
		return fmt.Sprintf("%s/%s", *input.AccountURL, *input.Filesystem)
	}
	return ""
}

func flattenArmWorkspaceAadAdmin(input *synapse.AadAdminProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	var tenantId, login, sid string
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	if input.Login != nil {
		login = *input.Login
	}
	if input.Sid != nil {
		sid = *input.Sid
	}
	return []interface{}{
		map[string]interface{}{
			"tenant_id": tenantId,
			"login":     login,
			"object_id": sid,
		},
	}
}

func flattenWorkspaceRepositoryConfiguration(config *synapse.WorkspaceRepositoryConfiguration) (repoTypeResult string, result []interface{}) {
	if config == nil {
		return "", make([]interface{}, 0)
	}

	if repoType := config.Type; repoType != nil {
		repo := map[string]interface{}{}

		if *repoType == workspaceVSTSConfiguration {
			if config.ProjectName != nil {
				repo["project_name"] = *config.ProjectName
			}
			if config.TenantID != nil {
				repo["tenant_id"] = config.TenantID.String()
			}
		} else if *repoType == workspaceGitHubConfiguration {
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

		return *repoType, []interface{}{repo}
	}

	return "", make([]interface{}, 0)
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

func flattenEncryptionDetails(encryption *synapse.EncryptionDetails) []interface{} {
	if encryption != nil {
		if cmk := encryption.Cmk; cmk != nil {
			if cmk.Key != nil {
				resultMap := map[string]interface{}{}
				resultMap["key_name"] = *cmk.Key.Name
				resultMap["key_versionless_id"] = *cmk.Key.KeyVaultURL
				return []interface{}{resultMap}
			}
		}

		// if cmk := encryption.Cmk; cmk != nil {
		// 	if key := cmk.Key; key != nil {
		// 		return key.Name, key.KeyVaultURL
		// 	}
		// }
	}

	return make([]interface{}, 0)
}

func expandIdentity(input []interface{}) (*synapse.ManagedIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	out := synapse.ManagedIdentity{
		Type: synapse.ResourceIdentityType(string(expanded.Type)),
	}
	return &out, nil
}

func flattenIdentity(input *synapse.ManagedIdentity) []interface{} {
	var config *identity.SystemAssigned

	if input != nil {
		config = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}

		if input.PrincipalID != nil {
			config.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			config.TenantId = input.TenantID.String()
		}
	}

	return identity.FlattenSystemAssigned(config)
}
