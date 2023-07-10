// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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
	Name              string              `tfschema:"name"`
	ResourceGroupName string              `tfschema:"resource_group_name"`
	Location          string              `tfschema:"location"`
	AutoShutdown      []AutoShutdown      `tfschema:"auto_shutdown"`
	ConnectionSetting []ConnectionSetting `tfschema:"connection_setting"`
	Security          []Security          `tfschema:"security"`
	Title             string              `tfschema:"title"`
	VirtualMachine    []VirtualMachine    `tfschema:"virtual_machine"`
	Network           []Network           `tfschema:"network"`
	Roster            []Roster            `tfschema:"roster"`
	Description       string              `tfschema:"description"`
	LabPlanId         string              `tfschema:"lab_plan_id"`
	Tags              map[string]string   `tfschema:"tags"`
}

type AutoShutdown struct {
	DisconnectDelay string                 `tfschema:"disconnect_delay"`
	IdleDelay       string                 `tfschema:"idle_delay"`
	NoConnectDelay  string                 `tfschema:"no_connect_delay"`
	ShutdownOnIdle  lab.ShutdownOnIdleMode `tfschema:"shutdown_on_idle"`
}

type ConnectionSetting struct {
	ClientRdpAccess lab.ConnectionType `tfschema:"client_rdp_access"`
	ClientSshAccess lab.ConnectionType `tfschema:"client_ssh_access"`
}

type Security struct {
	OpenAccessEnabled bool   `tfschema:"open_access_enabled"`
	RegistrationCode  string `tfschema:"registration_code"`
}

