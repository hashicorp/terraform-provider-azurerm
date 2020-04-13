package servicebus

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2018-01-01-preview/servicebus"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceBusNamespaceVirtualNetworkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusNamespaceVirtualNetworkRuleCreateUpdate,
		Read:   resourceArmServiceBusNamespaceVirtualNetworkRuleRead,
		Update: resourceArmServiceBusNamespaceVirtualNetworkRuleCreateUpdate,
		Delete: resourceArmServiceBusNamespaceVirtualNetworkRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusVirtualNetworkRuleName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceBusNamespaceName(),
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmServiceBusNamespaceVirtualNetworkRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	virtualNetworkSubnetId := d.Get("subnet_id").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetVirtualNetworkRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing ServiceBus Namespace Virtual Network Rule %q (ServiceBus Namespace: %q, Resource Group: %q): %+v", name, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_virtual_network_rule", *existing.ID)
		}
	}

	parameters := servicebus.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &servicebus.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID: utils.String(virtualNetworkSubnetId),
		},
	}

	if _, err := client.CreateOrUpdateVirtualNetworkRule(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("Error creating ServiceBus Namespace Virtual Network Rule %q (ServiceBus Namespace: %q, Resource Group: %q): %+v", name, namespaceName, resourceGroup, err)
	}

	resp, err := client.GetVirtualNetworkRule(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving ServiceBus Namespace Virtual Network Rule %q (ServiceBus Namespace: %q, Resource Group: %q): %+v", name, namespaceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmServiceBusNamespaceVirtualNetworkRuleRead(d, meta)
}

func resourceArmServiceBusNamespaceVirtualNetworkRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["virtualnetworkrules"]

	resp, err := client.GetVirtualNetworkRule(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading ServiceBus Namespace Virtual Network Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading ServiceBus Namespace Virtual Network Rule: %q (ServiceBus Namespace: %q, Resource Group: %q): %+v", name, namespaceName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("namespace_name", namespaceName)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.VirtualNetworkSubnetID)
	}

	return nil
}

func resourceArmServiceBusNamespaceVirtualNetworkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClientPreview
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["virtualnetworkrules"]

	if _, err = client.DeleteVirtualNetworkRule(ctx, resourceGroup, namespaceName, name); err != nil {
		return fmt.Errorf("Error issuing Azure ARM delete request of ServiceBus Namespace Virtual Network Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

/*
	This function checks the format of the ServiceBus Namespace Virtual Network Rule Name to make sure that
	it does not contain any potentially invalid values.
*/
func ValidateServiceBusNamespaceVirtualNetworkRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Cannot be empty
	if len(value) == 0 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be an empty string: %q", k, value))
	}

	// Cannot be more than 128 characters
	if len(value) > 128 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 128 characters: %q", k, value))
	}

	// Must only contain alphanumeric characters or hyphens
	if !regexp.MustCompile(`^[A-Za-z0-9-]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q can only contain alphanumeric characters and hyphens: %q",
			k, value))
	}

	// Cannot end in a hyphen
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen: %q", k, value))
	}

	// Cannot start with a number or hyphen
	if regexp.MustCompile(`^[0-9-]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with a number or hyphen: %q", k, value))
	}

	// There are multiple returns in the case that there is more than one invalid
	// case applied to the name.
	return warnings, errors
}
