// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint_custom_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmCdnEndpointCustomDomainDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customdomains.ParseEndpointCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}
