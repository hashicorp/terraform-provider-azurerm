package logic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"math"
	"reflect"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/go-azure-helpers/framework/convert"
	"github.com/hashicorp/go-azure-helpers/framework/identity"
	"github.com/hashicorp/go-azure-helpers/framework/location"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	rmidentity "github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name fw_logic_app_standard -properties "name,resource_group_name" -service-package-name logic -known-values "subscription_id:data.Subscriptions.Primary"

type FwLogicAppStandardResource struct{}

var _ sdk.FrameworkWrappedResourceWithUpdate = &FwLogicAppStandardResource{}

var _ sdk.FrameworkWrappedResourceWithConfigValidators = &FwLogicAppStandardResource{}

func (r FwLogicAppStandardResource) ModelObject() interface{} {
	return new(FwLogicAppStandardResourceModel)
}

func (r FwLogicAppStandardResource) ResourceType() string {
	return "azurerm_fw_logic_app_standard"
}

func (r FwLogicAppStandardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, metadata sdk.ResourceMetadata) {
	if req.ID == "" {
		resourceIdentity := &FwLogicAppStandardResourceIdentityModel{}
		req.Identity.Get(ctx, resourceIdentity)
		id := pointer.To(commonids.NewAppServiceID(resourceIdentity.SubscriptionId, resourceIdentity.ResourceGroupName, resourceIdentity.Name))
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id.ID())...)
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r FwLogicAppStandardResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			commonschema.Name: schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validate.LogicAppStandardName,
					},
				},
			},

			commonschema.ResourceGroupName: commonschema.ResourceGroupNameAttribute(),

			commonschema.Location: location.LocationAttribute(),

			"app_service_plan_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateAppServicePlanID,
					},
				},
			},

			"app_settings": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				CustomType:  typehelpers.NewMapTypeOf[types.String](ctx),
			},

			"use_extension_bundle": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(true),
			},

			"bundle_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedStringDefault("[1.*, 2.0.0]"),
			},

			"client_affinity_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(false),
			},

			"client_certificate_mode": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringInSlice([]string{
							"Required",
							"Optional",
						}, false),
					},
				},
			},

			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(true),
			},

			"https_only": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(false),
			},

			"storage_account_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: storageValidate.StorageAccountName,
					},
				},
			},

			"storage_account_access_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"public_network_access": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedStringDefault(helpers.PublicNetworkAccessEnabled),
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringInSlice([]string{
							helpers.PublicNetworkAccessEnabled,
							helpers.PublicNetworkAccessDisabled,
						}, false),
					},
				},
			},

			"storage_account_share_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},

			"version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedStringDefault("~4"),
			},

			"virtual_network_subnet_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateSubnetID,
					},
				},
			},

			"custom_domain_verification_id": schema.StringAttribute{
				Computed: true,
			},

			"default_hostname": schema.StringAttribute{
				Computed: true,
			},

			"kind": schema.StringAttribute{
				Computed: true,
			},

			"outbound_ip_addresses": schema.StringAttribute{
				Computed: true,
			},

			"possible_outbound_ip_addresses": schema.StringAttribute{
				Computed: true,
			},

			commonschema.Tags: commonschema.TagsResourceAttribute(ctx),
		},

		Blocks: map[string]schema.Block{
			"identity": identity.IdentityResourceBlockSchema(ctx),

			"site_config": schema.ListNestedBlock{
				// TODO - Computed blocks? This is always returned, defaults below don't help as the len on the list changes from 0 (config) -> 1 (response) :(
				CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppStandardSiteConfigModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"always_on": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"ftps_state": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									webapps.PossibleValuesForFtpsState()...,
								),
							},
						},

						"http2_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"linux_fx_version": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},

						"min_tls_version": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(string(webapps.SupportedTlsVersionsOnePointTwo), string(webapps.SupportedTlsVersionsOnePointThree)),
							},
						},

						"pre_warmed_instance_count": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							Validators: []validator.Int64{
								int64validator.Between(0, 20),
							},
						},

						"scm_use_main_ip_restriction": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"scm_min_tls_version": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: validation.StringInSlice([]string{
										string(webapps.SupportedTlsVersionsOnePointTwo),
									}, false),
								},
							},
						},

						"scm_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									string(webapps.ScmTypeBitbucketGit),
									string(webapps.ScmTypeBitbucketHg),
									string(webapps.ScmTypeCodePlexGit),
									string(webapps.ScmTypeCodePlexHg),
									string(webapps.ScmTypeDropbox),
									string(webapps.ScmTypeExternalGit),
									string(webapps.ScmTypeExternalHg),
									string(webapps.ScmTypeGitHub),
									string(webapps.ScmTypeLocalGit),
									string(webapps.ScmTypeNone),
									string(webapps.ScmTypeOneDrive),
									string(webapps.ScmTypeTfs),
									string(webapps.ScmTypeVSO),
									string(webapps.ScmTypeVSTSRM),
								),
							},
						},

						"use_32_bit_worker_process": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(true),
						},

						"websockets_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"health_check_path": schema.StringAttribute{
							Optional: true,
						},

						"elastic_instance_minimum": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							Validators: []validator.Int64{int64validator.Between(0, 20)},
						},

						"app_scale_limit": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							Validators: []validator.Int64{int64validator.AtLeast(0)},
						},

						"runtime_scale_monitoring_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"dotnet_framework_version": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedStringDefault("v4.0"),
							Validators: []validator.String{
								stringvalidator.OneOf(
									"v4.0",
									"v5.0",
									"v6.0",
									"v8.0",
								),
							},
						},

						"vnet_route_all_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},

						"auto_swap_slot_name": schema.StringAttribute{
							Computed: true,
						},
					},
					Blocks: map[string]schema.Block{
						"cors": schema.ListNestedBlock{
							CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppCORSSettingsModel](ctx),
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"allowed_origins": schema.SetAttribute{
										ElementType:         types.StringType,
										Required:            true,
										Description:         "",
										MarkdownDescription: "",
									},

									"support_credentials": schema.BoolAttribute{
										Optional: true,
										Computed: true,
										Default:  typehelpers.NewWrappedBoolDefault(false),
									},
								},
							},
						},

						"ip_restriction": ipRestrictionCommonSchema(ctx),

						"scm_ip_restriction": ipRestrictionCommonSchema(ctx),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				// PlanModifiers: []planmodifier.List{
				// 	// TODO - Can we use this for computed sub properties?
				// },
			},

			"connection_string": schema.SetNestedBlock{
				CustomType: typehelpers.NewSetNestedObjectTypeOf[FwLogicAppStandardConnectionStringsModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},

						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(webapps.PossibleValuesForConnectionStringType()...),
							},
						},

						"value": schema.StringAttribute{
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

			"site_credential": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppStandardSiteCredentialsModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Computed: true,
						},
						"password": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Description:         "",
				MarkdownDescription: "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
		},
	}
}

