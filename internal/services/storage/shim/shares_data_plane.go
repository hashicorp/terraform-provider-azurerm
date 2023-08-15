// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
)

type DataPlaneStorageShareWrapper struct {
	client *shares.Client
}

func NewDataPlaneStorageShareWrapper(client *shares.Client) StorageShareWrapper {
	return DataPlaneStorageShareWrapper{
		client: client,
	}
}

func (w DataPlaneStorageShareWrapper) Create(ctx context.Context, _, accountName, shareName string, input shares.CreateInput) error {
	timeout, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context is missing a timeout")
	}

	resp, err := w.client.Create(ctx, accountName, shareName, input)
	if err == nil {
		return nil
	}

	// If we fail due to previous delete still in progress, then we can retry
	if utils.ResponseWasConflict(resp) && strings.Contains(err.Error(), "ShareBeingDeleted") {
		stateConf := &pluginsdk.StateChangeConf{
			Pending:        []string{"waitingOnDelete"},
			Target:         []string{"succeeded"},
			Refresh:        w.createRefreshFunc(ctx, accountName, shareName, input),
			PollInterval:   10 * time.Second,
			NotFoundChecks: 180,
			Timeout:        time.Until(timeout),
		}

		_, err := stateConf.WaitForStateContext(ctx)
		return err
	}

	// otherwise it's a legit error, so raise it
	return err
}

func (w DataPlaneStorageShareWrapper) Delete(ctx context.Context, _, accountName, shareName string) error {
	deleteSnapshots := true
	_, err := w.client.Delete(ctx, accountName, shareName, deleteSnapshots)
	return err
}

func (w DataPlaneStorageShareWrapper) Exists(ctx context.Context, _, accountName, shareName string) (*bool, error) {
	existing, err := w.client.GetProperties(ctx, accountName, shareName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil, nil
		}

		return nil, err
	}

	return utils.Bool(true), nil
}

func (w DataPlaneStorageShareWrapper) Get(ctx context.Context, _, accountName, shareName string) (*StorageShareProperties, error) {
	props, err := w.client.GetProperties(ctx, accountName, shareName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			return nil, nil
		}

		return nil, err
	}

	acls, err := w.client.GetACL(ctx, accountName, shareName)
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

func (w DataPlaneStorageShareWrapper) UpdateACLs(ctx context.Context, _, accountName, shareName string, acls []shares.SignedIdentifier) error {
	_, err := w.client.SetACL(ctx, accountName, shareName, acls)
	return err
}

func (w DataPlaneStorageShareWrapper) UpdateMetaData(ctx context.Context, _, accountName, shareName string, metaData map[string]string) error {
	_, err := w.client.SetMetaData(ctx, accountName, shareName, metaData)
	return err
}

func (w DataPlaneStorageShareWrapper) UpdateQuota(ctx context.Context, _, accountName, shareName string, quotaGB int) error {
	_, err := w.client.SetProperties(ctx, accountName, shareName, shares.ShareProperties{
		QuotaInGb: &quotaGB,
	})
	return err
}

func (w DataPlaneStorageShareWrapper) UpdateTier(ctx context.Context, _, accountname, shareName string, tier shares.AccessTier) error {
	props := shares.ShareProperties{
		AccessTier: &tier,
	}
	_, err := w.client.SetProperties(ctx, accountname, shareName, props)
	return err
}

func (w DataPlaneStorageShareWrapper) createRefreshFunc(ctx context.Context, accountName string, shareName string, input shares.CreateInput) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := w.client.Create(ctx, accountName, shareName, input)
		if err != nil {
			if !utils.ResponseWasConflict(resp) {
				return nil, "", err
			}

			if utils.ResponseWasConflict(resp) && strings.Contains(err.Error(), "ShareBeingDeleted") {
				return nil, "waitingOnDelete", nil
			}
		}

		return "succeeded", "succeeded", nil
	}
}
