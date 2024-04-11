// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	kvValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsFunctionAppResource struct{}

type WindowsFunctionAppModel struct {
	Name               string `tfschema:"name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	ServicePlanId      string `tfschema:"service_plan_id"`
	StorageAccountName string `tfschema:"storage_account_name"`

	StorageAccountKey       string `tfschema:"storage_account_access_key"`
	StorageUsesMSI          bool   `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key
	StorageKeyVaultSecretID string `tfschema:"storage_key_vault_secret_id"`

	AppSettings                      map[string]string                      `tfschema:"app_settings"`
	StickySettings                   []helpers.StickySettings               `tfschema:"sticky_settings"`
	AuthSettings                     []helpers.AuthSettings                 `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings               `tfschema:"auth_settings_v2"`
	Backup                           []helpers.Backup                       `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging                   bool                                   `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled                bool                                   `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                                 `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                                 `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                []helpers.ConnectionString             `tfschema:"connection_string"`
	DailyMemoryTimeQuota             int64                                  `tfschema:"daily_memory_time_quota"`
	Enabled                          bool                                   `tfschema:"enabled"`
	FunctionExtensionsVersion        string                                 `tfschema:"functions_extension_version"`
	ForceDisableContentShare         bool                                   `tfschema:"content_share_force_disabled"`
	HttpsOnly                        bool                                   `tfschema:"https_only"`
	KeyVaultReferenceIdentityID      string                                 `tfschema:"key_vault_reference_identity_id"`
	PublicNetworkAccess              bool                                   `tfschema:"public_network_access_enabled"`
	SiteConfig                       []helpers.SiteConfigWindowsFunctionApp `tfschema:"site_config"`
	StorageAccounts                  []helpers.StorageAccount               `tfschema:"storage_account"`
	Tags                             map[string]string                      `tfschema:"tags"`
	VirtualNetworkSubnetID           string                                 `tfschema:"virtual_network_subnet_id"`
	ZipDeployFile                    string                                 `tfschema:"zip_deploy_file"`
	PublishingDeployBasicAuthEnabled bool                                   `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled    bool                                   `tfschema:"ftp_publish_basic_authentication_enabled"`

	// Computed
	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	HostingEnvId                  string   `tfschema:"hosting_environment_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

var _ sdk.ResourceWithUpdate = WindowsFunctionAppResource{}

var _ sdk.ResourceWithCustomImporter = WindowsFunctionAppResource{}

var _ sdk.ResourceWithCustomizeDiff = WindowsFunctionAppResource{}

var _ sdk.ResourceWithStateMigration = WindowsFunctionAppResource{}

func (r WindowsFunctionAppResource) ModelObject() interface{} {
	return &WindowsFunctionAppModel{}
}

func (r WindowsFunctionAppResource) ResourceType() string {
	return "azurerm_windows_function_app"
}

func (r WindowsFunctionAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateFunctionAppID
}

func (r WindowsFunctionAppResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: commonids.ValidateAppServicePlanID,
			Description:  "The ID of the App Service Plan within which to create this Function App",
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
				"storage_uses_managed_identity",
				"storage_key_vault_secret_id",
			},
			Description: "The access key which will be used to access the storage account for the Function App.",
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

		"backup": helpers.BackupSchema(),

		"builtin_logging_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting",
		},

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the function app use Client Certificates",
		},

		"client_certificate_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  webapps.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.ClientCertModeOptional),
				string(webapps.ClientCertModeRequired),
				string(webapps.ClientCertModeOptionalInteractiveUser),
			}, false),
			Description: "The mode of the Function App's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser` ",
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Paths to exclude when using client certificates, separated by ;",
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"daily_memory_time_quota": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps in Consumption Plans.",
		},

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Is the Windows Function App enabled.",
		},

		"content_share_force_disabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force disable the content share settings.",
		},

		"functions_extension_version": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
			Description: "The runtime version associated with the Function App.",
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Can the Function App only be accessed via HTTPS?",
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			Description:  "The User Assigned Identity to use for Key Vault access.",
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"site_config": helpers.SiteConfigSchemaWindowsFunctionApp(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"tags": tags.Schema(),

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
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Function App. **Note:** Using this value requires `WEBSITE_RUN_FROM_PACKAGE=1` to be set on the App in `app_settings`.",
		},
	}
}

