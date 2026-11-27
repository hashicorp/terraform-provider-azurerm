// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_private_cloud

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	vmwareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwarePrivateCloudDelete(d *pluginsdk.ResourceData, meta any) error {
	client := vmwareClient.From(meta).PrivateCloudClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := privateclouds.ParsePrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
