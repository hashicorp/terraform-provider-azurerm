package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SignalRKeys struct {
	PrimaryConnectionString   *string `json:"primaryConnectionString,omitempty"`
	PrimaryKey                *string `json:"primaryKey,omitempty"`
	SecondaryConnectionString *string `json:"secondaryConnectionString,omitempty"`
	SecondaryKey              *string `json:"secondaryKey,omitempty"`
}
