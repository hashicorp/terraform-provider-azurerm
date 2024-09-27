// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const AzureFirewallPolicyResourceName = "azurerm_firewall_policy"

func resourceFirewallPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallPolicyCreateUpdate,
		Read:   resourceFirewallPolicyRead,
		Update: resourceFirewallPolicyCreateUpdate,
		Delete: resourceFirewallPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := firewallpolicies.ParseFirewallPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceFirewallPolicySchema(),
	}
}

func resourceFirewallPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := firewallpolicies.NewFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id, firewallpolicies.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if resp.Model != nil {
			return tf.ImportAsExistsError("azurerm_firewall_policy", id.ID())
		}
	}

	props := firewallpolicies.FirewallPolicy{
		Properties: &firewallpolicies.FirewallPolicyPropertiesFormat{
			ThreatIntelMode:      pointer.To(firewallpolicies.AzureFirewallThreatIntelMode(d.Get("threat_intelligence_mode").(string))),
			ThreatIntelWhitelist: expandFirewallPolicyThreatIntelWhitelist(d.Get("threat_intelligence_allowlist").([]interface{})),
			DnsSettings:          expandFirewallPolicyDNSSetting(d.Get("dns").([]interface{})),
			IntrusionDetection:   expandFirewallPolicyIntrusionDetection(d.Get("intrusion_detection").([]interface{})),
			TransportSecurity:    expandFirewallPolicyTransportSecurity(d.Get("tls_certificate").([]interface{})),
			Insights:             expandFirewallPolicyInsights(d.Get("insights").([]interface{})),
			ExplicitProxy:        expandFirewallPolicyExplicitProxy(d.Get("explicit_proxy").([]interface{})),
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	// api will error if TypeNone is passed in
	if expandedIdentity.Type != identity.TypeNone {
		props.Identity = expandedIdentity
	}

	if id, ok := d.GetOk("base_policy_id"); ok {
		props.Properties.BasePolicy = &firewallpolicies.SubResource{Id: utils.String(id.(string))}
	}

	if v, ok := d.GetOk("sku"); ok {
		props.Properties.Sku = &firewallpolicies.FirewallPolicySku{
			Tier: pointer.To(firewallpolicies.FirewallPolicySkuTier(v.(string))),
		}
	}

	if v, ok := d.GetOk("sql_redirect_allowed"); ok {
		props.Properties.Sql = &firewallpolicies.FirewallPolicySQL{
			AllowSqlRedirect: utils.Bool(v.(bool)),
		}
	}

	if v, ok := d.GetOk("private_ip_ranges"); ok {
		privateIPRanges := utils.ExpandStringSlice(v.([]interface{}))
		props.Properties.Snat = &firewallpolicies.FirewallPolicySNAT{
			PrivateRanges: privateIPRanges,
		}
	}

	if v, ok := d.GetOk("auto_learn_private_ranges_enabled"); ok {
		if props.Properties.Snat == nil {
			props.Properties.Snat = &firewallpolicies.FirewallPolicySNAT{}
		}
		if v.(bool) {
			props.Properties.Snat.AutoLearnPrivateRanges = pointer.To(firewallpolicies.AutoLearnPrivateRangesModeEnabled)
		}
	}

	locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)

	if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceFirewallPolicyRead(d, meta)
}

func resourceFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallpolicies.ParseFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, firewallpolicies.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			basePolicyID := ""
			if props.BasePolicy != nil && props.BasePolicy.Id != nil {
				basePolicyID = *props.BasePolicy.Id
			}
			d.Set("base_policy_id", basePolicyID)

			d.Set("threat_intelligence_mode", string(pointer.From(props.ThreatIntelMode)))

			if sku := props.Sku; sku != nil {
				d.Set("sku", string(pointer.From(sku.Tier)))
			}

			if err := d.Set("threat_intelligence_allowlist", flattenFirewallPolicyThreatIntelWhitelist(props.ThreatIntelWhitelist)); err != nil {
				return fmt.Errorf(`setting "threat_intelligence_allowlist": %+v`, err)
			}

			if err := d.Set("dns", flattenFirewallPolicyDNSSetting(props.DnsSettings)); err != nil {
				return fmt.Errorf(`setting "dns": %+v`, err)
			}

			if err := d.Set("intrusion_detection", flattenFirewallPolicyIntrusionDetection(props.IntrusionDetection)); err != nil {
				return fmt.Errorf(`setting "intrusion_detection": %+v`, err)
			}

			if err := d.Set("tls_certificate", flattenFirewallPolicyTransportSecurity(props.TransportSecurity)); err != nil {
				return fmt.Errorf(`setting "tls_certificate": %+v`, err)
			}

			if err := d.Set("child_policies", flattenNetworkSubResourceID(props.ChildPolicies)); err != nil {
				return fmt.Errorf(`setting "child_policies": %+v`, err)
			}

			if err := d.Set("firewalls", flattenNetworkSubResourceID(props.Firewalls)); err != nil {
				return fmt.Errorf(`setting "firewalls": %+v`, err)
			}

			if err := d.Set("rule_collection_groups", flattenNetworkSubResourceID(props.RuleCollectionGroups)); err != nil {
				return fmt.Errorf(`setting "rule_collection_groups": %+v`, err)
			}

			var privateIPRanges []interface{}
			var isAutoLearnPrivateRangeEnabled bool
			if props.Snat != nil {
				privateIPRanges = utils.FlattenStringSlice(props.Snat.PrivateRanges)
				isAutoLearnPrivateRangeEnabled = pointer.From(props.Snat.AutoLearnPrivateRanges) == firewallpolicies.AutoLearnPrivateRangesModeEnabled
			}
			if err := d.Set("private_ip_ranges", privateIPRanges); err != nil {
				return fmt.Errorf("setting `private_ip_ranges`: %+v", err)
			}

			if err := d.Set("auto_learn_private_ranges_enabled", isAutoLearnPrivateRangeEnabled); err != nil {
				return fmt.Errorf("setting `auto_learn_private_ranges_enabled`: %+v", err)
			}

			if err := d.Set("insights", flattenFirewallPolicyInsights(props.Insights)); err != nil {
				return fmt.Errorf(`setting "insights": %+v`, err)
			}

			proxySettings := flattenFirewallPolicyExplicitProxy(props.ExplicitProxy)
			if err := d.Set("explicit_proxy", proxySettings); err != nil {
				return fmt.Errorf("setting `explicit_proxy`: %+v", err)
			}

			if props.Sql != nil && props.Sql.AllowSqlRedirect != nil {
				if err := d.Set("sql_redirect_allowed", props.Sql.AllowSqlRedirect); err != nil {
					return fmt.Errorf("setting `sql_redirect_allowed`: %+v", err)
				}
			}
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallpolicies.ParseFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandFirewallPolicyThreatIntelWhitelist(input []interface{}) *firewallpolicies.FirewallPolicyThreatIntelWhitelist {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &firewallpolicies.FirewallPolicyThreatIntelWhitelist{
		IPAddresses: utils.ExpandStringSlice(raw["ip_addresses"].(*pluginsdk.Set).List()),
		Fqdns:       utils.ExpandStringSlice(raw["fqdns"].(*pluginsdk.Set).List()),
	}

	return output
}

func expandFirewallPolicyDNSSetting(input []interface{}) *firewallpolicies.DnsSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &firewallpolicies.DnsSettings{
		Servers:     utils.ExpandStringSlice(raw["servers"].([]interface{})),
		EnableProxy: utils.Bool(raw["proxy_enabled"].(bool)),
	}

	return output
}

func expandFirewallPolicyIntrusionDetection(input []interface{}) *firewallpolicies.FirewallPolicyIntrusionDetection {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	var signatureOverrides []firewallpolicies.FirewallPolicyIntrusionDetectionSignatureSpecification
	for _, v := range raw["signature_overrides"].([]interface{}) {
		overrides := v.(map[string]interface{})
		signatureOverrides = append(signatureOverrides, firewallpolicies.FirewallPolicyIntrusionDetectionSignatureSpecification{
			Id:   utils.String(overrides["id"].(string)),
			Mode: pointer.To(firewallpolicies.FirewallPolicyIntrusionDetectionStateType(overrides["state"].(string))),
		})
	}

	var trafficBypass []firewallpolicies.FirewallPolicyIntrusionDetectionBypassTrafficSpecifications

	for _, v := range raw["traffic_bypass"].([]interface{}) {
		bypass := v.(map[string]interface{})
		trafficBypass = append(trafficBypass, firewallpolicies.FirewallPolicyIntrusionDetectionBypassTrafficSpecifications{
			Name:                 utils.String(bypass["name"].(string)),
			Description:          utils.String(bypass["description"].(string)),
			Protocol:             pointer.To(firewallpolicies.FirewallPolicyIntrusionDetectionProtocol(bypass["protocol"].(string))),
			SourceAddresses:      utils.ExpandStringSlice(bypass["source_addresses"].(*pluginsdk.Set).List()),
			DestinationAddresses: utils.ExpandStringSlice(bypass["destination_addresses"].(*pluginsdk.Set).List()),
			DestinationPorts:     utils.ExpandStringSlice(bypass["destination_ports"].(*pluginsdk.Set).List()),
			SourceIPGroups:       utils.ExpandStringSlice(bypass["source_ip_groups"].(*pluginsdk.Set).List()),
			DestinationIPGroups:  utils.ExpandStringSlice(bypass["destination_ip_groups"].(*pluginsdk.Set).List()),
		})
	}

	var privateRanges []string
	for _, v := range raw["private_ranges"].([]interface{}) {
		privateRanges = append(privateRanges, v.(string))
	}

	return &firewallpolicies.FirewallPolicyIntrusionDetection{
		Mode: pointer.To(firewallpolicies.FirewallPolicyIntrusionDetectionStateType(raw["mode"].(string))),
		Configuration: &firewallpolicies.FirewallPolicyIntrusionDetectionConfiguration{
			SignatureOverrides:    &signatureOverrides,
			PrivateRanges:         &privateRanges,
			BypassTrafficSettings: &trafficBypass,
		},
	}
}

