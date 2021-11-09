package configurationstores

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type ApiKey struct {
	ConnectionString *string `json:"connectionString,omitempty"`
	Id               *string `json:"id,omitempty"`
	LastModified     *string `json:"lastModified,omitempty"`
	Name             *string `json:"name,omitempty"`
	ReadOnly         *bool   `json:"readOnly,omitempty"`
	Value            *string `json:"value,omitempty"`
}

func (o ApiKey) GetLastModifiedAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o ApiKey) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
