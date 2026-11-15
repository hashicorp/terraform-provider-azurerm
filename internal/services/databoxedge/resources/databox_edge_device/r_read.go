// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databox_edge_device

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataboxEdge.DeviceClient

			id, err := devices.ParseDataBoxEdgeDeviceID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parse: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[INFO] %s was not found - removing from state", id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := EdgeDeviceModel{
				Name:              id.DataBoxEdgeDeviceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil {
					state.DeviceProperties = flattenDeviceProperties(props)
				}
				state.SkuName = flattenDeviceSku(model.Sku)
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDeviceProperties(input *devices.DataBoxEdgeDeviceProperties) []DevicePropertiesModel {
	output := make([]DevicePropertiesModel, 0)
	configuredRoleTypes := make([]string, 0)

	var status string
	var culture string
	var hcsVersion string
	var capacity int64
	var model string
	var softwareVersion string
	var deviceType string
	var serialNumber string
	var timeZone string

	if input != nil {
		o := DevicePropertiesModel{}
		if input.ConfiguredRoleTypes != nil {
			for _, item := range *input.ConfiguredRoleTypes {
				configuredRoleTypes = append(configuredRoleTypes, string(item))
			}
			o.ConfiguredRoleTypes = configuredRoleTypes
		}

		if v := input.DataBoxEdgeDeviceStatus; v != nil && *v != "" {
			status = string(*v)
			o.Status = status
		}

		if input.Culture != nil {
			culture = *input.Culture
			o.Culture = culture
		}

		if input.DeviceHcsVersion != nil {
			hcsVersion = *input.DeviceHcsVersion
			o.HcsVersion = hcsVersion
		}

		if input.DeviceLocalCapacity != nil {
			capacity = *input.DeviceLocalCapacity
			o.Capacity = capacity
		}

		if input.DeviceModel != nil {
			model = *input.DeviceModel
			o.Model = model
		}

		if input.DeviceSoftwareVersion != nil {
			softwareVersion = *input.DeviceSoftwareVersion
			o.SoftwareVersion = softwareVersion
		}

		if v := input.DeviceType; v != nil && *v != "" {
			deviceType = string(*v)
			o.Type = deviceType
		}

		if input.NodeCount != nil {
			o.NodeCount = *input.NodeCount
		}

		if input.SerialNumber != nil {
			serialNumber = *input.SerialNumber
			o.SerialNumber = serialNumber
		}

		if input.TimeZone != nil {
			timeZone = *input.TimeZone
			o.TimeZone = timeZone
		}

		output = append(output, o)
	}

	return output
}

func flattenDeviceSku(input *devices.Sku) string {
	if input == nil {
		return ""
	}

	var name devices.SkuName
	var tier devices.SkuTier

	if v := input.Name; v != nil && *v != "" {
		name = *v
	}

	if v := input.Tier; v != nil && *v != "" {
		tier = *v
	} else {
		tier = devices.SkuTierStandard
	}

	return fmt.Sprintf("%s-%s", name, tier)
}
