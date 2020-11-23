package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
)

// NOTE: since we're using `go run` for these tools all of the code needs to live within the main.go

func main() {
	f := flag.NewFlagSet("example", flag.ExitOnError)

	resourceName := f.String("name", "", "The name of the Data Source/Resource which should be generated")
	brandName := f.String("brand-name", "", "The friendly/brand name of this Data Source/Resource (e.g. Resource Group)")
	resourceId := f.String("resource-id", "", "An Azure Resource ID showing an example of how to Import this Resource")
	resourceType := f.String("type", "", "Whether this is a Data Source (data) or a Resource (resource)")
	websitePath := f.String("website-path", "", "The relative path to the website folder")

	_ = f.Parse(os.Args[1:])

	quitWithError := func(message string) {
		log.Print(message)
		os.Exit(1)
	}

	if resourceName == nil || *resourceName == "" {
		quitWithError("The name of the Data Source/Resource must be specified via `-name`")
		return
	}

	if brandName == nil || *brandName == "" {
		quitWithError("The friendly/brannd name of the Data Source/Resource must be specified via `-brand`")
		return
	}

	if resourceType == nil || *resourceType == "" {
		quitWithError("The type of the Data Source/Resource must be specified via `-type`")
		return
	}

	if *resourceType != "data" && *resourceType != "resource" {
		quitWithError("The type of the Data Source/Resource specified via `-type` must be either `data` or `resource`")
		return
	}

	if websitePath == nil || *websitePath == "" {
		quitWithError("The Relative Website Path must be specified via `-website-path`")
		return
	}

	isResource := *resourceType == "resource"
	if isResource && (resourceId == nil || *resourceId == "") {
		quitWithError("An example of an Azure Resource ID must be specified via `-resource-id` when scaffolding for a Resource")
		return
	}

	if err := run(*resourceName, *brandName, resourceId, isResource, *websitePath); err != nil {
		panic(err)
	}
}

func run(resourceName, brandName string, resourceId *string, isResource bool, websitePath string) error {
	content, err := getContent(resourceName, brandName, resourceId, isResource)
	if err != nil {
		return fmt.Errorf("Error building content: %s", err)
	}

	return saveContent(resourceName, websitePath, *content, isResource)
}

func getContent(resourceName, brandName string, resourceId *string, isResource bool) (*string, error) {
	generator := documentationGenerator{
		resourceName: resourceName,
		brandName:    brandName,
		resourceId:   resourceId,
		isDataSource: !isResource,
	}

	if !isResource {
		for _, service := range provider.SupportedTypedServices() {
			for _, ds := range service.SupportedDataSources() {
				if ds.ResourceType() == resourceName {
					wrapper := sdk.NewDataSourceWrapper(ds)
					dsWrapper, err := wrapper.DataSource()
					if err != nil {
						return nil, fmt.Errorf("wrapping Data Source %q: %+v", ds.ResourceType(), err)
					}

					generator.resource = dsWrapper
					generator.websiteCategories = service.WebsiteCategories()
					break
				}
			}
		}
		for _, service := range provider.SupportedUntypedServices() {
			for key, ds := range service.SupportedDataSources() {
				if key == resourceName {
					generator.resource = ds
					generator.websiteCategories = service.WebsiteCategories()
					break
				}
			}
		}

		if generator.resource == nil {
			return nil, fmt.Errorf("Data Source %q was not registered!", resourceName)
		}
	} else {
		for _, service := range provider.SupportedTypedServices() {
			for _, rs := range service.SupportedResources() {
				if rs.ResourceType() == resourceName {
					wrapper := sdk.NewResourceWrapper(rs)
					rsWrapper, err := wrapper.Resource()
					if err != nil {
						return nil, fmt.Errorf("wrapping Resource %q: %+v", rs.ResourceType(), err)
					}

					generator.resource = rsWrapper
					generator.websiteCategories = service.WebsiteCategories()
					break
				}
			}
		}
		for _, service := range provider.SupportedUntypedServices() {
			for key, rs := range service.SupportedResources() {
				if key == resourceName {
					generator.resource = rs
					generator.websiteCategories = service.WebsiteCategories()
					break
				}
			}
		}

		if generator.resource == nil {
			return nil, fmt.Errorf("Resource %q was not registered!", resourceName)
		}
	}

	docs := generator.generate()
	return &docs, nil
}

