package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SignalRNetworkACLs struct {
	DefaultAction    *ACLAction            `json:"defaultAction,omitempty"`
	PrivateEndpoints *[]PrivateEndpointACL `json:"privateEndpoints,omitempty"`
	PublicNetwork    *NetworkACL           `json:"publicNetwork,omitempty"`
}
