package networkrulesets

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type NetworkRuleSetProperties struct {
	DefaultAction               *DefaultAction                  `json:"defaultAction,omitempty"`
	IpRules                     *[]NWRuleSetIpRules             `json:"ipRules,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccessFlag        `json:"publicNetworkAccess,omitempty"`
	TrustedServiceAccessEnabled *bool                           `json:"trustedServiceAccessEnabled,omitempty"`
	VirtualNetworkRules         *[]NWRuleSetVirtualNetworkRules `json:"virtualNetworkRules,omitempty"`
}
