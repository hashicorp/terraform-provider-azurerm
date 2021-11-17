package usagedetails

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type DownloadURL struct {
	DownloadUrl *string `json:"downloadUrl,omitempty"`
	ValidTill   *string `json:"validTill,omitempty"`
}

func (o DownloadURL) GetValidTillAsTime() (*time.Time, error) {
	if o.ValidTill == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ValidTill, "2006-01-02T15:04:05Z07:00")
}

func (o DownloadURL) SetValidTillAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ValidTill = &formatted
}
