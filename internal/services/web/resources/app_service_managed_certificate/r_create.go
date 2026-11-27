// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_managed_certificate

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceManagedCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).CertificatesClient
	appServiceClient := webClient.From(meta).WebAppsClient
	subscriptionID := account.From(meta).SubscriptionId

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	chbID, err := webapps.ParseHostNameBindingID(d.Get("custom_hostname_binding_id").(string))
	if err != nil {
		return err
	}

	appServiceID := commonids.NewAppServiceID(subscriptionID, chbID.ResourceGroupName, chbID.SiteName)
	appService, err := appServiceClient.Get(ctx, appServiceID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", appServiceID, err)
	}

	if appService.Model == nil || appService.Model.Properties == nil || appService.Model.Properties.ServerFarmId == nil {
		return fmt.Errorf("retrieving `serverFarmId` from %s", appServiceID)
	}

	appServicePlanID, err := commonids.ParseAppServicePlanIDInsensitively(*appService.Model.Properties.ServerFarmId)
	if err != nil {
		return err
	}

	id := certificates.NewCertificateID(subscriptionID, appServicePlanID.ResourceGroupName, chbID.HostNameBindingName)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %w", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_app_service_managed_certificate", id.ID())
		}
	}

	certificate := certificates.Certificate{
		Properties: &certificates.CertificateProperties{
			CanonicalName: pointer.To(chbID.HostNameBindingName),
			ServerFarmId:  pointer.To(appServicePlanID.ID()),
			Password:      new(string),
		},
		Location: location.Normalize(appService.Model.Location),
	}

	resp, err := client.CreateOrUpdate(ctx, id, certificate)
	if err != nil {
		return fmt.Errorf("creating %s: %w", id, err)
	}

	d.SetId(id.ID())

	// API may return a 202, however, the Location header returned does not return a ProvisioningState when polled
	// causing the provider to poll until timeout.
	if response.WasStatusCode(resp.HttpResponse, http.StatusAccepted) {
		poller := pollers.NewPoller(custompollers.NewAppServiceManagedCertificateCreatePoller(client, id), 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling %s: %w", id, err)
		}
	}

	// An API issue prevents setting tags using the PUT operation, so we'll patch them in after
	// https://github.com/Azure/azure-rest-api-specs/issues/14529
	t := certificates.CertificatePatchResource{
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.Update(ctx, id, t); err != nil {
		return fmt.Errorf("creating `tags` for %s: %w", id, err)
	}

	return resourceAppServiceManagedCertificateRead(d, meta)
}
