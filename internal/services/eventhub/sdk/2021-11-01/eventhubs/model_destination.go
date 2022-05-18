package eventhubs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type Destination struct {
	Name       *string                `json:"name,omitempty"`
	Properties *DestinationProperties `json:"properties,omitempty"`
}
