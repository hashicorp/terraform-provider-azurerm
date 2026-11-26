// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			id, err := containerapps.ParseContainerAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retrieving properties for %s for update: %+v", *id, err)
			}

			if model.Properties.Configuration == nil {
				model.Properties.Configuration = &containerapps.Configuration{}
			}

			// Delta-updates need the secrets back from the list API, or we'll end up removing them or erroring out.
			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil || secretsResp.Model == nil {
				if !response.WasStatusCode(secretsResp.HttpResponse, http.StatusNoContent) {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", *id, err)
				}
			}
			model.Properties.Configuration.Secrets = helpers.UnpackContainerSecretsCollection(secretsResp.Model)

			if metadata.ResourceData.HasChange("revision_mode") {
				model.Properties.Configuration.ActiveRevisionsMode = pointer.ToEnum[containerapps.ActiveRevisionsMode](state.RevisionMode)
			}

			if metadata.ResourceData.HasChange("ingress") {
				model.Properties.Configuration.Ingress = helpers.ExpandContainerAppIngress(state.Ingress, id.ContainerAppName)
			}

			if metadata.ResourceData.HasChange("registry") {
				model.Properties.Configuration.Registries, err = helpers.ExpandContainerAppRegistries(state.Registries)
				if err != nil {
					return fmt.Errorf("invalid registry config for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("max_inactive_revisions") {
				model.Properties.Configuration.MaxInactiveRevisions = pointer.To(state.MaxInactiveRevisions)
			}

			if metadata.ResourceData.HasChange("dapr") {
				model.Properties.Configuration.Dapr = helpers.ExpandContainerAppDapr(state.Dapr)
			}

			if metadata.ResourceData.HasChange("template") {
				if model.Properties.Template == nil {
					model.Properties.Template = &containerapps.Template{}
				}
				allProbesRemoved := helpers.ContainerAppProbesRemoved(metadata)
				if allProbesRemoved {
					containers := *model.Properties.Template.Containers
					containers[0].Probes = pointer.To(make([]containerapps.ContainerAppProbe, 0))
					model.Properties.Template.Containers = &containers
				}
			}

			if metadata.ResourceData.HasChange("secret") {
				model.Properties.Configuration.Secrets, err = helpers.ExpandContainerSecrets(state.Secrets)
				if err != nil {
					return fmt.Errorf("invalid secrets config for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return err
				}
				model.Identity = pointer.To(identity.LegacySystemAndUserAssignedMap(*ident))
			}

			if metadata.ResourceData.HasChange("workload_profile_name") {
				model.Properties.WorkloadProfileName = pointer.To(state.WorkloadProfileName)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = tags.Expand(state.Tags)
			}

			model.Properties.Template = helpers.ExpandContainerAppTemplate(state.Template, metadata)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
