package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryLinkedServiceSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryLinkedServiceSQLServerCreateOrUpdate,
		Read:   resourceArmDataFactoryLinkedServiceSQLServerRead,
		Update: resourceArmDataFactoryLinkedServiceSQLServerCreateOrUpdate,
		Delete: resourceArmDataFactoryLinkedServiceSQLServerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO add validation
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

			"resource_group_name": resourceGroupNameSchema(),

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"connection_string": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: azureRmDataFactoryLinkedServiceConnectionStringDiff,
			},

			"integration_runtime_name": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceArmDataFactoryLinkedServiceSQLServerCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryLinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_sql_server_linked_service", *existing.ID)
		}
	}

	sqlServerProperties := &datafactory.SQLServerLinkedServiceTypeProperties{
		ConnectionString: d.Get("connection_string").(string),
	}

	description := d.Get("description").(string)
	parameters := expandDataFactoryLinkedServiceParameters(d)
	additionalProperties := expandDataFactoryLinkedServiceAdditionalProperties(d)
	annotations := expandDataFactoryLinkedServiceAnnotations(d)

	sqlServerLinkedService := &datafactory.SQLServerLinkedService{
		Description:                          &description,
		SQLServerLinkedServiceTypeProperties: sqlServerProperties,
		Type:                                 datafactory.TypeSQLServer,
		ConnectVia:                           expandDataFactoryLinkedServiceIntegrationRuntime(d),
		Parameters:                           parameters,
		Annotations:                          annotations,
		AdditionalProperties:                 additionalProperties,
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: sqlServerLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryLinkedServiceSQLServerRead(d, meta)
}

func resourceArmDataFactoryLinkedServiceSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryLinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

		return fmt.Errorf("Error retrieving Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	sqlServer, ok := resp.Properties.AsSQLServerLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeSQLServer, *resp.Type)
	}

	d.Set("additional_properties", sqlServer.AdditionalProperties)

	if sqlServer.Description != nil {
		d.Set("description", *sqlServer.Description)
	}

	annotations := flattenDataFactoryLinkedServiceAnnotations(sqlServer.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryLinkedServiceParameters(sqlServer.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := sqlServer.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", *connectVia.ReferenceName)
		}
	}

	if properties := sqlServer.SQLServerLinkedServiceTypeProperties; properties != nil {
		if properties.ConnectionString != nil {
			val, ok := properties.ConnectionString.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping connection string %q since it's not a string", val)
			} else {
				d.Set("connection_string", properties.ConnectionString.(string))
			}
		}
	}

	return nil
}

func resourceArmDataFactoryLinkedServiceSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryLinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["linkedservices"]

	response, err := client.Delete(ctx, resourceGroup, dataFactoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory Linked Service SQL Server %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}

// Because the password isn't returned from the api in the connection string, we'll check all
// but the password string and return true if they match.
func azureRmDataFactoryLinkedServiceConnectionStringDiff(k, old string, new string, d *schema.ResourceData) bool {
	oldSplit := strings.Split(strings.ToLower(old), ";")
	newSplit := strings.Split(strings.ToLower(new), ";")

	sort.Strings(oldSplit)
	sort.Strings(newSplit)

	// We need to remove the password from the new string since it isn't returned from the api
	for i, v := range newSplit {
		if strings.HasPrefix(v, "password") {
			newSplit = append(newSplit[:i], newSplit[i+1:]...)
		}
	}

	if len(oldSplit) != len(newSplit) {
		return false
	}

	// We'll error out if we find any differences between the old and the new connectiong strings
	for i := range oldSplit {
		if !strings.EqualFold(oldSplit[i], newSplit[i]) {
			return false
		}
	}

	return true
}

func expandDataFactoryLinkedServiceIntegrationRuntime(d *schema.ResourceData) *datafactory.IntegrationRuntimeReference {
	if v, ok := d.GetOk("integration_runtime_name"); ok {
		integrationRuntimeName := v.(string)
		typeString := "IntegrationRuntimeReference"

		return &datafactory.IntegrationRuntimeReference{
			ReferenceName: &integrationRuntimeName,
			Type:          &typeString,
		}
	}
	return nil
}

func expandDataFactoryLinkedServiceParameters(d *schema.ResourceData) map[string]*datafactory.ParameterSpecification {
	input := d.Get("parameters").(map[string]interface{})

	output := make(map[string]*datafactory.ParameterSpecification)

	for k, v := range input {
		output[k] = &datafactory.ParameterSpecification{
			Type:         datafactory.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func expandDataFactoryLinkedServiceAdditionalProperties(d *schema.ResourceData) map[string]interface{} {
	input := d.Get("additional_properties").(map[string]interface{})

	output := make(map[string]interface{})

	for k, v := range input {
		output[k] = v
	}

	return output
}

func flattenDataFactoryLinkedServiceParameters(input map[string]*datafactory.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func expandDataFactoryLinkedServiceAnnotations(d *schema.ResourceData) *[]interface{} {
	input := d.Get("annotations").([]interface{})
	annotations := make([]interface{}, 0)

	for _, annotation := range input {
		annotations = append(annotations, annotation.(string))
	}

	return &annotations
}

func flattenDataFactoryLinkedServiceAnnotations(input *[]interface{}) []string {
	annotations := make([]string, 0)
	if input != nil {
		for _, annotation := range *input {
			val, ok := annotation.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping annotation %q since it's not a string", val)
			}
			annotations = append(annotations, val)
		}
	}
	return annotations
}
