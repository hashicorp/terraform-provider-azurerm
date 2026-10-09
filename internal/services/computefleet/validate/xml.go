// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"encoding/xml"
	"fmt"
)

func Xml(i interface{}, k string) (warnings []string, errors []error) {
	v := i.(string)
	if err := xml.Unmarshal([]byte(v), new(interface{})); err != nil {
		errors = append(errors, fmt.Errorf("%q is not valid XML: %+v", k, err))
	}

	return
}
