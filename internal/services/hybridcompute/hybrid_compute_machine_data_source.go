package hybridcompute

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-03-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type HybridComputeMachineModel struct {
	Name                       string                    `tfschema:"name"`
	ResourceGroupName          string                    `tfschema:"resource_group_name"`
	AgentConfiguration         []AgentConfigurationModel `tfschema:"agent_configuration"`
	ClientPublicKey            string                    `tfschema:"client_public_key"`
	CloudMetadata              []CloudMetadataModel      `tfschema:"cloud_metadata"`
	DetectedProperties         map[string]string         `tfschema:"detected_properties"`
	Location                   string                    `tfschema:"location"`
	LocationData               []LocationDataModel       `tfschema:"location_data"`
	MssqlDiscovered            bool                      `tfschema:"mssql_discovered"`
	OsProfile                  []OSProfileModel          `tfschema:"os_profile"`
	OsType                     string                    `tfschema:"os_type"`
	ParentClusterResourceId    string                    `tfschema:"parent_cluster_resource_id"`
	PrivateLinkScopeResourceId string                    `tfschema:"private_link_scope_resource_id"`
	ServiceStatuses            []ServiceStatusesModel    `tfschema:"service_status"`
	Tags                       map[string]string         `tfschema:"tags"`
	VmId                       string                    `tfschema:"vm_id"`
	AdFqdn                     string                    `tfschema:"ad_fqdn"`
	AgentVersion               string                    `tfschema:"agent_version"`
	DisplayName                string                    `tfschema:"display_name"`
	DnsFqdn                    string                    `tfschema:"dns_fqdn"`
	DomainName                 string                    `tfschema:"domain_name"`
	ErrorDetails               []ErrorDetailModel        `tfschema:"error_details"`
	LastStatusChange           string                    `tfschema:"last_status_change"`
	MachineFqdn                string                    `tfschema:"machine_fqdn"`
	OsName                     string                    `tfschema:"os_name"`
	OsSku                      string                    `tfschema:"os_sku"`
	OsVersion                  string                    `tfschema:"os_version"`
	Status                     machines.StatusTypes      `tfschema:"status"`
	VmUuid                     string                    `tfschema:"vm_uuid"`
}

type AgentConfigurationModel struct {
	ExtensionsAllowList       []ConfigurationExtensionModel `tfschema:"extensions_allow_list"`
	ExtensionsBlockList       []ConfigurationExtensionModel `tfschema:"extensions_block_list"`
	ExtensionsEnabled         bool                          `tfschema:"extensions_enabled"`
	GuestConfigurationEnabled bool                          `tfschema:"guest_configuration_enabled"`
	IncomingConnectionsPorts  []string                      `tfschema:"incoming_connections_ports"`
	ProxyBypass               []string                      `tfschema:"proxy_bypass"`
	ProxyUrl                  string                        `tfschema:"proxy_url"`
}

type ConfigurationExtensionModel struct {
	Publisher string `tfschema:"publisher"`
	Type      string `tfschema:"type"`
}

type CloudMetadataModel struct {
	Provider string `tfschema:"provider"`
}

type LocationDataModel struct {
	City            string `tfschema:"city"`
	CountryOrRegion string `tfschema:"country_or_region"`
	District        string `tfschema:"district"`
	Name            string `tfschema:"name"`
}

type OSProfileModel struct {
	ComputerName         string                               `tfschema:"computer_name"`
	LinuxConfiguration   []OSProfileLinuxConfigurationModel   `tfschema:"linux_configuration"`
	WindowsConfiguration []OSProfileWindowsConfigurationModel `tfschema:"windows_configuration"`
}

type OSProfileLinuxConfigurationModel struct {
	PatchSettings []PatchSettingsModel `tfschema:"patch_settings"`
}

type PatchSettingsModel struct {
	AssessmentMode machines.AssessmentModeTypes `tfschema:"assessment_mode"`
	PatchMode      machines.PatchModeTypes      `tfschema:"patch_mode"`
}

type OSProfileWindowsConfigurationModel struct {
	PatchSettings []PatchSettingsModel `tfschema:"patch_settings"`
}

