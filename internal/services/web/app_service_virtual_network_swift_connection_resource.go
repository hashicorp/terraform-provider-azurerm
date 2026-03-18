// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	networkpoller "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceVirtualNetworkSwiftConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceVirtualNetworkSwiftConnectionCreate,
		Read:   resourceAppServiceVirtualNetworkSwiftConnectionRead,
		Update: resourceAppServiceVirtualNetworkSwiftConnectionUpdate,
		Delete: resourceAppServiceVirtualNetworkSwiftConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkSwiftConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"app_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceID,
			},
			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},
		},
	}
}

func resourceAppServiceVirtualNetworkSwiftConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appID, err := commonids.ParseAppServiceID(d.Get("app_service_id").(string))
	if err != nil {
		return err
	}

	subnetID, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, *appID)
	if err != nil {
		return fmt.Errorf("checking for presence of Swift Network Connection for %s: %w", appID, err)
	}

	if existing.Model != nil && existing.Model.Properties != nil && pointer.From(existing.Model.Properties.SubnetResourceId) != "" {
		return tf.ImportAsExistsError("azurerm_app_service_virtual_network_swift_connection", pointer.From(existing.Model.Properties.SubnetResourceId))
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

	pollerType := networkpoller.NewVirtualNetworkAndSubnetProvisioningSucceededPoller(meta.(*clients.Client).Network, subnetID)
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

	d.SetId(swiftVirtualNetworkId.ID())

	return resourceAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}

func resourceAppServiceVirtualNetworkSwiftConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}
	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, appID)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appID, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil || existing.Model.Properties.SubnetResourceId == nil {
		d.SetId("")
		return nil
	}

	subnetID, err := commonids.ParseSubnetID(*existing.Model.Properties.SubnetResourceId)
	if err != nil {
		return err
	}

	d.Set("subnet_id", subnetID.ID())
	d.Set("app_service_id", appID.ID())

	return nil
}

func resourceAppServiceVirtualNetworkSwiftConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}
	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, appID)
	if err != nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appID, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: model was nil", appID)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: model was nil", appID)
	}

	subnetID, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetID.SubnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetID.SubnetName, network.SubnetResourceName)

	existing.Model.Properties.SubnetResourceId = pointer.To(d.Get("subnet_id").(string))

	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck(ctx, appID, *existing.Model); err != nil {
		return fmt.Errorf("creating association between %s and %s: %w", appID, subnetID, err)
	}

	pollerType := networkpoller.NewVirtualNetworkAndSubnetProvisioningSucceededPoller(meta.(*clients.Client).Network, subnetID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling for completion of association between %s and %s: %w", appID, subnetID, err)
	}

	return resourceAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}

func resourceAppServiceVirtualNetworkSwiftConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}
	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, appID)
	if err != nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil || existing.Model.Properties.SubnetResourceId == nil {
		// assume deleted
		d.SetId("")
		return nil
	}

	subnetID, err := commonids.ParseSubnetID(*existing.Model.Properties.SubnetResourceId)
	if err != nil {
		return err
	}

	locks.ByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetID.SubnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetID.SubnetName, network.SubnetResourceName)

	if _, err := client.DeleteSwiftVirtualNetwork(ctx, appID); err != nil {
		return fmt.Errorf("deleting Swift Network Connection for %s: %w", id, err)
	}

	return nil
}
