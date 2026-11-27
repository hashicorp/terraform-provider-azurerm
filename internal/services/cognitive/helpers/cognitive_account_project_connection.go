// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func CognitiveAccountProjectConnectionImporter(expectedAuthType projectconnectionresource.ConnectionAuthType) sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Cognitive.ProjectConnectionResourceClient

		id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.ProjectConnectionsGet(ctx, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return metadata.MarkAsGone(id)
			}
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if err := ValidateCognitiveAccountProjectConnectionAuthType(resp.Model, expectedAuthType); err != nil {
			return fmt.Errorf("importing %s: %+v", *id, err)
		}

		return nil
	}
}

func ValidateCognitiveAccountProjectConnectionAuthType(model *projectconnectionresource.ConnectionPropertiesV2BasicResource, expectedAuthType projectconnectionresource.ConnectionAuthType) error {
	if model == nil || model.Properties == nil {
		return errors.New("connection `properties` were nil")
	}

	if !CognitiveAccountProjectConnectionHasExpectedAuthType(model, expectedAuthType) {
		authType := model.Properties.ConnectionPropertiesV2().AuthType
		return fmt.Errorf("connection has auth type `%s`, expected `%s`", authType, expectedAuthType)
	}
	return nil
}

func CognitiveAccountProjectConnectionHasExpectedAuthType(model *projectconnectionresource.ConnectionPropertiesV2BasicResource, expectedAuthType projectconnectionresource.ConnectionAuthType) bool {
	if model == nil || model.Properties == nil {
		return false
	}
	return model.Properties.ConnectionPropertiesV2().AuthType == expectedAuthType
}

func FlattenAccountProjectConnectionMetadata(priorMetadata map[string]string, apiMetadata *map[string]string) map[string]string {
	// The API returns additional metadata fields beyond those configured (e.g. `ApiVersion`,
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
