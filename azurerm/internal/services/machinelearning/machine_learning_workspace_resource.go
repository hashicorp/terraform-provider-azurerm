package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO -- remove this type when issue https://github.com/Azure/azure-rest-api-specs/issues/13546 is resolved
type WorkspaceSku string

const (
	Basic WorkspaceSku = "Basic"
	// TODO -- remove Enterprise in 3.0 which has been deprecated here: https://docs.microsoft.com/en-us/azure/machine-learning/concept-workspace#what-happened-to-enterprise-edition
	Enterprise WorkspaceSku = "Enterprise"
)

func resourceMachineLearningWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMachineLearningWorkspaceCreate,
		Read:   resourceMachineLearningWorkspaceRead,
		Update: resourceMachineLearningWorkspaceUpdate,
		Delete: resourceMachineLearningWorkspaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"application_insights_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// TODO -- use the custom validation function of application insights
				ValidateFunc: azure.ValidateResourceID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.VaultID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// TODO -- use the custom validation function of storage account, when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				ValidateFunc: azure.ValidateResourceID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(machinelearningservices.SystemAssigned),
							}, false),
						},
						"principal_id": {
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

			"container_registry_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				// TODO -- use the custom validation function of container registry
				ValidateFunc: azure.ValidateResourceID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"friendly_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"high_business_impact": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Basic",
				ValidateFunc: validation.StringInSlice([]string{
					string(Basic),
					// TODO -- remove Enterprise in 3.0 which has been deprecated here: https://docs.microsoft.com/en-us/azure/machine-learning/concept-workspace#what-happened-to-enterprise-edition
					string(Enterprise),
				}, true),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMachineLearningWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing AML Workspace %q (Resource Group %q): %s", name, resGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_workspace", *existing.ID)
	}

	workspace := machinelearningservices.Workspace{
		Name:     utils.String(name),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku: &machinelearningservices.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
			Tier: utils.String(d.Get("sku_name").(string)),
		},
		Identity: expandMachineLearningWorkspaceIdentity(d.Get("identity").([]interface{})),
		WorkspaceProperties: &machinelearningservices.WorkspaceProperties{
			StorageAccount:      utils.String(d.Get("storage_account_id").(string)),
			ApplicationInsights: utils.String(d.Get("application_insights_id").(string)),
			KeyVault:            utils.String(d.Get("key_vault_id").(string)),
		},
	}

	if v, ok := d.GetOk("description"); ok {
		workspace.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		workspace.FriendlyName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("container_registry_id"); ok {
		workspace.ContainerRegistry = utils.String(v.(string))
	}

	if v, ok := d.GetOk("high_business_impact"); ok {
		workspace.HbiWorkspace = utils.Bool(v.(bool))
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, workspace)
	if err != nil {
		return fmt.Errorf("creating Machine Learning Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Machine Learning Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewWorkspaceID(subscriptionId, resGroup, name)
	d.SetId(id.ID())

	return resourceMachineLearningWorkspaceRead(d, meta)
}

func resourceMachineLearningWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.WorkspaceProperties; props != nil {
		d.Set("application_insights_id", props.ApplicationInsights)
		d.Set("storage_account_id", props.StorageAccount)
		d.Set("key_vault_id", props.KeyVault)
		d.Set("container_registry_id", props.ContainerRegistry)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("high_business_impact", props.HbiWorkspace)
	}

	if err := d.Set("identity", flattenMachineLearningWorkspaceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("flattening identity on Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMachineLearningWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	update := machinelearningservices.WorkspaceUpdateParameters{
		WorkspacePropertiesUpdateParameters: &machinelearningservices.WorkspacePropertiesUpdateParameters{},
	}

	if d.HasChange("sku_name") {
		skuName := d.Get("sku_name").(string)
		update.Sku = &machinelearningservices.Sku{
			Name: &skuName,
			Tier: &skuName,
		}
	}

	if d.HasChange("description") {
		update.WorkspacePropertiesUpdateParameters.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("friendly_name") {
		update.WorkspacePropertiesUpdateParameters.FriendlyName = utils.String(d.Get("friendly_name").(string))
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, update); err != nil {
		return fmt.Errorf("updating Machine Learning Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceMachineLearningWorkspaceRead(d, meta)
}

func resourceMachineLearningWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Machine Learning Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandMachineLearningWorkspaceIdentity(input []interface{}) *machinelearningservices.Identity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &machinelearningservices.Identity{
		Type: machinelearningservices.ResourceIdentityType(v["type"].(string)),
	}
}

func flattenMachineLearningWorkspaceIdentity(identity *machinelearningservices.Identity) []interface{} {
	if identity == nil {
		return []interface{}{}
	}

	principalID := ""
	if identity.PrincipalID != nil {
		principalID = *identity.PrincipalID
	}

	tenantID := ""
	if identity.TenantID != nil {
		tenantID = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(identity.Type),
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
