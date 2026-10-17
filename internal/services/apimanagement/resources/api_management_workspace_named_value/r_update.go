// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_named_value

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namedvalue.ParseWorkspaceNamedValueID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.WorkspaceNamedValueGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := namedvalue.NamedValueCreateContract{
				Properties: &namedvalue.NamedValueCreateContractProperties{
					DisplayName: existing.Model.Properties.DisplayName,
					Secret:      existing.Model.Properties.Secret,
					Tags:        existing.Model.Properties.Tags,
					Value:       existing.Model.Properties.Value,
				},
			}

			if existing.Model.Properties.KeyVault != nil {
				parameters.Properties.KeyVault = &namedvalue.KeyVaultContractCreateProperties{
					SecretIdentifier: existing.Model.Properties.KeyVault.SecretIdentifier,
					IdentityClientId: existing.Model.Properties.KeyVault.IdentityClientId,
				}
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("secret") {
				parameters.Properties.Secret = pointer.To(model.Secret)
			}

			if metadata.ResourceData.HasChange("value") {
				if model.Value != "" {
					parameters.Properties.Value = pointer.To(model.Value)
				} else {
					parameters.Properties.Value = nil
				}
			}

			if metadata.ResourceData.HasChange("value_from_key_vault") {
				parameters.Properties.KeyVault = expandApiManagementWorkspaceNamedValueKeyVault(model.ValueFromKeyVault)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Properties.Tags = pointer.To(model.Tags)
			}

			if err := client.WorkspaceNamedValueCreateOrUpdateThenPoll(ctx, *id, parameters, namedvalue.DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
