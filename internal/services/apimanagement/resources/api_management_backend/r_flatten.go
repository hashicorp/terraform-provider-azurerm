// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_backend

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
)

func flattenApiManagementBackendCircuitBreaker(input *backend.BackendCircuitBreaker) []any {
	results := make([]any, 0)
	if input == nil || input.Rules == nil {
		return results
	}

	for _, rule := range *input.Rules {
		result := make(map[string]any)
		result["name"] = pointer.From(rule.Name)
		result["trip_duration"] = pointer.From(rule.TripDuration)
		result["accept_retry_after_enabled"] = pointer.From(rule.AcceptRetryAfter)
		result["failure_condition"] = flattenApiManagementBackendCircuitBreakerFailureCondition(rule.FailureCondition)
		results = append(results, result)
	}

	return results
}

func flattenApiManagementBackendCircuitBreakerStatusCodeRanges(input *[]backend.FailureStatusCodeRange) []any {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		result := make(map[string]any)
		result["min"] = pointer.From(item.Min)
		result["max"] = pointer.From(item.Max)
		results = append(results, result)
	}

	return results
}

func flattenApiManagementBackendCircuitBreakerFailureCondition(input *backend.CircuitBreakerFailureCondition) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	result := make(map[string]any)

	result["count"] = pointer.From(input.Count)
	result["percentage"] = pointer.From(input.Percentage)
	result["interval_duration"] = pointer.From(input.Interval)
	result["status_code_range"] = flattenApiManagementBackendCircuitBreakerStatusCodeRanges(input.StatusCodeRanges)
	result["error_reasons"] = pointer.From(input.ErrorReasons)

	return append(results, result)
}

func flattenApiManagementBackendCredentials(input *backend.BackendCredentialsContract) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	result := make(map[string]any)
	result["authorization"] = flattenApiManagementBackendCredentialsAuthorization(input.Authorization)
	if input.Certificate != nil {
		result["certificate"] = *input.Certificate
	}
	result["header"] = flattenApiManagementBackendCredentialsObject(input.Header)
	result["query"] = flattenApiManagementBackendCredentialsObject(input.Query)
	return append(results, result)
}

func flattenApiManagementBackendCredentialsObject(input *map[string][]string) map[string]any {
	results := make(map[string]any)
	if input == nil {
		return results
	}
	for k, v := range *input {
		results[k] = strings.Join(v, ",")
	}
	return results
}

func flattenApiManagementBackendCredentialsAuthorization(input *backend.BackendAuthorizationHeaderCredentials) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	result := make(map[string]any)
	if parameter := input.Parameter; parameter != "" {
		result["parameter"] = parameter
	}
	if scheme := input.Scheme; scheme != "" {
		result["scheme"] = scheme
	}
	return append(results, result)
}

func flattenApiManagementBackendProxy(input *backend.BackendProxyContract) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	result := make(map[string]any)
	if password := input.Password; password != nil {
		result["password"] = *password
	}
	if url := input.Url; url != "" {
		result["url"] = url
	}
	result["username"] = pointer.From(input.Username)
	return append(results, result)
}

func flattenApiManagementBackendServiceFabricCluster(input *backend.BackendServiceFabricClusterProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	result := make(map[string]any)
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

func flattenApiManagementBackendServiceFabricClusterServerX509Names(input *[]backend.X509CertificateName) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	for _, certificateName := range *input {
		result := make(map[string]any)
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

func flattenApiManagementBackendTls(input *backend.BackendTlsProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	result := make(map[string]any)
	result["validate_certificate_chain"] = pointer.From(input.ValidateCertificateChain)
	result["validate_certificate_name"] = pointer.From(input.ValidateCertificateName)
	return append(results, result)
}
