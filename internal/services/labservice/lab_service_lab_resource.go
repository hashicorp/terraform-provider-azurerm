package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LabServicesLabModel struct {
	Name                  string                       `tfschema:"name"`
	ResourceGroupName     string                       `tfschema:"resource_group_name"`
	AutoShutdownProfile   []AutoShutdownProfileModel   `tfschema:"auto_shutdown_profile"`
	ConnectionProfile     []ConnectionProfileModel     `tfschema:"connection_profile"`
	Description           string                       `tfschema:"description"`
	LabPlanId             string                       `tfschema:"lab_plan_id"`
	Location              string                       `tfschema:"location"`
	NetworkProfile        []LabNetworkProfileModel     `tfschema:"network_profile"`
	RosterProfile         []RosterProfileModel         `tfschema:"roster_profile"`
	SecurityProfile       []SecurityProfileModel       `tfschema:"security_profile"`
	Tags                  map[string]string            `tfschema:"tags"`
	Title                 string                       `tfschema:"title"`
	VirtualMachineProfile []VirtualMachineProfileModel `tfschema:"virtual_machine_profile"`
	State                 lab.LabState                 `tfschema:"state"`
}

type AutoShutdownProfileModel struct {
	DisconnectDelay          string                 `tfschema:"disconnect_delay"`
	IdleDelay                string                 `tfschema:"idle_delay"`
	NoConnectDelay           string                 `tfschema:"no_connect_delay"`
	ShutdownOnDisconnect     lab.EnableState        `tfschema:"shutdown_on_disconnect"`
	ShutdownOnIdle           lab.ShutdownOnIdleMode `tfschema:"shutdown_on_idle"`
	ShutdownWhenNotConnected lab.EnableState        `tfschema:"shutdown_when_not_connected"`
}

type ConnectionProfileModel struct {
	ClientRdpAccess lab.ConnectionType `tfschema:"client_rdp_access"`
	ClientSshAccess lab.ConnectionType `tfschema:"client_ssh_access"`
	WebRdpAccess    lab.ConnectionType `tfschema:"web_rdp_access"`
	WebSshAccess    lab.ConnectionType `tfschema:"web_ssh_access"`
}

type LabNetworkProfileModel struct {
	LoadBalancerId string `tfschema:"load_balancer_id"`
	PublicIPId     string `tfschema:"public_ip_id"`
	SubnetId       string `tfschema:"subnet_id"`
}

type RosterProfileModel struct {
	ActiveDirectoryGroupId string `tfschema:"active_directory_group_id"`
	LmsInstance            string `tfschema:"lms_instance"`
	LtiClientId            string `tfschema:"lti_client_id"`
	LtiContextId           string `tfschema:"lti_context_id"`
	LtiRosterEndpoint      string `tfschema:"lti_roster_endpoint"`
}

type SecurityProfileModel struct {
	OpenAccess       lab.EnableState `tfschema:"open_access"`
	RegistrationCode string          `tfschema:"registration_code"`
}

type VirtualMachineProfileModel struct {
	AdditionalCapabilities []VirtualMachineAdditionalCapabilitiesModel `tfschema:"additional_capabilities"`
	AdminUser              []CredentialsModel                          `tfschema:"admin_user"`
	CreateOption           lab.CreateOption                            `tfschema:"create_option"`
	ImageReference         []ImageReferenceModel                       `tfschema:"image_reference"`
	NonAdminUser           []CredentialsModel                          `tfschema:"non_admin_user"`
	OsType                 lab.OsType                                  `tfschema:"os_type"`
	Sku                    []SkuModel                                  `tfschema:"sku"`
	UsageQuota             string                                      `tfschema:"usage_quota"`
	UseSharedPassword      lab.EnableState                             `tfschema:"use_shared_password"`
}

type VirtualMachineAdditionalCapabilitiesModel struct {
	InstallGpuDrivers lab.EnableState `tfschema:"install_gpu_drivers"`
}

type CredentialsModel struct {
	Password string `tfschema:"password"`
	Username string `tfschema:"username"`
}

