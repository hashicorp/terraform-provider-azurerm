// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

type StorageShareWrapper interface {
	Create(ctx context.Context, shareName string, input shares.CreateInput) error
	Delete(ctx context.Context, shareName string) error
	Exists(ctx context.Context, shareName string) (*bool, error)
	Get(ctx context.Context, shareName string) (*StorageShareProperties, error)
	UpdateACLs(ctx context.Context, shareName string, input shares.SetAclInput) error
	UpdateMetaData(ctx context.Context, shareName string, metaData map[string]string) error
	UpdateQuota(ctx context.Context, shareName string, quotaGB int) error
	UpdateTier(ctx context.Context, shareName string, tier shares.AccessTier) error
}

type StorageShareProperties struct {
	ACLs            []shares.SignedIdentifier
	MetaData        map[string]string
	QuotaGB         int
	EnabledProtocol shares.ShareProtocol
	AccessTier      *shares.AccessTier
}
