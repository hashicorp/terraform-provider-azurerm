// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pool

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../tools/generator-tests resourceidentity -resource-name batch_pool -service-package-name batch -properties "name,resource_group_name,batch_account_name:account_name" -known-values "subscription_id:data.Subscriptions.Primary"

const ResourceName = "azurerm_batch_pool"

func Resource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBatchCreate,
		Read:   resourceBatchPoolRead,
		Update: resourceBatchUpdate,
		Delete: resourceBatchPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&pool.PoolId{}),

		Schema: resourceBatchPoolSchema(),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&pool.PoolId{}),
		},
	}
}
