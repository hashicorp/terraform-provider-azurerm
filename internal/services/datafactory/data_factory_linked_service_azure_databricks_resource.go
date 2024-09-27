// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-04-01-preview/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryLinkedServiceAzureDatabricks() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceDatabricksCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceDatabricksRead,
		Update: resourceDataFactoryLinkedServiceDatabricksCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceDatabricksDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, importDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeAzureDatabricks)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			// Authentication types
			"msi_work_space_resource_id": {
				// TODO: rename this to `msi_workspace_id` in v4.0
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
				ExactlyOneOf: []string{"access_token", "msi_work_space_resource_id", "key_vault_password"},
			},

			"access_token": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"access_token", "msi_work_space_resource_id", "key_vault_password"},
			},

			"key_vault_password": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// TODO use LinkedServiceDataSetName and NestedItemName validate here and in other linked service resources
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
					},
				},
				ExactlyOneOf: []string{"access_token", "msi_work_space_resource_id", "key_vault_password"},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"adb_domain": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Cluster types [existing cluster, new cluster, interactive pools]
			"existing_cluster_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
			},

			"new_cluster_config": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"custom_tags": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"min_number_of_workers": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      "1",
							ValidateFunc: validation.IntBetween(1, 25000),
						},
						"max_number_of_workers": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 25000),
						},
						"cluster_version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"spark_config": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"spark_environment_variables": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"log_destination": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"init_scripts": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"driver_node_type": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
			},

			"instance_pool": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"min_number_of_workers": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 25000),
						},
						"max_number_of_workers": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 25000),
						},
						"instance_pool_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"cluster_version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
			},

			"integration_runtime_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceDataFactoryLinkedServiceDatabricksCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory Databricks %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_databricks", id.ID())
		}
	}

	var databricksProperties *datafactory.AzureDatabricksLinkedServiceTypeProperties

	// Check if the MSI authentication block is set
	msiAuth := d.Get("msi_work_space_resource_id").(string)
	accessTokenAuth := d.Get("access_token").(string)
	accessTokenKeyVaultAuth := d.Get("key_vault_password").([]interface{})

	// Set the properties based on the authentication type that was provided
	if msiAuth != "" {
		databricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			Authentication:      "MSI",
			WorkspaceResourceID: msiAuth,
		}
	}
	if accessTokenAuth != "" {
		// Assign the access token in the properties block
		databricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			AccessToken: &datafactory.SecureString{
				Value: utils.String(accessTokenAuth),
				Type:  datafactory.TypeSecureString,
			},
		}
	}

	if len(accessTokenKeyVaultAuth) > 0 && accessTokenKeyVaultAuth[0] != nil {
		databricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			AccessToken: expandAzureKeyVaultSecretReference(accessTokenKeyVaultAuth),
		}
	}

	// Set the other type properties
	databricksProperties.Domain = d.Get("adb_domain").(string)

	if v, ok := d.GetOk("existing_cluster_id"); ok {
		databricksProperties.ExistingClusterID = v.(string)
	}

	if v, ok := d.GetOk("instance_pool"); ok && v.([]interface{})[0] != nil {
		instancePoolMap := v.([]interface{})[0].(map[string]interface{})

		if data := instancePoolMap["instance_pool_id"]; data != nil {
			databricksProperties.InstancePoolID = data
		}

		if data := instancePoolMap["cluster_version"]; data != nil {
			databricksProperties.NewClusterVersion = data
		}

		if minWorkersProperty := instancePoolMap["min_number_of_workers"]; minWorkersProperty != nil {
			maxWorkersProperty := instancePoolMap["max_number_of_workers"]
			if numOfWorkersProperty, err := buildNumberOfWorkersProperties(minWorkersProperty, maxWorkersProperty); err == nil {
				databricksProperties.NewClusterNumOfWorker = numOfWorkersProperty
			} else {
				return fmt.Errorf("expanding `instance_pool`: +%v", err)
			}
		}
	}

	if v, ok := d.GetOk("new_cluster_config"); ok && v.([]interface{})[0] != nil {
		newClusterMap := v.([]interface{})[0].(map[string]interface{})

		if data := newClusterMap["cluster_version"]; data != nil {
			databricksProperties.NewClusterVersion = data
		}

		if minWorkersProperty := newClusterMap["min_number_of_workers"]; minWorkersProperty != nil {
			maxWorkersProperty := newClusterMap["max_number_of_workers"]
			if numOfWorkersProperty, err := buildNumberOfWorkersProperties(minWorkersProperty, maxWorkersProperty); err == nil {
				databricksProperties.NewClusterNumOfWorker = numOfWorkersProperty
			} else {
				return fmt.Errorf("expanding `new_cluster_config`: +%v", err)
			}
		}

		if data := newClusterMap["node_type"]; data != nil {
			databricksProperties.NewClusterNodeType = data
		}

		if data := newClusterMap["driver_node_type"]; data != nil {
			databricksProperties.NewClusterDriverNodeType = data
		}

		if data := newClusterMap["log_destination"]; data != nil {
			databricksProperties.NewClusterLogDestination = data
		}

		if newClusterMap["spark_config"] != nil {
			if sparkConfig := newClusterMap["spark_config"].(map[string]interface{}); len(sparkConfig) > 0 {
				databricksProperties.NewClusterSparkConf = sparkConfig
			}
		}

		if newClusterMap["spark_environment_variables"] != nil {
			if sparkEnvVars := newClusterMap["spark_environment_variables"].(map[string]interface{}); len(sparkEnvVars) > 0 {
				databricksProperties.NewClusterSparkEnvVars = sparkEnvVars
			}
		}

		if newClusterMap["custom_tags"] != nil {
			if customTags := newClusterMap["custom_tags"].(map[string]interface{}); len(customTags) > 0 {
				databricksProperties.NewClusterCustomTags = customTags
			}
		}

		initScripts := newClusterMap["init_scripts"]
		databricksProperties.NewClusterInitScripts = &initScripts
	}

	databricksLinkedService := &datafactory.AzureDatabricksLinkedService{
		Description: utils.String(d.Get("description").(string)),
		AzureDatabricksLinkedServiceTypeProperties: databricksProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureDatabricks,
	}

	if v, ok := d.GetOk("parameters"); ok {
		databricksLinkedService.Parameters = expandLinkedServiceParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		databricksLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		databricksLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		databricksLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: databricksLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Azure Databricks %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceDatabricksRead(d, meta)
}

