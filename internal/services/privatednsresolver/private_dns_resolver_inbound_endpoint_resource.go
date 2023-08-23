// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverInboundEndpointModel struct {
	Name                 string                 `tfschema:"name"`
	PrivateDNSResolverId string                 `tfschema:"private_dns_resolver_id"`
	IPConfigurations     []IPConfigurationModel `tfschema:"ip_configurations"`
	Location             string                 `tfschema:"location"`
	Tags                 map[string]string      `tfschema:"tags"`
}

type IPConfigurationModel struct {
	PrivateIPAddress          string                              `tfschema:"private_ip_address"`
	PrivateIPAllocationMethod inboundendpoints.IPAllocationMethod `tfschema:"private_ip_allocation_method"`
	SubnetId                  string                              `tfschema:"subnet_id"`
}

type PrivateDNSResolverInboundEndpointResource struct{}

var _ sdk.ResourceWithUpdate = PrivateDNSResolverInboundEndpointResource{}

func (r PrivateDNSResolverInboundEndpointResource) ResourceType() string {
	return "azurerm_private_dns_resolver_inbound_endpoint"
}

func (r PrivateDNSResolverInboundEndpointResource) ModelObject() interface{} {
	return &PrivateDNSResolverInboundEndpointModel{}
}

func (r PrivateDNSResolverInboundEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return inboundendpoints.ValidateInboundEndpointID
}

func (r PrivateDNSResolverInboundEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"private_dns_resolver_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: dnsresolvers.ValidateDnsResolverID,
		},

		"ip_configurations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"private_ip_allocation_method": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(inboundendpoints.IPAllocationMethodDynamic),
						ValidateFunc: validation.StringInSlice([]string{
							string(inboundendpoints.IPAllocationMethodDynamic),
						}, false),
					},
				},
			},
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverInboundEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverInboundEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverInboundEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient
			dnsResolverId, err := dnsresolvers.ParseDnsResolverID(model.PrivateDNSResolverId)
			if err != nil {
				return err
			}

			id := inboundendpoints.NewInboundEndpointID(dnsResolverId.SubscriptionId, dnsResolverId.ResourceGroupName, dnsResolverId.DnsResolverName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &inboundendpoints.InboundEndpoint{
				Location:   location.Normalize(model.Location),
				Properties: inboundendpoints.InboundEndpointProperties{},
				Tags:       &model.Tags,
			}

			iPConfigurationsValue := expandIPConfigurationModel(model.IPConfigurations)

			if iPConfigurationsValue != nil {
				properties.Properties.IPConfigurations = *iPConfigurationsValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties, inboundendpoints.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PrivateDNSResolverInboundEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient

			id, err := inboundendpoints.ParseInboundEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverInboundEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("ip_configurations") {
				iPConfigurationsValue := expandIPConfigurationModel(model.IPConfigurations)

				if iPConfigurationsValue != nil {
					properties.Properties.IPConfigurations = *iPConfigurationsValue
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, inboundendpoints.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverInboundEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient

			id, err := inboundendpoints.ParseInboundEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PrivateDNSResolverInboundEndpointModel{
				Name:                 id.InboundEndpointName,
				PrivateDNSResolverId: dnsresolvers.NewDnsResolverID(id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName).ID(),
				Location:             location.Normalize(model.Location),
			}

			properties := &model.Properties

			state.IPConfigurations = flattenIPConfigurationModel(&properties.IPConfigurations)
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PrivateDNSResolverInboundEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.InboundEndpointsClient

			id, err := inboundendpoints.ParseInboundEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, inboundendpoints.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			log.Printf("[DEBUG] waiting for %s to be deleted", id)
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Pending"},
				Target:                    []string{"Succeeded"},
				Refresh:                   dnsResolverInboundEndpointDeleteRefreshFunc(ctx, client, id),
				MinTimeout:                1 * time.Minute,
				Timeout:                   time.Until(deadline),
				ContinuousTargetOccurence: 3,
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become deleted: %+v", id, err)
			}

			return nil
		},
	}
}

func dnsResolverInboundEndpointDeleteRefreshFunc(ctx context.Context, client *inboundendpoints.InboundEndpointsClient, id *inboundendpoints.InboundEndpointId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		existing, err := client.Get(ctx, *id)
		if err != nil {
			if response.WasNotFound(existing.HttpResponse) {
				return existing, "Succeeded", nil
			}
			return existing, "", err
		}
		return existing, "Pending", fmt.Errorf("checking for existing %s: %+v", id, err)
	}
}

func expandIPConfigurationModel(inputList []IPConfigurationModel) *[]inboundendpoints.IPConfiguration {
	var outputList []inboundendpoints.IPConfiguration
	for _, v := range inputList {
		input := v
		output := inboundendpoints.IPConfiguration{}

		if input.PrivateIPAllocationMethod != "" {
			output.PrivateIPAllocationMethod = &input.PrivateIPAllocationMethod
		}

		if input.PrivateIPAddress != "" {
			output.PrivateIPAddress = &input.PrivateIPAddress
		}

		output.Subnet = inboundendpoints.SubResource{
			Id: input.SubnetId,
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenIPConfigurationModel(inputList *[]inboundendpoints.IPConfiguration) []IPConfigurationModel {
	var outputList []IPConfigurationModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := IPConfigurationModel{}

		if input.PrivateIPAddress != nil {
			output.PrivateIPAddress = *input.PrivateIPAddress
		}

		if input.PrivateIPAllocationMethod != nil {
			output.PrivateIPAllocationMethod = *input.PrivateIPAllocationMethod
		}

		output.SubnetId = input.Subnet.Id

		outputList = append(outputList, output)
	}

	return outputList
}
