package migration

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
)

func CustomHttpsConfigurationV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"frontend_endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"custom_https_provisioning_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"custom_https_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"minimum_tls_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provisioning_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provisioning_substate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_key_vault_certificate_secret_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"azure_key_vault_certificate_secret_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"azure_key_vault_certificate_vault_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func CustomHttpsConfigurationV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// this was: fmt.Sprintf("%s/customHttpsConfiguration/%s", frontEndEndpointId, frontendEndpointName
	oldId := rawState["id"].(string)
	oldParsedId, err := azure.ParseAzureResourceID(oldId)
	if err != nil {
		return rawState, err
	}

	resourceGroup := oldParsedId.ResourceGroup
	frontdoorName := ""
	frontendEndpointName := ""
	for key, value := range oldParsedId.Path {
		if strings.EqualFold(key, "frontdoors") {
			frontdoorName = value
			continue
		}

		if strings.EqualFold(key, "frontendEndpoints") {
			frontendEndpointName = value
			continue
		}
	}

	if frontdoorName == "" {
		return rawState, fmt.Errorf("couldn't find the `frontdoors` segment in the old resource id %q", oldId)
	}

	if frontendEndpointName == "" {
		return rawState, fmt.Errorf("couldn't find the `frontendEndpoints` segment in the old resource id %q", oldId)
	}

	newId := parse.NewFrontendEndpointID(parse.NewFrontDoorID(oldParsedId.SubscriptionID, resourceGroup, frontdoorName), frontendEndpointName)
	newIdStr := newId.ID("")

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

	rawState["id"] = newIdStr

	return rawState, nil
}