func (r FwLogicAppStandardResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.RequiredTogether(
			path.MatchRoot("use_extension_bundle"),
			path.MatchRoot("bundle_version"),
		),
	}
}

func (r FwLogicAppStandardResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse, metadata sdk.ResourceMetadata, decodedPlan interface{}) {
	client := metadata.Client.AppService.WebAppsClient
	resourcesClient := metadata.Client.AppService.ResourceProvidersClient
	env := metadata.Client.Account.Environment
	data := sdk.AssertResourceModelType[FwLogicAppStandardResourceModel](decodedPlan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id := commonids.NewAppServiceID(metadata.SubscriptionId, data.ResourceGroupName.ValueString(), data.Name.ValueString())

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking for presence of existing Logic App: %+v", err), err.Error())
			return
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		metadata.ResourceRequiresImport("azurerm_fw_resource_group", id, resp)
		return
	}

	storageAccountDomainSuffix, ok := env.Storage.DomainSuffix()
	if !ok {
		sdk.SetResponseErrorDiagnostic(resp, "finding storage account domain suffix", fmt.Sprintf("could not determine the domain suffix for storage accounts in environment %q: %+v", env.Name, env.Storage))
		return
	}

	availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
		Name: id.SiteName,
		Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
	}

	available, err := resourcesClient.CheckNameAvailability(ctx, commonids.NewSubscriptionID(metadata.SubscriptionId), availabilityRequest)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking if name %q was available", id.SiteName), err.Error())
		return
	}

	if available.Model == nil || available.Model.NameAvailable == nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking if name %q was available", id.SiteName), "response `model` was nil")
		return
	}

	if !*available.Model.NameAvailable {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("the name %q used for the Logic App Standard needs to be globally unique and isn't available", id.SiteName), pointer.From(available.Model.Message))
		return
	}

	siteConfig := FwLogicAppStandardSiteConfigModel{}
	if !data.SiteConfig.IsNull() && len(data.SiteConfig.Elements()) > 0 {
		siteConfigList := make([]FwLogicAppStandardSiteConfigModel, 0)
		resp.Diagnostics.Append(data.SiteConfig.ListValue.ElementsAs(ctx, &siteConfigList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		siteConfig = siteConfigList[0]
	}

	kind := logicAppKindWorkflowApp
	if !siteConfig.LinuxFXVersion.IsNull() && siteConfig.LinuxFXVersion.String() != "" {
		kind = logicAppKindLinuxContainerWorkflowApp
	}

	appSettingsBase := data.buildBaseAppSettings(*storageAccountDomainSuffix)

	sc, diags := expandSiteConfigLogicAppStandard(ctx, siteConfig)
	if diags.HasError() {
		return
	}
	if sc != nil {
		sc.PublicNetworkAccess = data.PublicNetworkAccess.ValueStringPointer()
	}

	sc.AppSettings = expandAppSettingsLogicAppStandard(ctx, data.AppSettings.MapValue, appSettingsBase...)

	siteEnvelope := webapps.Site{
		Kind:     pointer.To(kind),
		Location: location.Normalize(data.Location.ValueString()),
		Properties: &webapps.SiteProperties{
			ServerFarmId:          data.AppServicePlanID.ValueStringPointer(),
			Enabled:               data.Enabled.ValueBoolPointer(),
			ClientAffinityEnabled: data.ClientAffinityEnabled.ValueBoolPointer(),
			ClientCertEnabled:     pointer.To(data.ClientCertificateMode.ValueString() != ""),
			HTTPSOnly:             data.HTTPSOnly.ValueBoolPointer(),
			PublicNetworkAccess:   data.PublicNetworkAccess.ValueStringPointer(),
			SiteConfig:            sc,
		},
	}

	if !data.Identity.IsNull() {
		ident := &rmidentity.SystemAndUserAssignedMap{}
		identity.ExpandToSystemAndUserAssignedMap(ctx, data.Identity, ident, &diags)
		if diags.HasError() {
			return
		}

		siteEnvelope.Identity = ident
	}

	if data.ClientCertificateMode.String() != "" {
		siteEnvelope.Properties.ClientCertMode = pointer.ToEnum[webapps.ClientCertMode](data.ClientCertificateMode.ValueString())
	}

	if !data.VirtualNetworkSubnetID.IsNull() {
		siteEnvelope.Properties.VirtualNetworkSubnetId = data.VirtualNetworkSubnetID.ValueStringPointer()
	}

	tags, diags := commonschema.ExpandTags(data.Tags.MapValue)
	if diags.HasError() {
		sdk.AppendResponseErrorDiagnostic(resp, diags)
		return
	}

	siteEnvelope.Tags = tags

	if err = client.CreateOrUpdateThenPoll(ctx, id, siteEnvelope); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("creating Logic %s", id), err)
	}

	data.ID = types.StringValue(id.ID())

	// Get computed values to satisfy known after Apply requirements
	existing, err = client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			metadata.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	if model := existing.Model; model != nil {
		data.Kind = types.StringPointerValue(model.Kind)
		if props := model.Properties; props != nil {
			data.OutboundIPAddresses = types.StringPointerValue(props.OutboundIPAddresses)
			data.PossibleOutboundIPAddresses = types.StringPointerValue(props.PossibleOutboundIPAddresses)
			data.DefaultHostname = types.StringPointerValue(props.DefaultHostName)
			data.CustomDomainVerificationID = types.StringPointerValue(props.CustomDomainVerificationId)
		}
		identity.FlattenFromSystemAndUserAssignedMap(ctx, model.Identity, &data.Identity, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

	}

	configResp, err := client.GetConfiguration(ctx, id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving site_config for %s:", id), err)
	}
	if model := configResp.Model; model != nil {
		sc, diags := flattenSiteConfigLogicAppStandard(ctx, model.Properties)
		if diags.HasError() {
			sdk.AppendResponseErrorDiagnostic(resp, diags)
			return
		}
		data.SiteConfig = sc
	}

	readLogicAppStandardAppSettings(ctx, &id, client, data, resp)
}

