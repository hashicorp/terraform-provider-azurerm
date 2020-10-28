package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	validateNetwork "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventHubNamespaceNetworkRuleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventHubNamespaceNetworkRuleSetCreateUpdate,
		Read:   resourceEventHubNamespaceNetworkRuleSetRead,
		Update: resourceEventHubNamespaceNetworkRuleSetCreateUpdate,
		Delete: resourceEventHubNamespaceNetworkRuleSetDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.EventHubNamespaceNetworkRuleSetID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"default_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(eventhub.Allow),
				ValidateFunc: validation.StringInSlice([]string{
					string(eventhub.Allow),
					string(eventhub.Deny),
				}, false),
			},

			"trusted_service_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ip_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"network_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      networkRuleHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNetwork.SubnetID,
							// The subnet ID returned from the service will have `resourceGroup/{resourceGroupName}` all in lower cases...
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"ignore_missing_vnet_service_endpoint": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceEventHubNamespaceNetworkRuleSetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetNetworkRuleSet(ctx, resourceGroup, namespaceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed to check for presence of existing Service Bus Namespace Network Rule Set (Namespace %q / Resource Group %q): %+v", namespaceName, resourceGroup, err)
			}
		}

		// This resource is unique to the corresponding service bus namespace.
		// It will be created automatically along with the namespace, therefore we check whether this resource is identical to a "deleted" one
		if !CheckNetworkRuleNullified(existing) {
			return tf.ImportAsExistsError("azurerm_EventHub_namespace_network_rule_set", *existing.ID)
		}
	}

	parameters := eventhub.NetworkRuleSet{
		NetworkRuleSetProperties: &eventhub.NetworkRuleSetProperties{
			DefaultAction:       eventhub.DefaultAction(d.Get("default_action").(string)),
			VirtualNetworkRules: expandEventHubNamespaceVirtualNetworkRules(d.Get("network_rules").(*schema.Set).List()),
			IPRules:             expandEventHubNamespaceIPRules(d.Get("ip_rules").(*schema.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resourceGroup, namespaceName, parameters); err != nil {
		return fmt.Errorf("failed to create Service Bus Namespace Network Rule Set (Namespace %q / Resource Group %q): %+v", namespaceName, resourceGroup, err)
	}

	resp, err := client.GetNetworkRuleSet(ctx, resourceGroup, namespaceName)
	if err != nil {
		return fmt.Errorf("failed to retrieve Service Bus Namespace Network Rule Set (Namespace %q / Resource Group %q): %+v", namespaceName, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Service Bus Namespace Network Rule Set (Namespace %q / Resource Group %q) ID", namespaceName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceEventHubNamespaceNetworkRuleSetRead(d, meta)
}

func resourceEventHubNamespaceNetworkRuleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventHubNamespaceNetworkRuleSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Service Bus Namespace Network Rule Set %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read Service Bus Namespace Network Rule Set %q (Namespace %q / Resource Group %q): %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.NetworkRuleSetProperties; props != nil {
		d.Set("default_action", string(props.DefaultAction))

		if err := d.Set("network_rules", schema.NewSet(networkRuleHash, flattenEventHubNamespaceVirtualNetworkRules(props.VirtualNetworkRules))); err != nil {
			return fmt.Errorf("failed to set `network_rules`: %+v", err)
		}

		if err := d.Set("ip_rules", flattenEventHubNamespaceIPRules(props.IPRules)); err != nil {
			return fmt.Errorf("failed to set `ip_rules`: %+v", err)
		}
	}

	return nil
}

func resourceEventHubNamespaceNetworkRuleSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventHubNamespaceNetworkRuleSetID(d.Id())
	if err != nil {
		return err
	}

	// A network rule is unique to a namespace, this rule cannot be deleted.
	// Therefore we here are just disabling it by setting the default_action to allow and remove all its rules and masks

	parameters := eventhub.NetworkRuleSet{
		NetworkRuleSetProperties: &eventhub.NetworkRuleSetProperties{
			DefaultAction: eventhub.Deny,
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName, parameters); err != nil {
		return fmt.Errorf("failed to delete Service Bus Namespace Network Rule Set %q (Namespace %q / Resource Group %q): %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	return nil
}

func expandEventHubNamespaceVirtualNetworkRules(input []interface{}) *[]eventhub.NWRuleSetVirtualNetworkRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]eventhub.NWRuleSetVirtualNetworkRules, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, eventhub.NWRuleSetVirtualNetworkRules{
			Subnet: &eventhub.Subnet{
				ID: utils.String(raw["subnet_id"].(string)),
			},
			IgnoreMissingVnetServiceEndpoint: utils.Bool(raw["ignore_missing_vnet_service_endpoint"].(bool)),
		})
	}

	return &result
}

func flattenEventHubNamespaceVirtualNetworkRules(input *[]eventhub.NWRuleSetVirtualNetworkRules) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		subnetId := ""
		if v.Subnet != nil && v.Subnet.ID != nil {
			subnetId = *v.Subnet.ID
		}

		ignore := false
		if v.IgnoreMissingVnetServiceEndpoint != nil {
			ignore = *v.IgnoreMissingVnetServiceEndpoint
		}

		result = append(result, map[string]interface{}{
			"subnet_id":                            subnetId,
			"ignore_missing_vnet_service_endpoint": ignore,
		})
	}

	return result
}

func expandEventHubNamespaceIPRules(input []interface{}) *[]eventhub.NWRuleSetIPRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]eventhub.NWRuleSetIPRules, 0)
	for _, v := range input {
		result = append(result, eventhub.NWRuleSetIPRules{
			IPMask: utils.String(v.(string)),
			Action: eventhub.NetworkRuleIPActionAllow,
		})
	}

	return &result
}

func flattenEventHubNamespaceIPRules(input *[]eventhub.NWRuleSetIPRules) []interface{} {
	result := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return result
	}

	for _, v := range *input {
		if v.IPMask != nil {
			result = append(result, *v.IPMask)
		}
	}

	return result
}

func networkRuleHash(input interface{}) int {
	v := input.(map[string]interface{})

	// we are just taking subnet_id into the hash function and ignore the ignore_missing_vnet_service_endpoint to ensure there would be no duplicates of subnet id
	// the service returns this ID with segment resourceGroup and resource group name all in lower cases, to avoid unnecessary diff, we extract this ID and reconstruct this hash code
	return set.HashStringIgnoreCase(v["subnet_id"])
}

func CheckNetworkRuleNullified(resp eventhub.NetworkRuleSet) bool {
	if resp.ID == nil || *resp.ID == "" {
		return true
	}
	if resp.NetworkRuleSetProperties == nil {
		return true
	}
	if resp.DefaultAction != eventhub.Deny {
		return false
	}
	if resp.VirtualNetworkRules != nil && len(*resp.VirtualNetworkRules) > 0 {
		return false
	}
	if resp.IPRules != nil && len(*resp.IPRules) > 0 {
		return false
	}
	return true
}
