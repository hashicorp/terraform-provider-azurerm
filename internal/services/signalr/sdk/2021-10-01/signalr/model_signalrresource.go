package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SignalRResource struct {
	Id         *string                           `json:"id,omitempty"`
	Identity   *identity.SystemOrUserAssignedMap `json:"identity,omitempty"`
	Kind       *ServiceKind                      `json:"kind,omitempty"`
	Location   *string                           `json:"location,omitempty"`
	Name       *string                           `json:"name,omitempty"`
	Properties *SignalRProperties                `json:"properties,omitempty"`
	Sku        *ResourceSku                      `json:"sku,omitempty"`
	SystemData *SystemData                       `json:"systemData,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
	Type       *string                           `json:"type,omitempty"`
}
