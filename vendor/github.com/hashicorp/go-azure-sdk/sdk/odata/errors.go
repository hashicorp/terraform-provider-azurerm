// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package odata

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	ErrorAddedObjectReferencesAlreadyExist      = "One or more added object references already exist"
	ErrorCannotDeleteOrUpdateEnabledEntitlement = "Permission (scope or role) cannot be deleted or updated unless disabled first"
	ErrorConflictingObjectPresentInDirectory    = "A conflicting object with one or more of the specified property values is present in the directory"
	ErrorNotValidReferenceUpdate                = "Not a valid reference update"
	ErrorPropertyValuesAreInvalid               = "One or more property values specified are invalid"
	ErrorResourceDoesNotExist                   = "Resource '.+' does not exist or one of its queried reference-property objects are not present"
	ErrorRemovedObjectReferencesDoNotExist      = "One or more removed object references do not exist"
	ErrorServicePrincipalAppInOtherTenant       = "When using this permission, the backing application of the service principal being created must in the local tenant"
	ErrorServicePrincipalInvalidAppId           = "The appId '.+' of the service principal does not reference a valid application object"
	ErrorUnknownUnsupportedQuery                = "UnknownError: Unsupported Query"
)

// Error is used to unmarshal an API error message.
type Error struct {
	Code            *string          `json:"code"`
	Date            *string          `json:"date"`
	Message         *string          `json:"-"`
	RawMessage      *json.RawMessage `json:"message"` // sometimes a string, sometimes an object :/
	ClientRequestId *string          `json:"client-request-id"`
	RequestId       *string          `json:"request-id"`

	InnerError *Error `json:"innerError"` // nested errors

	Details *[]struct {
		Code   *string `json:"code"`
		Target *string `json:"target"`
	} `json:"details"`

	Values *[]struct {
		Item  string `json:"item"`
		Value string `json:"value"`
	} `json:"values"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	// Perform unmarshalling using a local type
	type error Error
	var e2 error
	if err := json.Unmarshal(data, &e2); err != nil {
		return err
	}
	*e = Error(e2)

	// Unmarshal the message, which can be a plain string or an object wrapping a message
	if raw := e.RawMessage; raw != nil && len(*raw) > 0 {
		switch string((*raw)[0]) {
		case "\"":
			var s string
			if err := json.Unmarshal(*raw, &s); err != nil {
				return err
			}
			e.Message = &s
		case "{":
			var m map[string]interface{}
			if err := json.Unmarshal(*raw, &m); err != nil {
				return err
			}
			if v, ok := m["value"]; ok {
				if s, ok := v.(string); ok {
					e.Message = &s
				}
			}
		default:
			return fmt.Errorf("unrecognised error message: %#v", string(*raw))
		}
	}
	return nil
}

func (e Error) String() string {
	sl := make([]string, 0)
	if e.Code != nil {
		sl = append(sl, *e.Code)
	}
	if e.Message != nil {
		sl = append(sl, *e.Message)
	}
	if e.InnerError != nil {
		if is := e.InnerError.String(); is != "" {
			sl = append(sl, is)
		}
	}
	return strings.Join(sl, ": ")
}

func (e Error) Match(errorText string) bool {
	re := regexp.MustCompile(errorText)
	return re.MatchString(e.String())
}
