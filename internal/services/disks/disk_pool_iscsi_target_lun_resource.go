// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package disks

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = DiskPoolIscsiTargetLunModel{}
var _ sdk.ResourceWithDeprecationAndNoReplacement = DiskPoolIscsiTargetLunModel{}

type DiskPoolIscsiTargetLunModel struct {
	IscsiTargetId           string `tfschema:"iscsi_target_id"`
	ManagedDiskAttachmentId string `tfschema:"disk_pool_managed_disk_attachment_id"`
	Name                    string `tfschema:"name"`
	Lun                     int64  `tfschema:"lun"`
}

func (DiskPoolIscsiTargetLunModel) DeprecationMessage() string {
	return "The `azurerm_disk_pool_iscsi_target_lun` resource is deprecated and will be removed in v4.0 of the AzureRM Provider."
}

func (d DiskPoolIscsiTargetLunModel) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"iscsi_target_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iscsitargets.ValidateIscsiTargetID,
		},
		"disk_pool_managed_disk_attachment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DiskPoolManagedDiskAttachment,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 90),
				validation.StringMatch(regexp.MustCompile(`^[A-Za-z\d.\-_]*[A-Za-z\d]$`), "supported characters include [0-9A-Za-z-_.]; name should end with an alphanumeric character"),
			),
		},
	}
}

func (d DiskPoolIscsiTargetLunModel) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"lun": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
	}
}

func (d DiskPoolIscsiTargetLunModel) ModelObject() interface{} {
	return &DiskPoolIscsiTargetLunModel{}
}

func (d DiskPoolIscsiTargetLunModel) ResourceType() string {
	return "azurerm_disk_pool_iscsi_target_lun"
}

func (d DiskPoolIscsiTargetLunModel) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			m := DiskPoolIscsiTargetLunModel{}
			err := metadata.Decode(&m)
			if err != nil {
				return err
			}
			iscsiTargetId, err := iscsitargets.ParseIscsiTargetID(m.IscsiTargetId)
			if err != nil {
				return err
			}
			attachmentId, err := parse.DiskPoolManagedDiskAttachmentID(m.ManagedDiskAttachmentId)
			if err != nil {
				return err
			}
			id := parse.NewDiskPoolIscsiTargetLunId(*iscsiTargetId, attachmentId.ManagedDiskId)

			locks.ByID(iscsiTargetId.ID())
			defer locks.UnlockByID(iscsiTargetId.ID())

			client := metadata.Client.Disks.DisksPoolIscsiTargetClient
			resp, err := client.Get(ctx, *iscsiTargetId)
			if err != nil {
				return fmt.Errorf("retrieving %q: %+v", iscsiTargetId, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("malformed Iscsi Target response %q : %+v", iscsiTargetId.ID(), resp)
			}
			var luns []iscsitargets.IscsiLun
			if resp.Model.Properties.Luns != nil {
				luns = *resp.Model.Properties.Luns
			}
			for _, lun := range luns {
				if lun.ManagedDiskAzureResourceId == attachmentId.ManagedDiskId.ID() {
					return metadata.ResourceRequiresImport(d.ResourceType(), id)
				}
			}

			luns = append(luns, iscsitargets.IscsiLun{
				ManagedDiskAzureResourceId: attachmentId.ManagedDiskId.ID(),
				Name:                       m.Name,
			})
			sort.Slice(luns, func(i, j int) bool {
				return luns[i].ManagedDiskAzureResourceId < luns[j].ManagedDiskAzureResourceId
			})
			patch := iscsitargets.IscsiTargetUpdate{
				Properties: iscsitargets.IscsiTargetUpdateProperties{
					Luns: &luns,
				},
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			err = d.RetryError(time.Until(deadline), "waiting for creation DisksPool iscsi target", id.ID(), func() error {
				return client.UpdateThenPoll(ctx, *iscsiTargetId, patch)
			})
			if err != nil {
				return err
			}
			metadata.SetID(id)
			return nil
		},
	}
}

