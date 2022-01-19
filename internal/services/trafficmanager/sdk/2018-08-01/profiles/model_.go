package profiles

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)



type  struct {
	First *string `json:"first,omitempty"`
	Last *string `json:"last,omitempty"`
	Scope *int64 `json:"scope,omitempty"`
}







