package networkrulesets

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type NWRuleSetIpRules struct {
	Action *NetworkRuleIPAction `json:"action,omitempty"`
	IpMask *string              `json:"ipMask,omitempty"`
}
