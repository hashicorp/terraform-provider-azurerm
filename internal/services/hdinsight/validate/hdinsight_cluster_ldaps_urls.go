// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"net/url"
)

func HDInsightClusterLdapsUrls(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	ldapsUrl, err := url.Parse(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %q", k, v))
		return
	}

	if ldapsUrl.Scheme != "ldaps" {
		errors = append(errors, fmt.Errorf(`%s should start with "ldaps://"`, k))
		return
	}

	return
}
