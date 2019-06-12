package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementBackendCreateUpdate,
		Read:   resourceArmApiManagementBackendRead,
		Update: resourceArmApiManagementBackendCreateUpdate,
		Delete: resourceArmApiManagementBackendDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementBackendName,
			},

			"api_management_name": azure.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"credentials": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorization": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
									"scheme": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
						"certificate": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"header": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"query": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"properties": {},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.BackendProtocolHTTP),
					string(apimanagement.BackendProtocolSoap),
				}, false),
			},

			"proxy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"tls": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"validate_certificate_chain": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"validate_certificate_name": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmApiManagementBackendCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apimgmt.BackendClient
	ctx := meta.(*ArmClient).StopContext
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing backend %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_backend", *existing.ID)
		}
	}

	credentialsRaw := d.Get("credentials").([]interface{})
	credentials := expandApiManagementBackendCredentials(credentialsRaw)
	description := d.Get("description").(string)
	protocol := d.Get("protocol").(string)
	proxyRaw := d.Get("proxy").([]interface{})
	proxy := expandApiManagementBackendProxy(proxyRaw)
	resourceID := d.Get("resource_id").(string)
	title := d.Get("title").(string)
	tlsRaw := d.Get("tls").([]interface{})
	tls := expandApiManagementBackendTls(tlsRaw)
	url := d.Get("url").(string)

	backendContract := apimanagement.BackendContract{
		BackendContractProperties: &apimanagement.BackendContractProperties{
			Credentials: credentials,
			Description: utils.String(description),
			// Properties:  "",
			Protocol:   apimanagement.BackendProtocol(protocol),
			Proxy:      proxy,
			ResourceID: utils.String(resourceID),
			Title:      utils.String(title),
			TLS:        tls,
			URL:        utils.String(url),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, backendContract, ""); err != nil {
		return fmt.Errorf("Error creating/updating backend %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving backend %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for backend %q (API Management Service %q / Resource Group %q)", name, serviceName, resourceGroup)
	}

	d.SetId(*read.ID)
	return resourceArmApiManagementBackendRead(d, meta)
}

func resourceArmApiManagementBackendRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).apimgmt.BackendClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["backends"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] backend %q (API Management Service %q / Resource Group %q) does not exist - removing from state!", name, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving backend %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.BackendContractProperties; props != nil {
		d.Set("credentials", "")
		d.Set("description", props.Description)
		// d.Set("properties", "")
		d.Set("protocol", props.Protocol)
	}
}

func resourceArmApiManagementBackendDelete(d *schema.ResourceData, meta interface{}) error {
}

func expandApiManagementBackendCredentials(input []interface{}) *apimanagement.BackendCredentialsContract {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	authorizationRaw := v["authorization"].([]interface{})
	authorization := expandApiManagementBackendCredentialsAuthorization(authorizationRaw)
	certificate := v["certificate"].([]interface{})
	headerRaw := v["header"].(map[string]interface{})
	header := expandApiManagementBackendCredentialsObject(headerRaw)
	queryRaw := v["query"].(map[string]interface{})
	query := expandApiManagementBackendCredentialsObject(queryRaw)
	contract := apimanagement.BackendCredentialsContract{
		Authorization: authorization,
		Certificate:   utils.ExpandStringSlice(certificate),
		Header:        *header,
		Query:         *query,
	}
	return &contract
}

func expandApiManagementBackendCredentialsAuthorization(input []interface{}) *apimanagement.BackendAuthorizationHeaderCredentials {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	parameter := v["parameter"].(string)
	scheme := v["scheme"].(string)
	credentials := apimanagement.BackendAuthorizationHeaderCredentials{
		Parameter: utils.String(parameter),
		Scheme:    utils.String(scheme),
	}
	return &credentials
}

func expandApiManagementBackendCredentialsObject(input map[string]interface{}) *map[string][]string {
	output := make(map[string][]string)
	for k, v := range input {
		output[k] = v.([]string)
	}
	return &output
}

func expandApiManagementBackendProperties(input []interface{}) *apimanagement.BackendProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
}

func expandApiManagementBackendProxy(input []interface{}) *apimanagement.BackendProxyContract {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	password := v["password"].(string)
	url := v["url"].(string)
	username := v["username"].(string)
	contract := apimanagement.BackendProxyContract{
		Password: utils.String(password),
		URL:      utils.String(url),
		Username: utils.String(username),
	}
	return &contract
}

func expandApiManagementBackendTls(input []interface{}) *apimanagement.BackendTLSProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	validateCertificateChain := v["validate_certificate_chain"].(bool)
	validateCertificateName := v["validate_certificate_name"].(bool)
	properties := apimanagement.BackendTLSProperties{
		ValidateCertificateChain: utils.Bool(validateCertificateChain),
		ValidateCertificateName:  utils.Bool(validateCertificateName),
	}
	return &properties
}
