package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type UpstreamTemplate struct {
	Auth            *UpstreamAuthSettings `json:"auth,omitempty"`
	CategoryPattern *string               `json:"categoryPattern,omitempty"`
	EventPattern    *string               `json:"eventPattern,omitempty"`
	HubPattern      *string               `json:"hubPattern,omitempty"`
	UrlTemplate     string                `json:"urlTemplate"`
}
