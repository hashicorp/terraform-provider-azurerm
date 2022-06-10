package namespacesprivatelinkresources

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type PrivateLinkResourcesListResult struct {
	NextLink *string                `json:"nextLink,omitempty"`
	Value    *[]PrivateLinkResource `json:"value,omitempty"`
}
