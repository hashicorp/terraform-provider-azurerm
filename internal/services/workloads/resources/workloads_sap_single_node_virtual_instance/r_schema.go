// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_single_node_virtual_instance

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SAPVirtualInstanceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"app_location": commonschema.Location(),

		"sap_fqdn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(2, 34),
		},

		"single_server_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_resource_group_name": commonschema.ResourceGroupName(),

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: networkValidate.SubnetID,
					},

					"virtual_machine_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"image": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"offer": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"publisher": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ForceNew: true,
												ValidateFunc: validation.StringInSlice([]string{
													"RedHat",
													"SUSE",
												}, false),
											},

											"sku": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"version": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"os_profile": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"admin_username": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringLenBetween(1, 64),
											},

											"ssh_private_key": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												Sensitive:    true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"ssh_public_key": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"virtual_machine_size": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"database_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPDatabaseType(), false),
					},

					"disk_volume_configuration": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"volume_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										"backup",
										"hana/data",
										"hana/log",
										"hana/shared",
										"os",
										"usr/sap",
									}, false),
								},

								"number_of_disks": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"size_in_gb": {
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
								},

								"sku_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForDiskSkuName(), false),
								},
							},
						},
					},

					"secondary_ip_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"virtual_machine_resource_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_disk": {
									Type:         pluginsdk.TypeSet,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"volume_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ForceNew: true,
												ValidateFunc: validation.StringInSlice([]string{
													"default",
												}, false),
											},

											"names": {
												Type:     pluginsdk.TypeList,
												Required: true,
												ForceNew: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 80),
												},
											},
										},
									},
								},

								"host_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringLenBetween(1, 13),
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
								},

								"network_interface_names": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: networkValidate.NetworkInterfaceName,
									},
								},

								"os_disk_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringLenBetween(1, 80),
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
								},

								"virtual_machine_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: computeValidate.VirtualMachineName,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
								},
							},
						},
					},
				},
			},
		},

		"environment": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPEnvironmentType(), false),
		},

		"sap_product": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPProductType(), false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
		},

		"managed_resources_network_access_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(sapvirtualinstances.ManagedResourcesNetworkAccessTypePublic),
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForManagedResourcesNetworkAccessType(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if v := rd.Get("single_server_configuration.0.disk_volume_configuration"); v != nil {
				diskVolumes := v.(*pluginsdk.Set).List()
				if hasDuplicateVolumeNameForSAPSingleNodeVirtualInstance(diskVolumes) {
					return errors.New("`volume_name` cannot be duplicated")
				}
			}

			return nil
		},
	}
}
