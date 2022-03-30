package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SignalRUsage struct {
	CurrentValue *int64            `json:"currentValue,omitempty"`
	Id           *string           `json:"id,omitempty"`
	Limit        *int64            `json:"limit,omitempty"`
	Name         *SignalRUsageName `json:"name,omitempty"`
	Unit         *string           `json:"unit,omitempty"`
}
