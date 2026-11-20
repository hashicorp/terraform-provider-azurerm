// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := appplatform.NewGatewayID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GatewaysGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			service, err := client.ServicesGet(ctx, *springId)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", springId, err)
			}
			if service.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", springId)
			}
			if service.Model.Sku == nil || service.Model.Sku.Name == nil || service.Model.Sku.Tier == nil {
				return fmt.Errorf("invalid `sku` for %s", springId)
			}

			gatewayResource := appplatform.GatewayResource{
				Properties: &appplatform.GatewayProperties{
					ClientAuth:              expandGatewayClientAuth(model.ClientAuthorization),
					ApiMetadataProperties:   expandGatewayGatewayAPIMetadataProperties(model.ApiMetadata),
					Apms:                    expandGatewayApms(model.ApplicationPerformanceMonitoringIds),
					ApmTypes:                expandGatewayGatewayApmTypes(model.ApplicationPerformanceMonitoringTypes),
					CorsProperties:          expandGatewayGatewayCorsProperties(model.Cors),
					EnvironmentVariables:    expandGatewayGatewayEnvironmentVariables(model.EnvironmentVariables, model.SensitiveEnvironmentVariables),
					HTTPSOnly:               pointer.To(model.HttpsOnly),
					Public:                  pointer.To(model.PublicNetworkAccessEnabled),
					ResponseCacheProperties: expandGatewayResponseCacheProperties(model),
					ResourceRequests:        expandGatewayGatewayResourceRequests(model.Quota),
					SsoProperties:           expandGatewaySsoProperties(model.Sso),
				},
				Sku: &appplatform.Sku{
					Name:     service.Model.Sku.Name,
					Tier:     service.Model.Sku.Tier,
					Capacity: pointer.To(model.InstanceCount),
				},
			}

			if err := client.GatewaysCreateOrUpdateCallbackThenPoll(ctx, id, gatewayResource, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandGatewayGatewayAPIMetadataProperties(input []ApiMetadataModel) *appplatform.GatewayApiMetadataProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.GatewayApiMetadataProperties{
		Title:         pointer.To(v.Title),
		Description:   pointer.To(v.Description),
		Documentation: pointer.To(v.DocumentationUrl),
		Version:       pointer.To(v.Version),
		ServerURL:     pointer.To(v.ServerURL),
	}
}

func expandGatewayGatewayCorsProperties(input []CorsModel) *appplatform.GatewayCorsProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.GatewayCorsProperties{
		AllowedOrigins:        pointer.To(v.AllowedOrigins),
		AllowedOriginPatterns: pointer.To(v.AllowedOriginPatterns),
		AllowedMethods:        pointer.To(v.AllowedMethods),
		AllowedHeaders:        pointer.To(v.AllowedHeaders),
		MaxAge:                pointer.To(v.MaxAgeSeconds),
		AllowCredentials:      pointer.To(v.CredentialsAllowed),
		ExposedHeaders:        pointer.To(v.ExposedHeaders),
	}
}

func expandGatewayGatewayResourceRequests(input []QuotaModel) *appplatform.GatewayResourceRequests {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.GatewayResourceRequests{
		Cpu:    pointer.To(v.Cpu),
		Memory: pointer.To(v.Memory),
	}
}

func expandGatewaySsoProperties(input []GatewaySsoModel) *appplatform.SsoProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.SsoProperties{
		Scope:        pointer.To(v.Scope),
		ClientId:     pointer.To(v.ClientId),
		ClientSecret: pointer.To(v.ClientSecret),
		IssuerUri:    pointer.To(v.IssuerUri),
	}
}

func expandGatewayGatewayApmTypes(input []string) *[]appplatform.ApmType {
	if len(input) == 0 {
		return nil
	}
	out := make([]appplatform.ApmType, 0)
	for _, v := range input {
		out = append(out, appplatform.ApmType(v))
	}
	return &out
}

func expandGatewayGatewayEnvironmentVariables(env map[string]string, secrets map[string]string) *appplatform.GatewayPropertiesEnvironmentVariables {
	if len(env) == 0 && len(secrets) == 0 {
		return nil
	}

	return &appplatform.GatewayPropertiesEnvironmentVariables{
		Properties: pointer.To(env),
		Secrets:    pointer.To(secrets),
	}
}

func expandGatewayClientAuth(input []ClientAuthorizationModel) *appplatform.GatewayPropertiesClientAuth {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	verificationEnabled := appplatform.GatewayCertificateVerificationDisabled
	if v.VerificationEnabled {
		verificationEnabled = appplatform.GatewayCertificateVerificationEnabled
	}
	return &appplatform.GatewayPropertiesClientAuth{
		Certificates:            pointer.To(v.CertificateIds),
		CertificateVerification: pointer.To(verificationEnabled),
	}
}

func expandGatewayResponseCacheProperties(input SpringCloudGatewayModel) appplatform.GatewayResponseCacheProperties {
	if len(input.LocalResponseCachePerRoute) != 0 {
		return appplatform.GatewayLocalResponseCachePerRouteProperties{
			Size:       pointer.To(input.LocalResponseCachePerRoute[0].Size),
			TimeToLive: pointer.To(input.LocalResponseCachePerRoute[0].TimeToLive),
		}
	}
	if len(input.LocalResponseCachePerInstance) != 0 {
		return appplatform.GatewayLocalResponseCachePerInstanceProperties{
			Size:       pointer.To(input.LocalResponseCachePerInstance[0].Size),
			TimeToLive: pointer.To(input.LocalResponseCachePerInstance[0].TimeToLive),
		}
	}
	return appplatform.RawGatewayResponseCachePropertiesImpl{}
}

func expandGatewayApms(input []string) *[]appplatform.ApmReference {
	if len(input) == 0 {
		return nil
	}
	out := make([]appplatform.ApmReference, 0)
	for _, v := range input {
		out = append(out, appplatform.ApmReference{
			ResourceId: v,
		})
	}
	return pointer.To(out)
}
