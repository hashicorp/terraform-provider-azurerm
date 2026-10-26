// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_service_endpoint_storage_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/serviceendpointpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetServiceEndpointStoragePolicyCreate,
		Read:   resourceSubnetServiceEndpointStoragePolicyRead,
		Update: resourceSubnetServiceEndpointStoragePolicyUpdate,
		Delete: resourceSubnetServiceEndpointStoragePolicyDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&serviceendpointpolicies.ServiceEndpointPolicyId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&serviceendpointpolicies.ServiceEndpointPolicyId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: subnetServiceEndpointStoragePolicySchema(),
	}
}
