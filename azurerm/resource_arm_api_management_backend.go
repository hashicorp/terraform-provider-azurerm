package azurerm

import (
	"fmt"
	"log"
	"strings"

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
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"query": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 2000),
			},

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
				ValidateFunc: validation.StringLenBetween(1, 2000),
			},

			"service_fabric_cluster": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_certificate_thumbprint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"management_endpoints": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"max_partition_resolution_retries": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"server_certificate_thumbprints": {
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"service_fabric_cluster.0.server_x509_name"},
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"server_x509_name": {
							Type:          schema.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"service_fabric_cluster.0.server_certificate_thumbprints"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issuer_certificate_thumbprint": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
					},
				},
			},

			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 300),
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
	client := meta.(*ArmClient).apiManagement.BackendClient
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
	protocol := d.Get("protocol").(string)
	proxyRaw := d.Get("proxy").([]interface{})
	proxy := expandApiManagementBackendProxy(proxyRaw)
	tlsRaw := d.Get("tls").([]interface{})
	tls := expandApiManagementBackendTls(tlsRaw)
	url := d.Get("url").(string)

	backendContract := apimanagement.BackendContract{
		BackendContractProperties: &apimanagement.BackendContractProperties{
			Credentials: credentials,
			Protocol:    apimanagement.BackendProtocol(protocol),
			Proxy:       proxy,
			TLS:         tls,
			URL:         utils.String(url),
		},
	}
	if description, ok := d.GetOk("description"); ok {
		backendContract.BackendContractProperties.Description = utils.String(description.(string))
	}
	if resourceID, ok := d.GetOk("resource_id"); ok {
		backendContract.BackendContractProperties.ResourceID = utils.String(resourceID.(string))
	}
	if title, ok := d.GetOk("title"); ok {
		backendContract.BackendContractProperties.Title = utils.String(title.(string))
	}

	if serviceFabricClusterRaw, ok := d.GetOk("service_fabric_cluster"); ok {
		err, serviceFabricCluster := expandApiManagementBackendServiceFabricCluster(serviceFabricClusterRaw.([]interface{}))
		if err != nil {
			return err
		}
		backendContract.BackendContractProperties.Properties = &apimanagement.BackendProperties{
			ServiceFabricCluster: serviceFabricCluster,
		}
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
	client := meta.(*ArmClient).apiManagement.BackendClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
		d.Set("description", props.Description)
		d.Set("protocol", string(props.Protocol))
		d.Set("resource_id", props.ResourceID)
		d.Set("title", props.Title)
		d.Set("url", props.URL)
		if err := d.Set("credentials", flattenApiManagementBackendCredentials(props.Credentials)); err != nil {
			return fmt.Errorf("Error setting `credentials`: %s", err)
		}
		if err := d.Set("proxy", flattenApiManagementBackendProxy(props.Proxy)); err != nil {
			return fmt.Errorf("Error setting `proxy`: %s", err)
		}
		if properties := props.Properties; properties != nil {
			if err := d.Set("service_fabric_cluster", flattenApiManagementBackendServiceFabricCluster(properties.ServiceFabricCluster)); err != nil {
				return fmt.Errorf("Error setting `service_fabric_cluster`: %s", err)
			}
		}
		if err := d.Set("tls", flattenApiManagementBackendTls(props.TLS)); err != nil {
			return fmt.Errorf("Error setting `tls`: %s", err)
		}
	}

	return nil
}

func resourceArmApiManagementBackendDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.BackendClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["backends"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting backend %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
		}
	}

	return nil
}

