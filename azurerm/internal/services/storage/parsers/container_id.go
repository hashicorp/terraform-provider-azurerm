package parsers

import (
	"fmt"
	"net/url"
	"strings"
)

type ContainerID struct {
	ContainerName string
	AccountName   string
}

func ParseContainerID(id string) (*ContainerID, error) {
	if id == "" {
		return nil, fmt.Errorf("`id` is empty")
	}

	uri, err := url.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("Error parsing ID as a URL: %s", err)
	}

	accountName, err := getAccountNameFromEndpoint(uri.Host)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Name: %s", err)
	}

	containerName := strings.TrimPrefix(uri.Path, "/")
	return &ContainerID{
		AccountName:   *accountName,
		ContainerName: containerName,
	}, nil
}