func expandFirewallPolicyTransportSecurity(input []interface{}) *firewallpolicies.FirewallPolicyTransportSecurity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &firewallpolicies.FirewallPolicyTransportSecurity{
		CertificateAuthority: &firewallpolicies.FirewallPolicyCertificateAuthority{
			KeyVaultSecretId: utils.String(raw["key_vault_secret_id"].(string)),
			Name:             utils.String(raw["name"].(string)),
		},
	}
}

func expandFirewallPolicyInsights(input []interface{}) *firewallpolicies.FirewallPolicyInsights {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &firewallpolicies.FirewallPolicyInsights{
		IsEnabled:             utils.Bool(raw["enabled"].(bool)),
		RetentionDays:         utils.Int64(int64(raw["retention_in_days"].(int))),
		LogAnalyticsResources: expandFirewallPolicyLogAnalyticsResources(raw["default_log_analytics_workspace_id"].(string), raw["log_analytics_workspace"].([]interface{})),
	}

	return output
}

func expandFirewallPolicyExplicitProxy(input []interface{}) *firewallpolicies.ExplicitProxy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	if raw == nil {
		return nil
	}

	output := &firewallpolicies.ExplicitProxy{
		EnableExplicitProxy: utils.Bool(raw["enabled"].(bool)),
		HTTPPort:            utils.Int64(int64(raw["http_port"].(int))),
		HTTPSPort:           utils.Int64(int64(raw["https_port"].(int))),
		PacFilePort:         utils.Int64(int64(raw["pac_file_port"].(int))),
		PacFile:             utils.String(raw["pac_file"].(string)),
	}

	if val, ok := raw["enable_pac_file"]; ok {
		output.EnablePacFile = utils.Bool(val.(bool))
	}

	return output
}

func expandFirewallPolicyLogAnalyticsResources(defaultWorkspaceId string, workspaces []interface{}) *firewallpolicies.FirewallPolicyLogAnalyticsResources {
	output := &firewallpolicies.FirewallPolicyLogAnalyticsResources{
		DefaultWorkspaceId: &firewallpolicies.SubResource{
			Id: &defaultWorkspaceId,
		},
	}

	var workspaceList []firewallpolicies.FirewallPolicyLogAnalyticsWorkspace
	for _, workspace := range workspaces {
		workspace := workspace.(map[string]interface{})
		workspaceList = append(workspaceList, firewallpolicies.FirewallPolicyLogAnalyticsWorkspace{
			Region: utils.String(location.Normalize(workspace["firewall_location"].(string))),
			WorkspaceId: &firewallpolicies.SubResource{
				Id: utils.String(workspace["id"].(string)),
			},
		})
	}
	if workspaceList != nil {
		output.Workspaces = &workspaceList
	}

	return output
}

func flattenFirewallPolicyThreatIntelWhitelist(input *firewallpolicies.FirewallPolicyThreatIntelWhitelist) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"ip_addresses": utils.FlattenStringSlice(input.IPAddresses),
			"fqdns":        utils.FlattenStringSlice(input.Fqdns),
		},
	}
}

func flattenFirewallPolicyDNSSetting(input *firewallpolicies.DnsSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	proxyEnabled := false
	if input.EnableProxy != nil {
		proxyEnabled = *input.EnableProxy
	}

	return []interface{}{
		map[string]interface{}{
			"servers":       utils.FlattenStringSlice(input.Servers),
			"proxy_enabled": proxyEnabled,
		}}
}

