package account

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
)

type AccountPropertiesSystemData struct {
	CreatedAt          *string             `json:"createdAt,omitempty"`
	CreatedBy          *string             `json:"createdBy,omitempty"`
	CreatedByType      *CreatedByType      `json:"createdByType,omitempty"`
	LastModifiedAt     *string             `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string             `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *LastModifiedByType `json:"lastModifiedByType,omitempty"`
}

func (o AccountPropertiesSystemData) GetCreatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o AccountPropertiesSystemData) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o AccountPropertiesSystemData) GetLastModifiedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.LastModifiedAt, "2006-01-02T15:04:05Z07:00")
}

func (o AccountPropertiesSystemData) SetLastModifiedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedAt = &formatted
}
