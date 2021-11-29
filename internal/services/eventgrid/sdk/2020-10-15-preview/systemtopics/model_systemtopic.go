package systemtopics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SystemTopic struct {
	Id         *string                          `json:"id,omitempty"`
	Identity   *identity.SystemUserAssignedList `json:"identity,omitempty"`
	Location   string                           `json:"location"`
	Name       *string                          `json:"name,omitempty"`
	Properties *SystemTopicProperties           `json:"properties,omitempty"`
	SystemData *SystemData                      `json:"systemData,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
	Type       *string                          `json:"type,omitempty"`
}
