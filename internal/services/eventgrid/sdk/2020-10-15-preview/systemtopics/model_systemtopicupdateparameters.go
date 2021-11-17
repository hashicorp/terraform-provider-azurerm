package systemtopics

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type SystemTopicUpdateParameters struct {
	Identity *identity.SystemUserAssignedIdentityMap `json:"identity,omitempty"`
	Tags     *map[string]string                      `json:"tags,omitempty"`
}
