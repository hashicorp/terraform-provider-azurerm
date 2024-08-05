// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusNamespaceNetworkRuleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusNamespaceNetworkRuleSetCreateUpdate,
		Read:   resourceServiceBusNamespaceNetworkRuleSetRead,
		Update: resourceServiceBusNamespaceNetworkRuleSetCreateUpdate,
		Delete: resourceServiceBusNamespaceNetworkRuleSetDelete,

		DeprecationMessage: "The `azurerm_servicebus_namespace_network_rule_set` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use `network_rule_set` inside the `azurerm_servicebus_namespace` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := namespaces.ParseNamespaceID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NamespaceNetworkRuleSetV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServicebusNamespaceNetworkRuleSetSchema(),
	}
}

func resourceServicebusNamespaceNetworkRuleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		//lintignore: S013
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"default_action": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(namespaces.DefaultActionAllow),
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.DefaultActionAllow),
				string(namespaces.DefaultActionDeny),
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
						ValidateFunc: commonids.ValidateSubnetID,
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
	}
}

func resourceServiceBusNamespaceNetworkRuleSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.GetNetworkRuleSet(ctx, *id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %s: %+v", id, err)
			}
		}

		// This resource is unique to the corresponding service bus namespace.
		// It will be created automatically along with the namespace, therefore we check whether this resource is identical to a "deleted" one
		if model := existing.Model; model != nil {
			if !CheckNetworkRuleNullified(*model) {
				return tf.ImportAsExistsError("azurerm_servicebus_namespace_network_rule_set", id.ID())
			}
		}
	}

	defaultAction := namespaces.DefaultAction(d.Get("default_action").(string))
	vnetRule := expandServiceBusNamespaceVirtualNetworkRules(d.Get("network_rules").(*pluginsdk.Set).List())
	ipRule := expandServiceBusNamespaceIPRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	publicNetworkAcc := "Disabled"
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAcc = "Enabled"
	}

	// API doesn't accept "Deny" to be set for "default_action" if no "ip_rules" or "network_rules" is defined and returns no error message to the user
	if defaultAction == namespaces.DefaultActionDeny && vnetRule == nil && ipRule == nil {
		return fmt.Errorf(" The default action of %s can only be set to `Allow` if no `ip_rules` or `network_rules` is set", id)
	}

	publicNetworkAccess := namespaces.PublicNetworkAccessFlag(publicNetworkAcc)

	parameters := namespaces.NetworkRuleSet{
		Properties: &namespaces.NetworkRuleSetProperties{
			DefaultAction:               &defaultAction,
			VirtualNetworkRules:         vnetRule,
			IPRules:                     ipRule,
			PublicNetworkAccess:         &publicNetworkAccess,
			TrustedServiceAccessEnabled: utils.Bool(d.Get("trusted_services_allowed").(bool)),
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, *id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceServiceBusNamespaceNetworkRuleSetRead(d, meta)
}

func resourceServiceBusNamespaceNetworkRuleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetNetworkRuleSet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("namespace_id", id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			defaultAction := ""
			if v := props.DefaultAction; v != nil {
				defaultAction = string(*v)
			}
			d.Set("default_action", defaultAction)
			d.Set("trusted_services_allowed", props.TrustedServiceAccessEnabled)
			publicNetworkAccess := "Enabled"
			if v := props.PublicNetworkAccess; v != nil {
				publicNetworkAccess = string(*v)
			}
			d.Set("public_network_access_enabled", strings.EqualFold(publicNetworkAccess, "Enabled"))

			if err := d.Set("network_rules", pluginsdk.NewSet(networkRuleHash, flattenServiceBusNamespaceVirtualNetworkRules(props.VirtualNetworkRules))); err != nil {
				return fmt.Errorf("failed to set `network_rules`: %+v", err)
			}

			if err := d.Set("ip_rules", flattenServiceBusNamespaceIPRules(props.IPRules)); err != nil {
				return fmt.Errorf("failed to set `ip_rules`: %+v", err)
			}
		}
	}

	return nil
}

func resourceServiceBusNamespaceNetworkRuleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	// A network rule is unique to a namespace, this rule cannot be deleted.
	// Therefore we here are just disabling it by setting the default_action to allow and remove all its rules and masks

	defaultAction := namespaces.DefaultActionAllow
	parameters := namespaces.NetworkRuleSet{
		Properties: &namespaces.NetworkRuleSetProperties{
			DefaultAction: &defaultAction,
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, *id, parameters); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandServiceBusNamespaceVirtualNetworkRules(input []interface{}) *[]namespaces.NWRuleSetVirtualNetworkRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]namespaces.NWRuleSetVirtualNetworkRules, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, namespaces.NWRuleSetVirtualNetworkRules{
			Subnet: &namespaces.Subnet{
				Id: raw["subnet_id"].(string),
			},
			IgnoreMissingVnetServiceEndpoint: utils.Bool(raw["ignore_missing_vnet_service_endpoint"].(bool)),
		})
	}

	return &result
}

func flattenServiceBusNamespaceVirtualNetworkRules(input *[]namespaces.NWRuleSetVirtualNetworkRules) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		subnetId := ""
		if v.Subnet != nil && v.Subnet.Id != "" {
			subnetId = v.Subnet.Id
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

func expandServiceBusNamespaceIPRules(input []interface{}) *[]namespaces.NWRuleSetIPRules {
	if len(input) == 0 {
		return nil
	}

	action := namespaces.NetworkRuleIPActionAllow
	result := make([]namespaces.NWRuleSetIPRules, 0)
	for _, v := range input {
		result = append(result, namespaces.NWRuleSetIPRules{
			IPMask: utils.String(v.(string)),
			Action: &action,
		})
	}

	return &result
}

func flattenServiceBusNamespaceIPRules(input *[]namespaces.NWRuleSetIPRules) []interface{} {
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

func CheckNetworkRuleNullified(resp namespaces.NetworkRuleSet) bool {
	if resp.Id == nil || *resp.Id == "" {
		return true
	}

	if props := resp.Properties; props != nil {
		if *props.DefaultAction != namespaces.DefaultActionAllow {
			return false
		}

		if props.VirtualNetworkRules != nil && len(*props.VirtualNetworkRules) > 0 {
			return false
		}

		if props.IPRules != nil && len(*props.IPRules) > 0 {
			return false
		}
	}

	return true
}
