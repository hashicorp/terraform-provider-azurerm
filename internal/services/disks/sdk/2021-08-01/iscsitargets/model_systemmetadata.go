package iscsitargets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type SystemMetadata struct {
	CreatedAt          *string        `json:"createdAt,omitempty"`
	CreatedBy          *string        `json:"createdBy,omitempty"`
	CreatedByType      *CreatedByType `json:"createdByType,omitempty"`
	LastModifiedAt     *string        `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string        `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

func (o SystemMetadata) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o SystemMetadata) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o SystemMetadata) GetLastModifiedAtAsTime() (*time.Time, error) {
	if o.LastModifiedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedAt, "2006-01-02T15:04:05Z07:00")
}

func (o SystemMetadata) SetLastModifiedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedAt = &formatted
}
