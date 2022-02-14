package servicebus

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2021-06-01-preview/servicebus"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	validateNetwork "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// the only allowed value at this time
var namespaceNetworkRuleSetName = "default"

func resourceServiceBusNamespaceNetworkRuleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusNamespaceNetworkRuleSetCreateUpdate,
		Read:   resourceServiceBusNamespaceNetworkRuleSetRead,
		Update: resourceServiceBusNamespaceNetworkRuleSetCreateUpdate,
		Delete: resourceServiceBusNamespaceNetworkRuleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NamespaceNetworkRuleSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			// TODO 3.0 - Make it required
			"namespace_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validate.NamespaceID,
				ConflictsWith: []string{"namespace_name", "resource_group_name"},
			},

			// TODO 3.0 - Remove in favor of namespace_id
			"namespace_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validate.NamespaceName,
				Deprecated:    `Deprecated in favor of "namespace_id"`,
				ConflictsWith: []string{"namespace_id"},
			},

			// TODO 3.0 - Remove in favor of namespace_id
			"resource_group_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  azure.ValidateResourceGroupName,
				Deprecated:    `Deprecated in favor of "namespace_id"`,
				ConflictsWith: []string{"namespace_id"},
			},

			"default_action": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(servicebus.DefaultActionAllow),
				ValidateFunc: validation.StringInSlice([]string{
					string(servicebus.DefaultActionAllow),
					string(servicebus.DefaultActionDeny),
				}, false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ip_rules": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"trusted_services_allowed": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"network_rules": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Set:      networkRuleHash,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validateNetwork.SubnetID,
							// The subnet ID returned from the service will have `resourceGroup/{resourceGroupName}` all in lower cases...
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"ignore_missing_vnet_service_endpoint": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceServiceBusNamespaceNetworkRuleSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var resourceId parse.NamespaceNetworkRuleSetId
	if namespaceIdLit := d.Get("namespace_id").(string); namespaceIdLit != "" {
		namespaceId, _ := parse.NamespaceID(namespaceIdLit)
		resourceId = parse.NewNamespaceNetworkRuleSetID(namespaceId.SubscriptionId, namespaceId.ResourceGroup, namespaceId.Name, namespaceNetworkRuleSetName)
	} else {
		resourceId = parse.NewNamespaceNetworkRuleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), namespaceNetworkRuleSetName)
	}

	if d.IsNewResource() {
		existing, err := client.GetNetworkRuleSet(ctx, resourceId.ResourceGroup, resourceId.NamespaceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of existing %s: %+v", resourceId, err)
			}
		}

		// This resource is unique to the corresponding service bus namespace.
		// It will be created automatically along with the namespace, therefore we check whether this resource is identical to a "deleted" one
		if !CheckNetworkRuleNullified(existing) {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_network_rule_set", resourceId.ID())
		}
	}

	defaultAction := servicebus.DefaultAction(d.Get("default_action").(string))
	vnetRule := expandServiceBusNamespaceVirtualNetworkRules(d.Get("network_rules").(*pluginsdk.Set).List())
	ipRule := expandServiceBusNamespaceIPRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	publicNetworkAcc := "Disabled"
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAcc = "Enabled"
	}

	// API doesn't accept "Deny" to be set for "default_action" if no "ip_rules" or "network_rules" is defined and returns no error message to the user
	// TODO: The check won't be needed when 2021-11-01 API is released since service team will fail the update with bad request in that version
	if defaultAction == servicebus.DefaultActionDeny && vnetRule == nil && ipRule == nil {
		return fmt.Errorf(" The default action of %s can only be set to `Allow` if no `ip_rules` or `network_rules` is set", resourceId)
	}

	parameters := servicebus.NetworkRuleSet{
		NetworkRuleSetProperties: &servicebus.NetworkRuleSetProperties{
			DefaultAction:               defaultAction,
			VirtualNetworkRules:         vnetRule,
			IPRules:                     ipRule,
			PublicNetworkAccess:         servicebus.PublicNetworkAccessFlag(publicNetworkAcc),
			TrustedServiceAccessEnabled: utils.Bool(d.Get("trusted_services_allowed").(bool)),
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resourceId.ResourceGroup, resourceId.NamespaceName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceServiceBusNamespaceNetworkRuleSetRead(d, meta)
}

func resourceServiceBusNamespaceNetworkRuleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceNetworkRuleSetID(d.Id())
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
		return fmt.Errorf("failed to read Service Bus Namespace Network Rule Set %q (Namespace %q / Resource Group %q): %+v", id.NetworkrulesetName, id.NamespaceName, id.ResourceGroup, err)
	}

	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("namespace_id", parse.NewNamespaceID(id.SubscriptionId, id.ResourceGroup, id.NamespaceName).ID())

	if props := resp.NetworkRuleSetProperties; props != nil {
		d.Set("default_action", string(props.DefaultAction))
		d.Set("trusted_services_allowed", props.TrustedServiceAccessEnabled)
		d.Set("public_network_access_enabled", strings.EqualFold(string(props.PublicNetworkAccess), "Enabled"))

		if err := d.Set("network_rules", pluginsdk.NewSet(networkRuleHash, flattenServiceBusNamespaceVirtualNetworkRules(props.VirtualNetworkRules))); err != nil {
			return fmt.Errorf("failed to set `network_rules`: %+v", err)
		}

		if err := d.Set("ip_rules", flattenServiceBusNamespaceIPRules(props.IPRules)); err != nil {
			return fmt.Errorf("failed to set `ip_rules`: %+v", err)
		}
	}

	return nil
}

func resourceServiceBusNamespaceNetworkRuleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceNetworkRuleSetID(d.Id())
	if err != nil {
		return err
	}

	// A network rule is unique to a namespace, this rule cannot be deleted.
	// Therefore we here are just disabling it by setting the default_action to allow and remove all its rules and masks

	parameters := servicebus.NetworkRuleSet{
		NetworkRuleSetProperties: &servicebus.NetworkRuleSetProperties{
			DefaultAction: servicebus.DefaultActionAllow,
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName, parameters); err != nil {
		return fmt.Errorf("failed to delete Service Bus Namespace Network Rule Set %q (Namespace %q / Resource Group %q): %+v", id.NetworkrulesetName, id.NamespaceName, id.ResourceGroup, err)
	}

	return nil
}

func expandServiceBusNamespaceVirtualNetworkRules(input []interface{}) *[]servicebus.NWRuleSetVirtualNetworkRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]servicebus.NWRuleSetVirtualNetworkRules, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, servicebus.NWRuleSetVirtualNetworkRules{
			Subnet: &servicebus.Subnet{
				ID: utils.String(raw["subnet_id"].(string)),
			},
			IgnoreMissingVnetServiceEndpoint: utils.Bool(raw["ignore_missing_vnet_service_endpoint"].(bool)),
		})
	}

	return &result
}

func flattenServiceBusNamespaceVirtualNetworkRules(input *[]servicebus.NWRuleSetVirtualNetworkRules) []interface{} {
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

func expandServiceBusNamespaceIPRules(input []interface{}) *[]servicebus.NWRuleSetIPRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]servicebus.NWRuleSetIPRules, 0)
	for _, v := range input {
		result = append(result, servicebus.NWRuleSetIPRules{
			IPMask: utils.String(v.(string)),
			Action: servicebus.NetworkRuleIPActionAllow,
		})
	}

	return &result
}

func flattenServiceBusNamespaceIPRules(input *[]servicebus.NWRuleSetIPRules) []interface{} {
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

func CheckNetworkRuleNullified(resp servicebus.NetworkRuleSet) bool {
	if resp.ID == nil || *resp.ID == "" {
		return true
	}
	if resp.NetworkRuleSetProperties == nil {
		return true
	}
	if resp.DefaultAction != servicebus.DefaultActionAllow {
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
