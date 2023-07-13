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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverOutboundEndpointModel struct {
	Name                 string            `tfschema:"name"`
	PrivateDNSResolverId string            `tfschema:"private_dns_resolver_id"`
	Location             string            `tfschema:"location"`
	SubnetId             string            `tfschema:"subnet_id"`
	Tags                 map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverOutboundEndpointResource struct{}

var _ sdk.ResourceWithUpdate = PrivateDNSResolverOutboundEndpointResource{}

func (r PrivateDNSResolverOutboundEndpointResource) ResourceType() string {
	return "azurerm_private_dns_resolver_outbound_endpoint"
}

func (r PrivateDNSResolverOutboundEndpointResource) ModelObject() interface{} {
	return &PrivateDNSResolverOutboundEndpointModel{}
}

func (r PrivateDNSResolverOutboundEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outboundendpoints.ValidateOutboundEndpointID
}

func (r PrivateDNSResolverOutboundEndpointResource) Arguments() map[string]*pluginsdk.Schema {
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

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverOutboundEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverOutboundEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverOutboundEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.OutboundEndpointsClient
			dnsResolverId, err := dnsresolvers.ParseDnsResolverID(model.PrivateDNSResolverId)
			if err != nil {
				return err
			}

			id := outboundendpoints.NewOutboundEndpointID(dnsResolverId.SubscriptionId, dnsResolverId.ResourceGroupName, dnsResolverId.DnsResolverName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &outboundendpoints.OutboundEndpoint{
				Location: location.Normalize(model.Location),
				Properties: outboundendpoints.OutboundEndpointProperties{
					Subnet: outboundendpoints.SubResource{
						Id: model.SubnetId,
					},
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties, outboundendpoints.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PrivateDNSResolverOutboundEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.OutboundEndpointsClient

			id, err := outboundendpoints.ParseOutboundEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverOutboundEndpointModel
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

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, outboundendpoints.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverOutboundEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.OutboundEndpointsClient

			id, err := outboundendpoints.ParseOutboundEndpointID(metadata.ResourceData.Id())
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

			state := PrivateDNSResolverOutboundEndpointModel{
				Name:                 id.OutboundEndpointName,
				PrivateDNSResolverId: dnsresolvers.NewDnsResolverID(id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName).ID(),
				Location:             location.Normalize(model.Location),
			}

			properties := &model.Properties
			state.SubnetId = properties.Subnet.Id
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PrivateDNSResolverOutboundEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.OutboundEndpointsClient

			id, err := outboundendpoints.ParseOutboundEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, outboundendpoints.DeleteOperationOptions{}); err != nil {
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
				Refresh:                   dnsResolverOutboundEndpointDeleteRefreshFunc(ctx, client, id),
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

func dnsResolverOutboundEndpointDeleteRefreshFunc(ctx context.Context, client *outboundendpoints.OutboundEndpointsClient, id *outboundendpoints.OutboundEndpointId) pluginsdk.StateRefreshFunc {
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
