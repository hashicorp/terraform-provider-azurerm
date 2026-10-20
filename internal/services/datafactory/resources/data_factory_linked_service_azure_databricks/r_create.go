// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_databricks

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceDatabricksCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Data Factory Databricks %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_databricks", id.ID())
			}
		}
	}

	var databricksProperties *datafactory.AzureDatabricksLinkedServiceTypeProperties

	// Check if the MSI authentication block is set
	msiAuth := d.Get("msi_workspace_id")
	accessTokenAuth := d.Get("access_token").(string)
	accessTokenKeyVaultAuth := d.Get("key_vault_password").([]any)

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
				Value: pointer.To(accessTokenAuth),
				Type:  datafactory.TypeSecureString,
			},
		}
	}

	if len(accessTokenKeyVaultAuth) > 0 && accessTokenKeyVaultAuth[0] != nil {
		databricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			AccessToken: helpers.ExpandAzureKeyVaultSecretReference(accessTokenKeyVaultAuth),
		}
	}

	// Set the other type properties
	databricksProperties.Domain = d.Get("adb_domain").(string)

	if v, ok := d.GetOk("existing_cluster_id"); ok {
		databricksProperties.ExistingClusterID = v.(string)
	}

	if v, ok := d.GetOk("instance_pool"); ok && v.([]any)[0] != nil {
		instancePoolMap := v.([]any)[0].(map[string]any)

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

	if v, ok := d.GetOk("new_cluster_config"); ok && v.([]any)[0] != nil {
		newClusterMap := v.([]any)[0].(map[string]any)

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
			if sparkConfig := newClusterMap["spark_config"].(map[string]any); len(sparkConfig) > 0 {
				databricksProperties.NewClusterSparkConf = sparkConfig
			}
		}

		if newClusterMap["spark_environment_variables"] != nil {
			if sparkEnvVars := newClusterMap["spark_environment_variables"].(map[string]any); len(sparkEnvVars) > 0 {
				databricksProperties.NewClusterSparkEnvVars = sparkEnvVars
			}
		}

		if newClusterMap["custom_tags"] != nil {
			if customTags := newClusterMap["custom_tags"].(map[string]any); len(customTags) > 0 {
				databricksProperties.NewClusterCustomTags = customTags
			}
		}

		databricksProperties.NewClusterInitScripts = pointer.To(newClusterMap["init_scripts"])
	}

	databricksLinkedService := &datafactory.AzureDatabricksLinkedService{
		Description: pointer.To(d.Get("description").(string)),
		AzureDatabricksLinkedServiceTypeProperties: databricksProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureDatabricks,
	}

	if v, ok := d.GetOk("parameters"); ok {
		databricksLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		databricksLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		databricksLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		databricksLinkedService.Annotations = pointer.To(v.([]any))
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
