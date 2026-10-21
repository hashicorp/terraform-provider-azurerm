// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var networkInterfaceResourceName = "azurerm_network_interface"

func RegisterResource() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceNetworkInterfaceCreate,
		Read:   resourceNetworkInterfaceRead,
		Update: resourceNetworkInterfaceUpdate,
		Delete: resourceNetworkInterfaceDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.NetworkInterfaceId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.NetworkInterfaceId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: networkInterfaceSchema(),
	}

	if !features.SixPointOh() {
		r.Schema["auxiliary_mode"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(networkinterfaces.PossibleValuesForNetworkInterfaceAuxiliaryMode(), false),
			RequiredWith: []string{"auxiliary_sku"},
		}

		r.Schema["auxiliary_sku"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(networkinterfaces.PossibleValuesForNetworkInterfaceAuxiliarySku(), false),
			RequiredWith: []string{"auxiliary_mode"},
		}
	}

	return r
}
