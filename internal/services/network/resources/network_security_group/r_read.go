// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_group

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkSecurityGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkSecurityGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, networksecuritygroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := resourceNetworkSecurityGroupFlatten(d, id, resp.Model); err != nil {
		return fmt.Errorf("flattening %s: %+v", id, err)
	}

	return nil
}

func validateSecurityRule(sgRule map[string]any) error {
	var err *multierror.Error

	sourcePortRange := sgRule["source_port_range"].(string)
	sourcePortRanges := sgRule["source_port_ranges"].(*pluginsdk.Set)
	destinationPortRange := sgRule["destination_port_range"].(string)
	destinationPortRanges := sgRule["destination_port_ranges"].(*pluginsdk.Set)
	sourceAddressPrefix := sgRule["source_address_prefix"].(string)
	sourceAddressPrefixes := sgRule["source_address_prefixes"].(*pluginsdk.Set)
	destinationAddressPrefix := sgRule["destination_address_prefix"].(string)
	destinationAddressPrefixes := sgRule["destination_address_prefixes"].(*pluginsdk.Set)

	if sourcePortRange != "" && sourcePortRanges.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"source_port_range\" and \"source_port_ranges\" can be used per security rule",
		))
	}
	if destinationPortRange != "" && destinationPortRanges.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"destination_port_range\" and \"destination_port_ranges\" can be used per security rule",
		))
	}
	if sourceAddressPrefix != "" && sourceAddressPrefixes.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"source_address_prefix\" and \"source_address_prefixes\" can be used per security rule",
		))
	}
	if destinationAddressPrefix != "" && destinationAddressPrefixes.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"destination_address_prefix\" and \"destination_address_prefixes\" can be used per security rule",
		))
	}

	return err.ErrorOrNil()
}

// hashNetworkSecurityRule implements a hash function for the `security_rule` property,
// mainly to normalize the casing for `source_application_security_group_ids` and `destination_application_security_group_ids`.
func hashNetworkSecurityRule(input any) int {
	var buf bytes.Buffer
	if m, ok := input.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["description"].(string))
		buf.WriteString(m["protocol"].(string))

		buf.WriteString(m["source_port_range"].(string))
		if set := m["source_port_ranges"].(*pluginsdk.Set); set != nil {
			buf.WriteString(set.GoString())
		}

		buf.WriteString(m["destination_port_range"].(string))
		if set := m["destination_port_ranges"].(*pluginsdk.Set); set != nil {
			buf.WriteString(set.GoString())
		}

		buf.WriteString(m["source_address_prefix"].(string))
		if set := m["source_address_prefixes"].(*pluginsdk.Set); set != nil {
			buf.WriteString(set.GoString())
		}

		buf.WriteString(m["destination_address_prefix"].(string))
		if set := m["destination_address_prefixes"].(*pluginsdk.Set); set != nil {
			buf.WriteString(set.GoString())
		}

		if set := m["source_application_security_group_ids"].(*pluginsdk.Set); set != nil {
			for _, elem := range set.List() {
				buf.WriteString(strings.ToLower(elem.(string)))
			}
		}

		if set := m["destination_application_security_group_ids"].(*pluginsdk.Set); set != nil {
			for _, elem := range set.List() {
				buf.WriteString(strings.ToLower(elem.(string)))
			}
		}

		buf.WriteString(m["access"].(string))
		buf.WriteString(m["direction"].(string))
		buf.WriteString(strconv.Itoa(m["priority"].(int)))
	}
	return pluginsdk.HashString(buf.String())
}
