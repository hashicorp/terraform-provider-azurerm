// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_custom_domain

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func apiManagementCustomDomainRead(d *pluginsdk.ResourceData, meta any) error {
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

	apiMgmtId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)

	resp, err := client.Get(ctx, apiMgmtId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_management_id", apiMgmtId.ID())

	if resp.Model != nil && resp.Model.Properties.HostnameConfigurations != nil {
		configs := flattenApiManagementHostnameConfiguration(resp.Model.Properties.HostnameConfigurations, d, *resp.Model.Name, *apimHostNameSuffix)
		for _, config := range configs {
			for key, v := range config.(map[string]any) {
				// lintignore:R001
				if err := d.Set(key, v); err != nil {
					return fmt.Errorf("setting `hostname_configuration` %q: %+v", key, err)
				}
			}
		}
	}

	return nil
}

func flattenApiManagementHostnameConfiguration(input *[]apimanagementservice.HostnameConfiguration, d *pluginsdk.ResourceData, name, apimHostNameSuffix string) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	managementResults := make([]any, 0)
	portalResults := make([]any, 0)
	developerPortalResults := make([]any, 0)
	gatewayResults := make([]any, 0)
	scmResults := make([]any, 0)

	for _, config := range *input {
		output := make(map[string]any)

		// There'll always be a default custom domain with hostName "apim_name.azure-api.net" and Type "Proxy", which should be ignored
		if config.HostName == strings.ToLower(name)+"."+apimHostNameSuffix && config.Type == apimanagementservice.HostnameTypeProxy {
			continue
		}

		output["host_name"] = config.HostName
		output["negotiate_client_certificate"] = pointer.From(config.NegotiateClientCertificate)
		output["key_vault_certificate_id"] = pointer.From(config.KeyVaultId)
		output["ssl_keyvault_identity_client_id"] = pointer.From(config.IdentityClientId)

		var configType string
		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagementservice.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			output["default_ssl_binding"] = pointer.From(config.DefaultSslBinding)
			gatewayResults = append(gatewayResults, output)
			configType = "gateway"

		case strings.ToLower(string(apimanagementservice.HostnameTypeManagement)):
			managementResults = append(managementResults, output)
			configType = "management"

		case strings.ToLower(string(apimanagementservice.HostnameTypePortal)):
			portalResults = append(portalResults, output)
			configType = "portal"

		case strings.ToLower(string(apimanagementservice.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)
			configType = "developer_portal"

		case strings.ToLower(string(apimanagementservice.HostnameTypeScm)):
			scmResults = append(scmResults, output)
			configType = "scm"
		}

		if configType != "" {
			if valsRaw, ok := d.GetOk(configType); ok {
				vals := valsRaw.([]any)
				schemaz.CopyCertificateAndPassword(vals, config.HostName, output)
			}
		}
	}

	res := map[string]any{
		"management":       managementResults,
		"portal":           portalResults,
		"developer_portal": developerPortalResults,
		"scm":              scmResults,
		"gateway":          gatewayResults,
	}

	return []any{res}
}
