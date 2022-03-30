package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type PrivateLinkServiceConnectionState struct {
	ActionsRequired *string                             `json:"actionsRequired,omitempty"`
	Description     *string                             `json:"description,omitempty"`
	Status          *PrivateLinkServiceConnectionStatus `json:"status,omitempty"`
}