type ImageReferenceModel struct {
	ExactVersion string `tfschema:"exact_version"`
	Id           string `tfschema:"id"`
	Offer        string `tfschema:"offer"`
	Publisher    string `tfschema:"publisher"`
	Sku          string `tfschema:"sku"`
	Version      string `tfschema:"version"`
}

type SkuModel struct {
	Capacity int64       `tfschema:"capacity"`
	Family   string      `tfschema:"family"`
	Name     string      `tfschema:"name"`
	Size     string      `tfschema:"size"`
	Tier     lab.SkuTier `tfschema:"tier"`
}

type LabServicesLabResource struct{}

var _ sdk.ResourceWithUpdate = LabServicesLabResource{}

func (r LabServicesLabResource) ResourceType() string {
	return "azurerm_lab_services_lab"
}

func (r LabServicesLabResource) ModelObject() interface{} {
	return &LabServicesLabModel{}
}

func (r LabServicesLabResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return lab.ValidateLabID
}

func (r LabServicesLabResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"auto_shutdown_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disconnect_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"idle_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"no_connect_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"shutdown_on_disconnect": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.EnableStateEnabled),
							string(lab.EnableStateDisabled),
						}, false),
					},

					"shutdown_on_idle": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ShutdownOnIdleModeNone),
							string(lab.ShutdownOnIdleModeUserAbsence),
							string(lab.ShutdownOnIdleModeLowUsage),
						}, false),
					},

					"shutdown_when_not_connected": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.EnableStateEnabled),
							string(lab.EnableStateDisabled),
						}, false),
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
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypePublic),
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypeNone),
						}, false),
					},

					"client_ssh_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypeNone),
							string(lab.ConnectionTypePublic),
						}, false),
					},

					"web_rdp_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypePublic),
							string(lab.ConnectionTypePrivate),
							string(lab.ConnectionTypeNone),
						}, false),
					},

					"web_ssh_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.ConnectionTypeNone),
							string(lab.ConnectionTypePublic),
							string(lab.ConnectionTypePrivate),
						}, false),
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"lab_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"load_balancer_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"public_ip_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"roster_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_directory_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"lms_instance": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"lti_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"lti_context_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"lti_roster_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
					"open_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.EnableStateEnabled),
							string(lab.EnableStateDisabled),
						}, false),
					},

					"registration_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"title": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_machine_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_capabilities": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"install_gpu_drivers": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(lab.EnableStateEnabled),
										string(lab.EnableStateDisabled),
									}, false),
								},
							},
						},
					},

					"admin_user": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"username": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"create_option": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
								"exact_version": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"offer": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"publisher": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"sku": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"version": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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
								"password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"username": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"os_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.OsTypeWindows),
							string(lab.OsTypeLinux),
						}, false),
					},

					"sku": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"capacity": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"family": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"size": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"tier": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(lab.SkuTierFree),
										string(lab.SkuTierBasic),
										string(lab.SkuTierStandard),
										string(lab.SkuTierPremium),
									}, false),
								},
							},
						},
					},

					"usage_quota": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"use_shared_password": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(lab.EnableStateEnabled),
							string(lab.EnableStateDisabled),
						}, false),
					},
				},
			},
		},
	}
}

