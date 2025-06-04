// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppResource struct{}

type ContainerAppModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroup        string `tfschema:"resource_group_name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`
	Location             string `tfschema:"location"`

	RevisionMode string                      `tfschema:"revision_mode"`
	Ingress      []helpers.Ingress           `tfschema:"ingress"`
	Registries   []helpers.Registry          `tfschema:"registry"`
	Secrets      []helpers.Secret            `tfschema:"secret"`
	Dapr         []helpers.Dapr              `tfschema:"dapr"`
	Template     []helpers.ContainerTemplate `tfschema:"template"`

	Identity             []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	WorkloadProfileName  string                                     `tfschema:"workload_profile_name"`
	MaxInactiveRevisions int64                                      `tfschema:"max_inactive_revisions"`
	Tags                 map[string]interface{}                     `tfschema:"tags"`

	OutboundIpAddresses        []string `tfschema:"outbound_ip_addresses"`
	LatestRevisionName         string   `tfschema:"latest_revision_name"`
	LatestRevisionFqdn         string   `tfschema:"latest_revision_fqdn"`
	CustomDomainVerificationId string   `tfschema:"custom_domain_verification_id"`
}

var _ sdk.ResourceWithUpdate = ContainerAppResource{}

var _ sdk.ResourceWithCustomizeDiff = ContainerAppResource{}

func (r ContainerAppResource) ModelObject() interface{} {
	return &ContainerAppModel{}
}

func (r ContainerAppResource) ResourceType() string {
	return "azurerm_container_app"
}

func (r ContainerAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return containerapps.ValidateContainerAppID
}

func (r ContainerAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerAppName,
			Description:  "The name for this Container App.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to host this Container App.",
		},

		"template": helpers.ContainerTemplateSchema(),

		"revision_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(containerapps.ActiveRevisionsModeSingle),
				string(containerapps.ActiveRevisionsModeMultiple),
			}, false),
		},

		"ingress": helpers.ContainerAppIngressSchema(),

		"registry": helpers.ContainerAppRegistrySchema(),

		"secret": helpers.SecretsSchema(),

		"dapr": helpers.ContainerDaprSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"workload_profile_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"max_inactive_revisions": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ContainerAppResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"latest_revision_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the latest Container Revision.",
		},

		"latest_revision_fqdn": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The FQDN of the Latest Revision of the Container App.",
		},

		"custom_domain_verification_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The ID of the Custom Domain Verification for this Container App.",
		},
	}
}

func (r ContainerAppResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
						MaxInactiveRevisions: pointer.FromInt64(app.MaxInactiveRevisions),
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

			containerApp.Properties.Configuration.ActiveRevisionsMode = pointer.To(containerapps.ActiveRevisionsMode(app.RevisionMode))

			if err := client.CreateOrUpdateThenPoll(ctx, id, containerApp); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppResource) Read() sdk.ResourceFunc {
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
							state.RevisionMode = string(pointer.From(config.ActiveRevisionsMode))
						}
						state.Ingress = helpers.FlattenContainerAppIngress(config.Ingress, id.ContainerAppName)
						state.Registries = helpers.FlattenContainerAppRegistries(config.Registries)
						state.Dapr = helpers.FlattenContainerAppDapr(config.Dapr)
						state.MaxInactiveRevisions = pointer.ToInt64(config.MaxInactiveRevisions)
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

func (r ContainerAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			id, err := containerapps.ParseContainerAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppResource) Update() sdk.ResourceFunc {
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
				return fmt.Errorf("retreiving properties for %s for update: %+v", *id, err)
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
				model.Properties.Configuration.ActiveRevisionsMode = pointer.To(containerapps.ActiveRevisionsMode(state.RevisionMode))
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
				model.Properties.Configuration.MaxInactiveRevisions = pointer.FromInt64(state.MaxInactiveRevisions)
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

func (r ContainerAppResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}
			var app ContainerAppModel
			if err := metadata.DecodeDiff(&app); err != nil {
				return err
			}
			// Ingress traffic weight validations
			if len(app.Ingress) != 0 {
				ingress := app.Ingress[0]

				for i, tw := range ingress.TrafficWeights {
					if !tw.LatestRevision && tw.RevisionSuffix == "" {
						return fmt.Errorf("`either ingress.0.traffic_weight.%[1]d.revision_suffix` or `ingress.0.traffic_weight.%[1]d.latest_revision` should be specified", i)
					}
				}
			}

			for _, s := range app.Secrets {
				if s.KeyVaultSecretId != "" && s.Identity == "" {
					return fmt.Errorf("secret %s must supply identity for key vault secret id", s.Name)
				}
				if s.KeyVaultSecretId == "" && s.Identity != "" {
					return fmt.Errorf("secret %s must supply key vault secret id when specifying identity", s.Name)
				}
			}
			return nil
		},
	}
}
