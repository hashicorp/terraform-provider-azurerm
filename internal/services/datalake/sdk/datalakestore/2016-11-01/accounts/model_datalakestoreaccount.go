package accounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type DataLakeStoreAccount struct {
	Id         *string                         `json:"id,omitempty"`
	Identity   *identity.SystemAssigned        `json:"identity,omitempty"`
	Location   *string                         `json:"location,omitempty"`
	Name       *string                         `json:"name,omitempty"`
	Properties *DataLakeStoreAccountProperties `json:"properties,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
	Type       *string                         `json:"type,omitempty"`
}