func (r FwLogicAppStandardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, metadata sdk.ResourceMetadata, decodedState interface{}) {
	client := metadata.Client.AppService.WebAppsClient
	state := sdk.AssertResourceModelType[FwLogicAppStandardResourceModel](decodedState, resp)

	id, err := commonids.ParseLogicAppId(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing azurerm_fw_logic_app_standard ID", err)
		return
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			metadata.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	state.Name = types.StringValue(id.SiteName)
	state.ResourceGroupName = types.StringValue(id.ResourceGroupName)

	if model := existing.Model; model != nil {
		state.Location = types.StringValue(location.Normalize(model.Location))
		state.Kind = types.StringPointerValue(model.Kind)

		identity.FlattenFromSystemAndUserAssignedMap(ctx, model.Identity, &state.Identity, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if props := model.Properties; props != nil {
			servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
			if err != nil {
				sdk.SetResponseErrorDiagnostic(resp, "ID Parsing Error", err)
			}
			state.AppServicePlanID = types.StringValue(servicePlanId.ID())
			state.Enabled = types.BoolPointerValue(props.Enabled)
			state.DefaultHostname = types.StringPointerValue(props.DefaultHostName)
			state.HTTPSOnly = types.BoolPointerValue(props.HTTPSOnly)
			state.OutboundIPAddresses = types.StringPointerValue(props.OutboundIPAddresses)
			state.PossibleOutboundIPAddresses = types.StringPointerValue(props.PossibleOutboundIPAddresses)
			state.ClientAffinityEnabled = types.BoolPointerValue(props.ClientAffinityEnabled)
			state.CustomDomainVerificationID = types.StringPointerValue(props.CustomDomainVerificationId)
			state.VirtualNetworkSubnetID = types.StringPointerValue(props.VirtualNetworkSubnetId)
			state.PublicNetworkAccess = types.StringPointerValue(props.PublicNetworkAccess)
			state.ClientCertificateMode = types.StringNull()
			if state.ClientAffinityEnabled.ValueBool() {
				certMode := string(pointer.From(props.ClientCertMode))
				state.ClientCertificateMode = types.StringValue(certMode)
			}
		}
	}

	configResp, err := client.GetConfiguration(ctx, *id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving site_config for %s:", id), err)
	}
	if model := configResp.Model; model != nil {
		sc, diags := flattenSiteConfigLogicAppStandard(ctx, model.Properties)
		if diags.HasError() {
			sdk.AppendResponseErrorDiagnostic(resp, diags)
			return
		}
		state.SiteConfig = sc
	}

	readLogicAppStandardAppSettings(ctx, id, client, state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, *id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Listing Connection Strings for %s", *id), err)
	}

	if model := connectionStringsResp.Model; model != nil {
		convert.Flatten(ctx, model.Properties, &state.ConnectionStrings, &resp.Diagnostics) // Need to set manually, doesn't match Struct Names
		if resp.Diagnostics.HasError() {
			return
		}
	}

	siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("Listing Publishing Credentials for %s", *id), err)
		return
	}

	convert.Flatten(ctx, siteCredentials.Properties, &state.SiteCredentials, &resp.Diagnostics) // Need to set manually, doesn't match Struct Names - try the struct tagging?
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r FwLogicAppStandardResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse, metadata sdk.ResourceMetadata, decodedPlan interface{}, decodedState interface{}) {
	// client := metadata.Client.AppService.WebAppsClient
	// plan := assertFwLogicAppStandardResourceModel(decodedPlan, resp)
	// state := assertFwLogicAppStandardResourceModel(decodedState, resp)

}

