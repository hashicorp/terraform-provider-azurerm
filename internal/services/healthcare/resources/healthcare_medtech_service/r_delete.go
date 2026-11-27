// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	healthcareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisMedTechServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := healthcareClient.From(meta).HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := iotconnectors.ParseIotConnectorID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
