package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LinuxFunctionAppResource struct{}

type LinuxFunctionAppModel struct {
	Name               string `tfschema:"name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	ServicePlanId      string `tfschema:"service_plan_id"`
	StorageAccountName string `tfschema:"storage_account_name"`

	StorageAccountKey string `tfschema:"storage_account_access_key"`
	StorageUsesMSI    bool   `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key

	AppSettings               map[string]string                    `tfschema:"app_settings"`
	AuthSettings              []helpers.AuthSettings               `tfschema:"auth_settings"`
	Backup                    []helpers.Backup                     `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging            bool                                 `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled         bool                                 `tfschema:"client_cert_enabled"`
	ClientCertMode            string                               `tfschema:"client_cert_mode"`
	ConnectionStrings         []helpers.ConnectionString           `tfschema:"connection_string"`
	Enabled                   bool                                 `tfschema:"enabled"`
	FunctionExtensionsVersion string                               `tfschema:"functions_extension_version"`
	ForceDisableContentShare  bool                                 `tfschema:"force_disable_content_share"`
	HttpsOnly                 bool                                 `tfschema:"https_only"`
	Identity                  []helpers.Identity                   `tfschema:"identity"`
	SiteConfig                []helpers.SiteConfigLinuxFunctionApp `tfschema:"site_config"`
	Tags                      map[string]string                    `tfschema:"tags"`

	// Computed
	CustomDomainVerificationId    string                   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string                   `tfschema:"default_hostname"`
	Kind                          string                   `tfschema:"kind"`
	OutboundIPAddresses           string                   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string                 `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string                   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string                 `tfschema:"possible_outbound_ip_address_list"`
	SiteCredentials               []helpers.SiteCredential `tfschema:"site_credential"`

	// TODO - WEBSITE_RUN_FROM_PACKAGE (https://docs.microsoft.com/en-us/azure/azure-functions/run-functions-from-deployment-package) ??
	DailyMemoryTimeQuota int `tfschema:"daily_memory_time_quota"` // TODO - Value ignored in for linux apps, even in Consumption plans?
}

var _ sdk.ResourceWithUpdate = LinuxFunctionAppResource{}

var _ sdk.ResourceWithCustomImporter = LinuxFunctionAppResource{}

var _ sdk.ResourceWithCustomizeDiff = LinuxFunctionAppResource{}

func (r LinuxFunctionAppResource) ModelObject() interface{} {
	return &LinuxFunctionAppModel{}
}

func (r LinuxFunctionAppResource) ResourceType() string {
	return "azurerm_linux_function_app"
}

func (r LinuxFunctionAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FunctionAppID
}

func (r LinuxFunctionAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": location.Schema(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanID,
			// TODO - CustomizeDiff for ForceNew when Consumption Plan
			// Probably needs locking on the plan while this completes or core will see it as orphaned?
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true, // TODO - Check this?
			ValidateFunc: storageValidate.StorageAccountName,
		},

		"storage_account_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true, // TODO - Uncomment this
			ValidateFunc: validation.NoZeroValues,
			ExactlyOneOf: []string{
				"storage_uses_managed_identity",
				"storage_account_access_key",
			},
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ExactlyOneOf: []string{
				"storage_uses_managed_identity",
				"storage_account_access_key",
			},
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"backup": helpers.BackupSchema(),

		"builtin_logging_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_cert_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  web.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice([]string{
				string(web.ClientCertModeOptional),
				string(web.ClientCertModeRequired),
				string(web.ClientCertModeOptionalInteractiveUser),
			}, false),
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"daily_memory_time_quota": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"force_disable_content_share": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"functions_extension_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "~3",
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"identity": helpers.IdentitySchema(),

		"site_config": helpers.SiteConfigSchemaLinuxFunctionApp(),

		"tags": tags.Schema(),
	}
}

func (r LinuxFunctionAppResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r LinuxFunctionAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var functionApp LinuxFunctionAppModel

			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewFunctionAppID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)

			servicePlanId, err := parse.ServicePlanID(functionApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			sendContentSettings := !functionApp.ForceDisableContentShare
			if planSku := servicePlan.Sku; planSku != nil && planSku.Tier != nil {
				switch tier := *planSku.Tier; strings.ToLower(tier) {
				case "dynamic": // Consumption Plan modifications to request
					sendContentSettings = false
				case "elastic": // ElasticPremium Plan modifications to request?
				case "basic": // App Service Plan modifications to request?
					sendContentSettings = false
				case "standard":
					sendContentSettings = false
				case "premiumv2", "premiumv3":
					sendContentSettings = false
				}
			} else {
				return fmt.Errorf("determining plan type for Linux %s: %v", id, err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			availabilityRequest := web.ResourceNameAvailabilityRequest{
				Name: utils.String(functionApp.Name),
				Type: web.CheckNameResourceTypesMicrosoftWebsites,
			}

			if ase := servicePlan.HostingEnvironmentProfile; ase != nil {
				// Attempt to check the ASE for the appropriate suffix for the name availability request.
				// This varies between internal and external ASE Types, and potentially has other names in other clouds
				// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
				nameSuffix := "appserviceenvironment.net"
				if ase.ID != nil {
					aseId, err := parse.AppServiceEnvironmentID(*ase.ID)
					nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
					if err != nil {
						metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
					} else {
						existingASE, err := aseClient.Get(ctx, aseId.ResourceGroup, aseId.HostingEnvironmentName)
						if err != nil {
							metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
						} else if props := existingASE.AppServiceEnvironment; props != nil && props.DNSSuffix != nil && *props.DNSSuffix != "" {
							nameSuffix = *props.DNSSuffix
						}
					}
				}

				availabilityRequest.Name = utils.String(fmt.Sprintf("%s.%s", functionApp.Name, nameSuffix))
				availabilityRequest.IsFqdn = utils.Bool(true)
			}

			checkName, err := client.CheckNameAvailability(ctx, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for Linux %s: %+v", id, err)
			}
			if checkName.NameAvailable != nil && !*checkName.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Message)
			}

			storageString := functionApp.StorageAccountName
			if !functionApp.StorageUsesMSI {
				storageString = fmt.Sprintf(helpers.StorageStringFmt, functionApp.StorageAccountName, functionApp.StorageAccountKey, metadata.Client.Account.Environment.StorageEndpointSuffix)
			}
			siteConfig, err := helpers.ExpandSiteConfigLinuxFunctionApp(functionApp.SiteConfig, nil, metadata, functionApp.FunctionExtensionsVersion, storageString, functionApp.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding site_config for Linux %s: %+v", id, err)
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
				suffix := uuid.New().String()[0:4]
				functionApp.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(functionApp.Name), suffix)
				functionApp.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
			}

			siteConfig.LinuxFxVersion = helpers.EncodeFunctionAppLinuxFxVersion(functionApp.SiteConfig[0].ApplicationStack)
			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionApp.AppSettings)

			siteEnvelope := web.Site{
				Location: utils.String(functionApp.Location),
				Tags:     tags.FromTypedObject(functionApp.Tags),
				Kind:     utils.String("functionapp,linux"),
				SiteProperties: &web.SiteProperties{
					ServerFarmID:         utils.String(functionApp.ServicePlanId),
					Enabled:              utils.Bool(functionApp.Enabled),
					HTTPSOnly:            utils.Bool(functionApp.HttpsOnly),
					SiteConfig:           siteConfig,
					ClientCertEnabled:    utils.Bool(functionApp.ClientCertEnabled),
					ClientCertMode:       web.ClientCertMode(functionApp.ClientCertMode),
					DailyMemoryTimeQuota: utils.Int32(int32(functionApp.DailyMemoryTimeQuota)), // TODO - Investigate, setting appears silently ignored on Linux Function Apps?
				},
			}

			if identity := helpers.ExpandIdentity(functionApp.Identity); identity != nil {
				siteEnvelope.Identity = identity
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("creating Linux %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Linux %s: %+v", id, err)
			}

			updateFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("updating properties of Linux %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Linux %s: %+v", id, err)
			}

			backupConfig := helpers.ExpandBackupConfig(functionApp.Backup)
			if backupConfig.BackupRequestProperties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Linux %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionApp.AuthSettings)
			if auth.SiteAuthSettingsProperties != nil {
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Linux %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LinuxFunctionAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionApp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(functionApp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			if functionApp.SiteProperties == nil {
				return fmt.Errorf("reading properties of Linux %s", id)
			}
			props := *functionApp.SiteProperties

			appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Linux %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for Linux %s: %+v", id, err)
			}

			if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Site Publishing Credential information for Linux %s: %+v", id, err)
			}
			siteCredentials, err := siteCredentialsFuture.Result(*client)
			if err != nil {
				return fmt.Errorf("reading Site Publishing Credential information for Linux %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Linux %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Linux %s: %+v", id, err)
				}
			}

			state := LinuxFunctionAppModel{
				Name:                 id.SiteName,
				ResourceGroup:        id.ResourceGroup,
				ServicePlanId:        utils.NormalizeNilableString(props.ServerFarmID),
				Location:             location.NormalizeNilable(functionApp.Location),
				ClientCertMode:       string(functionApp.ClientCertMode),
				DailyMemoryTimeQuota: int(utils.NormaliseNilableInt32(props.DailyMemoryTimeQuota)),
				Tags:                 tags.ToTypedObject(functionApp.Tags),
				Kind:                 utils.NormalizeNilableString(functionApp.Kind),
			}

			if functionApp.Enabled != nil {
				state.Enabled = *functionApp.Enabled
			}

			if identity := helpers.FlattenIdentity(functionApp.Identity); identity != nil {
				state.Identity = identity
			}

			configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
			}

			siteConfig, err := helpers.FlattenSiteConfigLinuxFunctionApp(configResp.SiteConfig)
			if err != nil {
				return fmt.Errorf("reading Site Config for Linux %s: %+v", id, err)
			}
			state.SiteConfig = []helpers.SiteConfigLinuxFunctionApp{*siteConfig}

			state.unpackLinuxFunctionAppSettings(appSettingsResp)

			state.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			state.AuthSettings = helpers.FlattenAuthSettings(auth)

			state.Backup = helpers.FlattenBackupConfig(backup)

			if functionApp.HTTPSOnly != nil {
				state.HttpsOnly = *functionApp.HTTPSOnly
			}

			if functionApp.ClientCertEnabled != nil {
				state.ClientCertEnabled = *functionApp.ClientCertEnabled
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LinuxFunctionAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Linux %s", *id)

			deleteMetrics := true
			deleteEmptyServerFarm := false
			if _, err := client.Delete(ctx, id.ResourceGroup, id.SiteName, &deleteMetrics, &deleteEmptyServerFarm); err != nil {
				return fmt.Errorf("deleting Linux %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r LinuxFunctionAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.FunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LinuxFunctionAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Linux %s: %v", id, err)
			}

			// Some service plan updates are allowed - see customiseDiff for exceptions
			if metadata.ResourceData.HasChange("service_plan_id") {
				existing.SiteProperties.ServerFarmID = utils.String(state.ServicePlanId)
			}

			if metadata.ResourceData.HasChange("enabled") {
				existing.SiteProperties.Enabled = utils.Bool(state.Enabled)
			}

			if metadata.ResourceData.HasChange("https_only") {
				existing.SiteProperties.HTTPSOnly = utils.Bool(state.HttpsOnly)
			}

			if metadata.ResourceData.HasChange("client_cert_enabled") {
				existing.SiteProperties.ClientCertEnabled = utils.Bool(state.ClientCertEnabled)
			}

			if metadata.ResourceData.HasChange("client_cert_mode") {
				existing.SiteProperties.ClientCertMode = web.ClientCertMode(state.ClientCertMode)
			}

			if metadata.ResourceData.HasChange("identity") {
				existing.Identity = helpers.ExpandIdentity(state.Identity)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			storageString := ""
			if !state.StorageUsesMSI {
				storageString = fmt.Sprintf(helpers.StorageStringFmt, state.StorageAccountName, state.StorageAccountKey, metadata.Client.Account.Environment.StorageEndpointSuffix)
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigLinuxFunctionApp(state.SiteConfig, existing.SiteConfig, metadata, state.FunctionExtensionsVersion, storageString, state.StorageUsesMSI)
			if state.BuiltinLogging {
				if state.AppSettings == nil && !state.StorageUsesMSI {
					state.AppSettings = make(map[string]string)
				}
				state.AppSettings["AzureWebJobsDashboard"] = storageString
			}

			if metadata.ResourceData.HasChange("site_config") {
				if err != nil {
					return fmt.Errorf("expanding Site Config for Linux %s: %+v", id, err)
				}
				existing.SiteConfig = siteConfig
			}

			if metadata.ResourceData.HasChange("site_config.0.application_stack") {
				existing.SiteConfig.LinuxFxVersion = helpers.EncodeFunctionAppLinuxFxVersion(state.SiteConfig[0].ApplicationStack)
			}

			existing.SiteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, state.AppSettings)

			updateFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, existing)
			if err != nil {
				return fmt.Errorf("updating Linux %s: %+v", id, err)
			}
			if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to update %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = map[string]*web.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStrings(ctx, id.ResourceGroup, id.SiteName, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				if _, err := client.UpdateAuthSettings(ctx, id.ResourceGroup, id.SiteName, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate := helpers.ExpandBackupConfig(state.Backup)
				if backupUpdate.BackupRequestProperties == nil {
					if _, err := client.DeleteBackupConfiguration(ctx, id.ResourceGroup, id.SiteName); err != nil {
						return fmt.Errorf("removing Backup Settings for Linux %s: %+v", id, err)
					}
				} else {
					if _, err := client.UpdateBackupConfiguration(ctx, id.ResourceGroup, id.SiteName, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Linux %s: %+v", id, err)
					}
				}
			}

			return nil
		},
	}
}

func (r LinuxFunctionAppResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.AppService.WebAppsClient
		servicePlanClient := metadata.Client.AppService.ServicePlanClient

		id, err := parse.FunctionAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
		if err != nil || site.SiteProperties == nil {
			return fmt.Errorf("reading Linux %s: %+v", id, err)
		}
		props := site.SiteProperties
		if props.ServerFarmID == nil {
			return fmt.Errorf("determining Service Plan ID for Linux %s: %+v", id, err)
		}
		servicePlanId, err := parse.ServicePlanID(*props.ServerFarmID)
		if err != nil {
			return err
		}

		sp, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
		if err != nil || sp.Kind == nil {
			return fmt.Errorf("reading Service Plan for Linux %s: %+v", id, err)
		}
		if !strings.Contains(strings.ToLower(*sp.Kind), "linux") && !strings.Contains(strings.ToLower(*sp.Kind), "elastic") && !strings.Contains(strings.ToLower(*sp.Kind), "functionapp") {
			return fmt.Errorf("specified Service Plan is not a Linux Functionapp plan")
		}

		return nil
	}
}

func (r LinuxFunctionAppResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			rd := metadata.ResourceDiff

			if rd.HasChange("service_plan_id") {
				currentPlanIdRaw, newPlanIdRaw := rd.GetChange("service_plan_id")
				newPlanId, err := parse.ServicePlanID(newPlanIdRaw.(string))
				if err != nil {
					return fmt.Errorf("reading new plan id %+v", err)
				}

				var currentTierIsDynamic, newTierIsDynamic, newTierIsBasic bool

				newPlan, err := client.Get(ctx, newPlanId.ResourceGroup, newPlanId.ServerfarmName)
				if err != nil {
					return fmt.Errorf("could not read new Service Plan to check tier %s: %+v", newPlanId, err)
				}
				if planSku := newPlan.Sku; planSku != nil {
					if tier := planSku.Tier; tier != nil {
						newTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						newTierIsBasic = strings.EqualFold(*tier, "basic")
					}
				}

				// Service Plans can only be updated in place when both New and Existing are not Dynamic
				if currentPlanIdRaw.(string) != "" {
					currentPlanId, err := parse.ServicePlanID(currentPlanIdRaw.(string))
					if err != nil {
						return fmt.Errorf("reading existing plan id %+v", err)
					}

					currentPlan, err := client.Get(ctx, currentPlanId.ResourceGroup, currentPlanId.ServerfarmName)
					if err != nil {
						return fmt.Errorf("could not read old Service Plan to check tier %s: %+v", currentPlanId, err)
					}

					if planSku := currentPlan.Sku; planSku != nil {
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
				if _, ok := rd.GetOk("backup"); ok && newTierIsDynamic {
					return fmt.Errorf("cannot specify backup configuration for Dynamic tier Service Plans, Standard or higher is required")
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsBasic {
					return fmt.Errorf("cannot specify backup configuration for Basic tier Service Plans, Standard or higher is required")
				}
			}
			return nil
		},
	}
}

func (m *LinuxFunctionAppModel) unpackLinuxFunctionAppSettings(input web.StringDictionary) {
	if input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)
	var appStack helpers.ApplicationStackLinuxFunctionApp
	var dockerSettings helpers.ApplicationStackDocker
	m.BuiltinLogging = false

	for k, v := range input.Properties {
		switch k {
		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = utils.NormalizeNilableString(v)

		case "WEBSITE_NODE_DEFAULT_VERSION": // Note - This is only set if it's not the default of 12, but we collect it from LinuxFxVersion so can discard it here
		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
		case "WEBSITE_CONTENTSHARE":
		case "FUNCTIONS_WORKER_RUNTIME":
			appStack.CustomHandler = strings.EqualFold(*v, "custom")

		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.RegistryURL = utils.NormalizeNilableString(v)

		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.RegistryUsername = utils.NormalizeNilableString(v)

		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.RegistryPassword = utils.NormalizeNilableString(v)

		case "WEBSITES_ENABLE_APP_SERVICE_STORAGE": // TODO - Support this as a configurable bool, default `false` - Ref: https://docs.microsoft.com/en-us/azure/app-service/faq-app-service-linux#i-m-using-my-own-custom-container--i-want-the-platform-to-mount-an-smb-share-to-the---home---directory-

		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = utils.NormalizeNilableString(v)

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = utils.NormalizeNilableString(v)

		case "AzureWebJobsStorage":
			m.StorageAccountName, m.StorageAccountKey = helpers.ParseWebJobsStorageString(v)

		case "AzureWebJobsDashboard":
			m.BuiltinLogging = true
		default:
			appSettings[k] = utils.NormalizeNilableString(v)
		}
	}
	m.AppSettings = appSettings

}
