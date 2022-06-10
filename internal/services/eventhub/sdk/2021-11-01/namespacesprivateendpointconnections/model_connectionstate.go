package namespacesprivateendpointconnections

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ConnectionState struct {
	Description *string                      `json:"description,omitempty"`
	Status      *PrivateLinkConnectionStatus `json:"status,omitempty"`
}
