// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

var _ error = &storageError{}

type storageError struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

func (e *storageError) Error() string {
	out := e.Message
	if e.Code != "" {
		out = fmt.Sprintf("%s: %s", e.Code, out)
	}
	return out
}

var _ client.ResponseErrorParser = &ErrorParser{}

// ErrorParser is a custom error parse for Azure Storage errors
type ErrorParser struct{}

func (e *ErrorParser) FromResponse(resp *http.Response) error {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse response body")
	}
	resp.Body.Close()
	respBody = bytes.TrimPrefix(respBody, []byte("\xef\xbb\xbf"))

	res := storageError{}
	if err = xml.Unmarshal(respBody, &res); err != nil {
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	return &res
}
