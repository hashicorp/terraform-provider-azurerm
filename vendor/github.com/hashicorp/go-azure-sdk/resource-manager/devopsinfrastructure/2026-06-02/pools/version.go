package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2026-06-02"

func userAgent() string {
	return "hashicorp/go-azure-sdk/pools/2026-06-02"
}

func AzureAPIVersion() string {
	return defaultApiVersion
}
