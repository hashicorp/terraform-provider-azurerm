package appservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/identity"
	"github.com/hashicorp/go-azure-helpers/framework/location"
	"github.com/hashicorp/go-azure-helpers/framework/resourceid"
	"github.com/hashicorp/go-azure-helpers/framework/tags"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/sdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
)

type FWStaticWebAppResource struct {
	sdk.ResourceMetadata
}

var _ sdk.FrameworkResource = &FWStaticWebAppResource{}

var _ sdk.FrameworkResourceWithValidateConfig = &FWStaticWebAppResource{}

type FWStaticWebAppResourceModel struct {
	Name                            types.String `tfsdk:"name"`
	ResourceGroupName               types.String `tfsdk:"resource_group_name"`
	Location                        types.String `tfsdk:"location"`
	ConfigurationFileChangesEnabled types.Bool   `tfsdk:"configuration_file_changes_enabled"`
	PreviewEnvironmentsEnabled      types.Bool   `tfsdk:"preview_environments_enabled"`
	SkuTier                         types.String `tfsdk:"sku_tier"`
	SkuSize                         types.String `tfsdk:"sku_size"`
	AppSettings                     types.Map    `tfsdk:"app_settings"`
	Tags                            types.Map    `tfsdk:"tags"`

	BasicAuth types.List `tfsdk:"basic_auth"`
	Identity  types.List `tfsdk:"identity"`

	ApiKey          types.String `tfsdk:"api_key"`
	DefaultHostName types.String `tfsdk:"default_host_name"`

	Id types.String `tfsdk:"id"`
}

func (r *FWStaticWebAppResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, configResponse *resource.ValidateConfigResponse) {
	data := FWStaticWebAppResourceModel{}
	configResponse.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if configResponse.Diagnostics.HasError() {
		return
	}

	if strings.EqualFold(data.SkuTier.ValueString(), string(resourceproviders.SkuNameFree)) && strings.EqualFold(data.SkuSize.ValueString(), string(resourceproviders.SkuNameFree)) {
		if len(data.BasicAuth.Elements()) > 0 {
			configResponse.Diagnostics.AddError("config validation error", "basic_auth cannot be used with the Free tier of Static Web Apps")
			return
		}
		if len(data.Identity.Elements()) > 0 {
			configResponse.Diagnostics.AddError("config validation error", "identities cannot be used with the Free tier of Static Web Apps")
			return
		}
	}
}

type FWStaticWebAppAuthModel struct {
	Password     types.String `tfsdk:"password"`
	Environments types.String `tfsdk:"environments"`
}

var FWStaticWebAppAuthModelAttributes = map[string]attr.Type{
	"password":     types.StringType,
	"environments": types.StringType,
}

func NewFWStaticWebAppResource() resource.Resource {
	return &FWStaticWebAppResource{}
}

func (r *FWStaticWebAppResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "azurerm_fw_static_web_app"
}

