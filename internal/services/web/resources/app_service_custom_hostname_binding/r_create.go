// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_custom_hostname_binding

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCustomHostnameBindingCreate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := webapps.NewHostNameBindingID(account.From(meta).SubscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string), d.Get("hostname").(string))

	locks.ByName(id.SiteName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.SiteName, appServiceCustomHostnameBindingResourceName)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetHostNameBinding(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %w", id, err)
			}
			return tf.ImportAsExistsError("azurerm_app_service_custom_hostname_binding", id.ID())
		}
	}

	payload := webapps.HostNameBinding{
		Properties: &webapps.HostNameBindingProperties{
			SiteName: pointer.To(id.SiteName),
		},
	}

	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)
	if sslState != "" {
		if thumbprint == "" {
			return errors.New("`thumbprint` must be specified when `ssl_state` is set")
		}

		payload.Properties.SslState = pointer.ToEnum[webapps.SslState](sslState)
	}

	if thumbprint != "" {
		if sslState == "" {
			return errors.New("`ssl_state` must be specified when `thumbprint` is set")
		}

		payload.Properties.Thumbprint = pointer.To(thumbprint)
	}

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceAppServiceCustomHostnameBindingRead(d, meta)
}
