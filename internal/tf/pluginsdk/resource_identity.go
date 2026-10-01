// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// These functions support generating the resource identity schema for the following types of identities and resources
// * Hierarchical IDs (untyped and typed resources)
// * Composite Resource IDs (untyped and typed resources)

// ResourceTypeForIdentity is used to select different schema generation behaviours depending on the type of resource/resource ID
type ResourceTypeForIdentity int

const (
	ResourceTypeForIdentityDefault = iota
	ResourceTypeForIdentityVirtual
)

func identityType(idType []ResourceTypeForIdentity) ResourceTypeForIdentity {
	var t ResourceTypeForIdentity

	switch len(idType) {
	case 0:
		t = ResourceTypeForIdentityDefault
	case 1:
		t = idType[0]
	default:
		panic(fmt.Sprintf("expected a maximum of one value for the `idType` argument, got %d", len(idType)))
	}

	return t
}

// SegmentTypeSupported contains a list of segments that should be used to construct the resource identity schema
// this list will need to be extended to support hierarchical resource IDs for management groups or resources
// that begin with a different prefix to /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1
func SegmentTypeSupported(segment resourceids.SegmentType) bool {
	supportedSegmentTypes := []resourceids.SegmentType{
		resourceids.ConstantSegmentType,
		resourceids.SubscriptionIdSegmentType,
		resourceids.ResourceGroupSegmentType,
		resourceids.UserSpecifiedSegmentType,
		resourceids.ResourceIDSegmentType,
		resourceids.ScopeSegmentType,
	}

	return slices.Contains(supportedSegmentTypes, segment)
}

func SegmentName(segment resourceids.Segment, idType ResourceTypeForIdentity, numSegments, idx int) (name string) {
	switch idType {
	case ResourceTypeForIdentityVirtual:
		name = ToSnakeCase(segment.Name)
	default:
		// For the last segment, if it's a `*Name` field, we generate it as `name` rather than snake casing the segment's name
		if (idx+1) == numSegments && strings.HasSuffix(segment.Name, "Name") {
			return "name"
		}
		name = ToSnakeCase(segment.Name)
	}

	return normaliseSegmentName(name)
}

// GenerateIdentitySchema generates the resource identity schema from the resource ID type
func GenerateIdentitySchema(id resourceids.ResourceId, idType ...ResourceTypeForIdentity) func() map[string]*schema.Schema {
	return func() map[string]*schema.Schema {
		return identitySchema(id, identityType(idType))
	}
}

// GenerateCompositeIdentitySchema generates the resource identity schema for a composite resource ID with optional custom segment names
func GenerateCompositeIdentitySchema(id resourceids.ResourceId, customNames ...string) func() map[string]*schema.Schema {
	return func() map[string]*schema.Schema {
		return identitySchema(id, ResourceTypeForIdentityDefault, customNames...)
	}
}

func identitySchema(id resourceids.ResourceId, idType ResourceTypeForIdentity, customNames ...string) map[string]*schema.Schema {
	idSchema := make(map[string]*schema.Schema)

	segments := compositeSegments(id)
	numSegments := len(segments)
	for idx, segment := range segments {
		if SegmentTypeSupported(segment.Type) {
			name := SegmentName(segment, idType, numSegments, idx)
			if idx < len(customNames) && customNames[idx] != "" {
				name = customNames[idx]
			}

			idSchema[name] = &schema.Schema{
				Type:              schema.TypeString,
				RequiredForImport: true,
			}
		}
	}
	return idSchema
}

