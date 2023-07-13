// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubDPSCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSCertificateCreate,
		Read:   resourceIotHubDPSCertificateRead,
		Update: resourceIotHubDPSCertificateUpdate,
		Delete: resourceIotHubDPSCertificateDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dpscertificate.ParseCertificateID(id)
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
				ValidateFunc: validate.IoTHubDpsCertificateName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"iot_dps_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"certificate_content": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
			},

			"is_verified": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceIotHubDPSCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dpscertificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iot_dps_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, dpscertificate.GetOperationOptions{IfMatch: utils.String("")})
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing IoT Device Provisioning Service Certificate %s: %+v", id.String(), err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_iothub_dps_certificate", id.ID())
	}

	certificate := dpscertificate.CertificateResponse{
		Properties: &dpscertificate.CertificateProperties{
			Certificate: utils.String(d.Get("certificate_content").(string)),
		},
	}
	if d.Get("is_verified").(bool) {
		certificate.Properties.IsVerified = utils.Bool(true)
	}

	if _, err := client.CreateOrUpdate(ctx, id, certificate, dpscertificate.CreateOrUpdateOperationOptions{IfMatch: utils.String("")}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSCertificateRead(d, meta)
}

func resourceIotHubDPSCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dpscertificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, dpscertificate.GetOperationOptions{IfMatch: utils.String("")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("iot_dps_name", id.ProvisioningServiceName)
	// We are unable to set `certificate_content` since it is not returned from the API
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			isVerified := false
			if props.IsVerified != nil {
				isVerified = *props.IsVerified
			}
			d.Set("is_verified", isVerified)
		}
	}

	return nil
}

func resourceIotHubDPSCertificateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dpscertificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iot_dps_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, dpscertificate.GetOperationOptions{IfMatch: utils.String("")})
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	etag := ""
	if existing.Model != nil && existing.Model.Etag != nil {
		etag = *existing.Model.Etag
	}

	if d.HasChange("is_verified") {
		existing.Model.Properties.IsVerified = utils.Bool(d.Get("is_verified").(bool))
	}

	if d.HasChange("certificate_content") {
		existing.Model.Properties.Certificate = utils.String(d.Get("certificate_content").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, *existing.Model, dpscertificate.CreateOrUpdateOperationOptions{IfMatch: utils.String(etag)}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceIotHubDPSCertificateRead(d, meta)
}

func resourceIotHubDPSCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dpscertificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, dpscertificate.GetOperationOptions{IfMatch: utils.String("")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model != nil && resp.Model.Etag == nil {
		return fmt.Errorf("deleting %s because Etag is nil", id)
	}

	deleteOptions := dpscertificate.DeleteOperationOptions{
		IfMatch:         resp.Model.Etag,
		CertificateName: utils.String(id.CertificateName),
	}
	if _, err := client.Delete(ctx, *id, deleteOptions); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
