// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_ai_foundry"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.MachineLearning.Workspaces
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Kind == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if !strings.EqualFold(*resp.Model.Kind, "Hub") {
			return fmt.Errorf("importing %s: specified workspace is not of kind `Hub`, got `%s`", id, *resp.Model.Kind)
		}

		return nil
	}
}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomImporter = Resource{}
