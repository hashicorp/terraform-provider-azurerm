package version

import (
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/api-version-lint/sdk"
	"gopkg.in/yaml.v3"
)

type versionException struct {
	Module                  string  `yaml:"module"`
	Service                 string  `yaml:"service"`
	Version                 string  `yaml:"version"`
	StableVersionTargetDate *string `yaml:"stableVersionTargetDate,omitempty"`
	ResponsibleIndividual   *string `yaml:"responsibleIndividual,omitempty"`
}

func ParseHistoricalExceptions(path string) ([]Version, error) {
	yaml, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return parseExceptions(yaml, true)
}

func ParseExceptions(path string) ([]Version, error) {
	yaml, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return parseExceptions(yaml, false)
}

func parseExceptions(yamlBytes []byte, isHistorical bool) ([]Version, error) {
	exceptions := []versionException{}
	if err := yaml.Unmarshal(yamlBytes, &exceptions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	if !sort.SliceIsSorted(exceptions, func(i int, j int) bool {
		if exceptions[i].Module != exceptions[j].Module {
			return exceptions[i].Module < exceptions[j].Module
		}
		if exceptions[i].Service != exceptions[j].Service {
			return exceptions[i].Service < exceptions[j].Service
		}
		return exceptions[i].Version < exceptions[j].Version
	}) {
		return nil, fmt.Errorf("entries has to be sorted alphabetically by module, service and version")
	}

	validModules := make(map[string]bool, len(sdk.SdkTypes))
	for _, sdkType := range sdk.SdkTypes {
		validModules[sdkType.Module] = true
	}

	versions := make([]Version, 0, len(exceptions))
	errors := []string{}

	for _, e := range exceptions {
		e := trimSpaces(e)

		if errs := validateVersionException(e, isHistorical, validModules); len(errs) > 0 {
			errors = append(errors, errs...)
		} else {
			versions = append(versions, Version{
				Module:  e.Module,
				Service: e.Service,
				Version: e.Version,
			})
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("- %s", strings.Join(errors, "\n- "))
	}

	return versions, nil
}

func validateVersionException(e versionException, isHistorical bool, validModules map[string]bool) (errors []string) {
	if e.Module == "" {
		errors = append(errors, "module is required")
	} else if !validModules[e.Module] {
		validModulesStr := make([]string, 0, len(validModules))
		for k := range validModules {
			validModulesStr = append(validModulesStr, k)
		}
		errors = append(errors, fmt.Sprintf("unsupported sdk module %q\nvalid modules are: %q", e.Module, strings.Join(validModulesStr, `", "`)))
	} else if e.Service == "" {
		errors = append(errors, fmt.Sprintf("module %q: service is required", e.Module))
	} else if e.Version == "" {
		errors = append(errors, fmt.Sprintf("module %q, service %q: version is required", e.Module, e.Service))
	}

	if len(errors) > 0 {
		return
	}

	if !isHistorical {
		if pointer.From(e.StableVersionTargetDate) == "" {
			errors = append(errors, fmt.Sprintf("module %q, service %q, version %q: stableVersionTargetDate is required", e.Module, e.Service, e.Version))
		} else if !validRFC3339DateOnly(pointer.From(e.StableVersionTargetDate)) {
			errors = append(errors, fmt.Sprintf("module %q, service %q, version %q: invalid stableVersionTargetDate %q, has to be in YYYY-MM-DD format", e.Module, e.Service, e.Version, pointer.From(e.StableVersionTargetDate)))
		}

		if pointer.From(e.ResponsibleIndividual) == "" {
			errors = append(errors, fmt.Sprintf("module %q, service %q, version %q: responsibleIndividual is required", e.Module, e.Service, e.Version))
		} else if !validResponsibleIndividual(pointer.From(e.ResponsibleIndividual)) {
			errors = append(errors, fmt.Sprintf("module %q, service %q, version %q: invalid responsibleIndividual %q, has to be an email or a `github.com/yourname` GitHub handle", e.Module, e.Service, e.Version, pointer.From(e.ResponsibleIndividual)))
		}
	}

	return
}

func validRFC3339DateOnly(s string) bool {
	_, err := time.Parse(time.DateOnly, s)
	return err == nil
}

var githubHandleRegex = regexp.MustCompile(`(?i)^github.com/[a-z0-9-_]+$`)

func validResponsibleIndividual(s string) bool {
	_, err := mail.ParseAddress(s)
	if err == nil {
		return true
	}
	return githubHandleRegex.MatchString(s)
}

func trimSpaces(e versionException) versionException {
	res := versionException{
		Module:  strings.TrimSpace(e.Module),
		Service: strings.TrimSpace(e.Service),
		Version: strings.TrimSpace(e.Version),
	}
	if e.StableVersionTargetDate != nil {
		res.StableVersionTargetDate = pointer.To(strings.TrimSpace(*e.StableVersionTargetDate))
	}
	if e.ResponsibleIndividual != nil {
		res.ResponsibleIndividual = pointer.To(strings.TrimSpace(*e.ResponsibleIndividual))
	}
	return res
}
