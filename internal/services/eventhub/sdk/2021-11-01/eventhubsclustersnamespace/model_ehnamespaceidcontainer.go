package eventhubsclustersnamespace

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type EHNamespaceIdContainer struct {
	Id *string `json:"id,omitempty"`
}
