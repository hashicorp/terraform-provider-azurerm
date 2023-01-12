package hybridkubernetes

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func ExpandIdentity(input []IdentityModel) (*identity.SystemAssigned, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("error: Identity should be defined")
	}

	return &identity.SystemAssigned{
		Type: identity.TypeSystemAssigned,
	}, nil
}

func FlattenIdentity(input *identity.SystemAssigned) ([]IdentityModel, error) {
	if input == nil {
		return nil, fmt.Errorf("error: Identity is missing")
	}

	return []IdentityModel{
		{Type: string(input.Type)},
	}, nil
}
