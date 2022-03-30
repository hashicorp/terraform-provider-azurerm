package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type UpstreamAuthSettings struct {
	ManagedIdentity *ManagedIdentitySettings `json:"managedIdentity,omitempty"`
	Type            *UpstreamAuthType        `json:"type,omitempty"`
}
