// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_public_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	appserviceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServicePublicCertificateDelete(d *pluginsdk.ResourceData, meta any) error {
	client := appserviceClient.From(meta).WebAppsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
