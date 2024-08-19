// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recaser

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// ReCase tries to determine the type of Resource ID defined in `input` to be able to re-case it from
func ReCase(input string) string {
	return reCaseWithIds(input, knownResourceIds)
}

// reCaseWithIds tries to determine the type of Resource ID defined in `input` to be able to re-case it based on an input list of Resource IDs
// this is a "best-effort" function and can return the input unmodified. Functionality of this method is intended to be
// limited to resource IDs that have been registered with the package via the RegisterResourceId() function at init.
// However, some common static segments are corrected even when a corresponding ID type is not present.
func reCaseWithIds(input string, ids map[string]resourceids.ResourceId) string {
	result, err := reCaseKnownId(input, ids)
	if err == nil {
		return pointer.From(result)
	}

	output := input

	// if we didn't find a matching id then re-case these known segments for best effort
	segmentsToFix := []string{
		"/subscriptions/",
		"/resourceGroups/",
		"/managementGroups/",
		"/tenants/",
	}

	for _, segment := range segmentsToFix {
		output = fixSegment(output, segment)
	}

	return output
}

// ReCaseKnownId attempts to correct the casing on the static segments of an Azure resourceId. Functionality of this
// method is intended to be limited to resource IDs that have been registered with the package via the
// RegisterResourceId() function at init.
func ReCaseKnownId(input string) (*string, error) {
	return reCaseKnownId(input, knownResourceIds)
}

func reCaseKnownId(input string, ids map[string]resourceids.ResourceId) (*string, error) {
	output := input
	parsed := false
	key, ok := buildInputKey(input)
	if ok {
		id := ids[*key]
		if id != nil {
			var parseError error
			output, parseError = parseId(id, input)
			if parseError != nil {
				return &output, fmt.Errorf("fixing case for ID '%s': %+v", input, parseError)
			}
			parsed = true
		} else {
			for _, v := range PotentialScopeValues() {
				trimmedKey := strings.TrimPrefix(*key, v)
				if id = knownResourceIds[trimmedKey]; id != nil {
					var parseError error
					output, parseError = parseId(id, input)
					if parseError != nil {
						return &output, fmt.Errorf("fixing case for ID '%s': %+v", input, parseError)
					} else {
						parsed = true
						break
					}
				}
				// We have some cases where an erroneous trailing '/' causes problems. These may be data errors in specs, or API responses.
				// Either way, we can try and compensate for it.
				if id = knownResourceIds[strings.TrimPrefix(*key, strings.TrimSuffix(v, "/"))]; id != nil {
					var parseError error
					output, parseError = parseId(id, input)
					if parseError != nil {
						return &output, fmt.Errorf("fixing case for ID '%s': %+v", input, parseError)
					} else {
						parsed = true
						break
					}
				}
			}
		}
	}

	if !parsed {
		return &output, fmt.Errorf("could not determine ID type for '%s', or ID type not supported", input)
	}
	return &output, nil
}

// parseId uses the specified ResourceId to parse the input and returns the id string with correct casing
func parseId(id resourceids.ResourceId, input string) (string, error) {

	// we need to take a local copy of id to work against else we're mutating the original
	localId := id

	parser := resourceids.NewParserFromResourceIdType(localId)
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return input, err
	}

	if scope := parsed.Parsed["scope"]; scope != "" {
		parsed.Parsed["scope"] = reCaseWithIds(scope, knownResourceIds)
	}

	if err = id.FromParseResult(*parsed); err != nil {
		return input, err
	}
	input = id.ID()

	return input, err
}

// fixSegment searches the input id string for a specified segment case-insensitively
// and returns the input string with the casing corrected on the segment
func fixSegment(input, segment string) string {
	if strings.Contains(strings.ToLower(input), strings.ToLower(segment)) {
		re := regexp.MustCompile(fmt.Sprintf("(?i)%s", segment))
		input = re.ReplaceAllString(input, segment)
	}
	return input
}

// buildInputKey takes an input id string and removes user-specified values from it
// so it can be used as a key to extract the correct id from knownResourceIds
func buildInputKey(input string) (*string, bool) {

	// Attempt to determine if this is just missing a leading slash and prepend it if it seems to be
	if !strings.HasPrefix(input, "/") {
		if len(input) == 0 || !strings.Contains(input, "/") {
			return nil, false
		}

		input = "/" + input
	}

	output := ""

	segments := strings.Split(input, "/")
	// iterate through the segments extracting any that are not user inputs
	// and append them together to make a key
	// eg "/subscriptions/1111/resourceGroups/group1/providers/Microsoft.BotService/botServices/botServiceValue" will become:
	// "/subscriptions//resourceGroups//providers/Microsoft.BotService/botServices/"
	if len(segments)%2 != 0 {
		for i := 1; len(segments) > i; i++ {
			if i%2 != 0 {
				key := segments[i]
				output = fmt.Sprintf("%s/%s/", output, key)

				// if the current segment is a providers segment, then we should append the next segment to the key
				// as this is not a user input segment
				if strings.EqualFold(key, "providers") && len(segments) >= i+2 {
					value := segments[i+1]
					output = fmt.Sprintf("%s%s", output, value)
				}
			}
		}
	}
	output = strings.ToLower(output)
	return &output, true
}

// PotentialScopeValues returns a list of possible ScopeSegment values from all registered ID types
// This is a best effort process, limited to scope targets that are prefixed with '/subscriptions/' or '/providers/'
func PotentialScopeValues() []string {
	result := make([]string, 0)
	for k := range knownResourceIds {
		if strings.HasPrefix(k, "/subscriptions/") || strings.HasPrefix(k, "/providers/") {
			result = append(result, k)
		}
	}

	return result
}

// ResourceIdTypeFromResourceId takes a Azure Resource ID as a string and attempts to return the corresponding
// resourceids.ResourceId type. If a matching resourceId is not found in the supported/registered resourceId types then
// a `nil` value is returned.
func ResourceIdTypeFromResourceId(input string) resourceids.ResourceId {
	key, ok := buildInputKey(input)
	if ok {
		id := knownResourceIds[*key]
		if id != nil {
			result := reflect.New(reflect.TypeOf(id).Elem())
			return result.Interface().(resourceids.ResourceId)
		} else {
			for _, v := range PotentialScopeValues() {
				trimmedKey := strings.TrimPrefix(*key, v)
				if id = knownResourceIds[trimmedKey]; id != nil {
					result := reflect.New(reflect.TypeOf(id).Elem())
					return result.Interface().(resourceids.ResourceId)
				}
				if id = knownResourceIds[strings.TrimPrefix(*key, strings.TrimSuffix(v, "/"))]; id != nil {
					result := reflect.New(reflect.TypeOf(id).Elem())
					return result.Interface().(resourceids.ResourceId)
				}
			}
		}
	}

	return nil
}
