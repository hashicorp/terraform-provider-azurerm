// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

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

func resourceDataFactoryLinkedServiceCosmosDb() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceCosmosDbCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceCosmosDbRead,
		Update: resourceDataFactoryLinkedServiceCosmosDbCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceCosmosDbDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, importDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeCosmosDb)),

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

			"connection_string": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"account_endpoint", "account_key"},
				DiffSuppressFunc: azureRmDataFactoryLinkedServiceConnectionStringDiff,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"account_endpoint": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"connection_string"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"account_key": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"connection_string"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"database": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

func resourceDataFactoryLinkedServiceCosmosDbCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for presence of existing Data Factory CosmosDb %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_cosmosdb", id.ID())
		}
	}

	cosmosdbProperties := &datafactory.CosmosDbLinkedServiceTypeProperties{}

	endpoint := d.Get("account_endpoint").(string)
	accountKey := d.Get("account_key").(string)
	databaseName := d.Get("database").(string)

	isAccountDetailUsed := endpoint != "" && accountKey != "" && databaseName != ""

	if isAccountDetailUsed {
		accountKeySecureString := datafactory.SecureString{
			Value: &accountKey,
			Type:  datafactory.TypeSecureString,
		}
		cosmosdbProperties.AccountEndpoint = endpoint
		cosmosdbProperties.AccountKey = accountKeySecureString
		cosmosdbProperties.Database = databaseName
	} else {
		connectionString := d.Get("connection_string").(string)
		connectionStringSecureString := datafactory.SecureString{
			Value: &connectionString,
			Type:  datafactory.TypeSecureString,
		}
		cosmosdbProperties.ConnectionString = connectionStringSecureString
		cosmosdbProperties.Database = databaseName
	}

	description := d.Get("description").(string)

	cosmosdbLinkedService := &datafactory.CosmosDbLinkedService{
		Description:                         &description,
		CosmosDbLinkedServiceTypeProperties: cosmosdbProperties,
		Type:                                datafactory.TypeBasicLinkedServiceTypeCosmosDb,
	}

	if v, ok := d.GetOk("parameters"); ok {
		cosmosdbLinkedService.Parameters = expandLinkedServiceParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		cosmosdbLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		cosmosdbLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		cosmosdbLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: cosmosdbLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory CosmosDb %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceCosmosDbRead(d, meta)
}

func resourceDataFactoryLinkedServiceCosmosDbRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("retrieving Data Factory CosmosDB %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	cosmosdb, ok := resp.Properties.AsCosmosDbLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory CosmosDb %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeCosmosDb, *resp.Type)
	}

	d.Set("additional_properties", cosmosdb.AdditionalProperties)
	d.Set("description", cosmosdb.Description)

	annotations := flattenDataFactoryAnnotations(cosmosdb.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	parameters := flattenLinkedServiceParameters(cosmosdb.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := cosmosdb.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	accountEndpoint := cosmosdb.CosmosDbLinkedServiceTypeProperties.AccountEndpoint
	if accountEndpoint != "" {
		d.Set("account_endpoint", accountEndpoint)
	}

	databaseName := cosmosdb.CosmosDbLinkedServiceTypeProperties.Database
	d.Set("database", databaseName)

	return nil
}

func resourceDataFactoryLinkedServiceCosmosDbDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("deleting Data Factory CosmosDb %s: %+v", *id, err)
		}
	}

	return nil
}
