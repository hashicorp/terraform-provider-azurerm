// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Custom IP Prefix.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CustomIpPrefixName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"cidr": {
			Description:  "The `cidr` of the Custom IP Prefix, either IPv4 or IPv6.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"parent_custom_ip_prefix_id": {
			Description:  "The ID of the parent prefix. Only needed when creating a regional/child IPv6 prefix.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: customipprefixes.ValidateCustomIPPrefixID,
		},

		"roa_validity_end_date": {
			Description: "The expiration date of the Route Origin Authorization (ROA) document which has been filed with the Routing Internet Registry (RIR) for this prefix. The expected format is `YYYY-MM-DD`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: func(i any, k string) (warnings []string, errors []error) {
				v, ok := i.(string)
				if !ok {
					errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
					return warnings, errors
				}

				if _, err := time.Parse("2006-01-02", v); err != nil {
					errors = append(errors, fmt.Errorf("expected %q to be a valid date in the format YYYY-MM-DD, got %q: %+v", k, i, err))
				}

				return warnings, errors
			},
		},

		"wan_validation_signed_message": {
			Description:  "The signed base64-encoded authorization message, which will be sent to Microsoft for WAN verification. Refer to [Azure documentation](https://learn.microsoft.com/azure/virtual-network/ip-services/create-custom-ip-address-prefix-cli#certificate-readiness) for more details about the process for your RIR.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"commissioning_enabled": {
			Description: "The custom IP prefix should be commissioned after provisioning in Azure.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"internet_advertising_disabled": {
			Description: "The custom IP prefix should not be publicly advertised on the Internet when commissioned (regional commissioning feature).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
