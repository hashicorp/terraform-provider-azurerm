package azurerm

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryLinkedServiceMySQL() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryLinkedServiceMySQLCreateOrUpdate,
		Read:   resourceArmDataFactoryLinkedServiceMySQLRead,
		Update: resourceArmDataFactoryLinkedServiceMySQLCreateOrUpdate,
		Delete: resourceArmDataFactoryLinkedServiceMySQLDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMDataFactoryLinkedServiceDatasetName,
			},

			"data_factory_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"connection_string": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: azureRmDataFactoryLinkedServiceConnectionStringDiff,
				ValidateFunc:     validate.NoEmptyStrings,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"integration_runtime_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
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
			},
		},
	}
}

func resourceArmDataFactoryLinkedServiceMySQLCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.LinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_mysql", *existing.ID)
		}
	}

	connectionString := d.Get("connection_string").(string)
	secureString := datafactory.SecureString{
		Value: &connectionString,
		Type:  datafactory.TypeSecureString,
	}

	mysqlProperties := &datafactory.MySQLLinkedServiceTypeProperties{
		ConnectionString: &secureString,
	}

	description := d.Get("description").(string)

	mysqlLinkedService := &datafactory.MySQLLinkedService{
		Description:                      &description,
		MySQLLinkedServiceTypeProperties: mysqlProperties,
		Type:                             datafactory.TypeMySQL,
	}

	if v, ok := d.GetOk("parameters"); ok {
		mysqlLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		mysqlLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		mysqlLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		mysqlLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: mysqlLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryLinkedServiceMySQLRead(d, meta)
}

func resourceArmDataFactoryLinkedServiceMySQLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.LinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

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

		return fmt.Errorf("Error retrieving Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	mysql, ok := resp.Properties.AsMySQLLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeMySQL, *resp.Type)
	}

	d.Set("additional_properties", mysql.AdditionalProperties)

	if mysql.Description != nil {
		d.Set("description", *mysql.Description)
	}

	annotations := flattenDataFactoryAnnotations(mysql.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(mysql.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := mysql.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceArmDataFactoryLinkedServiceMySQLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.LinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

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
			return fmt.Errorf("Error deleting Data Factory Linked Service MySQL %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
