package account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type AccountUpdateParameters struct {
	Identity   *identity.SystemOrUserAssignedMap `json:"identity,omitempty"`
	Properties *AccountProperties                `json:"properties,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
}
