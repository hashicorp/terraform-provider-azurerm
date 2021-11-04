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
	DisksPoolId   string //TODO: use strong type id
	ManagedDiskId string
}

// //TODO: use strong type id as parameter
func NewStorageDisksPoolManagedDiskAttachmentId(diskPoolId, managedDiskId string) StorageDisksPoolManagedDiskAttachmentId {
	return StorageDisksPoolManagedDiskAttachmentId{
		DisksPoolId:   diskPoolId,
		ManagedDiskId: managedDiskId,
	}
}

func (d StorageDisksPoolManagedDiskAttachmentId) ID() string {
	return fmt.Sprintf("%s%s%s", d.DisksPoolId, storageDiskPoolManagedDiskAttachmentIdSeparator, d.ManagedDiskId)
}

func StorageDisksPoolManagedDiskAttachmentID(input string) (*StorageDisksPoolManagedDiskAttachmentId, error) {
	if !strings.Contains(input, storageDiskPoolManagedDiskAttachmentIdSeparator) {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	parts := strings.Split(input, storageDiskPoolManagedDiskAttachmentIdSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	poolId := parts[0]
	diskId := parts[1]

	if _, err := StorageDisksPoolID(poolId); err != nil {
		return nil, fmt.Errorf("malformed disks pool id: %q, %v", poolId, err)
	}
	if _, err := computeParse.ManagedDiskID(diskId); err != nil {
		return nil, fmt.Errorf("malformed disk id: %q, %v", diskId, err)
	}
	id := NewStorageDisksPoolManagedDiskAttachmentId(poolId, diskId)
	return &id, nil
}
