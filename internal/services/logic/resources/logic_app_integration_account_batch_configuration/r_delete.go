// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_batch_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountBatchConfigurationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountBatchConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountbatchconfigurations.ParseBatchConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
