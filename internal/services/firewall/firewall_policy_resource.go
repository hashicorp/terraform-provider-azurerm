package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	logAnalytiscValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const azureFirewallPolicyResourceName = "azurerm_firewall_policy"

func resourceFirewallPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallPolicyCreateUpdate,
		Read:   resourceFirewallPolicyRead,
		Update: resourceFirewallPolicyCreateUpdate,
		Delete: resourceFirewallPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallPolicyID(id)
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
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall_policy", *resp.ID)
		}
	}

	expandedIdentity, err := expandFirewallPolicyIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	props := network.FirewallPolicy{
		FirewallPolicyPropertiesFormat: &network.FirewallPolicyPropertiesFormat{
			ThreatIntelMode:      network.AzureFirewallThreatIntelMode(d.Get("threat_intelligence_mode").(string)),
			ThreatIntelWhitelist: expandFirewallPolicyThreatIntelWhitelist(d.Get("threat_intelligence_allowlist").([]interface{})),
			DNSSettings:          expandFirewallPolicyDNSSetting(d.Get("dns").([]interface{})),
			IntrusionDetection:   expandFirewallPolicyIntrusionDetection(d.Get("intrusion_detection").([]interface{})),
			TransportSecurity:    expandFirewallPolicyTransportSecurity(d.Get("tls_certificate").([]interface{})),
			Insights:             expandFirewallPolicyInsights(d.Get("insights").([]interface{})),
		},
		Identity: expandedIdentity,
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if id, ok := d.GetOk("base_policy_id"); ok {
		props.FirewallPolicyPropertiesFormat.BasePolicy = &network.SubResource{ID: utils.String(id.(string))}
	}

	if v, ok := d.GetOk("sku"); ok {
		props.FirewallPolicyPropertiesFormat.Sku = &network.FirewallPolicySku{
			Tier: network.FirewallPolicySkuTier(v.(string)),
		}
	}

	if v, ok := d.GetOk("private_ip_ranges"); ok {
		privateIPRanges := utils.ExpandStringSlice(v.([]interface{}))
		props.FirewallPolicyPropertiesFormat.Snat = &network.FirewallPolicySNAT{
			PrivateRanges: privateIPRanges,
		}
	}

	locks.ByName(id.Name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.Name, azureFirewallPolicyResourceName)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceFirewallPolicyRead(d, meta)
}

func resourceFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Firewall Policy %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if prop := resp.FirewallPolicyPropertiesFormat; prop != nil {
		basePolicyID := ""
		if resp.BasePolicy != nil && resp.BasePolicy.ID != nil {
			basePolicyID = *resp.BasePolicy.ID
		}
		d.Set("base_policy_id", basePolicyID)

		d.Set("threat_intelligence_mode", string(prop.ThreatIntelMode))

		if sku := prop.Sku; sku != nil {
			d.Set("sku", string(sku.Tier))
		}

		if err := d.Set("threat_intelligence_allowlist", flattenFirewallPolicyThreatIntelWhitelist(resp.ThreatIntelWhitelist)); err != nil {
			return fmt.Errorf(`setting "threat_intelligence_allowlist": %+v`, err)
		}

		if err := d.Set("dns", flattenFirewallPolicyDNSSetting(prop.DNSSettings)); err != nil {
			return fmt.Errorf(`setting "dns": %+v`, err)
		}

		if err := d.Set("intrusion_detection", flattenFirewallPolicyIntrusionDetection(resp.IntrusionDetection)); err != nil {
			return fmt.Errorf(`setting "intrusion_detection": %+v`, err)
		}

		if err := d.Set("tls_certificate", flattenFirewallPolicyTransportSecurity(prop.TransportSecurity)); err != nil {
			return fmt.Errorf(`setting "tls_certificate": %+v`, err)
		}

		if err := d.Set("child_policies", flattenNetworkSubResourceID(prop.ChildPolicies)); err != nil {
			return fmt.Errorf(`setting "child_policies": %+v`, err)
		}

		if err := d.Set("firewalls", flattenNetworkSubResourceID(prop.Firewalls)); err != nil {
			return fmt.Errorf(`setting "firewalls": %+v`, err)
		}

		if err := d.Set("rule_collection_groups", flattenNetworkSubResourceID(prop.RuleCollectionGroups)); err != nil {
			return fmt.Errorf(`setting "rule_collection_groups": %+v`, err)
		}

		var privateIPRanges []interface{}
		if prop.Snat != nil {
			privateIPRanges = utils.FlattenStringSlice(prop.Snat.PrivateRanges)
		}
		if err := d.Set("private_ip_ranges", privateIPRanges); err != nil {
			return fmt.Errorf("setting `private_ip_ranges`: %+v", err)
		}

		if err := d.Set("insights", flattenFirewallPolicyInsights(prop.Insights)); err != nil {
			return fmt.Errorf(`setting "insights": %+v`, err)
		}
	}

	flattenedIdentity, err := flattenFirewallPolicyIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.Name, azureFirewallPolicyResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandFirewallPolicyThreatIntelWhitelist(input []interface{}) *network.FirewallPolicyThreatIntelWhitelist {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &network.FirewallPolicyThreatIntelWhitelist{
		IPAddresses: utils.ExpandStringSlice(raw["ip_addresses"].(*pluginsdk.Set).List()),
		Fqdns:       utils.ExpandStringSlice(raw["fqdns"].(*pluginsdk.Set).List()),
	}

	return output
}

func expandFirewallPolicyDNSSetting(input []interface{}) *network.DNSSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &network.DNSSettings{
		Servers:     utils.ExpandStringSlice(raw["servers"].([]interface{})),
		EnableProxy: utils.Bool(raw["proxy_enabled"].(bool)),
	}

	return output
}

func expandFirewallPolicyIntrusionDetection(input []interface{}) *network.FirewallPolicyIntrusionDetection {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	var signatureOverrides []network.FirewallPolicyIntrusionDetectionSignatureSpecification
	for _, v := range raw["signature_overrides"].([]interface{}) {
		overrides := v.(map[string]interface{})
		signatureOverrides = append(signatureOverrides, network.FirewallPolicyIntrusionDetectionSignatureSpecification{
			ID:   utils.String(overrides["id"].(string)),
			Mode: network.FirewallPolicyIntrusionDetectionStateType(overrides["state"].(string)),
		})
	}

	var trafficBypass []network.FirewallPolicyIntrusionDetectionBypassTrafficSpecifications

	for _, v := range raw["traffic_bypass"].([]interface{}) {
		bypass := v.(map[string]interface{})
		trafficBypass = append(trafficBypass, network.FirewallPolicyIntrusionDetectionBypassTrafficSpecifications{
			Name:                 utils.String(bypass["name"].(string)),
			Description:          utils.String(bypass["description"].(string)),
			Protocol:             network.FirewallPolicyIntrusionDetectionProtocol(bypass["protocol"].(string)),
			SourceAddresses:      utils.ExpandStringSlice(bypass["source_addresses"].(*pluginsdk.Set).List()),
			DestinationAddresses: utils.ExpandStringSlice(bypass["destination_addresses"].(*pluginsdk.Set).List()),
			DestinationPorts:     utils.ExpandStringSlice(bypass["destination_ports"].(*pluginsdk.Set).List()),
			SourceIPGroups:       utils.ExpandStringSlice(bypass["source_ip_groups"].(*pluginsdk.Set).List()),
			DestinationIPGroups:  utils.ExpandStringSlice(bypass["destination_ip_groups"].(*pluginsdk.Set).List()),
		})
	}

	return &network.FirewallPolicyIntrusionDetection{
		Mode: network.FirewallPolicyIntrusionDetectionStateType(raw["mode"].(string)),
		Configuration: &network.FirewallPolicyIntrusionDetectionConfiguration{
			SignatureOverrides:    &signatureOverrides,
			BypassTrafficSettings: &trafficBypass,
		},
	}
}

