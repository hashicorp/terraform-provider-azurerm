// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/recordsets"
)

func ValidateRecordTypeID(recordType recordsets.RecordType) func(interface{}, string) ([]string, []error) {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		parsed, err := recordsets.ParseRecordTypeID(input.(string))
		if err != nil {
			errors = append(errors, err)
			return
		}
		if parsed.RecordType != recordType {
			errors = append(errors, fmt.Errorf("this resource only supports '%q' records", recordType))
			return
		}
		return
	}
}
