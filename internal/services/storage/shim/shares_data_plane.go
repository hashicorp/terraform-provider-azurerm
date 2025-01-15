// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

type DataPlaneStorageShareWrapper struct {
	client *shares.Client
}

func NewDataPlaneStorageShareWrapper(client *shares.Client) StorageShareWrapper {
	return DataPlaneStorageShareWrapper{
		client: client,
	}
}

func (w DataPlaneStorageShareWrapper) Create(ctx context.Context, shareName string, input shares.CreateInput) error {
	if _, err := w.client.Create(ctx, shareName, input); err != nil {
		return fmt.Errorf("creating share: %+v", err)
	}
	return nil
}

func (w DataPlaneStorageShareWrapper) Delete(ctx context.Context, shareName string) error {
	input := shares.DeleteInput{
		DeleteSnapshots: true,
	}
	_, err := w.client.Delete(ctx, shareName, input)
	return err
}

func (w DataPlaneStorageShareWrapper) Exists(ctx context.Context, shareName string) (*bool, error) {
	existing, err := w.client.GetProperties(ctx, shareName)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, err
	}
	return pointer.To(true), nil
}

func (w DataPlaneStorageShareWrapper) Get(ctx context.Context, shareName string) (*StorageShareProperties, error) {
	props, err := w.client.GetProperties(ctx, shareName)
	if err != nil {
		if props.HttpResponse == nil {
			return nil, pollers.PollingDroppedConnectionError{
				Message: err.Error(),
			}
		}
		if response.WasNotFound(props.HttpResponse) {
			return nil, nil
		}

		return nil, err
	}

	acls, err := w.client.GetACL(ctx, shareName)
	if err != nil {
		return nil, err
	}

	return &StorageShareProperties{
		MetaData:        props.MetaData,
		QuotaGB:         props.QuotaInGB,
		ACLs:            acls.SignedIdentifiers,
		EnabledProtocol: props.EnabledProtocol,
		AccessTier:      props.AccessTier,
	}, nil
}

func (w DataPlaneStorageShareWrapper) UpdateACLs(ctx context.Context, shareName string, input shares.SetAclInput) error {
	_, err := w.client.SetACL(ctx, shareName, input)
	return err
}

func (w DataPlaneStorageShareWrapper) UpdateMetaData(ctx context.Context, shareName string, metaData map[string]string) error {
	input := shares.SetMetaDataInput{
		MetaData: metaData,
	}
	_, err := w.client.SetMetaData(ctx, shareName, input)
	return err
}

func (w DataPlaneStorageShareWrapper) UpdateQuota(ctx context.Context, shareName string, quotaGB int) error {
	_, err := w.client.SetProperties(ctx, shareName, shares.ShareProperties{
		QuotaInGb: &quotaGB,
	})
	return err
}

func (w DataPlaneStorageShareWrapper) UpdateTier(ctx context.Context, shareName string, tier shares.AccessTier) error {
	props := shares.ShareProperties{
		AccessTier: &tier,
	}
	_, err := w.client.SetProperties(ctx, shareName, props)
	return err
}
