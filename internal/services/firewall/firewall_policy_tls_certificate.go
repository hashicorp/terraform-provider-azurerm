package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const azureFirewallPolicyTLSCertificateResourceName = "azurerm_firewall_policy_tls_certificate"

func resourceFirewallPolicyTLSCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallPolicyTLSCertificateCreateUpdate,
		Read:   resourceFirewallPolicyTLSCertificateRead,
		Update: resourceFirewallPolicyTLSCertificateCreateUpdate,
		Delete: resourceFirewallPolicyTLSCertificateDelete,

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

		Schema: map[string]*pluginsdk.Schema{
			"firewall_policy_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallPolicyName(),
			},

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
		},
	}
}

func resourceFirewallPolicyTLSCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	firewall_policy_id, err := parse.FirewallPolicyID(d.Get("firewall_policy_id").(string))

	resp, err := client.Get(ctx, firewall_policy_id.ResourceGroup, firewall_policy_id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing Firewall Policy %q (Resource Group %q): %+v",
				firewall_policy_id.Name, firewall_policy_id.ResourceGroup, err)
		}
	}

	resp.TransportSecurity = expandFirewallPolicyTransportSecurity(d.Get("tls_certificate").([]interface{}))

	d.SetId(*resp.ID)

	return resourceFirewallPolicyTLSCertificateRead(d, meta)
}

func resourceFirewallPolicyTLSCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	}

	if err := d.Set("identity", flattenFirewallPolicyIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("flattening identity on Firewall Policy %q (Resource Group %q): %+v",
			id.Name, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceFirewallPolicyTLSCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, azureFirewallPolicyTLSCertificateResourceName)
	defer locks.UnlockByName(id.Name, azureFirewallPolicyTLSCertificateResourceName)

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

func flattenFirewallPolicyTransportSecurity(input *network.FirewallPolicyTransportSecurity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_secret_id": input.CertificateAuthority.KeyVaultSecretID,
			"name":                input.CertificateAuthority.Name,
		},
	}
}
