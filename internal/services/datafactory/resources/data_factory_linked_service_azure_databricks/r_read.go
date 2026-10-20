// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_databricks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceDatabricksRead(d *pluginsdk.ResourceData, meta any) error {
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
		if response.WasNotFound(resp.Response.Response) {
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
			d.Set("msi_workspace_id", props.WorkspaceResourceID)
		} else if accessToken := props.AccessToken; accessToken != nil {
			// We only process AzureKeyVaultSecreReference because a string based access token is masked with asterisks in the GET response
			// so we can't set it
			if keyVaultPassword, ok := accessToken.AsAzureKeyVaultSecretReference(); ok {
				if err := d.Set("key_vault_password", helpers.FlattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
					return fmt.Errorf("setting `key_vault_password`: %+v", err)
				}
			}
		}

		instancePoolArray := make([]any, 0)
		newClusterArray := make([]any, 0)
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

			instancePoolMap := map[string]any{
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

			newClusterMap := map[string]any{
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

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(databricks.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(databricks.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := databricks.ConnectVia; connectVia != nil {
		d.Set("integration_runtime_name", connectVia.ReferenceName)
	}

	return nil
}

func buildNumberOfWorkersProperties(minWorkersProperty any, maxWorkersProperty any) (string, any) {
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
		err = fmt.Errorf("number of workers property has unknown format: %s", numberOfWorkersProperty)
	}

	return min, max, err
}
