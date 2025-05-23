// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/web/2022-09-01/web"
)

const (
	StorageStringFmt = "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s"
)

type FunctionAppFlexConsumptionResource struct{}

type FunctionAppFlexConsumptionModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location      string `tfschema:"location"`
	ServicePlanId string `tfschema:"service_plan_id"`

	Enabled                          bool                       `tfschema:"enabled"`
	AppSettings                      map[string]string          `tfschema:"app_settings"`
	StickySettings                   []helpers.StickySettings   `tfschema:"sticky_settings"`
	AuthSettings                     []helpers.AuthSettings     `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings   `tfschema:"auth_settings_v2"`
	ClientCertEnabled                bool                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                     `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                     `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                []helpers.ConnectionString `tfschema:"connection_string"`
	PublicNetworkAccess              bool                       `tfschema:"public_network_access_enabled"`
	HttpsOnly                        bool                       `tfschema:"https_only"`
	VirtualNetworkSubnetID           string                     `tfschema:"virtual_network_subnet_id"`
	ZipDeployFile                    string                     `tfschema:"zip_deploy_file"`
	PublishingDeployBasicAuthEnabled bool                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`

	StorageContainerType          string                                         `tfschema:"storage_container_type"`
	StorageContainerEndpoint      string                                         `tfschema:"storage_container_endpoint"`
	StorageAuthType               string                                         `tfschema:"storage_authentication_type"`
	StorageAccessKey              string                                         `tfschema:"storage_access_key"`
	StorageUserAssignedIdentityID string                                         `tfschema:"storage_user_assigned_identity_id"`
	RuntimeName                   string                                         `tfschema:"runtime_name"`
	RuntimeVersion                string                                         `tfschema:"runtime_version"`
	MaximumInstanceCount          int64                                          `tfschema:"maximum_instance_count"`
	InstanceMemoryInMB            int64                                          `tfschema:"instance_memory_in_mb"`
	AlwaysReady                   []FunctionAppAlwaysReady                       `tfschema:"always_ready"`
	SiteConfig                    []helpers.SiteConfigFunctionAppFlexConsumption `tfschema:"site_config"`
	Identity                      []identity.ModelSystemAssignedUserAssigned     `tfschema:"identity"`
	Tags                          map[string]string                              `tfschema:"tags"`

	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	HostingEnvId                  string   `tfschema:"hosting_environment_id"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

type FunctionAppAlwaysReady struct {
	Name          string `tfschema:"name"`
	InstanceCount int64  `tfschema:"instance_count"`
}

var _ sdk.ResourceWithUpdate = FunctionAppFlexConsumptionResource{}

func (r FunctionAppFlexConsumptionResource) ModelObject() interface{} {
	return &FunctionAppFlexConsumptionModel{}
}

func (r FunctionAppFlexConsumptionResource) ResourceType() string {
	return "azurerm_function_app_flex_consumption"
}

func (r FunctionAppFlexConsumptionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateFunctionAppID
}

func (r FunctionAppFlexConsumptionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Function App.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
			Description:  "The ID of the App Service Plan within which to create this Function App",
		},

		"storage_container_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.FunctionsDeploymentStorageTypeBlobContainer),
			}, false),
			Description: "The type of the storage container where the function app's code is hosted. Only `blobContainer` is supported currently.",
		},

		"storage_container_endpoint": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The endpoint of the storage container where the function app's code is hosted.",
		},

		"storage_authentication_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.AuthenticationTypeSystemAssignedIdentity),
				string(webapps.AuthenticationTypeStorageAccountConnectionString),
				string(webapps.AuthenticationTypeUserAssignedIdentity),
			}, false),
		},

		"storage_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"runtime_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.RuntimeNameDotnetNegativeisolated),
				string(webapps.RuntimeNameJava),
				string(webapps.RuntimeNameNode),
				string(webapps.RuntimeNamePowershell),
				string(webapps.RuntimeNamePython),
				string(webapps.RuntimeNameCustom),
			}, false),
		},

		"runtime_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"instance_memory_in_mb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  2048,
		},

		"maximum_instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.IntBetween(40, 1000),
		},

		// the name is always being lower-cased by the api: https://github.com/Azure/azure-rest-api-specs/issues/33095
		"always_ready": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"instance_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
				},
			},
		},

		"site_config": helpers.SiteConfigSchemaFunctionAppFlexConsumption(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A map of key-value pairs for [App Settings](https://docs.microsoft.com/en-us/azure/azure-functions/functions-app-settings) and custom values.",
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the function app use Client Certificates",
		},

		"client_certificate_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  web.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice([]string{
				string(web.ClientCertModeOptional),
				string(web.ClientCertModeRequired),
				string(web.ClientCertModeOptionalInteractiveUser),
			}, false),
			Description: "The mode of the Function App's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser` ",
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Paths to exclude when using client certificates, separated by ;",
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Is the Function App enabled.",
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Can the Function App only be accessed via HTTPS?",
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Function App. **Note:** Using this value requires either `WEBSITE_RUN_FROM_PACKAGE=1` or `SCM_DO_BUILD_DURING_DEPLOYMENT=true` to be set on the App in `app_settings`.",
		},

		"tags": tags.Schema(),
	}
}

func (r FunctionAppFlexConsumptionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}

func (r FunctionAppFlexConsumptionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			resourcesClient := metadata.Client.AppService.ResourceProvidersClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var functionAppFlexConsumption FunctionAppFlexConsumptionModel
			if err := metadata.Decode(&functionAppFlexConsumption); err != nil {
				return err
			}

			id := commonids.NewAppServiceID(subscriptionId, functionAppFlexConsumption.ResourceGroup, functionAppFlexConsumption.Name)

			servicePlanId, err := commonids.ParseAppServicePlanID(functionAppFlexConsumption.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", servicePlanId, err)
			}

			var planSKU *string
			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: functionAppFlexConsumption.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			if servicePlanModel := servicePlan.Model; servicePlanModel != nil {
				if sku := servicePlanModel.Sku; sku != nil && sku.Name != nil {
					planSKU = sku.Name
				}
			}

			isFlexConsumptionSku := helpers.PlanIsFlexConsumption(planSKU)
			if !isFlexConsumptionSku {
				return fmt.Errorf("the sku name is %s which is not valid for a flex consumption function app", *planSKU)
			}

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			checkName, err := resourcesClient.CheckNameAvailability(ctx, commonids.NewSubscriptionID(subscriptionId), availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for %s: %+v", id, err)
			}
			if model := checkName.Model; model != nil && model.NameAvailable != nil && !*model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *model.Message)
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(functionAppFlexConsumption.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			blobContainerType := webapps.FunctionsDeploymentStorageType(functionAppFlexConsumption.StorageContainerType)
			storageDeployment := &webapps.FunctionsDeployment{
				Storage: &webapps.FunctionsDeploymentStorage{
					Type:  &blobContainerType,
					Value: &functionAppFlexConsumption.StorageContainerEndpoint,
				},
			}
			storageAuthType := webapps.AuthenticationType(functionAppFlexConsumption.StorageAuthType)
			storageConnStringForFCApp := "DEPLOYMENT_STORAGE_CONNECTION_STRING"
			endpoint := strings.TrimPrefix(functionAppFlexConsumption.StorageContainerEndpoint, "https://")
			var storageString string
			if storageNameIndex := strings.Index(endpoint, "."); storageNameIndex != -1 {
				storageName := endpoint[:storageNameIndex]
				storageString = fmt.Sprintf(StorageStringFmt, storageName, functionAppFlexConsumption.StorageAccessKey, *storageDomainSuffix)
			} else {
				return fmt.Errorf("retrieving storage container endpoint error, the expected format is https://storagename.blob.core.windows.net/containername, the received value is %s", functionAppFlexConsumption.StorageContainerEndpoint)
			}
			storageAuth := webapps.FunctionsDeploymentStorageAuthentication{
				Type: &storageAuthType,
			}

			if functionAppFlexConsumption.StorageAuthType == string(webapps.AuthenticationTypeStorageAccountConnectionString) {
				if functionAppFlexConsumption.StorageAccessKey == "" {
					return fmt.Errorf("the storage account access key must be specified when using the storage key based access")
				}
			} else {
				storageConnStringForFCApp = ""
				if functionAppFlexConsumption.StorageAuthType == string(webapps.AuthenticationTypeUserAssignedIdentity) {
					if functionAppFlexConsumption.StorageUserAssignedIdentityID == "" {
						return fmt.Errorf("the user assigned identity id must be specified when using the user assigned identity to access the storage account")
					}
					storageAuth.UserAssignedIdentityResourceId = &functionAppFlexConsumption.StorageUserAssignedIdentityID
				}
			}

			storageAuth.StorageAccountConnectionStringName = &storageConnStringForFCApp
			storageDeployment.Storage.Authentication = &storageAuth
			runtimeName := webapps.RuntimeName(functionAppFlexConsumption.RuntimeName)
			runtime := webapps.FunctionsRuntime{
				Name:    &runtimeName,
				Version: &functionAppFlexConsumption.RuntimeVersion,
			}

			alwaysReady, err := ExpandAlwaysReadyConfiguration(functionAppFlexConsumption.AlwaysReady, functionAppFlexConsumption.MaximumInstanceCount)
			if err != nil {
				return fmt.Errorf("expanding `always_ready` for %s: %+v", id, err)
			}

			scaleAndConcurrencyConfig := webapps.FunctionsScaleAndConcurrency{
				AlwaysReady:          alwaysReady,
				InstanceMemoryMB:     &functionAppFlexConsumption.InstanceMemoryInMB,
				MaximumInstanceCount: &functionAppFlexConsumption.MaximumInstanceCount,
			}

			flexFunctionAppConfig := &webapps.FunctionAppConfig{
				Deployment:          storageDeployment,
				Runtime:             &runtime,
				ScaleAndConcurrency: &scaleAndConcurrencyConfig,
			}

			siteConfig, err := helpers.ExpandSiteConfigFunctionFlexConsumptionApp(functionAppFlexConsumption.SiteConfig, nil, metadata, false, storageString, storageConnStringForFCApp)
			if err != nil {
				return fmt.Errorf("expanding `site_config` for %s: %+v", id, err)
			}

			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionAppFlexConsumption.AppSettings)

			siteEnvelope := webapps.Site{
				Location: location.Normalize(functionAppFlexConsumption.Location),
				Tags:     pointer.To(functionAppFlexConsumption.Tags),
				Kind:     pointer.To("functionapp,linux"),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:      pointer.To(functionAppFlexConsumption.ServicePlanId),
					Enabled:           pointer.To(functionAppFlexConsumption.Enabled),
					SiteConfig:        siteConfig,
					HTTPSOnly:         pointer.To(functionAppFlexConsumption.HttpsOnly),
					FunctionAppConfig: flexFunctionAppConfig,
					ClientCertEnabled: pointer.To(functionAppFlexConsumption.ClientCertEnabled),
					ClientCertMode:    pointer.To(webapps.ClientCertMode(functionAppFlexConsumption.ClientCertMode)),
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !functionAppFlexConsumption.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if functionAppFlexConsumption.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(functionAppFlexConsumption.VirtualNetworkSubnetID)
			}

			if functionAppFlexConsumption.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(functionAppFlexConsumption.ClientCertExclusionPaths)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if !functionAppFlexConsumption.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			stickySettings := helpers.ExpandStickySettings(functionAppFlexConsumption.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionAppFlexConsumption.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettings(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(functionAppFlexConsumption.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionAppFlexConsumption.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for %s: %+v", id, err)
				}
			}

			if functionAppFlexConsumption.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, id, functionAppFlexConsumption.ZipDeployFile); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func (r FunctionAppFlexConsumptionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionAppFlexConsumption, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionAppFlexConsumption.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving App Settings for %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Connection String information for %s: %+v", id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
			if err != nil || stickySettings.Model == nil {
				return fmt.Errorf("retrieving Sticky Settings for %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Auth Settings for %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving authV2 settings for %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, *id); err != nil && basicAuthWebDeployResp.Model != nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			model := functionAppFlexConsumption.Model
			if model == nil {
				return fmt.Errorf("function app %s : model is nil", id)
			}
			flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			state := FunctionAppFlexConsumptionModel{
				Name:                             id.SiteName,
				ResourceGroup:                    id.ResourceGroupName,
				Location:                         location.Normalize(model.Location),
				ConnectionStrings:                helpers.FlattenConnectionStrings(connectionStrings.Model),
				StickySettings:                   helpers.FlattenStickySettings(stickySettings.Model.Properties),
				SiteCredentials:                  helpers.FlattenSiteCredentials(siteCredentials),
				AuthSettings:                     helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:                   helpers.FlattenAuthV2Settings(authV2),
				PublishingDeployBasicAuthEnabled: basicAuthWebDeploy,
				Tags:                             pointer.From(model.Tags),
				Kind:                             pointer.From(model.Kind),
				Identity:                         pointer.From(flattenedIdentity),
			}

			if props := model.Properties; props != nil {
				state.Enabled = pointer.From(props.Enabled)
				state.ClientCertMode = string(pointer.From(props.ClientCertMode))
				state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
				state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
				state.DefaultHostname = pointer.From(props.DefaultHostName)
				state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

				servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
				if err != nil {
					return err
				}
				state.ServicePlanId = servicePlanId.ID()

				if v := props.OutboundIPAddresses; v != nil {
					state.OutboundIPAddresses = *v
					state.OutboundIPAddressList = strings.Split(*v, ",")
				}

				if v := props.PossibleOutboundIPAddresses; v != nil {
					state.PossibleOutboundIPAddresses = *v
					state.PossibleOutboundIPAddressList = strings.Split(*v, ",")
				}

				configResp, err := client.GetConfiguration(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigFunctionAppFlexConsumption(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("retrieving Site Config for %s: %+v", id, err)
				}
				state.SiteConfig = []helpers.SiteConfigFunctionAppFlexConsumption{*siteConfig}

				if functionAppConfig := props.FunctionAppConfig; functionAppConfig != nil {
					if faConfigDeployment := functionAppConfig.Deployment; faConfigDeployment != nil && faConfigDeployment.Storage != nil {
						storageConfig := *faConfigDeployment.Storage
						state.StorageContainerType = string(pointer.From(storageConfig.Type))
						state.StorageContainerEndpoint = pointer.From(storageConfig.Value)
						if storageConfig.Authentication != nil && storageConfig.Authentication.Type != nil {
							state.StorageAuthType = string(pointer.From(storageConfig.Authentication.Type))
							if storageConfig.Authentication.UserAssignedIdentityResourceId != nil {
								state.StorageUserAssignedIdentityID = pointer.From(storageConfig.Authentication.UserAssignedIdentityResourceId)
							}
						}
					}

					if faConfigRuntime := functionAppConfig.Runtime; faConfigRuntime != nil {
						state.RuntimeName = string(pointer.From(faConfigRuntime.Name))
						state.RuntimeVersion = pointer.From(faConfigRuntime.Version)
					}

					if faConfigScale := functionAppConfig.ScaleAndConcurrency; faConfigScale != nil {
						state.AlwaysReady = FlattenAlwaysReadyConfiguration(faConfigScale.AlwaysReady)
						state.InstanceMemoryInMB = pointer.From(faConfigScale.InstanceMemoryMB)
						state.MaximumInstanceCount = pointer.From(faConfigScale.MaximumInstanceCount)
					}
				}

				state.unpackFunctionAppFlexConsumptionSettings(*appSettingsResp.Model)

				state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)

				state.HttpsOnly = pointer.From(props.HTTPSOnly)

				if props.VirtualNetworkSubnetId != nil && pointer.From(props.VirtualNetworkSubnetId) != "" {
					subnetId, err := commonids.ParseSubnetID(*props.VirtualNetworkSubnetId)
					if err != nil {
						return err
					}
					state.VirtualNetworkSubnetID = subnetId.ID()
				}

				// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
				if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
					state.ZipDeployFile = deployFile
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r FunctionAppFlexConsumptionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			delOptions := webapps.DeleteOperationOptions{
				DeleteEmptyServerFarm: pointer.To(false),
				DeleteMetrics:         pointer.To(false),
			}
			if _, err = client.Delete(ctx, *id, delOptions); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r FunctionAppFlexConsumptionResource) Update() sdk.ResourceFunc {
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

			var state FunctionAppFlexConsumptionModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(state.Enabled)
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if metadata.ResourceData.HasChange("client_certificate_enabled") {
				model.Properties.ClientCertEnabled = pointer.To(state.ClientCertEnabled)
			}

			if metadata.ResourceData.HasChange("client_certificate_mode") {
				model.Properties.ClientCertMode = pointer.To(webapps.ClientCertMode(state.ClientCertMode))
			}

			if metadata.ResourceData.HasChange("client_certificate_exclusion_paths") {
				model.Properties.ClientCertExclusionPaths = pointer.To(state.ClientCertExclusionPaths)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			var storageString string
			if state.StorageContainerEndpoint != "" || metadata.ResourceData.HasChange("storage_container_endpoint") {
				endpoint := strings.TrimPrefix(state.StorageContainerEndpoint, "https://")
				model.Properties.FunctionAppConfig.Deployment.Storage.Value = pointer.To(state.StorageContainerEndpoint)
				if storageNameIndex := strings.Index(endpoint, "."); storageNameIndex != -1 {
					storageName := endpoint[:storageNameIndex]
					storageString = fmt.Sprintf(StorageStringFmt, storageName, state.StorageAccessKey, *storageDomainSuffix)
				} else {
					return fmt.Errorf("retrieving storage container endpoint error, the expected format is https://storagename.blob.core.windows.net/containername, the received value is %s", state.StorageContainerEndpoint)
				}
			}

			storageConnStringForFCApp := "DEPLOYMENT_STORAGE_CONNECTION_STRING"
			if metadata.ResourceData.HasChange("storage_authentication_type") {
				storageAuthType := webapps.AuthenticationType(state.StorageAuthType)
				storageAuth := webapps.FunctionsDeploymentStorageAuthentication{
					Type: &storageAuthType,
				}
				if state.StorageAuthType == string(webapps.AuthenticationTypeStorageAccountConnectionString) {
					if state.StorageAccessKey == "" {
						return fmt.Errorf("the storage account access key must be specified when using the storage key based access")
					}
					storageAuth.StorageAccountConnectionStringName = pointer.To(storageConnStringForFCApp)
				} else {
					storageConnStringForFCApp = ""
					if state.StorageAuthType == string(webapps.AuthenticationTypeUserAssignedIdentity) {
						if state.StorageUserAssignedIdentityID == "" {
							return fmt.Errorf("the user assigned identity id must be specified when using the user assigned identity to access the storage account")
						}
						storageAuth.UserAssignedIdentityResourceId = &state.StorageUserAssignedIdentityID
					}
				}
				model.Properties.FunctionAppConfig.Deployment.Storage.Authentication = &storageAuth
			}

			if metadata.ResourceData.HasChange("storage_user_assigned_identity_id") {
				model.Properties.FunctionAppConfig.Deployment.Storage.Authentication.UserAssignedIdentityResourceId = &state.StorageUserAssignedIdentityID
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings, so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigFunctionFlexConsumptionApp(state.SiteConfig, model.Properties.SiteConfig, metadata, false, storageString, storageConnStringForFCApp)
			if err != nil {
				return fmt.Errorf("expanding Site Config for %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("site_config") {
				model.Properties.SiteConfig = siteConfig
			}

			if metadata.ResourceData.HasChange("always_ready") {
				arc, err := ExpandAlwaysReadyConfiguration(state.AlwaysReady, state.MaximumInstanceCount)
				if err != nil {
					return fmt.Errorf("expanding `always_ready` for %s: %+v", id, err)
				}
				model.Properties.FunctionAppConfig.ScaleAndConcurrency.AlwaysReady = arc
			}

			if metadata.ResourceData.HasChange("runtime_name") {
				runtimeName := webapps.RuntimeName(state.RuntimeName)
				model.Properties.FunctionAppConfig.Runtime.Name = pointer.To(runtimeName)
			}

			if metadata.ResourceData.HasChange("runtime_version") {
				model.Properties.FunctionAppConfig.Runtime.Version = pointer.To(state.RuntimeVersion)
			}

			model.Properties.SiteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, state.AppSettings)

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				pna := helpers.PublicNetworkAccessEnabled
				if !state.PublicNetworkAccess {
					pna = helpers.PublicNetworkAccessDisabled
				}

				// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
				model.Properties.PublicNetworkAccess = pointer.To(pna)
				model.Properties.SiteConfig.PublicNetworkAccess = model.Properties.PublicNetworkAccess
			}

			if metadata.ResourceData.HasChange("https_only") {
				model.Properties.HTTPSOnly = pointer.To(state.HttpsOnly)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if _, err := client.UpdateConfiguration(ctx, *id, webapps.SiteConfigResource{Properties: model.Properties.SiteConfig}); err != nil {
				return fmt.Errorf("updating Site Config for %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = pointer.To(map[string]webapps.ConnStringValueTypePair{})
				}
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("sticky_settings") {
				emptySlice := make([]string, 0)
				stickySettings := helpers.ExpandStickySettings(state.StickySettings)
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: &webapps.SlotConfigNames{
						AppSettingNames:       &emptySlice,
						ConnectionStringNames: &emptySlice,
					},
				}

				if stickySettings != nil {
					if stickySettings.AppSettingNames != nil {
						stickySettingsUpdate.Properties.AppSettingNames = stickySettings.AppSettingNames
					}
					if stickySettings.ConnectionStringNames != nil {
						stickySettingsUpdate.Properties.ConnectionStringNames = stickySettings.ConnectionStringNames
					}
				}

				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				// (@jackofallops) - in the case of a removal of this block, we need to zero these settings
				if authUpdate.Properties == nil {
					authUpdate.Properties = &webapps.SiteAuthSettingsProperties{
						Enabled:                           pointer.To(false),
						ClientSecret:                      pointer.To(""),
						ClientSecretSettingName:           pointer.To(""),
						ClientSecretCertificateThumbprint: pointer.To(""),
						GoogleClientSecret:                pointer.To(""),
						FacebookAppSecret:                 pointer.To(""),
						GitHubClientSecret:                pointer.To(""),
						TwitterConsumerSecret:             pointer.To(""),
						MicrosoftAccountClientSecret:      pointer.To(""),
					}
				}
				if _, err := client.UpdateAuthSettings(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("site_config.0.app_service_logs") {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(state.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, *id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, state.ZipDeployFile); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (m *FunctionAppFlexConsumptionModel) unpackFunctionAppFlexConsumptionSettings(input webapps.StringDictionary) {
	if input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)

	for k, v := range *input.Properties {
		switch k {
		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = v

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v

		case "AzureWebJobsStorage":
			_, m.StorageAccessKey = helpers.ParseWebJobsStorageString(v)

		case "WEBSITE_HEALTHCHECK_MAXPINGFAILURES":
			i, _ := strconv.Atoi(v)
			m.SiteConfig[0].HealthCheckEvictionTime = int64(i)

		case "DEPLOYMENT_STORAGE_CONNECTION_STRING":
			// Filter out - not user faced

		default:
			appSettings[k] = v
		}
	}
	m.AppSettings = appSettings
}

func ExpandAlwaysReadyConfiguration(input []FunctionAppAlwaysReady, maximumInstanceCount int64) (*[]webapps.FunctionsAlwaysReadyConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}
	var totalInstanceCount int64
	arList := make([]webapps.FunctionsAlwaysReadyConfig, 0)
	for _, v := range input {
		totalInstanceCount += v.InstanceCount
		arList = append(arList, webapps.FunctionsAlwaysReadyConfig{
			Name:          pointer.To(v.Name),
			InstanceCount: pointer.To(v.InstanceCount),
		})
	}

	if totalInstanceCount > maximumInstanceCount {
		return nil, fmt.Errorf("the total number of always-ready instances should not exceed the maximum scale out limit")
	}

	return &arList, nil
}

func FlattenAlwaysReadyConfiguration(alwaysReady *[]webapps.FunctionsAlwaysReadyConfig) []FunctionAppAlwaysReady {
	if alwaysReady == nil || len(*alwaysReady) == 0 {
		return []FunctionAppAlwaysReady{}
	}

	alwaysReadyList := make([]FunctionAppAlwaysReady, 0)

	for _, v := range *alwaysReady {
		alwaysReadyList = append(alwaysReadyList, FunctionAppAlwaysReady{
			Name:          pointer.From(v.Name),
			InstanceCount: pointer.From(v.InstanceCount),
		})
	}

	return alwaysReadyList
}
