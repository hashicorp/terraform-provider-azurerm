package checknameavailabilitydisasterrecoveryconfigs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type CheckNameAvailabilityResult struct {
	Message       *string            `json:"message,omitempty"`
	NameAvailable *bool              `json:"nameAvailable,omitempty"`
	Reason        *UnavailableReason `json:"reason,omitempty"`
}
