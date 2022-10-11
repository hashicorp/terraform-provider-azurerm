package containerapps

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerAppResource struct{}

type ContainerAppModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroup        string `tfschema:"resource_group_name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`
	Location             string `tfschema:"location"`

	RevisionMode string                      `tfschema:"revision_mode"`
	Ingress      []helpers.Ingress           `tfschema:"ingress"`
	Registries   []helpers.Registry          `tfschema:"registry"` // TODO - Need to do some ACR exploration
	Secrets      []helpers.Secret            `tfschema:"secret"`
	Dapr         []helpers.Dapr              `tfschema:"dapr"`
	Template     []helpers.ContainerTemplate `tfschema:"template"`

	// Identity identity.LegacySystemAndUserAssignedMap `tfschema:"identity"` // TODO - when the basics are working...

	Tags map[string]interface{} `tfschema:"tags"`

	OutboundIpAddresses        []string `tfschema:"outbound_ip_addresses"`
	LatestRevisionName         string   `tfschema:"latest_revision_name"`
	LatestRevisionFqdn         string   `tfschema:"latest_revision_fqdn"`
	CustomDomainVerificationId string   `tfschema:"custom_domain_verification_id"`

	// TODO - Expose SystemData ?
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
			ValidateFunc: helpers.ValidateContainerAppName,
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
			}, true),
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"ingress": helpers.ContainerAppIngressSchema(),

		"registry": helpers.ContainerAppRegistrySchema(),

		"secret": helpers.SecretsSchema(),

		"dapr": helpers.ContainerDaprSchema(),

		// "identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

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
			Description: "The Custom Domain Verification ID for the Container App.",
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
			if err != nil || env.Model == nil {
				return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
			}

			containerApp := containerapps.ContainerApp{
				Name:     utils.String(app.Name),
				Location: env.Model.Location,
				Properties: &containerapps.ContainerAppProperties{
					Configuration: &containerapps.Configuration{
						Ingress:    helpers.ExpandContainerAppIngress(app.Ingress, id.ContainerAppName),
						Dapr:       helpers.ExpandContainerAppDapr(app.Dapr),
						Secrets:    helpers.ExpandContainerSecrets(app.Secrets),
						Registries: helpers.ExpandContainerAppRegistries(app.Registries),
					},
					ManagedEnvironmentId: utils.String(app.ManagedEnvironmentId),
					Template:             helpers.ExpandContainerAppTemplate(app.Template, metadata),
				},
				// Identity: &app.Identity,
				Tags: tags.Expand(app.Tags),
			}

			revisionMode := containerapps.ActiveRevisionsMode(app.RevisionMode)
			containerApp.Properties.Configuration.ActiveRevisionsMode = &revisionMode

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
				state.Location = model.Location
				state.Tags = tags.Flatten(model.Tags)
				// state.Identity = []identity.LegacySystemAndUserAssignedMap{*model.Identity}

				if props := model.Properties; props != nil {
					envId, err := managedenvironments.ParseManagedEnvironmentIDInsensitively(utils.NormalizeNilableString(props.ManagedEnvironmentId))
					if err != nil {
						return err
					}
					state.ManagedEnvironmentId = envId.ID()
					state.Template = helpers.FlattenContainerAppTemplate(props.Template)
					if config := props.Configuration; config != nil {
						if config.ActiveRevisionsMode != nil {
							state.RevisionMode = strings.ToLower(string(*config.ActiveRevisionsMode))
						}
						state.Ingress = helpers.FlattenContainerAppIngress(config.Ingress, id.ContainerAppName)
						state.Registries = helpers.FlattenContainerAppRegistries(config.Registries)
						state.Dapr = helpers.FlattenContainerAppDapr(config.Dapr)
					}
					state.LatestRevisionName = utils.NormalizeNilableString(props.LatestRevisionName)
					state.LatestRevisionFqdn = utils.NormalizeNilableString(props.LatestRevisionFqdn)
					state.CustomDomainVerificationId = utils.NormalizeNilableString(props.CustomDomainVerificationId)
					state.OutboundIpAddresses = *props.OutboundIPAddresses
				}
				// state.Identity = *model.Identity
			}

			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil || secretsResp.Model == nil {
				if secretsResp.HttpResponse == nil || secretsResp.HttpResponse.StatusCode != http.StatusNoContent {
					return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
				}
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

			resp, err := client.Delete(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			if err := resp.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting for deletion of %s", *id)
			}

			return nil
		},
	}
}

func (r ContainerAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
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
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := existing.Model

			// Delta-updates need the secrets back from the list API, or we'll end up removing them or erroring out.
			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil || secretsResp.Model == nil {
				if secretsResp.HttpResponse == nil || secretsResp.HttpResponse.StatusCode != http.StatusNoContent {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", *id, err)
				}
			}
			model.Properties.Configuration.Secrets = helpers.UnpackContainerSecretsCollection(secretsResp.Model)

			if metadata.ResourceData.HasChange("revision_mode") {
				revisionMode := containerapps.ActiveRevisionsMode(state.RevisionMode)
				model.Properties.Configuration.ActiveRevisionsMode = &revisionMode
			}

			if metadata.ResourceData.HasChange("ingress") {
				model.Properties.Configuration.Ingress = helpers.ExpandContainerAppIngress(state.Ingress, id.ContainerAppName)
			}

			if metadata.ResourceData.HasChange("registry") {
				model.Properties.Configuration.Registries = helpers.ExpandContainerAppRegistries(state.Registries)

			}

			if metadata.ResourceData.HasChange("dapr") {
				model.Properties.Configuration.Dapr = helpers.ExpandContainerAppDapr(state.Dapr)

			}

			if metadata.ResourceData.HasChange("template") {
				allProbesRemoved := helpers.ContainerAppProbesRemoved(metadata)
				if allProbesRemoved {
					nilProbes := make([]containerapps.ContainerAppProbe, 0)
					containers := *model.Properties.Template.Containers
					containers[0].Probes = &nilProbes
					model.Properties.Template.Containers = &containers
				}
			}

			if metadata.ResourceData.HasChange("secret") {
				model.Properties.Configuration.Secrets = helpers.ExpandContainerSecrets(state.Secrets)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = tags.Expand(state.Tags)
			}

			model.Properties.Template = helpers.ExpandContainerAppTemplate(state.Template, metadata)

			// Zero R/O - API rejects the request if eny of these are set
			model.SystemData = nil
			model.Properties.OutboundIPAddresses = nil

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
			if metadata.ResourceDiff != nil && metadata.ResourceDiff.HasChange("secret") {
				stateSecretsRaw, configSecretsRaw := metadata.ResourceDiff.GetChange("secret")
				stateSecrets := stateSecretsRaw.([]interface{})
				configSecrets := configSecretsRaw.([]interface{})
				// Check there's not less
				if len(configSecrets) < len(stateSecrets) {
					return fmt.Errorf("cannot remove secrets from Container Apps at this time. Please see `https://github.com/microsoft/azure-container-apps/issues/395` for more details")
				}
				// Check secrets names in state are all present in config, the values don't matter
				if len(stateSecrets) > 0 {
					for _, s := range stateSecrets {
						found := false
						for _, c := range configSecrets {
							if s.(map[string]interface{})["name"] == c.(map[string]interface{})["name"] {
								found = true
								break
							}
							if !found {
								return fmt.Errorf("previously configured secret %q was removed. Removing secrets is not supported at this time, see `https://github.com/microsoft/azure-container-apps/issues/395` for more details", s.(map[string]interface{})["name"])
							}
						}
					}
				}
			}
			return nil
		},
	}
}
