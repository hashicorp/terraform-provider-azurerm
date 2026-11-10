// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			id, err := cognitiveservicesprojects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			accountId := cognitiveservicesprojects.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)
			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			var model CognitiveAccountProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.ProjectsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("description") {
				if payload.Properties == nil {
					payload.Properties = &cognitiveservicesprojects.ProjectProperties{}
				}
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("display_name") {
				if payload.Properties == nil {
					payload.Properties = &cognitiveservicesprojects.ProjectProperties{}
				}
				payload.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if err := client.ProjectsCreateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
