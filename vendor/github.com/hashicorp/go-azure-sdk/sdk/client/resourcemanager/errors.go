// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var _ error = Error{}

type Error struct {
	ActivityId string
	Code       string
	Message    string
	Status     string

	FullHttpBody string
}

func (e Error) Error() string {
	return fmt.Sprintf(`the Azure API returned the following error:

Status: %q
Code: %q
Message: %q
Activity Id: %q

---

API Response:

----[start]----
%s
-----[end]-----
`, e.Status, e.Code, e.Message, e.ActivityId, e.FullHttpBody)
}

// parseErrorFromApiResponse parses the error from the API Response
// into an Error type, which allows for better surfacing of errors
func parseErrorFromApiResponse(response http.Response) (*Error, error) {
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %+v", err)
	}
	response.Body.Close()

	respBody = bytes.TrimPrefix(respBody, []byte("\xef\xbb\xbf"))
	response.Body = io.NopCloser(bytes.NewBuffer(respBody))

	// there's a number of internal Azure error types, we should attempt unmarshalling into each
	// for now we're implementing the simple case we can add to
	var err1 lroErrorType1
	if err = json.Unmarshal(respBody, &err1); err == nil && err1.Id != "" && err1.Error.Code != "" && err1.Error.Message != "" {
		e := Error{
			Code:         err1.Error.Code,
			Message:      err1.Error.Message,
			Status:       err1.Status,
			FullHttpBody: string(respBody),
		}

		// given inconsistencies between different API's, this isn't pretty, but avoids crashes whilst best-efforting it
		for _, info := range err1.Error.AdditionalInfo {
			additionalInfo, ok := info.(map[string]interface{})
			if !ok {
				continue
			}

			typeVal, ok := additionalInfo["type"].(string)
			if !ok {
				continue
			}

			if strings.EqualFold(typeVal, "ActivityId") {
				infoBlock, ok := additionalInfo["info"].(map[string]interface{})
				if !ok {
					continue
				}

				for k, v := range infoBlock {
					if strings.EqualFold(k, "ActivityId") {
						if val, ok := v.(string); ok {
							e.ActivityId = val
						}
					}
				}
			}
		}

		return &e, nil
	}

	var err2 resourceManagerErrorType1
	if err = json.Unmarshal(respBody, &err2); err == nil && err2.Code != "" && err2.Message != "" {
		return &Error{
			Code:         err2.Code,
			Message:      err2.Message,
			Status:       "Unknown",
			FullHttpBody: string(respBody),
		}, nil
	}

	var err3 resourceManagerErrorType2
	if err = json.Unmarshal(respBody, &err3); err == nil && err3.Error.Code != "" && err3.Error.Message != "" {
		activityId := ""
		code := err3.Error.Code
		messages := []string{
			err3.Error.Message,
		}
		for _, v := range err3.Error.Details {
			code = v.Code
			if v.PossibleCauses != "" {
				messages = append(messages, fmt.Sprintf("Possible Causes: %q", v.PossibleCauses))
			}
			if v.RecommendedAction != "" {
				messages = append(messages, fmt.Sprintf("Recommended Action: %q", v.RecommendedAction))
			}
			if v.ActivityId != "" {
				activityId = v.ActivityId
			}
			break
		}
		return &Error{
			ActivityId:   activityId,
			Code:         code,
			Message:      strings.Join(messages, "\n"),
			Status:       "Unknown",
			FullHttpBody: string(respBody),
		}, nil
	}

	var err4 resourceManagerErrorType3
	if err = json.Unmarshal(respBody, &err4); err == nil && err4.Status != "" && err4.Error.Message != "" {
		activityId := ""
		code := err4.Status
		messages := []string{
			err4.Error.Message,
		}
		return &Error{
			ActivityId:   activityId,
			Code:         code,
			Message:      strings.Join(messages, "\n"),
			Status:       err4.Status,
			FullHttpBody: string(respBody),
		}, nil
	}

	return &Error{
		Code:         "internal-error",
		Message:      "Couldn't parse Azure API Response into a friendly error - please see the original HTTP Response for more details (and file a bug so we can fix this!).",
		Status:       "Unknown",
		FullHttpBody: string(respBody),
	}, nil
}

type lroErrorType1 struct {
	Error struct {
		Code           string        `json:"code"`
		Message        string        `json:"message"`
		AdditionalInfo []interface{} `json:"additionalInfo"`
	} `json:"error"`
	Id     string `json:"id"`
	Status string `json:"status"`
}

type resourceManagerErrorType1 struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []struct {
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"details"`
}

type resourceManagerErrorType2 struct {
	Error struct {
		Code    string `json:"code"`
		Details []struct {
			ActivityId        string `json:"activityId"`
			Code              string `json:"code"`
			Message           string `json:"message"`
			PossibleCauses    string `json:"possibleCauses"`
			RecommendedAction string `json:"recommendedAction"`
		} `json:"details"`
		Message string `json:"message"`
	} `json:"error"`
}

type resourceManagerErrorType3 struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ResourceId string `json:"resourceId"`
	Status     string `json:"status"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Error      struct {
		Message string `json:"message"`
	} `json:"error"`
}
