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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	kvValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WindowsFunctionAppSlotResource struct{}

type WindowsFunctionAppSlotModel struct {
	Name                             string                                     `tfschema:"name"`
	FunctionAppID                    string                                     `tfschema:"function_app_id"`
	ServicePlanID                    string                                     `tfschema:"service_plan_id"`
	StorageAccountName               string                                     `tfschema:"storage_account_name"`
	StorageAccountKey                string                                     `tfschema:"storage_account_access_key"`
	StorageUsesMSI                   bool                                       `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key
	StorageKeyVaultSecretID          string                                     `tfschema:"storage_key_vault_secret_id"`
	AppSettings                      map[string]string                          `tfschema:"app_settings"`
	AuthSettings                     []helpers.AuthSettings                     `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings                   `tfschema:"auth_settings_v2"`
	Backup                           []helpers.Backup                           `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging                   bool                                       `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled                bool                                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                                     `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                                     `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                []helpers.ConnectionString                 `tfschema:"connection_string"`
	DailyMemoryTimeQuota             int64                                      `tfschema:"daily_memory_time_quota"`
	Enabled                          bool                                       `tfschema:"enabled"`
	FunctionExtensionsVersion        string                                     `tfschema:"functions_extension_version"`
	ForceDisableContentShare         bool                                       `tfschema:"content_share_force_disabled"`
	HttpsOnly                        bool                                       `tfschema:"https_only"`
	KeyVaultReferenceIdentityID      string                                     `tfschema:"key_vault_reference_identity_id"`
	PublicNetworkAccess              bool                                       `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled bool                                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled    bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	Identity                         []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	SiteConfig                       []helpers.SiteConfigWindowsFunctionAppSlot `tfschema:"site_config"`
	Tags                             map[string]string                          `tfschema:"tags"`
	CustomDomainVerificationId       string                                     `tfschema:"custom_domain_verification_id"`
	HostingEnvId                     string                                     `tfschema:"hosting_environment_id"`
	DefaultHostname                  string                                     `tfschema:"default_hostname"`
	Kind                             string                                     `tfschema:"kind"`
	OutboundIPAddresses              string                                     `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList            []string                                   `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses      string                                     `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList    []string                                   `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials                  []helpers.SiteCredential                   `tfschema:"site_credential"`
	StorageAccounts                  []helpers.StorageAccount                   `tfschema:"storage_account"`
	VirtualNetworkSubnetID           string                                     `tfschema:"virtual_network_subnet_id"`
	VnetImagePullEnabled             bool                                       `tfschema:"vnet_image_pull_enabled,addedInNextMajorVersion"`
}

var _ sdk.ResourceWithUpdate = WindowsFunctionAppSlotResource{}

var _ sdk.ResourceWithStateMigration = WindowsFunctionAppSlotResource{}

func (r WindowsFunctionAppSlotResource) ModelObject() interface{} {
	return &WindowsFunctionAppSlotModel{}
}

func (r WindowsFunctionAppSlotResource) ResourceType() string {
	return "azurerm_windows_function_app_slot"
}

func (r WindowsFunctionAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateSlotID
}

func (r WindowsFunctionAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Windows Function App Slot.",
		},

		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
			Description:  "The ID of the Windows Function App this Slot is a member of.",
		},

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			Description:  "The backend storage account name which will be used by this Function App Slot.",
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
			Description: "The access key which will be used to access the storage account for the Function App Slot.",
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			},
			Description: "Should the Function App Slot use its Managed Identity to access storage?",
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
			Description: "Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting.",
		},

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the Function App Slot use Client Certificates.",
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
			Description:  "The mode of the Function App Slot's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser`.",
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
			Description: "Is the Windows Function App Slot enabled.",
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
			Description: "The runtime version associated with the Function App Slot.",
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Can the Function App Slot only be accessed via HTTPS?",
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

		"site_config": helpers.SiteConfigSchemaWindowsFunctionAppSlot(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"tags": tags.Schema(),

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},
	}
	if features.FourPointOhBeta() {
		s["vnet_image_pull_enabled"] = &pluginsdk.Schema{
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is container image pull over virtual network enabled? Defaults to `false`.",
		}
	}
	return s
}

