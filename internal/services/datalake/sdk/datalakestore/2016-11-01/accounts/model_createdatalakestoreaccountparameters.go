package accounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type CreateDataLakeStoreAccountParameters struct {
	Identity   *identity.SystemAssigned              `json:"identity,omitempty"`
	Location   string                                `json:"location"`
	Properties *CreateDataLakeStoreAccountProperties `json:"properties,omitempty"`
	Tags       *map[string]string                    `json:"tags,omitempty"`
}
