// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"
)

type AppServiceSourceControlTokenId struct {
	Type string
}

const (
	gitHubTokenId       = "/providers/Microsoft.Web/sourceControls/%s"
	gitHubTokenIdPrefix = "/providers/Microsoft.Web/sourceControls/"
)

func (id AppServiceSourceControlTokenId) String() string {
	return fmt.Sprintf("App Service Source Control Token for %s", id.Type)
}

func (id AppServiceSourceControlTokenId) ID() string {
	return fmt.Sprintf(gitHubTokenId, id.Type)
}

func NewAppServiceSourceControlTokenID(input string) AppServiceSourceControlTokenId {
	return AppServiceSourceControlTokenId{
		Type: input,
	}
}

func AppServiceSourceControlTokenID(input string) (*AppServiceSourceControlTokenId, error) {
	if !strings.HasPrefix(input, gitHubTokenIdPrefix) {
		return nil, fmt.Errorf("could not parse source control token ID, expected %s{type}, got %s", gitHubTokenIdPrefix, input)
	}

	return &AppServiceSourceControlTokenId{
		Type: strings.TrimPrefix(input, gitHubTokenIdPrefix),
	}, nil
}
