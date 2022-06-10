package eventhubsclusters

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ClusterSku struct {
	Capacity *int64         `json:"capacity,omitempty"`
	Name     ClusterSkuName `json:"name"`
}
