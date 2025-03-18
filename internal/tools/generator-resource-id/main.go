// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var packagesUsingAlias = map[string]struct{}{
	"advisor":          {},
	"analysisservices": {},
	"appconfiguration": {},
}

func main() {
	servicePackagePath := flag.String("path", "", "The relative path to the service package")
	name := flag.String("name", "", "The name of this Resource Type")
	id := flag.String("id", "", "An example of this Resource ID")
	rewrite := flag.Bool("rewrite", false, "Should this Resource ID be parsed insensitively, to workaround an API bug?")
	showHelp := flag.Bool("help", false, "Display this message")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	if err := run(*servicePackagePath, *name, *id, *rewrite); err != nil {
		panic(err)
	}
}

func run(servicePackagePath, name, id string, shouldRewrite bool) error {
	servicePackage, err := parseServicePackageName(servicePackagePath)
	if err != nil {
		return fmt.Errorf("determining Service Package Name for %q: %+v", servicePackagePath, err)
	}

	parsersPath := path.Join(servicePackagePath, "/parse")
	if err := os.Mkdir(parsersPath, 0o755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("creating parse directory at %q: %+v", parsersPath, err)
	}

	validatorPath := path.Join(servicePackagePath, "/validate")
	if err := os.Mkdir(validatorPath, 0o755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("creating validate directory at %q: %+v", validatorPath, err)
	}

	fileName := convertToSnakeCase(name)
	validatorFileName := fmt.Sprintf("%s_id", fileName)
	if strings.HasSuffix(fileName, "_test") {
		// e.g. "webtest" in applicationInsights
		fileName += "_id"
	}
	resourceId, err := NewResourceID(name, *servicePackage, id)
	if err != nil {
		return err
	}

	generator := ResourceIdGenerator{
		ResourceId:    *resourceId,
		ShouldRewrite: shouldRewrite,
	}

	parserFilePath := fmt.Sprintf("%s/%s.go", parsersPath, fileName)
	if err := goFmtAndWriteToFile(parserFilePath, generator.Code()); err != nil {
		return fmt.Errorf("generating Parser at %q: %+v", parserFilePath, err)
	}

	parserTestsFilePath := fmt.Sprintf("%s/%s_test.go", parsersPath, fileName)
	if err := goFmtAndWriteToFile(parserTestsFilePath, generator.TestCode()); err != nil {
		return fmt.Errorf("generating Parser Tests at %q: %+v", parserTestsFilePath, err)
	}

	validatorFilePath := fmt.Sprintf("%s/%s.go", validatorPath, validatorFileName)
	if err := goFmtAndWriteToFile(validatorFilePath, generator.ValidatorCode()); err != nil {
		return fmt.Errorf("generating Validator at %q: %+v", validatorFilePath, err)
	}

	validatorTestsFilePath := fmt.Sprintf("%s/%s_test.go", validatorPath, validatorFileName)
	if err := goFmtAndWriteToFile(validatorTestsFilePath, generator.ValidatorTestCode()); err != nil {
		return fmt.Errorf("generating Validator Tests at %q: %+v", validatorTestsFilePath, err)
	}

	return nil
}

func parseServicePackageName(relativePath string) (*string, error) {
	path := relativePath
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		path = abs
	}

	// we do this replacement to avoid the case that on windows machine, the absolute path are using the path separator of \ instead of /
	path = strings.ReplaceAll(path, "\\", "/")
	segments := strings.Split(path, "/")
	serviceIndex := -1
	for i, v := range segments {
		if strings.EqualFold(v, "services") {
			serviceIndex = i
			break
		}
	}

	if serviceIndex == -1 {
		return nil, fmt.Errorf("`services` segment was not found")
	}

	if len(segments) <= serviceIndex {
		return nil, fmt.Errorf("not enough segments")
	}

	servicePackageName := segments[serviceIndex+1]
	return &servicePackageName, nil
}

