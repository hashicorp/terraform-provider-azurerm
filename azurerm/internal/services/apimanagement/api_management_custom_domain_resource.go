package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var apiManagementCustomDomainResourceName = "azurerm_api_management_custom_domain"

func resourceApiManagementCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: apiManagementCustomDomainCreateUpdate,
		Read:   apiManagementCustomDomainRead,
		Update: apiManagementCustomDomainCreateUpdate,
		Delete: apiManagementCustomDomainDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
	id, err := parse.ApiManagementID(apiManagementID)
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	existing, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return fmt.Errorf("finding API Management (API Management %q / Resource Group %q): %s", serviceName, resourceGroup, err)
	}

	if d.IsNewResource() {
		if existing.ServiceProperties != nil && existing.ServiceProperties.HostnameConfigurations != nil && len(*existing.ServiceProperties.HostnameConfigurations) > 1 {
			return tf.ImportAsExistsError(apiManagementCustomDomainResourceName, *existing.ID)
		}
	}

	existing.ServiceProperties.HostnameConfigurations = expandApiManagementCustomDomains(d)

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for API Management Service %q (Resource Group: %q) to become ready", serviceName, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Updating", "Unknown"},
		Target:                    []string{"Succeeded", "Ready"},
		Refresh:                   apiManagementRefreshFunc(ctx, client, serviceName, resourceGroup),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 6,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for API Management Service %q (Resource Group: %q) to become ready: %+v", serviceName, resourceGroup, err)
	}

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

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for API Management Service %q (Resource Group: %q) to become ready", serviceName, resourceGroup)
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for API Management Service %q (Resource Group: %q) to become ready: %+v", serviceName, resourceGroup, err)
	}

	customDomainsID := fmt.Sprintf("%s/customDomains/default", *read.ID)
	d.SetId(customDomainsID)

	return apiManagementCustomDomainRead(d, meta)
}

func apiManagementCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	environment := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
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

	d.Set("api_management_id", resp.ID)

	if resp.ServiceProperties != nil && resp.ServiceProperties.HostnameConfigurations != nil {
		apimHostNameSuffix := environment.APIManagementHostNameSuffix
		configs := flattenApiManagementHostnameConfiguration(resp.ServiceProperties.HostnameConfigurations, d, *resp.Name, apimHostNameSuffix)
		for _, config := range configs {
			for key, v := range config.(map[string]interface{}) {
				// lintignore:R001
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

	id, err := parse.CustomDomainID(d.Id())
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

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for API Management Service %q (Resource Group: %q) to become ready", serviceName, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Updating", "Unknown"},
		Target:                    []string{"Succeeded", "Ready"},
		Refresh:                   apiManagementRefreshFunc(ctx, client, serviceName, resourceGroup),
		MinTimeout:                1 * time.Minute,
		Timeout:                   d.Timeout(schema.TimeoutDelete),
		ContinuousTargetOccurence: 6,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for API Management Service %q (Resource Group: %q) to become ready: %+v", serviceName, resourceGroup, err)
	}

	log.Printf("[DEBUG] Deleting API Management Custom Domain (API Management %q / Resource Group %q)", serviceName, resourceGroup)

	resp.ServiceProperties.HostnameConfigurations = nil

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, resp); err != nil {
		return fmt.Errorf("deleting Custom Domain (API Management %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for API Management Service %q (Resource Group: %q) to become ready", serviceName, resourceGroup)
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for API Management Service %q (Resource Group: %q) to become ready: %+v", serviceName, resourceGroup, err)
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

func flattenApiManagementHostnameConfiguration(input *[]apimanagement.HostnameConfiguration, d *schema.ResourceData, name, apimHostNameSuffix string) []interface{} {
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

		// There'll always be a default custom domain with hostName "apim_name.azure-api.net" and Type "Proxy", which should be ignored
		if *config.HostName == strings.ToLower(name)+"."+apimHostNameSuffix && config.Type == apimanagement.HostnameTypeProxy {
			continue
		}

		if config.HostName != nil {
			output["host_name"] = *config.HostName
		}

		if config.NegotiateClientCertificate != nil {
			output["negotiate_client_certificate"] = *config.NegotiateClientCertificate
		}

		if config.KeyVaultID != nil {
			output["key_vault_id"] = *config.KeyVaultID
		}

		var configType string
		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagement.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			if config.DefaultSslBinding != nil {
				output["default_ssl_binding"] = *config.DefaultSslBinding
			}
			proxyResults = append(proxyResults, output)
			configType = "proxy"

		case strings.ToLower(string(apimanagement.HostnameTypeManagement)):
			managementResults = append(managementResults, output)
			configType = "management"

		case strings.ToLower(string(apimanagement.HostnameTypePortal)):
			portalResults = append(portalResults, output)
			configType = "portal"

		case strings.ToLower(string(apimanagement.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)
			configType = "developer_portal"

		case strings.ToLower(string(apimanagement.HostnameTypeScm)):
			scmResults = append(scmResults, output)
			configType = "scm"
		}

		if configType != "" {
			if valsRaw, ok := d.GetOk(configType); ok {
				vals := valsRaw.([]interface{})
				schemaz.CopyCertificateAndPassword(vals, *config.HostName, output)
			}
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
