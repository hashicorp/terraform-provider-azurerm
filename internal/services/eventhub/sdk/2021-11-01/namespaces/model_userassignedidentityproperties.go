package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type UserAssignedIdentityProperties struct {
	UserAssignedIdentity *string `json:"userAssignedIdentity,omitempty"`
}
