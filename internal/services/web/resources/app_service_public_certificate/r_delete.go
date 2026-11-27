// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_public_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServicePublicCertificateDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapps.ParsePublicCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeletePublicCertificate(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting App Service Public Certificate %s: %+v)", *id, err)
		}
	}

	return nil
}