func convertToSnakeCase(input string) string {
	splitIdxMap := map[int]struct{}{}
	var lastChar rune
	for idx, char := range input {
		switch {
		case idx == 0:
			splitIdxMap[idx] = struct{}{}
		case unicode.IsUpper(lastChar) == unicode.IsUpper(char):
		case unicode.IsUpper(lastChar):
			splitIdxMap[idx-1] = struct{}{}
		case unicode.IsUpper(char):
			splitIdxMap[idx] = struct{}{}
		}
		lastChar = char
	}
	splitIdx := make([]int, 0, len(splitIdxMap))
	for idx := range splitIdxMap {
		splitIdx = append(splitIdx, idx)
	}
	sort.Ints(splitIdx)

	inputRunes := []rune(input)
	out := make([]string, len(splitIdx))
	for i := range splitIdx {
		if i == len(splitIdx)-1 {
			out[i] = strings.ToLower(string(inputRunes[splitIdx[i]:]))
			continue
		}
		out[i] = strings.ToLower(string(inputRunes[splitIdx[i]:splitIdx[i+1]]))
	}
	return strings.Join(out, "_")
}

type ResourceIdSegment struct {
	// ArgumentName is the name which should be used when this segment is used in an Argument
	ArgumentName string

	// FieldName is the name which should be used for this segment when referenced in a Field
	FieldName string

	// SegmentKey is the Segment used for this in the Resource ID e.g. `resourceGroups`
	SegmentKey string

	// SegmentValue is the value for this segment used in the Resource ID
	SegmentValue string
}

type ResourceId struct {
	TypeName string
	IDFmt    string
	IDRaw    string

	ServicePackageName string
	TestPackageSuffix  string

	HasResourceGroup  bool
	HasSubscriptionId bool
	Segments          []ResourceIdSegment // this has to be a slice not a map since we care about the order
}

