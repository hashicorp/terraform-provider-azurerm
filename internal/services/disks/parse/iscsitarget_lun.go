// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
)

const iscsiTargetLunSeparator = "/lun|"

var _ resourceids.Id = DiskPoolIscsiTargetLunId{}

type DiskPoolIscsiTargetLunId struct {
	IscsiTargetId iscsitargets.IscsiTargetId
	ManagedDiskId disks.DiskId
}

func NewDiskPoolIscsiTargetLunId(iscsiTargetId iscsitargets.IscsiTargetId, managedDiskId disks.DiskId) DiskPoolIscsiTargetLunId {
	return DiskPoolIscsiTargetLunId{
		IscsiTargetId: iscsiTargetId,
		ManagedDiskId: managedDiskId,
	}
}

func (id DiskPoolIscsiTargetLunId) ID() string {
	return fmt.Sprintf("%s%s%s", id.IscsiTargetId.ID(), iscsiTargetLunSeparator, id.ManagedDiskId.ID())
}

func (id DiskPoolIscsiTargetLunId) String() string {
	components := []string{
		fmt.Sprintf("Iscsi Target %q", id.IscsiTargetId.String()),
		fmt.Sprintf("Managed Disk %q", id.ManagedDiskId.String()),
	}
	return fmt.Sprintf("Disk Pool Iscsi Target Lun: %s", strings.Join(components, " / "))
}

func IscsiTargetLunID(input string) (*DiskPoolIscsiTargetLunId, error) {
	if !strings.Contains(input, iscsiTargetLunSeparator) {
		return nil, fmt.Errorf("malformed iscsi target lun id:%q", input)
	}
	parts := strings.Split(input, iscsiTargetLunSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed iscsi target lun id:%q", input)
	}

	iscsiTargetId, err := iscsitargets.ParseIscsiTargetID(parts[0])
	if iscsiTargetId == nil {
		return nil, fmt.Errorf("malformed iscsi target lun id:%q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed iscsi target id: %q, %v", input, err)
	}
	managedDiskId, err := disks.ParseDiskID(parts[1])
	if managedDiskId == nil {
		return nil, fmt.Errorf("malformed iscsi target id: %q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed iscsi target id: %q, %v", input, err)
	}
	id := NewDiskPoolIscsiTargetLunId(*iscsiTargetId, *managedDiskId)
	return &id, nil
}
