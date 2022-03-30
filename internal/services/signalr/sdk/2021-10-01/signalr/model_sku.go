package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type Sku struct {
	Capacity     *SkuCapacity `json:"capacity,omitempty"`
	ResourceType *string      `json:"resourceType,omitempty"`
	Sku          *ResourceSku `json:"sku,omitempty"`
}