func (r LabServicesLabResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r LabServicesLabResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServicesLabModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LabServices.LabClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := lab.NewLabID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &lab.Lab{
				Location:   location.Normalize(model.Location),
				Properties: lab.LabProperties{},
				Tags:       &model.Tags,
			}

			autoShutdownProfileValue, err := expandAutoShutdownProfileModel(model.AutoShutdownProfile)
			if err != nil {
				return err
			}

			if autoShutdownProfileValue != nil {
				properties.Properties.AutoShutdownProfile = *autoShutdownProfileValue
			}

			connectionProfileValue, err := expandConnectionProfileModel(model.ConnectionProfile)
			if err != nil {
				return err
			}

			if connectionProfileValue != nil {
				properties.Properties.ConnectionProfile = *connectionProfileValue
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.LabPlanId != "" {
				properties.Properties.LabPlanId = &model.LabPlanId
			}

			rosterProfileValue, err := expandRosterProfileModel(model.RosterProfile)
			if err != nil {
				return err
			}

			properties.Properties.RosterProfile = rosterProfileValue

			securityProfileValue, err := expandSecurityProfileModel(model.SecurityProfile)
			if err != nil {
				return err
			}

			if securityProfileValue != nil {
				properties.Properties.SecurityProfile = *securityProfileValue
			}

			if model.Title != "" {
				properties.Properties.Title = &model.Title
			}

			virtualMachineProfileValue, err := expandVirtualMachineProfileModel(model.VirtualMachineProfile)
			if err != nil {
				return err
			}

			if virtualMachineProfileValue != nil {
				properties.Properties.VirtualMachineProfile = *virtualMachineProfileValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LabServicesLabResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabServices.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServicesLabModel
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
				autoShutdownProfileValue, err := expandAutoShutdownProfileModel(model.AutoShutdownProfile)
				if err != nil {
					return err
				}

				if autoShutdownProfileValue != nil {
					properties.Properties.AutoShutdownProfile = *autoShutdownProfileValue
				}
			}

			if metadata.ResourceData.HasChange("connection_profile") {
				connectionProfileValue, err := expandConnectionProfileModel(model.ConnectionProfile)
				if err != nil {
					return err
				}

				if connectionProfileValue != nil {
					properties.Properties.ConnectionProfile = *connectionProfileValue
				}
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

			if metadata.ResourceData.HasChange("roster_profile") {
				rosterProfileValue, err := expandRosterProfileModel(model.RosterProfile)
				if err != nil {
					return err
				}

				properties.Properties.RosterProfile = rosterProfileValue
			}

			if metadata.ResourceData.HasChange("security_profile") {
				securityProfileValue, err := expandSecurityProfileModel(model.SecurityProfile)
				if err != nil {
					return err
				}

				if securityProfileValue != nil {
					properties.Properties.SecurityProfile = *securityProfileValue
				}
			}

			if metadata.ResourceData.HasChange("title") {
				if model.Title != "" {
					properties.Properties.Title = &model.Title
				} else {
					properties.Properties.Title = nil
				}
			}

			if metadata.ResourceData.HasChange("virtual_machine_profile") {
				virtualMachineProfileValue, err := expandVirtualMachineProfileModel(model.VirtualMachineProfile)
				if err != nil {
					return err
				}

				if virtualMachineProfileValue != nil {
					properties.Properties.VirtualMachineProfile = *virtualMachineProfileValue
				}
			}

			properties.SystemData = nil

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

func (r LabServicesLabResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabServices.LabClient

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

			state := LabServicesLabModel{
				Name:              id.LabName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties
			autoShutdownProfileValue, err := flattenAutoShutdownProfileModel(&properties.AutoShutdownProfile)
			if err != nil {
				return err
			}

			state.AutoShutdownProfile = autoShutdownProfileValue

			connectionProfileValue, err := flattenConnectionProfileModel(&properties.ConnectionProfile)
			if err != nil {
				return err
			}

			state.ConnectionProfile = connectionProfileValue

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.LabPlanId != nil {
				state.LabPlanId = *properties.LabPlanId
			}

			networkProfileValue, err := flattenLabNetworkProfileModel(properties.NetworkProfile)
			if err != nil {
				return err
			}

			state.NetworkProfile = networkProfileValue

			rosterProfileValue, err := flattenRosterProfileModel(properties.RosterProfile)
			if err != nil {
				return err
			}

			state.RosterProfile = rosterProfileValue

			securityProfileValue, err := flattenSecurityProfileModel(&properties.SecurityProfile)
			if err != nil {
				return err
			}

			state.SecurityProfile = securityProfileValue

			if properties.State != nil {
				state.State = *properties.State
			}

			if properties.Title != nil {
				state.Title = *properties.Title
			}

			virtualMachineProfileValue, err := flattenVirtualMachineProfileModel(&properties.VirtualMachineProfile)
			if err != nil {
				return err
			}

			state.VirtualMachineProfile = virtualMachineProfileValue
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServicesLabResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabServices.LabClient

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

func expandAutoShutdownProfileModel(inputList []AutoShutdownProfileModel) (*lab.AutoShutdownProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.AutoShutdownProfile{
		ShutdownOnDisconnect:     &input.ShutdownOnDisconnect,
		ShutdownOnIdle:           &input.ShutdownOnIdle,
		ShutdownWhenNotConnected: &input.ShutdownWhenNotConnected,
	}

	if input.DisconnectDelay != "" {
		output.DisconnectDelay = &input.DisconnectDelay
	}

	if input.IdleDelay != "" {
		output.IdleDelay = &input.IdleDelay
	}

	if input.NoConnectDelay != "" {
		output.NoConnectDelay = &input.NoConnectDelay
	}

	return &output, nil
}

func expandConnectionProfileModel(inputList []ConnectionProfileModel) (*lab.ConnectionProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.ConnectionProfile{
		ClientRdpAccess: &input.ClientRdpAccess,
		ClientSshAccess: &input.ClientSshAccess,
		WebRdpAccess:    &input.WebRdpAccess,
		WebSshAccess:    &input.WebSshAccess,
	}

	return &output, nil
}

func expandLabNetworkProfileModel(inputList []LabNetworkProfileModel) (*lab.LabNetworkProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.LabNetworkProfile{}

	if input.LoadBalancerId != "" {
		output.LoadBalancerId = &input.LoadBalancerId
	}

	if input.PublicIPId != "" {
		output.PublicIPId = &input.PublicIPId
	}

	if input.SubnetId != "" {
		output.SubnetId = &input.SubnetId
	}

	return &output, nil
}

func expandRosterProfileModel(inputList []RosterProfileModel) (*lab.RosterProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.RosterProfile{}

	if input.ActiveDirectoryGroupId != "" {
		output.ActiveDirectoryGroupId = &input.ActiveDirectoryGroupId
	}

	if input.LmsInstance != "" {
		output.LmsInstance = &input.LmsInstance
	}

	if input.LtiClientId != "" {
		output.LtiClientId = &input.LtiClientId
	}

	if input.LtiContextId != "" {
		output.LtiContextId = &input.LtiContextId
	}

	if input.LtiRosterEndpoint != "" {
		output.LtiRosterEndpoint = &input.LtiRosterEndpoint
	}

	return &output, nil
}

func expandSecurityProfileModel(inputList []SecurityProfileModel) (*lab.SecurityProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.SecurityProfile{
		OpenAccess: &input.OpenAccess,
	}

	return &output, nil
}

func expandVirtualMachineProfileModel(inputList []VirtualMachineProfileModel) (*lab.VirtualMachineProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.VirtualMachineProfile{
		CreateOption:      input.CreateOption,
		OsType:            &input.OsType,
		UsageQuota:        input.UsageQuota,
		UseSharedPassword: &input.UseSharedPassword,
	}

	additionalCapabilitiesValue, err := expandVirtualMachineAdditionalCapabilitiesModel(input.AdditionalCapabilities)
	if err != nil {
		return nil, err
	}

	output.AdditionalCapabilities = additionalCapabilitiesValue

	adminUserValue, err := expandCredentialsModel(input.AdminUser)
	if err != nil {
		return nil, err
	}

	if adminUserValue != nil {
		output.AdminUser = *adminUserValue
	}

	imageReferenceValue, err := expandImageReferenceModel(input.ImageReference)
	if err != nil {
		return nil, err
	}

	if imageReferenceValue != nil {
		output.ImageReference = *imageReferenceValue
	}

	nonAdminUserValue, err := expandCredentialsModel(input.NonAdminUser)
	if err != nil {
		return nil, err
	}

	output.NonAdminUser = nonAdminUserValue

	skuValue, err := expandSkuModel(input.Sku)
	if err != nil {
		return nil, err
	}

	if skuValue != nil {
		output.Sku = *skuValue
	}

	return &output, nil
}

func expandVirtualMachineAdditionalCapabilitiesModel(inputList []VirtualMachineAdditionalCapabilitiesModel) (*lab.VirtualMachineAdditionalCapabilities, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.VirtualMachineAdditionalCapabilities{
		InstallGpuDrivers: &input.InstallGpuDrivers,
	}

	return &output, nil
}

func expandCredentialsModel(inputList []CredentialsModel) (*lab.Credentials, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.Credentials{
		Username: input.Username,
	}

	if input.Password != "" {
		output.Password = &input.Password
	}

	return &output, nil
}

func expandImageReferenceModel(inputList []ImageReferenceModel) (*lab.ImageReference, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.ImageReference{}

	if input.Id != "" {
		output.Id = &input.Id
	}

	if input.Offer != "" {
		output.Offer = &input.Offer
	}

	if input.Publisher != "" {
		output.Publisher = &input.Publisher
	}

	if input.Sku != "" {
		output.Sku = &input.Sku
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output, nil
}

func expandSkuModel(inputList []SkuModel) (*lab.Sku, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := lab.Sku{
		Capacity: &input.Capacity,
		Name:     input.Name,
		Tier:     &input.Tier,
	}

	if input.Family != "" {
		output.Family = &input.Family
	}

	if input.Size != "" {
		output.Size = &input.Size
	}

	return &output, nil
}

func flattenAutoShutdownProfileModel(input *lab.AutoShutdownProfile) ([]AutoShutdownProfileModel, error) {
	var outputList []AutoShutdownProfileModel
	if input == nil {
		return outputList, nil
	}

	output := AutoShutdownProfileModel{}

	if input.DisconnectDelay != nil {
		output.DisconnectDelay = *input.DisconnectDelay
	}

	if input.IdleDelay != nil {
		output.IdleDelay = *input.IdleDelay
	}

	if input.NoConnectDelay != nil {
		output.NoConnectDelay = *input.NoConnectDelay
	}

	if input.ShutdownOnDisconnect != nil {
		output.ShutdownOnDisconnect = *input.ShutdownOnDisconnect
	}

	if input.ShutdownOnIdle != nil {
		output.ShutdownOnIdle = *input.ShutdownOnIdle
	}

	if input.ShutdownWhenNotConnected != nil {
		output.ShutdownWhenNotConnected = *input.ShutdownWhenNotConnected
	}

	return append(outputList, output), nil
}

func flattenConnectionProfileModel(input *lab.ConnectionProfile) ([]ConnectionProfileModel, error) {
	var outputList []ConnectionProfileModel
	if input == nil {
		return outputList, nil
	}

	output := ConnectionProfileModel{}

	if input.ClientRdpAccess != nil {
		output.ClientRdpAccess = *input.ClientRdpAccess
	}

	if input.ClientSshAccess != nil {
		output.ClientSshAccess = *input.ClientSshAccess
	}

	if input.WebRdpAccess != nil {
		output.WebRdpAccess = *input.WebRdpAccess
	}

	if input.WebSshAccess != nil {
		output.WebSshAccess = *input.WebSshAccess
	}

	return append(outputList, output), nil
}

func flattenLabNetworkProfileModel(input *lab.LabNetworkProfile) ([]LabNetworkProfileModel, error) {
	var outputList []LabNetworkProfileModel
	if input == nil {
		return outputList, nil
	}

	output := LabNetworkProfileModel{}

	if input.LoadBalancerId != nil {
		output.LoadBalancerId = *input.LoadBalancerId
	}

	if input.PublicIPId != nil {
		output.PublicIPId = *input.PublicIPId
	}

	if input.SubnetId != nil {
		output.SubnetId = *input.SubnetId
	}

	return append(outputList, output), nil
}

func flattenRosterProfileModel(input *lab.RosterProfile) ([]RosterProfileModel, error) {
	var outputList []RosterProfileModel
	if input == nil {
		return outputList, nil
	}

	output := RosterProfileModel{}

	if input.ActiveDirectoryGroupId != nil {
		output.ActiveDirectoryGroupId = *input.ActiveDirectoryGroupId
	}

	if input.LmsInstance != nil {
		output.LmsInstance = *input.LmsInstance
	}

	if input.LtiClientId != nil {
		output.LtiClientId = *input.LtiClientId
	}

	if input.LtiContextId != nil {
		output.LtiContextId = *input.LtiContextId
	}

	if input.LtiRosterEndpoint != nil {
		output.LtiRosterEndpoint = *input.LtiRosterEndpoint
	}

	return append(outputList, output), nil
}

func flattenSecurityProfileModel(input *lab.SecurityProfile) ([]SecurityProfileModel, error) {
	var outputList []SecurityProfileModel
	if input == nil {
		return outputList, nil
	}

	output := SecurityProfileModel{}

	if input.OpenAccess != nil {
		output.OpenAccess = *input.OpenAccess
	}

	if input.RegistrationCode != nil {
		output.RegistrationCode = *input.RegistrationCode
	}

	return append(outputList, output), nil
}

func flattenVirtualMachineProfileModel(input *lab.VirtualMachineProfile) ([]VirtualMachineProfileModel, error) {
	var outputList []VirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}

	output := VirtualMachineProfileModel{
		CreateOption: input.CreateOption,
		UsageQuota:   input.UsageQuota,
	}

	additionalCapabilitiesValue, err := flattenVirtualMachineAdditionalCapabilitiesModel(input.AdditionalCapabilities)
	if err != nil {
		return nil, err
	}

	output.AdditionalCapabilities = additionalCapabilitiesValue

	adminUserValue, err := flattenCredentialsModel(&input.AdminUser)
	if err != nil {
		return nil, err
	}

	output.AdminUser = adminUserValue

	imageReferenceValue, err := flattenImageReferenceModel(&input.ImageReference)
	if err != nil {
		return nil, err
	}

	output.ImageReference = imageReferenceValue

	nonAdminUserValue, err := flattenCredentialsModel(input.NonAdminUser)
	if err != nil {
		return nil, err
	}

	output.NonAdminUser = nonAdminUserValue

	if input.OsType != nil {
		output.OsType = *input.OsType
	}

	skuValue, err := flattenSkuModel(&input.Sku)
	if err != nil {
		return nil, err
	}

	output.Sku = skuValue

	if input.UseSharedPassword != nil {
		output.UseSharedPassword = *input.UseSharedPassword
	}

	return append(outputList, output), nil
}

func flattenVirtualMachineAdditionalCapabilitiesModel(input *lab.VirtualMachineAdditionalCapabilities) ([]VirtualMachineAdditionalCapabilitiesModel, error) {
	var outputList []VirtualMachineAdditionalCapabilitiesModel
	if input == nil {
		return outputList, nil
	}

	output := VirtualMachineAdditionalCapabilitiesModel{}

	if input.InstallGpuDrivers != nil {
		output.InstallGpuDrivers = *input.InstallGpuDrivers
	}

	return append(outputList, output), nil
}

func flattenCredentialsModel(input *lab.Credentials) ([]CredentialsModel, error) {
	var outputList []CredentialsModel
	if input == nil {
		return outputList, nil
	}

	output := CredentialsModel{
		Username: input.Username,
	}

	if input.Password != nil {
		output.Password = *input.Password
	}

	return append(outputList, output), nil
}

func flattenImageReferenceModel(input *lab.ImageReference) ([]ImageReferenceModel, error) {
	var outputList []ImageReferenceModel
	if input == nil {
		return outputList, nil
	}

	output := ImageReferenceModel{}

	if input.ExactVersion != nil {
		output.ExactVersion = *input.ExactVersion
	}

	if input.Id != nil {
		output.Id = *input.Id
	}

	if input.Offer != nil {
		output.Offer = *input.Offer
	}

	if input.Publisher != nil {
		output.Publisher = *input.Publisher
	}

	if input.Sku != nil {
		output.Sku = *input.Sku
	}

	if input.Version != nil {
		output.Version = *input.Version
	}

	return append(outputList, output), nil
}

func flattenSkuModel(input *lab.Sku) ([]SkuModel, error) {
	var outputList []SkuModel
	if input == nil {
		return outputList, nil
	}

	output := SkuModel{
		Name: input.Name,
	}

	if input.Capacity != nil {
		output.Capacity = *input.Capacity
	}

	if input.Family != nil {
		output.Family = *input.Family
	}

	if input.Size != nil {
		output.Size = *input.Size
	}

	if input.Tier != nil {
		output.Tier = *input.Tier
	}

	return append(outputList, output), nil
}
