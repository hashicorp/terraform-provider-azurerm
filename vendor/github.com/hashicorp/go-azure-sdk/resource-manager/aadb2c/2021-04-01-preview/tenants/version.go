package tenants

// Copyright IBM Corp. 2023, 2026 All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2021-04-01-preview"

func userAgent() string {
	return "hashicorp/go-azure-sdk/tenants/2021-04-01-preview"
}

func AzureAPIVersion() string {
	return defaultApiVersion
}
