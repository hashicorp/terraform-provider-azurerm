package systemtopics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SystemTopicUpdateParameters struct {
	Identity *identity.SystemUserAssignedList `json:"identity,omitempty"`
	Tags     *map[string]string               `json:"tags,omitempty"`
}
