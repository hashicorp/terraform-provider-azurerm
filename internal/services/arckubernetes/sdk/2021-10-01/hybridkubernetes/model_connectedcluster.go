package hybridkubernetes

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ConnectedCluster struct {
	Id         *string                    `json:"id,omitempty"`
	Identity   identity.SystemAssigned    `json:"identity"`
	Location   string                     `json:"location"`
	Name       *string                    `json:"name,omitempty"`
	Properties ConnectedClusterProperties `json:"properties"`
	SystemData *SystemData                `json:"systemData,omitempty"`
	Tags       *map[string]string         `json:"tags,omitempty"`
	Type       *string                    `json:"type,omitempty"`
}
