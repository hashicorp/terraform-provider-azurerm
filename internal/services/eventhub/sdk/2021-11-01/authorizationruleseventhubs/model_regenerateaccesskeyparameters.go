package authorizationruleseventhubs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type RegenerateAccessKeyParameters struct {
	Key     *string `json:"key,omitempty"`
	KeyType KeyType `json:"keyType"`
}
