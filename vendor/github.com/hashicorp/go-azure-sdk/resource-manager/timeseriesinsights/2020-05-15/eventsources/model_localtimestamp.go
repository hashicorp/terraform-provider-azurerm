package eventsources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalTimestamp struct {
	Format         *LocalTimestampFormat         `json:"format,omitempty"`
	TimeZoneOffset *LocalTimestampTimeZoneOffset `json:"timeZoneOffset,omitempty"`
}
