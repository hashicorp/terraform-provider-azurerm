package azuresdkhacks

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/encodings"
)

// Hacking the SDK model

var _ encodings.Codec = JpgImage{}

type JpgImage struct {
	KeyFrameInterval *string                  `json:"keyFrameInterval,omitempty"`
	Layers           *[]JpgLayer              `json:"layers,omitempty"`
	Range            *string                  `json:"range,omitempty"`
	SpriteColumn     *int64                   `json:"spriteColumn,omitempty"`
	Start            string                   `json:"start"`
	Step             *string                  `json:"step,omitempty"`
	StretchMode      *encodings.StretchMode   `json:"stretchMode,omitempty"`
	SyncMode         *encodings.VideoSyncMode `json:"syncMode,omitempty"`

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = JpgImage{}

func (s JpgImage) MarshalJSON() ([]byte, error) {
	type wrapper JpgImage
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JpgImage: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JpgImage: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JpgImage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JpgImage: %+v", err)
	}

	return encoded, nil
}
