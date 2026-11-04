// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_customer_managed_key

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			log.Printf(`[INFO] Customer Managed Keys cannot be removed from Data Protection Backup Vaults once added. To remove the Customer Managed Key, delete and recreate the parent Data Protection Backup Vault`)
			return nil
		},
	}
}
