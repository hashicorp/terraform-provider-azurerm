// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_public_certificate

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	appserviceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServicePublicCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := appserviceClient.From(meta).WebAppsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := webapps.ParsePublicCertificateID(d.Id())
	if err != nil {
		return err
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("could not determine context deadline for Read for %s", id)
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
		Refresh: func() (any, string, error) {
			resp, err := client.GetPublicCertificate(ctx, *id)
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

	resp, err := readWait.WaitForStateContext(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "couldn't find resource") {
			log.Printf("[DEBUG] App Service Public Certificate %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on App Service Public Certificate %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("app_service_name", id.SiteName)
	d.Set("certificate_name", id.PublicCertificateName)

	if cert, ok := resp.(webapps.GetPublicCertificateOperationResponse); ok {
		if model := cert.Model; model != nil {
			if properties := model.Properties; properties != nil {
				d.Set("certificate_location", pointer.FromEnum(properties.PublicCertificateLocation))
				d.Set("blob", pointer.From(properties.Blob))
				d.Set("thumbprint", pointer.From(properties.Thumbprint))
			}
		}
	}

	return nil
}
