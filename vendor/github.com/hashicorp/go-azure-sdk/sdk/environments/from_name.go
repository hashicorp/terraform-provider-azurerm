// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package environments

import (
	"fmt"
	"strings"
)

func FromName(name string) (*Environment, error) {
	switch strings.ToLower(name) {
	case "china":
		return AzureChina(), nil

	case "canary":
		return AzurePublicCanary(), nil

	case "global", "public":
		return AzurePublic(), nil

	case "usgovernment", "usgovernmentl4":
		return AzureUSGovernment(), nil

	case "dod", "usgovernmentl5":
		return AzureUSGovernmentL5(), nil
	}

	return nil, fmt.Errorf("no environment was found with the name %q", name)
}
