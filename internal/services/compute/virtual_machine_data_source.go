// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualMachineRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"private_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"public_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"power_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	networkInterfacesClient := meta.(*clients.Client).Network.NetworkInterfacesClient
	publicIPAddressesClient := meta.(*clients.Client).Network.PublicIPAddresses

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachines.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	options := virtualmachines.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachines.InstanceViewTypesInstanceView)
	resp, err := client.Get(ctx, id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if instance := props.InstanceView; instance != nil {
				if statues := instance.Statuses; statues != nil {
					for _, status := range *statues {
						if status.Code != nil && strings.HasPrefix(strings.ToLower(*status.Code), "powerstate/") {
							d.Set("power_state", strings.SplitN(*status.Code, "/", 2)[1])
						}
					}
				}
			}
			connectionInfo := retrieveConnectionInformation(ctx, networkInterfacesClient, publicIPAddressesClient, props)
			err = d.Set("private_ip_address", connectionInfo.primaryPrivateAddress)
			if err != nil {
				return err
			}
			err = d.Set("private_ip_addresses", connectionInfo.privateAddresses)
			if err != nil {
				return err
			}
			err = d.Set("public_ip_address", connectionInfo.primaryPublicAddress)
			if err != nil {
				return err
			}
			err = d.Set("public_ip_addresses", connectionInfo.publicAddresses)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
