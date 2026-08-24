// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func cognitiveAccountProjectConnectionImporter(expectedAuthType projectconnectionresource.ConnectionAuthType, resourceType string) sdk.ResourceRunFunc {
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

		if model := resp.Model; model != nil && model.Properties != nil {
			if authType := model.Properties.ConnectionPropertiesV2().AuthType; authType != expectedAuthType {
				return fmt.Errorf("connection %s has auth type `%s` and cannot be managed by `%s`", *id, authType, resourceType)
			}
		}

		return nil
	}
}
