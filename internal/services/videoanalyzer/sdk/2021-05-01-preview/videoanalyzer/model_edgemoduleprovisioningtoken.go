package videoanalyzer

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type EdgeModuleProvisioningToken struct {
	ExpirationDate *string `json:"expirationDate,omitempty"`
	Token          *string `json:"token,omitempty"`
}

func (o EdgeModuleProvisioningToken) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o EdgeModuleProvisioningToken) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}
