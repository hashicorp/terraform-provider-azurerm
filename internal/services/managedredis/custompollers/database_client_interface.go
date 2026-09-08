// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
)

type DatabasesClientInterface interface {
	Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error)
}
