package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	providers "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

func deleteItemsProvisionedByTemplate(ctx context.Context, client *client.Client, properties resources.DeploymentPropertiesExtended) error {
	if properties.Providers == nil {
		return fmt.Errorf("`properties.Providers` was nil - insufficient data to clean up this Template Deployment")
	}
	if properties.OutputResources == nil {
		return fmt.Errorf("`properties.OutputResources` was nil - insufficient data to clean up this Template Deployment")
	}

	providersClient := client.ProvidersClient
	resourcesClient := client.ResourcesClient

	log.Printf("[DEBUG] Determining the API Versions used for Resources provisioned in this Template..")
	resourceProviderApiVersions, err := determineResourceProviderAPIVersionsForResources(ctx, providersClient, *properties.Providers)
	if err != nil {
		return fmt.Errorf("determining API Versions for Resource Providers: %+v", err)
	}

	log.Printf("[DEBUG] Deleting the resources provisioned in this Template..")
	nestedResources := *properties.OutputResources

	// NOTE: this likely wants splitting out into a parallel loop which retries resources which fail deletion
	// for example when A depends on B - try deleting them both then loop, and if both subsequentally fail then error
	// but this seems sufficient for now
	for _, nestedResource := range nestedResources {
		if nestedResource.ID == nil {
			continue
		}

		parsedId, err := azure.ParseAzureResourceID(*nestedResource.ID)
		if err != nil {
			return fmt.Errorf("parsing ID %q from Template Output to delete it: %+v", *nestedResource.ID, err)
		}

		resourceProviderApiVersion, ok := (*resourceProviderApiVersions)[parsedId.Provider]
		if !ok {
			return fmt.Errorf("API version information for RP %q was not found", parsedId.Provider)
		}

		log.Printf("[DEBUG] Deleting Nested Resource %q..", *nestedResource.ID)
		future, err := resourcesClient.DeleteByID(ctx, *nestedResource.ID, resourceProviderApiVersion)
		if err != nil {
			if resp := future.Response(); resp != nil && resp.StatusCode == http.StatusNotFound {
				log.Printf("[DEBUG] Nested Resource %q has been deleted.. continuing..", *nestedResource.ID)
				continue
			}
			return fmt.Errorf("deleting Nested Resource %q: %+v", *nestedResource.ID, err)
		}

		log.Printf("[DEBUG] Waiting for Deletion of Nested Resource %q..", *nestedResource.ID)
		if err := future.WaitForCompletionRef(ctx, resourcesClient.Client); err != nil {
			return fmt.Errorf("waiting for deletion of Nested Resource %q: %+v", *nestedResource.ID, err)
		}

		log.Printf("[DEBUG] Deleted Nested Resource %q.", *nestedResource.ID)
	}

	return nil
}

func determineResourceProviderAPIVersionsForResources(ctx context.Context, client *providers.ProvidersClient, providers []resources.Provider) (*map[string]string, error) {
	resourceProviderApiVersions := make(map[string]string)

	for _, provider := range providers {
		if provider.Namespace == nil {
			continue
		}

		resourceProviderName := *provider.Namespace
		providerResp, err := client.Get(ctx, resourceProviderName, "")
		if err != nil {
			return nil, fmt.Errorf("retrieving Resource Provider MetaData for %q: %+v", resourceProviderName, err)
		}
		if providerResp.ResourceTypes == nil {
			return nil, fmt.Errorf("`resourceTypes` was nil for Resource Provider %q", resourceProviderName)
		}

		for _, resourceType := range *provider.ResourceTypes {
			resourceTypeName := *resourceType.ResourceType
			availableResourceTypes := *providerResp.ResourceTypes
			apiVersion := findApiVersionForResourceType(resourceTypeName, availableResourceTypes)
			if apiVersion == nil {
				return nil, fmt.Errorf("unable to determine API version for Resource Type %q (Resource Provider %q)", resourceTypeName, resourceProviderName)
			}

			// NOTE: there's an enhancement in that not all RP's necessarily offer everything in every version
			// but the majority do, so this is likely sufficient for now
			resourceProviderApiVersions[resourceProviderName] = *apiVersion
			break
		}
	}

	return &resourceProviderApiVersions, nil
}

func findApiVersionForResourceType(resourceType string, availableResourceTypes []providers.ProviderResourceType) *string {
	for _, item := range availableResourceTypes {
		if item.ResourceType == nil || item.APIVersions == nil {
			continue
		}

		if strings.HasPrefix(resourceType, *item.ResourceType) {
			apiVersions := *item.APIVersions
			apiVersion := apiVersions[0]
			return &apiVersion
		}
	}

	return nil
}
