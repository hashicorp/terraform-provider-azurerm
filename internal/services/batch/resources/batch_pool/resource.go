// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func Resource() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
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

		Schema: resourceBatch_poolSchema(),
	}

	resource.Identity = &schema.ResourceIdentity{
		SchemaFunc: pluginsdk.GenerateIdentitySchema(&pool.PoolId{}),
	}

	return resource
}
