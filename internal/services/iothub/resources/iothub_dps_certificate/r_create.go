// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	iothubClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	client := iothubClient.From(meta).DPSCertificateClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := dpscertificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iot_dps_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, dpscertificate.GetOperationOptions{IfMatch: pointer.To("")})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing IoT Device Provisioning Service Certificate %s: %+v", id.String(), err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iothub_dps_certificate", id.ID())
		}
	}

	certificate := dpscertificate.CertificateResponse{
		Properties: &dpscertificate.CertificateProperties{
			Certificate: pointer.To(d.Get("certificate_content").(string)),
		},
	}
	if d.Get("is_verified").(bool) {
		certificate.Properties.IsVerified = pointer.To(true)
	}

	if _, err := client.CreateOrUpdate(ctx, id, certificate, dpscertificate.CreateOrUpdateOperationOptions{IfMatch: pointer.To("")}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSCertificateRead(d, meta)
}
