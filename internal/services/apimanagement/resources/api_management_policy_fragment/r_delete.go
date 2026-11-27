// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_policy_fragment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyFragmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).PolicyFragmentClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id, policyfragment.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %s", *id, err)
	}

	return nil
}