func (r *FWStaticWebAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},

			"resource_group_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},

			"location": location.LocationAttribute(),

			"configuration_file_changes_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "",
				MarkdownDescription: "",
				Default:             typehelpers.NewWrappedBoolDefault(false),
			},

			"preview_environments_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "",
				MarkdownDescription: "",
				Default:             typehelpers.NewWrappedBoolDefault(false),
			},

			"sku_tier": schema.StringAttribute{
				Optional:            true,
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(resourceproviders.SkuNameStandard),
						string(resourceproviders.SkuNameFree),
					),
				},
				Default:  typehelpers.NewWrappedStringDefault(resourceproviders.SkuNameFree),
				Computed: true,
			},

			"sku_size": schema.StringAttribute{
				Optional:            true,
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(resourceproviders.SkuNameStandard),
						string(resourceproviders.SkuNameFree),
					),
				},
				Default:  typehelpers.NewWrappedStringDefault(resourceproviders.SkuNameFree),
				Computed: true,
			},

			"app_settings": schema.MapAttribute{
				ElementType:         basetypes.StringType{},
				Optional:            true,
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.KeysAre(
						stringvalidator.LengthAtLeast(1),
					),
				},
			},

			"tags": schema.MapAttribute{
				ElementType:         basetypes.StringType{},
				Optional:            true,
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
				},
			},

			"api_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"default_host_name": schema.StringAttribute{
				Computed: true,
			},

			// id not implicit
			"id": resourceid.IDAttribute(),

			// Blocks as attributes - Protocol v6 required :(
			//"basic_auth": schema.ListNestedAttribute{
			//	NestedObject: schema.NestedAttributeObject{
			//		Attributes: map[string]schema.Attribute{
			//			"password": schema.StringAttribute{
			//				Required:            true,
			//				Sensitive:           true,
			//				Description:         "",
			//				MarkdownDescription: "",
			//				Validators: []validator.String{
			//					typehelpers.WrappedStringValidator{
			//						Func: validate.StaticWebAppPassword,
			//					},
			//				},
			//			},
			//			"environments": schema.StringAttribute{
			//				Required:            true,
			//				Description:         "",
			//				MarkdownDescription: "",
			//				Validators: []validator.String{
			//					stringvalidator.OneOf(
			//						helpers.EnvironmentsTypeAllEnvironments,
			//						helpers.EnvironmentsTypeStagingEnvironments,
			//						helpers.EnvironmentsTypeSpecifiedEnvironments,
			//					),
			//				},
			//			},
			//		},
			//	},
			//	Validators: []validator.List{
			//		listvalidator.SizeAtMost(1),
			//	},
			//},
			//
			//"identity": typehelpers.IdentitySchemaAttribute(),
		},

		Blocks: map[string]schema.Block{
			"basic_auth": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"password": schema.StringAttribute{
							Required:            true,
							Sensitive:           true,
							Description:         "",
							MarkdownDescription: "",
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: validate.StaticWebAppPassword,
								},
							},
						},

						"environments": schema.StringAttribute{
							Required:            true,
							Description:         "",
							MarkdownDescription: "",
							Validators: []validator.String{
								stringvalidator.OneOf(
									helpers.EnvironmentsTypeAllEnvironments,
									helpers.EnvironmentsTypeStagingEnvironments,
									helpers.EnvironmentsTypeSpecifiedEnvironments,
								),
							},
						},
					},
				},
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"identity": identity.IdentitySchemaBlock(),
		},

		Description:         "",
		MarkdownDescription: "",
		Version:             0,
	}
}

func (r *FWStaticWebAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *FWStaticWebAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.TimeoutCreate)
	defer cancel()

	client := r.Client.AppService.StaticSitesClient

	data := FWStaticWebAppResourceModel{}

	if ok := r.DecodeCreate(ctx, req, resp, &data); !ok {
		return
	}

	id := staticsites.NewStaticSiteID(r.SubscriptionId, data.ResourceGroupName.ValueString(), data.Name.ValueString())

	existing, err := client.GetStaticSite(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking for presence of existing %s", id), err.Error())
			return
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		sdk.SetResponseErrorDiagnostic(resp, "Import Error", tf.ImportAsExistsError("azurerm_fw_static_web_app", id.ID()).Error())
		return
	}

	envelope := staticsites.StaticSiteARMResource{
		Location: location.Normalize(data.Location.ValueString()),
		Sku: &staticsites.SkuDescription{
			Name: data.SkuSize.ValueStringPointer(),
			Tier: data.SkuTier.ValueStringPointer(),
		},
	}

	t, diags := tags.Expand(data.Tags)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	envelope.Tags = t

	// Identity
	ident, diags := identity.ExpandSystemAndUserAssignedMap(data.Identity)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	envelope.Identity = ident

	props := &staticsites.StaticSite{
		AllowConfigFileUpdates:   data.ConfigurationFileChangesEnabled.ValueBoolPointer(),
		StagingEnvironmentPolicy: pointer.To(staticsites.StagingEnvironmentPolicyEnabled),
	}

	if !data.PreviewEnvironmentsEnabled.ValueBool() {
		props.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyDisabled)
	}

	envelope.Properties = props

	if err = client.CreateOrUpdateStaticSiteThenPoll(ctx, id, envelope); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("creating %s", id), fmt.Sprintf("%+v", err))
		//resp.Diagnostics.AddError(fmt.Sprintf("creating %s", id), fmt.Sprintf("%+v", err))
		return
	}

	data.Id = types.StringValue(id.ID())

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.Id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.AppSettings.IsNull() && !data.AppSettings.IsUnknown() {
		expanded, d := typehelpers.ExpandMapPointer[string](data.AppSettings)
		if d.HasError() {
			resp.Diagnostics.Append(d...)
			return
		}
		appSettings := staticsites.StringDictionary{
			Properties: expanded,
		}

		if _, err = client.CreateOrUpdateStaticSiteAppSettings(ctx, id, appSettings); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("updating app settings for %s", id), fmt.Sprintf("%+v", err))
			return
		}
	}

	if !data.BasicAuth.IsNull() && !data.BasicAuth.IsUnknown() {
		sdkHackClient := sdkhacks.NewStaticWebAppClient(client)

		auths := make([]FWStaticWebAppAuthModel, 0)
		diags = typehelpers.ExpandList(data.BasicAuth, &auths)
		if diags.HasError() {
			return
		}

		authProps := staticsites.StaticSiteBasicAuthPropertiesARMResource{
			Properties: &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
				ApplicableEnvironmentsMode: auths[0].Environments.ValueString(),
				Password:                   auths[0].Password.ValueStringPointer(),
				SecretState:                pointer.To("Password"),
			},
		}

		if _, err = sdkHackClient.CreateOrUpdateBasicAuth(ctx, id, authProps); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("setting basic auth on %s", id), fmt.Sprintf("%+v", err))
		}
	}

	// Computed values are expected to be known at the end of the Create, and an additional read is not performed!
	read, err := client.GetStaticSite(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", id), fmt.Sprintf("%+v", err))
	}
	if model := read.Model; model != nil {
		if readProps := model.Properties; readProps != nil {
			data.ConfigurationFileChangesEnabled = types.BoolPointerValue(readProps.AllowConfigFileUpdates)
			data.DefaultHostName = types.StringPointerValue(readProps.DefaultHostname)
		}
	}

	sec, err := client.ListStaticSiteSecrets(ctx, id)
	if err != nil || sec.Model == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving secrets for %s", id), fmt.Sprintf("%+v", err))
		return
	}

	if secProps := sec.Model.Properties; secProps != nil {
		propsMap := pointer.From(secProps)
		apiKey := ""
		apiKey = propsMap["apiKey"]
		data.ApiKey = types.StringValue(apiKey)
	}

	r.EncodeCreate(ctx, resp, &data)
}

