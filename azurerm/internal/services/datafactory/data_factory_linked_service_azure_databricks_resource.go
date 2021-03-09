package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	databricksValidator "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryLinkedServiceAzureDatabricks() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataFactoryLinkedServiceDatabricksCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceDatabricksRead,
		Update: resourceDataFactoryLinkedServiceDatabricksCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceDatabricksDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMDataFactoryLinkedServiceDatasetName,
			},

			"data_factory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			// Authentication types
			"authentication_msi": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_resource_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: databricksValidator.WorkspaceID,
						},
					},
				},
				ExactlyOneOf: []string{"authentication_access_token", "authentication_msi", "authentication_key_vault_password"},
			},

			"authentication_access_token": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_token": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"authentication_access_token", "authentication_msi", "authentication_key_vault_password"},
			},

			"authentication_key_vault_password": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"linked_service_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"secret_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"authentication_access_token", "authentication_msi", "authentication_key_vault_password"},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"adb_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Cluster types [existing cluster, new cluster, interactive pools]
			"existing_cluster_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
			},

			"new_cluster_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"custom_tags": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						// Consider changing this to min and adding an optional max since autoscaling uses the format min:max (e.g. 1:10)
						"number_of_workers": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "1",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"cluster_version": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"spark_config": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"spark_environment_variables": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"log_destination": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"init_scripts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"driver_node_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
			},

			"instance_pool": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// Consider changing this to min and adding an optional max since autoscaling uses the format min:max (e.g. 1:10)
						"number_of_workers": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "1",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"instance_pool_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"cluster_version": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
			},

			"integration_runtime_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"additional_properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDataFactoryLinkedServiceDatabricksCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service Databricks %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_data_lake_storage_gen2", *existing.ID)
		}
	}

	var DatabricksProperties *datafactory.AzureDatabricksLinkedServiceTypeProperties

	// Check if the MSI authentication block is set
	msiAuth := d.Get("authentication_msi").([]interface{})
	accessTokenAuth := d.Get("authentication_access_token").([]interface{})
	accessTokenKeyVaultAuth := d.Get("authentication_key_vault_password").([]interface{})

	if len(msiAuth) > 0 {
		// MSI Auth map has data
		workspaceResourceID := msiAuth[0].(map[string]interface{})["workspace_resource_id"].(string)

		DatabricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			Authentication:      "MSI",
			WorkspaceResourceID: workspaceResourceID,
		}
	}
	if len(accessTokenAuth) > 0 {
		// Setup authentication using access tokens

		accessToken := accessTokenAuth[0].(map[string]interface{})["access_token"].(string)

		accessTokenAsSecureString := datafactory.SecureString{
			Value: utils.String(accessToken),
			Type:  datafactory.TypeSecureString,
		}

		// Assign the access token in the properties block
		DatabricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			AccessToken: &accessTokenAsSecureString,
		}
	} else {
		DatabricksProperties = &datafactory.AzureDatabricksLinkedServiceTypeProperties{
			AccessToken: expandAzureKeyVaultPassword(accessTokenKeyVaultAuth),
		}
	}

	// Set the other type properties
	// Domain
	DatabricksProperties.Domain = d.Get("adb_domain").(string)

	// Check if the
	if v, ok := d.GetOk("existing_cluster_id"); ok {
		DatabricksProperties.ExistingClusterID = v.(string)
	}

	if v, ok := d.GetOk("instance_pool"); ok {
		instancePoolMap := v.([]interface{})[0].(map[string]interface{})

		// Process the instance pool identifier
		if data := instancePoolMap["instance_pool_id"]; data != nil {
			DatabricksProperties.InstancePoolID = data
		}

		// Process the cluster version
		if data := instancePoolMap["cluster_version"]; data != nil {
			DatabricksProperties.NewClusterVersion = data
		}

		// Process the number of workers
		if data := instancePoolMap["number_of_workers"]; data != nil {
			DatabricksProperties.NewClusterNumOfWorker = data
		}
	} else if v, ok := d.GetOk("new_cluster_config"); ok {
		newClusterMap := v.([]interface{})[0].(map[string]interface{})

		// Process the cluster version
		if data := newClusterMap["cluster_version"]; data != nil {
			DatabricksProperties.NewClusterVersion = data
		}

		// Process the number of workers
		if data := newClusterMap["number_of_workers"]; data != nil {
			DatabricksProperties.NewClusterNumOfWorker = data
		}

		// Process the node type
		if data := newClusterMap["node_type"]; data != nil {
			DatabricksProperties.NewClusterNodeType = data
		}

		if data := newClusterMap["driver_node_type"]; data != nil {
			DatabricksProperties.NewClusterDriverNodeType = data
		}

		if data := newClusterMap["log_destination"]; data != nil {
			DatabricksProperties.NewClusterLogDestination = data
		}

		if sparkConfig := newClusterMap["spark_config"].(map[string]interface{}); len(sparkConfig) > 0 {
			DatabricksProperties.NewClusterSparkConf = sparkConfig
		}

		if sparkEnvVars := newClusterMap["spark_environment_variables"].(map[string]interface{}); len(sparkEnvVars) > 0 {
			DatabricksProperties.NewClusterSparkEnvVars = sparkEnvVars
		}

		if customTags := newClusterMap["custom_tags"].(map[string]interface{}); len(customTags) > 0 {
			DatabricksProperties.NewClusterCustomTags = customTags
		}

		initScripts := newClusterMap["init_scripts"]
		DatabricksProperties.NewClusterInitScripts = &initScripts
	}

	DatabricksLinkedService := &datafactory.AzureDatabricksLinkedService{
		Description: utils.String(d.Get("description").(string)),
		AzureDatabricksLinkedServiceTypeProperties: DatabricksProperties,
		Type: datafactory.TypeAzureDatabricks,
	}

	if v, ok := d.GetOk("parameters"); ok {
		DatabricksLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		DatabricksLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		DatabricksLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		DatabricksLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: DatabricksLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service Azure Databricks %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service Databricks %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service Databricks %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryLinkedServiceDatabricksRead(d, meta)
}

func resourceDataFactoryLinkedServiceDatabricksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["linkedservices"]

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Linked Service Databricks %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	Databricks, ok := resp.Properties.AsAzureDatabricksLinkedService()

	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service Databricks %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeAzureDatabricks, *resp.Type)
	}

	// Check the properties and verify is authentication is set to MSI
	if props := Databricks.AzureDatabricksLinkedServiceTypeProperties; props != nil {
		if props.Authentication != nil && props.Authentication == "MSI" {
			authenticationMsi := make(map[string]interface{})
			authenticationMsi["workspace_resource_id"] = props.WorkspaceResourceID
			authenticationMsiArray := []interface{}{authenticationMsi}
			d.Set("authentication_msi", authenticationMsiArray)
		} else if props.AccessToken != nil {
			// Check the data type of the access token so we know how to process it.

			if accessToken := props.AccessToken; accessToken != nil {
				if keyVaultPassword, ok := accessToken.AsAzureKeyVaultSecretReference(); ok {
					if err := d.Set("authentication_key_vault_password", flattenAzureKeyVaultPassword(keyVaultPassword)); err != nil {
						return fmt.Errorf("setting `authentication_key_vault_password`: %+v", err)
					}
				}
			}
		}

		// Process the domain
		if props.Domain != nil {
			d.Set("adb_domain", props.Domain)
		}

		// Process the cluster information
		if props.ExistingClusterID != nil {
			d.Set("existing_cluster_id", props.ExistingClusterID)
		} else if id := props.InstancePoolID; id != nil {
			// Process the values for instance pool configuration
			numOfWorkers := props.NewClusterNumOfWorker
			clusterVersion := props.NewClusterVersion

			instancePoolMap := map[string]interface{}{
				"instance_pool_id":  id,
				"number_of_workers": numOfWorkers,
				"cluster_version":   clusterVersion,
			}

			instancePoolArray := []interface{}{instancePoolMap}
			d.Set("instance_pool", instancePoolArray)
		} else {
			// Process assuming it's a new cluster config
			numOfWorkers := props.NewClusterNumOfWorker
			clusterVersion := props.NewClusterVersion
			nodeType := props.NewClusterNodeType

			newClusterMap := map[string]interface{}{
				"number_of_workers": numOfWorkers,
				"cluster_version":   clusterVersion,
				"node_type":         nodeType,
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

			// Set the ResourceData with the map
			newClusterArray := []interface{}{newClusterMap}

			d.Set("new_cluster_config", newClusterArray)
		}
	}

	d.Set("additional_properties", Databricks.AdditionalProperties)

	if Databricks.Description != nil {
		d.Set("description", Databricks.Description)
	}

	annotations := flattenDataFactoryAnnotations(Databricks.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(Databricks.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := Databricks.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceDataFactoryLinkedServiceDatabricksDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["linkedservices"]

	response, err := client.Delete(ctx, resourceGroup, dataFactoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory Linked Service Databricks %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}
	return nil
}