func NewResourceID(typeName, servicePackageName, resourceId string) (*ResourceId, error) {
	// split the string, but remove the prefix of `/` since it's an empty segment
	split := strings.Split(strings.TrimPrefix(resourceId, "/"), "/")
	if len(split)%2 != 0 {
		return nil, fmt.Errorf("segments weren't divisible by 2: %q", resourceId)
	}

	segments := make([]ResourceIdSegment, 0)
	for i := 0; i < len(split); i += 2 {
		key := split[i]
		value := split[i+1]

		// the RP shouldn't be transformed
		if key == "providers" {
			r := regexp.MustCompile(`^Microsoft.[A-Z][A-Za-z]+$`)
			if !r.MatchString(value) {
				return nil, fmt.Errorf("the resource provider in the id must begin with upper case got: %s", value)
			}
			continue
		}

		segmentBuilder := func(key, value string, hasSubscriptionId bool) ResourceIdSegment {
			toCamelCase := func(input string) string {
				// lazy but it works
				out := make([]rune, 0)
				for i, char := range azure.TitleCase(input) {
					if i == 0 {
						out = append(out, unicode.ToLower(char))
						continue
					}

					out = append(out, char)
				}
				return string(out)
			}

			rewritten := fmt.Sprintf("%sName", key)
			segment := ResourceIdSegment{
				FieldName:    azure.TitleCase(rewritten),
				ArgumentName: toCamelCase(rewritten),
				SegmentKey:   key,
				SegmentValue: value,
			}

			if strings.EqualFold(key, "resourceGroups") {
				segment.FieldName = "ResourceGroup"
				segment.ArgumentName = "resourceGroup"
				return segment
			}

			if key == "subscriptions" && !hasSubscriptionId {
				segment.FieldName = "SubscriptionId"
				segment.ArgumentName = "subscriptionId"
				return segment
			}

			if strings.HasSuffix(key, "s") {
				// TODO: in time this could be worth a series of overrides

				// handles "GallerieName" and `DataFactoriesName`
				if strings.HasSuffix(key, "ies") {
					key = strings.TrimSuffix(key, "ies")
					key = fmt.Sprintf("%sy", key)
				}
				switch {
				case strings.HasSuffix(key, "sses"):
					// handles `PublicIPAddressesName`
					key = strings.TrimSuffix(key, "sses")
					key = fmt.Sprintf("%sss", key)
				case strings.HasSuffix(key, "xes"):
					// handles `CustomIPPrefixeName`
					key = strings.TrimSuffix(key, "xes")
					key = fmt.Sprintf("%sx", key)
				default:
					key = strings.TrimSuffix(key, "s")
				}

				if strings.EqualFold(key, typeName) {
					segment.FieldName = "Name"
					segment.ArgumentName = "name"
				} else {
					// remove {Thing}s and make that {Thing}Name
					rewritten = fmt.Sprintf("%sName", key)
					segment.FieldName = azure.TitleCase(rewritten)
					segment.ArgumentName = toCamelCase(rewritten)
				}
			}

			return segment
		}

		// handle multiple 'subscriptions' segments, ala ServiceBus Subscription
		hasSubscriptionId := false
		for _, v := range segments {
			if v.FieldName == "SubscriptionId" {
				hasSubscriptionId = true
				break
			}
		}

		segment := segmentBuilder(key, value, hasSubscriptionId)
		segments = append(segments, segment)
	}

	// finally build up the format string based on this information
	fmtString := resourceId
	hasResourceGroup := false
	hasSubscriptionId := false
	for _, segment := range segments {
		if strings.EqualFold(segment.SegmentKey, "subscriptions") {
			hasSubscriptionId = true
		}
		if strings.EqualFold(segment.SegmentKey, "resourceGroups") {
			hasResourceGroup = true
		}

		// has to be double-escaped since this is a fmtstring
		fmtString = strings.Replace(fmtString, segment.SegmentValue, "%s", 1)
	}

	packageSuffix := ""
	if _, ok := packagesUsingAlias[servicePackageName]; ok {
		packageSuffix = "_test"
	}

	return &ResourceId{
		IDFmt:              fmtString,
		IDRaw:              resourceId,
		HasResourceGroup:   hasResourceGroup,
		HasSubscriptionId:  hasSubscriptionId,
		Segments:           segments,
		ServicePackageName: servicePackageName,
		TypeName:           typeName,
		TestPackageSuffix:  packageSuffix,
	}, nil
}

type ResourceIdGenerator struct {
	ResourceId

	ShouldRewrite bool
}

func (id ResourceIdGenerator) Code() string {
	return fmt.Sprintf(`
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

%s
%s
%s
%s
%s
%s
`, id.codeForType(), id.codeForConstructor(), id.codeForDescription(), id.codeForFormatter(), id.codeForParser(), id.codeForParserInsensitive())
}

func (id ResourceIdGenerator) codeForType() string {
	fields := make([]string, 0)
	for _, segment := range id.Segments {
		fields = append(fields, fmt.Sprintf("\t%s\tstring", segment.FieldName))
	}
	fieldStr := strings.Join(fields, "\n")
	return fmt.Sprintf(`
type %[1]sId struct {
%[2]s
}
`, id.TypeName, fieldStr)
}

func (id ResourceIdGenerator) codeForConstructor() string {
	arguments := make([]string, 0)
	assignments := make([]string, 0)

	for _, segment := range id.Segments {
		arguments = append(arguments, segment.ArgumentName)
		assignments = append(assignments, fmt.Sprintf("\t\t%s:\t%s,", segment.FieldName, segment.ArgumentName))
	}

	argumentsStr := strings.Join(arguments, ", ")
	assignmentsStr := strings.Join(assignments, "\n")
	return fmt.Sprintf(`
func New%[1]sID(%[2]s string) %[1]sId {
	return %[1]sId{
%[3]s
	}
}
`, id.TypeName, argumentsStr, assignmentsStr)
}

