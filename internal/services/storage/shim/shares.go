// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
)

type StorageShareWrapper interface {
	Create(ctx context.Context, resourceGroup, accountName, shareName string, input shares.CreateInput) error
	Delete(ctx context.Context, resourceGroup, accountName, shareName string) error
	Exists(ctx context.Context, resourceGroup, accountName, shareName string) (*bool, error)
	Get(ctx context.Context, resourceGroup, accountName, shareName string) (*StorageShareProperties, error)
	UpdateACLs(ctx context.Context, resourceGroup, accountName, shareName string, acls []shares.SignedIdentifier) error
	UpdateMetaData(ctx context.Context, resourceGroup, accountName, shareName string, metaData map[string]string) error
	UpdateQuota(ctx context.Context, resourceGroup, accountName, shareName string, quotaGB int) error
	UpdateTier(ctx context.Context, resourceGroup, accountName, shareName string, tier shares.AccessTier) error
}

type StorageShareProperties struct {
	ACLs            []shares.SignedIdentifier
	MetaData        map[string]string
	QuotaGB         int
	EnabledProtocol shares.ShareProtocol
	AccessTier      *shares.AccessTier
}
