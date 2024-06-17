// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	sqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryIntegrationRuntimeAzureSsis() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeAzureSsisCreateUpdate,
		Read:   resourceDataFactoryIntegrationRuntimeAzureSsisRead,
		Update: resourceDataFactoryIntegrationRuntimeAzureSsisCreateUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeAzureSsisDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationRuntimeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
					`Invalid name for Managed Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(datafactory.IntegrationRuntimeEditionStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.IntegrationRuntimeEditionStandard),
					string(datafactory.IntegrationRuntimeEditionEnterprise),
				}, false),
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(datafactory.IntegrationRuntimeLicenseTypeLicenseIncluded),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.IntegrationRuntimeLicenseTypeLicenseIncluded),
					string(datafactory.IntegrationRuntimeLicenseTypeBasePrice),
				}, false),
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
							ValidateFunc:  sqlValidate.ValidateMsSqlElasticPoolName,
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
		},
	}
}

func resourceDataFactoryIntegrationRuntimeAzureSsisCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	subscriptionId := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewIntegrationRuntimeID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_azure_ssis", id.ID())
		}
	}

	description := d.Get("description").(string)
	managedIntegrationRuntime := datafactory.ManagedIntegrationRuntime{
		Description: &description,
		Type:        datafactory.TypeBasicIntegrationRuntimeTypeManaged,
		ManagedIntegrationRuntimeTypeProperties: &datafactory.ManagedIntegrationRuntimeTypeProperties{
			ComputeProperties:      expandDataFactoryIntegrationRuntimeAzureSsisComputeProperties(d),
			SsisProperties:         expandDataFactoryIntegrationRuntimeAzureSsisProperties(d),
			CustomerVirtualNetwork: expandDataFactoryIntegrationRuntimeCustomerVirtualNetwork(d.Get("express_vnet_integration").([]interface{})),
		},
	}

	basicIntegrationRuntime, _ := managedIntegrationRuntime.AsBasicIntegrationRuntime()

	integrationRuntime := datafactory.IntegrationRuntimeResource{
		Name:       &id.Name,
		Properties: basicIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, integrationRuntime, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryIntegrationRuntimeAzureSsisRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeAzureSsisRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	managedIntegrationRuntime, convertSuccess := resp.Properties.AsManagedIntegrationRuntime()
	if !convertSuccess {
		return fmt.Errorf("converting integration runtime to Azure-SSIS %s", *id)
	}

	if managedIntegrationRuntime.Description != nil {
		d.Set("description", managedIntegrationRuntime.Description)
	}

	if computeProps := managedIntegrationRuntime.ComputeProperties; computeProps != nil {
		if location := computeProps.Location; location != nil {
			d.Set("location", location)
		}

		if nodeSize := computeProps.NodeSize; nodeSize != nil {
			d.Set("node_size", nodeSize)
		}

		if numberOfNodes := computeProps.NumberOfNodes; numberOfNodes != nil {
			d.Set("number_of_nodes", numberOfNodes)
		}

		if maxParallelExecutionsPerNode := computeProps.MaxParallelExecutionsPerNode; maxParallelExecutionsPerNode != nil {
			d.Set("max_parallel_executions_per_node", maxParallelExecutionsPerNode)
		}

		if err := d.Set("vnet_integration", flattenDataFactoryIntegrationRuntimeAzureSsisVnetIntegration(computeProps.VNetProperties)); err != nil {
			return fmt.Errorf("setting `vnet_integration`: %+v", err)
		}

		if err := d.Set("copy_compute_scale", flattenDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(computeProps.CopyComputeScaleProperties)); err != nil {
			return fmt.Errorf("setting `copy_compute_scale`: %+v", err)
		}

		if err := d.Set("pipeline_external_compute_scale", flattenDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScaleProperties(computeProps.PipelineExternalComputeScaleProperties)); err != nil {
			return fmt.Errorf("setting `pipeline_external_compute_scale`: %+v", err)
		}
	}

	if ssisProps := managedIntegrationRuntime.SsisProperties; ssisProps != nil {
		d.Set("edition", string(ssisProps.Edition))
		d.Set("license_type", string(ssisProps.LicenseType))

		if err := d.Set("catalog_info", flattenDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(ssisProps.CatalogInfo, d)); err != nil {
			return fmt.Errorf("setting `catalog_info`: %+v", err)
		}

		if err := d.Set("credential_name", flattenDataFactoryIntegrationRuntimeUserAssignedCredential(ssisProps.Credential)); err != nil {
			return fmt.Errorf("setting `credential_name`: %+v", err)
		}

		if err := d.Set("custom_setup_script", flattenDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(ssisProps.CustomSetupScriptProperties, d)); err != nil {
			return fmt.Errorf("setting `custom_setup_script`: %+v", err)
		}

		if err := d.Set("express_custom_setup", flattenDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(ssisProps.ExpressCustomSetupProperties, d)); err != nil {
			return fmt.Errorf("setting `express_custom_setup`: %+v", err)
		}

		if err := d.Set("package_store", flattenDataFactoryIntegrationRuntimeAzureSsisPackageStore(ssisProps.PackageStores)); err != nil {
			return fmt.Errorf("setting `package_store`: %+v", err)
		}

		if err := d.Set("proxy", flattenDataFactoryIntegrationRuntimeAzureSsisProxy(ssisProps.DataProxyProperties)); err != nil {
			return fmt.Errorf("setting `proxy`: %+v", err)
		}
	}

	if err := d.Set("express_vnet_integration", flattenDataFactoryIntegrationRuntimeCustomerVnetIntegration(managedIntegrationRuntime.CustomerVirtualNetwork)); err != nil {
		return fmt.Errorf("setting `express_vnet_integration`: %+v", err)
	}

	return nil
}

func resourceDataFactoryIntegrationRuntimeAzureSsisDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandDataFactoryIntegrationRuntimeAzureSsisComputeProperties(d *pluginsdk.ResourceData) *datafactory.IntegrationRuntimeComputeProperties {
	location := azure.NormalizeLocation(d.Get("location").(string))
	computeProperties := datafactory.IntegrationRuntimeComputeProperties{
		Location:                     &location,
		NodeSize:                     utils.String(d.Get("node_size").(string)),
		NumberOfNodes:                utils.Int32(int32(d.Get("number_of_nodes").(int))),
		MaxParallelExecutionsPerNode: utils.Int32(int32(d.Get("max_parallel_executions_per_node").(int))),
	}

	if vnetIntegrations, ok := d.GetOk("vnet_integration"); ok && len(vnetIntegrations.([]interface{})) > 0 {
		vnetProps := vnetIntegrations.([]interface{})[0].(map[string]interface{})
		if vnetId := vnetProps["vnet_id"].(string); len(vnetId) > 0 {
			computeProperties.VNetProperties = &datafactory.IntegrationRuntimeVNetProperties{
				VNetID: utils.String(vnetId),
				Subnet: utils.String(vnetProps["subnet_name"].(string)),
			}
		}
		if subnetId := vnetProps["subnet_id"].(string); len(subnetId) > 0 {
			computeProperties.VNetProperties = &datafactory.IntegrationRuntimeVNetProperties{
				SubnetID: utils.String(subnetId),
			}
		}

		if publicIPs := vnetProps["public_ips"].([]interface{}); len(publicIPs) > 0 {
			computeProperties.VNetProperties.PublicIPs = utils.ExpandStringSlice(publicIPs)
		}
	}

	if copyComputeScales, ok := d.GetOk("copy_compute_scale"); ok && len(copyComputeScales.([]interface{})) > 0 {
		copyComputeScale := copyComputeScales.([]interface{})[0].(map[string]interface{})
		if v := copyComputeScale["data_integration_unit"].(int); v != 0 {
			if computeProperties.CopyComputeScaleProperties == nil {
				computeProperties.CopyComputeScaleProperties = &datafactory.CopyComputeScaleProperties{}
			}
			computeProperties.CopyComputeScaleProperties.DataIntegrationUnit = pointer.To(int32(copyComputeScale["data_integration_unit"].(int)))
		}
		if v := copyComputeScale["time_to_live"].(int); v != 0 {
			if computeProperties.CopyComputeScaleProperties == nil {
				computeProperties.CopyComputeScaleProperties = &datafactory.CopyComputeScaleProperties{}
			}
			computeProperties.CopyComputeScaleProperties.TimeToLive = pointer.To(int32(copyComputeScale["time_to_live"].(int)))
		}
	}

	if pipelineExternalComputeScales, ok := d.GetOk("pipeline_external_compute_scale"); ok && len(pipelineExternalComputeScales.([]interface{})) > 0 {
		pipelineExternalComputeScale := pipelineExternalComputeScales.([]interface{})[0].(map[string]interface{})
		if v := pipelineExternalComputeScale["number_of_external_nodes"].(int); v != 0 {
			if computeProperties.PipelineExternalComputeScaleProperties == nil {
				computeProperties.PipelineExternalComputeScaleProperties = &datafactory.PipelineExternalComputeScaleProperties{}
			}
			computeProperties.PipelineExternalComputeScaleProperties.NumberOfExternalNodes = pointer.To(int32(pipelineExternalComputeScale["number_of_external_nodes"].(int)))
		}
		if v := pipelineExternalComputeScale["number_of_pipeline_nodes"].(int); v != 0 {
			if computeProperties.PipelineExternalComputeScaleProperties == nil {
				computeProperties.PipelineExternalComputeScaleProperties = &datafactory.PipelineExternalComputeScaleProperties{}
			}
			computeProperties.PipelineExternalComputeScaleProperties.NumberOfPipelineNodes = pointer.To(int32(pipelineExternalComputeScale["number_of_pipeline_nodes"].(int)))
		}
		if v := pipelineExternalComputeScale["time_to_live"].(int); v != 0 {
			if computeProperties.PipelineExternalComputeScaleProperties == nil {
				computeProperties.PipelineExternalComputeScaleProperties = &datafactory.PipelineExternalComputeScaleProperties{}
			}
			computeProperties.PipelineExternalComputeScaleProperties.TimeToLive = pointer.To(int32(pipelineExternalComputeScale["time_to_live"].(int)))
		}
	}

	return &computeProperties
}

func expandDataFactoryIntegrationRuntimeAzureSsisProperties(d *pluginsdk.ResourceData) *datafactory.IntegrationRuntimeSsisProperties {
	ssisProperties := &datafactory.IntegrationRuntimeSsisProperties{
		LicenseType:                  datafactory.IntegrationRuntimeLicenseType(d.Get("license_type").(string)),
		DataProxyProperties:          expandDataFactoryIntegrationRuntimeAzureSsisProxy(d.Get("proxy").([]interface{})),
		Edition:                      datafactory.IntegrationRuntimeEdition(d.Get("edition").(string)),
		ExpressCustomSetupProperties: expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(d.Get("express_custom_setup").([]interface{})),
		PackageStores:                expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(d.Get("package_store").([]interface{})),
	}

	if credentialName := d.Get("credential_name"); credentialName.(string) != "" {
		ssisProperties.Credential = &datafactory.CredentialReference{
			ReferenceName: utils.String(credentialName.(string)),
			Type:          utils.String("CredentialReference"),
		}
	}

	if catalogInfos, ok := d.GetOk("catalog_info"); ok && len(catalogInfos.([]interface{})) > 0 {
		catalogInfo := catalogInfos.([]interface{})[0].(map[string]interface{})

		// the property `elastic_pool_name` and `pricing_tier` share the same prop `CatalogPricingTier` in request and response.
		var pricingTier datafactory.IntegrationRuntimeSsisCatalogPricingTier
		if elasticPoolName := catalogInfo["elastic_pool_name"]; elasticPoolName != nil && elasticPoolName.(string) != "" {
			pricingTier = datafactory.IntegrationRuntimeSsisCatalogPricingTier(formatDataFactoryIntegrationRuntimeElasticPool(elasticPoolName.(string)))
		} else {
			pricingTier = datafactory.IntegrationRuntimeSsisCatalogPricingTier(catalogInfo["pricing_tier"].(string))
		}

		ssisProperties.CatalogInfo = &datafactory.IntegrationRuntimeSsisCatalogInfo{
			CatalogServerEndpoint: utils.String(catalogInfo["server_endpoint"].(string)),
			CatalogPricingTier:    pricingTier,
		}

		if adminUserName := catalogInfo["administrator_login"]; adminUserName.(string) != "" {
			ssisProperties.CatalogInfo.CatalogAdminUserName = utils.String(adminUserName.(string))
		}

		if adminPassword := catalogInfo["administrator_password"]; adminPassword.(string) != "" {
			ssisProperties.CatalogInfo.CatalogAdminPassword = &datafactory.SecureString{
				Value: utils.String(adminPassword.(string)),
				Type:  datafactory.TypeSecureString,
			}
		}

		if dualStandbyPairName := catalogInfo["dual_standby_pair_name"].(string); dualStandbyPairName != "" {
			ssisProperties.CatalogInfo.DualStandbyPairName = utils.String(dualStandbyPairName)
		}
	}

	if customSetupScripts, ok := d.GetOk("custom_setup_script"); ok && len(customSetupScripts.([]interface{})) > 0 {
		customSetupScript := customSetupScripts.([]interface{})[0].(map[string]interface{})

		sasToken := &datafactory.SecureString{
			Value: utils.String(customSetupScript["sas_token"].(string)),
			Type:  datafactory.TypeSecureString,
		}

		ssisProperties.CustomSetupScriptProperties = &datafactory.IntegrationRuntimeCustomSetupScriptProperties{
			BlobContainerURI: utils.String(customSetupScript["blob_container_uri"].(string)),
			SasToken:         sasToken,
		}
	}

	return ssisProperties
}

func expandDataFactoryIntegrationRuntimeAzureSsisProxy(input []interface{}) *datafactory.IntegrationRuntimeDataProxyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})

	result := &datafactory.IntegrationRuntimeDataProxyProperties{
		ConnectVia: &datafactory.EntityReference{
			Type:          datafactory.IntegrationRuntimeEntityReferenceTypeIntegrationRuntimeReference,
			ReferenceName: utils.String(raw["self_hosted_integration_runtime_name"].(string)),
		},
		StagingLinkedService: &datafactory.EntityReference{
			Type:          datafactory.IntegrationRuntimeEntityReferenceTypeLinkedServiceReference,
			ReferenceName: utils.String(raw["staging_storage_linked_service_name"].(string)),
		},
	}
	if path := raw["path"].(string); len(path) > 0 {
		result.Path = utils.String(path)
	}
	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(input []interface{}) *[]datafactory.BasicCustomSetupBase {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})

	result := make([]datafactory.BasicCustomSetupBase, 0)
	if env := raw["environment"].(map[string]interface{}); len(env) > 0 {
		for k, v := range env {
			result = append(result, &datafactory.EnvironmentVariableSetup{
				Type: datafactory.TypeBasicCustomSetupBaseTypeEnvironmentVariableSetup,
				EnvironmentVariableSetupTypeProperties: &datafactory.EnvironmentVariableSetupTypeProperties{
					VariableName:  utils.String(k),
					VariableValue: utils.String(v.(string)),
				},
			})
		}
	}
	if powershellVersion := raw["powershell_version"].(string); powershellVersion != "" {
		result = append(result, &datafactory.AzPowerShellSetup{
			Type: datafactory.TypeBasicCustomSetupBaseTypeAzPowerShellSetup,
			AzPowerShellSetupTypeProperties: &datafactory.AzPowerShellSetupTypeProperties{
				Version: utils.String(powershellVersion),
			},
		})
	}
	if components := raw["component"].([]interface{}); len(components) > 0 {
		for _, item := range components {
			raw := item.(map[string]interface{})

			var license datafactory.BasicSecretBase
			if v := raw["license"].(string); v != "" {
				license = &datafactory.SecureString{
					Type:  datafactory.TypeSecureString,
					Value: utils.String(v),
				}
			} else {
				license = expandDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(raw["key_vault_license"].([]interface{}))
			}

			result = append(result, &datafactory.ComponentSetup{
				Type: datafactory.TypeBasicCustomSetupBaseTypeComponentSetup,
				LicensedComponentSetupTypeProperties: &datafactory.LicensedComponentSetupTypeProperties{
					ComponentName: utils.String(raw["name"].(string)),
					LicenseKey:    license,
				},
			})
		}
	}
	if cmdKeys := raw["command_key"].([]interface{}); len(cmdKeys) > 0 {
		for _, item := range cmdKeys {
			raw := item.(map[string]interface{})

			var password datafactory.BasicSecretBase
			if v := raw["password"].(string); v != "" {
				password = &datafactory.SecureString{
					Type:  datafactory.TypeSecureString,
					Value: utils.String(v),
				}
			} else {
				password = expandDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(raw["key_vault_password"].([]interface{}))
			}

			result = append(result, &datafactory.CmdkeySetup{
				Type: datafactory.TypeBasicCustomSetupBaseTypeCmdkeySetup,
				CmdkeySetupTypeProperties: &datafactory.CmdkeySetupTypeProperties{
					TargetName: utils.String(raw["target_name"].(string)),
					UserName:   utils.String(raw["user_name"].(string)),
					Password:   password,
				},
			})
		}
	}

	return &result
}

func expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(input []interface{}) *[]datafactory.PackageStore {
	if len(input) == 0 {
		return nil
	}

	result := make([]datafactory.PackageStore, 0)
	for _, item := range input {
		raw := item.(map[string]interface{})
		result = append(result, datafactory.PackageStore{
			Name: utils.String(raw["name"].(string)),
			PackageStoreLinkedService: &datafactory.EntityReference{
				Type:          datafactory.IntegrationRuntimeEntityReferenceTypeLinkedServiceReference,
				ReferenceName: utils.String(raw["linked_service_name"].(string)),
			},
		})
	}
	return &result
}

func expandDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(input []interface{}) *datafactory.AzureKeyVaultSecretReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	reference := &datafactory.AzureKeyVaultSecretReference{
		SecretName: raw["secret_name"].(string),
		Store: &datafactory.LinkedServiceReference{
			Type:          utils.String("LinkedServiceReference"),
			ReferenceName: utils.String(raw["linked_service_name"].(string)),
		},
		Type: datafactory.TypeAzureKeyVaultSecret,
	}
	if v := raw["secret_version"].(string); v != "" {
		reference.SecretVersion = v
	}
	if v := raw["parameters"].(map[string]interface{}); len(v) > 0 {
		reference.Store.Parameters = v
	}
	return reference
}

func expandDataFactoryIntegrationRuntimeCustomerVirtualNetwork(input []interface{}) *datafactory.IntegrationRuntimeCustomerVirtualNetwork {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})
	return &datafactory.IntegrationRuntimeCustomerVirtualNetwork{
		SubnetID: utils.String(raw["subnet_id"].(string)),
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisVnetIntegration(vnetProperties *datafactory.IntegrationRuntimeVNetProperties) []interface{} {
	if vnetProperties == nil {
		return []interface{}{}
	}

	var vnetId, subnetName, subnetId string
	if vnetProperties.VNetID != nil {
		vnetId = *vnetProperties.VNetID
	}
	if vnetProperties.Subnet != nil {
		subnetName = *vnetProperties.Subnet
	}
	if vnetProperties.SubnetID != nil {
		subnetId = *vnetProperties.SubnetID
	}

	return []interface{}{
		map[string]interface{}{
			"vnet_id":     vnetId,
			"subnet_id":   subnetId,
			"subnet_name": subnetName,
			"public_ips":  utils.FlattenStringSlice(vnetProperties.PublicIPs),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(ssisProperties *datafactory.IntegrationRuntimeSsisCatalogInfo, d *pluginsdk.ResourceData) []interface{} {
	if ssisProperties == nil {
		return []interface{}{}
	}

	var serverEndpoint, catalogAdminUserName, administratorPassword, dualStandbyPairName string
	if ssisProperties.CatalogServerEndpoint != nil {
		serverEndpoint = *ssisProperties.CatalogServerEndpoint
	}
	if ssisProperties.CatalogAdminUserName != nil {
		catalogAdminUserName = *ssisProperties.CatalogAdminUserName
	}
	if ssisProperties.DualStandbyPairName != nil {
		dualStandbyPairName = *ssisProperties.DualStandbyPairName
	}

	var pricingTier, elasticPoolName string
	elasticPoolName, elasticPoolNameMatched := parseDataFactoryIntegrationRuntimeElasticPool(string(ssisProperties.CatalogPricingTier))
	if !elasticPoolNameMatched {
		pricingTier = string(ssisProperties.CatalogPricingTier)
	}

	// read back
	if adminPassword, ok := d.GetOk("catalog_info.0.administrator_password"); ok {
		administratorPassword = adminPassword.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"server_endpoint":        serverEndpoint,
			"pricing_tier":           pricingTier,
			"elastic_pool_name":      elasticPoolName,
			"administrator_login":    catalogAdminUserName,
			"administrator_password": administratorPassword,
			"dual_standby_pair_name": dualStandbyPairName,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisProxy(input *datafactory.IntegrationRuntimeDataProxyProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var path, selfHostedIntegrationRuntimeName, stagingStorageLinkedServiceName string
	if input.Path != nil {
		path = *input.Path
	}
	if input.ConnectVia != nil && input.ConnectVia.ReferenceName != nil {
		selfHostedIntegrationRuntimeName = *input.ConnectVia.ReferenceName
	}
	if input.StagingLinkedService != nil && input.StagingLinkedService.ReferenceName != nil {
		stagingStorageLinkedServiceName = *input.StagingLinkedService.ReferenceName
	}
	return []interface{}{
		map[string]interface{}{
			"path":                                 path,
			"self_hosted_integration_runtime_name": selfHostedIntegrationRuntimeName,
			"staging_storage_linked_service_name":  stagingStorageLinkedServiceName,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeUserAssignedCredential(credentialProperties *datafactory.CredentialReference) *string {
	if credentialProperties == nil {
		return nil
	}

	return credentialProperties.ReferenceName
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(customSetupScriptProperties *datafactory.IntegrationRuntimeCustomSetupScriptProperties, d *pluginsdk.ResourceData) []interface{} {
	if customSetupScriptProperties == nil {
		return []interface{}{}
	}

	customSetupScript := map[string]string{
		"blob_container_uri": *customSetupScriptProperties.BlobContainerURI,
	}

	if sasToken, ok := d.GetOk("custom_setup_script.0.sas_token"); ok {
		customSetupScript["sas_token"] = sasToken.(string)
	}

	return []interface{}{customSetupScript}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisPackageStore(input *[]datafactory.PackageStore) []interface{} {
	if input == nil {
		return nil
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		var name, linkedServiceName string
		if item.Name != nil {
			name = *item.Name
		}
		if item.PackageStoreLinkedService != nil && item.PackageStoreLinkedService.ReferenceName != nil {
			linkedServiceName = *item.PackageStoreLinkedService.ReferenceName
		}

		result = append(result, map[string]interface{}{
			"name":                name,
			"linked_service_name": linkedServiceName,
		})
	}
	return result
}

func flattenDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(input *[]datafactory.BasicCustomSetupBase, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	// retrieve old state
	oldState := make(map[string]interface{})
	if arr := d.Get("express_custom_setup").([]interface{}); len(arr) > 0 {
		oldState = arr[0].(map[string]interface{})
	}
	oldComponents := make([]interface{}, 0)
	if rawComponent, ok := oldState["component"]; ok {
		if v := rawComponent.([]interface{}); len(v) > 0 {
			oldComponents = v
		}
	}
	oldCmdKey := make([]interface{}, 0)
	if rawCmdKey, ok := oldState["command_key"]; ok {
		if v := rawCmdKey.([]interface{}); len(v) > 0 {
			oldCmdKey = v
		}
	}

	env := make(map[string]interface{})
	powershellVersion := ""
	components := make([]interface{}, 0)
	cmdkeys := make([]interface{}, 0)
	for _, item := range *input {
		switch v := item.(type) {
		case datafactory.AzPowerShellSetup:
			if v.Version != nil {
				powershellVersion = *v.Version
			}
		case datafactory.ComponentSetup:
			var name string
			if v.ComponentName != nil {
				name = *v.ComponentName
			}
			var keyVaultLicense *datafactory.AzureKeyVaultSecretReference
			if v.LicenseKey != nil {
				if reference, ok := v.LicenseKey.AsAzureKeyVaultSecretReference(); ok {
					keyVaultLicense = reference
				}
			}
			components = append(components, map[string]interface{}{
				"name":              name,
				"key_vault_license": flattenDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(keyVaultLicense),
				"license": readBackSensitiveValue(oldComponents, "license", map[string]string{
					"name": name,
				}),
			})
		case datafactory.EnvironmentVariableSetup:
			if v.VariableName != nil && v.VariableValue != nil {
				env[*v.VariableName] = *v.VariableValue
			}
		case datafactory.CmdkeySetup:
			var name, userName string
			if v.TargetName != nil {
				if v, ok := v.TargetName.(string); ok {
					name = v
				}
			}
			if v.UserName != nil {
				if v, ok := v.UserName.(string); ok {
					userName = v
				}
			}
			var keyVaultPassword *datafactory.AzureKeyVaultSecretReference
			if v.Password != nil {
				if reference, ok := v.Password.AsAzureKeyVaultSecretReference(); ok {
					keyVaultPassword = reference
				}
			}
			cmdkeys = append(cmdkeys, map[string]interface{}{
				"target_name": name,
				"user_name":   userName,
				"password": readBackSensitiveValue(oldCmdKey, "password", map[string]string{
					"target_name": name,
					"user_name":   userName,
				}),
				"key_vault_password": flattenDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(keyVaultPassword),
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"environment":        env,
			"powershell_version": powershellVersion,
			"component":          components,
			"command_key":        cmdkeys,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(input *datafactory.AzureKeyVaultSecretReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	var linkedServiceName, secretName, secretVersion string
	var parameters map[string]interface{}
	if input.SecretName != nil {
		if v, ok := input.SecretName.(string); ok {
			secretName = v
		}
	}
	if input.SecretVersion != nil {
		if v, ok := input.SecretVersion.(string); ok {
			secretVersion = v
		}
	}
	if input.Store != nil {
		if input.Store.ReferenceName != nil {
			linkedServiceName = *input.Store.ReferenceName
		}
		if input.Store.Parameters != nil {
			parameters = input.Store.Parameters
		}
	}
	return []interface{}{
		map[string]interface{}{
			"linked_service_name": linkedServiceName,
			"parameters":          parameters,
			"secret_name":         secretName,
			"secret_version":      secretVersion,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeCustomerVnetIntegration(input *datafactory.IntegrationRuntimeCustomerVirtualNetwork) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	subnetId := ""
	if input.SubnetID != nil {
		subnetId = *input.SubnetID
	}
	return []interface{}{
		map[string]interface{}{
			"subnet_id": subnetId,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScaleProperties(input *datafactory.PipelineExternalComputeScaleProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"number_of_external_nodes": pointer.From(input.NumberOfPipelineNodes),
			"number_of_pipeline_nodes": pointer.From(input.NumberOfPipelineNodes),
			"time_to_live":             pointer.From(input.TimeToLive),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(input *datafactory.CopyComputeScaleProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"data_integration_unit": pointer.From(input.DataIntegrationUnit),
			"time_to_live":          pointer.From(input.TimeToLive),
		},
	}
}

func readBackSensitiveValue(input []interface{}, propertyName string, filters map[string]string) string {
	if len(input) == 0 {
		return ""
	}
	for _, item := range input {
		raw := item.(map[string]interface{})
		found := true
		for k, v := range filters {
			if raw[k].(string) != v {
				found = false
				break
			}
		}
		if found {
			return raw[propertyName].(string)
		}
	}
	return ""
}

func formatDataFactoryIntegrationRuntimeElasticPool(input string) string {
	return fmt.Sprintf(`ELASTIC_POOL(name="%s")`, input)
}

func parseDataFactoryIntegrationRuntimeElasticPool(input string) (string, bool) {
	matches := regexp.MustCompile(`^ELASTIC_POOL\(name="(.+)"\)$`).FindStringSubmatch(input)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}