func expandApiManagementBackendCredentials(input []interface{}) *apimanagement.BackendCredentialsContract {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	contract := apimanagement.BackendCredentialsContract{}
	if authorizationRaw := v["authorization"]; authorizationRaw != nil {
		authorization := expandApiManagementBackendCredentialsAuthorization(authorizationRaw.([]interface{}))
		contract.Authorization = authorization
	}
	if certificate := v["certificate"]; certificate != nil {
		certificates := utils.ExpandStringSlice(certificate.([]interface{}))
		if certificates != nil && len(*certificates) > 0 {
			contract.Certificate = certificates
		}
	}
	if headerRaw := v["header"]; headerRaw != nil {
		header := expandApiManagementBackendCredentialsObject(headerRaw.(map[string]interface{}))
		contract.Header = *header
	}
	if queryRaw := v["query"]; queryRaw != nil {
		query := expandApiManagementBackendCredentialsObject(queryRaw.(map[string]interface{}))
		contract.Query = *query
	}
	return &contract
}

func expandApiManagementBackendCredentialsAuthorization(input []interface{}) *apimanagement.BackendAuthorizationHeaderCredentials {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	credentials := apimanagement.BackendAuthorizationHeaderCredentials{}
	if parameter := v["parameter"]; parameter != nil {
		credentials.Parameter = utils.String(parameter.(string))
	}
	if scheme := v["scheme"]; scheme != nil {
		credentials.Scheme = utils.String(scheme.(string))
	}
	return &credentials
}

func expandApiManagementBackendCredentialsObject(input map[string]interface{}) *map[string][]string {
	output := make(map[string][]string)
	for k, v := range input {
		output[k] = strings.Split(v.(string), ",")
	}
	return &output
}

func expandApiManagementBackendProxy(input []interface{}) *apimanagement.BackendProxyContract {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	contract := apimanagement.BackendProxyContract{}
	if password := v["password"]; password != nil {
		contract.Password = utils.String(password.(string))
	}
	if url := v["url"]; url != nil {
		contract.URL = utils.String(url.(string))
	}
	if username := v["username"]; username != nil {
		contract.Username = utils.String(username.(string))
	}
	return &contract
}

func expandApiManagementBackendServiceFabricCluster(input []interface{}) (error, *apimanagement.BackendServiceFabricClusterProperties) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})
	clientCertificatethumbprint := v["client_certificate_thumbprint"].(string)
	managementEndpoints := v["management_endpoints"].(*schema.Set).List()
	maxPartitionResolutionRetries := int32(v["max_partition_resolution_retries"].(int))
	properties := apimanagement.BackendServiceFabricClusterProperties{
		ClientCertificatethumbprint:   utils.String(clientCertificatethumbprint),
		ManagementEndpoints:           utils.ExpandStringSlice(managementEndpoints),
		MaxPartitionResolutionRetries: utils.Int32(maxPartitionResolutionRetries),
	}
	serverCertificateThumbprintsUnset := true
	serverX509NamesUnset := true
	if serverCertificateThumbprints := v["server_certificate_thumbprints"]; serverCertificateThumbprints != nil {
		properties.ServerCertificateThumbprints = utils.ExpandStringSlice(serverCertificateThumbprints.(*schema.Set).List())
		serverCertificateThumbprintsUnset = false
	}
	if serverX509Names := v["server_x509_name"]; serverX509Names != nil {
		properties.ServerX509Names = expandApiManagementBackendServiceFabricClusterServerX509Names(serverX509Names.(*schema.Set).List())
		serverX509NamesUnset = false
	}
	if serverCertificateThumbprintsUnset && serverX509NamesUnset {
		return fmt.Errorf("One of `server_certificate_thumbprints` or `server_x509_name` must be set"), nil
	}
	return nil, &properties
}

func expandApiManagementBackendServiceFabricClusterServerX509Names(input []interface{}) *[]apimanagement.X509CertificateName {
	results := make([]apimanagement.X509CertificateName, 0)
	for _, certificateName := range input {
		v := certificateName.(map[string]interface{})
		result := apimanagement.X509CertificateName{
			IssuerCertificateThumbprint: utils.String(v["issuer_certificate_thumbprint"].(string)),
			Name:                        utils.String(v["name"].(string)),
		}
		results = append(results, result)
	}
	return &results
}

