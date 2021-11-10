package videoanalyzer

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type ListProvisioningTokenInput struct {
	ExpirationDate string `json:"expirationDate"`
}

func (o ListProvisioningTokenInput) GetExpirationDateAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o ListProvisioningTokenInput) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = formatted
}
