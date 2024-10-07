// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiVersionSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiVersionSetCreateUpdate,
		Read:   resourceApiManagementApiVersionSetRead,
		Update: resourceApiManagementApiVersionSetCreateUpdate,
		Delete: resourceApiManagementApiVersionSetDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apiversionset.ParseApiVersionSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiVersionSetV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"versioning_scheme": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apiversionset.VersioningSchemeHeader),
					string(apiversionset.VersioningSchemeQuery),
					string(apiversionset.VersioningSchemeSegment),
				}, false),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"version_header_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"version_query_name"},
			},

			"version_query_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"version_header_name"},
			},
		},
	}
}

func resourceApiManagementApiVersionSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apiversionset.NewApiVersionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_version_set", id.ID())
		}
	}

	versioningScheme := apiversionset.VersioningScheme(d.Get("versioning_scheme").(string))
	parameters := apiversionset.ApiVersionSetContract{
		Properties: &apiversionset.ApiVersionSetContractProperties{
			DisplayName:      d.Get("display_name").(string),
			VersioningScheme: versioningScheme,
			Description:      pointer.To(d.Get("description").(string)),
		},
	}

	var headerSet, querySet bool
	if v, ok := d.GetOk("version_header_name"); ok {
		headerSet = v.(string) != ""
		parameters.Properties.VersionHeaderName = pointer.To(v.(string))
	}
	if v, ok := d.GetOk("version_query_name"); ok {
		querySet = v.(string) != ""
		parameters.Properties.VersionQueryName = pointer.To(v.(string))
	}

	switch schema := versioningScheme; schema {
	case apiversionset.VersioningSchemeHeader:
		if !headerSet {
			return fmt.Errorf("`version_header_name` must be set if `versioning_schema` is `Header`")
		}
		if querySet {
			return fmt.Errorf("`version_query_name` can not be set if `versioning_schema` is `Header`")
		}

	case apiversionset.VersioningSchemeQuery:
		if headerSet {
			return fmt.Errorf("`version_header_name` can not be set if `versioning_schema` is `Query`")
		}
		if !querySet {
			return fmt.Errorf("`version_query_name` must be set if `versioning_schema` is `Query`")
		}

	case apiversionset.VersioningSchemeSegment:
		if headerSet {
			return fmt.Errorf("`version_header_name` can not be set if `versioning_schema` is `Segment`")
		}
		if querySet {
			return fmt.Errorf("`version_query_name` can not be set if `versioning_schema` is `Segment`")
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apiversionset.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiVersionSetRead(d, meta)
}

func resourceApiManagementApiVersionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apiversionset.ParseApiVersionSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VersionSetId)
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

func resourceApiManagementApiVersionSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apiversionsets.ParseApiVersionSetID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.ApiVersionSetDelete(ctx, *id, apiversionsets.ApiVersionSetDeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