func (id ResourceIdGenerator) codeForDescription() string {
	makeHumanReadable := func(input string) string {
		chars := make([]rune, 0)
		for i, c := range input {
			if unicode.IsUpper(c) && i+1 < len(input) && unicode.IsLower(rune(input[i+1])) {
				chars = append(chars, ' ')
			}

			chars = append(chars, c)
		}
		out := string(chars)
		return strings.TrimSpace(out)
	}

	formatKeys := make([]string, 0)
	for _, segment := range id.Segments {
		if segment.FieldName == "SubscriptionId" {
			continue
		}

		humanReadableKey := makeHumanReadable(segment.FieldName)
		formatKeys = append(formatKeys, fmt.Sprintf("\t\tfmt.Sprintf(\"%[1]s %%q\", id.%[2]s),", humanReadableKey, segment.FieldName))
	}

	reversedKeys := make([]string, 0)
	for i := len(formatKeys); i != 0; i-- {
		reversedKeys = append(reversedKeys, formatKeys[i-1])
	}

	formatKeysString := strings.Join(reversedKeys, "\n")
	return fmt.Sprintf(`
func (id %[1]sId) String() string {
	segments := []string{
%s
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%%s: (%%s)", %[3]q, segmentsStr)
}
`, id.TypeName, formatKeysString, makeHumanReadable(id.TypeName))
}

func (id ResourceIdGenerator) codeForFormatter() string {
	formatKeys := make([]string, 0)
	for _, segment := range id.Segments {
		formatKeys = append(formatKeys, fmt.Sprintf("id.%s", segment.FieldName))
	}
	formatKeysString := strings.Join(formatKeys, ", ")
	return fmt.Sprintf(`
func (id %[1]sId) ID() string {
	fmtString := %[2]q
	return fmt.Sprintf(fmtString, %[3]s)
}
`, id.TypeName, id.IDFmt, formatKeysString)
}

func (id ResourceIdGenerator) codeForParser() string {
	directAssignments := make([]string, 0)
	if id.HasSubscriptionId {
		directAssignments = append(directAssignments, "\t\tSubscriptionId: id.SubscriptionID,")
	}
	if id.HasResourceGroup {
		directAssignments = append(directAssignments, "\t\tResourceGroup: id.ResourceGroup,")
	}
	directAssignmentsStr := strings.Join(directAssignments, "\n")

	parserStatements := make([]string, 0)
	for _, segment := range id.Segments {
		isSubscription := strings.EqualFold(segment.FieldName, "SubscriptionId") && id.HasSubscriptionId
		isResourceGroup := strings.EqualFold(segment.FieldName, "ResourceGroup") && id.HasResourceGroup
		if isSubscription || isResourceGroup {
			parserStatements = append(parserStatements, fmt.Sprintf(`
	if resourceId.%[1]s == "" {
		return nil, errors.New("ID was missing the '%[2]s' element")
	}
`, segment.FieldName, segment.SegmentKey))
			continue
		}

		fmtString := "\tif resourceId.%[1]s, err = id.PopSegment(\"%[2]s\"); err != nil {\n\t\treturn nil, err\n\t}"
		parserStatements = append(parserStatements, fmt.Sprintf(fmtString, segment.FieldName, segment.SegmentKey))
	}
	parserStatementsStr := strings.Join(parserStatements, "\n")
	return fmt.Sprintf(`
// %[1]sID parses a %[1]s ID into an %[1]sId struct 
func %[1]sID(input string) (*%[1]sId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %%q as an %[1]s ID: %%+v", input, err)
	}

	resourceId := %[1]sId{
%[2]s
	}

%[3]s

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
`, id.TypeName, directAssignmentsStr, parserStatementsStr)
}

