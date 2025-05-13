package data

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	"github.com/spf13/afero"
)

var serviceDirPattern = "%s/internal/services/%s"

type Service struct {
	Name           string
	Path           string
	APIsbyResource map[string][]API
}

func NewService(fs afero.Fs, providerDir string, service any, serviceName string) (*Service, error) {
	labelFunc := func(s string) string {
		return strings.ReplaceAll(strings.TrimPrefix(s, "service/"), "-", "")
	}
	nameFunc := func(s string) string {
		var result []rune
		for _, r := range s {
			if unicode.IsLetter(r) {
				result = append(result, r)
			}
		}

		return strings.ToLower(string(result))
	}

	// Check if serviceName exists in ServiceFolderWorkaround
	if n, ok := ServiceFolderWorkaround[serviceName]; ok {
		serviceName = n
	}
	names := make([]string, 0)

	// ideally all service registrations have `AssociatedGitHubLabel` or some other function that tracks the service folder/internal name
	// e.g. "mssql" rather than the return of `Name()`: `Microsoft SQL Server / Azure SQL`
	// but for now this works
	switch s := service.(type) {
	case sdk.UntypedServiceRegistrationWithAGitHubLabel:
		names = append(names, nameFunc(serviceName), labelFunc(s.AssociatedGitHubLabel()))
	case sdk.TypedServiceRegistrationWithAGitHubLabel:
		names = append(names, nameFunc(serviceName), labelFunc(s.AssociatedGitHubLabel()))
	case sdk.UntypedServiceRegistration, sdk.FrameworkTypedServiceRegistration, sdk.TypedServiceRegistration:
		names = append(names, nameFunc(serviceName))
	default:
		return nil, fmt.Errorf("unexpected service type `%T`", s)
	}

	for _, n := range names {
		path := fmt.Sprintf(serviceDirPattern, providerDir, n)
		if util.DirExists(fs, path) {
			return &Service{
				Name: n,
				Path: path,
			}, nil
		}
	}

	return nil, errors.New("no service directory found")
}
