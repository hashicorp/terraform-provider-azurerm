package parse

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"strings"
)

const storageDiskPoolManagedDiskAttachmentIdSeparator = "/managedDisks|"

var _ resourceid.Formatter = StorageDisksPoolManagedDiskAttachmentId{}

type StorageDisksPoolManagedDiskAttachmentId struct {
	DisksPoolId   StorageDisksPoolId
	ManagedDiskId computeParse.ManagedDiskId
}

func NewStorageDisksPoolManagedDiskAttachmentId(diskPoolId StorageDisksPoolId, managedDiskId computeParse.ManagedDiskId) StorageDisksPoolManagedDiskAttachmentId {
	return StorageDisksPoolManagedDiskAttachmentId{
		DisksPoolId:   diskPoolId,
		ManagedDiskId: managedDiskId,
	}
}

func (d StorageDisksPoolManagedDiskAttachmentId) ID() string {
	return fmt.Sprintf("%s%s%s", d.DisksPoolId.ID(), storageDiskPoolManagedDiskAttachmentIdSeparator, d.ManagedDiskId.ID())
}

func StorageDisksPoolManagedDiskAttachmentID(input string) (*StorageDisksPoolManagedDiskAttachmentId, error) {
	if !strings.Contains(input, storageDiskPoolManagedDiskAttachmentIdSeparator) {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	parts := strings.Split(input, storageDiskPoolManagedDiskAttachmentIdSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}

	poolId, err := StorageDisksPoolID(parts[0])
	if err != nil {
		return nil, fmt.Errorf("malformed disks pool id: %q, %v", poolId.ID(), err)
	}
	diskId, err := computeParse.ManagedDiskID(parts[1])
	if err != nil {
		return nil, fmt.Errorf("malformed disk id: %q, %v", diskId.ID(), err)
	}
	id := NewStorageDisksPoolManagedDiskAttachmentId(*poolId, *diskId)
	return &id, nil
}