func resourceDataFactoryLinkedServiceDatabricksRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
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

		return fmt.Errorf("retrieving Data Factory Databricks %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	databricks, ok := resp.Properties.AsAzureDatabricksLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Databricks %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureDatabricks, *resp.Type)
	}

	// Check the properties and verify if authentication is set to MSI
	if props := databricks.AzureDatabricksLinkedServiceTypeProperties; props != nil {
		d.Set("adb_domain", props.Domain)

		if props.Authentication != nil && props.Authentication == "MSI" {
			d.Set("msi_work_space_resource_id", props.WorkspaceResourceID)
		} else if accessToken := props.AccessToken; accessToken != nil {
			// We only process AzureKeyVaultSecreReference because a string based access token is masked with asterisks in the GET response
			// so we can't set it
			if keyVaultPassword, ok := accessToken.AsAzureKeyVaultSecretReference(); ok {
				if err := d.Set("key_vault_password", flattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
					return fmt.Errorf("setting `key_vault_password`: %+v", err)
				}
			}
		}

		instancePoolArray := make([]interface{}, 0)
		newClusterArray := make([]interface{}, 0)
		if props.ExistingClusterID != nil {
			if err := d.Set("existing_cluster_id", props.ExistingClusterID); err != nil {
				return fmt.Errorf("setting `existing_cluster_id`: %+v", err)
			}
		} else if id := props.InstancePoolID; id != nil {
			numOfWorkers := props.NewClusterNumOfWorker
			clusterVersion := props.NewClusterVersion

			minWorkers, maxWorkers, err := parseNumberOfWorkersProperties(numOfWorkers.(string))
			if err != nil {
				return fmt.Errorf("setting `instance_pool`: %+v", err)
			}

			instancePoolMap := map[string]interface{}{
				"instance_pool_id":      id,
				"min_number_of_workers": minWorkers,
				"cluster_version":       clusterVersion,
			}

			if maxWorkers != 0 {
				instancePoolMap["max_number_of_workers"] = maxWorkers
			}

			instancePoolArray = append(instancePoolArray, instancePoolMap)
		} else {
			// Process assuming it's a new cluster config
			numOfWorkers := props.NewClusterNumOfWorker
			clusterVersion := props.NewClusterVersion
			nodeType := props.NewClusterNodeType

			minWorkers, maxWorkers, err := parseNumberOfWorkersProperties(numOfWorkers.(string))
			if err != nil {
				return fmt.Errorf("setting `new_cluster_config`: %+v", err)
			}

			newClusterMap := map[string]interface{}{
				"min_number_of_workers": minWorkers,
				"cluster_version":       clusterVersion,
				"node_type":             nodeType,
			}

			if maxWorkers != 0 {
				newClusterMap["max_number_of_workers"] = maxWorkers
			}

			// Retrieve all the optional arguments
			if data := props.NewClusterDriverNodeType; data != nil {
				newClusterMap["driver_node_type"] = data
			}

			if data := props.NewClusterLogDestination; data != nil {
				newClusterMap["log_destination"] = data
			}

			if data := props.NewClusterSparkConf; data != nil {
				newClusterMap["spark_config"] = data
			}

			if data := props.NewClusterCustomTags; data != nil {
				newClusterMap["custom_tags"] = data
			}

			if data := props.NewClusterSparkEnvVars; data != nil {
				newClusterMap["spark_environment_variables"] = data
			}

			if data := props.NewClusterInitScripts; data != nil {
				newClusterMap["init_scripts"] = data
			}
			newClusterArray = append(newClusterArray, newClusterMap)
		}
		if err := d.Set("new_cluster_config", newClusterArray); err != nil {
			return fmt.Errorf("setting `new_cluster_config`: %+v", err)
		}
		if err := d.Set("instance_pool", instancePoolArray); err != nil {
			return fmt.Errorf("setting `instance_pool`: %+v", err)
		}
	}

	d.Set("additional_properties", databricks.AdditionalProperties)
	d.Set("description", databricks.Description)

	annotations := flattenDataFactoryAnnotations(databricks.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	parameters := flattenLinkedServiceParameters(databricks.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := databricks.ConnectVia; connectVia != nil {
		d.Set("integration_runtime_name", connectVia.ReferenceName)
	}

	return nil
}

func resourceDataFactoryLinkedServiceDatabricksDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Data Factory Databricks %s: %+v", *id, err)
		}
	}
	return nil
}

