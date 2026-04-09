// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package odata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FromResponse parses a http.Response and returns an unmarshalled OData
// If no odata is present in the response, or the content type is invalid, returns nil
func FromResponse(resp *http.Response) (*OData, error) {
	if resp == nil {
		return nil, nil
	}

	var o OData

	// Check for json content before looking for odata metadata
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "application/json") {
		// Read the response body and close it
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Always reassign the response body
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

		if err != nil {
			return nil, fmt.Errorf("could not read response body: %s", err)
		}

		// Unmarshal odata
		if err := json.Unmarshal(respBody, &o); err != nil {
			return nil, err
		}

		return &o, nil
	}

	return nil, nil
}
