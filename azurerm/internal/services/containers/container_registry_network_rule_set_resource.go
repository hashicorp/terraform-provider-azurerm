package containers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2018-09-01/containerregistry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerRegistryNetworkRuleset() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerRegistryNetworkRulesetCreateUpdate,
		Read:   resourceArmContainerRegistryNetworkRulesetRead,
		Update: resourceArmContainerRegistryNetworkRulesetCreateUpdate,
		Delete: resourceArmContainerRegistryNetworkRulesetDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 2,

		//Setting custom timeouts
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"container_registry_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureRMContainerRegistryName,
			},

			"network_rule_set": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  containerregistry.DefaultActionAllow,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerregistry.DefaultActionAllow),
								string(containerregistry.DefaultActionDeny),
							}, false),
						},

						"ip_rule": {
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(containerregistry.Allow),
										}, false),
									},

									"ip_range": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.CIDR,
									},
								},
							},
						},

						"virtual_network": {
							Type:       schema.TypeSet,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(containerregistry.Allow),
										}, false),
									},

									"subnet_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmContainerRegistryNetworkRulesetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerRegistryName := d.Get("container_registry_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	containerRegistry, err := client.Get(ctx, resourceGroup, containerRegistryName)
	if err != nil {
		if utils.ResponseWasNotFound(containerRegistry.Response) {
			return fmt.Errorf("Container Registry %q (Resource Group %q) was not found", containerRegistryName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	//Generate a Resource ID manually as Azure Resource Mangager doesn't maintain an ID for this type of resource
	resourceId := fmt.Sprintf("%s/networkruleset/rule", *containerRegistry.ID)

	networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]interface{}))
	if networkRuleSet != nil && !strings.EqualFold(string(containerRegistry.Sku.Tier), string(containerregistry.Premium)) {
		return fmt.Errorf("`network_rule_set_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU plese set network_rule_set = []")
	}

	parameters := containerregistry.RegistryUpdateParameters{
		RegistryPropertiesUpdateParameters: &containerregistry.RegistryPropertiesUpdateParameters{
			NetworkRuleSet: networkRuleSet,
		},
	}

	future, err := client.Update(ctx, resourceGroup, containerRegistryName, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	return resourceArmContainerRegistryNetworkRulesetRead(d, meta)
}

func resourceArmContainerRegistryNetworkRulesetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	containerRegistryName := id.Path["registries"]

	containerRegistry, err := client.Get(ctx, resourceGroup, containerRegistryName)
	if err != nil {
		if utils.ResponseWasNotFound(containerRegistry.Response) {
			return fmt.Errorf("Azure Container Registry %q (Resource Group %q) was not found", containerRegistryName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Azure Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	d.Set("container_registry_name", containerRegistryName)
	d.Set("resource_group_name", resourceGroup)
	if location := containerRegistry.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if rules := containerRegistry.NetworkRuleSet; rules != nil {
		if err := d.Set("network_rule_set", flattenNetworkRuleSet(rules)); err != nil {
			return fmt.Errorf("Error setting `network_rule_set` for Azure Container Registry %q: %+v", containerRegistryName, err)
		}
	}

	return nil
}

func resourceArmContainerRegistryNetworkRulesetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	containerRegistryName := id.Path["registries"]

	containerRegistry, err := client.Get(ctx, resourceGroup, containerRegistryName)
	if err != nil {
		if utils.ResponseWasNotFound(containerRegistry.Response) {
			return fmt.Errorf("Azure Container Registry %q (Resource Group %q) was not found", containerRegistryName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Azure Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	if containerRegistry.NetworkRuleSet == nil {
		return nil
	}

	//Delete operation doesn't delete any of the resources, instead it updates the resource to initial configuration
	parameters := containerregistry.RegistryUpdateParameters{
		RegistryPropertiesUpdateParameters: &containerregistry.RegistryPropertiesUpdateParameters{
			NetworkRuleSet: &containerregistry.NetworkRuleSet{
				DefaultAction: containerregistry.DefaultActionAllow,
			},
		},
	}

	future, err := client.Update(ctx, resourceGroup, containerRegistryName, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
	}

	return nil
}
