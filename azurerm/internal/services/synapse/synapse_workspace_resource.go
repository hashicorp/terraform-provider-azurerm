package synapse

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/2019-06-01-preview/synapse"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSynapseWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceSynapseWorkspaceCreate,
		Read:   resourceSynapseWorkspaceRead,
		Update: resourceSynapseWorkspaceUpdate,
		Delete: resourceSynapseWorkspaceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"storage_data_lake_gen2_filesystem_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sql_administrator_login": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlAdministratorLoginName,
			},

			"sql_administrator_login_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"managed_virtual_network_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"aad_admin": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login": {
							Type:     schema.TypeString,
							Required: true,
						},

						"object_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"tenant_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"connectivity_endpoints": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"managed_resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"workspace_identity_control_for_sql_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSynapseWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	aadAdminClient := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_synapse_workspace", *existing.ID)
	}

	managedVirtualNetwork := ""
	if d.Get("managed_virtual_network_enabled").(bool) {
		managedVirtualNetwork = "default"
	}

	workspaceInfo := synapse.Workspace{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		WorkspaceProperties: &synapse.WorkspaceProperties{
			DefaultDataLakeStorage:        expandArmWorkspaceDataLakeStorageAccountDetails(d.Get("storage_data_lake_gen2_filesystem_id").(string)),
			ManagedVirtualNetwork:         utils.String(managedVirtualNetwork),
			SQLAdministratorLogin:         utils.String(d.Get("sql_administrator_login").(string)),
			SQLAdministratorLoginPassword: utils.String(d.Get("sql_administrator_login_password").(string)),
		},
		Identity: &synapse.ManagedIdentity{
			Type: synapse.ResourceIdentityTypeSystemAssigned,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, workspaceInfo)
	if err != nil {
		return fmt.Errorf("creating Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	aadAdmin := expandArmWorkspaceAadAdmin(d.Get("aad_admin").([]interface{}))
	if aadAdmin != nil {
		workspaceAadAdminsCreateOrUpdateFuture, err := aadAdminClient.CreateOrUpdate(ctx, resourceGroup, name, *aadAdmin)
		if err != nil {
			return fmt.Errorf("updating Synapse Workspace %q Sql Admin (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = workspaceAadAdminsCreateOrUpdateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on updating for Synapse Workspace %q Sql Admin (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	sqlControlSettings := expandIdentityControlSQLSettings(d.Get("workspace_identity_control_for_sql_enabled").(bool))
	_, err = identitySQLControlClient.CreateOrUpdate(ctx, resourceGroup, name, *sqlControlSettings)
	if err != nil {
		return fmt.Errorf("Granting workspace identity control for SQL pool: %+v", err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse Workspace %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceSynapseWorkspaceRead(d, meta)
}

func resourceSynapseWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	aadAdminClient := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
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
			log.Printf("[INFO] synapse %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	aadAdmin, err := aadAdminClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(aadAdmin.Response) {
			return fmt.Errorf("retrieving Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	sqlControlSettings, err := identitySQLControlClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving workspace identity control for SQL pool: %+v", err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if err := d.Set("identity", flattenArmWorkspaceManagedIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.WorkspaceProperties; props != nil {
		managedVirtualNetworkEnabled := false
		if props.ManagedVirtualNetwork != nil && strings.EqualFold(*props.ManagedVirtualNetwork, "default") {
			managedVirtualNetworkEnabled = true
		}
		d.Set("managed_virtual_network_enabled", managedVirtualNetworkEnabled)
		d.Set("storage_data_lake_gen2_filesystem_id", flattenArmWorkspaceDataLakeStorageAccountDetails(props.DefaultDataLakeStorage))
		d.Set("sql_administrator_login", props.SQLAdministratorLogin)
		d.Set("managed_resource_group_name", props.ManagedResourceGroupName)
		d.Set("connectivity_endpoints", utils.FlattenMapStringPtrString(props.ConnectivityEndpoints))
	}
	if err := d.Set("aad_admin", flattenArmWorkspaceAadAdmin(aadAdmin.AadAdminProperties)); err != nil {
		return fmt.Errorf("setting `aad_admin`: %+v", err)
	}
	if err := d.Set("workspace_identity_control_for_sql_enabled", flattenIdentityControlSQLSettings(sqlControlSettings)); err != nil {
		return fmt.Errorf("setting `workspace_identity_control_for_sql_enabled`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSynapseWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	aadAdminClient := meta.(*clients.Client).Synapse.WorkspaceAadAdminsClient
	identitySQLControlClient := meta.(*clients.Client).Synapse.WorkspaceManagedIdentitySQLControlSettingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChanges("tags", "sql_administrator_login_password") {
		workspacePatchInfo := synapse.WorkspacePatchInfo{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}
		if d.HasChange("sql_administrator_login_password") {
			workspacePatchInfo.WorkspacePatchProperties = &synapse.WorkspacePatchProperties{
				SQLAdministratorLoginPassword: utils.String(d.Get("sql_administrator_login_password").(string)),
			}
		}

		future, err := client.Update(ctx, id.ResourceGroup, id.Name, workspacePatchInfo)
		if err != nil {
			return fmt.Errorf("updating Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on updating future for Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if d.HasChange("aad_admin") {
		aadAdmin := expandArmWorkspaceAadAdmin(d.Get("aad_admin").([]interface{}))
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

	if d.HasChange("workspace_identity_control_for_sql_enabled") {
		sqlControlSettings := expandIdentityControlSQLSettings(d.Get("workspace_identity_control_for_sql_enabled").(bool))
		_, err = identitySQLControlClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *sqlControlSettings)
		if err != nil {
			return fmt.Errorf("Updating workspace identity control for SQL pool: %+v", err)
		}
	}

	return resourceSynapseWorkspaceRead(d, meta)
}

func resourceSynapseWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandArmWorkspaceDataLakeStorageAccountDetails(storageDataLakeGen2FilesystemId string) *synapse.DataLakeStorageAccountDetails {
	uri, _ := url.Parse(storageDataLakeGen2FilesystemId)
	return &synapse.DataLakeStorageAccountDetails{
		AccountURL: utils.String(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)), // https://storageaccountname.dfs.core.windows.net/filesystemname -> https://storageaccountname.dfs.core.windows.net
		Filesystem: utils.String(uri.Path[1:]),                                 // https://storageaccountname.dfs.core.windows.net/filesystemname -> filesystemname
	}
}

func expandArmWorkspaceAadAdmin(input []interface{}) *synapse.WorkspaceAadAdminInfo {
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

func flattenArmWorkspaceManagedIdentity(input *synapse.ManagedIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var principalId string
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}
	var tenantId string
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
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
