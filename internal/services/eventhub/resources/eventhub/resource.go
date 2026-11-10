// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name eventhub -properties "name" -compare-values "subscription_id:namespace_id,resource_group_name:namespace_id,namespace_name:namespace_id" -test-name standard

var eventHubResourceName = "azurerm_eventhub"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubCreate,
		Read:   resourceEventHubRead,
		Update: resourceEventHubUpdate,
		Delete: resourceEventHubDelete,

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&eventhubs.EventhubId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&eventhubs.EventhubId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: eventhubSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			if d.Id() == "" && d.Get("status").(string) == string(eventhubs.EntityStatusSendDisabled) {
				return fmt.Errorf("`status` cannot be set to `%s` when creating an Event Hub, it can only be set on an existing Event Hub", string(eventhubs.EntityStatusSendDisabled))
			}

			if capture := d.GetRawConfig().AsValueMap()["capture_description"]; capture.IsKnown() && !capture.IsNull() && capture.LengthInt() > 0 {
				if destination := capture.AsValueSlice()[0].AsValueMap()["destination"]; destination.IsKnown() && !destination.IsNull() && destination.LengthInt() > 0 {
					destinationMap := destination.AsValueSlice()[0].AsValueMap()
					authType := destinationMap["storage_authentication_type"]
					if authType.IsKnown() && !authType.IsNull() && authType.AsString() == string(eventhubs.CaptureIdentityTypeUserAssigned) && destinationMap["storage_authentication_id"].IsNull() {
						return fmt.Errorf("`storage_authentication_id` must be specified when `storage_authentication_type` is set to `UserAssigned`")
					}
				}
			}
			return nil
		}),
	}
}