func (id ResourceIdGenerator) codeForParserInsensitive() string {
	if !id.ShouldRewrite {
		// this only exists to workaround broken API's to patch those ID's, so shouldn't be used in most circumstances
		return ""
	}

	directAssignments := make([]string, 0)
	if id.HasSubscriptionId {
		directAssignments = append(directAssignments, "\t\tSubscriptionId: id.SubscriptionID,")
	}
	if id.HasResourceGroup {
		directAssignments = append(directAssignments, "\t\tResourceGroup: id.ResourceGroup,")
	}
	directAssignmentsStr := strings.Join(directAssignments, "\n")

	parserStatements := make([]string, 0)
	for _, segment := range id.Segments {
		isSubscription := strings.EqualFold(segment.FieldName, "SubscriptionId") && id.HasSubscriptionId
		isResourceGroup := strings.EqualFold(segment.FieldName, "ResourceGroup") && id.HasResourceGroup
		if isSubscription || isResourceGroup {
			parserStatements = append(parserStatements, fmt.Sprintf(`
	if resourceId.%[1]s == "" {
		return nil, errors.New("ID was missing the '%[2]s' element")
	}
`, segment.FieldName, segment.SegmentKey))
			continue
		}

		// NOTE: This becomes dramatically simpler long-term - but for now has to be long-winded
		// to avoid subtle changes to resources until this is threaded through everywhere
		fmtString := `
  // find the correct casing for the '%[2]s' segment
  %[2]sKey := "%[2]s"
  for key := range id.Path {
  	if strings.EqualFold(key, %[2]sKey) {
  		%[2]sKey = key
  		break
  	}
  }
  if resourceId.%[1]s, err = id.PopSegment(%[2]sKey); err != nil {
    return nil, err
  }
`
		parserStatements = append(parserStatements, fmt.Sprintf(fmtString, segment.FieldName, segment.SegmentKey))
	}
	parserStatementsStr := strings.Join(parserStatements, "\n")
	return fmt.Sprintf(`
// %[1]sIDInsensitively parses an %[1]s ID into an %[1]sId struct, insensitively
// This should only be used to parse an ID for rewriting, the %[1]sID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func %[1]sIDInsensitively(input string) (*%[1]sId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := %[1]sId{
%[2]s
	}

%[3]s

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
`, id.TypeName, directAssignmentsStr, parserStatementsStr)
}

func (id ResourceIdGenerator) TestCode() string {
	importLine := ""
	if id.TestPackageSuffix != "" {
		importLine = fmt.Sprintf("\"github.com/hashicorp/terraform-provider-azurerm/internal/services/%s/parse\"", id.ServicePackageName)
	}

	return fmt.Sprintf(`
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse%s

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	%s
)

%s
%s
%s
`, id.TestPackageSuffix, importLine, id.testCodeForFormatter(), id.testCodeForParser(), id.testCodeForParserInsensitive())
}

func (id ResourceIdGenerator) testCodeForFormatter() string {
	arguments := make([]string, 0)
	for _, segment := range id.Segments {
		arguments = append(arguments, fmt.Sprintf("%q", segment.SegmentValue))
	}
	argumentsStr := strings.Join(arguments, ", ")
	if id.TestPackageSuffix == "" {
		return fmt.Sprintf(`
var _ resourceids.Id = %[1]sId{}

func Test%[1]sIDFormatter(t *testing.T) {
	actual := New%[1]sID(%[2]s).ID()
	expected := %[3]q
	if actual != expected {
		t.Fatalf("Expected %%q but got %%q", expected, actual)
	}
}
`, id.TypeName, argumentsStr, id.IDRaw)
	}

	return fmt.Sprintf(`
var _ resourceid.Formatter = parse.%[1]sId{}

func Test%[1]sIDFormatter(t *testing.T) {
	actual := parse.New%[1]sID(%[2]s).ID()
	expected := %[3]q
	if actual != expected {
		t.Fatalf("Expected %%q but got %%q", expected, actual)
	}
}
`, id.TypeName, argumentsStr, id.IDRaw)
}