// ValidateResourceId validates that the provided inputId string is valid according to the resource ID definition.
// For composite resource IDs, it splits the ID on '|' and validates each component with its respective parser.
// For standard resource IDs, it uses resourceids.NewParserFromResourceIdType.
func ValidateResourceId(id resourceids.ResourceId, inputId string, insensitively bool) error {
	if isCompositeResourceId(id) {
		components := strings.Split(inputId, "|")
		segments := compositeSegments(id)
		if len(components) != len(segments) {
			return fmt.Errorf("parsing %q: expected %d resource IDs separated by '|', got %d", inputId, len(segments), len(components))
		}

		comps, hasComps := getCompositeResourceComponents(id)
		for i, compStr := range components {
			if strings.TrimSpace(compStr) == "" {
				return fmt.Errorf("parsing %q: component %d cannot be empty", inputId, i+1)
			}

			if hasComps && i < len(comps) {
				parser := resourceids.NewParserFromResourceIdType(comps[i])
				if _, err := parser.Parse(compStr, insensitively); err != nil {
					return fmt.Errorf("parsing %q component %d (%q): %w", inputId, i+1, compStr, err)
				}
			}
		}
		return nil
	}

	parser := resourceids.NewParserFromResourceIdType(id)
	if _, err := parser.Parse(inputId, insensitively); err != nil {
		return err
	}
	return nil
}

// ValidateResourceIdentityData validates the resource identity data provided by the user when performing a plannable
// import using resource identity
func ValidateResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId, idType ...ResourceTypeForIdentity) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	if isCompositeResourceId(id) {
		return validateCompositeResourceIdentityData(d, identity, id, identityType(idType))
	}

	var components []string
	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		switch segment.Type {
		case resourceids.StaticSegmentType, resourceids.ResourceProviderSegmentType:
			components = append(components, pointer.From(segment.FixedValue))

		case resourceids.ScopeSegmentType:
			name := SegmentName(segment, identityType(idType), numSegments, idx)

			field, ok := identity.GetOk(name)
			if !ok {
				return fmt.Errorf("getting %q in resource identity", name)
			}

			value, ok := field.(string)
			if !ok {
				return fmt.Errorf("converting %q to string", name)
			}

			if value == "" {
				return fmt.Errorf("%q cannot be empty", name)
			}

			if err := identity.Set(name, value); err != nil {
				return fmt.Errorf("error setting id: %+v", err)
			}

			trimmed := strings.Trim(value, "/")
			if trimmed == "" && value == "/" {
				components = append(components, "")
			} else {
				components = append(components, trimmed)
			}

		default:
			if SegmentTypeSupported(segment.Type) {
				name := SegmentName(segment, identityType(idType), numSegments, idx)

				field, ok := identity.GetOk(name)
				if !ok {
					return fmt.Errorf("getting %q in resource identity", name)
				}

				value, ok := field.(string)
				if !ok {
					return fmt.Errorf("converting %q to string", name)
				}

				if value == "" {
					return fmt.Errorf("%q cannot be empty", name)
				}

				if err := identity.Set(name, value); err != nil {
					return fmt.Errorf("error setting id: %+v", err)
				}

				components = append(components, strings.Trim(value, "/"))
			}
		}
	}

	identityString := "/" + strings.Join(components, "/")

	if err := ValidateResourceId(id, identityString, true); err != nil {
		return fmt.Errorf("parsing after building Resource ID: %s", err)
	}

	d.SetId(identityString)

	return nil
}

// ValidateCompositeResourceIdentityData validates composite resource identity data with optional custom segment names
func ValidateCompositeResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId, customNames ...string) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	return validateCompositeResourceIdentityData(d, identity, id, ResourceTypeForIdentityDefault, customNames...)
}

