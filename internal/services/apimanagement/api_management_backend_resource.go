// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend" // Explicitly uses the 2024-05-01 version so that the circuit breaker and load balancer pools functionality is available.
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azvalidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
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
				Description:  "Specifies the name of the API Management Backend. Changing this forces a new resource to be created.",
			},

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"credentials": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pool"},
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
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"pool"},
				ValidateFunc: validation.StringInSlice([]string{
					string(backend.BackendProtocolHTTP),
					string(backend.BackendProtocolSoap),
				}, false),
			},

			"proxy": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pool"},
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
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(1, 2000),
				ConflictsWith: []string{"pool"},
			},

			"service_fabric_cluster": {
				Type:          pluginsdk.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"pool"},
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
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pool"},
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
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"pool"},
			},

			"circuit_breaker_rule": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pool"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"accept_retry_after": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"trip_duration": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azvalidate.ISO8601Duration,
						},
						"failure_condition": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"count": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
										ExactlyOneOf: []string{"circuit_breaker_rule.0.failure_condition.0.count", "circuit_breaker_rule.0.failure_condition.0.percentage"},
									},
									"error_reasons": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"OperationNotFound",
												"SubscriptionKeyNotFound",
												"SubscriptionKeyInvalid",
												"ClientConnectionFailure",
												"BackendConnectionFailure",
												"ExpressionValueEvaluationFailure",
											}, false),
										},
									},
									"interval": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: azvalidate.ISO8601Duration,
									},
									"percentage": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 100),
										ExactlyOneOf: []string{"circuit_breaker_rule.0.failure_condition.0.count", "circuit_breaker_rule.0.failure_condition.0.percentage"},
									},
									"status_code_range": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"min": {
													Type:         pluginsdk.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(200, 599),
												},
												"max": {
													Type:         pluginsdk.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(200, 599),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"pool": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"circuit_breaker_rule", "url"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"service": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 30,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: backend.ValidateBackendID,
									},
									"priority": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"weight": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
					},
				},
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

	properties := new(backend.BackendContractProperties)

	// Pool type backends are very particular about what fields can and cannot be set
	if poolRaw, ok := d.GetOk("pool"); ok {
		properties.Type = pointer.To(backend.BackendTypePool)
		properties.Pool = expandApiManagementBackendPool(poolRaw.([]interface{}))
	} else {
		properties.Type = pointer.To(backend.BackendTypeSingle) // Set the type to Single if pool is not defined
		// Single type backends can have all the other fields set
		credentialsRaw := d.Get("credentials").([]interface{})
		properties.Credentials = expandApiManagementBackendCredentials(credentialsRaw)
		properties.Protocol = backend.BackendProtocol(d.Get("protocol").(string))
		proxyRaw := d.Get("proxy").([]interface{})
		properties.Proxy = expandApiManagementBackendProxy(proxyRaw)
		tlsRaw := d.Get("tls").([]interface{})
		properties.Tls = expandApiManagementBackendTls(tlsRaw)
		properties.Url = d.Get("url").(string)
		circuitBreakerRaw := d.Get("circuit_breaker_rule").([]interface{})
		properties.CircuitBreaker = expandApiManagementBackendCircuitBreaker(circuitBreakerRaw)

		if serviceFabricClusterRaw, ok := d.GetOk("service_fabric_cluster"); ok {
			err, serviceFabricCluster := expandApiManagementBackendServiceFabricCluster(serviceFabricClusterRaw.([]interface{}))
			if err != nil {
				return err
			}
			properties.Properties = &backend.BackendProperties{
				ServiceFabricCluster: serviceFabricCluster,
			}
		}
	}

	backendContract := backend.BackendContract{
		Properties: properties,
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

	// TODO, remove this debugging
	log.Printf("[DEBUG] sending properties: %+v", *backendContract.Properties)

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
			if err := d.Set("circuit_breaker_rule", flattenApiManagementBackendCircuitBreaker(props.CircuitBreaker)); err != nil {
				return fmt.Errorf("setting `circuit_breaker_rule`: %s", err)
			}
			if err := d.Set("pool", flattenApiManagementBackendPool(props.Pool)); err != nil {
				return fmt.Errorf("setting `pool`: %s", err)
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
		return errors.New("at least one of `client_certificate_thumbprint` and `client_certificate_id` must be set"), nil
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
		return errors.New("one of `server_certificate_thumbprints` or `server_x509_name` must be set"), nil
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

func expandApiManagementBackendCircuitBreaker(input []interface{}) *backend.BackendCircuitBreaker {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	circuitBreaker := backend.BackendCircuitBreaker{} // API & SDK have a circuit_breaker and a sub rules block. This is a lot of nesting so we "merge" the two levels into one object.
	rules := make([]backend.CircuitBreakerRule, 0)    // API requests a list so we need to provide a slice, even if it only contains one element.
	rule := backend.CircuitBreakerRule{}

	if name, ok := v["name"].(string); ok {
		rule.Name = pointer.To(name)
	}
	if acceptRetryAfter, ok := v["accept_retry_after"].(bool); ok {
		rule.AcceptRetryAfter = pointer.To(acceptRetryAfter)
	}
	if tripDuration, ok := v["trip_duration"].(string); ok {
		rule.TripDuration = pointer.To(tripDuration)
	}

	if failureConditionRaw := v["failure_condition"].([]interface{}); len(failureConditionRaw) > 0 {
		failureCondition := failureConditionRaw[0].(map[string]interface{})
		circuitBreakerFailureCondition := backend.CircuitBreakerFailureCondition{}

		if count, ok := failureCondition["count"].(int); ok && count != 0 {
			circuitBreakerFailureCondition.Count = pointer.To(int64(count))
		}

		if percentage, ok := failureCondition["percentage"].(int); ok && percentage != 0 {
			circuitBreakerFailureCondition.Percentage = pointer.To(int64(percentage))
		}

		if errorReasonsRaw := failureCondition["error_reasons"].([]interface{}); len(errorReasonsRaw) > 0 {
			errorReasons := make([]string, 0)
			for _, v := range errorReasonsRaw {
				errorReasons = append(errorReasons, v.(string))
			}
			circuitBreakerFailureCondition.ErrorReasons = &errorReasons
		}
		if interval, ok := failureCondition["interval"].(string); ok && interval != "" {
			circuitBreakerFailureCondition.Interval = pointer.To(interval)
		}

		if statusCodeRangesRaw := failureCondition["status_code_range"].([]interface{}); len(statusCodeRangesRaw) > 0 {
			statusCodeRanges := make([]backend.FailureStatusCodeRange, 0)
			for _, rangeRaw := range statusCodeRangesRaw {
				scRange := rangeRaw.(map[string]interface{})
				statusCodeRange := backend.FailureStatusCodeRange{}
				if min, ok := scRange["min"].(int); ok {
					statusCodeRange.Min = pointer.To(int64(min))
				}
				if max, ok := scRange["max"].(int); ok {
					statusCodeRange.Max = pointer.To(int64(max))
				}
				statusCodeRanges = append(statusCodeRanges, statusCodeRange)
			}
			circuitBreakerFailureCondition.StatusCodeRanges = &statusCodeRanges
		}
		rule.FailureCondition = &circuitBreakerFailureCondition
	}

	rules = append(rules, rule) // Add element to "list"
	circuitBreaker.Rules = &rules

	return &circuitBreaker
}

func expandApiManagementBackendPool(input []interface{}) *backend.BackendBaseParametersPool {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	pool := backend.BackendBaseParametersPool{}

	if servicesRaw := v["service"].([]interface{}); len(servicesRaw) > 0 {
		services := make([]backend.BackendPoolItem, 0)
		for _, serviceRaw := range servicesRaw {
			service := serviceRaw.(map[string]interface{})
			poolItem := backend.BackendPoolItem{
				Id: service["id"].(string),
			}

			if priority, ok := service["priority"].(int); ok {
				poolItem.Priority = pointer.To(int64(priority))
			}

			if weight, ok := service["weight"].(int); ok {
				poolItem.Weight = pointer.To(int64(weight))
			}

			services = append(services, poolItem)
		}
		pool.Services = &services
	}

	return &pool
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

func flattenApiManagementBackendCircuitBreaker(input *backend.BackendCircuitBreaker) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	if input.Rules == nil || len(*input.Rules) == 0 {
		return results
	}

	// We're only processing the first rule as the Terraform schema expects one circuit breaker rule
	rule := (*input.Rules)[0]
	result := make(map[string]interface{}) // API & SDK have a circuit_breaker and a sub rules block. This is a lot of nesting so we "merge" the two levels into one object.

	if rule.Name != nil {
		result["name"] = *rule.Name
	}

	if rule.AcceptRetryAfter != nil {
		result["accept_retry_after"] = *rule.AcceptRetryAfter
	}

	if rule.TripDuration != nil {
		result["trip_duration"] = *rule.TripDuration
	}

	if failureCondition := rule.FailureCondition; failureCondition != nil {
		failureConditionResult := make(map[string]interface{})

		if failureCondition.Count != nil {
			if *failureCondition.Count != 0 {
				failureConditionResult["count"] = int(*failureCondition.Count)
			}
		}

		if failureCondition.ErrorReasons != nil {
			failureConditionResult["error_reasons"] = *failureCondition.ErrorReasons
		}

		if failureCondition.Interval != nil {
			failureConditionResult["interval"] = *failureCondition.Interval
		}

		if failureCondition.Percentage != nil {
			if *failureCondition.Percentage != 0 {
				failureConditionResult["percentage"] = int(*failureCondition.Percentage)
			}
		}

		if statusCodeRanges := failureCondition.StatusCodeRanges; statusCodeRanges != nil {
			statusCodeRangesResult := make([]interface{}, 0)

			for _, statusCodeRange := range *statusCodeRanges {
				rangeResult := make(map[string]interface{})

				if statusCodeRange.Min != nil {
					rangeResult["min"] = int(*statusCodeRange.Min)
				}

				if statusCodeRange.Max != nil {
					rangeResult["max"] = int(*statusCodeRange.Max)
				}

				statusCodeRangesResult = append(statusCodeRangesResult, rangeResult)
			}

			failureConditionResult["status_code_range"] = statusCodeRangesResult
		}

		result["failure_condition"] = []interface{}{failureConditionResult}
	}

	return []interface{}{result}
}

func flattenApiManagementBackendPool(input *backend.BackendBaseParametersPool) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if services := input.Services; services != nil {
		servicesResult := make([]interface{}, 0)

		for _, service := range *services {
			serviceResult := make(map[string]interface{})

			serviceResult["id"] = service.Id

			if service.Priority != nil {
				serviceResult["priority"] = *service.Priority
			}

			if service.Weight != nil {
				serviceResult["weight"] = *service.Weight
			}

			servicesResult = append(servicesResult, serviceResult)
		}

		result["service"] = servicesResult
	}

	return []interface{}{result}
}
