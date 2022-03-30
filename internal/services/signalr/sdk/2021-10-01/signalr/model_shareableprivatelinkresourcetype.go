package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ShareablePrivateLinkResourceType struct {
	Name       *string                                 `json:"name,omitempty"`
	Properties *ShareablePrivateLinkResourceProperties `json:"properties,omitempty"`
}
