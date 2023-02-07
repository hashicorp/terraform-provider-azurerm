// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemdata

import "encoding/json"

var _ json.Marshaler = &SystemData{}

type SystemData struct {
	CreatedBy          string `json:"createdBy"`
	CreatedByType      string `json:"createdByType"`
	CreatedAt          string `json:"createdAt"`
	LastModifiedBy     string `json:"lastModifiedBy"`
	LastModifiedbyType string `json:"lastModifiedbyType"`
	LastModifiedAt     string `json:"lastModifiedAt"`
}

func (s *SystemData) MarshalJSON() ([]byte, error) {
	// SystemData is a Read Only type. If Systemdata is part of a request some Azure APIs will ignore it,
	// others will return HTTP 400. We're returning nothing on purpose to avoid the error.
	return json.Marshal(nil)
}