func (r *FWStaticWebAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.TimeoutRead)
	defer cancel()

	diags := diag.Diagnostics{}

	client := r.Client.AppService.StaticSitesClient

	state := FWStaticWebAppResourceModel{}

	if ok := r.DecodeRead(ctx, req, resp, &state); !ok {
		return
	}

	id, err := staticsites.ParseStaticSiteID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("parsing ID", err.Error())
		return
	}
	state.Name = types.StringValue(id.StaticSiteName)
	state.ResourceGroupName = types.StringValue(id.ResourceGroupName)

	staticSite, err := client.GetStaticSite(ctx, *id)
	if err != nil {
		if response.WasNotFound(staticSite.HttpResponse) {
			resp.Diagnostics.AddError(fmt.Sprintf("reading %s, removing from state", id.String()), err.Error())
			state.Id = types.StringValue("")
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", *id), fmt.Sprintf("%+v", err))
		return
	}

	if model := staticSite.Model; model != nil {
		state.Location = types.StringValue(location.Normalize(model.Location))

		t, diags := tags.Flatten(model.Tags)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.Tags = t

		ident, diags := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		state.Identity = ident

		if sku := model.Sku; sku != nil {
			state.SkuSize = types.StringPointerValue(sku.Name)
			state.SkuTier = types.StringPointerValue(sku.Tier)
		}

		if props := model.Properties; props != nil {
			state.ConfigurationFileChangesEnabled = types.BoolPointerValue(props.AllowConfigFileUpdates)
			state.DefaultHostName = types.StringPointerValue(props.DefaultHostname)
			state.PreviewEnvironmentsEnabled = types.BoolValue(pointer.From(props.StagingEnvironmentPolicy) == staticsites.StagingEnvironmentPolicyEnabled)
		}
	}

	sec, err := client.ListStaticSiteSecrets(ctx, *id)
	if err != nil || sec.Model == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving secrets for %s", *id), fmt.Sprintf("%+v", err))
		return
	}

	if secProps := sec.Model.Properties; secProps != nil {
		propsMap := pointer.From(secProps)
		apiKey := ""
		apiKey = propsMap["apiKey"]
		state.ApiKey = types.StringValue(apiKey)
	}

	appSettings, err := client.ListStaticSiteAppSettings(ctx, *id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving app_settings for %s", *id), err.Error())
		return
	}

	if appSettingsModel := appSettings.Model; appSettingsModel != nil {
		state.AppSettings, diags = typehelpers.FlattenMapPointer(appSettingsModel.Properties)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	sdkHackClient := sdkhacks.NewStaticWebAppClient(client)
	auth, err := sdkHackClient.GetBasicAuth(ctx, *id)
	if err != nil && !response.WasNotFound(auth.HttpResponse) { // If basic auth is not configured then this 404's
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving auth config for %s", *id), fmt.Sprintf("%+v", err))
	}
	if !response.WasNotFound(auth.HttpResponse) {
		if authModel := auth.Model; authModel != nil && authModel.Properties != nil && !strings.EqualFold(authModel.Properties.ApplicableEnvironmentsMode, helpers.EnvironmentsTypeSpecifiedEnvironments) {
			state.BasicAuth, diags = types.ListValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(FWStaticWebAppAuthModelAttributes), *authModel.Properties)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
		}
	}

	r.EncodeRead(ctx, resp, &state)
}

