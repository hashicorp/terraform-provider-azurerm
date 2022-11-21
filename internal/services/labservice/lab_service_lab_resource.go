package labservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LabServiceLabModel struct {
	Name                  string                  `tfschema:"name"`
	ResourceGroupName     string                  `tfschema:"resource_group_name"`
	Location              string                  `tfschema:"location"`
	AutoShutdownProfile   []AutoShutdownProfile   `tfschema:"auto_shutdown_profile"`
	ConnectionProfile     []ConnectionProfile     `tfschema:"connection_profile"`
	SecurityProfile       []SecurityProfile       `tfschema:"security_profile"`
	Title                 string                  `tfschema:"title"`
	VirtualMachineProfile []VirtualMachineProfile `tfschema:"virtual_machine_profile"`
	Description           string                  `tfschema:"description"`
	LabPlanId             string                  `tfschema:"lab_plan_id"`
	Tags                  map[string]string       `tfschema:"tags"`
}

type AutoShutdownProfile struct {
	DisconnectDelay                 string                 `tfschema:"disconnect_delay"`
	IdleDelay                       string                 `tfschema:"idle_delay"`
	NoConnectDelay                  string                 `tfschema:"no_connect_delay"`
	ShutdownOnDisconnectEnabled     bool                   `tfschema:"shutdown_on_disconnect_enabled"`
	ShutdownOnIdle                  lab.ShutdownOnIdleMode `tfschema:"shutdown_on_idle"`
	ShutdownEnabledWhenNotConnected bool                   `tfschema:"shutdown_enabled_when_not_connected"`
}

type ConnectionProfile struct {
	ClientRdpAccess lab.ConnectionType `tfschema:"client_rdp_access"`
	ClientSshAccess lab.ConnectionType `tfschema:"client_ssh_access"`
	WebRdpAccess    lab.ConnectionType `tfschema:"web_rdp_access"`
	WebSshAccess    lab.ConnectionType `tfschema:"web_ssh_access"`
}

type SecurityProfile struct {
	OpenAccessEnabled bool `tfschema:"open_access_enabled"`
}

type VirtualMachineProfile struct {
	AdditionalCapability  []AdditionalCapability `tfschema:"additional_capability"`
	AdminUser             []Credential           `tfschema:"admin_user"`
	CreateOption          lab.CreateOption       `tfschema:"create_option"`
	ImageReference        []ImageReference       `tfschema:"image_reference"`
	NonAdminUser          []Credential           `tfschema:"non_admin_user"`
	Sku                   []Sku                  `tfschema:"sku"`
	UsageQuota            string                 `tfschema:"usage_quota"`
	SharedPasswordEnabled bool                   `tfschema:"shared_password_enabled"`
}

type AdditionalCapability struct {
	GpuDriversInstalled bool `tfschema:"gpu_drivers_installed"`
}

type Credential struct {
	Password string `tfschema:"password"`
	Username string `tfschema:"username"`
}

