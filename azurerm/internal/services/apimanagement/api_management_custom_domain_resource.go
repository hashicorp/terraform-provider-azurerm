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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
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
			"api_management_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

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

	apiManagementID := d.Get("api_management_id").(string)
	id, err := azure.ParseAzureResourceID(apiManagementID)
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]

	existing, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return fmt.Errorf("finding API Management (API Management %q / Resource Group %q): %s", serviceName, resourceGroup, err)
	}

	if d.IsNewResource() {
		if existing.ServiceProperties != nil && existing.ServiceProperties.HostnameConfigurations != nil {
			return tf.ImportAsExistsError(apiManagementCustomDomainResourceName, *existing.ID)
		}
	}

	existing.ServiceProperties.HostnameConfigurations = expandApiManagementCustomDomains(d)

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, existing); err != nil {
		return fmt.Errorf("creating/updating Custom Domain (API Management %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return fmt.Errorf("retrieving Custom Domain (API Management %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read ID for Custom Domain (API Management %q / Resource Group %q)", serviceName, resourceGroup)
	}

	customDomainsID := fmt.Sprintf("%s/customDomains/default", *read.ID)
	d.SetId(customDomainsID)

	return apiManagementCustomDomainRead(d, meta)
}

func apiManagementCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiManagementCustomDomainID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	resp, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("API Management Service %q was not found in Resource Group %q - removing from state!", serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on API Management Service %q (Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_id", resp.ID)

	if resp.ServiceProperties != nil && resp.ServiceProperties.HostnameConfigurations != nil {
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

	id, err := parse.ApiManagementCustomDomainID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	resp, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("API Management Service %q was not found in Resource Group %q - removing from state!", serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on API Management Service %q (Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	log.Printf("[DEBUG] Deleting API Management Custom domain (API Management %q / Resource Group %q)", serviceName, resourceGroup)

	resp.ServiceProperties.HostnameConfigurations = nil

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, resp); err != nil {
		return fmt.Errorf("deleting Custom Domain (API Management %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	return nil
}

func expandApiManagementCustomDomains(input *schema.ResourceData) *[]apimanagement.HostnameConfiguration {
	results := make([]apimanagement.HostnameConfiguration, 0)

	if managementRawVal, ok := input.GetOk("management"); ok {
		vs := managementRawVal.([]interface{})
		for _, rawVal := range vs {
			v := rawVal.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.HostnameTypeManagement)
			results = append(results, output)
		}
	}
	if portalRawVal, ok := input.GetOk("portal"); ok {
		vs := portalRawVal.([]interface{})
		for _, rawVal := range vs {
			v := rawVal.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.HostnameTypePortal)
			results = append(results, output)
		}
	}
	if developerPortalRawVal, ok := input.GetOk("developer_portal"); ok {
		vs := developerPortalRawVal.([]interface{})
		for _, rawVal := range vs {
			v := rawVal.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.HostnameTypeDeveloperPortal)
			results = append(results, output)
		}
	}
	if proxyRawVal, ok := input.GetOk("proxy"); ok {
		vs := proxyRawVal.([]interface{})
		for _, rawVal := range vs {
			v := rawVal.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.HostnameTypeProxy)
			if value, ok := v["default_ssl_binding"]; ok {
				output.DefaultSslBinding = utils.Bool(value.(bool))
			}
			results = append(results, output)
		}
	}
	if scmRawVal, ok := input.GetOk("scm"); ok {
		vs := scmRawVal.([]interface{})
		for _, rawVal := range vs {
			v := rawVal.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.HostnameTypeScm)
			results = append(results, output)
		}
	}
	return &results
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

		snakeCaseConfigType := azure.ToSnakeCase(string(config.Type))
		if valsRaw, ok := d.GetOk(snakeCaseConfigType); ok {
			vals := valsRaw.([]interface{})
			azure.CopyCertificateAndPassword(vals, *config.HostName, output)
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
