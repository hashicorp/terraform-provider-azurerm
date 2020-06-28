package synapse

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/2019-06-01-preview/synapse"
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

func resourceArmSynapseWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSynapseWorkspaceCreate,
		Read:   resourceArmSynapseWorkspaceRead,
		Update: resourceArmSynapseWorkspaceUpdate,
		Delete: resourceArmSynapseWorkspaceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SynapseWorkspaceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SynapseWorkspaceName,
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
				Optional:  true,
				Sensitive: true,
			},

			"managed_virtual_network_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(synapse.ResourceIdentityTypeSystemAssigned),
							}, false),
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

			"connectivity_endpoints": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"managed_resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmSynapseWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
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
		Identity: expandArmWorkspaceManagedIdentity(d.Get("identity").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, workspaceInfo)
	if err != nil {
		return fmt.Errorf("creating Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Synapse Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse Workspace %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmSynapseWorkspaceRead(d, meta)
}

func resourceArmSynapseWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseWorkspaceID(d.Id())
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
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if err := d.Set("identity", flattenArmWorkspaceManagedIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting identity: %+v", err)
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
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSynapseWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

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
	return resourceArmSynapseWorkspaceRead(d, meta)
}

func resourceArmSynapseWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Synapse Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return waitForSynapseWorkspaceToBeDeleted(ctx, client, id.ResourceGroup, id.Name, d)
}

func waitForSynapseWorkspaceToBeDeleted(ctx context.Context, client *synapse.WorkspacesClient, resourceGroup, name string, d *schema.ResourceData) error {
	log.Printf("[DEBUG] Waiting for Synapse Workspace %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: synapseWorkspaceStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Synapse Workspace %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func synapseWorkspaceStatusCodeRefreshFunc(ctx context.Context, client *synapse.WorkspacesClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("retrieving Synapse Workspace %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of the Synapse Workspace %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandArmWorkspaceDataLakeStorageAccountDetails(storageDataLakeGen2FilesystemId string) *synapse.DataLakeStorageAccountDetails {
	uri, _ := url.Parse(storageDataLakeGen2FilesystemId)
	return &synapse.DataLakeStorageAccountDetails{
		AccountURL: utils.String(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)), // https://storageaccountname.dfs.core.windows.net/filesystemname -> https://storageaccountname.dfs.core.windows.net
		Filesystem: utils.String(uri.Path[1:]),                                 // https://storageaccountname.dfs.core.windows.net/filesystemname -> filesystemname
	}
}

func expandArmWorkspaceManagedIdentity(input []interface{}) *synapse.ManagedIdentity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &synapse.ManagedIdentity{
		Type: synapse.ResourceIdentityType(v["type"].(string)),
	}
}

func flattenArmWorkspaceManagedIdentity(input *synapse.ManagedIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t synapse.ResourceIdentityType
	if input.Type != "" {
		t = input.Type
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
			"type":         t,
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
