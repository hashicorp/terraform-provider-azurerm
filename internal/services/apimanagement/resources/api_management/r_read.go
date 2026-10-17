// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	signInClient := meta.(*clients.Client).ApiManagement.SignInClient
	signUpClient := meta.(*clients.Client).ApiManagement.SignUpClient
	delegationClient := meta.(*clients.Client).ApiManagement.DelegationSettingsClient
	tenantAccessClient := meta.(*clients.Client).ApiManagement.TenantAccessClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apimanagementservice.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	policyServiceId := policy.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
	policyClient := meta.(*clients.Client).ApiManagement.PolicyClient
	policy, err := policyClient.Get(ctx, policyServiceId, policy.GetOperationOptions{Format: pointer.To(policy.PolicyExportFormatXml)})
	if err != nil {
		if !response.WasNotFound(policy.HttpResponse) {
			return fmt.Errorf("retrieving Policy for %s: %+v", *id, err)
		}
	}

	d.Set("name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("publisher_email", model.Properties.PublisherEmail)
		d.Set("publisher_name", model.Properties.PublisherName)
		d.Set("notification_sender_email", pointer.From(model.Properties.NotificationSenderEmail))
		d.Set("gateway_url", pointer.From(model.Properties.GatewayURL))
		d.Set("gateway_regional_url", pointer.From(model.Properties.GatewayRegionalURL))
		d.Set("portal_url", pointer.From(model.Properties.PortalURL))
		d.Set("developer_portal_url", pointer.From(model.Properties.DeveloperPortalURL))
		d.Set("management_api_url", pointer.From(model.Properties.ManagementApiURL))
		d.Set("scm_url", pointer.From(model.Properties.ScmURL))
		d.Set("public_ip_addresses", pointer.From(model.Properties.PublicIPAddresses))
		d.Set("public_ip_address_id", pointer.From(model.Properties.PublicIPAddressId))
		d.Set("public_network_access_enabled", pointer.From(model.Properties.PublicNetworkAccess) == apimanagementservice.PublicNetworkAccessEnabled)
		d.Set("private_ip_addresses", pointer.From(model.Properties.PrivateIPAddresses))
		d.Set("virtual_network_type", pointer.From(model.Properties.VirtualNetworkType))
		d.Set("client_certificate_enabled", pointer.From(model.Properties.EnableClientCertificate))
		d.Set("gateway_disabled", pointer.From(model.Properties.DisableGateway))

		d.Set("certificate", flattenAPIManagementCertificates(d, model.Properties.Certificates))

		if model.Sku.Name != "" {
			if err := d.Set("security", flattenApiManagementSecurityCustomProperties(*model.Properties.CustomProperties, model.Sku.Name == apimanagementservice.SkuTypeConsumption)); err != nil {
				return fmt.Errorf("setting `security`: %+v", err)
			}
		}

		if err := d.Set("protocols", flattenApiManagementProtocolsCustomProperties(*model.Properties.CustomProperties)); err != nil {
			return fmt.Errorf("setting `protocols`: %+v", err)
		}

		if err := d.Set("hostname_configuration", flattenApiManagementHostnameConfigurations(model.Properties.HostnameConfigurations, d)); err != nil {
			return fmt.Errorf("setting `hostname_configuration`: %+v", err)
		}
		additionalLocation, err := flattenApiManagementAdditionalLocations(model.Properties.AdditionalLocations)
		if err != nil {
			return err
		}
		if err := d.Set("additional_location", additionalLocation); err != nil {
			return fmt.Errorf("setting `additional_location`: %+v", err)
		}

		virtualNetworkConfiguration, err := flattenApiManagementVirtualNetworkConfiguration(model.Properties.VirtualNetworkConfiguration)
		if err != nil {
			return err
		}
		if err := d.Set("virtual_network_configuration", virtualNetworkConfiguration); err != nil {
			return fmt.Errorf("setting `virtual_network_configuration`: %+v", err)
		}

		var minApiVersion string
		if model.Properties.ApiVersionConstraint != nil {
			minApiVersion = pointer.From(model.Properties.ApiVersionConstraint.MinApiVersion)
		}
		d.Set("min_api_version", minApiVersion)

		if err := d.Set("sku_name", flattenApiManagementServiceSkuName(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku_name`: %+v", err)
		}
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if model.Sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(model.Sku.Name), "V2") {
			signInSettingServiceId := signinsettings.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
			signInSettings, err := signInClient.Get(ctx, signInSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Sign In Settings for %s: %+v", *id, err)
			}
			if err := d.Set("sign_in", flattenApiManagementSignInSettings(*signInSettings.Model)); err != nil {
				return fmt.Errorf("setting `sign_in`: %+v", err)
			}

			signUpSettingServiceId := signupsettings.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
			signUpSettings, err := signUpClient.Get(ctx, signUpSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Sign Up Settings for %s: %+v", *id, err)
			}

			if err := d.Set("sign_up", flattenApiManagementSignUpSettings(*signUpSettings.Model)); err != nil {
				return fmt.Errorf("setting `sign_up`: %+v", err)
			}

			delegationSettingServiceId := delegationsettings.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
			delegationSettings, err := delegationClient.Get(ctx, delegationSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Delegation Settings for %s: %+v", *id, err)
			}

			delegationValidationKeyContract, err := delegationClient.ListSecrets(ctx, delegationSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Delegation Validation Key for %s: %+v", *id, err)
			}

			if err := d.Set("delegation", flattenApiManagementDelegationSettings(*delegationSettings.Model, *delegationValidationKeyContract.Model)); err != nil {
				return fmt.Errorf("setting `delegation`: %+v", err)
			}

			tenantAccessServiceId := tenantaccess.NewAccessID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, "access")
			tenantAccessInformationContract, err := tenantAccessClient.ListSecrets(ctx, tenantAccessServiceId)
			if err != nil {
				return fmt.Errorf("retrieving tenant access properties for %s: %+v", *id, err)
			}
			if err := d.Set("tenant_access", flattenApiManagementTenantAccessSettings(*tenantAccessInformationContract.Model)); err != nil {
				return fmt.Errorf("setting `tenant_access`: %+v", err)
			}
		} else {
			d.Set("sign_in", []any{})
			d.Set("sign_up", []any{})
			d.Set("delegation", []any{})
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

var (
	apimBackendProtocolSsl3                  = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30"
	apimBackendProtocolTls10                 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10"
	apimBackendProtocolTls11                 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11"
	apimFrontendProtocolSsl3                 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30"
	apimFrontendProtocolTls10                = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"
	apimFrontendProtocolTls11                = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"
	apimTripleDesCiphers                     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168"
	apimHttp2Protocol                        = "Microsoft.WindowsAzure.ApiManagement.Gateway.Protocols.Server.Http2"
	apimTlsEcdheEcdsaWithAes256CbcShaCiphers = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	apimTlsEcdheEcdsaWithAes128CbcShaCiphers = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	apimTlsEcdheRsaWithAes256CbcShaCiphers   = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	apimTlsEcdheRsaWithAes128CbcShaCiphers   = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	apimTlsRsaWithAes128GcmSha256Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_GCM_SHA256"
	apimTlsRsaWithAes256CbcSha256Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA256"
	apimTlsRsaWithAes256GcmSha384Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_GCM_SHA384"
	apimTlsRsaWithAes128CbcSha256Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA256"
	apimTlsRsaWithAes256CbcShaCiphers        = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA"
	apimTlsRsaWithAes128CbcShaCiphers        = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA"
)

func parseApiManagementNilableDictionary(input map[string]string, key string) bool {
	log.Printf("Parsing value for %q", key)

	v, ok := input[key]
	if !ok {
		log.Printf("%q was not found in the input - returning `false` as the default value", key)
		return false
	}

	val, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf(" parsing %q (key %q) as bool: %+v - assuming false", key, v, err)
		return false
	}

	return val
}

func optedOutOfRecoveringSoftDeletedApiManagementErrorFmt(name, location string) string {
	message := `
An existing soft-deleted API Management exists with the Name %q in the location %q, however
automatically recovering this API Management has been disabled via the "features" block.

Terraform can automatically recover the soft-deleted API Management when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import", or pick a different name/location.
`
	return fmt.Sprintf(message, name, location)
}
