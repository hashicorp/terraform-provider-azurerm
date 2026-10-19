package rules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorrelationFilter struct {
	ContentType           *string            `json:"contentType,omitempty"`
	CorrelationId         *string            `json:"correlationId,omitempty"`
	Label                 *string            `json:"label,omitempty"`
	MessageId             *string            `json:"messageId,omitempty"`
	Properties            *map[string]string `json:"properties,omitempty"`
	ReplyTo               *string            `json:"replyTo,omitempty"`
	ReplyToSessionId      *string            `json:"replyToSessionId,omitempty"`
	RequiresPreprocessing *bool              `json:"requiresPreprocessing,omitempty"`
	SessionId             *string            `json:"sessionId,omitempty"`
	To                    *string            `json:"to,omitempty"`
}
