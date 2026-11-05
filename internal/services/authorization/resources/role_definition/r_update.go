// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_definition

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			stateId, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			id := roledefinitions.NewScopedRoleDefinitionID(stateId.Scope, stateId.RoleID)

			var config RoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", stateId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", stateId)
			}

			model := *existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", stateId)
			}

			props := *model.Properties

			if metadata.ResourceData.HasChange("name") {
				props.RoleName = &config.Name
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = &config.Description
			}

			if metadata.ResourceData.HasChange("permissions") {
				props.Permissions = pointer.To(expandRoleDefinitionPermissions(config.Permissions))
			}

			if metadata.ResourceData.HasChange("assignable_scopes") {
				props.AssignableScopes = pointer.To(expandRoleDefinitionAssignableScopes(config))
			}

			model.Properties = &props

			resp, err := client.CreateOrUpdate(ctx, id, model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", stateId, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("updating %s: model was nil", stateId)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("updating %s: properties was nil", stateId)
			}

			updatedOn := resp.Model.Properties.UpdatedOn
			if updatedOn == nil {
				return fmt.Errorf("updating Role Definition %q (Scope %q): `properties.UpdatedOn` was nil", stateId.RoleID, stateId.Scope)
			}

			// "Updating" a role definition actually creates a new one and these get consolidated a few seconds later
			// where the "create date" and "update date" match for the newly created record
			// but eventually switch to being the old create date and the new update date
			// ergo we can for the old create date and the new updated date
			log.Printf("[DEBUG] Waiting for %s to settle down..", stateId)
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				ContinuousTargetOccurence: 12,
				Delay:                     60 * time.Second,
				MinTimeout:                10 * time.Second,
				Pending:                   []string{"Pending"},
				Target:                    []string{"Updated"},
				Refresh:                   roleDefinitionEventualConsistencyUpdate(ctx, client, id, *updatedOn),
				Timeout:                   time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to settle down: %+v", stateId, err)
			}

			return nil
		},
	}
}
