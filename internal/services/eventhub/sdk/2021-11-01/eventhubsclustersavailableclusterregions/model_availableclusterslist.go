package eventhubsclustersavailableclusterregions

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type AvailableClustersList struct {
	Value *[]AvailableCluster `json:"value,omitempty"`
}
