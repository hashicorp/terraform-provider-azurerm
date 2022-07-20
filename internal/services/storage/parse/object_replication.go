package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2021-04-01/objectreplicationpolicies"
)

// This is manual for concat two ids are not supported in auto-generation

type ObjectReplicationId struct {
	Src objectreplicationpolicies.ObjectReplicationPoliciesId
	Dst objectreplicationpolicies.ObjectReplicationPoliciesId
}

func NewObjectReplicationID(srcId, dstId objectreplicationpolicies.ObjectReplicationPoliciesId) ObjectReplicationId {
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
	srcId, err := objectreplicationpolicies.ParseObjectReplicationPoliciesID(ids[0])
	if err != nil {
		return nil, err
	}

	dstId, err := objectreplicationpolicies.ParseObjectReplicationPoliciesID(strings.TrimSuffix(ids[1], ";"))
	if err != nil {
		return nil, err
	}

	return &ObjectReplicationId{
		Src: *srcId,
		Dst: *dstId,
	}, nil
}
