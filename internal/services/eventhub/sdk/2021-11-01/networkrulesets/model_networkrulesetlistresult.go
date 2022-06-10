package networkrulesets

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type NetworkRuleSetListResult struct {
	NextLink *string           `json:"nextLink,omitempty"`
	Value    *[]NetworkRuleSet `json:"value,omitempty"`
}
