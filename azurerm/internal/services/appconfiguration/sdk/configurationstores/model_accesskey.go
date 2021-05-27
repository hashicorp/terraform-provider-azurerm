package configurationstores

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
)

type AccessKey struct {
	ConnectionString *string `json:"connectionString,omitempty"`
	ID               *string `json:"id,omitempty"`
	LastModified     *string `json:"lastModified,omitempty"`
	Name             *string `json:"name,omitempty"`
	ReadOnly         *bool   `json:"readOnly,omitempty"`
	Value            *string `json:"value,omitempty"`
}

func (o AccessKey) ListLastModifiedAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o AccessKey) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
