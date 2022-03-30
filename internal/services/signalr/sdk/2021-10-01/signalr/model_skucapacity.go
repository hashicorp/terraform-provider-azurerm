package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SkuCapacity struct {
	AllowedValues *[]int64   `json:"allowedValues,omitempty"`
	Default       *int64     `json:"default,omitempty"`
	Maximum       *int64     `json:"maximum,omitempty"`
	Minimum       *int64     `json:"minimum,omitempty"`
	ScaleType     *ScaleType `json:"scaleType,omitempty"`
}
