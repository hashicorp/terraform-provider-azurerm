// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type templateDeploymentDebugLevel string

const (
	debugLevelNone                          templateDeploymentDebugLevel = "none"
	debugLevelRequestContent                templateDeploymentDebugLevel = "requestContent"
	debugLevelResponseContent               templateDeploymentDebugLevel = "responseContent"
	debugLevelRequestContentResponseContent templateDeploymentDebugLevel = "requestContent, responseContent"
)

var templateDeploymentDebugLevels = []string{
	string(debugLevelNone),
	string(debugLevelRequestContent),
	string(debugLevelResponseContent),
	string(debugLevelRequestContentResponseContent),
}

func expandTemplateDeploymentDebugSetting(debugLevel string) *resources.DebugSetting {
	if debugLevel == "" {
		return &resources.DebugSetting{
			DetailLevel: utils.String(string(debugLevelNone)),
		}
	}

	return &resources.DebugSetting{
		DetailLevel: utils.String(debugLevel),
	}
}

func flattenTemplateDeploymentDebugSetting(input *resources.DebugSetting) string {
	if input == nil || input.DetailLevel == nil || strings.EqualFold(*input.DetailLevel, string(debugLevelNone)) {
		return ""
	}

	for _, key := range templateDeploymentDebugLevels {
		if strings.EqualFold(key, *input.DetailLevel) {
			return key
		}
	}

	return ""
}

func expandTemplateDeploymentBody(input string) (*map[string]interface{}, error) {
	var output map[string]interface{}

	if err := json.Unmarshal([]byte(input), &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func flattenTemplateDeploymentBody(input interface{}) (*string, error) {
	output := "{}" // since this should be json

	if input == nil {
		return &output, nil
	}

	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshalling json: %+v", err)
	}

	output = string(bytes)
	return &output, nil
}

func filterOutTemplateDeploymentParameters(input interface{}) interface{} {
	if input == nil {
		return nil
	}

	items, ok := input.(map[string]interface{})
	if !ok {
		// this is best-effort
		return input
	}

	/*
			Example:

		{
			"someParam": {
				"type": "String",
				"value": "ohhai"
			},
			"dnsLabelPrefix": {
				"reference": {
					"keyvault": {
						"id": "/some/id/that/doesnt/matter/right/now"
					},
					"secretName": "some-name"
				}
			}
		}
	*/

	output := make(map[string]interface{})
	for topLevelKey, topLevelValue := range items {
		if topLevelValue == nil {
			continue
		}

		// give us the original
		output[topLevelKey] = topLevelValue

		// then filter it if necessary
		if innerVals, ok := topLevelValue.(map[string]interface{}); ok {
			outputVals := make(map[string]interface{})
			for innerKey, innerValue := range innerVals {
				if strings.EqualFold("type", innerKey) {
					continue
				}

				outputVals[innerKey] = innerValue
			}
			output[topLevelKey] = outputVals
		}
	}

	return output
}

func deleteNestedResource(ctx context.Context, resourcesClient *resources.Client, resourceProviderApiVersions *map[string]string, nestedResource resources.Reference) error {
	parsedId, err := azure.ParseAzureResourceID(*nestedResource.ID)
	if err != nil {
		return fmt.Errorf("parsing ID %q from Template Output to delete it: %+v", *nestedResource.ID, err)
	}

	resourceProviderApiVersion, ok := (*resourceProviderApiVersions)[strings.ToLower(parsedId.Provider)]
	if !ok {
		resourceProviderApiVersion, ok = (*resourceProviderApiVersions)[strings.ToLower(parsedId.SecondaryProvider)]
		if !ok {
			return fmt.Errorf("API version information for RP %q (%q) was not found - nestedResource=%q", parsedId.Provider, parsedId.SecondaryProvider, *nestedResource.ID)
		}
	}

	log.Printf("[DEBUG] Deleting Nested Resource %q..", *nestedResource.ID)
	future, err := resourcesClient.DeleteByID(ctx, *nestedResource.ID, resourceProviderApiVersion)

	// NOTE: resourceProviderApiVersion is gotten from one of resource types of the provider.
	// When the provider has multiple resource types, it may cause API version mismatched.
	// For such error, try to get available API version from error code. Ugly but this seems sufficient for now
	if err != nil && strings.Contains(err.Error(), `Code="NoRegisteredProviderFound"`) {
		apiPat := regexp.MustCompile(`\d{4}-\d{2}-\d{2}(-preview)*`)
		matches := apiPat.FindAllStringSubmatch(err.Error(), -1)
		for _, match := range matches {
			if resourceProviderApiVersion != match[0] {
				future, err = resourcesClient.DeleteByID(ctx, *nestedResource.ID, match[0])
				break
			}
		}
	}

	if err != nil {
		if resp := future.Response(); resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Nested Resource %q has been deleted.. continuing..", *nestedResource.ID)
			return nil
		}

		return fmt.Errorf("deleting Nested Resource %q: %+v", *nestedResource.ID, err)
	}

	log.Printf("[DEBUG] Waiting for Deletion of Nested Resource %q..", *nestedResource.ID)
	if err := future.WaitForCompletionRef(ctx, resourcesClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Nested Resource %q: %+v", *nestedResource.ID, err)
	}

	log.Printf("[DEBUG] Deleted Nested Resource %q.", *nestedResource.ID)
	return nil
}

