// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/backend"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementBackend() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementBackendCreateUpdate,
		Read:   resourceApiManagementBackendRead,
		Update: resourceApiManagementBackendCreateUpdate,
		Delete: resourceApiManagementBackendDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := backend.ParseBackendID(id)
			return err
		}),

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
				ValidateFunc: validate.ApiManagementBackendName,
			},

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"credentials": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authorization": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"parameter": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										AtLeastOneOf: []string{"credentials.0.authorization.0.parameter", "credentials.0.authorization.0.scheme"},
									},
									"scheme": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										AtLeastOneOf: []string{"credentials.0.authorization.0.parameter", "credentials.0.authorization.0.scheme"},
									},
								},
							},
							AtLeastOneOf: []string{"credentials.0.authorization", "credentials.0.certificate", "credentials.0.header", "credentials.0.query"},
						},
						"certificate": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{"credentials.0.authorization", "credentials.0.certificate", "credentials.0.header", "credentials.0.query"},
						},
						"header": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{"credentials.0.authorization", "credentials.0.certificate", "credentials.0.header", "credentials.0.query"},
						},
						"query": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{"credentials.0.authorization", "credentials.0.certificate", "credentials.0.header", "credentials.0.query"},
						},
					},
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 2000),
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(backend.BackendProtocolHTTP),
					string(backend.BackendProtocolSoap),
				}, false),
			},

			"proxy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"password": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"username": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 2000),
			},

			"service_fabric_cluster": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_certificate_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.CertificateID,
						},

						"client_certificate_thumbprint": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"management_endpoints": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"max_partition_resolution_retries": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"server_certificate_thumbprints": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"service_fabric_cluster.0.server_x509_name"},
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"server_x509_name": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"service_fabric_cluster.0.server_certificate_thumbprints"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"issuer_certificate_thumbprint": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},

			"title": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 300),
			},

			"tls": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"validate_certificate_chain": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							AtLeastOneOf: []string{"tls.0.validate_certificate_chain", "tls.0.validate_certificate_name"},
						},
						"validate_certificate_name": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							AtLeastOneOf: []string{"tls.0.validate_certificate_chain", "tls.0.validate_certificate_name"},
						},
					},
				},
			},

			"url": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementBackendCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.BackendClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := backend.NewBackendID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_backend", id.ID())
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

	backendContract := backend.BackendContract{
		Properties: &backend.BackendContractProperties{
			Credentials: credentials,
			Protocol:    backend.BackendProtocol(protocol),
			Proxy:       proxy,
			Tls:         tls,
			Url:         url,
		},
	}
	if description, ok := d.GetOk("description"); ok {
		backendContract.Properties.Description = pointer.To(description.(string))
	}
	if resourceID, ok := d.GetOk("resource_id"); ok {
		backendContract.Properties.ResourceId = pointer.To(resourceID.(string))
	}
	if title, ok := d.GetOk("title"); ok {
		backendContract.Properties.Title = pointer.To(title.(string))
	}

	if serviceFabricClusterRaw, ok := d.GetOk("service_fabric_cluster"); ok {
		err, serviceFabricCluster := expandApiManagementBackendServiceFabricCluster(serviceFabricClusterRaw.([]interface{}))
		if err != nil {
			return err
		}
		backendContract.Properties.Properties = &backend.BackendProperties{
			ServiceFabricCluster: serviceFabricCluster,
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, backendContract, backend.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementBackendRead(d, meta)
}

func resourceApiManagementBackendRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.BackendClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := backend.ParseBackendID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("protocol", string(props.Protocol))
			d.Set("resource_id", pointer.From(props.ResourceId))
			d.Set("title", pointer.From(props.Title))
			d.Set("url", props.Url)
			if err := d.Set("credentials", flattenApiManagementBackendCredentials(props.Credentials)); err != nil {
				return fmt.Errorf("setting `credentials`: %s", err)
			}
			if err := d.Set("proxy", flattenApiManagementBackendProxy(props.Proxy)); err != nil {
				return fmt.Errorf("setting `proxy`: %s", err)
			}
			if properties := props.Properties; properties != nil {
				if err := d.Set("service_fabric_cluster", flattenApiManagementBackendServiceFabricCluster(properties.ServiceFabricCluster)); err != nil {
					return fmt.Errorf("setting `service_fabric_cluster`: %s", err)
				}
			}
			if err := d.Set("tls", flattenApiManagementBackendTls(props.Tls)); err != nil {
				return fmt.Errorf("setting `tls`: %s", err)
			}
		}
	}

	return nil
}

func resourceApiManagementBackendDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.BackendClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backend.ParseBackendID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, backend.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}

func expandApiManagementBackendCredentials(input []interface{}) *backend.BackendCredentialsContract {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	contract := backend.BackendCredentialsContract{}
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
		contract.Header = header
	}
	if queryRaw := v["query"]; queryRaw != nil {
		query := expandApiManagementBackendCredentialsObject(queryRaw.(map[string]interface{}))
		contract.Query = query
	}
	return &contract
}