func (r WindowsFunctionAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The identifier used by App Service to perform domain ownership verification via DNS TXT record.",
		},

		"default_hostname": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The default hostname of the Windows Function App Slot.",
		},

		"hosting_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Kind value for this Windows Function App Slot.",
		},

		"outbound_ip_addresses": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.",
		},

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of outbound IP addresses. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
		},

		"possible_outbound_ip_addresses": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`. For example `[\"52.23.25.3\", \"52.143.43.12\",\"52.143.43.17\"]`.",
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}

func (r WindowsFunctionAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var functionAppSlot WindowsFunctionAppSlotModel

			if err := metadata.Decode(&functionAppSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourceProvidersClient := metadata.Client.AppService.ResourceProvidersClient
			functionAppId, err := commonids.ParseFunctionAppID(functionAppSlot.FunctionAppID)
			if err != nil {
				return err
			}

			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := webapps.NewSlotID(subscriptionId, functionAppId.ResourceGroupName, functionAppId.SiteName, functionAppSlot.Name)

			functionApp, err := client.Get(ctx, *functionAppId)
			if err != nil || functionApp.Model == nil {
				return fmt.Errorf("retrieving parent Windows %s: %+v", *functionAppId, err)
			}

			var servicePlanId *commonids.AppServicePlanId
			differentServicePlanToParent := false
			if functionApp.Model.Properties == nil || functionApp.Model.Properties.ServerFarmId == nil {
				return fmt.Errorf("could not determine Service Plan ID for %s: %+v", id, err)
			}

			servicePlanId, err = commonids.ParseAppServicePlanIDInsensitively(*functionApp.Model.Properties.ServerFarmId)
			if err != nil {
				return err
			}

			if functionAppSlot.ServicePlanID != "" {
				newServicePlanId, err := commonids.ParseAppServicePlanID(functionAppSlot.ServicePlanID)
				if err != nil {
					return err
				}
				// we only set `service_plan_id` when it differs from the parent `service_plan_id` which is causing issues
				// https://github.com/hashicorp/terraform-provider-azurerm/issues/21024
				// we'll error here if the `service_plan_id` equals the parent `service_plan_id`
				if strings.EqualFold(newServicePlanId.ID(), servicePlanId.ID()) {
					return fmt.Errorf("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App")
				}

				servicePlanId = newServicePlanId
				differentServicePlanToParent = true
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: fmt.Sprintf("%s-%s", id.SiteName, id.SlotName),
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			var planSKU *string
			if model := servicePlan.Model; model != nil {
				if sku := model.Sku; sku != nil && sku.Name != nil {
					planSKU = sku.Name
				} else {
					return fmt.Errorf("could not determine Service Plan SKU type")
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
								metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionAppSlot.Name, servicePlanId)
							} else {
								existingASE, err := aseClient.Get(ctx, *aseId)
								if err != nil || existingASE.Model == nil {
									metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionAppSlot.Name, servicePlanId)
								} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
									nameSuffix = *props.DnsSuffix
								}
							}
						}

						availabilityRequest.Name = fmt.Sprintf("%s.%s", functionAppSlot.Name, nameSuffix)
						availabilityRequest.IsFqdn = pointer.To(true)
						if features.FourPointOhBeta() {
							if !functionAppSlot.VnetImagePullEnabled {
								return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app running in an app service environment")
							}
						}
					}
				}
			}
			// Only send for Dynamic and ElasticPremium
			sendContentSettings := (helpers.PlanIsConsumption(planSKU) || helpers.PlanIsElastic(planSKU)) && !functionAppSlot.ForceDisableContentShare

			existing, err := client.GetSlot(ctx, id)
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

			storageString := functionAppSlot.StorageAccountName
			if !functionAppSlot.StorageUsesMSI {
				if functionAppSlot.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, functionAppSlot.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, functionAppSlot.StorageAccountName, functionAppSlot.StorageAccountKey, *storageDomainSuffix)
				}
			}
			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionAppSlot(functionAppSlot.SiteConfig, nil, metadata, functionAppSlot.FunctionExtensionsVersion, storageString, functionAppSlot.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding site_config for Windows %s: %+v", id, err)
			}

			if functionAppSlot.BuiltinLogging {
				if functionAppSlot.AppSettings == nil {
					functionAppSlot.AppSettings = make(map[string]string)
				}
				if !functionAppSlot.StorageUsesMSI {
					functionAppSlot.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					functionAppSlot.AppSettings["AzureWebJobsDashboard__accountName"] = functionAppSlot.StorageAccountName
				}
			}

			if sendContentSettings {
				if functionAppSlot.AppSettings == nil {
					functionAppSlot.AppSettings = make(map[string]string)
				}
				if !functionAppSlot.StorageUsesMSI {
					suffix := uuid.New().String()[0:4]
					if _, present := functionAppSlot.AppSettings["WEBSITE_CONTENTSHARE"]; !present {
						functionAppSlot.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(functionAppSlot.Name), suffix)
					}
					if _, present := functionAppSlot.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"]; !present {
						functionAppSlot.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
					}
				} else {
					if _, present := functionAppSlot.AppSettings["AzureWebJobsStorage__accountName"]; !present {
						functionAppSlot.AppSettings["AzureWebJobsStorage__accountName"] = storageString
					}
				}
			}

			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionAppSlot.AppSettings)

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(functionAppSlot.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(functionApp.Model.Location),
				Tags:     pointer.To(functionAppSlot.Tags),
				Kind:     pointer.To("functionapp"),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					Enabled:              pointer.To(functionAppSlot.Enabled),
					HTTPSOnly:            pointer.To(functionAppSlot.HttpsOnly),
					SiteConfig:           siteConfig,
					ClientCertEnabled:    pointer.To(functionAppSlot.ClientCertEnabled),
					ClientCertMode:       pointer.To(webapps.ClientCertMode(functionAppSlot.ClientCertMode)),
					DailyMemoryTimeQuota: pointer.To(functionAppSlot.DailyMemoryTimeQuota),
					VnetRouteAllEnabled:  siteConfig.VnetRouteAllEnabled,
				},
			}
			if differentServicePlanToParent {
				siteEnvelope.Properties.ServerFarmId = pointer.To(servicePlanId.ID())
			}

			if features.FourPointOhBeta() {
				siteEnvelope.Properties.VnetImagePullEnabled = pointer.To(functionAppSlot.VnetImagePullEnabled)
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !functionAppSlot.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if functionAppSlot.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(functionAppSlot.VirtualNetworkSubnetID)
				siteEnvelope.Properties.ServerFarmId = pointer.To(servicePlanId.ID())
			}

			if functionAppSlot.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(functionAppSlot.KeyVaultReferenceIdentityID)
			}

			if functionAppSlot.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(functionAppSlot.ClientCertExclusionPaths)
			}

			if err = client.CreateOrUpdateSlotThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			if !functionAppSlot.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err = client.UpdateScmAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if !functionAppSlot.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: false,
					},
				}
				if _, err = client.UpdateFtpAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if err = client.CreateOrUpdateSlotThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("updating properties of Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			backupConfig, err := helpers.ExpandBackupConfig(functionAppSlot.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", id, err)
			}

			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfigurationSlot(ctx, id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionAppSlot.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettingsSlot(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Windows %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(functionAppSlot.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2Slot(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Windows %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(functionAppSlot.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id, *storageConfig); err != nil {
					if err != nil {
						return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", id, err)
					}
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionAppSlot.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStringsSlot(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if _, ok := metadata.ResourceData.GetOk("site_config.0.app_service_logs"); ok {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(functionAppSlot.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r WindowsFunctionAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			functionAppSlot, err := client.GetSlot(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionAppSlot.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			functionAppId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

			appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccountsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentialsSlot(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2Slot(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowedSlot(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowedSlot(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}
			state := WindowsFunctionAppSlotModel{
				Name:                             id.SlotName,
				FunctionAppID:                    functionAppId.ID(),
				PublishingFTPBasicAuthEnabled:    basicAuthFTP,
				PublishingDeployBasicAuthEnabled: basicAuthWebDeploy,
				ConnectionStrings:                helpers.FlattenConnectionStrings(connectionStrings.Model),
				SiteCredentials:                  helpers.FlattenSiteCredentials(siteCredentials),
				AuthSettings:                     helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:                   helpers.FlattenAuthV2Settings(authV2),
				Backup:                           helpers.FlattenBackupConfig(backup.Model),
				StorageAccounts:                  helpers.FlattenStorageAccounts(storageAccounts.Model),
			}
			if model := functionAppSlot.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				state.Kind = pointer.From(model.Kind)
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				state.Identity = pointer.From(flattenedIdentity)

				if props := model.Properties; props != nil {
					state.Enabled = pointer.From(props.Enabled)
					state.ClientCertMode = string(pointer.From(props.ClientCertMode))
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

					if features.FourPointOhBeta() {
						state.VnetImagePullEnabled = pointer.From(props.VnetImagePullEnabled)
					}
					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						state.HostingEnvId = pointer.From(hostingEnv.Id)
					}

					functionApp, err := client.Get(ctx, functionAppId)
					if err != nil {
						return fmt.Errorf("reading parent Function App for Windows %s: %+v", *id, err)
					}
					if functionApp.Model == nil || functionApp.Model.Properties == nil || functionApp.Model.Properties.ServerFarmId == nil {
						return fmt.Errorf("reading parent Function App Service Plan information for Windows %s: %+v", *id, err)
					}
					parentAppFarmId, err := commonids.ParseAppServicePlanIDInsensitively(*functionApp.Model.Properties.ServerFarmId)
					if err != nil {
						return fmt.Errorf("parsing parent Service Plan ID for %s: %+v", *id, err)
					}

					if slotPlanIdRaw := props.ServerFarmId; slotPlanIdRaw != nil && !strings.EqualFold(parentAppFarmId.ID(), *slotPlanIdRaw) {
						slotPlanId, err := commonids.ParseAppServicePlanIDInsensitively(*slotPlanIdRaw)
						if err != nil {
							return fmt.Errorf("parsing Service Plan ID for %s: %+v", *id, err)
						}
						state.ServicePlanID = slotPlanId.ID()
					}

					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						state.VirtualNetworkSubnetID = subnetId
					}

				}

				configResp, err := client.GetConfigurationSlot(ctx, *id)
				if err != nil {
					return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionAppSlot(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
				}
				state.SiteConfig = []helpers.SiteConfigWindowsFunctionAppSlot{*siteConfig}

				state.unpackWindowsFunctionAppSettings(appSettingsResp.Model, metadata)

				state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs.Model)

				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (r WindowsFunctionAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Windows %s", *id)

			delOpts := webapps.DeleteSlotOperationOptions{
				DeleteEmptyServerFarm: pointer.To(false),
				DeleteMetrics:         pointer.To(false),
			}

			if _, err := client.DeleteSlot(ctx, *id, delOpts); err != nil {
				return fmt.Errorf("deleting Windows %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WindowsFunctionAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state WindowsFunctionAppSlotModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetSlot(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			model := *existing.Model

			_, planSKU, err := helpers.ServicePlanInfoForAppSlot(ctx, metadata, *id)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("service_plan_id") {
				o, n := metadata.ResourceData.GetChange("service_plan_id")
				oldPlan, err := commonids.ParseAppServicePlanID(o.(string))
				if err != nil {
					return err
				}

				newPlan, err := commonids.ParseAppServicePlanID(n.(string))
				if err != nil {
					return err
				}
				locks.ByID(oldPlan.ID())
				defer locks.UnlockByID(oldPlan.ID())
				locks.ByID(newPlan.ID())
				defer locks.UnlockByID(newPlan.ID())
				if model.Properties == nil {
					return fmt.Errorf("updating Service Plan for Windows %s: Slot SiteProperties was nil", *id)
				}
				model.Properties.ServerFarmId = pointer.To(newPlan.ID())
			}

			// Only send for Dynamic and ElasticPremium
			sendContentSettings := (helpers.PlanIsConsumption(planSKU) || helpers.PlanIsElastic(planSKU)) && !state.ForceDisableContentShare

			// Some service plan updates are allowed - see customiseDiff for exceptions
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
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
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
					if _, err = client.DeleteSwiftVirtualNetworkSlot(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if metadata.ResourceData.HasChange("vnet_image_pull_enabled") && features.FourPointOhBeta() {
				model.Properties.VnetImagePullEnabled = pointer.To(state.VnetImagePullEnabled)
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err = client.UpdateAzureStorageAccountsSlot(ctx, *id, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", id, err)
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
				appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
				}
				if state.AppSettings == nil {
					state.AppSettings = make(map[string]string)
				}
				state.AppSettings = helpers.ParseContentSettings(appSettingsResp.Model, state.AppSettings)
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings, so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionAppSlot(state.SiteConfig, model.Properties.SiteConfig, metadata, state.FunctionExtensionsVersion, storageString, state.StorageUsesMSI)
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
				if err != nil {
					return fmt.Errorf("expanding Site Config for Windows %s: %+v", id, err)
				}
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

			if err = client.CreateOrUpdateSlotThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingFTPBasicAuthEnabled,
					},
				}
				if _, err = client.UpdateFtpAllowedSlot(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err = client.UpdateScmAllowedSlot(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if _, err = client.UpdateConfigurationSlot(ctx, *id, webapps.SiteConfigResource{Properties: siteConfig}); err != nil {
				return fmt.Errorf("updating Site Config for Windows %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = &map[string]webapps.ConnStringValueTypePair{}
				}
				if _, err = client.UpdateConnectionStringsSlot(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
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
				if _, err = client.UpdateAuthSettingsSlot(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				if _, err := client.UpdateAuthSettingsV2Slot(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.Properties == nil {
					if _, err = client.DeleteBackupConfigurationSlot(ctx, *id); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err = client.UpdateBackupConfigurationSlot(ctx, *id, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("site_config.0.app_service_logs") {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(state.SiteConfig[0].AppServiceLogs)
				if _, err = client.UpdateDiagnosticLogsConfigSlot(ctx, *id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (m *WindowsFunctionAppSlotModel) unpackWindowsFunctionAppSettings(input *webapps.StringDictionary, metadata sdk.ResourceMetaData) {
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

func (r WindowsFunctionAppSlotResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsFunctionAppSlotV0toV1{},
		},
	}
}

func (r WindowsFunctionAppSlotResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			appClient := metadata.Client.AppService.WebAppsClient
			rd := metadata.ResourceDiff
			if rd.HasChange("vnet_image_pull_enabled") {
				appId := rd.Get("function_app_id")
				if appId.(string) == "" {
					return nil
				}
				_, newValue := rd.GetChange("vnet_image_pull_enabled")
				functionAppId, err := commonids.ParseAppServiceID(appId.(string))
				if err != nil {
					return err
				}

				functionApp, err := appClient.Get(ctx, *functionAppId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", functionAppId, err)
				}
				if functionAppModel := functionApp.Model; functionAppModel != nil && functionAppModel.Properties != nil {
					if ase := functionAppModel.Properties.HostingEnvironmentProfile; ase != nil && ase.Id != nil && *(ase.Id) != "" && !newValue.(bool) {
						return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app slot running in an app service environment")
					}
				}
			}
			return nil
		},
	}
}
