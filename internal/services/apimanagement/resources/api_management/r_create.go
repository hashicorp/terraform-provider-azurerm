// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/deletedservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	apiClient := meta.(*clients.Client).ApiManagement.ApiClient
	deletedServicesClient := meta.(*clients.Client).ApiManagement.DeletedServicesClient
	productsClient := meta.(*clients.Client).ApiManagement.ProductsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sku := expandAzureRmApiManagementSkuName(d.Get("sku_name").(string))

	id := apimanagementservice.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)

	publicIpAddressId := d.Get("public_ip_address_id").(string)
	notificationSenderEmail := d.Get("notification_sender_email").(string)
	virtualNetworkType := d.Get("virtual_network_type").(string)

	customProperties, err := expandApiManagementCustomProperties(d, sku.Name == apimanagementservice.SkuTypeConsumption)
	if err != nil {
		return err
	}
	certificates := expandAzureRmApiManagementCertificates(d)

	publicNetworkAccess := apimanagementservice.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = apimanagementservice.PublicNetworkAccessDisabled
	}

	// before creating check to see if the resource exists in the soft delete state
	deletedServiceId := deletedservice.NewDeletedServiceID(id.SubscriptionId, location, id.ServiceName)
	softDeleted, err := deletedServicesClient.GetByName(ctx, deletedServiceId)
	if err != nil {
		// If Terraform lacks permission to read at the Subscription we'll get 403, not 404
		if !response.WasNotFound(softDeleted.HttpResponse) && !response.WasForbidden(softDeleted.HttpResponse) {
			return fmt.Errorf("checking for the presence of an existing Soft-Deleted API Management %q (Location %q): %+v", id.ServiceName, location, err)
		}
	}

	// if so, does the user want us to recover it?
	if !response.WasNotFound(softDeleted.HttpResponse) && !response.WasForbidden(softDeleted.HttpResponse) {
		if !meta.(*clients.Client).Features.ApiManagement.RecoverSoftDeleted {
			// this exists but the users opted out, so they must import this it out-of-band
			return errors.New(optedOutOfRecoveringSoftDeletedApiManagementErrorFmt(id.ServiceName, location))
		}

		// First recover the deleted API Management, since all other properties are ignored during a restore operation
		// (don't set the ID just yet to avoid tainting on failure)
		params := apimanagementservice.ApiManagementServiceResource{
			Location: location,
			Properties: apimanagementservice.ApiManagementServiceProperties{
				Restore: pointer.To(true),
			},
			Sku: sku,
		}

		// retry to restore service since there is an API issue : https://github.com/Azure/azure-rest-api-specs/issues/25262
		if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
			resp, err := client.CreateOrUpdate(ctx, id, params)
			if err != nil {
				if response.WasBadRequest(resp.HttpResponse) {
					return pluginsdk.RetryableError(err)
				}
				return pluginsdk.NonRetryableError(err)
			}
			if err := resp.Poller.PollUntilDone(ctx); err != nil {
				return pluginsdk.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("recovering %s: %+v", id, err)
		}
	}

	properties := apimanagementservice.ApiManagementServiceResource{
		Location: location,
		Properties: apimanagementservice.ApiManagementServiceProperties{
			PublisherName:       d.Get("publisher_name").(string),
			PublisherEmail:      d.Get("publisher_email").(string),
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
			CustomProperties:    pointer.To(customProperties),
			Certificates:        certificates,
		},
		Sku:  sku,
		Tags: tags.Expand(t),
	}

	if _, ok := d.GetOk("hostname_configuration"); ok {
		properties.Properties.HostnameConfigurations = expandAzureRmApiManagementHostnameConfigurations(d)
	}

	// intentionally not gated since we specify a default value (of None) in the expand, which we need on updates
	identityRaw := d.Get("identity").([]any)
	identity, err := identity.ExpandSystemAndUserAssignedMap(identityRaw)
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	properties.Identity = identity

	if _, ok := d.GetOk("additional_location"); ok {
		var err error
		properties.Properties.AdditionalLocations, err = expandAzureRmApiManagementAdditionalLocations(d, sku)
		if err != nil {
			return err
		}
	}

	if notificationSenderEmail != "" {
		properties.Properties.NotificationSenderEmail = pointer.To(notificationSenderEmail)
	}

	if virtualNetworkType != "" {
		properties.Properties.VirtualNetworkType = pointer.ToEnum[apimanagementservice.VirtualNetworkType](virtualNetworkType)

		if virtualNetworkType != string(apimanagementservice.VirtualNetworkTypeNone) {
			virtualNetworkConfiguration := expandAzureRmApiManagementVirtualNetworkConfigurations(d)
			if virtualNetworkConfiguration == nil {
				return fmt.Errorf("you must specify 'virtual_network_configuration' when 'virtual_network_type' is %q", virtualNetworkType)
			}
			properties.Properties.VirtualNetworkConfiguration = virtualNetworkConfiguration
		}
	}

	if publicIpAddressId != "" {
		if sku.Name != apimanagementservice.SkuTypePremium && sku.Name != apimanagementservice.SkuTypeDeveloper {
			if virtualNetworkType == string(apimanagementservice.VirtualNetworkTypeNone) {
				return errors.New("`public_ip_address_id` is only supported when sku type is `Developer` or `Premium`, and the APIM instance is deployed in a virtual network")
			}
		}
		properties.Properties.PublicIPAddressId = pointer.To(publicIpAddressId)
	}

	if d.HasChange("client_certificate_enabled") {
		enableClientCertificate := d.Get("client_certificate_enabled").(bool)
		if enableClientCertificate && sku.Name != apimanagementservice.SkuTypeConsumption {
			return errors.New("`client_certificate_enabled` is only supported when sku type is `Consumption`")
		}
		properties.Properties.EnableClientCertificate = pointer.To(enableClientCertificate)
	}

	gateWayDisabled := d.Get("gateway_disabled").(bool)
	if gateWayDisabled && len(*properties.Properties.AdditionalLocations) == 0 {
		return errors.New("`gateway_disabled` is only supported when `additional_location` is set")
	}
	properties.Properties.DisableGateway = pointer.To(gateWayDisabled)

	if v, ok := d.GetOk("min_api_version"); ok {
		properties.Properties.ApiVersionConstraint = &apimanagementservice.ApiVersionConstraint{
			MinApiVersion: pointer.To(v.(string)),
		}
	}

	if v := d.Get("zones").(*schema.Set).List(); len(v) > 0 {
		if sku.Name != apimanagementservice.SkuTypePremium {
			return errors.New("`zones` is only supported when sku type is `Premium`")
		}

		zones := zones.ExpandUntyped(v)
		properties.Zones = &zones
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// Remove sample products and APIs after creating (v3.0 behaviour)
	apiServiceId := api.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)

	listResp, err := apiClient.ListByService(ctx, apiServiceId, api.ListByServiceOperationOptions{})
	if err != nil {
		return fmt.Errorf("listing APIs after creation of %s: %+v", id, err)
	}
	if model := listResp.Model; model != nil {
		for _, contract := range *model {
			if contract.Id == nil {
				continue
			}
			apiId, err := api.ParseApiID(pointer.From(contract.Id))
			if err != nil {
				return fmt.Errorf("parsing API ID: %+v", err)
			}
			if delResp, err := apiClient.Delete(ctx, *apiId, api.DeleteOperationOptions{DeleteRevisions: pointer.To(true)}); err != nil {
				if !response.WasNotFound(delResp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *apiId, err)
				}
			}
		}
	}

	produceServiceId := product.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
	proListResp, err := productsClient.ListByService(ctx, produceServiceId, product.ListByServiceOperationOptions{})
	if err != nil {
		return fmt.Errorf("listing products after creation of %s: %+v", id, err)
	}
	if model := proListResp.Model; model != nil {
		for _, contract := range *model {
			if contract.Id == nil {
				continue
			}
			productId, err := product.ParseProductID(pointer.From(contract.Id))
			if err != nil {
				return fmt.Errorf("parsing product ID: %+v", err)
			}
			if delResp, err := productsClient.Delete(ctx, *productId, product.DeleteOperationOptions{DeleteSubscriptions: pointer.To(true)}); err != nil {
				if !response.WasNotFound(delResp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *productId, err)
				}
			}
		}
	}

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

	signUpSettingsRaw := d.Get("sign_up").([]any)
	if (sku.Name == apimanagementservice.SkuTypeConsumption || strings.Contains(string(sku.Name), "V2")) && len(signUpSettingsRaw) > 0 {
		return fmt.Errorf("`sign_up` is not supported for sku tiers `Consumption` and `V2`")
	}
	if sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(sku.Name), "V2") {
		signUpSettingServiceId := signupsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
		signUpSettings := expandApiManagementSignUpSettings(signUpSettingsRaw)
		signUpClient := meta.(*clients.Client).ApiManagement.SignUpClient
		if _, err := signUpClient.CreateOrUpdate(ctx, signUpSettingServiceId, signUpSettings, signupsettings.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf(" setting Sign Up settings for %s: %+v", id, err)
		}
	}

	delegationSettingsRaw := d.Get("delegation").([]any)
	if (sku.Name == apimanagementservice.SkuTypeConsumption || strings.Contains(string(sku.Name), "V2")) && len(delegationSettingsRaw) > 0 {
		return fmt.Errorf("`delegation` is not supported for sku tiers `Consumption` and `V2`")
	}
	if sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(sku.Name), "V2") && len(delegationSettingsRaw) > 0 {
		delegationSettingServiceId := delegationsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
		delegationSettings := expandApiManagementDelegationSettings(delegationSettingsRaw)
		delegationClient := meta.(*clients.Client).ApiManagement.DelegationSettingsClient
		if _, err := delegationClient.CreateOrUpdate(ctx, delegationSettingServiceId, delegationSettings, delegationsettings.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf(" setting Delegation settings for %s: %+v", id, err)
		}
	}

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

	return resourceApiManagementServiceRead(d, meta)
}
