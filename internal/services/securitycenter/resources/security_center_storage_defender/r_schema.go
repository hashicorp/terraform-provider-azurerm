// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_storage_defender

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (s Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"storage_account_id": {
			Description:  "The ID of the Storage Account the defender applied to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"override_subscription_settings_enabled": {
			Description: "Whether the settings defined for this Storage Account should override the settings defined for the subscription.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"malware_scanning_on_upload_enabled": {
			Description: "Whether On Upload malware scanning should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"malware_scanning_on_upload_cap_gb_per_month": {
			Description: "The max GB to be scanned per Month. Omit this property or set to `-1` if no capping is needed.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     -1,
			ValidateFunc: validation.Any(
				validation.IntAtLeast(1),
				validation.IntInSlice([]int{-1}),
			),
		},

		"sensitive_data_discovery_enabled": {
			Description: "Whether Sensitive Data Discovery should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"scan_results_event_grid_topic_id": commonschema.ResourceIDReferenceOptional(&topics.TopicId{}, "Event Grid Topic"),
	}
}

func (s Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