func deleteItemsProvisionedByTemplate(ctx context.Context, client *client.Client, properties resources.DeploymentPropertiesExtended, subscriptionId string) error {
	if properties.Providers == nil {
		return fmt.Errorf("`properties.Providers` was nil - insufficient data to clean up this Template Deployment")
	}
	if properties.OutputResources == nil {
		return fmt.Errorf("`properties.OutputResources` was nil - insufficient data to clean up this Template Deployment")
	}

	providersClient := client.ResourceProvidersClient
	resourcesClient := client.LegacyResourcesClient

	log.Printf("[DEBUG] Determining the API Versions used for Resources provisioned in this Template..")
	resourceProviderApiVersions, err := determineResourceProviderAPIVersionsForResources(ctx, providersClient, *properties.Providers, subscriptionId)
	if err != nil {
		return fmt.Errorf("determining API Versions for Resource Providers: %+v", err)
	}

	log.Printf("[DEBUG] Deleting the resources provisioned in this Template..")
	nestedResources := *properties.OutputResources
	deletedResources := make(map[string]bool)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("could not retrieve context deadline")
	}

	return pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
		deletedTimes := 0
		var errorList []error
		for _, nestedResource := range nestedResources {
			if nestedResource.ID == nil {
				continue
			}

			if _, exists := deletedResources[*nestedResource.ID]; exists {
				continue
			}

			err = deleteNestedResource(ctx, resourcesClient, resourceProviderApiVersions, nestedResource)
			if err != nil {
				errorList = append(errorList, err)
			} else {
				deletedResources[*nestedResource.ID] = true
				deletedTimes++
			}
		}

		if deletedTimes > 0 {
			return pluginsdk.RetryableError(fmt.Errorf("may exist nested resources to delete, retrying"))
		}

		// If `deletedTimes` is 0, it means all resources have been successfully deleted if the `errorList` is empty, or the remaining resources cannot be deleted
		if len(errorList) > 0 {
			return pluginsdk.NonRetryableError(fmt.Errorf("%+v", errorList[0]))
		}

		return nil
	})
}

func determineResourceProviderAPIVersionsForResources(ctx context.Context, client *providers.ProvidersClient, resourceProviders []resources.Provider, subscriptionId string) (*map[string]string, error) {
	resourceProviderApiVersions := make(map[string]string)

	for _, provider := range resourceProviders {
		if provider.Namespace == nil {
			continue
		}

		providerId := providers.NewSubscriptionProviderID(subscriptionId, *provider.Namespace)
		providerResp, err := client.Get(ctx, providerId, providers.DefaultGetOperationOptions())
		if err != nil {
			return nil, fmt.Errorf("retrieving MetaData for %s: %+v", providerId, err)
		}
		resourceTypes := make([]providers.ProviderResourceType, 0)
		if model := providerResp.Model; model != nil && model.ResourceTypes != nil {
			resourceTypes = *model.ResourceTypes
		}
		if len(resourceTypes) == 0 {
			return nil, fmt.Errorf("`resourceTypes` was nil/empty for %s", providerId)
		}

		for _, resourceType := range *provider.ResourceTypes {
			resourceTypeName := *resourceType.ResourceType
			availableResourceTypes := resourceTypes
			apiVersion := findApiVersionForResourceType(resourceTypeName, availableResourceTypes)
			if apiVersion == nil {
				return nil, fmt.Errorf("unable to determine API version for Resource Type %q (%s)", resourceTypeName, providerId)
			}

			// NOTE: there's an enhancement in that not all RP's necessarily offer everything in every version
			// but the majority do, so this is likely sufficient for now
			resourceProviderApiVersions[strings.ToLower(providerId.ProviderName)] = *apiVersion
			break
		}
	}

	return &resourceProviderApiVersions, nil
}

func findApiVersionForResourceType(resourceType string, availableResourceTypes []providers.ProviderResourceType) *string {
	for _, item := range availableResourceTypes {
		if item.ResourceType == nil || item.ApiVersions == nil {
			continue
		}

		isExactMatch := strings.EqualFold(resourceType, *item.ResourceType)
		isPrefixMatch := strings.HasPrefix(strings.ToLower(resourceType), strings.ToLower(*item.ResourceType))
		if isExactMatch || isPrefixMatch {
			apiVersions := *item.ApiVersions
			apiVersion := apiVersions[0]
			return &apiVersion
		}
	}

	return nil
}
