package loadtests

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type LoadTestResource struct {
	Id         *string                  `json:"id,omitempty"`
	Identity   *identity.SystemAssigned `json:"identity,omitempty"`
	Location   string                   `json:"location"`
	Name       *string                  `json:"name,omitempty"`
	Properties *LoadTestProperties      `json:"properties,omitempty"`
	SystemData *SystemData              `json:"systemData,omitempty"`
	Tags       *map[string]string       `json:"tags,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