type VirtualMachine struct {
	AdditionalCapabilityGpuDriversInstalled bool             `tfschema:"additional_capability_gpu_drivers_installed"`
	AdminUser                               []Credential     `tfschema:"admin_user"`
	CreateOption                            lab.CreateOption `tfschema:"create_option"`
	ImageReference                          []ImageReference `tfschema:"image_reference"`
	NonAdminUser                            []Credential     `tfschema:"non_admin_user"`
	Sku                                     []Sku            `tfschema:"sku"`
	UsageQuota                              string           `tfschema:"usage_quota"`
	SharedPasswordEnabled                   bool             `tfschema:"shared_password_enabled"`
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

type Network struct {
	SubnetId       string `tfschema:"subnet_id"`
	LoadBalancerId string `tfschema:"load_balancer_id"`
	PublicIPId     string `tfschema:"public_ip_id"`
}

type Roster struct {
	ActiveDirectoryGroupId string `tfschema:"active_directory_group_id"`
	LmsInstance            string `tfschema:"lms_instance"`
	LtiClientId            string `tfschema:"lti_client_id"`
	LtiContextId           string `tfschema:"lti_context_id"`
	LtiRosterEndpoint      string `tfschema:"lti_roster_endpoint"`
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

		"security": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"open_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"registration_code": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"title": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LabTitle,
		},

		"virtual_machine": {
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
									Sensitive:    true,
									ValidateFunc: validate.LabPassword,
								},
							},
						},
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
										"virtual_machine.0.image_reference.0.offer",
										"virtual_machine.0.image_reference.0.publisher",
										"virtual_machine.0.image_reference.0.sku",
										"virtual_machine.0.image_reference.0.version",
									},
								},

								"offer": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"virtual_machine.0.image_reference.0.id"},
								},

								"publisher": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"virtual_machine.0.image_reference.0.id"},
								},

								"sku": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"virtual_machine.0.image_reference.0.id"},
								},

								"version": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validate.LabImageVersion,
									ConflictsWith: []string{"virtual_machine.0.image_reference.0.id"},
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

								// Once `active_directory_group_id` is enabled, `capacity` wouldn't take effect, and it would be automatically set to the number of members in AAD Group by service API. So it has to use `ignore_changes` to avoid the difference of tf plan.
								// Service team confirmed they will fix this validation bug (https://github.com/Azure/azure-rest-api-specs/issues/22037) to allow `capacity` to be empty. Once the issue is fixed, `capacity` will be marked as `Optional`.
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
						Optional:     true,
						Default:      "PT0S",
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"additional_capability_gpu_drivers_installed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"create_option": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(lab.CreateOptionImage),
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.CreateOptionImage),
							string(lab.CreateOptionTemplateVM),
						}, false),
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
									ValidateFunc: validate.LabUsername,
								},

								"password": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validate.LabPassword,
								},
							},
						},
					},

					"shared_password_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"auto_shutdown": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disconnect_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"idle_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"no_connect_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"shutdown_on_idle": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ShutdownOnIdleModeLowUsage),
							string(lab.ShutdownOnIdleModeUserAbsence),
						}, false),
					},
				},
			},
		},

		"connection_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_rdp_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypePublic),
						}, false),
					},

					"client_ssh_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypePublic),
						}, false),
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

		"network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"load_balancer_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"public_ip_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"roster": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_directory_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},

					"lms_instance": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"lti_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"lti_context_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},

					"lti_roster_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r LabServiceLabResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServiceLabResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
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
					AutoShutdownProfile:   expandAutoShutdownProfile(model.AutoShutdown),
					ConnectionProfile:     expandConnectionProfile(model.ConnectionSetting),
					NetworkProfile:        expandNetworkProfile(model.Network, false, nil),
					RosterProfile:         expandRosterProfile(model.Roster),
					SecurityProfile:       expandSecurityProfile(model.Security),
					Title:                 &model.Title,
					VirtualMachineProfile: expandVirtualMachineProfile(model.VirtualMachine, false),
				},
				Tags: &model.Tags,
			}

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
		Timeout: 90 * time.Minute,
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

			props := resp.Model
			if props == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("auto_shutdown") {
				props.Properties.AutoShutdownProfile = expandAutoShutdownProfile(model.AutoShutdown)
			}

			if metadata.ResourceData.HasChange("connection_setting") {
				props.Properties.ConnectionProfile = expandConnectionProfile(model.ConnectionSetting)
			}

			if metadata.ResourceData.HasChange("security") {
				props.Properties.SecurityProfile = expandSecurityProfile(model.Security)
			}

			if metadata.ResourceData.HasChange("title") {
				props.Properties.Title = &model.Title
			}

			if metadata.ResourceData.HasChange("virtual_machine") {
				props.Properties.VirtualMachineProfile = expandVirtualMachineProfile(model.VirtualMachine, true)
			}

			if metadata.ResourceData.HasChange("network") {
				props.Properties.NetworkProfile = expandNetworkProfile(model.Network, true, props.Properties.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("roster") {
				props.Properties.RosterProfile = expandRosterProfile(model.Roster)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("lab_plan_id") {
				if model.LabPlanId != "" {
					props.Properties.LabPlanId = &model.LabPlanId
				} else {
					props.Properties.LabPlanId = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *props); err != nil {
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

			state := LabServiceLabModel{
				Name:              id.LabName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := resp.Model; model != nil {
				props := model.Properties

				state.Location = location.Normalize(model.Location)
				state.AutoShutdown = flattenAutoShutdownProfile(props.AutoShutdownProfile)
				state.ConnectionSetting = flattenConnectionProfile(props.ConnectionProfile)
				state.Network = flattenNetworkProfile(props.NetworkProfile)
				state.Roster = flattenRosterProfile(props.RosterProfile)
				state.Security = flattenSecurityProfile(props.SecurityProfile)
				state.VirtualMachine = flattenVirtualMachineProfile(props.VirtualMachineProfile, metadata.ResourceData)

				if props.Description != nil {
					state.Description = *props.Description
				}
				if props.LabPlanId != nil {
					state.LabPlanId = *props.LabPlanId
				}
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				if props.Title != nil {
					state.Title = *props.Title
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServiceLabResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
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

			if oldVal, newVal := rd.GetChange("virtual_machine.0.non_admin_user"); oldVal != nil && newVal != nil && (len(oldVal.([]interface{})) == 0 && len(newVal.([]interface{})) == 1) {
				if err := rd.ForceNew("virtual_machine.0.non_admin_user"); err != nil {
					return err
				}
			}

			if oldVal, newVal := rd.GetChange("virtual_machine.0.non_admin_user"); oldVal != nil && newVal != nil && (len(oldVal.([]interface{})) == 1 && len(newVal.([]interface{})) == 1) {
				if ok := rd.HasChange("virtual_machine.0.non_admin_user.0.username"); ok {
					if err := rd.ForceNew("virtual_machine.0.non_admin_user.0.username"); err != nil {
						return err
					}
				}

				if ok := rd.HasChange("virtual_machine.0.non_admin_user.0.password"); ok {
					if err := rd.ForceNew("virtual_machine.0.non_admin_user.0.password"); err != nil {
						return err
					}
				}
			}

			if oldVal, newVal := rd.GetChange("network"); oldVal != nil && newVal != nil && (len(oldVal.([]interface{})) == 0 && len(newVal.([]interface{})) == 1) {
				if err := rd.ForceNew("network"); err != nil {
					return err
				}
			}

			if oldVal, newVal := rd.GetChange("network"); oldVal != nil && newVal != nil && (len(oldVal.([]interface{})) == 1 && len(newVal.([]interface{})) == 1) {
				if ok := rd.HasChange("network.0.subnet_id"); ok {
					if err := rd.ForceNew("network.0.subnet_id"); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func expandAutoShutdownProfile(input []AutoShutdown) lab.AutoShutdownProfile {
	shutdownOnDisconnectEnabled := lab.EnableStateDisabled
	shutdownWhenNotConnectedEnabled := lab.EnableStateDisabled
	shutdownOnIdle := lab.ShutdownOnIdleModeNone

	if len(input) == 0 {
		return lab.AutoShutdownProfile{
			ShutdownOnDisconnect:     &shutdownOnDisconnectEnabled,
			ShutdownWhenNotConnected: &shutdownWhenNotConnectedEnabled,
			ShutdownOnIdle:           &shutdownOnIdle,
		}
	}

	autoShutdownProfile := input[0]
	result := lab.AutoShutdownProfile{}

	if autoShutdownProfile.DisconnectDelay != "" {
		shutdownOnDisconnectEnabled = lab.EnableStateEnabled
		result.DisconnectDelay = &autoShutdownProfile.DisconnectDelay
	}
	result.ShutdownOnDisconnect = &shutdownOnDisconnectEnabled

	if autoShutdownProfile.IdleDelay != "" {
		result.IdleDelay = &autoShutdownProfile.IdleDelay
	}

	if autoShutdownProfile.NoConnectDelay != "" {
		shutdownWhenNotConnectedEnabled = lab.EnableStateEnabled
		result.NoConnectDelay = &autoShutdownProfile.NoConnectDelay
	}
	result.ShutdownWhenNotConnected = &shutdownWhenNotConnectedEnabled

	if autoShutdownProfile.ShutdownOnIdle != "" {
		shutdownOnIdle = autoShutdownProfile.ShutdownOnIdle
	}
	result.ShutdownOnIdle = &shutdownOnIdle

	return result
}

func flattenAutoShutdownProfile(input lab.AutoShutdownProfile) []AutoShutdown {
	// default values
	shutdownOnDisconnectEnabled := input.ShutdownOnDisconnect != nil && *input.ShutdownOnDisconnect != lab.EnableStateDisabled
	shutdownWhenNotConnectedEnabled := input.ShutdownWhenNotConnected != nil && *input.ShutdownWhenNotConnected != lab.EnableStateDisabled
	shutdownOnIdleModeConfigured := input.ShutdownOnIdle != nil && *input.ShutdownOnIdle != lab.ShutdownOnIdleModeNone
	idleDelayConfigured := input.IdleDelay != nil
	if !shutdownOnDisconnectEnabled && !shutdownWhenNotConnectedEnabled && !shutdownOnIdleModeConfigured && !idleDelayConfigured {
		return []AutoShutdown{}
	}

	var autoShutdownProfiles []AutoShutdown
	autoShutdownProfile := AutoShutdown{}

	if input.DisconnectDelay != nil {
		autoShutdownProfile.DisconnectDelay = *input.DisconnectDelay
	}

	if input.IdleDelay != nil {
		autoShutdownProfile.IdleDelay = *input.IdleDelay
	}

	if input.NoConnectDelay != nil {
		autoShutdownProfile.NoConnectDelay = *input.NoConnectDelay
	}

	if shutdownOnIdle := input.ShutdownOnIdle; shutdownOnIdle != nil && *shutdownOnIdle != lab.ShutdownOnIdleModeNone {
		autoShutdownProfile.ShutdownOnIdle = *shutdownOnIdle
	}

	return append(autoShutdownProfiles, autoShutdownProfile)
}

func expandConnectionProfile(input []ConnectionSetting) lab.ConnectionProfile {
	clientRdpAccess := lab.ConnectionTypeNone
	clientSshAccess := lab.ConnectionTypeNone
	webRdpAccess := lab.ConnectionTypeNone
	webSshAccess := lab.ConnectionTypeNone

	if len(input) == 0 {
		return lab.ConnectionProfile{
			ClientRdpAccess: &clientRdpAccess,
			ClientSshAccess: &clientSshAccess,
			WebRdpAccess:    &webRdpAccess,
			WebSshAccess:    &webSshAccess,
		}
	}

	connectionProfile := input[0]
	result := lab.ConnectionProfile{
		WebRdpAccess: &webRdpAccess,
		WebSshAccess: &webSshAccess,
	}

	if connectionProfile.ClientRdpAccess != "" {
		clientRdpAccess = connectionProfile.ClientRdpAccess
	}
	result.ClientRdpAccess = &clientRdpAccess

	if connectionProfile.ClientSshAccess != "" {
		clientSshAccess = connectionProfile.ClientSshAccess
	}
	result.ClientSshAccess = &clientSshAccess

	return result
}

func flattenConnectionProfile(input lab.ConnectionProfile) []ConnectionSetting {
	if (input.ClientRdpAccess == nil || *input.ClientRdpAccess == lab.ConnectionTypeNone) && (input.ClientSshAccess == nil || *input.ClientSshAccess == lab.ConnectionTypeNone) {
		return []ConnectionSetting{}
	}

	connectionProfile := ConnectionSetting{}

	if clientRdpAccess := input.ClientRdpAccess; clientRdpAccess != nil && *clientRdpAccess != lab.ConnectionTypeNone {
		connectionProfile.ClientRdpAccess = *clientRdpAccess
	}

	if clientSshAccess := input.ClientSshAccess; clientSshAccess != nil && *clientSshAccess != lab.ConnectionTypeNone {
		connectionProfile.ClientSshAccess = *clientSshAccess
	}

	return []ConnectionSetting{
		connectionProfile,
	}
}

func expandSecurityProfile(input []Security) lab.SecurityProfile {
	result := lab.SecurityProfile{}

	if len(input) == 0 {
		return result
	}

	securityProfile := input[0]

	openAccessEnabled := lab.EnableStateEnabled
	if !securityProfile.OpenAccessEnabled {
		openAccessEnabled = lab.EnableStateDisabled
	}
	result.OpenAccess = &openAccessEnabled

	return result
}

func flattenSecurityProfile(input lab.SecurityProfile) []Security {
	var securityProfiles []Security
	securityProfile := Security{}

	if input.OpenAccess != nil {
		securityProfile.OpenAccessEnabled = *input.OpenAccess == lab.EnableStateEnabled
	}

	if input.RegistrationCode != nil {
		securityProfile.RegistrationCode = *input.RegistrationCode
	}

	return append(securityProfiles, securityProfile)
}

func expandVirtualMachineProfile(input []VirtualMachine, isUpdate bool) lab.VirtualMachineProfile {
	result := lab.VirtualMachineProfile{}

	if len(input) == 0 {
		return result
	}

	virtualMachineProfile := input[0]
	result.AdminUser = *expandCredential(virtualMachineProfile.AdminUser, isUpdate)
	result.CreateOption = virtualMachineProfile.CreateOption
	result.UsageQuota = virtualMachineProfile.UsageQuota
	result.Sku = expandSku(virtualMachineProfile.Sku)
	result.NonAdminUser = expandCredential(virtualMachineProfile.NonAdminUser, isUpdate)
	result.ImageReference = expandImageReference(virtualMachineProfile.ImageReference)

	sharedPasswordEnabled := lab.EnableStateDisabled
	if virtualMachineProfile.SharedPasswordEnabled {
		sharedPasswordEnabled = lab.EnableStateEnabled
	}
	result.UseSharedPassword = &sharedPasswordEnabled

	additionalCapabilityGpuDriversInstalled := lab.EnableStateDisabled
	if virtualMachineProfile.AdditionalCapabilityGpuDriversInstalled {
		additionalCapabilityGpuDriversInstalled = lab.EnableStateEnabled
	}
	result.AdditionalCapabilities = &lab.VirtualMachineAdditionalCapabilities{
		InstallGpuDrivers: &additionalCapabilityGpuDriversInstalled,
	}

	return result
}

func expandCredential(input []Credential, isUpdate bool) *lab.Credentials {
	if len(input) == 0 {
		return nil
	}

	credential := input[0]
	result := lab.Credentials{
		Username: credential.Username,
	}

	if !isUpdate && credential.Password != "" {
		result.Password = &credential.Password
	}

	return &result
}

func expandImageReference(input []ImageReference) lab.ImageReference {
	result := lab.ImageReference{}

	if len(input) == 0 {
		return result
	}

	imageReference := input[0]

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

	return result
}

func expandSku(input []Sku) lab.Sku {
	result := lab.Sku{}

	if len(input) == 0 {
		return result
	}

	sku := input[0]
	result.Name = sku.Name
	result.Capacity = &sku.Capacity

	return result
}

func flattenVirtualMachineProfile(input lab.VirtualMachineProfile, d *pluginsdk.ResourceData) []VirtualMachine {
	var virtualMachineProfiles []VirtualMachine

	virtualMachineProfile := VirtualMachine{
		AdminUser:      flattenCredential(&input.AdminUser, d.Get("virtual_machine.0.admin_user.0.password").(string)),
		CreateOption:   input.CreateOption,
		ImageReference: flattenImageReference(&input.ImageReference),
		Sku:            flattenSku(&input.Sku),
		UsageQuota:     input.UsageQuota,
	}

	if input.AdditionalCapabilities != nil && *input.AdditionalCapabilities.InstallGpuDrivers != "" {
		virtualMachineProfile.AdditionalCapabilityGpuDriversInstalled = *input.AdditionalCapabilities.InstallGpuDrivers == lab.EnableStateEnabled
	}

	if input.NonAdminUser != nil {
		virtualMachineProfile.NonAdminUser = flattenCredential(input.NonAdminUser, d.Get("virtual_machine.0.non_admin_user.0.password").(string))
	}

	if input.UseSharedPassword != nil {
		virtualMachineProfile.SharedPasswordEnabled = *input.UseSharedPassword == lab.EnableStateEnabled
	}

	return append(virtualMachineProfiles, virtualMachineProfile)
}

func flattenCredential(input *lab.Credentials, originalPassword string) []Credential {
	if input == nil {
		return []Credential{}
	}

	var credentials []Credential

	credential := Credential{
		Username: input.Username,
		Password: originalPassword,
	}

	return append(credentials, credential)
}

func flattenImageReference(input *lab.ImageReference) []ImageReference {
	if input == nil {
		return []ImageReference{}
	}

	var imageReferences []ImageReference
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

	return append(imageReferences, imageReference)
}

func flattenSku(input *lab.Sku) []Sku {
	if input == nil {
		return []Sku{}
	}

	var skus []Sku
	sku := Sku{
		Name:     input.Name,
		Capacity: *input.Capacity,
	}

	return append(skus, sku)
}

func expandNetworkProfile(input []Network, isUpdate bool, existingNetwork *lab.LabNetworkProfile) *lab.LabNetworkProfile {
	if len(input) == 0 {
		return nil
	}

	networkProfile := input[0]
	result := lab.LabNetworkProfile{}

	if networkProfile.SubnetId != "" {
		result.SubnetId = &networkProfile.SubnetId

		if isUpdate && existingNetwork != nil {
			result.LoadBalancerId = existingNetwork.LoadBalancerId
			result.PublicIPId = existingNetwork.PublicIPId
		}
	}

	return &result
}

func flattenNetworkProfile(input *lab.LabNetworkProfile) []Network {
	if input == nil {
		return []Network{}
	}

	var networkProfiles []Network
	networkProfile := Network{}

	if input.SubnetId != nil {
		networkProfile.SubnetId = *input.SubnetId
	}

	if input.LoadBalancerId != nil {
		networkProfile.LoadBalancerId = *input.LoadBalancerId
	}

	if input.PublicIPId != nil {
		networkProfile.PublicIPId = *input.PublicIPId
	}

	return append(networkProfiles, networkProfile)
}

func expandRosterProfile(input []Roster) *lab.RosterProfile {
	if len(input) == 0 {
		return nil
	}

	rosterProfile := input[0]
	result := lab.RosterProfile{}

	if rosterProfile.ActiveDirectoryGroupId != "" {
		result.ActiveDirectoryGroupId = &rosterProfile.ActiveDirectoryGroupId
	}

	if rosterProfile.LmsInstance != "" {
		result.LmsInstance = &rosterProfile.LmsInstance
	}

	if rosterProfile.LtiClientId != "" {
		result.LtiClientId = &rosterProfile.LtiClientId
	}

	if rosterProfile.LtiContextId != "" {
		result.LtiContextId = &rosterProfile.LtiContextId
	}

	if rosterProfile.LtiRosterEndpoint != "" {
		result.LtiRosterEndpoint = &rosterProfile.LtiRosterEndpoint
	}

	return &result
}

func flattenRosterProfile(input *lab.RosterProfile) []Roster {
	if input == nil {
		return []Roster{}
	}

	var rosterProfiles []Roster
	rosterProfile := Roster{}

	if input.ActiveDirectoryGroupId != nil {
		rosterProfile.ActiveDirectoryGroupId = *input.ActiveDirectoryGroupId
	}

	if input.LmsInstance != nil {
		rosterProfile.LmsInstance = *input.LmsInstance
	}

	if input.LtiClientId != nil {
		rosterProfile.LtiClientId = *input.LtiClientId
	}

	if input.LtiContextId != nil {
		rosterProfile.LtiContextId = *input.LtiContextId
	}

	if input.LtiRosterEndpoint != nil {
		rosterProfile.LtiRosterEndpoint = *input.LtiRosterEndpoint
	}

	return append(rosterProfiles, rosterProfile)
}
