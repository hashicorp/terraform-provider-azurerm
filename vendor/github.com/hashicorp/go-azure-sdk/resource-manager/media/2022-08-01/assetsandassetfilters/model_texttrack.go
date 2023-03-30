package assetsandassetfilters

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TrackBase = TextTrack{}

type TextTrack struct {
	DisplayName      *string      `json:"displayName,omitempty"`
	FileName         *string      `json:"fileName,omitempty"`
	HlsSettings      *HlsSettings `json:"hlsSettings,omitempty"`
	LanguageCode     *string      `json:"languageCode,omitempty"`
	PlayerVisibility *Visibility  `json:"playerVisibility,omitempty"`

	// Fields inherited from TrackBase
}

var _ json.Marshaler = TextTrack{}

func (s TextTrack) MarshalJSON() ([]byte, error) {
	type wrapper TextTrack
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TextTrack: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TextTrack: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.TextTrack"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TextTrack: %+v", err)
	}

	return encoded, nil
}
