// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

type StorageContainerWrapper interface {
	Create(ctx context.Context, containerName string, input containers.CreateInput) error
	Delete(ctx context.Context, containerName string) error
	Exists(ctx context.Context, containerName string) (*bool, error)
	Get(ctx context.Context, containerName string) (*StorageContainerProperties, error)
	UpdateAccessLevel(ctx context.Context, containerName string, level containers.AccessLevel) error
	UpdateMetaData(ctx context.Context, containerName string, metaData map[string]string) error
}

type StorageContainerProperties struct {
	AccessLevel                     containers.AccessLevel
	DefaultEncryptionScope          string
	EncryptionScopeOverrideDisabled bool
	MetaData                        map[string]string
	HasImmutabilityPolicy           bool
	HasLegalHold                    bool
}
