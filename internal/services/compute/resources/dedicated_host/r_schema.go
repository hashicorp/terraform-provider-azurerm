// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhosts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dedicatedHostSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DedicatedHostName(),
		},

		"dedicated_host_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateDedicatedHostGroupID,
		},

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"DADSv5-Type1",
				"DASv4-Type1",
				"DASv4-Type2",
				"DASv5-Type1",
				"DCSv2-Type1",
				"DDSv4-Type1",
				"DDSv4-Type2",
				"DDSv5-Type1",
				"DSv3-Type1",
				"DSv3-Type2",
				"DSv3-Type3",
				"DSv3-Type4",
				"DSv4-Type1",
				"DSv4-Type2",
				"DSv5-Type1",
				"EADSv5-Type1",
				"EASv4-Type1",
				"EASv4-Type2",
				"EASv5-Type1",
				"EDSv4-Type1",
				"EDSv4-Type2",
				"EDSv5-Type1",
				"ESv3-Type1",
				"ESv3-Type2",
				"ESv3-Type3",
				"ESv3-Type4",
				"ESv4-Type1",
				"ESv4-Type2",
				"ESv5-Type1",
				"FSv2-Type2",
				"FSv2-Type3",
				"FSv2-Type4",
				"FXmds-Type1",
				"LSv2-Type1",
				"LSv3-Type1",
				"MDMSv2MedMem-Type1",
				"MDSv2MedMem-Type1",
				"MMSv2MedMem-Type1",
				"MS-Type1",
				"MSm-Type1",
				"MSmv2-Type1",
				"MSv2-Type1",
				"MSv2MedMem-Type1",
				"NVASv4-Type1",
				"NVSv3-Type1",
			}, false),
		},

		"platform_fault_domain": {
			Type:     pluginsdk.TypeInt,
			ForceNew: true,
			Required: true,
		},

		"auto_replace_on_failure": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(dedicatedhosts.DedicatedHostLicenseTypesWindowsServerHybrid),
				string(dedicatedhosts.DedicatedHostLicenseTypesWindowsServerPerpetual),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
