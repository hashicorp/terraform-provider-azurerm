// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceTemplateSpecVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceTemplateSpecVersionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		//lintignore:S033
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"version": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecVersionName,
			},

			"template_body": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceTemplateSpecVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecsVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := templatespecversions.NewTemplateSpecVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string), d.Get("version").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		templateBodyRaw, err := flattenTemplateDeploymentBody(model.Properties.MainTemplate)
		if err != nil {
			return err
		}
		d.Set("template_body", templateBodyRaw)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
