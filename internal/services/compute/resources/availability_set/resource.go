// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package availability_set

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceAvailabilitySetCreateUpdate,
		Read:     resourceAvailabilitySetRead,
		Update:   resourceAvailabilitySetCreateUpdate,
		Delete:   resourceAvailabilitySetDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.AvailabilitySetId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.AvailabilitySetId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: availabilitySetSchema(),
	}
}
