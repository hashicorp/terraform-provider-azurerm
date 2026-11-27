// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_workspace

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterWorkspaceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.WorkspaceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Delete(ctx, securityCenterWorkspaceName)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription Workspace was not found: %v", err)
			return nil
		}

		return fmt.Errorf("deleting Security Center Workspace: %+v", err)
	}

	return nil
}