func validateCompositeResourceIdentityData(d *schema.ResourceData, identity *schema.IdentityData, id resourceids.ResourceId, idType ResourceTypeForIdentity, customNames ...string) error {
	segments := compositeSegments(id)
	numSegments := len(segments)
	values := make([]string, 0, numSegments)

	for idx, segment := range segments {
		if SegmentTypeSupported(segment.Type) {
			name := SegmentName(segment, idType, numSegments, idx)
			if idx < len(customNames) && customNames[idx] != "" {
				name = customNames[idx]
			}

			field, ok := identity.GetOk(name)
			if !ok {
				return fmt.Errorf("getting %q in resource identity", name)
			}

			value, ok := field.(string)
			if !ok {
				return fmt.Errorf("converting %q to string", name)
			}

			if value == "" {
				return fmt.Errorf("%q cannot be empty", name)
			}

			if err := identity.Set(name, value); err != nil {
				return fmt.Errorf("error setting id: %+v", err)
			}

			values = append(values, value)
		}
	}

	identityString := strings.Join(values, "|")

	if err := ValidateResourceId(id, identityString, true); err != nil {
		return fmt.Errorf("parsing after building Resource ID: %s", err)
	}

	d.SetId(identityString)
	return nil
}

// SetResourceIdentityData sets the resource identity data in state
func SetResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId, idType ...ResourceTypeForIdentity) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	if isCompositeResourceId(id) {
		return compositeResourceIdentityData(identity, id, identityType(idType))
	}

	if err := resourceIdentityData(identity, id, identityType(idType)); err != nil {
		return err
	}

	return nil
}

// SetCompositeResourceIdentityData sets the composite resource identity data in state with optional custom segment names
func SetCompositeResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId, customNames ...string) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	return compositeResourceIdentityData(identity, id, ResourceTypeForIdentityDefault, customNames...)
}

func compositeResourceIdentityData(identity *schema.IdentityData, id resourceids.ResourceId, idType ResourceTypeForIdentity, customNames ...string) error {
	components := strings.Split(id.ID(), "|")
	segments := compositeSegments(id)
	numSegments := len(segments)

	if len(components) != numSegments {
		return fmt.Errorf("parsing resource ID %q: expected %d parts separated by '|', got %d", id.ID(), numSegments, len(components))
	}

	comps, hasComps := getCompositeResourceComponents(id)
	for idx, segment := range segments {
		if SegmentTypeSupported(segment.Type) {
			name := SegmentName(segment, idType, numSegments, idx)
			if idx < len(customNames) && customNames[idx] != "" {
				name = customNames[idx]
			}

			val := components[idx]
			if hasComps && idx < len(comps) {
				parser := resourceids.NewParserFromResourceIdType(comps[idx])
				if _, err := parser.Parse(val, true); err != nil {
					return fmt.Errorf("parsing component %d of resource ID %s: %w", idx+1, id.ID(), err)
				}
			}

			if err := identity.Set(name, val); err != nil {
				return fmt.Errorf("setting `%s` in resource identity: %+v", name, err)
			}
		}
	}

	return nil
}

func resourceIdentityData(identity *schema.IdentityData, id resourceids.ResourceId, idType ResourceTypeForIdentity) error {
	parser := resourceids.NewParserFromResourceIdType(id)
	parsed, err := parser.Parse(id.ID(), true)
	if err != nil {
		return fmt.Errorf("parsing resource ID: %s", err)
	}

	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if SegmentTypeSupported(segment.Type) {
			name := SegmentName(segment, idType, numSegments, idx)

			field, ok := parsed.Parsed[segment.Name]
			if !ok {
				return fmt.Errorf("field `%s` was not found in the parsed resource ID %s", name, id)
			}

			if err = identity.Set(name, field); err != nil {
				return fmt.Errorf("setting `%s` in resource identity: %+v", name, err)
			}
		}
	}

	return nil
}

func isCompositeResourceId(id resourceids.ResourceId) bool {
	if id == nil {
		return false
	}
	if _, ok := getCompositeResourceComponents(id); ok {
		return true
	}
	for _, seg := range id.Segments() {
		if seg.Type == resourceids.ResourceIDSegmentType {
			return true
		}
	}
	return false
}

