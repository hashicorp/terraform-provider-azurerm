package disasterrecoveryconfigs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ArmDisasterRecoveryProperties struct {
	AlternateName                     *string               `json:"alternateName,omitempty"`
	PartnerNamespace                  *string               `json:"partnerNamespace,omitempty"`
	PendingReplicationOperationsCount *int64                `json:"pendingReplicationOperationsCount,omitempty"`
	ProvisioningState                 *ProvisioningStateDR  `json:"provisioningState,omitempty"`
	Role                              *RoleDisasterRecovery `json:"role,omitempty"`
}
