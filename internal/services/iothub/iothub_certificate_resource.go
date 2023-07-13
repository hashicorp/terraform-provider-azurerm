// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubCertificateCreate,
		Read:   resourceIotHubCertificateRead,
		Update: resourceIotHubCertificateUpdate,
		Delete: resourceIotHubCertificateDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubCertificateV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IotHubCertificateID(id)
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
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"iothub_name": {
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
				Default:  false,
			},
		},
	}
}

func resourceIotHubCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIotHubCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_iothub_certificate", id.ID())
	}

	certificate := devices.CertificateDescription{
		Properties: &devices.CertificateProperties{
			IsVerified:  utils.Bool(d.Get("is_verified").(bool)),
			Certificate: utils.String(d.Get("certificate_content").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName, certificate, ""); err != nil {
		return fmt.Errorf("creating  %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubCertificateRead(d, meta)
}

func resourceIotHubCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotHubCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("iothub_name", id.IotHubName)

	if props := resp.Properties; props != nil {
		d.Set("is_verified", props.IsVerified)
		d.Set("certificate_content", props.Certificate)
	}

	return nil
}

func resourceIotHubCertificateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIotHubCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		return fmt.Errorf("reading %s: %v", id, err)
	}

	etag := ""
	if existing.Etag != nil {
		etag = *existing.Etag
	}

	if d.HasChange("is_verified") {
		existing.Properties.IsVerified = utils.Bool(d.Get("is_verified").(bool))
	}

	if d.HasChange("certificate_content") {
		existing.Properties.Certificate = utils.String(d.Get("certificate_content").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName, existing, etag); err != nil {
		return fmt.Errorf("updating  %s: %+v", id, err)
	}

	return resourceIotHubCertificateRead(d, meta)
}

func resourceIotHubCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotHubCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Etag == nil {
		return fmt.Errorf("deleting  %s because Etag is nil", *id)
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName, *utils.String(*resp.Etag)); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}
