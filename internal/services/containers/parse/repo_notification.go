// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"
)

type RepositoryNotificationAction string

const (
	RepositoryNotificationActionPush   RepositoryNotificationAction = "push"
	RepositoryNotificationActionDelete RepositoryNotificationAction = "delete"
	RepositoryNotificationActionAny    RepositoryNotificationAction = "*"
)

var allowedRepositoryNotificationActions = []RepositoryNotificationAction{
	RepositoryNotificationActionPush,
	RepositoryNotificationActionDelete,
	RepositoryNotificationActionAny,
}

type RepositoryNotification struct {
	Artifact
	Action RepositoryNotificationAction
}

type Artifact struct {
	Name   string
	Tag    string
	Digest string
}

func ParseRepositoryNotification(v string) (*RepositoryNotification, error) {
	idx := strings.LastIndex(v, ":")
	if idx == -1 {
		return nil, fmt.Errorf(`no separator ":" found`)
	}
	if idx == len(v)-1 {
		return nil, fmt.Errorf(`malformed format: unexpected trailing ":"`)
	}
	artifact, err := parseArtifact(v[:idx])
	if err != nil {
		return nil, fmt.Errorf("parsing artifact %q: %w", v[:idx], err)
	}
	action := RepositoryNotificationAction(v[idx+1:])
	isAllowedAction := false
	for _, a := range allowedRepositoryNotificationActions {
		if a == action {
			isAllowedAction = true
			break
		}
	}
	if !isAllowedAction {
		return nil, fmt.Errorf("invalid action %q found", action)
	}
	return &RepositoryNotification{
		Artifact: *artifact,
		Action:   action,
	}, nil
}

func parseArtifact(v string) (*Artifact, error) {
	if idx := strings.Index(v, "@"); idx != -1 {
		if idx == len(v)-1 {
			return nil, fmt.Errorf(`malformed format: unexpected trailing "@"`)
		}
		return &Artifact{
			Name:   v[:idx],
			Digest: v[idx+1:],
		}, nil
	}
	if idx := strings.Index(v, ":"); idx != -1 {
		if idx == len(v)-1 {
			return nil, fmt.Errorf(`malformed format: unexpected trailing ":"`)
		}
		return &Artifact{
			Name: v[:idx],
			Tag:  v[idx+1:],
		}, nil
	}
	return &Artifact{Name: v}, nil
}

func (r RepositoryNotification) String() string {
	return fmt.Sprintf("%s:%s", r.Artifact, r.Action)
}

func (a Artifact) String() string {
	out := a.Name
	if a.Tag != "" {
		out += ":" + a.Tag
	}
	if a.Digest != "" {
		out += "@" + a.Digest
	}
	return out
}
