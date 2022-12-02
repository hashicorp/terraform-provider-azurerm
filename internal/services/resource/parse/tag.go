package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = TagId{}

type TagId struct {
	ResourceId string
	Key        string
}

func NewTagID(resourceId string, tagKey string) TagId {
	return TagId{
		ResourceId: resourceId,
		Key:        tagKey,
	}
}

func (id TagId) String() string {
	return fmt.Sprintf("Tag %s at scope %s", id.Key, id.ResourceId)
}

func (id TagId) ID() string {
	fmtString := "%s|%s"
	return fmt.Sprintf(fmtString, id.ResourceId, id.Key)
}

// TagID parses a Tag ID into an TagID struct
func TagID(input string) (*TagId, error) {
	segments := strings.Split(input, "|")

	id, err := resourceids.ParseAzureResourceID(segments[0])
	if err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	tagId := TagId{
		ResourceId: segments[0],
		Key:        segments[1],
	}

	if tagId.ResourceId == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceId' element")
	}

	if tagId.Key == "" {
		return nil, fmt.Errorf("ID was missing the 'key' element")
	}

	return &tagId, nil
}

func ValidateResourceTagID(input interface{}, _ string) (warnings []string, errors []error) {
	_, err := TagID(input.(string))

	errors = []error{}

	if err != nil {
		errors = append(errors, err)
	}

	return nil, errors
}