func flattenFirewallPolicyIntrusionDetection(input *firewallpolicies.FirewallPolicyIntrusionDetection) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	signatureOverrides := make([]interface{}, 0)
	trafficBypass := make([]interface{}, 0)

	if input.Configuration == nil {
		return []interface{}{
			map[string]interface{}{
				"mode":                string(pointer.From(input.Mode)),
				"signature_overrides": signatureOverrides,
				"traffic_bypass":      trafficBypass,
			},
		}
	}

	if overrides := input.Configuration.SignatureOverrides; overrides != nil {
		for _, override := range *overrides {
			id := ""
			if override.Id != nil {
				id = *override.Id
			}
			signatureOverrides = append(signatureOverrides, map[string]interface{}{
				"id":    id,
				"state": string(pointer.From(override.Mode)),
			})
		}
	}

	if bypasses := input.Configuration.BypassTrafficSettings; bypasses != nil {
		for _, bypass := range *bypasses {
			name := ""
			if bypass.Name != nil {
				name = *bypass.Name
			}

			description := ""
			if bypass.Description != nil {
				description = *bypass.Description
			}

			var sourceAddresses []string
			if bypass.SourceAddresses != nil {
				sourceAddresses = *bypass.SourceAddresses
			}

			var destinationAddresses []string
			if bypass.DestinationAddresses != nil {
				destinationAddresses = *bypass.DestinationAddresses
			}

			destinationPorts := make([]string, 0)
			if bypass.DestinationPorts != nil {
				destinationPorts = *bypass.DestinationPorts
			}

			sourceIPGroups := make([]string, 0)
			if bypass.SourceIPGroups != nil {
				sourceIPGroups = *bypass.SourceIPGroups
			}

			destinationIPGroups := make([]string, 0)
			if bypass.DestinationIPGroups != nil {
				destinationIPGroups = *bypass.DestinationIPGroups
			}

			trafficBypass = append(trafficBypass, map[string]interface{}{
				"name":                  name,
				"description":           description,
				"protocol":              string(pointer.From(bypass.Protocol)),
				"source_addresses":      sourceAddresses,
				"destination_addresses": destinationAddresses,
				"destination_ports":     destinationPorts,
				"source_ip_groups":      sourceIPGroups,
				"destination_ip_groups": destinationIPGroups,
			})
		}
	}
	var privateRanges []string
	if privates := input.Configuration.PrivateRanges; privates != nil {
		privateRanges = *privates
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                string(pointer.From(input.Mode)),
			"signature_overrides": signatureOverrides,
			"traffic_bypass":      trafficBypass,
			"private_ranges":      privateRanges,
		},
	}
}

func flattenFirewallPolicyTransportSecurity(input *firewallpolicies.FirewallPolicyTransportSecurity) []interface{} {
	if input == nil || input.CertificateAuthority == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_secret_id": input.CertificateAuthority.KeyVaultSecretId,
			"name":                input.CertificateAuthority.Name,
		},
	}
}

func flattenFirewallPolicyInsights(input *firewallpolicies.FirewallPolicyInsights) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var enabled bool
	if input.IsEnabled != nil {
		enabled = *input.IsEnabled
	}

	var retentionInDays int
	if input.RetentionDays != nil {
		retentionInDays = int(*input.RetentionDays)
	}

	defaultLogAnalyticsWorspaceId, logAnalyticsWorkspaces := flattenFirewallPolicyLogAnalyticsResources(input.LogAnalyticsResources)

	return []interface{}{
		map[string]interface{}{
			"enabled":                            enabled,
			"retention_in_days":                  retentionInDays,
			"default_log_analytics_workspace_id": defaultLogAnalyticsWorspaceId,
			"log_analytics_workspace":            logAnalyticsWorkspaces,
		},
	}
}

func flattenFirewallPolicyExplicitProxy(input *firewallpolicies.ExplicitProxy) (result []interface{}) {
	if input == nil {
		return
	}
	output := map[string]interface{}{
		"enabled":         input.EnableExplicitProxy,
		"http_port":       input.HTTPPort,
		"https_port":      input.HTTPSPort,
		"enable_pac_file": input.EnablePacFile,
		"pac_file_port":   input.PacFilePort,
		"pac_file":        input.PacFile,
	}
	return []interface{}{output}
}

