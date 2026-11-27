// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_virtual_network_swift_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/custompollers"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceVirtualNetworkSwiftConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	appID, err := commonids.ParseAppServiceID(d.Get("app_service_id").(string))
	if err != nil {
		return err
	}

	subnetID, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetSwiftVirtualNetworkConnection(ctx, *appID)
		if err != nil {
			return fmt.Errorf("checking for presence of Swift Network Connection for %s: %w", appID, err)
		}

		if existing.Model != nil && existing.Model.Properties != nil && pointer.From(existing.Model.Properties.SubnetResourceId) != "" {
			return tf.ImportAsExistsError("azurerm_app_service_virtual_network_swift_connection", pointer.From(existing.Model.Properties.SubnetResourceId))
		}
	}

	locks.ByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetID.SubnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetID.SubnetName, network.SubnetResourceName)

	connectionEnvelope := webapps.SwiftVirtualNetwork{
		Properties: &webapps.SwiftVirtualNetworkProperties{
			SubnetResourceId: pointer.To(d.Get("subnet_id").(string)),
		},
	}

	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck(ctx, *appID, connectionEnvelope); err != nil {
		return fmt.Errorf("creating association between %s and %s: %w", appID, subnetID, err)
	}

	pollerType := custompollers.NewVirtualNetworkAndSubnetProvisioningSucceededPoller(networkClient.From(meta), subnetID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling for completion of association between %s and %s: %w", appID, subnetID, err)
	}

	read, err := client.GetSwiftVirtualNetworkConnection(ctx, *appID)
	if err != nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appID, err)
	}

	if read.Model == nil || read.Model.Id == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: missing ID", appID)
	}

	swiftVirtualNetworkId, err := parse.VirtualNetworkSwiftConnectionID(*read.Model.Id)
	if err != nil {
		return err
	}

	// TODO: migrate to a typed resource ID
	d.SetId(swiftVirtualNetworkId.ID())

	return resourceAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}
