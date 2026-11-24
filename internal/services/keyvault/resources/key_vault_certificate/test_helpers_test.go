// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyVaultResource struct{}

func (KeyVaultResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseKeyVaultID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err := client.KeyVault.VaultsClient.Delete(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}
