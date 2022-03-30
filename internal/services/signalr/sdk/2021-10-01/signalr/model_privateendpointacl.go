package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type PrivateEndpointACL struct {
	Allow *[]SignalRRequestType `json:"allow,omitempty"`
	Deny  *[]SignalRRequestType `json:"deny,omitempty"`
	Name  string                `json:"name"`
}