func flattenFirewallPolicyLogAnalyticsResources(input *firewallpolicies.FirewallPolicyLogAnalyticsResources) (string, []interface{}) {
	if input == nil {
		return "", []interface{}{}
	}

	var defaultLogAnalyticsWorkspaceId string
	if input.DefaultWorkspaceId != nil && input.DefaultWorkspaceId.Id != nil {
		defaultLogAnalyticsWorkspaceId = *input.DefaultWorkspaceId.Id
	}

	var workspaceList []interface{}
	if input.Workspaces != nil {
		for _, workspace := range *input.Workspaces {
			loc := location.NormalizeNilable(workspace.Region)

			var id string
			if workspace.WorkspaceId != nil && workspace.WorkspaceId.Id != nil {
				id = *workspace.WorkspaceId.Id
			}

			workspaceList = append(workspaceList, map[string]interface{}{
				"id":                id,
				"firewall_location": loc,
			})
		}
	}

	return defaultLogAnalyticsWorkspaceId, workspaceList
}

func resourceFirewallPolicySchema() map[string]*pluginsdk.Schema {
	resource := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallPolicyName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(firewallpolicies.FirewallPolicySkuTierStandard),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(firewallpolicies.FirewallPolicySkuTierPremium),
				string(firewallpolicies.FirewallPolicySkuTierStandard),
				string(firewallpolicies.FirewallPolicySkuTierBasic),
			}, false),
		},

		"location": commonschema.Location(),

		"base_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: firewallpolicies.ValidateFirewallPolicyID,
		},

		"dns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"servers": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
					"proxy_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"threat_intelligence_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(firewallpolicies.AzureFirewallThreatIntelModeAlert),
			ValidateFunc: validation.StringInSlice([]string{
				string(firewallpolicies.AzureFirewallThreatIntelModeAlert),
				string(firewallpolicies.AzureFirewallThreatIntelModeDeny),
				string(firewallpolicies.AzureFirewallThreatIntelModeOff),
			}, false),
		},

		"threat_intelligence_allowlist": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPv4Address),
						},
						AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
					},
					"fqdns": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{"threat_intelligence_allowlist.0.ip_addresses", "threat_intelligence_allowlist.0.fqdns"},
					},
				},
			},
		},

		"intrusion_detection": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							string(firewallpolicies.FirewallPolicyIntrusionDetectionStateTypeOff),
							string(firewallpolicies.FirewallPolicyIntrusionDetectionStateTypeAlert),
							string(firewallpolicies.FirewallPolicyIntrusionDetectionStateTypeDeny),
						}, false),
						Optional: true,
					},
					"signature_overrides": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"state": {
									Type: pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										string(firewallpolicies.FirewallPolicyIntrusionDetectionStateTypeOff),
										string(firewallpolicies.FirewallPolicyIntrusionDetectionStateTypeAlert),
										string(firewallpolicies.FirewallPolicyIntrusionDetectionStateTypeDeny),
									}, false),
									Optional: true,
								},
								"id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"private_ranges": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"traffic_bypass": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"description": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"protocol": {
									Type:     pluginsdk.TypeString,
									Required: true,
									// protocol to be one of [ICMP ANY TCP UDP] but response may be "Any"
									DiffSuppressFunc: suppress.CaseDifference,
									ValidateFunc: validation.StringInSlice([]string{
										string(firewallpolicies.FirewallPolicyIntrusionDetectionProtocolICMP),
										string(firewallpolicies.FirewallPolicyIntrusionDetectionProtocolANY),
										string(firewallpolicies.FirewallPolicyIntrusionDetectionProtocolTCP),
										string(firewallpolicies.FirewallPolicyIntrusionDetectionProtocolUDP),
									}, true),
								},
								"source_addresses": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_addresses": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_ports": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"source_ip_groups": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"destination_ip_groups": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tls_certificate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_secret_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"insights": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"default_log_analytics_workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
					"retention_in_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"log_analytics_workspace": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: workspaces.ValidateWorkspaceID,
								},
								"firewall_location": commonschema.LocationWithoutForceNew(),
							},
						},
					},
				},
			},
		},

		"explicit_proxy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"http_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"https_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"enable_pac_file": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"pac_file_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 35536),
					},
					"pac_file": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sql_redirect_allowed": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"child_policies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"firewalls": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"rule_collection_groups": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_ip_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
					validation.IsIPv4Address,
				),
			},
		},

		"auto_learn_private_ranges_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}

	if !features.FourPointOhBeta() {
		resource["sku"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(firewallpolicies.FirewallPolicySkuTierPremium),
				string(firewallpolicies.FirewallPolicySkuTierStandard),
				string(firewallpolicies.FirewallPolicySkuTierBasic),
			}, false),
		}
	}

	return resource
}
