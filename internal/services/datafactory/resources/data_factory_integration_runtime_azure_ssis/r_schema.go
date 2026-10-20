// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryIntegrationRuntimeAzureSsisSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
				`Invalid name for Managed Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
			),
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"location": commonschema.Location(),

		"node_size": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard_D2_v3",
				"Standard_D4_v3",
				"Standard_D8_v3",
				"Standard_D16_v3",
				"Standard_D32_v3",
				"Standard_D64_v3",
				"Standard_E2_v3",
				"Standard_E4_v3",
				"Standard_E8_v3",
				"Standard_E16_v3",
				"Standard_E32_v3",
				"Standard_E64_v3",
				"Standard_D1_v2",
				"Standard_D2_v2",
				"Standard_D3_v2",
				"Standard_D4_v2",
				"Standard_A4_v2",
				"Standard_A8_v2",
			}, false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"number_of_nodes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 10),
		},

		"max_parallel_executions_per_node": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 16),
		},

		"credential_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"edition": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(integrationruntimes.IntegrationRuntimeEditionStandard),
			ValidateFunc: validation.StringInSlice(integrationruntimes.PossibleValuesForIntegrationRuntimeEdition(), false),
		},

		"copy_compute_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_integration_unit": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ValidateFunc: validation.All(
							validation.IntBetween(4, 256),
							validation.IntDivisibleBy(4),
						),
					},

					"time_to_live": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(5),
					},
				},
			},
		},

		"express_vnet_integration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"license_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(integrationruntimes.IntegrationRuntimeLicenseTypeLicenseIncluded),
			ValidateFunc: validation.StringInSlice(integrationruntimes.PossibleValuesForIntegrationRuntimeLicenseType(), false),
		},

		"vnet_integration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"vnet_integration.0.vnet_id", "vnet_integration.0.subnet_id"},
						ValidateFunc: azure.ValidateResourceID,
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"vnet_integration.0.vnet_id", "vnet_integration.0.subnet_id"},
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"subnet_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						RequiredWith: []string{"vnet_integration.0.vnet_id"},
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"public_ips": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 2,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidatePublicIPAddressID,
						},
					},
				},
			},
		},

		"custom_setup_script": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_container_uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sas_token": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"catalog_info": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server_endpoint": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"administrator_login": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"administrator_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"pricing_tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Basic",
							"S0", "S1", "S2", "S3", "S4", "S6", "S7", "S9", "S12",
							"P1", "P2", "P4", "P6", "P11", "P15",
							"GP_S_Gen5_1", "GP_S_Gen5_2", "GP_S_Gen5_4", "GP_S_Gen5_6", "GP_S_Gen5_8", "GP_S_Gen5_10", "GP_S_Gen5_12", "GP_S_Gen5_14", "GP_S_Gen5_16", "GP_S_Gen5_18", "GP_S_Gen5_20", "GP_S_Gen5_24", "GP_S_Gen5_32", "GP_S_Gen5_40",
							"GP_Gen5_2", "GP_Gen5_4", "GP_Gen5_6", "GP_Gen5_8", "GP_Gen5_10", "GP_Gen5_12", "GP_Gen5_14", "GP_Gen5_16", "GP_Gen5_18", "GP_Gen5_20", "GP_Gen5_24", "GP_Gen5_32", "GP_Gen5_40", "GP_Gen5_80",
							"BC_Gen5_2", "BC_Gen5_4", "BC_Gen5_6", "BC_Gen5_8", "BC_Gen5_10", "BC_Gen5_12", "BC_Gen5_14", "BC_Gen5_16", "BC_Gen5_18", "BC_Gen5_20", "BC_Gen5_24", "BC_Gen5_32", "BC_Gen5_40", "BC_Gen5_80",
							"HS_Gen5_2", "HS_Gen5_4", "HS_Gen5_6", "HS_Gen5_8", "HS_Gen5_10", "HS_Gen5_12", "HS_Gen5_14", "HS_Gen5_16", "HS_Gen5_18", "HS_Gen5_20", "HS_Gen5_24", "HS_Gen5_32", "HS_Gen5_40", "HS_Gen5_80",
						}, false),
						ConflictsWith: []string{"catalog_info.0.elastic_pool_name"},
					},
					"elastic_pool_name": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validate.ValidateMsSqlElasticPoolName,
						ConflictsWith: []string{"catalog_info.0.pricing_tier"},
					},
					"dual_standby_pair_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"express_custom_setup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"environment": {
						Type:         pluginsdk.TypeMap,
						Optional:     true,
						AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"powershell_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"command_key": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"target_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"user_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"key_vault_password": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"linked_service_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"secret_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"parameters": {
												Type:     pluginsdk.TypeMap,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"secret_version": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},

					"component": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"license": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"key_vault_license": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"linked_service_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"secret_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"parameters": {
												Type:     pluginsdk.TypeMap,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"secret_version": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
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

		"package_store": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"pipeline_external_compute_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"number_of_external_nodes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 10),
					},

					"number_of_pipeline_nodes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 10),
					},

					"time_to_live": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(5),
					},
				},
			},
		},

		"proxy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"self_hosted_integration_runtime_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"staging_storage_linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
