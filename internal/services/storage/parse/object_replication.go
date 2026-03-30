// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/objectreplicationpolicyoperationgroup"
)

// This is manual for concat two ids are not supported in auto-generation

var _ resourceids.Id = ObjectReplicationId{}

type ObjectReplicationId struct {
	Src objectreplicationpolicyoperationgroup.ObjectReplicationPolicyId
	Dst objectreplicationpolicyoperationgroup.ObjectReplicationPolicyId
}

func NewObjectReplicationID(srcId, dstId objectreplicationpolicyoperationgroup.ObjectReplicationPolicyId) ObjectReplicationId {
	return ObjectReplicationId{
		Src: srcId,
		Dst: dstId,
	}
}

func (id ObjectReplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Source %q", id.Src),
		fmt.Sprintf("Destination %q", id.Dst),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Object Replication", segmentsStr)
}

func (id ObjectReplicationId) ID() string {
	return fmt.Sprintf("%s;%s", id.Src.ID(), id.Dst.ID())
}

// ObjectReplicationID parses a ObjectReplication ID into an ObjectReplicationId struct
func ObjectReplicationID(input string) (*ObjectReplicationId, error) {
	ids := strings.Split(input, ";")
	if len(ids) != 2 {
		return nil, fmt.Errorf("storage Object Replication Id is composed as format `sourceId;destinationId`")
	}
	srcId, err := objectreplicationpolicyoperationgroup.ParseObjectReplicationPolicyID(ids[0])
	if err != nil {
		return nil, err
	}

	dstId, err := objectreplicationpolicyoperationgroup.ParseObjectReplicationPolicyID(strings.TrimSuffix(ids[1], ";"))
	if err != nil {
		return nil, err
	}

	return &ObjectReplicationId{
		Src: *srcId,
		Dst: *dstId,
	}, nil
}
