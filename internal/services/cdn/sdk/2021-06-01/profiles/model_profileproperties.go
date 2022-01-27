package profiles

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ProfileProperties struct {
	FrontDoorId                  *string                            `json:"frontDoorId,omitempty"`
	Identity                     *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	OriginResponseTimeoutSeconds *int64                             `json:"originResponseTimeoutSeconds,omitempty"`
	ProvisioningState            *string                            `json:"provisioningState,omitempty"`
	ResourceState                *ProfileResourceState              `json:"resourceState,omitempty"`
}