func (r WindowsFunctionAppResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r WindowsFunctionAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var functionApp WindowsFunctionAppModel

			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourceProvidersClient := metadata.Client.AppService.ResourceProvidersClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			baseId := commonids.NewAppServiceID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)
			id, err := commonids.ParseFunctionAppID(baseId.ID())
			if err != nil {
				return err
			}

			servicePlanId, err := commonids.ParseAppServicePlanID(functionApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: functionApp.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			var planSKU *string
			if model := servicePlan.Model; model != nil {
				if sku := model.Sku; sku != nil && sku.Name != nil {
					planSKU = sku.Name
				}
				if model.Properties != nil {
					if ase := model.Properties.HostingEnvironmentProfile; ase != nil {
						// Attempt to check the ASE for the appropriate suffix for the name availability request.
						// This varies between internal and external ASE Types, and potentially has other names in other clouds
						// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
						nameSuffix := "appserviceenvironment.net"
						if ase.Id != nil {
							aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
							nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
							if err != nil {
								metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
							} else {
								existingASE, err := aseClient.Get(ctx, *aseId)
								if err != nil || existingASE.Model == nil {
									metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
								} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
									nameSuffix = *props.DnsSuffix
								}
							}
						}

						availabilityRequest.Name = fmt.Sprintf("%s.%s", functionApp.Name, nameSuffix)
						availabilityRequest.IsFqdn = pointer.To(true)
					}
				}
			}
			// Only send for Dynamic and ElasticPremium
			sendContentSettings := (helpers.PlanIsConsumption(planSKU) || helpers.PlanIsElastic(planSKU)) && !functionApp.ForceDisableContentShare

			existing, err := client.Get(ctx, *id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			subscriptionID := commonids.NewSubscriptionID(subscriptionId)
			checkName, err := resourceProvidersClient.CheckNameAvailability(ctx, subscriptionID, availabilityRequest)
			if err != nil || checkName.Model == nil {
				return fmt.Errorf("checking name availability for Windows %s: %+v", id, err)
			}
			if checkName.Model.NameAvailable != nil && !*checkName.Model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Model.Message)
			}

			storageString := functionApp.StorageAccountName
			if !functionApp.StorageUsesMSI {
				if functionApp.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, functionApp.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, functionApp.StorageAccountName, functionApp.StorageAccountKey, *storageDomainSuffix)
				}
			}

			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionApp(functionApp.SiteConfig, nil, metadata, functionApp.FunctionExtensionsVersion, storageString, functionApp.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding site_config for Windows %s: %+v", id, err)
			}

			if functionApp.BuiltinLogging {
				if functionApp.AppSettings == nil {
					functionApp.AppSettings = make(map[string]string)
				}
				if !functionApp.StorageUsesMSI {
					functionApp.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					functionApp.AppSettings["AzureWebJobsDashboard__accountName"] = functionApp.StorageAccountName
				}
			}

			if sendContentSettings {
				if functionApp.AppSettings == nil {
					functionApp.AppSettings = make(map[string]string)
				}
				if !functionApp.StorageUsesMSI {
					suffix := uuid.New().String()[0:4]
					_, contentOverVnetEnabled := functionApp.AppSettings["WEBSITE_CONTENTOVERVNET"]
					_, contentSharePresent := functionApp.AppSettings["WEBSITE_CONTENTSHARE"]
					if _, contentShareConnectionStringPresent := functionApp.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"]; !contentShareConnectionStringPresent {
						functionApp.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
					}

					if !contentSharePresent {
						if contentOverVnetEnabled {
							return fmt.Errorf("the app_setting WEBSITE_CONTENTSHARE must be specified and set to a valid share when WEBSITE_CONTENTOVERVNET is specified")
						}
						functionApp.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(functionApp.Name), suffix)
					}
				} else {
					if _, present := functionApp.AppSettings["AzureWebJobsStorage__accountName"]; !present {
						functionApp.AppSettings["AzureWebJobsStorage__accountName"] = storageString
					}
				}
			}

			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionApp.AppSettings)

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(functionApp.Location),
				Tags:     pointer.To(functionApp.Tags),
				Kind:     pointer.To("functionapp"),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:         pointer.To(functionApp.ServicePlanId),
					Enabled:              pointer.To(functionApp.Enabled),
					HTTPSOnly:            pointer.To(functionApp.HttpsOnly),
					SiteConfig:           siteConfig,
					ClientCertEnabled:    pointer.To(functionApp.ClientCertEnabled),
					ClientCertMode:       pointer.To(webapps.ClientCertMode(functionApp.ClientCertMode)),
					DailyMemoryTimeQuota: pointer.To(functionApp.DailyMemoryTimeQuota),
					VnetRouteAllEnabled:  siteConfig.VnetRouteAllEnabled,
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !functionApp.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if functionApp.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(functionApp.VirtualNetworkSubnetID)
			}

			if functionApp.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(functionApp.KeyVaultReferenceIdentityID)
			}

			if functionApp.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(functionApp.ClientCertExclusionPaths)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, siteEnvelope); err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			if !functionApp.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if !functionApp.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err := client.UpdateFtpAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, siteEnvelope); err != nil {
				return fmt.Errorf("updating properties of Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			stickySettings := helpers.ExpandStickySettings(functionApp.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Windows %s: %+v", id, err)
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(functionApp.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", id, err)
			}
			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, *id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionApp.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettings(ctx, *id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Windows %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(functionApp.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, *id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Windows %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(functionApp.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageConfig); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if _, ok := metadata.ResourceData.GetOk("site_config.0.app_service_logs"); ok {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(functionApp.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, *id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			if functionApp.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, functionApp.ZipDeployFile); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (r WindowsFunctionAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionApp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionApp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", *id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", *id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Windows %s: %+v", *id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfiguration(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowed(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			if model := functionApp.Model; model != nil {
				state := WindowsFunctionAppModel{
					Name:          id.SiteName,
					ResourceGroup: id.ResourceGroupName,
					Location:      location.Normalize(model.Location),
					Tags:          pointer.From(model.Tags),
					Kind:          utils.NormalizeNilableString(model.Kind),
				}

				if props := model.Properties; props != nil {

					state.Enabled = pointer.From(props.Enabled)
					state.ClientCertMode = string(pointer.From(props.ClientCertMode))
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					state.StickySettings = helpers.FlattenStickySettings(stickySettings.Model.Properties)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

					servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
					if err != nil {
						return err
					}
					state.ServicePlanId = servicePlanId.ID()

					state.PublishingFTPBasicAuthEnabled = basicAuthFTP
					state.PublishingDeployBasicAuthEnabled = basicAuthWebDeploy

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						hostingEnvId, err := parse.AppServiceEnvironmentIDInsensitively(*hostingEnv.Id)
						if err != nil {
							return err
						}
						state.HostingEnvId = hostingEnvId.ID()
					}

					if v := props.OutboundIPAddresses; v != nil {
						state.OutboundIPAddresses = *v
						state.OutboundIPAddressList = strings.Split(*v, ",")
					}

					if v := props.PossibleOutboundIPAddresses; v != nil {
						state.PossibleOutboundIPAddresses = *v
						state.PossibleOutboundIPAddressList = strings.Split(*v, ",")
					}

					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						state.VirtualNetworkSubnetID = subnetId
					}

				}
				configResp, err := client.GetConfiguration(ctx, *id)
				if err != nil {
					return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionApp(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
				}

				state.SiteConfig = []helpers.SiteConfigWindowsFunctionApp{*siteConfig}

				state.unpackWindowsFunctionAppSettings(appSettingsResp.Model, metadata)

				state.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings.Model)

				state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

				state.AuthSettings = helpers.FlattenAuthSettings(auth.Model)

				state.AuthV2Settings = helpers.FlattenAuthV2Settings(authV2)

				state.Backup = helpers.FlattenBackupConfig(backup.Model)

				state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs.Model)

				state.StorageAccounts = helpers.FlattenStorageAccounts(storageAccounts.Model)

				// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
				if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
					state.ZipDeployFile = deployFile
				}
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				if err := metadata.ResourceData.Set("identity", flattenedIdentity); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (r WindowsFunctionAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Windows %s", *id)

			delOptions := webapps.DeleteOperationOptions{
				DeleteEmptyServerFarm: pointer.To(false),
				DeleteMetrics:         pointer.To(false),
			}
			if _, err = client.Delete(ctx, *id, delOptions); err != nil {
				return fmt.Errorf("deleting Windows %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WindowsFunctionAppResource) Update() sdk.ResourceFunc {
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

			var state WindowsFunctionAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			model := *existing.Model

			var serviceFarmId string
			if metadata.ResourceData.HasChange("service_plan_id") {
				serviceFarmId = state.ServicePlanId
				existing.Model.Properties.ServerFarmId = pointer.To(serviceFarmId)
			}

			_, planSKU, err := helpers.ServicePlanInfoForApp(ctx, metadata, *id)
			if err != nil {
				return err
			}

			// Some service plan updates are allowed - see customiseDiff for exceptions
			if metadata.ResourceData.HasChange("service_plan_id") {
				model.Properties.ServerFarmId = pointer.To(state.ServicePlanId)
				servicePlanId, err := commonids.ParseAppServicePlanID(state.ServicePlanId)
				if err != nil {
					return err
				}

				servicePlanClient := metadata.Client.AppService.ServicePlanClient
				servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
				if err != nil {
					return fmt.Errorf("reading new service plan (%s) for Windows %s: %+v", servicePlanId, id, err)
				}

				if servicePlan.Model != nil {
					if sku := servicePlan.Model.Sku; sku != nil && sku.Name != nil {
						planSKU = sku.Name
					}
				}
			}

			// Only send for ElasticPremium and consumption plan
			sendContentSettings := (helpers.PlanIsConsumption(planSKU) || helpers.PlanIsElastic(planSKU)) && !state.ForceDisableContentShare

			// Some service plan updates are allowed - see customiseDiff for exceptions
			if metadata.ResourceData.HasChange("service_plan_id") {
				model.Properties.ServerFarmId = pointer.To(state.ServicePlanId)
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(state.Enabled)
			}

			if metadata.ResourceData.HasChange("https_only") {
				model.Properties.HTTPSOnly = pointer.To(state.HttpsOnly)
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
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") {
				model.Properties.KeyVaultReferenceIdentity = pointer.To(state.KeyVaultReferenceIdentityID)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
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

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", *id, err)
				}
			}

			storageString := state.StorageAccountName
			if !state.StorageUsesMSI {
				if state.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, state.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, state.StorageAccountName, state.StorageAccountKey, *storageDomainSuffix)
				}
			}

			if sendContentSettings {
				appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading App Settings for Windows %s: %+v", *id, err)
				}
				if state.AppSettings == nil {
					state.AppSettings = make(map[string]string)
				}
				state.AppSettings = helpers.ParseContentSettings(appSettingsResp.Model, state.AppSettings)

				if !state.StorageUsesMSI {
					suffix := uuid.New().String()[0:4]
					_, contentOverVnetEnabled := state.AppSettings["WEBSITE_CONTENTOVERVNET"]
					_, contentSharePresent := state.AppSettings["WEBSITE_CONTENTSHARE"]
					if _, contentShareConnectionStringPresent := state.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"]; !contentShareConnectionStringPresent {
						state.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
					}

					if !contentSharePresent {
						if contentOverVnetEnabled {
							return fmt.Errorf("the value of WEBSITE_CONTENTSHARE must be set to a predefined share when the storage account is restricted to a virtual network")
						}
						state.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(state.Name), suffix)
					}
				} else {
					if _, present := state.AppSettings["AzureWebJobsStorage__accountName"]; !present {
						state.AppSettings["AzureWebJobsStorage__accountName"] = storageString
					}
				}
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings, so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionApp(state.SiteConfig, model.Properties.SiteConfig, metadata, state.FunctionExtensionsVersion, storageString, state.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding Site Config for Windows %s: %+v", id, err)
			}

			if state.BuiltinLogging {
				if state.AppSettings == nil && !state.StorageUsesMSI {
					state.AppSettings = make(map[string]string)
				}
				if !state.StorageUsesMSI {
					state.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					state.AppSettings["AzureWebJobsDashboard__accountName"] = state.StorageAccountName
				}
			}

			if metadata.ResourceData.HasChange("site_config") {
				model.Properties.SiteConfig = siteConfig
				model.Properties.VnetRouteAllEnabled = model.Properties.SiteConfig.VnetRouteAllEnabled
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

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingFTPBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateFtpAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
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
				return fmt.Errorf("updating Site Config for Windows %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = &map[string]webapps.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
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
					return fmt.Errorf("updating Sticky Settings for Windows %s: %+v", id, err)
				}
			}

			updateLogs := false

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
					updateLogs = true
				}
				if _, err := client.UpdateAuthSettings(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Windows %s: %+v", id, err)
				}
				updateLogs = true
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.Properties == nil {
					if _, err := client.DeleteBackupConfiguration(ctx, *id); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err := client.UpdateBackupConfiguration(ctx, *id, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("site_config.0.app_service_logs") || updateLogs {
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

func (r WindowsFunctionAppResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.AppService.WebAppsClient
		servicePlanClient := metadata.Client.AppService.ServicePlanClient

		id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, *id)
		if err != nil || site.Model == nil || site.Model.Properties == nil {
			return fmt.Errorf("reading Windows %s: %+v", id, err)
		}
		props := site.Model.Properties
		if props.ServerFarmId == nil {
			return fmt.Errorf("determining Service Plan ID for Windows %s: %+v", id, err)
		}
		servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
		if err != nil {
			return err
		}

		sp, err := servicePlanClient.Get(ctx, *servicePlanId)
		if err != nil || sp.Model == nil || sp.Model.Kind == nil {
			return fmt.Errorf("reading Service Plan for Windows %s: %+v", id, err)
		}

		if strings.Contains(strings.ToLower(*sp.Model.Kind), "Windows") {
			return fmt.Errorf("specified Service Plan is not a Windows Functionapp plan")
		}

		return nil
	}
}

func (r WindowsFunctionAppResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			rd := metadata.ResourceDiff

			if rd.HasChange("service_plan_id") {
				currentPlanIdRaw, newPlanIdRaw := rd.GetChange("service_plan_id")
				if newPlanIdRaw.(string) == "" {
					// Plans creating a new service_plan inline will be empty as `Computed` known after apply
					return nil
				}
				newPlanId, err := commonids.ParseAppServicePlanID(newPlanIdRaw.(string))
				if err != nil {
					return fmt.Errorf("reading new plan id %+v", err)
				}

				var currentTierIsDynamic, newTierIsDynamic, newTierIsBasic bool

				newPlan, err := client.Get(ctx, *newPlanId)
				if err != nil || newPlan.Model == nil {
					return fmt.Errorf("could not read new Service Plan to check tier %s: %+v", newPlanId, err)
				}
				if planSku := newPlan.Model.Sku; planSku != nil {
					if tier := planSku.Tier; tier != nil {
						newTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						newTierIsBasic = strings.EqualFold(*tier, "basic")
					}
				}

				if _, ok := rd.GetOk("backup"); ok && newTierIsDynamic {
					return fmt.Errorf("cannot specify backup configuration for Dynamic tier Service Plans, Standard or higher is required")
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsBasic {
					return fmt.Errorf("cannot specify backup configuration for Basic tier Service Plans, Standard or higher is required")
				}

				if strings.EqualFold(currentPlanIdRaw.(string), newPlanIdRaw.(string)) || currentPlanIdRaw == "" {
					// State migration escape for correcting case in serverFarms
					// change of case here will not move the app to a new Service Plan
					return nil
				}
				// Service Plans can only be updated in place when both New and Existing are not Dynamic
				if currentPlanIdRaw.(string) != "" {
					currentPlanId, err := commonids.ParseAppServicePlanID(currentPlanIdRaw.(string))
					if err != nil {
						return fmt.Errorf("reading existing plan id %+v", err)
					}

					currentPlan, err := client.Get(ctx, *currentPlanId)
					if err != nil || currentPlan.Model == nil {
						return fmt.Errorf("could not read current Service Plan to check tier %s: %+v", currentPlanId, err)
					}

					if planSku := currentPlan.Model.Sku; planSku != nil {
						if tier := planSku.Tier; tier != nil {
							currentTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						}
					}

					if currentTierIsDynamic || newTierIsDynamic {
						if err := rd.ForceNew("service_plan_id"); err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
	}
}

func (m *WindowsFunctionAppModel) unpackWindowsFunctionAppSettings(input *webapps.StringDictionary, metadata sdk.ResourceMetaData) {
	if input == nil || input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)
	var dockerSettings helpers.ApplicationStackDocker
	m.BuiltinLogging = false

	for k, v := range *input.Properties {
		switch k {
		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = v

		case "WEBSITE_NODE_DEFAULT_VERSION":
			if len(m.SiteConfig[0].ApplicationStack) == 0 {
				m.SiteConfig[0].ApplicationStack = []helpers.ApplicationStackWindowsFunctionApp{{}}
			}
			m.SiteConfig[0].ApplicationStack[0].NodeVersion = v
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"); ok {
				appSettings[k] = v
			}

		case "WEBSITE_CONTENTSHARE":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTSHARE"); ok {
				appSettings[k] = v
			}

		case "WEBSITE_HTTPLOGGING_RETENTION_DAYS":
		case "FUNCTIONS_WORKER_RUNTIME":
			if _, ok := metadata.ResourceData.GetOk("app_settings.FUNCTIONS_WORKER_RUNTIME"); ok {
				appSettings[k] = v
			}
			switch v {
			case "dotnet-isolated":
				m.SiteConfig[0].ApplicationStack[0].DotNetIsolated = true
			case "custom":
				m.SiteConfig[0].ApplicationStack[0].CustomHandler = true

			}

		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.RegistryURL = v

		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.RegistryUsername = v

		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.RegistryPassword = v

		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = v

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v

		case "AzureWebJobsStorage":
			if strings.HasPrefix(v, "@Microsoft.KeyVault") {
				trimmed := strings.TrimPrefix(strings.TrimSuffix(v, ")"), "@Microsoft.KeyVault(SecretUri=")
				m.StorageKeyVaultSecretID = trimmed
			} else {
				m.StorageAccountName, m.StorageAccountKey = helpers.ParseWebJobsStorageString(v)
			}

		case "AzureWebJobsDashboard":
			m.BuiltinLogging = true

		case "WEBSITE_HEALTHCHECK_MAXPINGFAILURES":
			i, _ := strconv.Atoi(v)
			m.SiteConfig[0].HealthCheckEvictionTime = int64(i)

		case "AzureWebJobsStorage__accountName":
			m.StorageUsesMSI = true
			m.StorageAccountName = v

		case "AzureWebJobsDashboard__accountName":
			m.BuiltinLogging = true

		case "WEBSITE_VNET_ROUTE_ALL":
			// Filter out - handled by site_config setting `vnet_route_all_enabled`

		default:
			appSettings[k] = v
		}
	}

	m.AppSettings = appSettings
}

func (r WindowsFunctionAppResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsFunctionAppV0toV1{},
		},
	}
}
