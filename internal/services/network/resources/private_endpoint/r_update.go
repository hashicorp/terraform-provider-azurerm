// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateEndpointUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateendpoints.ParsePrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	// Ensure we don't overwrite the existing ApplicationSecurityGroups
	existing, err := client.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model.Properties` was nil", *id)
	}

	applicationSecurityGroupAssociation := existing.Model.Properties.ApplicationSecurityGroups
	location := location.Normalize(d.Get("location").(string))
	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]any)
	privateServiceConnections := d.Get("private_service_connection").([]any)
	ipConfigurations := d.Get("ip_configuration").([]any)
	subnetId := d.Get("subnet_id").(string)
	customNicName := d.Get("custom_network_interface_name").(string)

	// Certain child resources like those for cognitive account lock the parent resource to make sure we don't try and update the parent when it's not ready.
	// Due to that, we'll lock on that resource id to try and prevent those type of errors.
	privateLinkServiceConnections, privateLinkServiceConnectionIds := expandPrivateLinkEndpointServiceConnection(privateServiceConnections, false)
	manualPrivateLinkServiceConnections, manualPrivateLinkServiceConnectionIds := expandPrivateLinkEndpointServiceConnection(privateServiceConnections, true)
	privateLinkServiceConnectionIds = append(privateLinkServiceConnectionIds, manualPrivateLinkServiceConnectionIds...)
	locks.MultipleByID(pointer.To(privateLinkServiceConnectionIds))
	defer locks.UnlockMultipleByID(pointer.To(privateLinkServiceConnectionIds))

	// TODO: in future it'd be nice to support conditional updates here, but one problem at a time
	parameters := privateendpoints.PrivateEndpoint{
		Location:         pointer.To(location),
		ExtendedLocation: helpers.ExpandEdgeZoneModel(d.Get("edge_zone").(string)),
		Properties: &privateendpoints.PrivateEndpointProperties{
			ApplicationSecurityGroups:           applicationSecurityGroupAssociation,
			PrivateLinkServiceConnections:       privateLinkServiceConnections,
			ManualPrivateLinkServiceConnections: manualPrivateLinkServiceConnections,
			Subnet: &privateendpoints.Subnet{
				Id: pointer.To(subnetId),
			},
			IPConfigurations:           expandPrivateEndpointIPConfigurations(ipConfigurations),
			CustomNetworkInterfaceName: pointer.To(customNicName),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err = validatePrivateLinkServiceId(*parameters.Properties.PrivateLinkServiceConnections); err != nil {
		return err
	}
	if err = validatePrivateLinkServiceId(*parameters.Properties.ManualPrivateLinkServiceConnections); err != nil {
		return err
	}

	if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if err = client.CreateOrUpdateThenPoll(ctx, *id, parameters); err != nil {
			switch {
			case strings.EqualFold(err.Error(), "is missing required parameter 'group Id'"):
				{
					return &pluginsdk.RetryError{
						Err: fmt.Errorf("updating %s due to missing 'group Id', ensure that the 'subresource_names' type is populated: %+v", id, err),
					}
				}
			case strings.Contains(err.Error(), "PrivateLinkServiceId Invalid private link service id"):
				{
					return &pluginsdk.RetryError{
						Err:       fmt.Errorf("creating Private Endpoint %s: %+v", id, err),
						Retryable: true,
					}
				}
			default:
				return &pluginsdk.RetryError{
					Err: fmt.Errorf("updating %s: %+v", id, err),
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// 1 Private Endpoint can have 1 Private DNS Zone Group - so to update we need to Delete & Recreate
	if d.HasChange("private_dns_zone_group") {
		existingDnsZoneGroups, err := retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx, dnsClient, *id)
		if err != nil {
			return err
		}

		newDnsZoneGroups := d.Get("private_dns_zone_group").([]any)
		newDnsZoneName := ""
		idHasBeenChanged := false
		if len(newDnsZoneGroups) > 0 {
			groupRaw := newDnsZoneGroups[0].(map[string]any)
			newDnsZoneName = groupRaw["name"].(string)

			// it is possible to add or remove a private_dns_zone_id, but if an id is added at the same time as one as been removed and the name has not been changed
			// an existing entry is updated, which is not allowed, so we need to delete the existing private dns zone groups
			if d.HasChange("private_dns_zone_group.0.private_dns_zone_ids") {
				o, n := d.GetChange("private_dns_zone_group.0.private_dns_zone_ids")
				if len(o.([]any)) == len(n.([]any)) {
					idHasBeenChanged = true
				}
			}
		}

		needToRemove := newDnsZoneName == ""
		nameHasChanged := false
		if existingDnsZoneGroups != nil && newDnsZoneName != "" {
			needToRemove = len(*existingDnsZoneGroups) > 0 && len(newDnsZoneGroups) == 0

			// there should only be a single one, but there's no harm checking all returned
			for _, existing := range *existingDnsZoneGroups {
				if existing.PrivateDnsZoneGroupName != newDnsZoneName {
					nameHasChanged = true
					break
				}
			}
		}

		if needToRemove || nameHasChanged || idHasBeenChanged {
			if err := deletePrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, *id); err != nil {
				return err
			}
		}

		if len(privateDnsZoneGroup) > 0 {
			if err := createPrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, *id, privateDnsZoneGroup); err != nil {
				return err
			}
		}
	}

	return resourcePrivateEndpointRead(d, meta)
}
