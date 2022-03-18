package iscsitargets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

const iscsiTargetLunSeparator = "/lun|"

var _ resourceid.Formatter = DiskPoolIscsiTargetLunId{}

type DiskPoolIscsiTargetLunId struct {
	IscsiTargetId IscsiTargetId
	ManagedDiskId computeParse.ManagedDiskId
}

func NewDiskPoolIscsiTargetLunId(iscsiTargetId IscsiTargetId, managedDiskId computeParse.ManagedDiskId) DiskPoolIscsiTargetLunId {
	return DiskPoolIscsiTargetLunId{
		IscsiTargetId: iscsiTargetId,
		ManagedDiskId: managedDiskId,
	}
}

func (d DiskPoolIscsiTargetLunId) ID() string {
	return fmt.Sprintf("%s%s%s", d.IscsiTargetId.ID(), iscsiTargetLunSeparator, d.ManagedDiskId.ID())
}

func ParseIscsiTargetLunID(input string) (*DiskPoolIscsiTargetLunId, error) {
	if !strings.Contains(input, iscsiTargetLunSeparator) {
		return nil, fmt.Errorf("malformed iscsi target lun id:%q", input)
	}
	parts := strings.Split(input, iscsiTargetLunSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed iscsi target lun id:%q", input)
	}

	iscsiTargetId, err := ParseIscsiTargetID(parts[0])
	if iscsiTargetId == nil {
		return nil, fmt.Errorf("malformed iscsi target lun id:%q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed iscsi target id: %q, %v", input, err)
	}
	managedDiskId, err := computeParse.ManagedDiskID(parts[1])
	if managedDiskId == nil {
		return nil, fmt.Errorf("malformed iscsi target id: %q", input)
	}
	if err != nil {
		return nil, fmt.Errorf("malformed iscsi target id: %q, %v", input, err)
	}
	id := NewDiskPoolIscsiTargetLunId(*iscsiTargetId, *managedDiskId)
	return &id, nil
}
