// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Virtual Machine.",
		Read:        dataSourceVirtualMachineRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Virtual Machine.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"private_ip_address": {
				Description: "The Primary Private IP Address assigned to this Virtual Machine.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"private_ip_addresses": {
				Description: "A list of Private IP Addresses assigned to this Virtual Machine.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"public_ip_address": {
				Description: "The Primary Public IP Address assigned to this Virtual Machine.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"public_ip_addresses": {
				Description: "A list of the Public IP Addresses assigned to this Virtual Machine.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"power_state": {
				Description: "The power state of the Virtual Machine.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).VirtualMachinesClient
	subscriptionId := account.From(meta).SubscriptionId
	networkInterfacesClient := networkClient.From(meta).NetworkInterfacesClient
	publicIPAddressesClient := networkClient.From(meta).PublicIPAddresses

	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
			connectionInfo := helpers.RetrieveConnectionInformation(ctx, networkInterfacesClient, publicIPAddressesClient, props)
			if err = d.Set("private_ip_address", connectionInfo.PrimaryPrivateAddress); err != nil {
				return err
			}
			if err = d.Set("private_ip_addresses", connectionInfo.PrivateAddresses); err != nil {
				return err
			}
			if err = d.Set("public_ip_address", connectionInfo.PrimaryPublicAddress); err != nil {
				return err
			}
			if err = d.Set("public_ip_addresses", connectionInfo.PublicAddresses); err != nil {
				return err
			}
		}
	}
	return nil
}
