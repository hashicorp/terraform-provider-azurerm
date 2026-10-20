// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_kusto_database

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_share_dataset_kusto_database -service-package-name datashare -properties "name" -compare-values "subscription_id:share_id,resource_group_name:share_id,account_name:share_id,share_name:share_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataShareDataSetKustoDatabaseCreate,
		Read:   resourceDataShareDataSetKustoDatabaseRead,
		Delete: resourceDataShareDataSetKustoDatabaseDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&dataset.DataSetId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&dataset.DataSetId{}),
		},

		Schema: dataShareDatasetKustoDatabaseSchema(),
	}
}
