// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
)

// retrieveDiskEncryptionSetEncryptionType returns encryption type of the disk encryption set
func retrieveDiskEncryptionSetEncryptionType(ctx context.Context, client *diskencryptionsets.DiskEncryptionSetsClient, diskEncryptionSetId string) (*disks.EncryptionType, error) {
	id, err := diskencryptionsets.ParseDiskEncryptionSetID(diskEncryptionSetId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	var encryptionType *disks.EncryptionType

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.EncryptionType != nil {
			s := props.EncryptionType
			v := disks.EncryptionType(*s)
			encryptionType = &v
		}
	}

	if encryptionType == nil {
		return nil, fmt.Errorf("retrieving %s: EncryptionType was nil", *id)
	}

	return encryptionType, nil
}
