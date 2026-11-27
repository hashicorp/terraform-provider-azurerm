// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package arc_machine

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	hybridcomputeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/hybridcompute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ArcMachineModel struct {
	Name                       string                    `tfschema:"name"`
	ResourceGroupName          string                    `tfschema:"resource_group_name"`
	AgentConfiguration         []AgentConfigurationModel `tfschema:"agent"`
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
	AdFqdn                     string                    `tfschema:"active_directory_fqdn"`
	AgentVersion               string                    `tfschema:"agent_version"`
	DisplayName                string                    `tfschema:"display_name"`
	DnsFqdn                    string                    `tfschema:"dns_fqdn"`
	DomainName                 string                    `tfschema:"domain_name"`
	LastStatusChange           string                    `tfschema:"last_status_change_time"`
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
	ProxyURL                  string                        `tfschema:"proxy_url"`
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
	LinuxConfiguration   []OSProfileLinuxConfigurationModel   `tfschema:"linux"`
	WindowsConfiguration []OSProfileWindowsConfigurationModel `tfschema:"windows"`
}

type OSProfileLinuxConfigurationModel struct {
	PatchSettings []PatchSettingsModel `tfschema:"patch"`
}

type PatchSettingsModel struct {
	AssessmentMode machines.AssessmentModeTypes `tfschema:"assessment_mode"`
	PatchMode      machines.PatchModeTypes      `tfschema:"patch_mode"`
}

type OSProfileWindowsConfigurationModel struct {
	PatchSettings []PatchSettingsModel `tfschema:"patch"`
}

type ServiceStatusesModel struct {
	ExtensionService          []ServiceStatusModel `tfschema:"extension_service"`
	GuestConfigurationService []ServiceStatusModel `tfschema:"guest_configuration_service"`
}