type ImageReference struct {
	Id        string `tfschema:"id"`
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type Sku struct {
	Capacity int64  `tfschema:"capacity"`
	Name     string `tfschema:"name"`
}

type LabServiceLabResource struct{}

var _ sdk.ResourceWithUpdate = LabServiceLabResource{}

func (r LabServiceLabResource) ResourceType() string {
	return "azurerm_lab_service_lab"
}

func (r LabServiceLabResource) ModelObject() interface{} {
	return &LabServiceLabModel{}
}

func (r LabServiceLabResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return lab.ValidateLabID
}

func (r LabServiceLabResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LabName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"auto_shutdown_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disconnect_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"idle_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"no_connect_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"shutdown_on_disconnect_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"shutdown_on_idle": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ShutdownOnIdleModeUserAbsence),
							string(lab.ShutdownOnIdleModeLowUsage),
							string(lab.ShutdownOnIdleModeNone),
						}, false),
					},

					"shutdown_enabled_when_not_connected": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"connection_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_rdp_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypeNone),
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypePublic),
						}, false),
					},

					"client_ssh_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypeNone),
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypePublic),
						}, false),
					},

					"web_rdp_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypeNone),
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypePublic),
						}, false),
					},

					"web_ssh_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypeNone),
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypePublic),
						}, false),
					},
				},
			},
		},

		"security_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"open_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"title": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LabTitle,
		},

		"virtual_machine_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_user": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"username": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.LabUsername,
								},

								"password": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.LabPassword,
								},
							},
						},
					},

					"create_option": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.CreateOptionImage),
							string(lab.CreateOptionTemplateVM),
						}, false),
					},

					"image_reference": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: computeValidate.SharedImageID,
									ConflictsWith: []string{
										"virtual_machine_profile.0.image_reference.0.offer",
										"virtual_machine_profile.0.image_reference.0.publisher",
										"virtual_machine_profile.0.image_reference.0.sku",
										"virtual_machine_profile.0.image_reference.0.version",
									},
								},

								"offer": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"virtual_machine_profile.0.image_reference.0.id"},
								},

								"publisher": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"virtual_machine_profile.0.image_reference.0.id"},
								},

								"sku": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"virtual_machine_profile.0.image_reference.0.id"},
								},

								"version": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validate.LabImageVersion,
									ConflictsWith: []string{"virtual_machine_profile.0.image_reference.0.id"},
								},
							},
						},
					},

					"sku": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.LabSkuName,
								},

								"capacity": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 400),
								},
							},
						},
					},

					"usage_quota": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"shared_password_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						ForceNew: true,
					},

					"additional_capability": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"gpu_drivers_installed": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
									Default:  false,
								},
							},
						},
					},

					"non_admin_user": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"username": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.LabUsername,
								},

								"password": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.LabPassword,
								},
							},
						},
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.LabDescription,
		},

		"lab_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: labplan.ValidateLabPlanID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r LabServiceLabResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServiceLabResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServiceLabModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LabService.LabClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := lab.NewLabID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := &lab.Lab{
				Location: location.Normalize(model.Location),
				Properties: lab.LabProperties{
					Title: &model.Title,
				},
				Tags: &model.Tags,
			}

			autoShutdownProfile, err := expandAutoShutdownProfile(model.AutoShutdownProfile)
			if err != nil {
				return err
			}
			props.Properties.AutoShutdownProfile = *autoShutdownProfile

			connectionProfile, err := expandConnectionProfile(model.ConnectionProfile)
			if err != nil {
				return err
			}
			props.Properties.ConnectionProfile = *connectionProfile

			securityProfile, err := expandSecurityProfile(model.SecurityProfile)
			if err != nil {
				return err
			}
			props.Properties.SecurityProfile = *securityProfile

			virtualMachineProfile, err := expandVirtualMachineProfile(model.VirtualMachineProfile, true)
			if err != nil {
				return err
			}
			props.Properties.VirtualMachineProfile = *virtualMachineProfile

			if model.Description != "" {
				props.Properties.Description = &model.Description
			}

			if model.LabPlanId != "" {
				props.Properties.LabPlanId = &model.LabPlanId
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LabServiceLabResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServiceLabModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("auto_shutdown_profile") {
				autoShutdownProfile, err := expandAutoShutdownProfile(model.AutoShutdownProfile)
				if err != nil {
					return err
				}
				properties.Properties.AutoShutdownProfile = *autoShutdownProfile
			}

			if metadata.ResourceData.HasChange("connection_profile") {
				connectionProfile, err := expandConnectionProfile(model.ConnectionProfile)
				if err != nil {
					return err
				}
				properties.Properties.ConnectionProfile = *connectionProfile
			}

			if metadata.ResourceData.HasChange("security_profile") {
				securityProfile, err := expandSecurityProfile(model.SecurityProfile)
				if err != nil {
					return err
				}
				properties.Properties.SecurityProfile = *securityProfile
			}

			if metadata.ResourceData.HasChange("title") {
				if model.Title != "" {
					properties.Properties.Title = &model.Title
				} else {
					properties.Properties.Title = nil
				}
			}

			if metadata.ResourceData.HasChange("virtual_machine_profile") {
				virtualMachineProfile, err := expandVirtualMachineProfile(model.VirtualMachineProfile, false)
				if err != nil {
					return err
				}
				properties.Properties.VirtualMachineProfile = *virtualMachineProfile
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = &model.Description
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("lab_plan_id") {
				if model.LabPlanId != "" {
					properties.Properties.LabPlanId = &model.LabPlanId
				} else {
					properties.Properties.LabPlanId = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LabServiceLabResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := LabServiceLabModel{
				Name:              id.LabName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties

			autoShutdownProfile, err := flattenAutoShutdownProfile(&properties.AutoShutdownProfile)
			if err != nil {
				return err
			}
			state.AutoShutdownProfile = autoShutdownProfile

			connectionProfile, err := flattenConnectionProfile(&properties.ConnectionProfile)
			if err != nil {
				return err
			}
			state.ConnectionProfile = connectionProfile

			securityProfile, err := flattenSecurityProfile(&properties.SecurityProfile)
			if err != nil {
				return err
			}
			state.SecurityProfile = securityProfile

			if properties.Title != nil {
				state.Title = *properties.Title
			}

			virtualMachineProfile, err := flattenVirtualMachineProfile(&properties.VirtualMachineProfile, metadata.ResourceData)
			if err != nil {
				return err
			}
			state.VirtualMachineProfile = virtualMachineProfile

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.LabPlanId != nil {
				state.LabPlanId = *properties.LabPlanId
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServiceLabResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LabServiceLabResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if oldVal, newVal := rd.GetChange("virtual_machine_profile.0.non_admin_user"); len(oldVal.(*pluginsdk.Set).List()) == 0 && len(newVal.(*pluginsdk.Set).List()) == 1 {
				if err := rd.ForceNew("virtual_machine_profile.0.non_admin_user"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func expandAutoShutdownProfile(input []AutoShutdownProfile) (*lab.AutoShutdownProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	autoShutdownProfile := &input[0]
	result := lab.AutoShutdownProfile{}

	if autoShutdownProfile.DisconnectDelay != "" {
		result.DisconnectDelay = &autoShutdownProfile.DisconnectDelay
	}

	if autoShutdownProfile.IdleDelay != "" {
		result.IdleDelay = &autoShutdownProfile.IdleDelay
	}

	if autoShutdownProfile.NoConnectDelay != "" {
		result.NoConnectDelay = &autoShutdownProfile.NoConnectDelay
	}

	shutdownOnDisconnectEnabled := lab.EnableStateEnabled
	if !autoShutdownProfile.ShutdownOnDisconnectEnabled {
		shutdownOnDisconnectEnabled = lab.EnableStateDisabled
	}
	result.ShutdownOnDisconnect = &shutdownOnDisconnectEnabled

	if autoShutdownProfile.ShutdownOnIdle != "" {
		result.ShutdownOnIdle = &autoShutdownProfile.ShutdownOnIdle
	}

	shutdownEnabledWhenNotConnected := lab.EnableStateEnabled
	if !autoShutdownProfile.ShutdownEnabledWhenNotConnected {
		shutdownEnabledWhenNotConnected = lab.EnableStateDisabled
	}
	result.ShutdownWhenNotConnected = &shutdownEnabledWhenNotConnected

	return &result, nil
}

func flattenAutoShutdownProfile(input *lab.AutoShutdownProfile) ([]AutoShutdownProfile, error) {
	var autoShutdownProfiles []AutoShutdownProfile
	if input == nil {
		return autoShutdownProfiles, nil
	}

	autoShutdownProfile := AutoShutdownProfile{}

	if input.DisconnectDelay != nil {
		autoShutdownProfile.DisconnectDelay = *input.DisconnectDelay
	}

	if input.IdleDelay != nil {
		autoShutdownProfile.IdleDelay = *input.IdleDelay
	}

	if input.NoConnectDelay != nil {
		autoShutdownProfile.NoConnectDelay = *input.NoConnectDelay
	}

	if input.ShutdownOnDisconnect != nil {
		autoShutdownProfile.ShutdownOnDisconnectEnabled = *input.ShutdownOnDisconnect == lab.EnableStateEnabled
	}

	if input.ShutdownOnIdle != nil {
		autoShutdownProfile.ShutdownOnIdle = *input.ShutdownOnIdle
	}

	if input.ShutdownWhenNotConnected != nil {
		autoShutdownProfile.ShutdownEnabledWhenNotConnected = *input.ShutdownWhenNotConnected == lab.EnableStateEnabled
	}

	return append(autoShutdownProfiles, autoShutdownProfile), nil
}

func expandConnectionProfile(input []ConnectionProfile) (*lab.ConnectionProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	connectionProfile := &input[0]
	result := lab.ConnectionProfile{}

	if connectionProfile.ClientRdpAccess != "" {
		result.ClientRdpAccess = &connectionProfile.ClientRdpAccess
	}

	if connectionProfile.ClientSshAccess != "" {
		result.ClientSshAccess = &connectionProfile.ClientSshAccess
	}

	if connectionProfile.WebRdpAccess != "" {
		result.WebRdpAccess = &connectionProfile.WebRdpAccess
	}

	if connectionProfile.WebSshAccess != "" {
		result.WebSshAccess = &connectionProfile.WebSshAccess
	}

	return &result, nil
}

func flattenConnectionProfile(input *lab.ConnectionProfile) ([]ConnectionProfile, error) {
	var connectionProfiles []ConnectionProfile
	if input == nil {
		return connectionProfiles, nil
	}

	connectionProfile := ConnectionProfile{}

	if input.ClientRdpAccess != nil {
		connectionProfile.ClientRdpAccess = *input.ClientRdpAccess
	}

	if input.ClientSshAccess != nil {
		connectionProfile.ClientSshAccess = *input.ClientSshAccess
	}

	if input.WebRdpAccess != nil {
		connectionProfile.WebRdpAccess = *input.WebRdpAccess
	}

	if input.WebSshAccess != nil {
		connectionProfile.WebSshAccess = *input.WebSshAccess
	}

	return append(connectionProfiles, connectionProfile), nil
}

func expandSecurityProfile(input []SecurityProfile) (*lab.SecurityProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	securityProfile := &input[0]
	result := lab.SecurityProfile{}

	openAccessEnabled := lab.EnableStateEnabled
	if !securityProfile.OpenAccessEnabled {
		openAccessEnabled = lab.EnableStateDisabled
	}
	result.OpenAccess = &openAccessEnabled

	return &result, nil
}

func flattenSecurityProfile(input *lab.SecurityProfile) ([]SecurityProfile, error) {
	var securityProfiles []SecurityProfile
	if input == nil {
		return securityProfiles, nil
	}

	securityProfile := SecurityProfile{}

	if input.OpenAccess != nil {
		securityProfile.OpenAccessEnabled = *input.OpenAccess == lab.EnableStateEnabled
	}

	return append(securityProfiles, securityProfile), nil
}

func expandVirtualMachineProfile(input []VirtualMachineProfile, includePassword bool) (*lab.VirtualMachineProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	virtualMachineProfile := &input[0]
	result := lab.VirtualMachineProfile{
		CreateOption: virtualMachineProfile.CreateOption,
		UsageQuota:   virtualMachineProfile.UsageQuota,
	}

	sharedPasswordEnabled := lab.EnableStateEnabled
	if !virtualMachineProfile.SharedPasswordEnabled {
		sharedPasswordEnabled = lab.EnableStateDisabled
	}
	result.UseSharedPassword = &sharedPasswordEnabled

	additionalCapability, err := expandAdditionalCapability(virtualMachineProfile.AdditionalCapability)
	if err != nil {
		return nil, err
	}
	result.AdditionalCapabilities = additionalCapability

	adminUserValue, err := expandCredential(virtualMachineProfile.AdminUser, includePassword)
	if err != nil {
		return nil, err
	}
	result.AdminUser = *adminUserValue

	imageReferenceValue, err := expandImageReference(virtualMachineProfile.ImageReference)
	if err != nil {
		return nil, err
	}
	result.ImageReference = *imageReferenceValue

	nonAdminUserValue, err := expandCredential(virtualMachineProfile.NonAdminUser, includePassword)
	if err != nil {
		return nil, err
	}
	result.NonAdminUser = nonAdminUserValue

	sku, err := expandSku(virtualMachineProfile.Sku)
	if err != nil {
		return nil, err
	}
	result.Sku = *sku

	return &result, nil
}

func expandAdditionalCapability(input []AdditionalCapability) (*lab.VirtualMachineAdditionalCapabilities, error) {
	if len(input) == 0 {
		return nil, nil
	}

	additionalCapability := &input[0]
	result := lab.VirtualMachineAdditionalCapabilities{}

	gpuDriversInstalled := lab.EnableStateEnabled
	if !additionalCapability.GpuDriversInstalled {
		gpuDriversInstalled = lab.EnableStateDisabled
	}
	result.InstallGpuDrivers = &gpuDriversInstalled

	return &result, nil
}

func expandCredential(input []Credential, includePassword bool) (*lab.Credentials, error) {
	if len(input) == 0 {
		return nil, nil
	}

	credential := &input[0]
	result := lab.Credentials{
		Username: credential.Username,
	}

	if includePassword && credential.Password != "" {
		result.Password = &credential.Password
	}

	return &result, nil
}

func expandImageReference(input []ImageReference) (*lab.ImageReference, error) {
	if len(input) == 0 {
		return nil, nil
	}

	imageReference := &input[0]
	result := lab.ImageReference{}

	if imageReference.Id != "" {
		result.Id = &imageReference.Id
	}

	if imageReference.Offer != "" {
		result.Offer = &imageReference.Offer
	}

	if imageReference.Publisher != "" {
		result.Publisher = &imageReference.Publisher
	}

	if imageReference.Sku != "" {
		result.Sku = &imageReference.Sku
	}

	if imageReference.Version != "" {
		result.Version = &imageReference.Version
	}

	return &result, nil
}

func expandSku(input []Sku) (*lab.Sku, error) {
	if len(input) == 0 {
		return nil, nil
	}

	sku := &input[0]
	result := lab.Sku{
		Name:     sku.Name,
		Capacity: &sku.Capacity,
	}

	return &result, nil
}

func flattenVirtualMachineProfile(input *lab.VirtualMachineProfile, d *schema.ResourceData) ([]VirtualMachineProfile, error) {
	var virtualMachineProfiles []VirtualMachineProfile
	if input == nil {
		return virtualMachineProfiles, nil
	}

	virtualMachineProfile := VirtualMachineProfile{
		CreateOption: input.CreateOption,
		UsageQuota:   input.UsageQuota,
	}

	additionalCapability, err := flattenAdditionalCapability(input.AdditionalCapabilities)
	if err != nil {
		return nil, err
	}
	virtualMachineProfile.AdditionalCapability = additionalCapability

	adminUser, err := flattenCredential(&input.AdminUser, d.Get("virtual_machine_profile.0.admin_user.0.password").(string))
	if err != nil {
		return nil, err
	}
	virtualMachineProfile.AdminUser = adminUser

	imageReference, err := flattenImageReference(&input.ImageReference)
	if err != nil {
		return nil, err
	}
	virtualMachineProfile.ImageReference = imageReference

	nonAdminUser, err := flattenCredential(input.NonAdminUser, d.Get("virtual_machine_profile.0.non_admin_user.0.password").(string))
	if err != nil {
		return nil, err
	}
	virtualMachineProfile.NonAdminUser = nonAdminUser

	sku, err := flattenSku(&input.Sku)
	if err != nil {
		return nil, err
	}
	virtualMachineProfile.Sku = sku

	if input.UseSharedPassword != nil {
		virtualMachineProfile.SharedPasswordEnabled = *input.UseSharedPassword == lab.EnableStateEnabled
	}

	return append(virtualMachineProfiles, virtualMachineProfile), nil
}

func flattenAdditionalCapability(input *lab.VirtualMachineAdditionalCapabilities) ([]AdditionalCapability, error) {
	var additionalCapabilities []AdditionalCapability
	if input == nil {
		return additionalCapabilities, nil
	}

	additionalCapability := AdditionalCapability{}

	if input.InstallGpuDrivers != nil {
		additionalCapability.GpuDriversInstalled = *input.InstallGpuDrivers == lab.EnableStateEnabled
	}

	return append(additionalCapabilities, additionalCapability), nil
}

func flattenCredential(input *lab.Credentials, originalPassword string) ([]Credential, error) {
	var credentials []Credential
	if input == nil {
		return credentials, nil
	}

	credential := Credential{
		Username: input.Username,
		Password: originalPassword,
	}

	return append(credentials, credential), nil
}

func flattenImageReference(input *lab.ImageReference) ([]ImageReference, error) {
	var imageReferences []ImageReference
	if input == nil {
		return imageReferences, nil
	}

	imageReference := ImageReference{}

	if input.Id != nil {
		imageReference.Id = *input.Id
	}

	if input.Offer != nil {
		imageReference.Offer = *input.Offer
	}

	if input.Publisher != nil {
		imageReference.Publisher = *input.Publisher
	}

	if input.Sku != nil {
		imageReference.Sku = *input.Sku
	}

	if input.Version != nil {
		imageReference.Version = *input.Version
	}

	return append(imageReferences, imageReference), nil
}

func flattenSku(input *lab.Sku) ([]Sku, error) {
	var skus []Sku
	if input == nil {
		return skus, nil
	}

	sku := Sku{
		Name: input.Name,
	}

	if input.Capacity != nil {
		sku.Capacity = *input.Capacity
	}

	return append(skus, sku), nil
}
