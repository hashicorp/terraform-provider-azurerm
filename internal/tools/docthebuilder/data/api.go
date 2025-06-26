package data

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	log "github.com/sirupsen/logrus"
)

var (
	apiProviderRegex = regexp.MustCompile(`providers/([\w.]*)/`)
	apiVersionRegex  = regexp.MustCompile(`/(\d{4}-\d{2}-\d{2}(?:-preview)?)/`)
)

type API struct {
	Name     string
	URL      string // TODO: not currently used, Resource-manager apis: https://learn.microsoft.com/en-us/rest/api/<Name> -- pattern does not work for all
	Versions []string
}

func (a API) String() string {
	return fmt.Sprintf("* `%s`: %s", a.Name, strings.Join(a.Versions, ", "))
}

func methodsToAPIs(methods []sdkMethod) []API {
	apis := make(map[string]map[string]struct{})
	result := make([]API, 0)

	debugLog := func(m sdkMethod, msg string) {
		log.WithFields(log.Fields{
			"api_path": m.APIPath,
			"method":   m.MethodName,
			"package":  m.Pkg,
		}).Debug(strings.TrimSpace(msg) + " - skipping...")
	}

	for _, m := range methods {
		if m.APIPath == "" {
			// most of these will be methods like `ID` or `String` that we don't care about (thus skipping the logging of these)
			// however there are some genuine API calls in here that are not being parsed correctly
			if m.MethodName != "ID" && m.MethodName != "String" {
				debugLog(m, "sdkMethod object contained an empty API path")
			}
			continue
		}

		matches := apiProviderRegex.FindStringSubmatch(m.APIPath)
		if len(matches) != 2 {
			// TODO: handle cases without `provider/<>`
			/*
				e.g.
				/subscriptions/%s
				/subscriptions/%s/resourceGroups/%s
			*/
			debugLog(m, "did not find provider in API path")
			continue
		}

		apiProvider := matches[1]
		if _, ok := apis[apiProvider]; !ok {
			apis[apiProvider] = map[string]struct{}{}
		}

		matches = apiVersionRegex.FindStringSubmatch(m.Pkg.ID)
		if len(matches) != 2 {
			debugLog(m, "did not find API version in package ID")
			continue
		}

		apiVersion := matches[1]
		apis[apiProvider][apiVersion] = struct{}{}
	}

	for apiProvider, apiVersions := range apis {
		versions := util.MapKeys2Slice(apiVersions)
		sort.Sort(sort.Reverse(sort.StringSlice(versions)))

		result = append(result, API{
			Name:     apiProvider,
			Versions: versions,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}
