// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceManagedCertificate() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceAppServiceManagedCertificateCreate,
		Read:   resourceAppServiceManagedCertificateRead,
		Update: resourceAppServiceManagedCertificateUpdate,
		Delete: resourceAppServiceManagedCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := certificates.ParseCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"custom_hostname_binding_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: webapps.ValidateHostNameBindingID,
			},

			"canonical_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subject_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"host_names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"issuer": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"issue_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"expiration_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FivePointOh() {
		// Parse insensitively for 4.x matching existing behaviour, enforce casing in 5.0
		r.Schema["custom_hostname_binding_id"].ValidateFunc = func(input interface{}, key string) (warnings []string, errors []error) {
			v, ok := input.(string)
			if !ok {
				errors = append(errors, fmt.Errorf("expected %q to be a string", key))
				return
			}

			if _, err := webapps.ParseHostNameBindingIDInsensitively(v); err != nil {
				errors = append(errors, err)
			}

			return
		}
	}

	return r
}

func resourceAppServiceManagedCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	appServiceClient := meta.(*clients.Client).Web.WebAppsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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

	existing, err := client.Get(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %w", id, err)
		}
		return tf.ImportAsExistsError("azurerm_app_service_managed_certificate", id.ID())
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

	// API may return a 202, however, the Location header returned does not return a ProvisioningState when polled
	// causing the provider to poll until timeout.
	if response.WasStatusCode(resp.HttpResponse, http.StatusAccepted) {
		poller := pollers.NewPoller(custompollers.NewAppServiceManagedCertificateCreatePoller(client, id), 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling %s: %w", id, err)
		}
	}

	d.SetId(id.ID())

	// An API issue prevents setting tags using the PUT operation, so we'll patch them in after
	// https://github.com/Azure/azure-rest-api-specs/issues/14529
	tags := certificates.CertificatePatchResource{
		Tags: utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, id, tags); err != nil {
		return fmt.Errorf("creating `tags` for %s: %w", id, err)
	}

	return resourceAppServiceManagedCertificateRead(d, meta)
}

func resourceAppServiceManagedCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("tags", model.Tags)
		if props := model.Properties; props != nil {
			d.Set("canonical_name", props.CanonicalName)
			d.Set("friendly_name", props.FriendlyName)
			d.Set("subject_name", props.SubjectName)
			d.Set("host_names", props.HostNames)
			d.Set("issuer", props.Issuer)
			d.Set("issue_date", props.IssueDate)
			d.Set("expiration_date", props.ExpirationDate)
			d.Set("thumbprint", props.Thumbprint)
		}
	}

	return nil
}

func resourceAppServiceManagedCertificateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Get(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	payload := certificates.CertificatePatchResource{
		Tags: utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %w", id, err)
	}

	return resourceAppServiceManagedCertificateRead(d, meta)
}

func resourceAppServiceManagedCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %w", id, err)
	}

	return nil
}
