// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_custom_keys

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			var model CognitiveAccountProjectConnectionCustomKeysModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			projectId, err := projectconnectionresource.ParseProjectID(model.CognitiveAccountProjectId)
			if err != nil {
				return err
			}

			locks.ByID(projectId.ID())
			defer locks.UnlockByID(projectId.ID())

			id := projectconnectionresource.NewProjectConnectionID(projectId.SubscriptionId, projectId.ResourceGroupName, projectId.AccountName, projectId.ProjectName, model.Name)
			existing, err := client.ProjectConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
				if err := helpers.ValidateCognitiveAccountProjectConnectionAuthType(existing.Model, projectconnectionresource.ConnectionAuthTypeCustomKeys, r.ResourceType()); err != nil {
					return fmt.Errorf("checking existing %s: %+v", id, err)
				}
			}

			properties := projectconnectionresource.CustomKeysConnectionProperties{
				AuthType: projectconnectionresource.ConnectionAuthTypeCustomKeys,
				Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
				Target:   pointer.To(model.Target),
				Credentials: &projectconnectionresource.CustomKeys{
					Keys: pointer.To(model.CustomKeys),
				},
			}
			if len(model.Metadata) > 0 {
				properties.Metadata = pointer.To(model.Metadata)
			}

			connection := projectconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}
			if _, err := client.ProjectConnectionsCreate(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