func (r FwLogicAppStandardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, metadata sdk.ResourceMetadata, decodedState interface{}) {
	client := metadata.Client.AppService.WebAppsClient
	state := sdk.AssertResourceModelType[FwLogicAppStandardResourceModel](decodedState, resp)
	id, err := commonids.ParseLogicAppId(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "ID parsing error", err)
		return
	}

	if _, err = client.Delete(ctx, *id, webapps.DeleteOperationOptions{
		DeleteEmptyServerFarm: pointer.To(false),
		DeleteMetrics:         pointer.To(false),
	}); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("deleting %s:", *id), err.Error())
	}
}

func (r FwLogicAppStandardResource) Identity() (id resourceids.ResourceId, idType []sdk.ResourceTypeForIdentity) {
	return &commonids.LogicAppId{}, []sdk.ResourceTypeForIdentity{sdk.ResourceTypeForIdentityDefault}
}

func (r *FwLogicAppStandardResourceModel) buildBaseAppSettings(storageAccountDomainSuffix string) []webapps.NameValuePair {
	storageConnectionString := fmt.Sprintf(storageConnectionFmt, r.StorageAccountName.ValueString(), r.StorageAccountAccessKey.ValueString(), storageAccountDomainSuffix)
	contentSharePropVal := fmt.Sprintf("%s-content", strings.ToLower(r.Name.ValueString()))
	if v := r.StorageAccountShareName.ValueString(); v != "" {
		contentSharePropVal = v
	}

	result := []webapps.NameValuePair{
		{
			Name:  pointer.To(storagePropName),
			Value: pointer.To(storageConnectionString),
		},
		{
			Name:  pointer.To(functionVersionPropName),
			Value: pointer.To(r.Version.ValueString()),
		},
		{
			Name:  pointer.To(appKindPropName),
			Value: pointer.To(appKindPropValue),
		},
		{
			Name:  pointer.To(contentSharePropName),
			Value: pointer.To(contentSharePropVal),
		},
		{
			Name:  pointer.To(contentFileConnStringPropName),
			Value: pointer.To(storageConnectionString),
		},
	}

	if r.UseExtensionBundle.ValueBool() {
		result = append(result, webapps.NameValuePair{
			Name:  pointer.To(extensionBundlePropName),
			Value: pointer.To(extensionBundleName),
		})
		result = append(result, webapps.NameValuePair{
			Name:  pointer.To(extensionBundleVersionPropName),
			Value: pointer.To(r.BundleVersion.ValueString()),
		})
	}

	return result
}

