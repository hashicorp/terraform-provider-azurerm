package appservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2024-03-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	kvValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LinuxFunctionAppOnContainerResource struct{}

type LinuxFunctionAppOnContainerModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroup        string `tfschema:"resource_group_name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	StorageAccountName          string `tfschema:"storage_account_name"`
	StorageAccountKey           string `tfschema:"storage_account_access_key"`
	StorageKeyVaultSecretID     string `tfschema:"storage_key_vault_secret_id"`
	StorageUsesMSI              bool   `tfschema:"storage_uses_managed_identity"`
	KeyVaultReferenceIdentityID string `tfschema:"key_vault_reference_identity_id"`

	AppSettings               map[string]string `tfschema:"app_settings"`
	FunctionExtensionsVersion string            `tfschema:"functions_extension_version"`

	Registries     []helpers.Registry                              `tfschema:"registry"`
	ContainerImage string                                          `tfschema:"container_image"`
	SiteConfig     []helpers.SiteConfigLinuxFunctionAppOnContainer `tfschema:"site_config"`
	Identity       []identity.ModelSystemAssignedUserAssigned      `tfschema:"identity"`

	Tags map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = LinuxFunctionAppOnContainerResource{}

func (r LinuxFunctionAppOnContainerResource) ModelObject() interface{} {
	return &LinuxFunctionAppOnContainerModel{}
}

func (r LinuxFunctionAppOnContainerResource) ResourceType() string {
	return "azurerm_linux_function_app_on_container"
}

func (r LinuxFunctionAppOnContainerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateFunctionAppID
}

func (r LinuxFunctionAppOnContainerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerizedFunctionAppName,
			Description:  "Specifies the name of the Function App.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to host this Container App.",
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			Description:  "The backend storage account name which will be used by this Function App.",
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

		"storage_account_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			ConflictsWith: []string{
				"storage_key_vault_secret_id",
			},
			Description: "The access key which will be used to access the storage account for the Function App.",
		},

		"storage_key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: kvValidate.NestedItemIdWithOptionalVersion,
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
			Description: "The Key Vault Secret ID, including version, that contains the Connection String to connect to the storage account for this Function App.",
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			},
			Description: "Should the Function App use its Managed Identity to access storage?",
		},

		"functions_extension_version": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
			Description: "The runtime version associated with the Function App.",
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A map of key-value pairs for [App Settings](https://docs.microsoft.com/en-us/azure/azure-functions/functions-app-settings) and custom values.",
		},

		"registry": helpers.RegistrySchemaLinuxFunctionAppOnContainer(),

		"container_image": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The name of the Container Image that includes image tag",
		},

		"site_config": helpers.SiteConfigSchemaLinuxFunctionAppOnContainer(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			Description:  "The User Assigned Identity to use for Key Vault access.",
		},

		"tags": tags.Schema(),
	}
}

func (r LinuxFunctionAppOnContainerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r LinuxFunctionAppOnContainerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}
			var linuxFunctionAppOnContainer LinuxFunctionAppOnContainerModel

			if err := metadata.Decode(&linuxFunctionAppOnContainer); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourcesClient := metadata.Client.AppService.ResourceProvidersClient
			containerEnvClient := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := commonids.NewAppServiceID(subscriptionId, linuxFunctionAppOnContainer.ResourceGroup, linuxFunctionAppOnContainer.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: linuxFunctionAppOnContainer.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			envId, err := managedenvironments.ParseManagedEnvironmentID(linuxFunctionAppOnContainer.ManagedEnvironmentId)
			if err != nil {
				return fmt.Errorf("parsing Container App Environment ID for %s: %+v", id, err)
			}

			env, err := containerEnvClient.Get(ctx, *envId)
			if err != nil {
				return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
			}

			subId := commonids.NewSubscriptionID(subscriptionId)

			checkName, err := resourcesClient.CheckNameAvailability(ctx, subId, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for Linux %s: %+v", id, err)
			}
			if model := checkName.Model; model != nil && model.NameAvailable != nil && !*model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *model.Message)
			}

			storageString := linuxFunctionAppOnContainer.StorageAccountName
			if !linuxFunctionAppOnContainer.StorageUsesMSI {
				if linuxFunctionAppOnContainer.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, linuxFunctionAppOnContainer.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, linuxFunctionAppOnContainer.StorageAccountName, linuxFunctionAppOnContainer.StorageAccountKey, *storageDomainSuffix)
				}
			}

			siteConfig := helpers.ExpandSiteConfigLinuxFunctionAppOnContainer(linuxFunctionAppOnContainer.SiteConfig, nil, metadata, linuxFunctionAppOnContainer.Registries[0], linuxFunctionAppOnContainer.FunctionExtensionsVersion, storageString, linuxFunctionAppOnContainer.StorageUsesMSI)
			siteConfig.LinuxFxVersion = helpers.EncodeLinuxFunctionAppOnContainerRegistryImage(linuxFunctionAppOnContainer.Registries, linuxFunctionAppOnContainer.ContainerImage)

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(linuxFunctionAppOnContainer.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(env.Model.Location),
				Kind:     pointer.To("functionapp,linux,container,azurecontainerapps"),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					SiteConfig:           siteConfig,
					ManagedEnvironmentId: pointer.To(linuxFunctionAppOnContainer.ManagedEnvironmentId),
				},
				Tags: pointer.To(linuxFunctionAppOnContainer.Tags),
			}
			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, linuxFunctionAppOnContainer.AppSettings)
			if linuxFunctionAppOnContainer.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(linuxFunctionAppOnContainer.KeyVaultReferenceIdentityID)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, siteEnvelope); err != nil {
				if strings.Contains(err.Error(), "expected either `provisioningState` or `status` to be returned from the LRO API but both were empty") {
					stateConf := &pluginsdk.StateChangeConf{
						Delay:                     5 * time.Minute,
						Pending:                   []string{"204"},
						Target:                    []string{"200"},
						Refresh:                   functionAcaStateRefreshFunc(ctx, client, id),
						MinTimeout:                15 * time.Second,
						ContinuousTargetOccurence: 10,
						Timeout:                   metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate),
					}
					if _, err := stateConf.WaitForStateContext(ctx); err != nil {
						return fmt.Errorf("waiting for creation of %s: %+v", id, err)
					}
				} else {
					return fmt.Errorf("waiting creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LinuxFunctionAppOnContainerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			functionAppOnContainer, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionAppOnContainer.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			if model := functionAppOnContainer.Model; model != nil {
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state := LinuxFunctionAppOnContainerModel{
					Name:          id.SiteName,
					ResourceGroup: id.ResourceGroupName,
					Identity:      pointer.From(flattenedIdentity),
					Tags:          pointer.From(model.Tags),
				}

				if props := model.Properties; props != nil {
					state.ManagedEnvironmentId = pointer.From(props.ManagedEnvironmentId)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					configResp, err := client.GetConfiguration(ctx, *id)
					if err != nil {
						return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
					}

					if configRespModel := configResp.Model; configRespModel != nil && configRespModel.Properties != nil {
						state.Identity = pointer.From(flattenedIdentity)
						siteConfig := helpers.FlattenSiteConfigLinuxFunctionAppOnContainer(configRespModel.Properties)
						state.SiteConfig = []helpers.SiteConfigLinuxFunctionAppOnContainer{*siteConfig}
						state.ContainerImage, state.Registries, err = helpers.DecodeLinuxFunctionAppOnContainerRegistryImage(configRespModel.Properties.LinuxFxVersion, appSettingsResp.Model)
						if err != nil {
							return fmt.Errorf("flattening linuxFxVersion: %s", err)
						}
					}

					state.unpackLinuxFunctionAppOnContainerAppSettings(*appSettingsResp.Model, metadata)
				}
				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (r LinuxFunctionAppOnContainerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			metadata.Logger.Infof("deleting Linux %s", *id)

			delOptions := webapps.DeleteOperationOptions{
				DeleteEmptyServerFarm: pointer.To(false),
				DeleteMetrics:         pointer.To(false),
			}
			if _, err := client.Delete(ctx, *id, delOptions); err != nil {
				return fmt.Errorf("deleting Linux %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r LinuxFunctionAppOnContainerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LinuxFunctionAppOnContainerModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Linux %s: %v", id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("reading of Linux %s for update", id)
			}
			model := *existing.Model

			// below properties is not supported
			model.Properties.DefaultHostName = nil
			model.Properties.State = nil
			model.Properties.ResourceConfig = nil

			storageString := state.StorageAccountName
			if !state.StorageUsesMSI {
				if state.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, state.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, state.StorageAccountName, state.StorageAccountKey, *storageDomainSuffix)
				}
			}

			siteConfig := helpers.ExpandSiteConfigLinuxFunctionAppOnContainer(state.SiteConfig, model.Properties.SiteConfig, metadata, state.Registries[0], state.FunctionExtensionsVersion, storageString, state.StorageUsesMSI)
			if metadata.ResourceData.HasChange("site_config") {
				model.Properties.SiteConfig = siteConfig
			}

			model.Properties.SiteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, state.AppSettings)

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("registry") || metadata.ResourceData.HasChange("container_image") {
				model.Properties.SiteConfig.LinuxFxVersion = helpers.EncodeLinuxFunctionAppOnContainerRegistryImage(state.Registries, state.ContainerImage)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				if strings.Contains(err.Error(), "expected either `provisioningState` or `status` to be returned from the LRO API but both were empty") {
					stateConf := &pluginsdk.StateChangeConf{
						Delay:                     5 * time.Minute,
						Pending:                   []string{"204"},
						Target:                    []string{"200"},
						Refresh:                   functionAcaStateRefreshFunc(ctx, client, *id),
						MinTimeout:                15 * time.Second,
						ContinuousTargetOccurence: 10,
						Timeout:                   metadata.ResourceData.Timeout(pluginsdk.TimeoutUpdate),
					}
					if _, err := stateConf.WaitForStateContext(ctx); err != nil {
						return fmt.Errorf("waiting for update of %s: %+v", id, err)
					}
				} else {
					return fmt.Errorf("waiting for update of %s: %+v", id, err)
				}
			}

			if _, err := client.UpdateConfiguration(ctx, *id, webapps.SiteConfigResource{Properties: model.Properties.SiteConfig}); err != nil {
				return fmt.Errorf("updating Site Config for Linux %s: %s", id, err)
			}

			return nil
		},
	}
}

func (m *LinuxFunctionAppOnContainerModel) unpackLinuxFunctionAppOnContainerAppSettings(input webapps.StringDictionary, metadata sdk.ResourceMetaData) {
	if input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)

	var dockerSettings helpers.Registry

	for k, v := range *input.Properties {
		switch k {
		case "WEBSITE_AUTH_ENCRYPTION_KEY":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_AUTH_ENCRYPTION_KEY"); ok {
				appSettings[k] = v
			}

		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = v

		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"); ok {
				appSettings[k] = v
			}

		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.Server = v
		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.UserName = v
		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.Password = v
		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v
		case "AzureWebJobsStorage":
			if strings.HasPrefix(v, "@Microsoft.KeyVault") {
				trimmed := strings.TrimPrefix(strings.TrimSuffix(v, ")"), "@Microsoft.KeyVault(SecretUri=")
				m.StorageKeyVaultSecretID = trimmed
			} else {
				m.StorageAccountName, m.StorageAccountKey = helpers.ParseWebJobsStorageString(v)
			}
		case "AzureWebJobsStorage__accountName":
			m.StorageUsesMSI = true
			m.StorageAccountName = v

		default:
			appSettings[k] = v
		}
	}
	m.AppSettings = appSettings
}

func functionAcaStateRefreshFunc(ctx context.Context, client *webapps.WebAppsClient, id commonids.AppServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
	}
}
