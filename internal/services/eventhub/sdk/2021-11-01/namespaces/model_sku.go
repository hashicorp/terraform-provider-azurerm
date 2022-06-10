package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type Sku struct {
	Capacity *int64   `json:"capacity,omitempty"`
	Name     SkuName  `json:"name"`
	Tier     *SkuTier `json:"tier,omitempty"`
}
