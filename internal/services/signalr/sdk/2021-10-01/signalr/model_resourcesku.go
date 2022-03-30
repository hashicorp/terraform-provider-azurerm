package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ResourceSku struct {
	Capacity *int64          `json:"capacity,omitempty"`
	Family   *string         `json:"family,omitempty"`
	Name     string          `json:"name"`
	Size     *string         `json:"size,omitempty"`
	Tier     *SignalRSkuTier `json:"tier,omitempty"`
}