func expandFirewallPolicyTransportSecurity(input []interface{}) *network.FirewallPolicyTransportSecurity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &network.FirewallPolicyTransportSecurity{
		CertificateAuthority: &network.FirewallPolicyCertificateAuthority{
			KeyVaultSecretID: utils.String(raw["key_vault_secret_id"].(string)),
			Name:             utils.String(raw["name"].(string)),
		},
	}
}

func expandFirewallPolicyIdentity(input []interface{}) (*network.ManagedServiceIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	if expanded.Type == identity.TypeNone {
		return nil, nil
	}

	out := network.ManagedServiceIdentity{
		PrincipalID:            nil,
		TenantID:               nil,
		Type:                   network.ResourceIdentityType(string(expanded.Type)),
		UserAssignedIdentities: nil,
	}
	if expanded.Type == identity.TypeUserAssigned {
		out.UserAssignedIdentities = make(map[string]*network.ManagedServiceIdentityUserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &network.ManagedServiceIdentityUserAssignedIdentitiesValue{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func expandFirewallPolicyInsights(input []interface{}) *network.FirewallPolicyInsights {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &network.FirewallPolicyInsights{
		IsEnabled:             utils.Bool(raw["enabled"].(bool)),
		RetentionDays:         utils.Int32(int32(raw["retention_in_days"].(int))),
		LogAnalyticsResources: expandFirewallPolicyLogAnalyticsResources(raw["default_log_analytics_workspace_id"].(string), raw["log_analytics_workspace"].([]interface{})),
	}

	return output
}

func expandFirewallPolicyLogAnalyticsResources(defaultWorkspaceId string, workspaces []interface{}) *network.FirewallPolicyLogAnalyticsResources {
	output := &network.FirewallPolicyLogAnalyticsResources{
		DefaultWorkspaceID: &network.SubResource{
			ID: &defaultWorkspaceId,
		},
	}

	var workspaceList []network.FirewallPolicyLogAnalyticsWorkspace
	for _, workspace := range workspaces {
		workspace := workspace.(map[string]interface{})
		workspaceList = append(workspaceList, network.FirewallPolicyLogAnalyticsWorkspace{
			Region: utils.String(location.Normalize(workspace["firewall_location"].(string))),
			WorkspaceID: &network.SubResource{
				ID: utils.String(workspace["id"].(string)),
			},
		})
	}
	if workspaceList != nil {
		output.Workspaces = &workspaceList
	}

	return output
}

func flattenFirewallPolicyThreatIntelWhitelist(input *network.FirewallPolicyThreatIntelWhitelist) []interface{} {
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

func flattenFirewallPolicyDNSSetting(input *network.DNSSettings) []interface{} {
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

func flattenFirewallPolicyIntrusionDetection(input *network.FirewallPolicyIntrusionDetection) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	signatureOverrides := make([]interface{}, 0)
	trafficBypass := make([]interface{}, 0)

	if input.Configuration == nil {
		return []interface{}{
			map[string]interface{}{
				"mode":                string(input.Mode),
				"signature_overrides": signatureOverrides,
				"traffic_bypass":      trafficBypass,
			},
		}
	}

	if overrides := input.Configuration.SignatureOverrides; overrides != nil {
		for _, override := range *overrides {
			id := ""
			if override.ID != nil {
				id = *override.ID
			}
			signatureOverrides = append(signatureOverrides, map[string]interface{}{
				"id":    id,
				"state": string(override.Mode),
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

			sourceAddresses := make([]string, 0)
			if bypass.SourceAddresses != nil {
				sourceAddresses = *bypass.SourceAddresses
			}

			destinationAddresses := make([]string, 0)
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
				"protocol":              string(bypass.Protocol),
				"source_addresses":      sourceAddresses,
				"destination_addresses": destinationAddresses,
				"destination_ports":     destinationPorts,
				"source_ip_groups":      sourceIPGroups,
				"destination_ip_groups": destinationIPGroups,
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                string(input.Mode),
			"signature_overrides": signatureOverrides,
			"traffic_bypass":      trafficBypass,
		},
	}
}

func flattenFirewallPolicyTransportSecurity(input *network.FirewallPolicyTransportSecurity) []interface{} {
	if input == nil || input.CertificateAuthority == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_secret_id": input.CertificateAuthority.KeyVaultSecretID,
			"name":                input.CertificateAuthority.Name,
		},
	}
}

func flattenFirewallPolicyIdentity(input *network.ManagedServiceIdentity) (*[]interface{}, error) {
	var transition *identity.UserAssignedMap

	if input != nil {
		transition = &identity.UserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k, v := range input.UserAssignedIdentities {
			transition.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenUserAssignedMap(transition)
}

func flattenFirewallPolicyInsights(input *network.FirewallPolicyInsights) []interface{} {
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

func flattenFirewallPolicyLogAnalyticsResources(input *network.FirewallPolicyLogAnalyticsResources) (string, []interface{}) {
	if input == nil {
		return "", []interface{}{}
	}

	var defaultLogAnalyticsWorkspaceId string
	if input.DefaultWorkspaceID != nil && input.DefaultWorkspaceID.ID != nil {
		defaultLogAnalyticsWorkspaceId = *input.DefaultWorkspaceID.ID
	}

	var workspaceList []interface{}
	if input.Workspaces != nil {
		for _, workspace := range *input.Workspaces {
			loc := location.NormalizeNilable(workspace.Region)

			var id string
			if workspace.WorkspaceID != nil && workspace.WorkspaceID.ID != nil {
				id = *workspace.WorkspaceID.ID
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
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallPolicyName(),
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.FirewallPolicySkuTierPremium),
				string(network.FirewallPolicySkuTierStandard),
				string(network.FirewallPolicySkuTierBasic),
			}, false),
		},

		"location": commonschema.Location(),

		"base_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.FirewallPolicyID,
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
			Default:  string(network.AzureFirewallThreatIntelModeAlert),
			ValidateFunc: validation.StringInSlice([]string{
				string(network.AzureFirewallThreatIntelModeAlert),
				string(network.AzureFirewallThreatIntelModeDeny),
				string(network.AzureFirewallThreatIntelModeOff),
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
							string(network.FirewallPolicyIntrusionDetectionStateTypeOff),
							string(network.FirewallPolicyIntrusionDetectionStateTypeAlert),
							string(network.FirewallPolicyIntrusionDetectionStateTypeDeny),
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
										string(network.FirewallPolicyIntrusionDetectionStateTypeOff),
										string(network.FirewallPolicyIntrusionDetectionStateTypeAlert),
										string(network.FirewallPolicyIntrusionDetectionStateTypeDeny),
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
									ValidateFunc: validation.StringInSlice([]string{
										string(network.FirewallPolicyIntrusionDetectionProtocolICMP),
										string(network.FirewallPolicyIntrusionDetectionProtocolANY),
										string(network.FirewallPolicyIntrusionDetectionProtocolTCP),
										string(network.FirewallPolicyIntrusionDetectionProtocolUDP),
									}, false),
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

		"identity": commonschema.UserAssignedIdentityOptional(),

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
						ValidateFunc: logAnalytiscValidate.LogAnalyticsWorkspaceID,
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
									ValidateFunc: logAnalytiscValidate.LogAnalyticsWorkspaceID,
								},
								"firewall_location": commonschema.LocationWithoutForceNew(),
							},
						},
					},
				},
			},
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

		"tags": tags.SchemaEnforceLowerCaseKeys(),
	}
}
