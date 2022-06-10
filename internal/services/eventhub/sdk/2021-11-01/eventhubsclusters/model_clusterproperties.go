package eventhubsclusters

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ClusterProperties struct {
	CreatedAt *string `json:"createdAt,omitempty"`
	MetricId  *string `json:"metricId,omitempty"`
	Status    *string `json:"status,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
}
