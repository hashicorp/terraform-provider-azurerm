package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SignalRUsageName struct {
	LocalizedValue *string `json:"localizedValue,omitempty"`
	Value          *string `json:"value,omitempty"`
}
