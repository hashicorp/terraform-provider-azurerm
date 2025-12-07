// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/differror"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/template"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
	log "github.com/sirupsen/logrus"
)

type S002 struct{}

var (
	_ Rule = S002{}

	// \x60 is a backtick (`), a less strict regex would be nice, but this is basically the only way to find the brand name in a semi-reliable way
	fullTimeoutRegex    = regexp.MustCompile(`(?i)\x60(\w*)\x60[\s\t-]*\([\s\ta-z]*(\d*)[\s\ta-z]*\) used when \w*ing (?:the|this|a)(.*)`)
	partialTimeoutRegex = regexp.MustCompile(`(?i)(\d+)[\s\t]*(?:hours|minutes)?`)
)

func (r S002) ID() string {
	return "S002"
}

func (r S002) Name() string {
	return "Timeouts Section"
}

func (r S002) Description() string {
	return "validates the `Timeouts` section exists, contains correct timeout values, and timeouts are in `CRUD` order"
}

func (r S002) Run(rd *data.TerraformNodeData, fix bool) []error {
	if SkipRule(rd.Type, rd.Name, r.ID()) {
		return nil
	}

	if !rd.Document.Exists {
		return nil
	}

	errs := make([]error, 0)
	resourceTimeouts := timeoutSliceToMap(rd.Timeouts)

	var section *markdown.TimeoutsSection
	for _, sec := range rd.Document.Sections {
		if sec, ok := sec.(*markdown.TimeoutsSection); ok {
			section = sec
			break
		}
	}

	if len(resourceTimeouts) == 0 {
		// add an error?
		return errs
	}

	if section == nil {
		errs = append(errs, fmt.Errorf("%s: Missing Timeouts section", IdAndName(r)))

		if !fix {
			return errs
		}

		section = &markdown.TimeoutsSection{}
		content, err := template.Render(rd, section.Template())
		if err != nil {
			log.WithFields(log.Fields{
				"name": rd.Name,
				"type": rd.Type,
			}).Error(fmt.Errorf("%s: Failed to render template: %+v", IdAndName(r), err))
		}

		rd.Document.HasChange = true
		section.SetContent(content)
		sections, err := markdown.InsertAfterSection(section, rd.Document.Sections, &markdown.AttributesSection{})
		if err != nil {
			log.WithFields(log.Fields{
				"name": rd.Name,
				"type": rd.Type,
			}).Error(fmt.Errorf("%s: Failed to insert new templated section: %+v", IdAndName(r), err))
		}
		rd.Document.Sections = sections
	} else {
		content := section.GetContent()
		foundTimeouts := make(map[data.TimeoutType]int)
		timeoutBrandName := ""
		start, end := 0, 0

		for idx, line := range content {
			if partialTimeoutRegex.MatchString(line) {
				// track start and end of timeout lines, in case we need to insert a new timeout
				// we can insert at end and let the reorder func take care of the rest
				if start == 0 {
					start = idx
				}
				end = idx

				t := parseTimeout(line)
				if t == nil {
					errs = append(errs, fmt.Errorf("%s: Unable to parse timeout line (`%s`), this will require a manual fix", IdAndName(r), line))
					continue
				}
				timeoutBrandName = t.Name

				if _, ok := foundTimeouts[t.Type]; ok {
					errs = append(errs, fmt.Errorf("%s: Documentation contains a duplicate timeout", IdAndName(r)))

					if fix {
						rd.Document.HasChange = true
						content = slices.Delete(content, idx, idx+1)
						section.SetContent(content)
						continue
					}
				}
				foundTimeouts[t.Type] = idx

				if _, ok := resourceTimeouts[t.Type]; !ok {
					errs = append(errs, fmt.Errorf("%s: Documentation contains a timeout (%s) that is not present in the %s", IdAndName(r), t.Type, rd.Type))

					if fix {
						rd.Document.HasChange = true
						content = slices.Delete(content, idx, idx+1)
						section.SetContent(content)
						continue
					}

					return errs
				}

				expectedTimeout := resourceTimeouts[t.Type]
				expectedTimeout.Name = t.Name

				expected := expectedTimeout.String()

				if line != expected {
					errs = append(errs, differror.New(fmt.Sprintf("%s: Timeout line for `%s` not in expected format", IdAndName(r), t.Type), line, expected))

					if fix {
						rd.Document.HasChange = true
						content[idx] = expected
						section.SetContent(content)
					}
				}
			}
		}

		for _, t := range rd.Timeouts {
			if _, ok := foundTimeouts[t.Type]; !ok {
				errs = append(errs, fmt.Errorf("%s: Timeout line for `%s` missing in the documentation", IdAndName(r), t.Type))

				if fix {
					if timeoutBrandName != "" {
						t.Name = timeoutBrandName
					}
					rd.Document.HasChange = true
					content = slices.Insert(content, end+1, t.String())
					foundTimeouts[t.Type] = end + 1
					section.SetContent(content)
				}
			}
		}

		var orderChanged bool
		content, orderChanged = reorderTimeouts(content, foundTimeouts)
		if orderChanged {
			errs = append(errs, fmt.Errorf("%s: Timeouts are not ordered as expected (CRUD)", IdAndName(r)))

			if fix {
				rd.Document.HasChange = true
				section.SetContent(content)
			}
		}
	}

	return errs
}

func parseTimeout(line string) *data.Timeout {
	matches := fullTimeoutRegex.FindStringSubmatch(line)

	if len(matches) == 4 {
		d, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil
		}

		return &data.Timeout{
			Type:     data.TimeoutType(matches[1]),
			Duration: d,
			Name:     strings.TrimSpace(matches[3]),
		}
	}

	return nil
}

func reorderTimeouts(c []string, found map[data.TimeoutType]int) ([]string, bool) {
	if len(found) <= 1 {
		return c, false
	}

	expectedOrder := []data.TimeoutType{
		data.TimeoutTypeCreate,
		data.TimeoutTypeRead,
		data.TimeoutTypeUpdate,
		data.TimeoutTypeDelete,
	}

	originalContent := make([]string, len(c))
	copy(originalContent, c)

	orderedTimeouts := make([]data.TimeoutType, 0, len(found))
	for _, t := range expectedOrder {
		if _, ok := found[t]; ok {
			orderedTimeouts = append(orderedTimeouts, t)
		}
	}

	orderedIndexes := util.MapValues2Slice(found)
	slices.Sort(orderedIndexes)

	for idx, v := range orderedTimeouts {
		origIdx := orderedIndexes[idx]
		c[origIdx] = originalContent[found[v]]
	}

	return c, !reflect.DeepEqual(originalContent, c)
}

func timeoutSliceToMap(t []data.Timeout) map[data.TimeoutType]data.Timeout {
	result := make(map[data.TimeoutType]data.Timeout)

	for _, v := range t {
		result[v.Type] = v
	}

	return result
}
