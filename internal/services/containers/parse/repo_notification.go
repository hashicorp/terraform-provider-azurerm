package parse

import (
	"fmt"
	"strings"
)

type RepositoryNotificationAction string

const (
	RepositoryNotificationActionPush   RepositoryNotificationAction = "push"
	RepositoryNotificationActionDelete                              = "delete"
	RepositoryNotificationActionAny                                 = "*"
)

var allowedRepositoryNotificationActions = []RepositoryNotificationAction{
	RepositoryNotificationActionPush,
	RepositoryNotificationActionDelete,
	RepositoryNotificationActionAny,
}

type repositoryNotification struct {
	Artifact
	Action RepositoryNotificationAction
}

type Artifact struct {
	Name   string
	Tag    string
	Digest string
}

func RepositoryNotification(v string) (*repositoryNotification, error) {
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
	return &repositoryNotification{
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
