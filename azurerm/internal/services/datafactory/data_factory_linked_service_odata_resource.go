package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryLinkedServiceOData() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryLinkedServiceODataCreateUpdate,
		Read:   resourceArmDataFactoryLinkedServiceODataRead,
		Update: resourceArmDataFactoryLinkedServiceWebCreateUpdate,
		Delete: resourceArmDataFactoryLinkedServiceODataDelete,

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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.ODataAuthenticationTypeAadServicePrincipal),
					string(datafactory.ODataAuthenticationTypeAnonymous),
					string(datafactory.ODataAuthenticationTypeBasic),
					string(datafactory.ODataAuthenticationTypeWindows),
				}, false),
			},

			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"aad_resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tenant": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			"service_principal_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			"aad_service_principal_credential_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.ServicePrincipalCert),
					string(datafactory.ServicePrincipalKey),
				}, false),
			},

			"service_principal_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"service_principal_embedded_cert": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"service_principal_embedded_cert_password": {
				Type:         schema.TypeString,
				Optional:     true,
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

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmDataFactoryLinkedServiceODataCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service OData Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_odata", *existing.ID)
		}
	}

	description := d.Get("description").(string)

	odataLinkedService := &datafactory.ODataLinkedService{
		Description: &description,
		Type:        datafactory.TypeOData,
	}

	url := d.Get("url").(string)
	authenticationType := d.Get("authentication_type").(string)

	if authenticationType == string(datafactory.ODataAuthenticationTypeAadServicePrincipal) {

		servicePrincipalId := d.Get("service_principal_id").(string)
		aadServicePrincipalCredentialType := d.Get("aad_service_principal_credential_type").(string)
		tenant := d.Get("tenant").(string)
		aadResourceID := d.Get("aad_resource_id").(string)

		if aadServicePrincipalCredentialType == string(datafactory.ServicePrincipalCert) {
			servicePrincipalEmbeddedCertSecureString := datafactory.SecureString{
				Value: utils.String(d.Get("service_principal_embedded_cert").(string)),
				Type:  datafactory.TypeSecureString,
			}
			servicePrincipalEmbeddedCertPasswordSecureString := datafactory.SecureString{
				Value: utils.String(d.Get("service_principal_embedded_cert_password").(string)),
				Type:  datafactory.TypeSecureString,
			}
			servicePrincipalAuthProperties := &datafactory.ODataLinkedServiceTypeProperties{
				AuthenticationType:                   datafactory.ODataAuthenticationType(authenticationType),
				URL:                                  utils.String(url),
				ServicePrincipalID:                   utils.String(servicePrincipalId),
				AadServicePrincipalCredentialType:    datafactory.ODataAadServicePrincipalCredentialType(aadServicePrincipalCredentialType),
				ServicePrincipalEmbeddedCert:         &servicePrincipalEmbeddedCertSecureString,
				ServicePrincipalEmbeddedCertPassword: &servicePrincipalEmbeddedCertPasswordSecureString,
				Tenant:                               utils.String(tenant),
				AadResourceID:                        utils.String(aadResourceID),
			}
			odataLinkedService.ODataLinkedServiceTypeProperties = servicePrincipalAuthProperties

		} else if aadServicePrincipalCredentialType == string(datafactory.ServicePrincipalKey) {
			servicePrincipalKeySecureString := datafactory.SecureString{
				Value: utils.String(d.Get("service_principal_key").(string)),
				Type:  datafactory.TypeSecureString,
			}
			servicePrincipalAuthProperties := &datafactory.ODataLinkedServiceTypeProperties{
				AuthenticationType:                datafactory.ODataAuthenticationType(authenticationType),
				URL:                               utils.String(url),
				ServicePrincipalID:                utils.String(servicePrincipalId),
				AadServicePrincipalCredentialType: datafactory.ODataAadServicePrincipalCredentialType(aadServicePrincipalCredentialType),
				ServicePrincipalKey:               &servicePrincipalKeySecureString,
				Tenant:                            utils.String(tenant),
				AadResourceID:                     utils.String(aadResourceID),
			}
			odataLinkedService.ODataLinkedServiceTypeProperties = servicePrincipalAuthProperties
		}
	}

	if authenticationType == string(datafactory.ODataAuthenticationTypeAnonymous) {
		anonAuthProperties := &datafactory.ODataLinkedServiceTypeProperties{
			AuthenticationType: datafactory.ODataAuthenticationType(authenticationType),
			URL:                utils.String(url),
		}
		odataLinkedService.ODataLinkedServiceTypeProperties = anonAuthProperties
	}

	if authenticationType == string(datafactory.ODataAuthenticationTypeBasic) || authenticationType == string(datafactory.ODataAuthenticationTypeWindows) {
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		passwordSecureString := datafactory.SecureString{
			Value: &password,
			Type:  datafactory.TypeSecureString,
		}
		basicAuthProperties := &datafactory.ODataLinkedServiceTypeProperties{
			AuthenticationType: datafactory.ODataAuthenticationType(authenticationType),
			URL:                utils.String(url),
			UserName:           username,
			Password:           &passwordSecureString,
		}
		odataLinkedService.ODataLinkedServiceTypeProperties = basicAuthProperties
	}

	if v, ok := d.GetOk("parameters"); ok {
		odataLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		odataLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		odataLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		odataLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: odataLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service OData Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service OData Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service OData Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryLinkedServiceODataRead(d, meta)
}

func resourceArmDataFactoryLinkedServiceODataRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error retrieving Data Factory Linked Service Web %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	odata, ok := resp.Properties.AsODataLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service Web %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeWeb, *resp.Type)
	}

	props := odata.ODataLinkedServiceTypeProperties
	d.Set("authentication_type", props.AuthenticationType)
	d.Set("url", props.URL)
	if props.AuthenticationType == datafactory.ODataAuthenticationTypeBasic || props.AuthenticationType == datafactory.ODataAuthenticationTypeWindows {
		d.Set("username", props.UserName)
		d.Set("password", props.Password)
	}
	if props.AuthenticationType == datafactory.ODataAuthenticationTypeAadServicePrincipal {
		d.Set("aad_resource_id", props.AadResourceID)
		d.Set("tenant", props.Tenant)
		d.Set("service_principal_id", props.ServicePrincipalID)
		d.Set("aad_service_principal_credential_type", props.AadServicePrincipalCredentialType)
		switch props.AadServicePrincipalCredentialType {
		case datafactory.ServicePrincipalCert:
			d.Set("service_principal_embedded_cert", props.ServicePrincipalEmbeddedCert)
			d.Set("service_principal_embedded_cert_password", props.ServicePrincipalEmbeddedCertPassword)
		case datafactory.ServicePrincipalKey:
			d.Set("service_principal_key", props.ServicePrincipalKey)
		default:
			return fmt.Errorf("Unsupported `aad_service_principal_credential_type`: %+v", props.AadServicePrincipalCredentialType)
		}
	}

	d.Set("additional_properties", odata.AdditionalProperties)
	d.Set("description", odata.Description)

	annotations := flattenDataFactoryAnnotations(odata.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(odata.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := odata.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceArmDataFactoryLinkedServiceODataDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Data Factory Linked Service OData %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
