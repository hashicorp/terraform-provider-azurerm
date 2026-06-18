// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceAppServiceCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAppServiceCertificateRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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
}

func dataSourceAppServiceCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := certificates.NewCertificateID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("tags", model.Tags)

		if props := model.Properties; props != nil {
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