func getCompositeResourceComponents(id resourceids.ResourceId) ([]resourceids.ResourceId, bool) {
	if id == nil {
		return nil, false
	}
	val := reflect.ValueOf(id)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, false
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, false
	}

	firstField := val.FieldByName("First")
	secondField := val.FieldByName("Second")
	if !firstField.IsValid() || !secondField.IsValid() {
		return nil, false
	}

	getFirst := func(f reflect.Value) (resourceids.ResourceId, bool) {
		if !f.IsValid() {
			return nil, false
		}
		if f.Kind() == reflect.Pointer && f.IsNil() {
			elemType := f.Type().Elem()
			newVal := reflect.New(elemType)
			if rId, ok := reflect.TypeAssert[resourceids.ResourceId](newVal); ok {
				return rId, true
			}
		}
		if rId, ok := reflect.TypeAssert[resourceids.ResourceId](f); ok {
			return rId, true
		}
		return nil, false
	}

	first, ok1 := getFirst(firstField)
	second, ok2 := getFirst(secondField)
	if ok1 && ok2 {
		return []resourceids.ResourceId{first, second}, true
	}

	return nil, false
}

func compositeSegments(id resourceids.ResourceId) []resourceids.Segment {
	if comps, ok := getCompositeResourceComponents(id); ok && len(comps) == 2 {
		firstSegments := comps[0].Segments()
		secondSegments := comps[1].Segments()
		return []resourceids.Segment{
			resourceids.ResourceIDSegment("resourceID1", resourceids.BuildExpectedResourceId(firstSegments)),
			resourceids.ResourceIDSegment("resourceID2", resourceids.BuildExpectedResourceId(secondSegments)),
		}
	}
	return id.Segments()
}

// normaliseSegmentName replaces any known oddities when generating based on struct field names with predetermined values.
func normaliseSegmentName(input string) string {
	// For now, use complete values, if this becomes unmanageable investigate a more general replacement strategy
	replacements := map[string]string{
		"signal_r_name":    "signalr_name",
		"web_pub_sub_name": "web_pubsub_name",
		"resource_id_1":    "resource_id1",
		"resource_id_2":    "resource_id2",
	}

	if v, ok := replacements[input]; ok {
		return v
	}

	return input
}

// ToSnakeCase is a slightly altered version of `strcase.ToSnake()`
// the main difference is that it doesn't split patterns like `v2` to `v_2`
func ToSnakeCase(input string) string {
	delimiter := uint8('_')

	input = strings.TrimSpace(input)
	n := strings.Builder{}
	n.Grow(len(input) + 2) // nominal 2 bytes of extra space for inserted delimiters
	for i, v := range []byte(input) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if vIsCap {
			v += 'a'
			v -= 'A'
		}

		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		if i+1 < len(input) {
			next := input[i+1]
			vIsNum := v >= '0' && v <= '9'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			nextIsNum := next >= '0' && next <= '9'

			// if it looks like a version (e.g. `ServerGroupsv2Name`), group `v2`/`V2` together -> `_v2_`
			if (vIsLow || vIsCap) && nextIsNum && (v == 'v' || v == 'V') {
				// avoid duplicate delimiter
				if b := []byte(n.String()); len(b)-1 >= 0 && b[len(b)-1] != delimiter {
					n.WriteByte(delimiter)
				}
				n.WriteByte(v)
				continue
			}

			// add underscore if next letter case type is changed
			if (vIsCap && (nextIsLow || nextIsNum)) || (vIsLow && (nextIsCap || nextIsNum)) || (vIsNum && (nextIsCap || nextIsLow)) {
				if vIsCap && nextIsLow {
					if prevIsCap := i > 0 && input[i-1] >= 'A' && input[i-1] <= 'Z'; prevIsCap {
						n.WriteByte(delimiter)
					}
				}
				n.WriteByte(v)
				if vIsLow || vIsNum || nextIsNum {
					n.WriteByte(delimiter)
				}
				continue
			}
		}

		if v == ' ' || v == '_' || v == '-' || v == '.' {
			// replace space/underscore/hyphen/dot with delimiter
			n.WriteByte(delimiter)
		} else {
			n.WriteByte(v)
		}
	}

	return n.String()
}
