// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_version_set

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiVersionSetCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apiversionset.NewApiVersionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			return errors.New("`version_header_name` must be set if `versioning_schema` is `Header`")
		}
		if querySet {
			return errors.New("`version_query_name` can not be set if `versioning_schema` is `Header`")
		}

	case apiversionset.VersioningSchemeQuery:
		if headerSet {
			return errors.New("`version_header_name` can not be set if `versioning_schema` is `Query`")
		}
		if !querySet {
			return errors.New("`version_query_name` must be set if `versioning_schema` is `Query`")
		}

	case apiversionset.VersioningSchemeSegment:
		if headerSet {
			return errors.New("`version_header_name` can not be set if `versioning_schema` is `Segment`")
		}
		if querySet {
			return errors.New("`version_query_name` can not be set if `versioning_schema` is `Segment`")
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apiversionset.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiVersionSetRead(d, meta)
}
