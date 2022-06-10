package eventhubs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type CaptureDescription struct {
	Destination       *Destination                `json:"destination,omitempty"`
	Enabled           *bool                       `json:"enabled,omitempty"`
	Encoding          *EncodingCaptureDescription `json:"encoding,omitempty"`
	IntervalInSeconds *int64                      `json:"intervalInSeconds,omitempty"`
	SizeLimitInBytes  *int64                      `json:"sizeLimitInBytes,omitempty"`
	SkipEmptyArchives *bool                       `json:"skipEmptyArchives,omitempty"`
}