func (id ResourceIdGenerator) testCodeForParser() string {
	testCases := make([]string, 0)
	testCases = append(testCases, `
		{
			// empty
			Input: "",
			Error: true,
		},
`)
	assignmentChecks := make([]string, 0)
	for _, segment := range id.Segments {
		testCaseFmt := `
		{
			// missing %s
			Input: %q,
			Error: true,
		},`
		// missing the key
		resourceIdToThisPointIndex := strings.Index(id.IDRaw, segment.SegmentKey)
		resourceIdToThisPoint := id.IDRaw[0:resourceIdToThisPointIndex]
		testCases = append(testCases, fmt.Sprintf(testCaseFmt, segment.FieldName, resourceIdToThisPoint))

		// missing the value
		resourceIdToThisPointIndex = strings.Index(id.IDRaw, segment.SegmentValue)
		resourceIdToThisPoint = id.IDRaw[0:resourceIdToThisPointIndex]
		testCases = append(testCases, fmt.Sprintf(testCaseFmt, fmt.Sprintf("value for %s", segment.FieldName), resourceIdToThisPoint))

		assignmentsFmt := "\t\tif actual.%[1]s != v.Expected.%[1]s {\n\t\t\tt.Fatalf(\"Expected %%q but got %%q for %[1]s\", v.Expected.%[1]s, actual.%[1]s)\n\t\t}"
		assignmentChecks = append(assignmentChecks, fmt.Sprintf(assignmentsFmt, segment.FieldName))
	}

	// add a successful test case
	expectAssignments := make([]string, 0)
	for _, segment := range id.Segments {
		expectAssignments = append(expectAssignments, fmt.Sprintf("\t\t\t\t%s:\t%q,", segment.FieldName, segment.SegmentValue))
	}
	typeName := fmt.Sprintf("%sId", id.TypeName)
	if id.TestPackageSuffix != "" {
		typeName = fmt.Sprintf("parse.%s", typeName)
	}
	testCases = append(testCases, fmt.Sprintf(`
		{
			// valid
			Input: "%[1]s",
			Expected: &%[2]s{
%[3]s
			},
		},
`, id.IDRaw, typeName, strings.Join(expectAssignments, "\n")))

	// add an intentionally failing upper-cased test case
	testCases = append(testCases, fmt.Sprintf(`
		{
			// upper-cased
			Input: %q,
			Error: true,
		},`, strings.ToUpper(id.IDRaw)))

	testCasesStr := strings.Join(testCases, "\n")
	assignmentCheckStr := strings.Join(assignmentChecks, "\n")

	if id.TestPackageSuffix == "" {
		return fmt.Sprintf(`
func Test%[1]sID(t *testing.T) {
	testData := []struct {
		Input  string
		Error  bool
		Expected *%[1]sId
	}{%[2]s
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %%q", v.Input)

		actual, err := %[1]sID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %%s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

%[3]s
	}
}
`, id.TypeName, testCasesStr, assignmentCheckStr)
	}

	return fmt.Sprintf(`
func Test%[1]sID(t *testing.T) {
	testData := []struct {
		Input  string
		Error  bool
		Expected *parse.%[1]sId
	}{%[2]s
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %%q", v.Input)

		actual, err := parse.%[1]sID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %%s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

%[3]s
	}
}
`, id.TypeName, testCasesStr, assignmentCheckStr)
}

