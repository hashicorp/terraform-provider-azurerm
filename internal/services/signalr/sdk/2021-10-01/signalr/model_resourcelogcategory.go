package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ResourceLogCategory struct {
	Enabled *string `json:"enabled,omitempty"`
	Name    *string `json:"name,omitempty"`
}
