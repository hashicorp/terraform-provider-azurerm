// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_security_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceApplicationSecurityGroupFlatten(d *pluginsdk.ResourceData, id *applicationsecuritygroups.ApplicationSecurityGroupId, model *applicationsecuritygroups.ApplicationSecurityGroup) error {
	d.Set("name", id.ApplicationSecurityGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
