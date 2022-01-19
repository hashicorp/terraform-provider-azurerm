package loadtests

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type LoadTestResourcePatchRequestBody struct {
	Identity   *identity.SystemAssigned                    `json:"identity,omitempty"`
	Properties *LoadTestResourcePatchRequestBodyProperties `json:"properties,omitempty"`
	Tags       *interface{}                                `json:"tags,omitempty"`
}
