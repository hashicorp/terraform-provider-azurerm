// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			id, err := containerapps.ParseContainerAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppModel

			state.Name = id.ContainerAppName
			state.ResourceGroup = id.ResourceGroupName

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)
				if model.Identity != nil {
					ident, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To(identity.SystemAndUserAssignedMap(*model.Identity)))
					if err != nil {
						return err
					}
					state.Identity = pointer.From(ident)
				}

				if props := model.Properties; props != nil {
					envId, err := managedenvironments.ParseManagedEnvironmentIDInsensitively(pointer.From(props.ManagedEnvironmentId))
					if err != nil {
						return err
					}
					state.ManagedEnvironmentId = envId.ID()
					state.Template = helpers.FlattenContainerAppTemplate(props.Template)
					if config := props.Configuration; config != nil {
						if config.ActiveRevisionsMode != nil {
							state.RevisionMode = pointer.FromEnum(config.ActiveRevisionsMode)
						}
						state.Ingress = helpers.FlattenContainerAppIngress(config.Ingress, id.ContainerAppName)
						state.Registries = helpers.FlattenContainerAppRegistries(config.Registries)
						state.Dapr = helpers.FlattenContainerAppDapr(config.Dapr)
						state.MaxInactiveRevisions = pointer.From(config.MaxInactiveRevisions)
					}
					state.LatestRevisionName = pointer.From(props.LatestRevisionName)
					state.LatestRevisionFqdn = pointer.From(props.LatestRevisionFqdn)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.OutboundIpAddresses = pointer.From(props.OutboundIPAddresses)
					state.WorkloadProfileName = pointer.From(props.WorkloadProfileName)
				}
			}

			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
			}

			state.Secrets = helpers.FlattenContainerAppSecrets(secretsResp.Model)

			return metadata.Encode(&state)
		},
	}
}
