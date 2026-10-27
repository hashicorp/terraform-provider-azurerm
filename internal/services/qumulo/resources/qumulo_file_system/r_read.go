// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package qumulo_file_system

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2026-04-16/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Qumulo.FileSystemsClient

			id, err := filesystems.ParseFileSystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := FileSystemResourceSchema{
				Name:              id.FileSystemName,
				ResourceGroupName: id.ResourceGroupName,

				AdminPassword: config.AdminPassword,
				Email:         config.Email,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props := model.Properties
				state.OfferId = props.MarketplaceDetails.OfferId
				state.PlanId = props.MarketplaceDetails.PlanId
				state.PublisherId = pointer.From(props.MarketplaceDetails.PublisherId)
				state.StorageSku = props.StorageSku
				state.Zone = pointer.From(props.AvailabilityZone)

				subnetId, err := commonids.ParseSubnetIDInsensitively(props.DelegatedSubnetId)
				if err != nil {
					return err
				}
				state.SubnetId = subnetId.ID()
			}

			return metadata.Encode(&state)
		},
	}
}

type FileSystemResourceSchema struct {
	AdminPassword     string            `tfschema:"admin_password"`
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	OfferId           string            `tfschema:"offer_id"`
	PlanId            string            `tfschema:"plan_id"`
	PublisherId       string            `tfschema:"publisher_id"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	StorageSku        string            `tfschema:"storage_sku"`
	SubnetId          string            `tfschema:"subnet_id"`
	Tags              map[string]string `tfschema:"tags"`
	Email             string            `tfschema:"email"`
	Zone              string            `tfschema:"zone"`
}

func checkSubnet(ctx context.Context, rawSubnetId string, metadata sdk.ResourceMetaData) error {
	delegationAction := "Microsoft.Network/virtualNetworks/subnets/join/action"
	delegationName := "Qumulo.Storage/fileSystems"

	subnetId, err := commonids.ParseSubnetID(rawSubnetId)
	if err != nil {
		return err
	}

	subnet, err := metadata.Client.Network.Subnets.Get(ctx, *subnetId, subnets.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", subnetId, err)
	}

	if subnet.Model != nil && subnet.Model.Properties != nil && subnet.Model.Properties.Delegations != nil {
		for _, delegation := range *subnet.Model.Properties.Delegations {
			if props := delegation.Properties; props != nil && props.Actions != nil &&
				props.ServiceName != nil && strings.EqualFold(*props.ServiceName, delegationName) {
				for _, action := range *props.Actions {
					if strings.EqualFold(action, delegationAction) {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("subnet %q is missing action %q on delegation %q", rawSubnetId, delegationAction, delegationName)
}
