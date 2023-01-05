package assetsandassetfilters

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackBase = AudioTrack{}

type AudioTrack struct {
	BitRate      *int64        `json:"bitRate,omitempty"`
	DashSettings *DashSettings `json:"dashSettings,omitempty"`
	DisplayName  *string       `json:"displayName,omitempty"`
	FileName     *string       `json:"fileName,omitempty"`
	HlsSettings  *HlsSettings  `json:"hlsSettings,omitempty"`
	LanguageCode *string       `json:"languageCode,omitempty"`
	Mpeg4TrackId *int64        `json:"mpeg4TrackId,omitempty"`

	// Fields inherited from TrackBase
}

var _ json.Marshaler = AudioTrack{}

func (s AudioTrack) MarshalJSON() ([]byte, error) {
	type wrapper AudioTrack
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AudioTrack: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AudioTrack: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.AudioTrack"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AudioTrack: %+v", err)
	}

	return encoded, nil
}