func (d DiskPoolIscsiTargetLunModel) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.IscsiTargetLunID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			iscsiTargetId := id.IscsiTargetId

			client := metadata.Client.Disks.DisksPoolIscsiTargetClient
			resp, err := client.Get(ctx, iscsiTargetId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %+v", iscsiTargetId, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("malformed Iscsi Target response %q : %+v", iscsiTargetId.ID(), resp)
			}
			if resp.Model.Properties.Luns == nil {
				return metadata.MarkAsGone(id)
			}

			for _, lun := range *resp.Model.Properties.Luns {
				if lun.ManagedDiskAzureResourceId == id.ManagedDiskId.ID() {
					diskPoolId := diskpools.NewDiskPoolID(iscsiTargetId.SubscriptionId, iscsiTargetId.ResourceGroupName, iscsiTargetId.DiskPoolName)
					diskId, err := disks.ParseDiskIDInsensitively(lun.ManagedDiskAzureResourceId)
					if err != nil {
						return fmt.Errorf("invalid managed disk id in iscsi target response %q : %q", iscsiTargetId.ID(), lun.ManagedDiskAzureResourceId)
					}
					attachmentId := parse.NewDiskPoolManagedDiskAttachmentId(diskPoolId, *diskId)
					if lun.Lun == nil {
						return fmt.Errorf("malformed Iscsi Target response %q : %+v", iscsiTargetId.ID(), resp)
					}
					m := DiskPoolIscsiTargetLunModel{
						IscsiTargetId:           iscsiTargetId.ID(),
						Lun:                     *lun.Lun,
						ManagedDiskAttachmentId: attachmentId.ID(),
						Name:                    lun.Name,
					}
					return metadata.Encode(&m)
				}
			}

			return metadata.MarkAsGone(id)
		},
	}
}

func (d DiskPoolIscsiTargetLunModel) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			m := DiskPoolIscsiTargetLunModel{}
			err := metadata.Decode(&m)
			if err != nil {
				return err
			}
			iscsiTargetId, err := iscsitargets.ParseIscsiTargetID(m.IscsiTargetId)
			if err != nil {
				return err
			}
			attachmentId, err := parse.DiskPoolManagedDiskAttachmentID(m.ManagedDiskAttachmentId)
			if err != nil {
				return err
			}
			id := parse.NewDiskPoolIscsiTargetLunId(*iscsiTargetId, attachmentId.ManagedDiskId)

			locks.ByID(iscsiTargetId.ID())
			defer locks.UnlockByID(iscsiTargetId.ID())

			client := metadata.Client.Disks.DisksPoolIscsiTargetClient
			resp, err := client.Get(ctx, *iscsiTargetId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %q: %+v", iscsiTargetId, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("malformed Iscsi Target response %q : %+v", iscsiTargetId.ID(), resp)
			}
			if resp.Model.Properties.Luns == nil {
				return nil
			}
			luns := make([]iscsitargets.IscsiLun, 0)
			for _, lun := range *resp.Model.Properties.Luns {
				if lun.ManagedDiskAzureResourceId != id.ManagedDiskId.ID() {
					luns = append(luns, lun)
				}
			}
			sort.Slice(luns, func(i, j int) bool {
				return luns[i].ManagedDiskAzureResourceId < luns[j].ManagedDiskAzureResourceId
			})
			patch := iscsitargets.IscsiTargetUpdate{
				Properties: iscsitargets.IscsiTargetUpdateProperties{
					Luns: &luns,
				},
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			return d.RetryError(time.Until(deadline), "waiting for delete DisksPool iscsi target lun", id.ID(), func() error {
				return client.UpdateThenPoll(ctx, *iscsiTargetId, patch)
			})
		},
	}
}

func (d DiskPoolIscsiTargetLunModel) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DiskPoolIscsiTargetLunId
}

func (DiskPoolIscsiTargetLunModel) RetryError(timeout time.Duration, action string, id string, retryFunc func() error) error {
	return pluginsdk.Retry(timeout, func() *pluginsdk.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		// according to https://docs.microsoft.com/en-us/azure/virtual-machines/disks-pools-troubleshoot#common-failure-codes-when-enabling-iscsi-on-disk-pools the errors below are retryable.
		retryableErrors := []string{
			"GoalStateApplicationTimeoutError",
			"OngoingOperationInProgress",
		}
		for _, retryableError := range retryableErrors {
			if strings.Contains(err.Error(), retryableError) {
				return pluginsdk.RetryableError(fmt.Errorf("%s %s: %+v", action, id, err))
			}
		}
		return pluginsdk.NonRetryableError(fmt.Errorf("%s %s: %+v", action, id, err))
	})
}