func buildNumberOfWorkersProperties(minWorkersProperty interface{}, maxWorkersProperty interface{}) (string, interface{}) {
	var err error

	// Default settings
	minWorkers := minWorkersProperty.(int)
	workersConfig := fmt.Sprintf("%d", minWorkers)

	// If max workers are set, we'll assume they want to setup autoscaling and throw an error if the configuration looks invalid
	if maxWorkers := maxWorkersProperty.(int); maxWorkers > 0 {
		if maxWorkers < minWorkers {
			err = fmt.Errorf("`max_number_of_workers` [%d]` cannot be less than `min_number_of_workers` [%d]", maxWorkers, minWorkers)
		} else if maxWorkers > minWorkers {
			workersConfig = fmt.Sprintf("%d:%d", minWorkers, maxWorkers)
		}
	}
	return workersConfig, err
}

func parseNumberOfWorkersProperties(numberOfWorkersProperty string) (int, int, error) {
	// The number of workers should be either a fixed number (no autoscaling) or have a format of min:max if autoscaling is set.
	numOfWorkersParts := strings.Split(numberOfWorkersProperty, ":")

	var min int
	var max int
	var err error
	switch len(numOfWorkersParts) {
	case 1:
		min, err = strconv.Atoi(numOfWorkersParts[0])
	case 2:
		if min, err = strconv.Atoi(numOfWorkersParts[0]); err == nil {
			max, err = strconv.Atoi(numOfWorkersParts[1])
		}
	default:
		err = fmt.Errorf("Number of workers property has unknown format: %s", numberOfWorkersProperty)
	}

	return min, max, err
}
