// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_public_certificate

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServicePublicCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapps.NewPublicCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string), d.Get("certificate_name").(string))
	certificateLocation := d.Get("certificate_location").(string)
	blob := d.Get("blob").(string)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetPublicCertificate(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_app_service_public_certificate", id.ID())
		}
	}

	certificate := webapps.PublicCertificate{
		Properties: &webapps.PublicCertificateProperties{
			PublicCertificateLocation: pointer.ToEnum[webapps.PublicCertificateLocation](certificateLocation),
		},
	}

	if blob != "" {
		certificate.Properties.Blob = pointer.To(blob)
	}

	if _, err := client.CreateOrUpdatePublicCertificate(ctx, id, certificate); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	d.SetId(id.ID())

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
		Refresh: func() (any, string, error) {
			resp, err := client.GetPublicCertificate(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
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

	return resourceAppServicePublicCertificateRead(d, meta)
}
