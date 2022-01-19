package profiles

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)



type MonitorConfig struct {
	CustomHeaders *[] `json:"customHeaders,omitempty"`
	ExpectedStatusCodeRanges *[] `json:"expectedStatusCodeRanges,omitempty"`
	IntervalInSeconds *int64 `json:"intervalInSeconds,omitempty"`
	Path *string `json:"path,omitempty"`
	Port *int64 `json:"port,omitempty"`
	ProfileMonitorStatus *ProfileMonitorStatus `json:"profileMonitorStatus,omitempty"`
	Protocol *MonitorProtocol `json:"protocol,omitempty"`
	TimeoutInSeconds *int64 `json:"timeoutInSeconds,omitempty"`
	ToleratedNumberOfFailures *int64 `json:"toleratedNumberOfFailures,omitempty"`
}







