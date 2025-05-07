package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/staticcidrs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagerIpamPoolStaticCidrResource{}

type ManagerIpamPoolStaticCidrResource struct{}

func (ManagerIpamPoolStaticCidrResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticcidrs.ValidateStaticCidrID
}

func (ManagerIpamPoolStaticCidrResource) ResourceType() string {
	return "azurerm_network_manager_ipam_pool_static_cidr"
}

func (ManagerIpamPoolStaticCidrResource) ModelObject() interface{} {
	return &ManagerIpamPoolStaticCidrResourceModel{}
}

type ManagerIpamPoolStaticCidrResourceModel struct {
	AddressPrefixes []string `tfschema:"address_prefixes"`
	Description     string   `tfschema:"description"`
	IpAddressNumber string   `tfschema:"ip_address_number"`
	Name            string   `tfschema:"name"`
	IpamPoolId      string   `tfschema:"ipam_pool_id"`
}

func (ManagerIpamPoolStaticCidrResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"ipam_pool_id": commonschema.ResourceIDReferenceRequiredForceNew(&staticcidrs.IPamPoolId{}),

		"address_prefixes": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"address_prefixes", "ip_address_number"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"ip_address_number": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"address_prefixes", "ip_address_number"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[1-9]\d*$`),
				"`ip_address_number` must be a string that represents a positive number"),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (ManagerIpamPoolStaticCidrResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerIpamPoolStaticCidrResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticCidrs
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerIpamPoolStaticCidrResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			ipamPoolId, err := staticcidrs.ParseIPamPoolID(config.IpamPoolId)
			if err != nil {
				return err
			}

			id := staticcidrs.NewStaticCidrID(subscriptionId, ipamPoolId.ResourceGroupName, ipamPoolId.NetworkManagerName, ipamPoolId.IpamPoolName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := staticcidrs.StaticCidr{
				Name: pointer.To(config.Name),
				Properties: &staticcidrs.StaticCidrProperties{
					Description: pointer.To(config.Description),
				},
			}

			if len(config.AddressPrefixes) > 0 {
				payload.Properties.AddressPrefixes = pointer.To(config.AddressPrefixes)
			}

			if config.IpAddressNumber != "" {
				payload.Properties.NumberOfIPAddressesToAllocate = pointer.To(config.IpAddressNumber)
			}

			if _, err := client.Create(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagerIpamPoolStaticCidrResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticCidrs

			id, err := staticcidrs.ParseStaticCidrID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			ipamPoolId := staticcidrs.NewIPamPoolID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.IpamPoolName).ID()
			schema := ManagerIpamPoolStaticCidrResourceModel{
				Name:       id.StaticCidrName,
				IpamPoolId: ipamPoolId,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.AddressPrefixes = pointer.From(props.AddressPrefixes)
					schema.Description = pointer.From(props.Description)
					schema.IpAddressNumber = pointer.From(props.NumberOfIPAddressesToAllocate)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerIpamPoolStaticCidrResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticCidrs

			id, err := staticcidrs.ParseStaticCidrID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerIpamPoolStaticCidrResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` is nil", id)
			}

			parameters := resp.Model

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("address_prefixes") {
				if len(model.AddressPrefixes) > 0 {
					parameters.Properties.AddressPrefixes = pointer.To(model.AddressPrefixes)
					parameters.Properties.NumberOfIPAddressesToAllocate = pointer.To("")
				} else {
					parameters.Properties.AddressPrefixes = pointer.To([]string{})
				}
			}

			if metadata.ResourceData.HasChange("ip_address_number") {
				if model.IpAddressNumber != "" {
					parameters.Properties.NumberOfIPAddressesToAllocate = pointer.To(model.IpAddressNumber)
					parameters.Properties.AddressPrefixes = pointer.To([]string{})
				} else {
					parameters.Properties.NumberOfIPAddressesToAllocate = pointer.To("")
				}
			}

			if _, err := client.Create(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r ManagerIpamPoolStaticCidrResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticCidrs

			id, err := staticcidrs.ParseStaticCidrID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
