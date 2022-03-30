package signalr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type PrivateLinkResourceProperties struct {
	GroupId                           *string                             `json:"groupId,omitempty"`
	RequiredMembers                   *[]string                           `json:"requiredMembers,omitempty"`
	RequiredZoneNames                 *[]string                           `json:"requiredZoneNames,omitempty"`
	ShareablePrivateLinkResourceTypes *[]ShareablePrivateLinkResourceType `json:"shareablePrivateLinkResourceTypes,omitempty"`
}
