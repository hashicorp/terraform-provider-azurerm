package graph

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
)

type ObjectSubResourceId struct {
	objectId string
	subId    string
	Type     string
}

func (id ObjectSubResourceId) String() string {
	return id.objectId + "/" + id.Type + "/" + id.subId
}

func ParseObjectSubResourceId(idString, expectedType string) (ObjectSubResourceId, error) {
	parts := strings.Split(idString, "/")
	if len(parts) != 3 {
		return ObjectSubResourceId{}, fmt.Errorf("Object Resource ID should be in the format {objectId}/{keyId} - but got %q", idString)
	}

	id := ObjectSubResourceId{
		objectId: parts[0],
		Type:     parts[1],
		subId:    parts[2],
	}

	if _, err := uuid.ParseUUID(id.objectId); err != nil {
		return ObjectSubResourceId{}, fmt.Errorf("Object ID isn't a valid UUID (%q): %+v", id.objectId, err)
	}

	if id.Type == "" {
		return ObjectSubResourceId{}, fmt.Errorf("Type in {objectID}/{type}/{subID} should not blank")
	}

	if id.Type != expectedType {
		return ObjectSubResourceId{}, fmt.Errorf("Type in {objectID}/{type}/{subID} was expected to be %s, got %s", expectedType, parts[2])
	}

	if _, err := uuid.ParseUUID(id.subId); err != nil {
		return ObjectSubResourceId{}, fmt.Errorf("Object Sub Resource ID isn't a valid UUID (%q): %+v", id.subId, err)
	}

	return id, nil

}

func ObjectSubResourceIdFrom(objectId, typeId, subId string) ObjectSubResourceId {
	return ObjectSubResourceId{
		objectId: objectId,
		Type:     typeId,
		subId:    subId,
	}
}