func expandApiManagementBackendCredentialsAuthorization(input []interface{}) *backend.BackendAuthorizationHeaderCredentials {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	credentials := backend.BackendAuthorizationHeaderCredentials{}
	if parameter := v["parameter"]; parameter != nil {
		credentials.Parameter = parameter.(string)
	}
	if scheme := v["scheme"]; scheme != nil {
		credentials.Scheme = scheme.(string)
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

func expandApiManagementBackendProxy(input []interface{}) *backend.BackendProxyContract {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	contract := backend.BackendProxyContract{}
	if password := v["password"]; password != nil {
		contract.Password = pointer.To(password.(string))
	}
	if url := v["url"]; url != nil {
		contract.Url = url.(string)
	}
	if username := v["username"]; username != nil {
		contract.Username = pointer.To(username.(string))
	}
	return &contract
}

func expandApiManagementBackendServiceFabricCluster(input []interface{}) (error, *backend.BackendServiceFabricClusterProperties) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})
	managementEndpoints := v["management_endpoints"].(*pluginsdk.Set).List()
	maxPartitionResolutionRetries := int64(v["max_partition_resolution_retries"].(int))
	properties := backend.BackendServiceFabricClusterProperties{
		ManagementEndpoints:           pointer.From(utils.ExpandStringSlice(managementEndpoints)),
		MaxPartitionResolutionRetries: pointer.To(maxPartitionResolutionRetries),
	}

	if v2, ok := v["client_certificate_thumbprint"].(string); ok && v2 != "" {
		properties.ClientCertificatethumbprint = pointer.To(v2)
	}

	if v2, ok := v["client_certificate_id"].(string); ok && v2 != "" {
		properties.ClientCertificateId = pointer.To(v2)
	}

	if properties.ClientCertificateId == nil && properties.ClientCertificatethumbprint == nil {
		return fmt.Errorf("at least one of `client_certificate_thumbprint` and `client_certificate_id` must be set"), nil
	}

	serverCertificateThumbprintsUnset := true
	serverX509NamesUnset := true
	if serverCertificateThumbprints := v["server_certificate_thumbprints"]; serverCertificateThumbprints != nil {
		properties.ServerCertificateThumbprints = utils.ExpandStringSlice(serverCertificateThumbprints.(*pluginsdk.Set).List())
		serverCertificateThumbprintsUnset = false
	}
	if serverX509Names := v["server_x509_name"]; serverX509Names != nil {
		properties.ServerX509Names = expandApiManagementBackendServiceFabricClusterServerX509Names(serverX509Names.(*pluginsdk.Set).List())
		serverX509NamesUnset = false
	}
	if serverCertificateThumbprintsUnset && serverX509NamesUnset {
		return fmt.Errorf("One of `server_certificate_thumbprints` or `server_x509_name` must be set"), nil
	}
	return nil, &properties
}

func expandApiManagementBackendServiceFabricClusterServerX509Names(input []interface{}) *[]backend.X509CertificateName {
	results := make([]backend.X509CertificateName, 0)
	for _, certificateName := range input {
		v := certificateName.(map[string]interface{})
		result := backend.X509CertificateName{
			IssuerCertificateThumbprint: pointer.To(v["issuer_certificate_thumbprint"].(string)),
			Name:                        pointer.To(v["name"].(string)),
		}
		results = append(results, result)
	}
	return &results
}

func expandApiManagementBackendTls(input []interface{}) *backend.BackendTlsProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	properties := backend.BackendTlsProperties{}
	if validateCertificateChain := v["validate_certificate_chain"]; validateCertificateChain != nil {
		properties.ValidateCertificateChain = pointer.To(validateCertificateChain.(bool))
	}
	if validateCertificateName := v["validate_certificate_name"]; validateCertificateName != nil {
		properties.ValidateCertificateName = pointer.To(validateCertificateName.(bool))
	}
	return &properties
}

func flattenApiManagementBackendCredentials(input *backend.BackendCredentialsContract) []interface{} {
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

func flattenApiManagementBackendCredentialsObject(input *map[string][]string) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}
	for k, v := range *input {
		results[k] = strings.Join(v, ",")
	}
	return results
}

func flattenApiManagementBackendCredentialsAuthorization(input *backend.BackendAuthorizationHeaderCredentials) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if parameter := input.Parameter; parameter != "" {
		result["parameter"] = parameter
	}
	if scheme := input.Scheme; scheme != "" {
		result["scheme"] = scheme
	}
	return append(results, result)
}

func flattenApiManagementBackendProxy(input *backend.BackendProxyContract) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if password := input.Password; password != nil {
		result["password"] = *password
	}
	if url := input.Url; url != "" {
		result["url"] = url
	}
	result["username"] = pointer.From(input.Username)
	return append(results, result)
}

func flattenApiManagementBackendServiceFabricCluster(input *backend.BackendServiceFabricClusterProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	if clientCertificatethumbprint := input.ClientCertificatethumbprint; clientCertificatethumbprint != nil {
		result["client_certificate_thumbprint"] = *clientCertificatethumbprint
	}

	if input.ClientCertificateId != nil {
		result["client_certificate_id"] = *input.ClientCertificateId
	}

	if managementEndpoints := input.ManagementEndpoints; managementEndpoints != nil {
		result["management_endpoints"] = managementEndpoints
	}
	if maxPartitionResolutionRetries := input.MaxPartitionResolutionRetries; maxPartitionResolutionRetries != nil {
		result["max_partition_resolution_retries"] = int(pointer.From(input.MaxPartitionResolutionRetries))
	}
	result["server_certificate_thumbprints"] = pointer.From(input.ServerCertificateThumbprints)
	result["server_x509_name"] = flattenApiManagementBackendServiceFabricClusterServerX509Names(input.ServerX509Names)
	return append(results, result)
}

func flattenApiManagementBackendServiceFabricClusterServerX509Names(input *[]backend.X509CertificateName) []interface{} {
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

func flattenApiManagementBackendTls(input *backend.BackendTlsProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	result := make(map[string]interface{})
	result["validate_certificate_chain"] = pointer.From(input.ValidateCertificateChain)
	result["validate_certificate_name"] = pointer.From(input.ValidateCertificateName)
	return append(results, result)
}
