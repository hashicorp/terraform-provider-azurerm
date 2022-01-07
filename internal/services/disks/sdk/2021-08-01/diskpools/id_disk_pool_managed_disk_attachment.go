package diskpools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

const storageDiskPoolManagedDiskAttachmentIdSeparator = "/managedDisks|"

var _ resourceid.Formatter = DiskPoolManagedDiskAttachmentId{}

type DiskPoolManagedDiskAttachmentId struct {
	DiskPoolId    DiskPoolId
	ManagedDiskId computeParse.ManagedDiskId
}

func NewDiskPoolManagedDiskAttachmentId(diskPoolId DiskPoolId, managedDiskId computeParse.ManagedDiskId) DiskPoolManagedDiskAttachmentId {
	return DiskPoolManagedDiskAttachmentId{
		DiskPoolId:    diskPoolId,
		ManagedDiskId: managedDiskId,
	}
}

func (d DiskPoolManagedDiskAttachmentId) ID() string {
	return fmt.Sprintf("%s%s%s", d.DiskPoolId.ID(), storageDiskPoolManagedDiskAttachmentIdSeparator, d.ManagedDiskId.ID())
}

func DiskPoolManagedDiskAttachmentID(input string) (*DiskPoolManagedDiskAttachmentId, error) {
	if !strings.Contains(input, storageDiskPoolManagedDiskAttachmentIdSeparator) {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	parts := strings.Split(input, storageDiskPoolManagedDiskAttachmentIdSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}

	poolId, err := ParseDiskPoolID(parts[0])
	if poolId == nil {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed disks pool id: %q, %v", poolId.ID(), err)
	}
	diskId, err := computeParse.ManagedDiskID(parts[1])
	if diskId == nil {
		return nil, fmt.Errorf("malformed disks pool managed disk attachment id:%q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed disk id: %q, %v", diskId.ID(), err)
	}
	id := NewDiskPoolManagedDiskAttachmentId(*poolId, *diskId)
	return &id, nil
}
