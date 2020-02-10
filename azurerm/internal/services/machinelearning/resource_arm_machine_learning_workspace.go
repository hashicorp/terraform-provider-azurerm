package machinelearning

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2019-11-01/machinelearningservices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMachineLearningWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMachineLearningWorkspaceCreateUpdate,
		Read:   resourceArmMachineLearningWorkspaceRead,
		Update: resourceArmMachineLearningWorkspaceCreateUpdate,
		Delete: resourceArmMachineLearningWorkspaceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"key_vault": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"application_insights": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"container_registry": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"storage_account": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tier": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(machinelearningservices.SystemAssigned),
								"systemAssigned",
							}, true),
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
		},
	}
}

func resourceArmMachineLearningWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	existing, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existing Azure Machine Learning Workspace %q (Resource Group %q): %s", name, resGroup, err)
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_machine_learning_workspace", *existing.ID)
		}
	}

	workspace := machinelearningservices.Workspace{
		Name:                &name,
		Location:            &location,
		Tags:                tags.Expand(t),
		Sku:                 expandArmMachineLearningSku(d),
		Identity:            expandArmMachineLearningWorkspaceIdentity(d),
		WorkspaceProperties: expandWorkspaceProperties(d),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, workspace)
	if err != nil {
		return fmt.Errorf("Error during Azure Machine Learning Workspace creation %q in resource group (%q): %+v", name, resGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Azure Machine Learning Workspace %q (Resource Group %q) to be created..", name, resGroup)
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for create/update of Azure Machine Learning Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}
	log.Printf("[DEBUG] Azure Machine Learning Workspace %q (Resource Group %q) was created", name, resGroup)

	log.Printf("[DEBUG] Retrieving Azure Machine Learning Workspace %q (Resource Group %q)..", name, resGroup)
	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Machine Learning Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Error reading Azure Machine Learning Workspace %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmMachineLearningWorkspaceRead(d, meta)
}

func resourceArmMachineLearningWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["machineLearningServices"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WorkspaceProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("storage_account", props.StorageAccount)
		d.Set("discovery_url", props.DiscoveryURL)
		d.Set("container_registry", props.ContainerRegistry)
		d.Set("application_insights", props.ApplicationInsights)
		d.Set("key_vault", props.KeyVault)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMachineLearningWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["machineLearningServices"]

	_, err = client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandWorkspaceProperties(d *schema.ResourceData) *machinelearningservices.WorkspaceProperties {
	properties := machinelearningservices.WorkspaceProperties{}

	if description, hasDescription := d.GetOk("description"); hasDescription {
		de := description.(string)
		properties.Description = &de
	}

	if friendlyName, hasFriendlyName := d.GetOk("friendly_name"); hasFriendlyName {
		f := friendlyName.(string)
		properties.FriendlyName = &f
	}

	if storage, hasStorage := d.GetOk("storage_account"); hasStorage {
		s := storage.(string)
		properties.StorageAccount = &s
	}

	if containerRegistry, hasContainerRegistry := d.GetOk("container_registry"); hasContainerRegistry {
		c := containerRegistry.(string)
		properties.ContainerRegistry = &c
	}

	if applicationInsights, hasApplicationInsights := d.GetOk("application_insights"); hasApplicationInsights {
		i := applicationInsights.(string)
		properties.ApplicationInsights = &i
	}

	if keyVault, hasKeyVault := d.GetOk("key_vault"); hasKeyVault {
		k := keyVault.(string)
		properties.KeyVault = &k
	}

	return &properties
}

func expandArmMachineLearningWorkspaceIdentity(d *schema.ResourceData) *machinelearningservices.Identity {
	identities := d.Get("identity").([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := machinelearningservices.ResourceIdentityType(identity["type"].(string))
	return &machinelearningservices.Identity{
		Type: identityType,
	}
}

func expandArmMachineLearningSku(d *schema.ResourceData) *machinelearningservices.Sku {
	s := d.Get("sku").([]interface{})
	sku := s[0].(map[string]interface{})
	skuName := sku["name"].(string)
	skuTier := sku["tier"].(string)
	return &machinelearningservices.Sku{
		Name: &skuName,
		Tier: &skuTier,
	}
}
