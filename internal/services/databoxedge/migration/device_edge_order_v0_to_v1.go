// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/orders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DataBoxEdgeOrderV0ToV1{}

type DataBoxEdgeOrderV0ToV1 struct{}

func (DataBoxEdgeOrderV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"device_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"contact": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"company_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"emails": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"phone_number": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"status": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"info": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"comments": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"additional_details": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"last_update": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"shipment_address": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 3,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"city": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"country": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"state": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"postal_code": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"shipment_tracking": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"carrier_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"serial_number": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"tracking_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"tracking_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"shipment_history": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_details": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"comments": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_update": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"return_tracking": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"carrier_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"serial_number": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"tracking_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"tracking_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"serial_number": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (DataBoxEdgeOrderV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := parseLegacyResourceID(oldIdRaw)
		if err != nil {
			return rawState, fmt.Errorf("parsing the existing resource id %q: %+v", oldId, err)
		}
		newId := orders.NewDataBoxEdgeDeviceID(oldId.SubscriptionId, oldId.ResourceGroupName, oldId.DataBoxEdgeDeviceName).ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}

type legacyResourceId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DataBoxEdgeDeviceName string
}

func parseLegacyResourceID(input string) (*legacyResourceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := legacyResourceId{
		SubscriptionId:    id.SubscriptionID,
		ResourceGroupName: id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroupName == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DataBoxEdgeDeviceName, err = id.PopSegment("dataBoxEdgeDevices"); err != nil {
		return nil, err
	}

	// whilst there will be one set of segments left (`/orders/default`), we're intentionally ignoring them since we're unconcerned with them for now

	return &resourceId, nil
}
