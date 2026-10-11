package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2022-10-01-preview"

func userAgent() string {
	return "hashicorp/go-azure-sdk/dataconnectors/2022-10-01-preview"
}

func AzureAPIVersion() string {
	return defaultApiVersion
}