func (id ResourceIdGenerator) testCodeForParserInsensitive() string {
	if !id.ShouldRewrite {
		// this functionality isn't enabled by default
		return ""
	}

	testCases := make([]string, 0)
	testCases = append(testCases, `
		{
			// empty
			Input: "",
			Error: true,
		},
`)
	assignmentChecks := make([]string, 0)
	for _, segment := range id.Segments {
		testCaseFmt := `
		{
			// missing %s
			Input: %q,
			Error: true,
		},`
		// missing the key
		resourceIdToThisPointIndex := strings.Index(id.IDRaw, segment.SegmentKey)
		resourceIdToThisPoint := id.IDRaw[0:resourceIdToThisPointIndex]
		testCases = append(testCases, fmt.Sprintf(testCaseFmt, segment.FieldName, resourceIdToThisPoint))

		// missing the value
		resourceIdToThisPointIndex = strings.Index(id.IDRaw, segment.SegmentValue)
		resourceIdToThisPoint = id.IDRaw[0:resourceIdToThisPointIndex]
		testCases = append(testCases, fmt.Sprintf(testCaseFmt, fmt.Sprintf("value for %s", segment.FieldName), resourceIdToThisPoint))

		assignmentsFmt := "\t\tif actual.%[1]s != v.Expected.%[1]s {\n\t\t\tt.Fatalf(\"Expected %%q but got %%q for %[1]s\", v.Expected.%[1]s, actual.%[1]s)\n\t\t}"
		assignmentChecks = append(assignmentChecks, fmt.Sprintf(assignmentsFmt, segment.FieldName))
	}

	// add a successful test case
	expectAssignments := make([]string, 0)
	for _, segment := range id.Segments {
		expectAssignments = append(expectAssignments, fmt.Sprintf("\t\t\t\t%s:\t%q,", segment.FieldName, segment.SegmentValue))
	}
	testCases = append(testCases, fmt.Sprintf(`
		{
			// valid
			Input: "%[1]s",
			Expected: &%[2]sId{
%[3]s
			},
		},
`, id.IDRaw, id.TypeName, strings.Join(expectAssignments, "\n")))

	testCaseWithTransformation := func(testCaseName string, transform func(in string) string) string {
		resourceIdWithTransform := id.IDRaw
		for _, segment := range id.Segments {
			// we're not as concerned with these two for now
			if segment.FieldName == "SubscriptionId" || segment.FieldName == "ResourceGroup" {
				continue
			}

			transformedKey := transform(segment.SegmentKey)
			resourceIdWithTransform = strings.Replace(resourceIdWithTransform, segment.SegmentKey, transformedKey, 1)
		}

		typeName := fmt.Sprintf("%sId", id.TypeName)
		if id.TestPackageSuffix != "" {
			typeName = fmt.Sprintf("parse.%s", typeName)
		}
		return fmt.Sprintf(`
		{
			// %[4]s
			Input: "%[1]s",
			Expected: &%[2]s{
%[3]s
			},
		},`, resourceIdWithTransform, typeName, strings.Join(expectAssignments, "\n"), testCaseName)
	}

	testCases = append(testCases, testCaseWithTransformation("lower-cased segment names", strings.ToLower))
	testCases = append(testCases, testCaseWithTransformation("upper-cased segment names", strings.ToUpper))
	testCases = append(testCases, testCaseWithTransformation("mixed-cased segment names", func(in string) string {
		out := make([]rune, 0)
		for i, c := range in {
			if i%2 == 0 {
				out = append(out, unicode.ToUpper(c))
			} else {
				out = append(out, unicode.ToLower(c))
			}
		}
		return string(out)
	}))

	testCasesStr := strings.Join(testCases, "\n")
	assignmentCheckStr := strings.Join(assignmentChecks, "\n")

	if id.TestPackageSuffix == "" {
		return fmt.Sprintf(`
func Test%[1]sIDInsensitively(t *testing.T) {
	testData := []struct {
		Input  string
		Error  bool
		Expected *%[1]sId
	}{%[2]s
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %%q", v.Input)

		actual, err := %[1]sIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %%s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

%[3]s
	}
}
`, id.TypeName, testCasesStr, assignmentCheckStr)
	}

	return fmt.Sprintf(`
func Test%[1]sIDInsensitively(t *testing.T) {
	testData := []struct {
		Input  string
		Error  bool
		Expected *parse.%[1]sId
	}{%[2]s
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %%q", v.Input)

		actual, err := parse.%[1]sIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %%s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

%[3]s
	}
}
`, id.TypeName, testCasesStr, assignmentCheckStr)
}

func (id ResourceIdGenerator) ValidatorCode() string {
	return fmt.Sprintf(`
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/%[2]s/parse"
)

func %[1]sID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %%q to be a string", key))
		return
	}

	if _, err := parse.%[1]sID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
`, id.TypeName, id.ServicePackageName)
}

