// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

func resourceAppServicePublicCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServicePublicCertificateCreate,
		Read:   resourceAppServicePublicCertificateRead,
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
			"resource_group_name": commonschema.ResourceGroupName(),

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
				ForceNew: true,
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

func resourceAppServicePublicCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPublicCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string), d.Get("certificate_name").(string))
	certificateLocation := d.Get("certificate_location").(string)
	blob := d.Get("blob").(string)

	existing, err := client.GetPublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_app_service_public_certificate", id.ID())
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

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("could not determine context deadline for create for %s", id)
	}

	// (@jackofallops) - The ok on the create call above can in some cases return before the resource is retrievable by
	// the `GetPublicCertificate` call, so we'll check it is actually created before progressing to read to prevent
	// false negative removal there.
	createWait := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound"},
		Target:                    []string{"ok"},
		MinTimeout:                10 * time.Second,
		Timeout:                   time.Until(deadline),
		NotFoundChecks:            10,
		ContinuousTargetOccurence: 3,
		Refresh: func() (interface{}, string, error) {
			resp, err := client.GetPublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return nil, "notfound", nil
				} else {
					return nil, "error", err
				}
			}
			return resp, "ok", nil
		},
	}

	if _, err := createWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %s", id, err)
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

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("could not determine context deadline for create for %s", id)
	}

	// (@mbfrahry) - similar to what @jackofallops noted above, the Get call sometimes does not return the public certificate so we'll do a get multiple times to confirm
	// that it's not there before removing the resource from state
	readWait := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound"},
		Target:                    []string{"ok"},
		MinTimeout:                10 * time.Second,
		Timeout:                   time.Until(deadline),
		NotFoundChecks:            10,
		ContinuousTargetOccurence: 1,
		Refresh: func() (interface{}, string, error) {
			resp, err := client.GetPublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return nil, "notfound", nil
				} else {
					return nil, "error", err
				}
			}
			return resp, "ok", nil
		},
	}

	resp, err := readWait.WaitForStateContext(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "couldn't find resource") {
			log.Printf("[DEBUG] App Service Public Certificate %q (Resource Group %q, App Service %q) was not found - removing from state", id.Name, id.ResourceGroup, id.SiteName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on App Service Public Certificate %q (Resource Group %q, App Service %q): %+v", id.Name, id.ResourceGroup, id.SiteName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("app_service_name", id.SiteName)
	d.Set("certificate_name", id.Name)

	if model, ok := resp.(web.PublicCertificate); ok {
		if properties := model.PublicCertificateProperties; properties != nil {
			d.Set("certificate_location", properties.PublicCertificateLocation)
			d.Set("blob", base64.StdEncoding.EncodeToString(*properties.Blob))
			d.Set("thumbprint", properties.Thumbprint)
		}
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
