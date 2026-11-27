// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_tenant_configuration

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	portalClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePortalTenantConfigurationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := portalClient.From(meta).TenantConfigurationsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.PortalTenantConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
