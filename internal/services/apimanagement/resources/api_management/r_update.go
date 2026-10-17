// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sku := expandAzureRmApiManagementSkuName(d.Get("sku_name").(string))
	virtualNetworkType := d.Get("virtual_network_type").(string)
	virtualNetworkConfiguration := expandAzureRmApiManagementVirtualNetworkConfigurations(d)

	id := apimanagementservice.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if _, err := client.Get(ctx, id); err != nil {
		return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
	}

	props := apimanagementservice.ApiManagementServiceUpdateProperties{}
	payload := apimanagementservice.ApiManagementServiceUpdateParameters{}

	if d.HasChange("sku_name") {
		payload.Sku = pointer.To(sku)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("public_ip_address_id") {
		publicIpAddressId := d.Get("public_ip_address_id").(string)
		if publicIpAddressId != "" {
			if sku.Name != apimanagementservice.SkuTypePremium && sku.Name != apimanagementservice.SkuTypeDeveloper {
				if d.Get("virtual_network_type").(string) == string(apimanagementservice.VirtualNetworkTypeNone) {
					return fmt.Errorf("`public_ip_address_id` is only supported when sku type is `Developer` or `Premium`, and the APIM instance is deployed in a virtual network")
				}
			}
			props.PublicIPAddressId = pointer.To(publicIpAddressId)
		}
	}

	if d.HasChange("notification_sender_email") {
		props.NotificationSenderEmail = pointer.To(d.Get("notification_sender_email").(string))
	}

	if d.HasChange("virtual_network_type") {
		props.VirtualNetworkType = pointer.ToEnum[apimanagementservice.VirtualNetworkType](virtualNetworkType)
		if virtualNetworkType != string(apimanagementservice.VirtualNetworkTypeNone) {
			if virtualNetworkConfiguration == nil {
				return fmt.Errorf("you must specify 'virtual_network_configuration' when 'virtual_network_type' is %q", virtualNetworkType)
			}
			props.VirtualNetworkConfiguration = virtualNetworkConfiguration
		}
	}

	if d.HasChange("virtual_network_configuration") {
		props.VirtualNetworkConfiguration = virtualNetworkConfiguration
		if virtualNetworkType == string(apimanagementservice.VirtualNetworkTypeNone) {
			if virtualNetworkConfiguration != nil {
				return fmt.Errorf("you must specify 'virtual_network_type' when specifying 'virtual_network_configuration'")
			}
		}
	}

	if d.HasChanges("security", "protocols") {
		customProperties, err := expandApiManagementCustomProperties(d, sku.Name == apimanagementservice.SkuTypeConsumption)
		if err != nil {
			return err
		}
		props.CustomProperties = pointer.To(customProperties)
	}

	if d.HasChange("certificate") {
		props.Certificates = expandAzureRmApiManagementCertificates(d)
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := apimanagementservice.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = apimanagementservice.PublicNetworkAccessDisabled
		}

		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("publisher_name") {
		props.PublisherName = pointer.To(d.Get("publisher_name").(string))
	}

	if d.HasChange("publisher_email") {
		props.PublisherEmail = pointer.To(d.Get("publisher_email").(string))
	}

	if d.HasChange("hostname_configuration") {
		props.HostnameConfigurations = expandAzureRmApiManagementHostnameConfigurations(d)
	}

	// intentionally not gated since we specify a default value (of None) in the expand, which we need on updates
	if d.HasChange("identity") {
		identityRaw := d.Get("identity").([]any)
		identity, err := identity.ExpandSystemAndUserAssignedMap(identityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = identity
	}

	if d.HasChange("additional_location") {
		var err error
		props.AdditionalLocations, err = expandAzureRmApiManagementAdditionalLocations(d, sku)
		if err != nil {
			return err
		}
	}

	if d.HasChange("client_certificate_enabled") {
		enableClientCertificate := d.Get("client_certificate_enabled").(bool)
		if enableClientCertificate && sku.Name != apimanagementservice.SkuTypeConsumption {
			return errors.New("`client_certificate_enabled` is only supported when sku type is `Consumption`")
		}
		props.EnableClientCertificate = pointer.To(enableClientCertificate)
	}

	if d.HasChange("gateway_disabled") {
		gateWayDisabled := d.Get("gateway_disabled").(bool)
		if gateWayDisabled && props.AdditionalLocations != nil && len(*props.AdditionalLocations) == 0 {
			return errors.New("`gateway_disabled` is only supported when `additional_location` is set")
		}
		props.DisableGateway = pointer.To(gateWayDisabled)
	}

	if d.HasChange("min_api_version") {
		props.ApiVersionConstraint = &apimanagementservice.ApiVersionConstraint{}

		if v, ok := d.GetOk("min_api_version"); ok {
			props.ApiVersionConstraint.MinApiVersion = pointer.To(v.(string))
		}
	}

	if d.HasChange("zones") {
		if v := d.Get("zones").(*schema.Set).List(); len(v) > 0 {
			if sku.Name != apimanagementservice.SkuTypePremium {
				return errors.New("`zones` is only supported when sku type is `Premium`")
			}

			zones := zones.ExpandUntyped(v)
			payload.Zones = &zones
		}
	}

	payload.Properties = pointer.To(props)

	if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if d.HasChange("sign_in") {
		signInSettingsRaw := d.Get("sign_in").([]any)
		if (sku.Name == apimanagementservice.SkuTypeConsumption || strings.Contains(string(sku.Name), "V2")) && len(signInSettingsRaw) > 0 {
			return errors.New("`sign_in` is not supported for sku tiers `Consumption` and `V2`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(sku.Name), "V2") {
			signInSettingServiceId := signinsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
			signInSettings := expandApiManagementSignInSettings(signInSettingsRaw)
			signInClient := meta.(*clients.Client).ApiManagement.SignInClient
			if _, err := signInClient.CreateOrUpdate(ctx, signInSettingServiceId, signInSettings, signinsettings.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" setting Sign In settings for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("sign_up") {
		signUpSettingsRaw := d.Get("sign_up").([]any)
		if (sku.Name == apimanagementservice.SkuTypeConsumption || strings.Contains(string(sku.Name), "V2")) && len(signUpSettingsRaw) > 0 {
			return errors.New("`sign_up` is not supported for sku tiers `Consumption` and `V2`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(sku.Name), "V2") {
			signUpSettingServiceId := signupsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
			signUpSettings := expandApiManagementSignUpSettings(signUpSettingsRaw)
			signUpClient := meta.(*clients.Client).ApiManagement.SignUpClient
			if _, err := signUpClient.CreateOrUpdate(ctx, signUpSettingServiceId, signUpSettings, signupsettings.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" setting Sign Up settings for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("delegation") {
		delegationSettingsRaw := d.Get("delegation").([]any)
		if (sku.Name == apimanagementservice.SkuTypeConsumption || strings.Contains(string(sku.Name), "V2")) && len(delegationSettingsRaw) > 0 {
			return errors.New("`delegation` is not supported for sku tiers `Consumption` and `V2`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(sku.Name), "V2") && len(delegationSettingsRaw) > 0 {
			delegationSettingServiceId := delegationsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
			delegationSettings := expandApiManagementDelegationSettings(delegationSettingsRaw)
			delegationClient := meta.(*clients.Client).ApiManagement.DelegationSettingsClient
			if _, err := delegationClient.CreateOrUpdate(ctx, delegationSettingServiceId, delegationSettings, delegationsettings.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" setting Delegation settings for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("tenant_access") {
		tenantAccessRaw := d.Get("tenant_access").([]any)
		if (sku.Name == apimanagementservice.SkuTypeConsumption || strings.Contains(string(sku.Name), "V2")) && len(tenantAccessRaw) > 0 {
			return fmt.Errorf("`tenant_access` is not supported for sku tiers `Consumption` and `V2`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(sku.Name), "V2") && d.HasChange("tenant_access") {
			tenantAccessServiceId := tenantaccess.NewAccessID(subscriptionId, id.ResourceGroupName, id.ServiceName, "access")
			tenantAccessInformationParametersRaw := d.Get("tenant_access").([]any)
			tenantAccessInformationParameters := expandApiManagementTenantAccessSettings(tenantAccessInformationParametersRaw)
			tenantAccessClient := meta.(*clients.Client).ApiManagement.TenantAccessClient
			if _, err := tenantAccessClient.Update(ctx, tenantAccessServiceId, tenantAccessInformationParameters, tenantaccess.UpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" updating tenant access settings for %s: %+v", id, err)
			}
		}
	}

	return resourceApiManagementServiceRead(d, meta)
}
