// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
)

const storageDiskPoolManagedDiskAttachmentIdSeparator = "/managedDisks|"

var _ resourceids.Id = DiskPoolManagedDiskAttachmentId{}

type DiskPoolManagedDiskAttachmentId struct {
	DiskPoolId    diskpools.DiskPoolId
	ManagedDiskId disks.DiskId
}

func NewDiskPoolManagedDiskAttachmentId(diskPoolId diskpools.DiskPoolId, managedDiskId disks.DiskId) DiskPoolManagedDiskAttachmentId {
	return DiskPoolManagedDiskAttachmentId{
		DiskPoolId:    diskPoolId,
		ManagedDiskId: managedDiskId,
	}
}

func (id DiskPoolManagedDiskAttachmentId) ID() string {
	return fmt.Sprintf("%s%s%s", id.DiskPoolId.ID(), storageDiskPoolManagedDiskAttachmentIdSeparator, id.ManagedDiskId.ID())
}

func (id DiskPoolManagedDiskAttachmentId) String() string {
	components := []string{
		fmt.Sprintf("Disk Pool %q", id.DiskPoolId.String()),
		fmt.Sprintf("Managed Disk %q", id.ManagedDiskId.String()),
	}
	return fmt.Sprintf("Disk Pool Managed Disk Attachment: %s", strings.Join(components, " / "))
}

func DiskPoolManagedDiskAttachmentID(input string) (*DiskPoolManagedDiskAttachmentId, error) {
	if !strings.Contains(input, storageDiskPoolManagedDiskAttachmentIdSeparator) {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	parts := strings.Split(input, storageDiskPoolManagedDiskAttachmentIdSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}

	poolId, err := diskpools.ParseDiskPoolID(parts[0])
	if poolId == nil {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed disks pool id: %q, %v", poolId.ID(), err)
	}
	diskId, err := disks.ParseDiskID(parts[1])
	if diskId == nil {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed disk id: %q, %v", diskId.ID(), err)
	}
	id := NewDiskPoolManagedDiskAttachmentId(*poolId, *diskId)
	return &id, nil
}
