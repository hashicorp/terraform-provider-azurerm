// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceApiManagementApiVersionSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementApiVersionSetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildDataSourceName(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versioning_scheme": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version_header_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version_query_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementApiVersionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	newId := apiversionset.NewApiVersionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, newId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", newId, err)
		}

		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	if resp.Model != nil && (resp.Model.Id == nil || *resp.Model.Id == "") {
		return fmt.Errorf("retrieving API Version Set %q (API Management Service %q /Resource Group %q): ID was nil or empty", newId.VersionSetId, newId.ServiceName, newId.ResourceGroupName)
	}

	id, err := apiversionset.ParseApiVersionSetID(*resp.Model.Id)
	if err != nil {
		return fmt.Errorf("parsing API Version Set ID: %q", *id)
	}

	d.Set("name", id.VersionSetId)
	d.SetId(id.ID())

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("versioning_scheme", string(props.VersioningScheme))
			d.Set("version_header_name", pointer.From(props.VersionHeaderName))
			d.Set("version_query_name", pointer.From(props.VersionQueryName))
		}
	}

	return nil
}