func (r *FWStaticWebAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.TimeoutCreate)
	defer cancel()

	client := r.Client.AppService.StaticSitesClient

	plan := FWStaticWebAppResourceModel{}

	state := FWStaticWebAppResourceModel{}

	r.DecodeUpdate(ctx, req, resp, &plan, &state)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := staticsites.ParseStaticSiteID(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("parsing ID", err.Error())
		return
	}

	existing, err := client.GetStaticSite(ctx, *id)
	if err != nil || existing.Model == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", *id), fmt.Sprintf("%+v", err))
		return
	}

	model := *existing.Model

	if !plan.Identity.Equal(state.Identity) {
		ident, diags := identity.ExpandSystemAndUserAssignedMap(plan.Identity)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		model.Identity = ident
		if !state.SkuTier.Equal(plan.SkuTier) && strings.EqualFold(string(resourceproviders.SkuNameFree), plan.SkuTier.ValueString()) {
			if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, *id, model); err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("updating %s", *id), fmt.Sprintf("%+v", err))
				return
			}
		}
	}

	if !plan.SkuTier.Equal(state.SkuTier) || !plan.SkuSize.Equal(state.SkuSize) {
		model.Sku = &staticsites.SkuDescription{
			Size: plan.SkuSize.ValueStringPointer(),
			Tier: plan.SkuTier.ValueStringPointer(),
		}
	}

	if !plan.ConfigurationFileChangesEnabled.Equal(state.ConfigurationFileChangesEnabled) {
		if plan.ConfigurationFileChangesEnabled.ValueBool() {
			model.Properties.AllowConfigFileUpdates = plan.ConfigurationFileChangesEnabled.ValueBoolPointer()
		}
	}

	if !plan.PreviewEnvironmentsEnabled.Equal(state.PreviewEnvironmentsEnabled) {
		if plan.PreviewEnvironmentsEnabled.ValueBool() {
			model.Properties.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyEnabled)
		} else {
			model.Properties.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyDisabled)
		}
	}

	if !plan.Tags.Equal(state.Tags) {
		t, diags := tags.Expand(plan.Tags)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		model.Tags = t
	}

	if err = client.CreateOrUpdateStaticSiteThenPoll(ctx, *id, model); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("updating %s", id), fmt.Sprintf("%+v", err))
		return
	}
}

func (r *FWStaticWebAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.TimeoutDelete)
	defer cancel()

	client := r.Client.AppService.StaticSitesClient

	state := FWStaticWebAppResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := staticsites.ParseStaticSiteID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("parsing ID", err.Error())
		return
	}

	if err = client.DeleteStaticSiteThenPoll(ctx, *id); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("deleting %s", *id), fmt.Sprintf("%+v", err))
		return
	}
}

func (r *FWStaticWebAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type MyEphemeralResource struct {
	sdk.EphemeralResourceMetadata
}

var _ sdk.EphemeralResource = &MyEphemeralResource{}

type MyEphemeralResourceModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *MyEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "azurerm_my_ephemeral_resource"
}

func (m *MyEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ...
		},
		Blocks: map[string]schema.Block{
			// ...
		},
	}
}

func (m *MyEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	m.Defaults(req, resp)
}

func (m *MyEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	client := m.Client.SomeAzureService.FooClient

	var data MyEphemeralResourceModel

	if !m.DecodeOpen(ctx, req, resp, &data) {
		return
	}

	// do the thing...
}
