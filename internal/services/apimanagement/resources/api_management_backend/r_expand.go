// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_backend

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandApiManagementBackendCredentials(input []any) *backend.BackendCredentialsContract {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	contract := backend.BackendCredentialsContract{}
	if authorizationRaw := v["authorization"]; authorizationRaw != nil {
		contract.Authorization = expandApiManagementBackendCredentialsAuthorization(authorizationRaw.([]any))
	}
	if certificate := v["certificate"]; certificate != nil {
		certificates := pluginsdk.ExpandStringSlice(certificate.([]any))
		if certificates != nil && len(*certificates) > 0 {
			contract.Certificate = certificates
		}
	}
	if headerRaw := v["header"]; headerRaw != nil {
		contract.Header = expandApiManagementBackendCredentialsObject(headerRaw.(map[string]any))
	}
	if queryRaw := v["query"]; queryRaw != nil {
		contract.Query = expandApiManagementBackendCredentialsObject(queryRaw.(map[string]any))
	}
	return &contract
}

func expandApiManagementBackendCredentialsAuthorization(input []any) *backend.BackendAuthorizationHeaderCredentials {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	credentials := backend.BackendAuthorizationHeaderCredentials{}
	if parameter := v["parameter"]; parameter != nil {
		credentials.Parameter = parameter.(string)
	}
	if scheme := v["scheme"]; scheme != nil {
		credentials.Scheme = scheme.(string)
	}
	return &credentials
}

func expandApiManagementBackendCredentialsObject(input map[string]any) *map[string][]string {
	output := make(map[string][]string)
	for k, v := range input {
		output[k] = strings.Split(v.(string), ",")
	}
	return &output
}

func expandApiManagementBackendProxy(input []any) *backend.BackendProxyContract {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
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

func expandApiManagementBackendServiceFabricCluster(input []any) (*backend.BackendServiceFabricClusterProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]any)
	managementEndpoints := v["management_endpoints"].(*pluginsdk.Set).List()
	maxPartitionResolutionRetries := int64(v["max_partition_resolution_retries"].(int))
	properties := backend.BackendServiceFabricClusterProperties{
		ManagementEndpoints:           pointer.From(pluginsdk.ExpandStringSlice(managementEndpoints)),
		MaxPartitionResolutionRetries: pointer.To(maxPartitionResolutionRetries),
	}

	if v2, ok := v["client_certificate_thumbprint"].(string); ok && v2 != "" {
		properties.ClientCertificatethumbprint = pointer.To(v2)
	}

	if v2, ok := v["client_certificate_id"].(string); ok && v2 != "" {
		properties.ClientCertificateId = pointer.To(v2)
	}

	if properties.ClientCertificateId == nil && properties.ClientCertificatethumbprint == nil {
		return nil, errors.New("at least one of `client_certificate_thumbprint` and `client_certificate_id` must be set")
	}

	serverCertificateThumbprintsUnset := true
	serverX509NamesUnset := true
	if serverCertificateThumbprints := v["server_certificate_thumbprints"]; serverCertificateThumbprints != nil {
		properties.ServerCertificateThumbprints = pluginsdk.ExpandStringSlice(serverCertificateThumbprints.(*pluginsdk.Set).List())
		serverCertificateThumbprintsUnset = false
	}
	if serverX509Names := v["server_x509_name"]; serverX509Names != nil {
		properties.ServerX509Names = expandApiManagementBackendServiceFabricClusterServerX509Names(serverX509Names.(*pluginsdk.Set).List())
		serverX509NamesUnset = false
	}
	if serverCertificateThumbprintsUnset && serverX509NamesUnset {
		return nil, errors.New("one of `server_certificate_thumbprints` or `server_x509_name` must be set")
	}
	return &properties, nil
}

func expandApiManagementBackendServiceFabricClusterServerX509Names(input []any) *[]backend.X509CertificateName {
	results := make([]backend.X509CertificateName, 0)
	for _, certificateName := range input {
		v := certificateName.(map[string]any)
		result := backend.X509CertificateName{
			IssuerCertificateThumbprint: pointer.To(v["issuer_certificate_thumbprint"].(string)),
			Name:                        pointer.To(v["name"].(string)),
		}
		results = append(results, result)
	}
	return &results
}

func expandApiManagementBackendTls(input []any) *backend.BackendTlsProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	properties := backend.BackendTlsProperties{}
	if validateCertificateChain := v["validate_certificate_chain"]; validateCertificateChain != nil {
		properties.ValidateCertificateChain = pointer.To(validateCertificateChain.(bool))
	}
	if validateCertificateName := v["validate_certificate_name"]; validateCertificateName != nil {
		properties.ValidateCertificateName = pointer.To(validateCertificateName.(bool))
	}
	return &properties
}

func expandApiManagementBackendCircuitBreaker(input []any) *backend.BackendCircuitBreaker {
	if len(input) == 0 {
		return nil
	}

	rules := make([]backend.CircuitBreakerRule, 0)

	v := input[0].(map[string]any)
	rule := backend.CircuitBreakerRule{
		Name:             pointer.To(v["name"].(string)),
		TripDuration:     pointer.To(v["trip_duration"].(string)),
		FailureCondition: expandApiManagementBackendCircuitBreakerFailureCondition(v["failure_condition"].([]any)),
	}

	if acceptRetryAfter, ok := v["accept_retry_after_enabled"]; ok {
		rule.AcceptRetryAfter = pointer.To(acceptRetryAfter.(bool))
	}

	rules = append(rules, rule)

	return &backend.BackendCircuitBreaker{
		Rules: &rules,
	}
}

func expandApiManagementBackendCircuitBreakerFailureCondition(input []any) *backend.CircuitBreakerFailureCondition {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	condition := backend.CircuitBreakerFailureCondition{
		Interval: pointer.To(v["interval_duration"].(string)),
	}

	if count, ok := v["count"]; ok && count.(int) > 0 {
		condition.Count = pointer.To(int64(count.(int)))
	}

	if percentage, ok := v["percentage"]; ok && percentage.(int) > 0 {
		condition.Percentage = pointer.To(int64(percentage.(int)))
	}

	if statusCodeRanges, ok := v["status_code_range"]; ok {
		ranges := statusCodeRanges.([]any)
		if len(ranges) > 0 {
			condition.StatusCodeRanges = expandApiManagementBackendCircuitBreakerStatusCodeRanges(ranges)
		}
	}

	if errorReasons, ok := v["error_reasons"]; ok {
		reasons := errorReasons.([]any)
		if len(reasons) > 0 {
			condition.ErrorReasons = pluginsdk.ExpandStringSlice(reasons)
		}
	}

	return &condition
}

func expandApiManagementBackendCircuitBreakerStatusCodeRanges(input []any) *[]backend.FailureStatusCodeRange {
	if len(input) == 0 {
		return nil
	}

	codeRanges := make([]backend.FailureStatusCodeRange, 0)
	for _, item := range input {
		v := item.(map[string]any)
		codeRange := backend.FailureStatusCodeRange{
			Max: pointer.To(int64(v["max"].(int))),
			Min: pointer.To(int64(v["min"].(int))),
		}
		codeRanges = append(codeRanges, codeRange)
	}

	return &codeRanges
}