func expandSiteConfigLogicAppStandard(ctx context.Context, siteConfig FwLogicAppStandardSiteConfigModel) (*webapps.SiteConfig, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := &webapps.SiteConfig{
		AlwaysOn:                               siteConfig.AlwaysOn.ValueBoolPointer(),
		FtpsState:                              pointer.ToEnum[webapps.FtpsState](siteConfig.FTPSState.ValueString()),
		HTTP20Enabled:                          siteConfig.HTTP2Enabled.ValueBoolPointer(),
		LinuxFxVersion:                         siteConfig.LinuxFXVersion.ValueStringPointer(),
		MinTlsVersion:                          pointer.ToEnum[webapps.SupportedTlsVersions](siteConfig.MinTLSVersion.ValueString()),
		PreWarmedInstanceCount:                 siteConfig.PreWarmedInstanceCount.ValueInt64Pointer(),
		ScmIPSecurityRestrictionsUseMain:       siteConfig.SCMUseMainIPRestriction.ValueBoolPointer(),
		ScmMinTlsVersion:                       pointer.ToEnum[webapps.SupportedTlsVersions](siteConfig.SCMMinTLSVersion.ValueString()),
		ScmType:                                pointer.ToEnum[webapps.ScmType](siteConfig.SCMType.ValueString()),
		Use32BitWorkerProcess:                  siteConfig.Use32BitWorkerProcesses.ValueBoolPointer(),
		WebSocketsEnabled:                      siteConfig.WebSocketsEnabled.ValueBoolPointer(),
		HealthCheckPath:                        siteConfig.HealthCheckPath.ValueStringPointer(),
		MinimumElasticInstanceCount:            siteConfig.ElasticInstanceMinimum.ValueInt64Pointer(),
		FunctionAppScaleLimit:                  siteConfig.AppScaleLimit.ValueInt64Pointer(),
		FunctionsRuntimeScaleMonitoringEnabled: siteConfig.RuntimeScaleMonitoringEnabled.ValueBoolPointer(),
		NetFrameworkVersion:                    siteConfig.DotnetFrameworkVersion.ValueStringPointer(),
		VnetRouteAllEnabled:                    siteConfig.VNETRouteAllEnabled.ValueBoolPointer(),
		AutoSwapSlotName:                       siteConfig.AutoSwapSlotName.ValueStringPointer(),
	}

	cors, diags := expandCorsSettings(ctx, siteConfig.CORS)
	if diags.HasError() {
		return nil, diags
	}

	result.Cors = cors

	ipr, diags := expandIPRestrictions(ctx, siteConfig.IPRestriction)
	result.IPSecurityRestrictions = ipr
	if diags.HasError() {
		return nil, diags
	}

	scmIpr, diags := expandIPRestrictions(ctx, siteConfig.SCMIPRestriction)
	if diags.HasError() {
		return nil, diags
	}
	result.ScmIPSecurityRestrictions = scmIpr

	return result, diags
}

