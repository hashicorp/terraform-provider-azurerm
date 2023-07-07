// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBatchCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBatchCertificateRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CertificateName,
			},

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"public_data": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"format": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"thumbprint_algorithm": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBatchCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := certificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
	}
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			format := ""
			if v := props.Format; v != nil {
				format = string(*v)
			}
			d.Set("format", format)

			publicData := ""
			if v := props.PublicData; v != nil {
				publicData = *v
			}
			d.Set("public_data", publicData)

			thumbprint := ""
			if v := props.Thumbprint; v != nil {
				thumbprint = *v
			}
			d.Set("thumbprint", thumbprint)

			thumbprintAlgorithm := ""
			if v := props.ThumbprintAlgorithm; v != nil {
				thumbprintAlgorithm = *v
			}
			d.Set("thumbprint_algorithm", thumbprintAlgorithm)
		}
	}

	return nil
}
