package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryLinkedServiceKusto() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceKustoCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceKustoRead,
		Update: resourceDataFactoryLinkedServiceKustoCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceKustoDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validate.DataFactoryID,
			},

			"kusto_endpoint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"kusto_database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"use_managed_identity": {
				Type:         pluginsdk.TypeBool,
				Optional:     true,
				Default:      false,
				ExactlyOneOf: []string{"service_principal_id", "use_managed_identity"},
			},

			"service_principal_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"service_principal_key"},
				ExactlyOneOf: []string{"service_principal_id", "use_managed_identity"},
			},

			"service_principal_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"service_principal_id"},
			},

			"tenant": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"service_principal_id"},
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

func resourceDataFactoryLinkedServiceKustoCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := parse.DataFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroup, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_kusto", id.ID())
		}
	}

	kustoLinkedService := &datafactory.AzureDataExplorerLinkedService{
		AzureDataExplorerLinkedServiceTypeProperties: &datafactory.AzureDataExplorerLinkedServiceTypeProperties{
			Endpoint: d.Get("kusto_endpoint").(string),
			Database: d.Get("kusto_database_name").(string),
		},
		Description: utils.String(d.Get("description").(string)),
		Type:        datafactory.TypeBasicLinkedServiceTypeAzureDataExplorer,
	}

	if d.Get("use_managed_identity").(bool) {
		kustoLinkedService.AzureDataExplorerLinkedServiceTypeProperties = &datafactory.AzureDataExplorerLinkedServiceTypeProperties{
			Endpoint: d.Get("kusto_endpoint").(string),
			Database: d.Get("kusto_database_name").(string),
		}
	} else if v, ok := d.GetOk("service_principal_id"); ok {
		kustoLinkedService.AzureDataExplorerLinkedServiceTypeProperties = &datafactory.AzureDataExplorerLinkedServiceTypeProperties{
			Endpoint:           d.Get("kusto_endpoint").(string),
			Database:           d.Get("kusto_database_name").(string),
			ServicePrincipalID: v.(string),
			ServicePrincipalKey: &datafactory.SecureString{
				Value: utils.String(d.Get("service_principal_key").(string)),
				Type:  datafactory.TypeSecureString,
			},
			Tenant: utils.String(d.Get("tenant").(string)),
		}
	} else {
		return fmt.Errorf("one of Managed Identity and service principal authentication must be set")
	}

	if v, ok := d.GetOk("parameters"); ok {
		kustoLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		kustoLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		kustoLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		kustoLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: kustoLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceKustoRead(d, meta)
}

func resourceDataFactoryLinkedServiceKustoRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	linkedService, ok := resp.Properties.AsAzureDataExplorerLinkedService()
	if !ok {
		return fmt.Errorf("classifiying %s: Expected: %q", id, datafactory.TypeBasicLinkedServiceTypeAzureDataExplorer)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", parse.NewDataFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())
	d.Set("additional_properties", linkedService.AdditionalProperties)
	d.Set("description", linkedService.Description)
	if err := d.Set("annotations", flattenDataFactoryAnnotations(linkedService.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}
	if err := d.Set("parameters", flattenDataFactoryParameters(linkedService.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	integrationRuntimeName := ""
	if linkedService.ConnectVia != nil && linkedService.ConnectVia.ReferenceName != nil {
		integrationRuntimeName = *linkedService.ConnectVia.ReferenceName
	}
	d.Set("integration_runtime_name", integrationRuntimeName)

	if prop := linkedService.AzureDataExplorerLinkedServiceTypeProperties; prop != nil {
		d.Set("kusto_endpoint", prop.Endpoint)
		d.Set("kusto_database_name", prop.Database)
		d.Set("tenant", prop.Tenant)
		d.Set("service_principal_id", prop.ServicePrincipalID)

		useManagedIdentity := true
		if prop.ServicePrincipalID != nil {
			useManagedIdentity = false
		}
		d.Set("use_managed_identity", useManagedIdentity)
	}

	return nil
}

func resourceDataFactoryLinkedServiceKustoDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