type ServiceStatusModel struct {
	StartupType string `tfschema:"startup_type"`
	Status      string `tfschema:"status"`
}

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (a DataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Azure Arc machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (a DataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"agent": {
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
									Description: "Publisher of the extension.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"type": {
									Description: "Type of the extension.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
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
									Description: "Publisher of the extension.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"type": {
									Description: "Type of the extension.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},

					"extensions_enabled": {
						Description: "Whether the extension service is enabled or disabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"guest_configuration_enabled": {
						Description: "Specified whether the guest configuration service is enabled or disabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"incoming_connections_ports": {
						Description: "The list of ports that the agent will be able to listen on.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"proxy_bypass": {
						Description: "List of service names which should not use the specified proxy server.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"proxy_url": {
						Description: "The URL of the proxy to be used.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"active_directory_fqdn": {
			Description: "The AD fully qualified display name.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"agent_version": {
			Description: "The Azure Arc machine agent full version.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"client_public_key": {
			Description: "Public Key that the client provides to be used during initial resource onboarding.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"cloud_metadata": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"provider": {
						Description: "The cloud provider. For example `Azure`, `AWS` and `GCP`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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
			Description: "The Azure Arc machine display name.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"dns_fqdn": {
			Description: "The DNS fully qualified display name.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"domain_name": {
			Description: "The Windows domain name.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedIdentityComputed(),

		"last_status_change_time": {
			Description: "The time of the last status change.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"location": commonschema.LocationComputed(),

		"location_data": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"city": {
						Description: "The city or locality where the resource is located.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"country_or_region": {
						Description: "The country or region where the resource is located.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"district": {
						Description: "The district, state, or province where the resource is located.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"name": {
						Description: "A canonical name for the geographic or physical location.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"machine_fqdn": {
			Description: "The Azure Arc machine fully qualified display name.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"mssql_discovered": {
			Description: "Whether any MS SQL instance is discovered on the machine.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"os_name": {
			Description: "The Operating System running on the Azure Arc machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Description: "The host OS name of the Azure Arc machine.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"linux": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"patch": {
									Type:     pluginsdk.TypeList,
									Computed: true,

									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"assessment_mode": {
												Description: "The assessment mode.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},

											"patch_mode": {
												Description: "The patch mode.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},

					"windows": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"patch": {
									Type:     pluginsdk.TypeList,
									Computed: true,

									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"assessment_mode": {
												Description: "The assessment mode.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},

											"patch_mode": {
												Description: "The patch mode.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
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
			Description: "The Operating System product SKU.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"os_type": {
			Description: "The type of Operating System.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"os_version": {
			Description: "The version of Operating System running on the Azure Arc machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"parent_cluster_resource_id": {
			Description: "The resource ID of the parent cluster (Azure HCI) this machine is assigned to, if any.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_link_scope_resource_id": {
			Description: "The resource ID of the parent cluster (Azure HCI) this machine is assigned to, if any.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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
									Description: "The behaviour of the service when the Arc-enabled machine starts up.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"status": {
									Description: "The current status of the service.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
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
									Description: "The behaviour of the service when the Arc-enabled machine starts up.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"status": {
									Description: "The current status of the service.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"status": {
			Description: "The status of the Azure Arc machine agent.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"vm_id": {
			Description: "The Azure Arc machine unique ID.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"vm_uuid": {
			Description: "The Arc Machine's unique SMBIOS ID.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (a DataSource) ModelObject() any {
	return &DataSource{}
}

func (a DataSource) ResourceType() string {
	return "azurerm_arc_machine"
}

func (a DataSource) Description() string {
	return "Use this data source to access information about an existing Azure Arc machine."
}

func (a DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := hybridcomputeClient.From(metadata.Client).MachinesClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var arcMachineModel ArcMachineModel
			if err := metadata.Decode(&arcMachineModel); err != nil {
				return err
			}

			id := machines.NewMachineID(subscriptionId, arcMachineModel.ResourceGroupName, arcMachineModel.Name)

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

			state := ArcMachineModel{
				Name:              id.MachineName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if err := metadata.ResourceData.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
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

				state.CloudMetadata = flattenCloudMetadataModel(properties.CloudMetadata)

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

				if properties.ErrorDetails != nil {
					if len(*properties.ErrorDetails) > 0 {
						return fmt.Errorf("retrieving %s: error details: %+v", id, *properties.ErrorDetails)
					}
				}

				if properties.LastStatusChange != nil {
					state.LastStatusChange = *properties.LastStatusChange
				}

				state.LocationData = flattenLocationDataModel(properties.LocationData)

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

				state.OsProfile = flattenOSProfileModel(properties.OsProfile)

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

				state.ServiceStatuses = flattenServiceStatusesModel(properties.ServiceStatuses)

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
		return []AgentConfigurationModel{}, nil
	}

	output := AgentConfigurationModel{}

	output.ExtensionsAllowList = flattenConfigurationExtensionModel(input.ExtensionsAllowList)

	output.ExtensionsBlockList = flattenConfigurationExtensionModel(input.ExtensionsBlockList)

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

	if input.ProxyURL != nil {
		output.ProxyURL = *input.ProxyURL
	}

	return append(outputList, output), nil
}

func flattenConfigurationExtensionModel(inputList *[]machines.ConfigurationExtension) []ConfigurationExtensionModel {
	if inputList == nil {
		return []ConfigurationExtensionModel{}
	}

	outputList := make([]ConfigurationExtensionModel, 0, len(*inputList))
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
		return []CloudMetadataModel{}
	}

	output := CloudMetadataModel{}

	if input.Provider != nil {
		output.Provider = *input.Provider
	}

	return append(outputList, output)
}

func flattenLocationDataModel(input *machines.LocationData) []LocationDataModel {
	var outputList []LocationDataModel
	if input == nil {
		return []LocationDataModel{}
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
		return []OSProfileModel{}
	}

	output := OSProfileModel{}

	if input.ComputerName != nil {
		output.ComputerName = *input.ComputerName
	}

	output.LinuxConfiguration = flattenOSProfileLinuxConfigurationModel(input.LinuxConfiguration)

	output.WindowsConfiguration = flattenOSProfileWindowsConfigurationModel(input.WindowsConfiguration)

	return append(outputList, output)
}

func flattenOSProfileLinuxConfigurationModel(input *machines.OSProfileLinuxConfiguration) []OSProfileLinuxConfigurationModel {
	var outputList []OSProfileLinuxConfigurationModel
	if input == nil {
		return []OSProfileLinuxConfigurationModel{}
	}

	output := OSProfileLinuxConfigurationModel{}

	output.PatchSettings = flattenPatchSettingsModel(input.PatchSettings)

	return append(outputList, output)
}

func flattenPatchSettingsModel(input *machines.PatchSettings) []PatchSettingsModel {
	var outputList []PatchSettingsModel
	if input == nil {
		return []PatchSettingsModel{}
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
		return []OSProfileWindowsConfigurationModel{}
	}

	output := OSProfileWindowsConfigurationModel{}
	output.PatchSettings = flattenPatchSettingsModel(input.PatchSettings)

	return append(outputList, output)
}

func flattenServiceStatusesModel(input *machines.ServiceStatuses) []ServiceStatusesModel {
	var outputList []ServiceStatusesModel
	if input == nil {
		return []ServiceStatusesModel{}
	}

	output := ServiceStatusesModel{}

	output.ExtensionService = flattenServiceStatusModel(input.ExtensionService)

	output.GuestConfigurationService = flattenServiceStatusModel(input.GuestConfigurationService)

	return append(outputList, output)
}

func flattenServiceStatusModel(input *machines.ServiceStatus) []ServiceStatusModel {
	var outputList []ServiceStatusModel
	if input == nil {
		return []ServiceStatusModel{}
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
