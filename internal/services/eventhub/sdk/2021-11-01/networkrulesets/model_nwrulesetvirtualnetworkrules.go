package networkrulesets

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type NWRuleSetVirtualNetworkRules struct {
	IgnoreMissingVnetServiceEndpoint *bool   `json:"ignoreMissingVnetServiceEndpoint,omitempty"`
	Subnet                           *Subnet `json:"subnet,omitempty"`
}
