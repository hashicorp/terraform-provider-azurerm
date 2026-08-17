package sessionhost

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2025-10-10"

func userAgent() string {
	return "hashicorp/go-azure-sdk/sessionhost/2025-10-10"
}

func AzureAPIVersion() string {
	return defaultApiVersion
}