func expandCorsSettings(ctx context.Context, input typehelpers.ListNestedObjectValueOf[FwLogicAppCORSSettingsModel]) (*webapps.CorsSettings, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := &webapps.CorsSettings{}
	convert.Expand(ctx, input, result, &diags)

	return result, diags
}

func expandIPRestrictions(ctx context.Context, input typehelpers.ListNestedObjectValueOf[FwLogicAppIPRestrictionModel]) (*[]webapps.IPSecurityRestriction, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := &[]webapps.IPSecurityRestriction{}
	convert.Expand(ctx, input, result, &diags)

	return result, diags
}

func flattenSiteConfigLogicAppStandard(ctx context.Context, input *webapps.SiteConfig) (typehelpers.ListNestedObjectValueOf[FwLogicAppStandardSiteConfigModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := typehelpers.NewListNestedObjectValueOfNull[FwLogicAppStandardSiteConfigModel](ctx)

	if input == nil {
		return result, nil
	}

	sc := FwLogicAppStandardSiteConfigModel{
		AlwaysOn:                      types.BoolPointerValue(input.AlwaysOn),
		FTPSState:                     types.StringValue(pointer.FromEnum(input.FtpsState)),
		HTTP2Enabled:                  types.BoolPointerValue(input.HTTP20Enabled),
		LinuxFXVersion:                types.StringPointerValue(input.LinuxFxVersion),
		MinTLSVersion:                 types.StringValue(pointer.FromEnum(input.MinTlsVersion)),
		PreWarmedInstanceCount:        types.Int64PointerValue(input.PreWarmedInstanceCount),
		SCMUseMainIPRestriction:       types.BoolPointerValue(input.ScmIPSecurityRestrictionsUseMain),
		SCMMinTLSVersion:              types.StringValue(pointer.FromEnum(input.ScmMinTlsVersion)),
		SCMType:                       types.StringValue(pointer.FromEnum(input.ScmType)),
		Use32BitWorkerProcesses:       types.BoolPointerValue(input.Use32BitWorkerProcess),
		WebSocketsEnabled:             types.BoolPointerValue(input.WebSocketsEnabled),
		HealthCheckPath:               types.StringPointerValue(input.HealthCheckPath),
		ElasticInstanceMinimum:        types.Int64PointerValue(input.MinimumElasticInstanceCount),
		AppScaleLimit:                 types.Int64PointerValue(input.FunctionAppScaleLimit),
		RuntimeScaleMonitoringEnabled: types.BoolPointerValue(input.FunctionsRuntimeScaleMonitoringEnabled),
		DotnetFrameworkVersion:        types.StringPointerValue(input.NetFrameworkVersion),
		VNETRouteAllEnabled:           types.BoolPointerValue(input.VnetRouteAllEnabled),
		AutoSwapSlotName:              types.StringPointerValue(input.AutoSwapSlotName),
	}

	// CORS
	cors, diags := flattenCORSSettings(ctx, input.Cors)
	if diags.HasError() {
		return result, diags
	}
	sc.CORS = cors
	// IPRestriction
	ipr, diags := flattenIPRestriction(ctx, input.IPSecurityRestrictions)
	if diags.HasError() {
		return result, diags
	}
	sc.IPRestriction = ipr
	// SCMIPRestriction
	scmIpr, diags := flattenIPRestriction(ctx, input.ScmIPSecurityRestrictions)
	if diags.HasError() {
		return result, diags
	}
	sc.SCMIPRestriction = scmIpr

	// result, diags = typehelpers.NewListNestedObjectValueOfSlice(ctx, []*FwLogicAppStandardSiteConfigModel{sc})
	result = typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []FwLogicAppStandardSiteConfigModel{sc})

	return result, diags
}

