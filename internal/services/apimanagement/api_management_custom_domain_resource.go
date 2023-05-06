package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var apiManagementCustomDomainResourceName = "azurerm_api_management_custom_domain"

func resourceApiManagementCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: apiManagementCustomDomainCreateUpdate,
		Read:   apiManagementCustomDomainRead,
		Update: apiManagementCustomDomainCreateUpdate,
		Delete: apiManagementCustomDomainDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CustomDomainID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementID,
			},

			"management": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
				Elem: &pluginsdk.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"portal": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
				Elem: &pluginsdk.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"developer_portal": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
				Elem: &pluginsdk.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"scm": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
				Elem: &pluginsdk.Resource{
					Schema: apiManagementResourceHostnameSchema(),
				},
			},
			"gateway": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
				Elem: &pluginsdk.Resource{
					Schema: apiManagementResourceHostnameProxySchema(),
				},
			},
		},
	}
}

func apiManagementCustomDomainCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiMgmtId, err := parse.ApiManagementID(d.Get("api_management_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewCustomDomainID(apiMgmtId.SubscriptionId, apiMgmtId.ResourceGroup, apiMgmtId.ServiceName, "default")

	existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if d.IsNewResource() {
		if existing.ServiceProperties != nil && existing.ServiceProperties.HostnameConfigurations != nil && len(*existing.ServiceProperties.HostnameConfigurations) > 1 {
			return tf.ImportAsExistsError(apiManagementCustomDomainResourceName, *existing.ID)
		}
	}

	existing.ServiceProperties.HostnameConfigurations = expandApiManagementCustomDomains(d)

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Updating", "Unknown"},
		Target:                    []string{"Succeeded", "Ready"},
		Refresh:                   apiManagementRefreshFunc(ctx, client, id.ServiceName, id.ResourceGroup),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 6,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}

	// The API expects user assigned identities to be submitted with nil values
	if existing.Identity != nil {
		for k, v := range existing.Identity.UserAssignedIdentities {
			if v == nil {
				continue
			}
			existing.Identity.UserAssignedIdentities[k].ClientID = nil
			existing.Identity.UserAssignedIdentities[k].PrincipalID = nil
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, existing)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", id)
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}
	d.SetId(id.ID())

	return apiManagementCustomDomainRead(d, meta)
}

func apiManagementCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	environment := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimHostNameSuffix, ok := environment.ApiManagement.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine API Management domain suffix for environment %q", environment.Name)
	}

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	apiMgmtId := parse.NewApiManagementID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_management_id", apiMgmtId.ID())

	if resp.ServiceProperties != nil && resp.ServiceProperties.HostnameConfigurations != nil {
		configs := flattenApiManagementHostnameConfiguration(resp.ServiceProperties.HostnameConfigurations, d, *resp.Name, *apimHostNameSuffix)
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

func apiManagementCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Updating", "Unknown"},
		Target:                    []string{"Succeeded", "Ready"},
		Refresh:                   apiManagementRefreshFunc(ctx, client, id.ServiceName, id.ResourceGroup),
		MinTimeout:                1 * time.Minute,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 6,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", *id, err)
	}

	log.Printf("[DEBUG] Deleting %s", *id)

	resp.ServiceProperties.HostnameConfigurations = nil

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, resp)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", *id)
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", *id, err)
	}

	return nil
}

func expandApiManagementCustomDomains(input *pluginsdk.ResourceData) *[]apimanagement.HostnameConfiguration {
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

	if gatewayRawVal, ok := input.GetOk("gateway"); ok {
		vs := gatewayRawVal.([]interface{})
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

func flattenApiManagementHostnameConfiguration(input *[]apimanagement.HostnameConfiguration, d *pluginsdk.ResourceData, name, apimHostNameSuffix string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	managementResults := make([]interface{}, 0)
	portalResults := make([]interface{}, 0)
	developerPortalResults := make([]interface{}, 0)
	gatewayResults := make([]interface{}, 0)
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

		if config.IdentityClientID != nil {
			output["ssl_keyvault_identity_client_id"] = *config.IdentityClientID
		}

		var configType string
		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagement.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			if config.DefaultSslBinding != nil {
				output["default_ssl_binding"] = *config.DefaultSslBinding
			}
			gatewayResults = append(gatewayResults, output)
			configType = "gateway"

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

	res := map[string]interface{}{
		"management":       managementResults,
		"portal":           portalResults,
		"developer_portal": developerPortalResults,
		"scm":              scmResults,
		"gateway":          gatewayResults,
	}

	return []interface{}{res}
}
