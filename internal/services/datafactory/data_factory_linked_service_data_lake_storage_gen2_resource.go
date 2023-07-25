// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryLinkedServiceDataLakeStorageGen2() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceDataLakeStorageGen2CreateUpdate,
		Read:   resourceDataFactoryLinkedServiceDataLakeStorageGen2Read,
		Update: resourceDataFactoryLinkedServiceDataLakeStorageGen2CreateUpdate,
		Delete: resourceDataFactoryLinkedServiceDataLakeStorageGen2Delete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}, importDataFactoryLinkedService(datafactory.TypeBasicLinkedServiceTypeAzureBlobFS)),

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

			"url": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"use_managed_identity": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"service_principal_key", "service_principal_id", "storage_account_key", "tenant"},
				AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
			},

			"service_principal_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsUUID,
				RequiredWith:  []string{"service_principal_key", "tenant"},
				ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
				AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
			},

			"service_principal_key": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				RequiredWith:  []string{"service_principal_id", "tenant"},
				ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
				AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
			},

			"storage_account_key": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"service_principal_id", "service_principal_key", "use_managed_identity", "tenant"},
				AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
			},

			"tenant": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				RequiredWith:  []string{"service_principal_id", "service_principal_key"},
				ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
				AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
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

func resourceDataFactoryLinkedServiceDataLakeStorageGen2CreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for presence of existing Data Factory Data Lake Storage Gen2 %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_data_lake_storage_gen2", id.ID())
		}
	}

	var datalakeStorageGen2Properties *datafactory.AzureBlobFSLinkedServiceTypeProperties

	if d.Get("use_managed_identity").(bool) {
		datalakeStorageGen2Properties = &datafactory.AzureBlobFSLinkedServiceTypeProperties{
			URL: utils.String(d.Get("url").(string)),
		}
	} else if v, ok := d.GetOk("storage_account_key"); ok {
		datalakeStorageGen2Properties = &datafactory.AzureBlobFSLinkedServiceTypeProperties{
			URL: utils.String(d.Get("url").(string)),
			AccountKey: datafactory.SecureString{
				Value: utils.String(v.(string)),
				Type:  datafactory.TypeSecureString,
			},
		}
	} else {
		secureString := datafactory.SecureString{
			Value: utils.String(d.Get("service_principal_key").(string)),
			Type:  datafactory.TypeSecureString,
		}

		datalakeStorageGen2Properties = &datafactory.AzureBlobFSLinkedServiceTypeProperties{
			URL:                 utils.String(d.Get("url").(string)),
			ServicePrincipalID:  utils.String(d.Get("service_principal_id").(string)),
			Tenant:              utils.String(d.Get("tenant").(string)),
			ServicePrincipalKey: &secureString,
		}
	}

	datalakeStorageGen2LinkedService := &datafactory.AzureBlobFSLinkedService{
		Description:                            utils.String(d.Get("description").(string)),
		AzureBlobFSLinkedServiceTypeProperties: datalakeStorageGen2Properties,
		Type:                                   datafactory.TypeBasicLinkedServiceTypeAzureBlobFS,
	}

	if v, ok := d.GetOk("parameters"); ok {
		datalakeStorageGen2LinkedService.Parameters = expandLinkedServiceParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		datalakeStorageGen2LinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		datalakeStorageGen2LinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		datalakeStorageGen2LinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: datalakeStorageGen2LinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Data Lake Storage Gen2 %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceDataLakeStorageGen2Read(d, meta)
}

func resourceDataFactoryLinkedServiceDataLakeStorageGen2Read(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("retrieving Data Factory Data Lake Storage Gen2 %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	dataLakeStorageGen2, ok := resp.Properties.AsAzureBlobFSLinkedService()

	if !ok {
		return fmt.Errorf("classifying Data Factory Data Lake Storage Gen2 %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureBlobFS, *resp.Type)
	}

	if dataLakeStorageGen2.Tenant != nil {
		d.Set("tenant", dataLakeStorageGen2.Tenant)
	}

	if dataLakeStorageGen2.ServicePrincipalID != nil {
		d.Set("service_principal_id", dataLakeStorageGen2.ServicePrincipalID)
		d.Set("use_managed_identity", false)
	}

	if dataLakeStorageGen2.URL != nil {
		d.Set("url", dataLakeStorageGen2.URL)
	}

	d.Set("additional_properties", dataLakeStorageGen2.AdditionalProperties)

	if dataLakeStorageGen2.Description != nil {
		d.Set("description", dataLakeStorageGen2.Description)
	}

	annotations := flattenDataFactoryAnnotations(dataLakeStorageGen2.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	parameters := flattenLinkedServiceParameters(dataLakeStorageGen2.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := dataLakeStorageGen2.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceDataFactoryLinkedServiceDataLakeStorageGen2Delete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("deleting Data Factory Data Lake Storage Gen2 %s: %+v", *id, err)
		}
	}

	return nil
}
