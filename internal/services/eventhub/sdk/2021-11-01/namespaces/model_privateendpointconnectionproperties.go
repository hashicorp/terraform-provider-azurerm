package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type PrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpoint           `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *ConnectionState           `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *EndPointProvisioningState `json:"provisioningState,omitempty"`
}
