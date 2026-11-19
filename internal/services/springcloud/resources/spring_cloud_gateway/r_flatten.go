// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
)

func flattenGatewayGatewayAPIMetadataProperties(input *appplatform.GatewayApiMetadataProperties) []ApiMetadataModel {
	if input == nil {
		return make([]ApiMetadataModel, 0)
	}

	return []ApiMetadataModel{
		{
			Description:      pointer.From(input.Description),
			DocumentationUrl: pointer.From(input.Documentation),
			ServerURL:        pointer.From(input.ServerURL),
			Title:            pointer.From(input.Title),
			Version:          pointer.From(input.Version),
		},
	}
}

func flattenGatewayGatewayCorsProperties(input *appplatform.GatewayCorsProperties) []CorsModel {
	if input == nil {
		return make([]CorsModel, 0)
	}

	return []CorsModel{
		{
			CredentialsAllowed:    pointer.From(input.AllowCredentials),
			AllowedHeaders:        pointer.From(input.AllowedHeaders),
			AllowedMethods:        pointer.From(input.AllowedMethods),
			AllowedOrigins:        pointer.From(input.AllowedOrigins),
			AllowedOriginPatterns: pointer.From(input.AllowedOriginPatterns),
			ExposedHeaders:        pointer.From(input.ExposedHeaders),
			MaxAgeSeconds:         pointer.From(input.MaxAge),
		},
	}
}

func flattenGatewayGatewayResourceRequests(input *appplatform.GatewayResourceRequests) []QuotaModel {
	if input == nil {
		return make([]QuotaModel, 0)
	}
	return []QuotaModel{
		{
			Cpu:    pointer.From(input.Cpu),
			Memory: pointer.From(input.Memory),
		},
	}
}

func flattenGatewaySsoProperties(input *appplatform.SsoProperties, old []GatewaySsoModel) []GatewaySsoModel {
	if input == nil {
		return make([]GatewaySsoModel, 0)
	}

	oldItems := make(map[string]GatewaySsoModel)
	for _, item := range old {
		oldItems[item.IssuerUri] = item
	}

	issuerUri := pointer.From(input.IssuerUri)
	var clientId string
	var clientSecret string
	if oldItem, ok := oldItems[issuerUri]; ok {
		clientId = oldItem.ClientId
		clientSecret = oldItem.ClientSecret
	}
	return []GatewaySsoModel{
		{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			IssuerUri:    issuerUri,
			Scope:        pointer.From(input.Scope),
		},
	}
}

func flattenGatewayGatewayApmTypes(input *[]appplatform.ApmType) []string {
	if input == nil {
		return []string{}
	}
	out := make([]string, 0)
	for _, v := range *input {
		out = append(out, string(v))
	}
	return out
}

func flattenGatewayClientAuth(input *appplatform.GatewayPropertiesClientAuth) []ClientAuthorizationModel {
	if input == nil || input.Certificates == nil || len(*input.Certificates) == 0 {
		return make([]ClientAuthorizationModel, 0)
	}
	certificateIds := make([]string, 0)
	if input.Certificates != nil {
		for _, v := range *input.Certificates {
			certId, err := appplatform.ParseCertificateIDInsensitively(v)
			if err == nil {
				certificateIds = append(certificateIds, certId.ID())
			}
		}
	}
	verificationEnabled := input.CertificateVerification != nil && *input.CertificateVerification == appplatform.GatewayCertificateVerificationEnabled

	return []ClientAuthorizationModel{
		{
			CertificateIds:      certificateIds,
			VerificationEnabled: verificationEnabled,
		},
	}
}

func flattenGatewayLocalResponseCachePerRouteProperties(input appplatform.GatewayResponseCacheProperties) []ResponseCacheModel {
	if input == nil {
		return make([]ResponseCacheModel, 0)
	}
	if v, ok := input.(appplatform.GatewayLocalResponseCachePerRouteProperties); ok {
		return []ResponseCacheModel{
			{
				Size:       pointer.From(v.Size),
				TimeToLive: pointer.From(v.TimeToLive),
			},
		}
	}
	return make([]ResponseCacheModel, 0)
}

func flattenGatewayLocalResponseCachePerInstanceProperties(input appplatform.GatewayResponseCacheProperties) []ResponseCacheModel {
	if input == nil {
		return make([]ResponseCacheModel, 0)
	}
	if v, ok := input.(appplatform.GatewayLocalResponseCachePerInstanceProperties); ok {
		return []ResponseCacheModel{
			{
				Size:       pointer.From(v.Size),
				TimeToLive: pointer.From(v.TimeToLive),
			},
		}
	}
	return make([]ResponseCacheModel, 0)
}

func flattenGatewayApms(input *[]appplatform.ApmReference) ([]string, error) {
	out := make([]string, 0)
	if input == nil {
		return out, nil
	}
	for _, v := range *input {
		id, err := appplatform.ParseApmIDInsensitively(v.ResourceId)
		if err != nil {
			return nil, err
		}
		out = append(out, id.ID())
	}
	return out, nil
}
