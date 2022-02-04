package accounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type DataLakeAnalyticsAccountPropertiesBasic struct {
	AccountId         *string                         `json:"accountId,omitempty"`
	CreationTime      *string                         `json:"creationTime,omitempty"`
	Endpoint          *string                         `json:"endpoint,omitempty"`
	LastModifiedTime  *string                         `json:"lastModifiedTime,omitempty"`
	ProvisioningState *DataLakeAnalyticsAccountStatus `json:"provisioningState,omitempty"`
	State             *DataLakeAnalyticsAccountState  `json:"state,omitempty"`
}

func (o DataLakeAnalyticsAccountPropertiesBasic) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o DataLakeAnalyticsAccountPropertiesBasic) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o DataLakeAnalyticsAccountPropertiesBasic) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o DataLakeAnalyticsAccountPropertiesBasic) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
