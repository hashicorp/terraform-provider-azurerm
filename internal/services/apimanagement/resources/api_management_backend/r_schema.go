// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_backend

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementBackendSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the API Management backend.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApiManagementBackendName,
		},

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"circuit_breaker_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the circuit breaker rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$`),
							"`name` must be between 1 and 80 characters in length and may contain only numbers, letters, and hyphens (-) sign when preceded and followed by number or a letter.",
						),
					},
					"trip_duration": {
						Description:  "The duration for which the circuit remains open before retrying, in ISO 8601 format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.ISO8601Duration,
					},
					"accept_retry_after_enabled": {
						Description: "Whether the circuit breaker should honour `Retry-After` requests.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"failure_condition": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"interval_duration": {
									Description:  "The time window over which failures are counted, in ISO 8601 format.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.ISO8601Duration,
								},
								"count": {
									Description:  "The number of failures within the specified interval that will trigger the circuit breaker.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ExactlyOneOf: []string{"circuit_breaker_rule.0.failure_condition.0.count", "circuit_breaker_rule.0.failure_condition.0.percentage"},
									ValidateFunc: validation.IntBetween(1, 10000),
								},
								"error_reasons": {
									Description:  "A list of error reasons to consider as failures.",
									Type:         pluginsdk.TypeList,
									Optional:     true,
									MaxItems:     10,
									AtLeastOneOf: []string{"circuit_breaker_rule.0.failure_condition.0.status_code_range", "circuit_breaker_rule.0.failure_condition.0.error_reasons"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringLenBetween(1, 200),
									},
								},
								"percentage": {
									Description:  "The percentage of failures within the specified interval that will trigger the circuit breaker.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ExactlyOneOf: []string{"circuit_breaker_rule.0.failure_condition.0.count", "circuit_breaker_rule.0.failure_condition.0.percentage"},
									ValidateFunc: validation.IntBetween(1, 100),
								},
								"status_code_range": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									MaxItems:     10,
									AtLeastOneOf: []string{"circuit_breaker_rule.0.failure_condition.0.status_code_range", "circuit_breaker_rule.0.failure_condition.0.error_reasons"},
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"min": {
												Description:  "The minimum HTTP status code to consider as a failure.",
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntBetween(200, 599),
											},
											"max": {
												Description:  "The maximum HTTP status code to consider as a failure.",
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
									Description:  "The authentication Parameter value.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
									AtLeastOneOf: []string{"credentials.0.authorization.0.parameter", "credentials.0.authorization.0.scheme"},
								},
								"scheme": {
									Description:  "The authentication Scheme name.",
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
						Description: "A list of client certificate thumbprints to present to the backend host. The certificates must exist within the API Management Service.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{"credentials.0.authorization", "credentials.0.certificate", "credentials.0.header", "credentials.0.query"},
					},
					"header": {
						Description: "A mapping of header parameters to pass to the backend host. The keys are the header names and the values are a comma separated string of header values. This is converted to a list before being passed to the API.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{"credentials.0.authorization", "credentials.0.certificate", "credentials.0.header", "credentials.0.query"},
					},
					"query": {
						Description: "A mapping of query parameters to pass to the backend host. The keys are the query names and the values are a comma separated string of query values. This is converted to a list before being passed to the API.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
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
			Description:  "The description of the backend.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 2000),
		},

		"protocol": {
			Description:  "The protocol used by the backend host.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(backend.PossibleValuesForBackendProtocol(), false),
		},

		"proxy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"password": {
						Description:  "The password to connect to the proxy server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Description:  "The URL of the proxy server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"username": {
						Description:  "The username to connect to the proxy server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"resource_id": {
			Description:  "The management URI of the backend host in an external system. This URI can be the ARM Resource ID of Logic Apps, Function Apps or API Apps, or the management endpoint of a Service Fabric Cluster.",
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
						Description:  "The client certificate resource ID for the management endpoint.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.AsGeneratedID(certificate.ParseCertificateIDInsensitively),
					},

					"client_certificate_thumbprint": {
						Description:  "The client certificate thumbprint for the management endpoint.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"management_endpoints": {
						Description: "A list of cluster management endpoints.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"max_partition_resolution_retries": {
						Description: "The maximum number of retries when attempting resolve the partition.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},
					"server_certificate_thumbprints": {
						Description:   "A list of thumbprints of the server certificates of the Service Fabric Cluster.",
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
									Description:  "The thumbprint for the issuer of the certificate.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"name": {
									Description:  "The common name of the certificate.",
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
			Description:  "The title of the backend.",
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
						Description:  "Flag indicating whether SSL certificate chain validation should be done when using self-signed certificates for the backend host.",
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						AtLeastOneOf: []string{"tls.0.validate_certificate_chain", "tls.0.validate_certificate_name"},
					},
					"validate_certificate_name": {
						Description:  "Flag indicating whether SSL certificate name validation should be done when using self-signed certificates for the backend host.",
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						AtLeastOneOf: []string{"tls.0.validate_certificate_chain", "tls.0.validate_certificate_name"},
					},
				},
			},
		},

		"url": {
			Description:  "The backend host URL should be specified in the format `\"https://backend.com/api\"`, avoiding trailing slashes (/) to minimize misconfiguration risks. Azure API Management instance will append the backend resource name to this URL. This URL typically serves as the `base-url` in the [`set-backend-service`](https://learn.microsoft.com/azure/api-management/set-backend-service-policy) policy, enabling seamless transitions from frontend to backend.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
