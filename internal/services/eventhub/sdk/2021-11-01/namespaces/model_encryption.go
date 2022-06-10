package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type Encryption struct {
	KeySource                       *KeySource            `json:"keySource,omitempty"`
	KeyVaultProperties              *[]KeyVaultProperties `json:"keyVaultProperties,omitempty"`
	RequireInfrastructureEncryption *bool                 `json:"requireInfrastructureEncryption,omitempty"`
}