func flattenCORSSettings(ctx context.Context, settings *webapps.CorsSettings) (typehelpers.ListNestedObjectValueOf[FwLogicAppCORSSettingsModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := typehelpers.NewListNestedObjectValueOfNull[FwLogicAppCORSSettingsModel](ctx)

	if settings != nil {
		convert.Flatten(ctx, settings, &result, &diags)
	}

	return result, diags
}

func flattenIPRestriction(ctx context.Context, settings *[]webapps.IPSecurityRestriction) (typehelpers.ListNestedObjectValueOf[FwLogicAppIPRestrictionModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := typehelpers.NewListNestedObjectValueOfNull[FwLogicAppIPRestrictionModel](ctx)

	if settings != nil {
		if len(*settings) == 1 {
			// check if it's just the default value for the block
			def := webapps.IPSecurityRestriction{
				Action:      pointer.To("Allow"),
				Description: pointer.To("Allow all access"),
				IPAddress:   pointer.To("Any"),
				Name:        pointer.To("Allow all"),
				Priority:    pointer.To(int64(2147483647)),
			}
			if reflect.DeepEqual((*settings)[0], def) {
				return result, diags
			}
		}

		convert.Flatten(ctx, settings, &result, &diags)
	}

	return result, diags
}

func expandAppSettingsLogicAppStandard(ctx context.Context, input types.Map, baseSettings ...webapps.NameValuePair) *[]webapps.NameValuePair {
	result := make([]webapps.NameValuePair, 0)

	if !(input.IsUnknown() || input.IsNull()) {
		decode := make(map[string]string)
		input.ElementsAs(ctx, &decode, false)

		if len(decode) > 0 {
			for k, v := range decode {
				s := webapps.NameValuePair{
					Name:  pointer.To(k),
					Value: pointer.To(v),
				}
				result = append(result, s)
			}
		}
	}

	result = append(result, baseSettings...)

	return &result
}

func readLogicAppStandardAppSettings(ctx context.Context, id *commonids.AppServiceId, client *webapps.WebAppsClient, state *FwLogicAppStandardResourceModel, resp interface{}) {
	diags := &diag.Diagnostics{}
	appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("listing app settings: %s", id), err)
		return
	}

	if model := appSettingsResp.Model; model != nil {
		if props := model.Properties; props != nil && len(*props) > 0 {
			appSettings := *props

			connectionString := appSettings[storagePropName]
			connectionStringParts := strings.Split(connectionString, ";")
			for _, part := range connectionStringParts {
				switch {
				case strings.HasPrefix(part, "AccountName"):
					accountNameParts := strings.Split(part, "AccountName=")
					if len(accountNameParts) > 1 {
						state.StorageAccountName = types.StringValue(accountNameParts[1])
					}
				case strings.HasPrefix(part, "AccountKey"):
					accountKeyParts := strings.Split(part, "AccountKey=")
					if len(accountKeyParts) > 1 {
						state.StorageAccountAccessKey = types.StringValue(accountKeyParts[1])
					}
				}
			}

			state.Version = types.StringValue(appSettings[functionVersionPropName])

			_, useExtensionBundle := appSettings[extensionBundlePropName]
			state.UseExtensionBundle = types.BoolValue(useExtensionBundle)
			if useExtensionBundle {
				if v, ok := appSettings[extensionBundleVersionPropName]; ok {
					state.BundleVersion = types.StringValue(v)
				}
			}

			state.StorageAccountShareName = types.StringValue(appSettings[contentSharePropName])
			delete(appSettings, contentFileConnStringPropName)
			delete(appSettings, appKindPropName)
			delete(appSettings, extensionBundlePropName)
			delete(appSettings, extensionBundleVersionPropName)
			delete(appSettings, webJobsDashboardPropName)
			delete(appSettings, storagePropName)
			delete(appSettings, functionVersionPropName)
			delete(appSettings, contentSharePropName)

			convert.Flatten(ctx, appSettings, &state.AppSettings, diags)
			if diags.HasError() {
				sdk.AppendResponseErrorDiagnostic(resp, *diags)
				return
			}
		}
	}

	if state.AppSettings.IsUnknown() {
		state.AppSettings = typehelpers.NewMapValueOfNull[types.String](ctx)
	}
}