func expandApiManagementBackendTls(input []interface{}) *apimanagement.BackendTLSProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	properties := apimanagement.BackendTLSProperties{}
	if validateCertificateChain := v["validate_certificate_chain"]; validateCertificateChain != nil {
		properties.ValidateCertificateChain = utils.Bool(validateCertificateChain.(bool))
	}
	if validateCertificateName := v["validate_certificate_name"]; validateCertificateName != nil {
		properties.ValidateCertificateName = utils.Bool(validateCertificateName.(bool))
	}
	return &properties
}

func flattenApiManagementBackendCredentials(input *apimanagement.BackendCredentialsContract) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	result["authorization"] = flattenApiManagementBackendCredentialsAuthorization(input.Authorization)
	if input.Certificate != nil {
		result["certificate"] = *input.Certificate
	}
	result["header"] = flattenApiManagementBackendCredentialsObject(input.Header)
	result["query"] = flattenApiManagementBackendCredentialsObject(input.Query)
	return append(results, result)
}

func flattenApiManagementBackendCredentialsObject(input map[string][]string) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}
	for k, v := range input {
		results[k] = strings.Join(v, ",")
	}
	return results
}

func flattenApiManagementBackendCredentialsAuthorization(input *apimanagement.BackendAuthorizationHeaderCredentials) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if parameter := input.Parameter; parameter != nil {
		result["parameter"] = *parameter
	}
	if scheme := input.Scheme; scheme != nil {
		result["scheme"] = *scheme
	}
	return append(results, result)
}

func flattenApiManagementBackendProxy(input *apimanagement.BackendProxyContract) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if password := input.Password; password != nil {
		result["password"] = *password
	}
	if url := input.URL; url != nil {
		result["url"] = *url
	}
	if username := input.Username; username != nil {
		result["username"] = *username
	}
	return append(results, result)
}

func flattenApiManagementBackendServiceFabricCluster(input *apimanagement.BackendServiceFabricClusterProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if clientCertificatethumbprint := input.ClientCertificatethumbprint; clientCertificatethumbprint != nil {
		result["client_certificate_thumbprint"] = *clientCertificatethumbprint
	}
	if managementEndpoints := input.ManagementEndpoints; managementEndpoints != nil {
		result["management_endpoints"] = *managementEndpoints
	}
	if maxPartitionResolutionRetries := input.MaxPartitionResolutionRetries; maxPartitionResolutionRetries != nil {
		result["max_partition_resolution_retries"] = int(*maxPartitionResolutionRetries)
	}
	if serverCertificateThumbprints := input.ServerCertificateThumbprints; serverCertificateThumbprints != nil {
		result["server_certificate_thumbprints"] = *serverCertificateThumbprints
	}
	result["server_x509_name"] = flattenApiManagementBackendServiceFabricClusterServerX509Names(input.ServerX509Names)
	return append(results, result)
}

func flattenApiManagementBackendServiceFabricClusterServerX509Names(input *[]apimanagement.X509CertificateName) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	for _, certificateName := range *input {
		result := make(map[string]interface{})
		if issuerCertificateThumbprint := certificateName.IssuerCertificateThumbprint; issuerCertificateThumbprint != nil {
			result["issuer_certificate_thumbprint"] = *issuerCertificateThumbprint
		}
		if name := certificateName.Name; name != nil {
			result["name"] = *name
		}
		results = append(results, result)
	}
	return results
}

func flattenApiManagementBackendTls(input *apimanagement.BackendTLSProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if validateCertificateChain := input.ValidateCertificateChain; validateCertificateChain != nil {
		result["validate_certificate_chain"] = *validateCertificateChain
	}
	if validateCertificateName := input.ValidateCertificateName; validateCertificateName != nil {
		result["validate_certificate_name"] = *validateCertificateName
	}
	return append(results, result)
}