type ServiceStatusesModel struct {
	ExtensionService          []ServiceStatusModel `tfschema:"extension_service"`
	GuestConfigurationService []ServiceStatusModel `tfschema:"guest_configuration_service"`
}

type ServiceStatusModel struct {
	StartupType string `tfschema:"startup_type"`
	Status      string `tfschema:"status"`
}

type ErrorDetailModel struct {
	AdditionalInfo []ErrorAdditionalInfoModel `tfschema:"additional_info"`
	Code           string                     `tfschema:"code"`
	Message        string                     `tfschema:"message"`
	Target         string                     `tfschema:"target"`
}

type ErrorAdditionalInfoModel struct {
	Info string `tfschema:"info"`
	Type string `tfschema:"type"`
}

type HybridComputeMachineDataSource struct{}

var _ sdk.DataSource = HybridComputeMachineDataSource{}

func (r HybridComputeMachineDataSource) ResourceType() string {
	return "azurerm_hybrid_compute_machine"
}

func (r HybridComputeMachineDataSource) ModelObject() interface{} {
	return &HybridComputeMachineModel{}
}

func (r HybridComputeMachineDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return machines.ValidateMachineID
}

func (r HybridComputeMachineDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r HybridComputeMachineDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"extensions_allow_list": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"publisher": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"extensions_block_list": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"publisher": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"extensions_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"guest_configuration_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"incoming_connections_ports": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"proxy_bypass": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"proxy_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"ad_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"agent_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"client_public_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cloud_metadata": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"provider": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"detected_properties": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"error_details": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_info": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"info": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"code": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"message": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"target": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedIdentityComputed(),

		"last_status_change": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"location_data": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"city": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"country_or_region": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"district": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"machine_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"mssql_discovered": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"os_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"linux_configuration": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"patch_settings": {
									Type:     pluginsdk.TypeList,
									Computed: true,

									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"assessment_mode": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"patch_mode": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"windows_configuration": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"patch_settings": {
									Type:     pluginsdk.TypeList,
									Computed: true,

									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"assessment_mode": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"patch_mode": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		"os_sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"os_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parent_cluster_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_link_scope_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_status": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"extension_service": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"startup_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"status": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"guest_configuration_service": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"startup_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"status": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vm_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vm_uuid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r HybridComputeMachineDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.MachinesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var hybridComputeMachineModel HybridComputeMachineModel
			if err := metadata.Decode(&hybridComputeMachineModel); err != nil {
				return err
			}

			id := machines.NewMachineID(subscriptionId, hybridComputeMachineModel.ResourceGroupName, hybridComputeMachineModel.Name)

			resp, err := client.Get(ctx, id, machines.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := HybridComputeMachineModel{
				Name:              id.MachineName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue := identity.FlattenSystemAssigned(model.Identity)

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				if properties.AdFqdn != nil {
					state.AdFqdn = *properties.AdFqdn
				}

				agentConfigurationValue, err := flattenAgentConfigurationModel(properties.AgentConfiguration)
				if err != nil {
					return err
				}

				state.AgentConfiguration = agentConfigurationValue

				if properties.AgentVersion != nil {
					state.AgentVersion = *properties.AgentVersion
				}

				if properties.ClientPublicKey != nil {
					state.ClientPublicKey = *properties.ClientPublicKey
				}

				cloudMetadataValue := flattenCloudMetadataModel(properties.CloudMetadata)

				state.CloudMetadata = cloudMetadataValue

				if properties.DetectedProperties != nil {
					state.DetectedProperties = *properties.DetectedProperties
				}

				if properties.DisplayName != nil {
					state.DisplayName = *properties.DisplayName
				}

				if properties.DnsFqdn != nil {
					state.DnsFqdn = *properties.DnsFqdn
				}

				if properties.DomainName != nil {
					state.DomainName = *properties.DomainName
				}

				errorDetailsValue := flattenErrorDetailModel(properties.ErrorDetails)

				state.ErrorDetails = errorDetailsValue

				if properties.LastStatusChange != nil {
					state.LastStatusChange = *properties.LastStatusChange
				}

				locationDataValue := flattenLocationDataModel(properties.LocationData)

				state.LocationData = locationDataValue

				if properties.MachineFqdn != nil {
					state.MachineFqdn = *properties.MachineFqdn
				}

				if properties.MssqlDiscovered != nil {
					state.MssqlDiscovered, err = strconv.ParseBool(*properties.MssqlDiscovered)
					if err != nil {
						return err
					}
				}

				if properties.OsName != nil {
					state.OsName = *properties.OsName
				}

				osProfileValue := flattenOSProfileModel(properties.OsProfile)

				state.OsProfile = osProfileValue

				if properties.OsSku != nil {
					state.OsSku = *properties.OsSku
				}

				if properties.OsType != nil {
					state.OsType = *properties.OsType
				}

				if properties.OsVersion != nil {
					state.OsVersion = *properties.OsVersion
				}

				if properties.ParentClusterResourceId != nil {
					state.ParentClusterResourceId = *properties.ParentClusterResourceId
				}

				if properties.PrivateLinkScopeResourceId != nil {
					state.PrivateLinkScopeResourceId = *properties.PrivateLinkScopeResourceId
				}

				serviceStatusesValue := flattenServiceStatusesModel(properties.ServiceStatuses)

				state.ServiceStatuses = serviceStatusesValue

				if properties.Status != nil {
					state.Status = *properties.Status
				}

				if properties.VMId != nil {
					state.VmId = *properties.VMId
				}

				if properties.VMUuid != nil {
					state.VmUuid = *properties.VMUuid
				}
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenAgentConfigurationModel(input *machines.AgentConfiguration) ([]AgentConfigurationModel, error) {
	var outputList []AgentConfigurationModel
	if input == nil {
		return outputList, nil
	}

	output := AgentConfigurationModel{}

	extensionsAllowListValue := flattenConfigurationExtensionModel(input.ExtensionsAllowList)

	output.ExtensionsAllowList = extensionsAllowListValue

	extensionsBlockListValue := flattenConfigurationExtensionModel(input.ExtensionsBlockList)

	output.ExtensionsBlockList = extensionsBlockListValue

	if input.ExtensionsEnabled != nil {
		parsedBool, err := strconv.ParseBool(*input.ExtensionsEnabled)
		if err != nil {
			return nil, err
		}
		output.ExtensionsEnabled = parsedBool
	}

	if input.GuestConfigurationEnabled != nil {
		parsedBool, err := strconv.ParseBool(*input.GuestConfigurationEnabled)
		if err != nil {
			return nil, err
		}
		output.GuestConfigurationEnabled = parsedBool
	}

	if input.IncomingConnectionsPorts != nil {
		output.IncomingConnectionsPorts = *input.IncomingConnectionsPorts
	}

	if input.ProxyBypass != nil {
		output.ProxyBypass = *input.ProxyBypass
	}

	if input.ProxyUrl != nil {
		output.ProxyUrl = *input.ProxyUrl
	}

	return append(outputList, output), nil
}

func flattenConfigurationExtensionModel(inputList *[]machines.ConfigurationExtension) []ConfigurationExtensionModel {
	var outputList []ConfigurationExtensionModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := ConfigurationExtensionModel{}

		if input.Publisher != nil {
			output.Publisher = *input.Publisher
		}

		if input.Type != nil {
			output.Type = *input.Type
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenCloudMetadataModel(input *machines.CloudMetadata) []CloudMetadataModel {
	var outputList []CloudMetadataModel
	if input == nil {
		return outputList
	}

	output := CloudMetadataModel{}

	if input.Provider != nil {
		output.Provider = *input.Provider
	}

	return append(outputList, output)
}

func flattenErrorDetailModel(inputList *[]machines.ErrorDetail) []ErrorDetailModel {
	var outputList []ErrorDetailModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := ErrorDetailModel{}

		additionalInfoValue := flattenErrorAdditionalInfoModel(input.AdditionalInfo)

		output.AdditionalInfo = additionalInfoValue

		if input.Code != nil {
			output.Code = *input.Code
		}

		if input.Message != nil {
			output.Message = *input.Message
		}

		if input.Target != nil {
			output.Target = *input.Target
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenErrorAdditionalInfoModel(inputList *[]machines.ErrorAdditionalInfo) []ErrorAdditionalInfoModel {
	var outputList []ErrorAdditionalInfoModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := ErrorAdditionalInfoModel{}

		if input.Info != nil && *input.Info != nil {

			infoValue, err := json.Marshal(*input.Info)
			if err != nil {
				return nil
			}

			output.Info = string(infoValue)
		}

		if input.Type != nil {
			output.Type = *input.Type
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenLocationDataModel(input *machines.LocationData) []LocationDataModel {
	var outputList []LocationDataModel
	if input == nil {
		return outputList
	}

	output := LocationDataModel{
		Name: input.Name,
	}

	if input.City != nil {
		output.City = *input.City
	}

	if input.CountryOrRegion != nil {
		output.CountryOrRegion = *input.CountryOrRegion
	}

	if input.District != nil {
		output.District = *input.District
	}

	return append(outputList, output)
}

func flattenOSProfileModel(input *machines.OSProfile) []OSProfileModel {
	var outputList []OSProfileModel
	if input == nil {
		return outputList
	}

	output := OSProfileModel{}

	if input.ComputerName != nil {
		output.ComputerName = *input.ComputerName
	}

	linuxConfigurationValue := flattenOSProfileLinuxConfigurationModel(input.LinuxConfiguration)

	output.LinuxConfiguration = linuxConfigurationValue

	windowsConfigurationValue := flattenOSProfileWindowsConfigurationModel(input.WindowsConfiguration)
	output.WindowsConfiguration = windowsConfigurationValue

	return append(outputList, output)
}

func flattenOSProfileLinuxConfigurationModel(input *machines.OSProfileLinuxConfiguration) []OSProfileLinuxConfigurationModel {
	var outputList []OSProfileLinuxConfigurationModel
	if input == nil {
		return outputList
	}

	output := OSProfileLinuxConfigurationModel{}

	patchSettingsValue := flattenPatchSettingsModel(input.PatchSettings)

	output.PatchSettings = patchSettingsValue

	return append(outputList, output)
}

func flattenPatchSettingsModel(input *machines.PatchSettings) []PatchSettingsModel {
	var outputList []PatchSettingsModel
	if input == nil {
		return outputList
	}

	output := PatchSettingsModel{}

	if input.AssessmentMode != nil {
		output.AssessmentMode = *input.AssessmentMode
	}

	if input.PatchMode != nil {
		output.PatchMode = *input.PatchMode
	}

	return append(outputList, output)
}

func flattenOSProfileWindowsConfigurationModel(input *machines.OSProfileWindowsConfiguration) []OSProfileWindowsConfigurationModel {
	var outputList []OSProfileWindowsConfigurationModel
	if input == nil {
		return outputList
	}

	output := OSProfileWindowsConfigurationModel{}
	patchSettingsValue := flattenPatchSettingsModel(input.PatchSettings)
	output.PatchSettings = patchSettingsValue

	return append(outputList, output)
}

func flattenServiceStatusesModel(input *machines.ServiceStatuses) []ServiceStatusesModel {
	var outputList []ServiceStatusesModel
	if input == nil {
		return outputList
	}

	output := ServiceStatusesModel{}

	extensionServiceValue := flattenServiceStatusModel(input.ExtensionService)
	output.ExtensionService = extensionServiceValue

	guestConfigurationServiceValue := flattenServiceStatusModel(input.GuestConfigurationService)
	output.GuestConfigurationService = guestConfigurationServiceValue

	return append(outputList, output)
}

func flattenServiceStatusModel(input *machines.ServiceStatus) []ServiceStatusModel {
	var outputList []ServiceStatusModel
	if input == nil {
		return outputList
	}

	output := ServiceStatusModel{}

	if input.StartupType != nil {
		output.StartupType = *input.StartupType
	}

	if input.Status != nil {
		output.Status = *input.Status
	}

	return append(outputList, output)
}
