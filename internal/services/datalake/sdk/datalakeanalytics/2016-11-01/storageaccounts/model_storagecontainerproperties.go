package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type StorageContainerProperties struct {
	LastModifiedTime *string `json:"lastModifiedTime,omitempty"`
}

func (o StorageContainerProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o StorageContainerProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