func (id ResourceIdGenerator) ValidatorTestCode() string {
	testCases := make([]string, 0)
	testCases = append(testCases, `
		{
			// empty
			Input: "",
			Valid: false,
		},
`)
	for _, segment := range id.Segments {
		testCaseFmt := `
		{
			// missing %s
			Input: %q,
			Valid: false,
		},`
		// missing the key
		resourceIdToThisPointIndex := strings.Index(id.IDRaw, segment.SegmentKey)
		resourceIdToThisPoint := id.IDRaw[0:resourceIdToThisPointIndex]
		testCases = append(testCases, fmt.Sprintf(testCaseFmt, segment.FieldName, resourceIdToThisPoint))

		// missing the value
		resourceIdToThisPointIndex = strings.Index(id.IDRaw, segment.SegmentValue)
		resourceIdToThisPoint = id.IDRaw[0:resourceIdToThisPointIndex]
		testCases = append(testCases, fmt.Sprintf(testCaseFmt, fmt.Sprintf("value for %s", segment.FieldName), resourceIdToThisPoint))
	}

	// add a successful test case
	testCases = append(testCases, fmt.Sprintf(`
		{
			// valid
			Input: %q,
			Valid: true,
		},
`, id.IDRaw))

	// add an intentionally failing upper-cased test case
	testCases = append(testCases, fmt.Sprintf(`
		{
			// upper-cased
			Input: %q,
			Valid: false,
		},`, strings.ToUpper(id.IDRaw)))

	testCasesStr := strings.Join(testCases, "\n")

	if id.TestPackageSuffix == "" {
		return fmt.Sprintf(`
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func Test%[1]sID(t *testing.T) {
	cases := []struct {
		Input    string
		Valid bool
	}{%[2]s
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %%s", tc.Input)
		_, errors := %[1]sID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %%t but got %%t", tc.Valid, valid)
		}
	}
}
`, id.TypeName, testCasesStr)
	}

	return fmt.Sprintf(`// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate%[1]s

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/%[4]s/validate"
)

func Test%[2]sID(t *testing.T) {
	cases := []struct {
		Input    string
		Valid bool
	}{%[3]s
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %%s", tc.Input)
		_, errors := validate.%[2]sID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %%t but got %%t", tc.Valid, valid)
		}
	}
}
`, id.TestPackageSuffix, id.TypeName, testCasesStr, id.ServicePackageName)
}

func goFmtAndWriteToFile(filePath, fileContents string) error {
	fmt, err := GolangCodeFormatter{}.Format(fileContents)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, []byte(*fmt), 0o644); err != nil {
		return err
	}

	return nil
}

type GolangCodeFormatter struct{}

func (f GolangCodeFormatter) Format(input string) (*string, error) {
	tmpfile, err := os.CreateTemp("", "temp-*.go")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %+v", err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	filePath := tmpfile.Name()

	if _, err := tmpfile.WriteString(input); err != nil {
		return nil, fmt.Errorf("writing contents to %q: %+v", filePath, err)
	}

	f.runGoFmt(filePath)
	f.runGoImports(filePath)

	contents, err := f.readFileContents(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading contents from %q: %+v", filePath, err)
	}

	return contents, nil
}

func (f GolangCodeFormatter) runGoFmt(filePath string) {
	cmd := exec.Command("gofmt", "-w", filePath)
	// intentionally not using these errors since the exit codes are kinda uninteresting
	_ = cmd.Start()
	_ = cmd.Wait()
}

func (f GolangCodeFormatter) runGoImports(filePath string) {
	cmd := exec.Command("goimports", "-w", filePath)
	// intentionally not using these errors since the exit codes are kinda uninteresting
	_ = cmd.Start()
	_ = cmd.Wait()
}

func (f GolangCodeFormatter) readFileContents(filePath string) (*string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contents := string(data)
	return &contents, nil
}
