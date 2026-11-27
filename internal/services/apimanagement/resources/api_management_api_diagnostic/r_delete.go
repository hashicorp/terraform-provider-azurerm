// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_diagnostic

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiDiagnosticDelete(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).ApiDiagnosticClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	diagnosticId, err := apidiagnostic.ParseApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	name := helpers.GetApiName(diagnosticId.ApiId)

	newId := apidiagnostic.NewApiDiagnosticID(diagnosticId.SubscriptionId, diagnosticId.ResourceGroupName, diagnosticId.ServiceName, name, diagnosticId.DiagnosticId)
	if resp, err := client.Delete(ctx, newId, apidiagnostic.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", newId, err)
		}
	}

	return nil
}
