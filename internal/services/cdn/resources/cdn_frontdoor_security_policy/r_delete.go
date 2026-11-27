// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_security_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontdoorSecurityPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := cdnClient.From(meta).FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, pointer.From(id)); err != nil {
		return fmt.Errorf("deleting %s: %+v", pointer.From(id), err)
	}

	return nil
}
