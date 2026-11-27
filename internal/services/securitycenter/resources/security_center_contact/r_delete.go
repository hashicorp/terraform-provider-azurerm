// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_contact

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterContactDelete(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).ContactsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.ContactID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.SecurityContactName); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
