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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryLinkedServiceSFTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataFactoryLinkedServiceSFTPCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceSFTPRead,
		Update: resourceDataFactoryLinkedServiceSFTPCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceSFTPDelete,

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

			"authentication_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"host": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

func resourceDataFactoryLinkedServiceSFTPCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service SFTP %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_sftp", *existing.ID)
		}
	}

	authenticationType := d.Get("authentication_type").(string)

	host := d.Get("host").(string)
	port := d.Get("port").(int)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	passwordSecureString := datafactory.SecureString{
		Value: &password,
		Type:  datafactory.TypeSecureString,
	}

	sftpProperties := &datafactory.SftpServerLinkedServiceTypeProperties{
		Host:               utils.String(host),
		Port:               port,
		AuthenticationType: datafactory.SftpAuthenticationType(authenticationType),
		UserName:           utils.String(username),
		Password:           &passwordSecureString,
	}

	description := d.Get("description").(string)

	sftpLinkedService := &datafactory.SftpServerLinkedService{
		Description:                           &description,
		SftpServerLinkedServiceTypeProperties: sftpProperties,
		Type:                                  datafactory.TypeSftp,
	}

	if v, ok := d.GetOk("parameters"); ok {
		sftpLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		sftpLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		sftpLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		sftpLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: sftpLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service SFTP Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service SFTP Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service SFTP Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryLinkedServiceSFTPRead(d, meta)
}

func resourceDataFactoryLinkedServiceSFTPRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error retrieving Data Factory Linked Service SFTP %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	sftp, ok := resp.Properties.AsSftpServerLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service SFTP %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeSftp, *resp.Type)
	}

	d.Set("authentication_type", sftp.AuthenticationType)
	d.Set("username", sftp.UserName)
	d.Set("port", sftp.Port)
	d.Set("host", sftp.Host)

	d.Set("additional_properties", sftp.AdditionalProperties)
	d.Set("description", sftp.Description)

	annotations := flattenDataFactoryAnnotations(sftp.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(sftp.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := sftp.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceDataFactoryLinkedServiceSFTPDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Data Factory Linked Service SFTP %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
