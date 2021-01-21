package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApplicationGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationGatewayRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceApplicationGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["applicationGateways"]

	applicationGateway, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(applicationGateway.Response) {
			log.Printf("[DEBUG] Application Gateway %q was not found in Resource Group %q - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Application Gateway %s: %+v", name, err)
	}

	if location := applicationGateway.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("zones", applicationGateway.Zones)

	identity, err := flattenRmApplicationGatewayIdentity(applicationGateway.Identity)
	if err != nil {
		return err
	}
	if err = d.Set("identity", identity); err != nil {
		return err
	}

	if props := applicationGateway.ApplicationGatewayPropertiesFormat; props != nil {
		if err = d.Set("authentication_certificate", flattenApplicationGatewayAuthenticationCertificates(props.AuthenticationCertificates, d)); err != nil {
			return fmt.Errorf("Error setting `authentication_certificate`: %+v", err)
		}

		if err = d.Set("trusted_root_certificate", flattenApplicationGatewayTrustedRootCertificates(props.TrustedRootCertificates, d)); err != nil {
			return fmt.Errorf("Error setting `trusted_root_certificate`: %+v", err)
		}

		if setErr := d.Set("backend_address_pool", flattenApplicationGatewayBackendAddressPools(props.BackendAddressPools)); setErr != nil {
			return fmt.Errorf("Error setting `backend_address_pool`: %+v", setErr)
		}

		backendHttpSettings, err := flattenApplicationGatewayBackendHTTPSettings(props.BackendHTTPSettingsCollection)
		if err != nil {
			return fmt.Errorf("Error flattening `backend_http_settings`: %+v", err)
		}
		if setErr := d.Set("backend_http_settings", backendHttpSettings); setErr != nil {
			return fmt.Errorf("Error setting `backend_http_settings`: %+v", setErr)
		}

		if setErr := d.Set("ssl_policy", flattenApplicationGatewaySslPolicy(props.SslPolicy)); setErr != nil {
			return fmt.Errorf("Error setting `ssl_policy`: %+v", setErr)
		}

		d.Set("enable_http2", props.EnableHTTP2)

		httpListeners, err := flattenApplicationGatewayHTTPListeners(props.HTTPListeners)
		if err != nil {
			return fmt.Errorf("Error flattening `http_listener`: %+v", err)
		}
		if setErr := d.Set("http_listener", httpListeners); setErr != nil {
			return fmt.Errorf("Error setting `http_listener`: %+v", setErr)
		}

		if setErr := d.Set("frontend_port", flattenApplicationGatewayFrontendPorts(props.FrontendPorts)); setErr != nil {
			return fmt.Errorf("Error setting `frontend_port`: %+v", setErr)
		}

		if setErr := d.Set("frontend_ip_configuration", flattenApplicationGatewayFrontendIPConfigurations(props.FrontendIPConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `frontend_ip_configuration`: %+v", setErr)
		}

		if setErr := d.Set("gateway_ip_configuration", flattenApplicationGatewayIPConfigurations(props.GatewayIPConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `gateway_ip_configuration`: %+v", setErr)
		}

		if setErr := d.Set("probe", flattenApplicationGatewayProbes(props.Probes)); setErr != nil {
			return fmt.Errorf("Error setting `probe`: %+v", setErr)
		}

		requestRoutingRules, err := flattenApplicationGatewayRequestRoutingRules(props.RequestRoutingRules)
		if err != nil {
			return fmt.Errorf("Error flattening `request_routing_rule`: %+v", err)
		}
		if setErr := d.Set("request_routing_rule", requestRoutingRules); setErr != nil {
			return fmt.Errorf("Error setting `request_routing_rule`: %+v", setErr)
		}

		redirectConfigurations, err := flattenApplicationGatewayRedirectConfigurations(props.RedirectConfigurations)
		if err != nil {
			return fmt.Errorf("Error flattening `redirect configuration`: %+v", err)
		}
		if setErr := d.Set("redirect_configuration", redirectConfigurations); setErr != nil {
			return fmt.Errorf("Error setting `redirect configuration`: %+v", setErr)
		}

		rewriteRuleSets := flattenApplicationGatewayRewriteRuleSets(props.RewriteRuleSets)
		if setErr := d.Set("rewrite_rule_set", rewriteRuleSets); setErr != nil {
			return fmt.Errorf("Error setting `rewrite_rule_set`: %+v", setErr)
		}

		if setErr := d.Set("sku", flattenApplicationGatewaySku(props.Sku)); setErr != nil {
			return fmt.Errorf("Error setting `sku`: %+v", setErr)
		}

		if setErr := d.Set("autoscale_configuration", flattenApplicationGatewayAutoscaleConfiguration(props.AutoscaleConfiguration)); setErr != nil {
			return fmt.Errorf("Error setting `autoscale_configuration`: %+v", setErr)
		}

		if setErr := d.Set("ssl_certificate", flattenApplicationGatewaySslCertificates(props.SslCertificates, d)); setErr != nil {
			return fmt.Errorf("Error setting `ssl_certificate`: %+v", setErr)
		}

		if setErr := d.Set("custom_error_configuration", flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `custom_error_configuration`: %+v", setErr)
		}

		urlPathMaps, err := flattenApplicationGatewayURLPathMaps(props.URLPathMaps)
		if err != nil {
			return fmt.Errorf("Error flattening `url_path_map`: %+v", err)
		}
		if setErr := d.Set("url_path_map", urlPathMaps); setErr != nil {
			return fmt.Errorf("Error setting `url_path_map`: %+v", setErr)
		}

		if setErr := d.Set("waf_configuration", flattenApplicationGatewayWafConfig(props.WebApplicationFirewallConfiguration)); setErr != nil {
			return fmt.Errorf("Error setting `waf_configuration`: %+v", setErr)
		}

		if props.FirewallPolicy != nil {
			d.Set("firewall_policy_id", props.FirewallPolicy.ID)
		}
	}

	return tags.FlattenAndSet(d, applicationGateway.Tags)
}
