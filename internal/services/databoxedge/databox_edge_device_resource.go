package databoxedge

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2020-12-01/databoxedge" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevicePropertiesModel struct {
	Capacity            int64    `tfschema:"capacity"`
	ConfiguredRoleTypes []string `tfschema:"configured_role_types"`
	Culture             string   `tfschema:"culture"`
	HcsVersion          string   `tfschema:"hcs_version"`
	Model               string   `tfschema:"model"`
	NodeCount           int32    `tfschema:"node_count"`
	SerialNumber        string   `tfschema:"serial_number"`
	SoftwareVersion     string   `tfschema:"software_version"`
	Status              string   `tfschema:"status"`
	TimeZone            string   `tfschema:"time_zone"`
	Type                string   `tfschema:"type"`
}
type EdgeDeviceModel struct {
	DeviceProperties  []DevicePropertiesModel `tfschema:"device_properties"`
	Location          string                  `tfschema:"location"`
	Name              string                  `tfschema:"name"`
	ResourceGroupName string                  `tfschema:"resource_group_name"`
	SkuName           string                  `tfschema:"sku_name"`
	Tags              map[string]string       `tfschema:"tags"`
}

type EdgeDeviceResource struct{}

var _ sdk.ResourceWithUpdate = EdgeDeviceResource{}

func (r EdgeDeviceResource) ModelObject() interface{} {
	return &EdgeDeviceModel{}
}

func (r EdgeDeviceResource) ResourceType() string {
	return "azurerm_databox_edge_device"
}

func (r EdgeDeviceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DeviceID
}

func (r EdgeDeviceResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataboxEdgeName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataboxEdgeDeviceSkuName,
		},

		"tags": tags.Schema(),
	}
}

func (r EdgeDeviceResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"device_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"configured_role_types": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"culture": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hcs_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"model": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"status": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"software_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"node_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"serial_number": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_zone": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r EdgeDeviceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.DataboxEdge.DeviceClient

			var metaModel EdgeDeviceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := parse.NewDeviceID(subscriptionId, metaModel.ResourceGroupName, metaModel.Name)
			// sdk method is Get(ctx context.Context, deviceName string, resourceGroupName string)
			existing, err := client.Get(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError("azurerm_databox_edge_device", id.ID())
			}

			dataBoxEdgeDevice := databoxedge.Device{
				Location: utils.String(location.Normalize(metaModel.Location)),
				Sku:      expandDeviceSku(metaModel.SkuName),
				Tags:     tags.FromTypedObject(metaModel.Tags),
			}

			future, err := client.CreateOrUpdate(ctx, id.DataBoxEdgeDeviceName, dataBoxEdgeDevice, id.ResourceGroup)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r EdgeDeviceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataboxEdge.DeviceClient

			id, err := parse.DeviceID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parse: %+v", err)
			}

			resp, err := client.Get(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					log.Printf("[INFO] %s was not found - removing from state", id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := EdgeDeviceModel{
				Name:              id.DataBoxEdgeDeviceName,
				ResourceGroupName: id.ResourceGroup,
				Location:          location.NormalizeNilable(resp.Location),
			}

			if props := resp.DeviceProperties; props != nil {
				state.DeviceProperties = flattenDeviceProperties(props)
				state.SkuName = flattenDeviceSku(resp.Sku)
				state.Tags = tags.ToTypedObject(resp.Tags)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func (r EdgeDeviceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataboxEdge.DeviceClient

			id, err := parse.DeviceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var metaModel EdgeDeviceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := databoxedge.DevicePatch{}
			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.FromTypedObject(metaModel.Tags)
			}

			if _, err := client.Update(ctx, id.DataBoxEdgeDeviceName, parameters, id.ResourceGroup); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r EdgeDeviceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataboxEdge.DeviceClient

			id, err := parse.DeviceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var metaModel EdgeDeviceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			future, err := client.Delete(ctx, id.DataBoxEdgeDeviceName, id.ResourceGroup)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func expandDeviceSku(input string) *databoxedge.Sku {
	if len(input) == 0 {
		return nil
	}

	v, err := parse.DataboxEdgeDeviceSkuName(input)
	if err != nil {
		return nil
	}

	return &databoxedge.Sku{
		Name: databoxedge.SkuName(v.Name),
		Tier: databoxedge.SkuTier(v.Tier),
	}
}

func flattenDeviceProperties(input *databoxedge.DeviceProperties) []DevicePropertiesModel {
	output := make([]DevicePropertiesModel, 0)
	configuredRoleTypes := make([]string, 0)

	var status string
	var culture string
	var hcsVersion string
	var capacity int64
	var model string
	var softwareVersion string
	var deviceType string
	var nodeCount int32
	var serialNumber string
	var timeZone string

	if input != nil {
		o := DevicePropertiesModel{}
		if input.ConfiguredRoleTypes != nil {
			for _, item := range *input.ConfiguredRoleTypes {
				configuredRoleTypes = append(configuredRoleTypes, (string)(item))
			}
			o.ConfiguredRoleTypes = configuredRoleTypes
		}

		if input.DataBoxEdgeDeviceStatus != "" {
			status = string(input.DataBoxEdgeDeviceStatus)
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

		if input.DeviceType != "" {
			deviceType = string(input.DeviceType)
			o.Type = deviceType
		}

		if input.NodeCount != nil {
			nodeCount = *input.NodeCount
			o.NodeCount = nodeCount
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

func flattenDeviceSku(input *databoxedge.Sku) string {
	if input == nil {
		return ""
	}

	var name databoxedge.SkuName
	var tier databoxedge.SkuTier

	if input.Name != "" {
		name = input.Name
	}

	if input.Tier != "" {
		tier = input.Tier
	} else {
		tier = databoxedge.Standard
	}

	skuName := fmt.Sprintf("%s-%s", name, tier)

	return skuName
}
