package cognitive

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func cognitiveAccountConnectionImporter(expectedAuthType accountconnectionresource.ConnectionAuthType, resourceType string) sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Cognitive.AccountConnectionResourceClient

		id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.AccountConnectionsGet(ctx, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return metadata.MarkAsGone(id)
			}
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if model := resp.Model; model != nil && model.Properties != nil {
			if authType := model.Properties.ConnectionPropertiesV2().AuthType; authType != expectedAuthType {
				return fmt.Errorf("connection %s has auth type `%s` and cannot be managed by `%s`", *id, authType, resourceType)
			}
		}

		return nil
	}
}

func flattenAccountConnectionMetadata(priorMetadata map[string]string, apiMetadata *map[string]string) map[string]string {
	// Some Connection APIs return additional metadata fields beyond those configured (e.g. `ApiVersion`,
	// `DeploymentApiVersion`). When prior configuration is known (Read/Update) only the configured
	// keys are surfaced to avoid diffs; otherwise (import or list) all API metadata fields are returned.
	apiMetadataValues := pointer.From(apiMetadata)
	if len(priorMetadata) == 0 {
		return apiMetadataValues
	}

	filtered := make(map[string]string)
	for configKey := range priorMetadata {
		for apiKey, apiValue := range apiMetadataValues {
			if strings.EqualFold(configKey, apiKey) {
				filtered[configKey] = apiValue
				break
			}
		}
	}
	return filtered
}