func saveContent(resourceName string, websitePath string, content string, isResource bool) error {
	resourceKind := "r"
	if !isResource {
		resourceKind = "d"
	}

	fileName := strings.TrimPrefix(resourceName, "azurerm_")
	outputFileName := fmt.Sprintf("%s/docs/%s/%s.html.markdown", websitePath, resourceKind, fileName)
	outputPath, err := filepath.Abs(outputFileName)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	if os.IsExist(err) {
		os.Remove(outputPath)
		file, err = os.Create(outputPath)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	content = strings.TrimSpace(content)
	_, _ = file.WriteString(content)
	return file.Sync()
}

type documentationGenerator struct {
	resource *schema.Resource

	// brandName is the marketing brand name used for this resource (e.g. Resource Group / App Service / Web Apps)
	brandName string

	// resourceName is the name of the resource e.g. `azurerm_resource_group`
	resourceName string

	// isDataSource defines if this is a Data Source (if not it's a Resource)
	isDataSource bool

	// resourceId is an example of the ID used by this Resource
	resourceId *string

	// websiteCategories is the list of categories available by this service definition
	websiteCategories []string
}

func (gen documentationGenerator) generate() string {
	title := gen.title()
	argumentsBlock := gen.argumentsBlock()
	attributesBlock := gen.attributesBlock()
	description := gen.description()
	exampleUsageBlock := gen.exampleUsageBlock()
	frontMatterBlock := gen.frontMatterBlock()
	importBlock := gen.importBlock()
	timeoutsBlock := gen.timeoutsBlock()

	template := fmt.Sprintf(`%s
# %s

%s.

## Example Usage

[][][]hcl
%s
[][][]

%s

%s

%s

%s`, frontMatterBlock, title, description, exampleUsageBlock, argumentsBlock, attributesBlock, timeoutsBlock, importBlock)
	return strings.ReplaceAll(template, "[][][]", "```")
}

// blocks
func (gen documentationGenerator) argumentsBlock() string {
	documentationForArguments := func(input map[string]*schema.Schema, onlyRequired, onlyOptional bool, blockName string) string {
		fields := ""

		for _, fieldName := range gen.sortFields(input) {
			field := input[fieldName]

			// nothing to see here, move along
			if !field.Optional && !field.Required {
				continue
			}

			if onlyRequired && !field.Required {
				continue
			}

			if onlyOptional && !field.Optional {
				continue
			}

			status := "Optional"
			if field.Required {
				status = "Required"
			}

			value := gen.buildDescriptionForArgument(fieldName, field, blockName)
			if len(field.ConflictsWith) > 0 {
				conflictingValues := make([]string, 0)
				for _, v := range field.ConflictsWith {
					conflictingValues = append(conflictingValues, fmt.Sprintf("`%s`", v))
				}

				value += fmt.Sprintf("Conflicts with %s", strings.Join(conflictingValues, ","))
			}
			if field.ForceNew {
				value += fmt.Sprintf(" Changing this forces a new %s to be created.", gen.brandName)
			}
			fields += fmt.Sprintf("* `%s` - (%s) %s\n\n", fieldName, status, value)
		}

		return fields
	}

	// first output the Required fields
	fields := documentationForArguments(gen.resource.Schema, true, false, "")
	// then prepare the Optional fields
	optionalFields := documentationForArguments(gen.resource.Schema, false, true, "")

	// assuming we have both optional & required fields - let's add a separarer
	if len(fields) > 0 && len(optionalFields) > 0 {
		fields += "---\n\n"
	}
	fields += optionalFields

	// first list all of the top-level fields / blocks alphabetically

	// then we need to collect a list of all block names, everywhere
	blockNames, blocks := gen.uniqueBlockNamesForArgument(gen.resource.Schema)

	for _, blockName := range blockNames {
		block := blocks[blockName]

		fields += "---\n\n"
		fields += fmt.Sprintf("A `%s` block supports the following:\n\n", blockName)
		// required
		fields += documentationForArguments(block, true, false, blockName)
		// optional
		fields += documentationForArguments(block, false, true, blockName)
	}

	fields = strings.TrimSuffix(fields, "\n\n")

	return fmt.Sprintf(`## Arguments Reference

The following arguments are supported:

%s`, fields)
}

func (gen documentationGenerator) attributesBlock() string {
	documentationForAttributes := func(input map[string]*schema.Schema, onlyComputed bool, blockName string) string {
		fields := ""

		// now list all of the top-level fields / blocks alphabetically
		for _, fieldName := range gen.sortFields(input) {
			field := input[fieldName]
			// when we're in a nested block there's no need to duplicate the fields
			if onlyComputed && !field.Computed {
				continue
			}
			if onlyComputed && (field.Optional || field.Required) {
				continue
			}

			value := gen.buildDescriptionForAttribute(fieldName, field, blockName)
			fields += fmt.Sprintf("* `%s` - %s\n\n", fieldName, value)
		}

		return fields
	}

	// present in everything
	fields := fmt.Sprintf("* `id` - The ID of the %s.\n\n", gen.brandName)

	// now list all of the top-level fields / blocks alphabetically
	fields += documentationForAttributes(gen.resource.Schema, true, "")

	// then we need to collect a list of all block names, everywhere
	blockNames, blocks := gen.uniqueBlockNamesForAttribute(gen.resource.Schema)

	for _, blockName := range blockNames {
		block := blocks[blockName]

		fields += "---\n\n"
		fields += fmt.Sprintf("A `%s` block exports the following:\n\n", blockName)
		fields += documentationForAttributes(block, false, blockName)
	}

	fields = strings.TrimSuffix(fields, "\n\n")

	return fmt.Sprintf(`## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

%s`, fields)
}

func (gen documentationGenerator) description() string {
	if gen.isDataSource {
		return fmt.Sprintf("Use this data source to access information about an existing %s", gen.brandName)
	}

	return fmt.Sprintf("Manages a %s", gen.brandName)
}

func (gen documentationGenerator) exampleUsageBlock() string {
	requiredFields := gen.requiredFieldsForExampleBlock(gen.resource.Schema, 1)

	if gen.isDataSource {
		return fmt.Sprintf(`data "%s" "example" {
%s
}

output "id" {
  value = data.%s.example.id
}`, gen.resourceName, requiredFields, gen.resourceName)
	}

	return fmt.Sprintf(`resource "%s" "example" {
%s
}`, gen.resourceName, requiredFields)
}

func (gen documentationGenerator) frontMatterBlock() string {
	category := "TODO"
	if len(gen.websiteCategories) > 0 {
		if len(gen.websiteCategories) == 1 {
			category = gen.websiteCategories[0]
		} else {
			category = fmt.Sprintf("TODO - pick from: %s", strings.Join(gen.websiteCategories, "|"))
		}
	}

	title := gen.title()
	var description string
	if gen.isDataSource {
		description = fmt.Sprintf("Gets information about an existing %s", gen.brandName)
	} else {
		description = fmt.Sprintf("Manages a %s", gen.brandName)
	}

	return fmt.Sprintf(`
---
subcategory: "%s"
layout: "azurerm"
page_title: "Azure Resource Manager: %s"
description: |-
  %s.
---
`, category, title, description)
}

func (gen documentationGenerator) importBlock() string {
	// data source don't support import
	if gen.isDataSource {
		return ""
	}

	template := fmt.Sprintf(`## Import

%ss can be imported using the []resource id[], e.g.

[][][]shell
terraform import %s.example %s
[][][]`, gen.brandName, gen.resourceName, *gen.resourceId)
	return strings.ReplaceAll(template, "[]", "`")
}

func (gen documentationGenerator) timeoutsBlock() string {
	if gen.resource.Timeouts == nil {
		return ""
	}
	timeouts := *gen.resource.Timeouts

	timeoutsBlurb := "The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:"

	timeoutToFriendlyText := func(duration time.Duration) string {
		hours := int(math.Floor(duration.Hours()))
		if hours > 0 {
			var hoursText string
			if hours > 1 {
				hoursText = fmt.Sprintf("%d hours", hours)
			} else {
				hoursText = "1 hour"
			}

			minutesRemaining := int(math.Floor(duration.Minutes())) % 60.0
			if minutesRemaining == 0 {
				return hoursText
			}

			var minutesText string
			if minutesRemaining > 1 {
				minutesText = fmt.Sprintf("%d minutes", minutesRemaining)
			} else {
				minutesText = "1 minute"
			}

			return fmt.Sprintf("%s and %s", hoursText, minutesText)
		}

		minutes := int(duration.Minutes())
		if minutes > 1 {
			return fmt.Sprintf("%d minutes", minutes)
		}

		return "1 minute"
	}

	timeoutsText := ""
	if timeouts.Create != nil {
		friendlyText := timeoutToFriendlyText(*timeouts.Create)
		timeoutsText += fmt.Sprintf("* `create` - (Defaults to %s) Used when creating the %s.\n", friendlyText, gen.brandName)
	}

	if timeouts.Read != nil {
		friendlyText := timeoutToFriendlyText(*timeouts.Read)
		timeoutsText += fmt.Sprintf("* `read` - (Defaults to %s) Used when retrieving the %s.\n", friendlyText, gen.brandName)
	}

	if timeouts.Update != nil {
		friendlyText := timeoutToFriendlyText(*timeouts.Update)
		timeoutsText += fmt.Sprintf("* `update` - (Defaults to %s) Used when updating the %s.\n", friendlyText, gen.brandName)
	}

	if timeouts.Delete != nil {
		friendlyText := timeoutToFriendlyText(*timeouts.Delete)
		timeoutsText += fmt.Sprintf("* `delete` - (Defaults to %s) Used when deleting the %s.\n", friendlyText, gen.brandName)
	}

	timeoutsText = strings.TrimSuffix(timeoutsText, "\n")
	return fmt.Sprintf(`## Timeouts

%s

%s`, timeoutsBlurb, timeoutsText)
}

func (gen documentationGenerator) title() string {
	if gen.isDataSource {
		return fmt.Sprintf("Data Source: %s", gen.resourceName)
	}

	return gen.resourceName
}

// helpers
func (gen documentationGenerator) blockIsBefore(name string, blockName string) bool {
	if blockName == "" {
		return false
	}

	items := []string{name, blockName}
	sort.Strings(items)
	return items[0] == name
}

func (gen documentationGenerator) buildIndentForExample(level int) string {
	out := ""
	for i := 0; i < level; i++ {
		out += "  "
	}
	return out
}

func (gen documentationGenerator) buildDescriptionForArgument(name string, field *schema.Schema, blockName string) string {
	if name == "name" {
		if blockName == "" {
			if gen.isDataSource {
				return fmt.Sprintf("The name of this %s.", gen.brandName)
			}

			return fmt.Sprintf("The name which should be used for this %s.", gen.brandName)
		} else {
			return "The name which should be used for this TODO."
		}
	}
	if name == "location" {
		if gen.isDataSource {
			return fmt.Sprintf("The Azure Region where the %s exists.", gen.brandName)
		}

		return fmt.Sprintf("The Azure Region where the %s should exist.", gen.brandName)
	}
	if name == "resource_group_name" {
		if gen.isDataSource {
			return fmt.Sprintf("The name of the Resource Group where the %s exists.", gen.brandName)
		}

		return fmt.Sprintf("The name of the Resource Group where the %s should exist.", gen.brandName)
	}
	if name == "tags" {
		return fmt.Sprintf("A mapping of tags which should be assigned to the %s.", gen.brandName)
	}

	if name == "enabled" || strings.HasSuffix(name, "_enabled") {
		return "Should the TODO be enabled?"
	}

	if strings.HasSuffix(name, "_id") {
		return "The ID of the TODO."
	}

	if field.Elem != nil {
		if _, ok := field.Elem.(*schema.Resource); ok {
			fmtBlock := func(name string, maxItem int, blockIsBefore bool) string {
				var head string
				if maxItem == 1 {
					head = fmt.Sprintf("A `%s` block", name)
				} else {
					head = fmt.Sprintf("One or more `%s` blocks", name)
				}

				var tail string
				if blockIsBefore {
					tail = "as defined above."
				} else {
					tail = "as defined below."
				}
				return head + " " + tail
			}
			return fmtBlock(name, field.MaxItems, gen.blockIsBefore(name, blockName))
		}
	}

	switch field.Type {
	case schema.TypeList, schema.TypeSet, schema.TypeMap:
		return "Specifies a list of TODO."
	}

	return "TODO."
}

func (gen documentationGenerator) buildDescriptionForAttribute(name string, field *schema.Schema, blockName string) string {
	if name == "name" {
		if blockName == "" {
			return fmt.Sprintf("The name of this %s.", gen.brandName)
		} else {
			return "The name of this TODO."
		}
	}
	if name == "location" {
		return fmt.Sprintf("The Azure Region where the %s exists.", gen.brandName)
	}
	if name == "resource_group_name" {
		return fmt.Sprintf("The name of the Resource Group where the %s is located.", gen.brandName)
	}
	if name == "tags" {
		return fmt.Sprintf("A mapping of tags assigned to the %s.", gen.brandName)
	}

	if name == "enabled" || strings.HasSuffix(name, "_enabled") {
		return "Is the TODO enabled?"
	}

	if strings.HasSuffix(name, "_id") {
		return "The ID of the TODO."
	}

	if field.Elem != nil {
		if _, ok := field.Elem.(*schema.Schema); ok {
			if gen.blockIsBefore(name, blockName) {
				return fmt.Sprintf("A `%s` block as defined above.", name)
			} else {
				return fmt.Sprintf("A `%s` block as defined below.", name)
			}
		}

		if _, ok := field.Elem.(*schema.Resource); ok {
			if gen.blockIsBefore(name, blockName) {
				return fmt.Sprintf("A `%s` block as defined above.", name)
			} else {
				return fmt.Sprintf("A `%s` block as defined below.", name)
			}
		}
	}
	if field.Type == schema.TypeList {
		if gen.blockIsBefore(name, blockName) {
			return fmt.Sprintf("A `%s` block as defined above.", name)
		} else {
			return fmt.Sprintf("A `%s` block as defined below.", name)
		}
	}

	return "TODO."
}

func (gen documentationGenerator) determineDefaultValueForExample(name string, field *schema.Schema) string {
	if field.Default != nil {
		if v, ok := field.Default.(bool); ok {
			return strconv.FormatBool(v)
		}

		if v, ok := field.Default.(int); ok {
			return fmt.Sprintf("%d", v)
		}

		if v, ok := field.Default.(string); ok {
			return v
		}
	}

	if field.Type == schema.TypeBool {
		return "false"
	}

	if field.Type == schema.TypeInt {
		return "42"
	}

	if field.Type == schema.TypeFloat {
		return "1.23456"
	}

	if name == "name" || strings.HasSuffix(name, "_name") {
		if gen.isDataSource {
			return "\"existing\""
		}

		return "\"example\""
	}
	if name == "location" {
		return "\"West Europe\""
	}
	if name == "resource_group_name" {
		return "\"example-resources\""
	}

	return "\"TODO\""
}

func (gen documentationGenerator) distinctBlockNames(input []string) []string {
	// this is a delightful hack to work around multiple blocks being a thing
	temp := make(map[string]struct{})
	for _, v := range input {
		temp[v] = struct{}{}
	}

	output := make([]string, 0)
	for k := range temp {
		output = append(output, k)
	}

	return output
}

func (gen documentationGenerator) processElementForExample(field string, indentLevel int, elem interface{}, isAttribute bool) string {
	indent := gen.buildIndentForExample(indentLevel)

	// it's an array of something, work out what
	if array, ok := elem.(*schema.Schema); ok {
		switch array.Type {
		case schema.TypeString:
			return fmt.Sprintf("%s%s = [ \"example\" ]\n", indent, field)

		case schema.TypeInt:
			return fmt.Sprintf("%s%s = [ 1 ]\n", indent, field)

		default:
			return "TODO"
		}
	}

	// otherwise it's a list so we're gonna have to go around
	if list, ok := elem.(*schema.Resource); ok && len(list.Schema) > 0 {
		innerFields := gen.requiredFieldsForExampleBlock(list.Schema, indentLevel+1)
		attributeSyntax := ""
		if isAttribute {
			attributeSyntax = " ="
		}
		return fmt.Sprintf("\n%s%s%s {\n%s  %s\n%s}\n", indent, field, attributeSyntax, innerFields, indent, indent)
	}

	// unless something's broken, since this is likely during provider dev it's likely things could be missing
	panic("Field %q has an Element but isn't a Set or List - double-check the schema")
}

func (gen documentationGenerator) requiredFieldsForExampleBlock(fields map[string]*schema.Schema, indentLevel int) string {
	indent := gen.buildIndentForExample(indentLevel)
	output := ""

	processField := func(name string, field *schema.Schema) string {
		value := gen.determineDefaultValueForExample(name, field)
		return fmt.Sprintf("%s%s = %s\n", indent, name, value)
	}

	// if we have a "name", "location" "resource_group_name" field output those first as per convention
	if v, ok := fields["name"]; ok && v.Required {
		output += processField("name", v)
	}
	if v, ok := fields["resource_group_name"]; ok && v.Required {
		output += processField("resource_group_name", v)
	}
	if v, ok := fields["location"]; ok && v.Required {
		output += processField("location", v)
	}

	for field, v := range fields {
		if !v.Required {
			continue
		}
		if field == "location" || field == "name" || field == "resource_group_name" {
			continue
		}

		if v.Elem != nil {
			isAttribute := v.ConfigMode == schema.SchemaConfigModeAttr
			output += gen.processElementForExample(field, indentLevel, v.Elem, isAttribute)
			continue
		}

		output += processField(field, v)
	}
	return strings.TrimSuffix(output, "\n")
}

func (gen documentationGenerator) sortFields(input map[string]*schema.Schema) []string {
	fieldNames := make([]string, 0)
	for field := range input {
		fieldNames = append(fieldNames, field)
	}
	sort.Strings(fieldNames)
	return fieldNames
}

func (gen documentationGenerator) uniqueBlockNamesForArgument(fields map[string]*schema.Schema) ([]string, map[string]map[string]*schema.Schema) {
	blockNames := make([]string, 0)
	blocks := make(map[string]map[string]*schema.Schema)

	for _, fieldName := range gen.sortFields(fields) {
		field := fields[fieldName]

		// compute-only fields can be omitted
		if field.Computed && !(field.Optional || field.Required) {
			continue
		}

		if field.Type != schema.TypeList && field.Type != schema.TypeSet {
			continue
		}

		if field.Elem == nil {
			continue
		}
		v, ok := field.Elem.(*schema.Resource)
		if !ok {
			continue
		}
		if v == nil {
			continue
		}

		// add this block
		blockNames = append(blockNames, fieldName)
		blocks[fieldName] = v.Schema

		// at this point we want to iterate over all the fields to determine which ones are nested blocks, then iterate over/aggregate those
		for _, innerElem := range v.Schema {
			if innerElem.Type != schema.TypeList && innerElem.Type != schema.TypeSet {
				continue
			}
			if field.Elem == nil {
				continue
			}

			innerV, ok := field.Elem.(*schema.Resource)
			if !ok {
				continue
			}
			if innerV == nil {
				continue
			}

			innerBlockNames, innerBlocks := gen.uniqueBlockNamesForArgument(innerV.Schema)
			for _, innerBlockName := range innerBlockNames {
				innerBlock := innerBlocks[innerBlockName]

				blockNames = append(blockNames, innerBlockName)
				blocks[innerBlockName] = innerBlock
			}
		}
	}

	blockNames = gen.distinctBlockNames(blockNames)
	sort.Strings(blockNames)

	return blockNames, blocks
}

func (gen documentationGenerator) uniqueBlockNamesForAttribute(fields map[string]*schema.Schema) ([]string, map[string]map[string]*schema.Schema) {
	blockNames := make([]string, 0)
	blocks := make(map[string]map[string]*schema.Schema)

	for _, fieldName := range gen.sortFields(fields) {
		field := fields[fieldName]

		// fields which are setable but aren't computed-only can be skipped
		if (field.Optional || field.Required) && !field.Computed {
			continue
		}

		// optional+computed blocks with fields which aren't computed shouldn't be documented for attributes
		if field.Optional && field.Computed {
			continue
		}

		if field.Type != schema.TypeList && field.Type != schema.TypeSet {
			continue
		}

		if field.Elem == nil {
			continue
		}
		v, ok := field.Elem.(*schema.Resource)
		if !ok {
			continue
		}
		if v == nil {
			continue
		}

		// add this block
		blockNames = append(blockNames, fieldName)
		blocks[fieldName] = v.Schema

		// at this point we want to iterate over all the fields to determine which ones are nested blocks, then iterate over/aggregate those
		for _, innerElem := range v.Schema {
			if innerElem.Type != schema.TypeList && innerElem.Type != schema.TypeSet {
				continue
			}
			if field.Elem == nil {
				continue
			}

			innerV, ok := field.Elem.(*schema.Resource)
			if !ok {
				continue
			}
			if innerV == nil {
				continue
			}

			innerBlockNames, innerBlocks := gen.uniqueBlockNamesForAttribute(innerV.Schema)
			for _, innerBlockName := range innerBlockNames {
				innerBlock := innerBlocks[innerBlockName]

				blockNames = append(blockNames, innerBlockName)
				blocks[innerBlockName] = innerBlock
			}
		}
	}

	blockNames = gen.distinctBlockNames(blockNames)
	sort.Strings(blockNames)

	return blockNames, blocks
}