func ipRestrictionCommonSchema(ctx context.Context) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppIPRestrictionModel](ctx),
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"ip_address": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"service_tag": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"virtual_network_subnet_id": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"name": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},

				"priority": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  typehelpers.NewWrappedInt64Default(65000),
					Validators: []validator.Int64{
						int64validator.Between(1, math.MaxInt32),
					},
				},

				"action": schema.StringAttribute{
					Default:  typehelpers.NewWrappedStringDefault("Allow"),
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						typehelpers.WrappedStringValidator{
							Func: validation.StringInSlice([]string{
								"Allow",
								"Deny",
							}, false),
						},
					},
				},
			},
			Blocks: map[string]schema.Block{
				"headers": schema.ListNestedBlock{
					CustomType: typehelpers.NewListNestedObjectTypeOf[FwLogicAppIPRestrictionHeadersModel](ctx),
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"x_forwarded_host": schema.SetAttribute{
								CustomType:  typehelpers.NewSetTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.SizeAtMost(8),
								},
							},

							"x_forwarded_for": schema.SetAttribute{
								CustomType:  typehelpers.NewSetTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.All(
										setvalidator.SizeAtMost(8),
										setvalidator.ValueStringsAre(
											typehelpers.WrappedStringValidator{
												Func: validation.IsCIDR,
											},
										),
									),
								},
							},

							"x_azure_fdid": schema.SetAttribute{
								CustomType:  typehelpers.NewSetTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.All(
										setvalidator.SizeAtMost(8),
										setvalidator.ValueStringsAre(
											typehelpers.WrappedStringValidator{
												Func: validation.IsUUID,
											},
										),
									),
								},
							},

							"x_fd_health_probe": schema.ListAttribute{
								CustomType:  typehelpers.NewListTypeOf[types.String](ctx),
								ElementType: types.StringType,
								Optional:    true,
								Validators: []validator.List{
									listvalidator.All(
										listvalidator.SizeAtMost(1),
										listvalidator.ValueStringsAre(
											stringvalidator.OneOf("0", "1"),
										),
									),
								},
							},
						},
					},
				},
			},
		},
	}
}
