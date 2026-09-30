// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name dedicated_host -service-package-name compute -properties "name" -compare-values "subscription_id:dedicated_host_group_id,resource_group_name:dedicated_host_group_id,host_group_name:dedicated_host_group_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDedicatedHostCreate,
		Read:   resourceDedicatedHostRead,
		Update: resourceDedicatedHostUpdate,
		Delete: resourceDedicatedHostDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.DedicatedHostId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.DedicatedHostId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dedicatedHostSchema(),
	}
}
