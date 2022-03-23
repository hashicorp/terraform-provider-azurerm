package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServicePublicCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServicePublicCertificateCreateUpdate,
		Read:   resourceAppServicePublicCertificateRead,
		Update: resourceAppServicePublicCertificateCreateUpdate,
		Delete: resourceAppServicePublicCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PublicCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"app_service_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
			},

			"certificate_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"certificate_location": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.PublicCertificateLocationLocalMachineMy),
					string(web.PublicCertificateLocationCurrentUserMy),
					string(web.PublicCertificateLocationUnknown),
				}, false),
			},

			"blob": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsBase64,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServicePublicCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPublicCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string), d.Get("certificate_name").(string))
	certificateLocation := d.Get("certificate_location").(string)
	blob := d.Get("blob").(string)

	if d.IsNewResource() {
		existing, err := client.GetPublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_app_service_public_certificate", id.ID())
		}
	}

	certificate := web.PublicCertificate{
		PublicCertificateProperties: &web.PublicCertificateProperties{
			PublicCertificateLocation: web.PublicCertificateLocation(certificateLocation),
		},
	}

	if blob != "" {
		decodedBlob, err := base64.StdEncoding.DecodeString(blob)
		if err != nil {
			return fmt.Errorf("could not decode blob: %+v", err)
		}
		certificate.PublicCertificateProperties.Blob = &decodedBlob
	}

	if _, err := client.CreateOrUpdatePublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name, certificate); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServicePublicCertificateRead(d, meta)
}

func resourceAppServicePublicCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PublicCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetPublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Public Certificate %q (Resource Group %q, App Service %q) was not found - removing from state", id.Name, id.ResourceGroup, id.SiteName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on App Service Public Certificate %q (Resource Group %q, App Service %q): %+v", id.Name, id.ResourceGroup, id.SiteName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("app_service_name", id.SiteName)
	d.Set("certificate_name", id.Name)

	if properties := resp.PublicCertificateProperties; properties != nil {
		d.Set("certificate_location", properties.PublicCertificateLocation)
		d.Set("blob", base64.StdEncoding.EncodeToString(*properties.Blob))
		d.Set("thumbprint", properties.Thumbprint)
	}

	return nil
}

func resourceAppServicePublicCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PublicCertificateID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Public Certificate %q (Resource Group %q, App Service %q)", id.Name, id.ResourceGroup, id.SiteName)

	resp, err := client.DeletePublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting App Service Public Certificate %q (Resource Group %q, App Servcice %q): %s)", id.Name, id.ResourceGroup, id.SiteName, err)
		}
	}

	return nil
}
