package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SignalRFeature struct {
	Flag       FeatureFlags       `json:"flag"`
	Properties *map[string]string `json:"properties,omitempty"`
	Value      string             `json:"value"`
}
