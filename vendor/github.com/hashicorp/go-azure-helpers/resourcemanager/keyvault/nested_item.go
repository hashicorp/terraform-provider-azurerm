package keyvault

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"regexp"
)

// TODO: add this to the go-azure-helpers PR

func ValidateNestedItemName(v interface{}, k string) (warnings []string, errors []error) {
	if !regexp.MustCompile(`^[0-9a-zA-Z-]+$`).MatchString(v.(string)) {
		errors = append(errors, fmt.Errorf("`%s` may only contain alphanumeric characters and dashes", k))
	}

	return warnings, errors
}

//

type NestedItemID struct {
	KeyVaultBaseURL string
	NestedItemType  NestedItemType
	Name            string
	Version         *string
}

func NewNestedItemID(keyVaultBaseURL string, nestedItemType NestedItemType, name string, version *string) (*NestedItemID, error) {
	if keyVaultBaseURL == "" {
		return nil, errors.New("expected a non-empty value for `keyVaultBaseURL`")
	}

	if name == "" {
		return nil, errors.New("expected a non-empty value for `name`")
	}

	keyVaultUrl, err := url.Parse(keyVaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing `%s`: %w", keyVaultBaseURL, err)
	}

	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	if nestedItemType == NestedItemTypeAny {
		return nil, fmt.Errorf("`NestedItemTypeAny` is not valid when creating a new nested item ID, please specify one of %s", strings.Join(PossibleNestedItemTypeValues(), ", "))
	}

	return &NestedItemID{
		KeyVaultBaseURL: strings.TrimSuffix(keyVaultUrl.String(), "/"),
		NestedItemType:  nestedItemType,
		Name:            name,
		Version:         version,
	}, nil
}

func (id NestedItemID) ID() string {
	segments := []string{
		id.KeyVaultBaseURL,
		string(id.NestedItemType),
		id.Name,
	}

	if id.Version != nil {
		segments = append(segments, *id.Version)
	}

	return strings.Join(segments, "/")
}

func (id NestedItemID) VersionlessID() string {
	segments := []string{
		id.KeyVaultBaseURL,
		string(id.NestedItemType),
		id.Name,
	}

	return strings.Join(segments, "/")
}

func (id NestedItemID) String() string {
	components := []string{
		fmt.Sprintf("Base URL %q", id.KeyVaultBaseURL),
		fmt.Sprintf("Nested Item Type %q", string(id.NestedItemType)),
		fmt.Sprintf("Name %q", id.Name),
	}

	if id.Version != nil {
		components = append(components, fmt.Sprintf("Version %q", *id.Version))
	}

	return fmt.Sprintf("Key Vault Nested Item (%s)", strings.Join(components, " / "))
}

func ParseNestedItemID(input string, versionType VersionType, nestedItemType NestedItemType) (*NestedItemID, error) {
	id, err := parseNestedItemID(input)
	if err != nil {
		return nil, err
	}

	if versionType == VersionTypeVersioned && id.Version == nil {
		return nil, fmt.Errorf("parsing `%s`: expected a versioned ID", input)
	}

	if versionType == VersionTypeVersionless && id.Version != nil {
		return nil, fmt.Errorf("parsing `%s`: expected a versionless ID", input)
	}

	if nestedItemType != id.NestedItemType && nestedItemType != NestedItemTypeAny {
		return nil, fmt.Errorf("parsing `%s`: expected `NestedItemType` to be `%s`, got `%s`", input, nestedItemType, id.NestedItemType)
	}

	return id, nil
}

func parseNestedItemID(input string) (*NestedItemID, error) {
	inputUrl, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	path := strings.TrimSuffix(strings.TrimPrefix(inputUrl.Path, "/"), "/")
	pathSegments := strings.Split(path, "/")

	if l := len(pathSegments); l != 2 && l != 3 {
		return nil, fmt.Errorf("expected 2 or 3 path segments, found %d segment(s) in `%s`", l, input)
	}

	var version *string
	if len(pathSegments) == 3 {
		version = &pathSegments[2]
	}

	return &NestedItemID{
		KeyVaultBaseURL: fmt.Sprintf("%s://%s", inputUrl.Scheme, inputUrl.Host),
		NestedItemType:  NestedItemType(pathSegments[0]),
		Name:            pathSegments[1],
		Version:         version,
	}, nil
}

// ValidateNestedItemID validates the provided ID based on the provided `VersionType` and `NestedItemType` constants
func ValidateNestedItemID(versionType VersionType, nestedItemType NestedItemType) func(input any, key string) (warnings []string, errors []error) {
	return func(input any, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected `%s` to be a string", key))
		}

		if _, err := ParseNestedItemID(v, versionType, nestedItemType); err != nil {
			errors = append(errors, err)
		}

		return
	}
}

// IsManagedHSM is a helper to determine whether the key vault URL is for a Managed HSM vault.
// This can be used to determine whether to set `managed_hsm_key_id` into state while this argument is in a deprecated state.
func (id NestedItemID) IsManagedHSM() bool {
	return strings.Contains(id.KeyVaultBaseURL, ".managedhsm.")
}
