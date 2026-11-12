// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridDomainDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.Domains
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domains.ParseDomainID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
