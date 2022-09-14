package tenants

import "fmt"

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2021-04-01-preview"

func userAgent() string {
	return fmt.Sprintf("hashicorp/go-azure-sdk/tenants/%s", defaultApiVersion)
}
