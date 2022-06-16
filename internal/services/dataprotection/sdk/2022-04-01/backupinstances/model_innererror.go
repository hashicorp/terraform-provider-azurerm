package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InnerError struct {
	AdditionalInfo     *map[string]string `json:"additionalInfo,omitempty"`
	Code               *string            `json:"code,omitempty"`
	EmbeddedInnerError *InnerError        `json:"embeddedInnerError,omitempty"`
}
