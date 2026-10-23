// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.ServicesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting Healthcare Service %q (Resource Group %q): %+v", id.ServiceName, id.ResourceGroupName, err)
	}
	return nil
}
