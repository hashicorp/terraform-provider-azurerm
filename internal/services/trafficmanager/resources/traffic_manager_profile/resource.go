// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_profile

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-params "Geographic"

const azureTrafficManagerProfileResourceName = "azurerm_traffic_manager_profile"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmTrafficManagerProfileCreate,
		Read:   resourceArmTrafficManagerProfileRead,
		Update: resourceArmTrafficManagerProfileUpdate,
		Delete: resourceArmTrafficManagerProfileDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&profiles.TrafficManagerProfileId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&profiles.TrafficManagerProfileId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: trafficManagerProfileSchema(),
	}
}
