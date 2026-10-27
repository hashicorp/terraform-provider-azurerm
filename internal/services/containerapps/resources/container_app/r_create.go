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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient
			environmentClient := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var app ContainerAppModel

			if err := metadata.Decode(&app); err != nil {
				return err
			}

			id := containerapps.NewContainerAppID(subscriptionId, app.ResourceGroup, app.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			envId, err := managedenvironments.ParseManagedEnvironmentID(app.ManagedEnvironmentId)
			if err != nil {
				return fmt.Errorf("parsing Container App Environment ID for %s: %+v", id, err)
			}

			env, err := environmentClient.Get(ctx, *envId)
			if err != nil {
				return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
			}

			registries, err := helpers.ExpandContainerAppRegistries(app.Registries)
			if err != nil {
				return fmt.Errorf("invalid registry config for %s: %+v", id, err)
			}

			secrets, err := helpers.ExpandContainerSecrets(app.Secrets)
			if err != nil {
				return fmt.Errorf("invalid secrets config for %s: %+v", id, err)
			}

			containerApp := containerapps.ContainerApp{
				Location: location.Normalize(env.Model.Location),
				Properties: &containerapps.ContainerAppProperties{
					Configuration: &containerapps.Configuration{
						Ingress:              helpers.ExpandContainerAppIngress(app.Ingress, id.ContainerAppName),
						Dapr:                 helpers.ExpandContainerAppDapr(app.Dapr),
						Secrets:              secrets,
						Registries:           registries,
						MaxInactiveRevisions: pointer.To(app.MaxInactiveRevisions),
					},
					ManagedEnvironmentId: pointer.To(app.ManagedEnvironmentId),
					Template:             helpers.ExpandContainerAppTemplate(app.Template, metadata),
					WorkloadProfileName:  pointer.To(app.WorkloadProfileName),
				},
				Tags: tags.Expand(app.Tags),
			}

			ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(app.Identity)
			if err != nil {
				return err
			}
			containerApp.Identity = pointer.To(identity.LegacySystemAndUserAssignedMap(*ident))

			containerApp.Properties.Configuration.ActiveRevisionsMode = pointer.ToEnum[containerapps.ActiveRevisionsMode](app.RevisionMode)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, containerApp, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
