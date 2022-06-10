package eventhubsclustersconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ClusterQuotaConfigurationProperties struct {
	Settings *map[string]string `json:"settings,omitempty"`
}
