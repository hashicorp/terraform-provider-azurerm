package resource

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ServicesDescription struct {
	Etag       *string                  `json:"etag,omitempty"`
	Id         *string                  `json:"id,omitempty"`
	Identity   *identity.SystemAssigned `json:"identity,omitempty"`
	Kind       Kind                     `json:"kind"`
	Location   string                   `json:"location"`
	Name       *string                  `json:"name,omitempty"`
	Properties *ServicesProperties      `json:"properties,omitempty"`
	SystemData *SystemData              `json:"systemData,omitempty"`
	Tags       *map[string]string       `json:"tags,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
