package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type PrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                          `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpoint                   `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *ProvisioningState                 `json:"provisioningState,omitempty"`
}
