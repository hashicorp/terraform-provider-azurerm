package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var apiManagementCustomDomainResourceName = "azurerm_api_management_custom_domain"

func resourceArmApiManagementCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: apiManagementCustomDomainCreateUpdate,
		Read:   apiManagementCustomDomainRead,
		Update: apiManagementCustomDomainCreateUpdate,
		Delete: apiManagementCustomDomainDelete,
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
			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"management": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "proxy", "scm"},
				Elem: &schema.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"portal": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "proxy", "scm"},
				Elem: &schema.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"developer_portal": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "proxy", "scm"},
				Elem: &schema.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"proxy": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "proxy", "scm"},
				Elem: &schema.Resource{
					Schema: apiManagementResourceHostnameProxySchema(),
				},
			},
			"scm": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "proxy", "scm"},
				Elem: &schema.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
		},
	}
}

func apiManagementCustomDomainCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for API Management Custom domain creation.")

	name := d.Get("api_management_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("finding API Management (API Management %q / Resource Group %q): %s", name, resourceGroup, err)
	}

	if d.IsNewResource() {
		if existing.ServiceProperties.HostnameConfigurations != nil {
			return tf.ImportAsExistsError(apiManagementCustomDomainResourceName, *existing.ID)
		}
	}

	existing.ServiceProperties.HostnameConfigurations = expandAzureRmApiManagementHostnameConfigurations(d)

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, existing); err != nil {
		return fmt.Errorf("creating/updating Custom Domain (API Management %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Custom Domain (API Management %q / Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read ID for Custom Domain (API Management %q / Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return apiManagementCustomDomainRead(d, meta)
}

func apiManagementCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["service"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("API Management Service %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", resp.Name)

	if props := resp.ServiceProperties.HostnameConfigurations; props != nil {
		configs := flattenApiManagementHostnameConfiguration(resp.ServiceProperties.HostnameConfigurations, d)
		for _, config := range configs {
			for key, v := range config.(map[string]interface{}) {
				if err := d.Set(key, v); err != nil {
					return fmt.Errorf("setting `hostname_configuration` %q: %+v", key, err)
				}
			}
		}
	}

	return nil
}

func apiManagementCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["service"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("API Management Service %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	log.Printf("[DEBUG] Deleting API Management Custom domain (API Management %q / Resource Group %q)", name, resourceGroup)

	resp.ServiceProperties.HostnameConfigurations = nil

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, resp); err != nil {
		return fmt.Errorf("deleting Custom Domain (API Management %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func flattenApiManagementHostnameConfiguration(input *[]apimanagement.HostnameConfiguration, d *schema.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	managementResults := make([]interface{}, 0)
	portalResults := make([]interface{}, 0)
	developerPortalResults := make([]interface{}, 0)
	proxyResults := make([]interface{}, 0)
	scmResults := make([]interface{}, 0)

	for _, config := range *input {
		output := make(map[string]interface{})

		if config.HostName != nil {
			output["host_name"] = *config.HostName
		}

		if config.NegotiateClientCertificate != nil {
			output["negotiate_client_certificate"] = *config.NegotiateClientCertificate
		}

		if config.KeyVaultID != nil {
			output["key_vault_id"] = *config.KeyVaultID
		}

		// Iterate through old state to find sensitive props not returned by API.
		// This must be done in order to avoid state diffs.
		// NOTE: this information won't be available during times like Import, so this is a best-effort.
		snakeCaseConfigType := azure.ToSnakeCase(string(config.Type))
		if valsRaw, ok := d.GetOk(snakeCaseConfigType); ok {
			vals := valsRaw.([]interface{})
			for _, val := range vals {
				oldConfig := val.(map[string]interface{})

				if oldConfig["host_name"] == *config.HostName {
					output["certificate_password"] = oldConfig["certificate_password"]
					output["certificate"] = oldConfig["certificate"]
				}
			}
		}

		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagement.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			if config.DefaultSslBinding != nil {
				output["default_ssl_binding"] = *config.DefaultSslBinding
			}
			proxyResults = append(proxyResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypeManagement)):
			managementResults = append(managementResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypePortal)):
			portalResults = append(portalResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypeScm)):
			scmResults = append(scmResults, output)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"management":       managementResults,
			"portal":           portalResults,
			"developer_portal": developerPortalResults,
			"proxy":            proxyResults,
			"scm":              scmResults,
		},
	}
}
