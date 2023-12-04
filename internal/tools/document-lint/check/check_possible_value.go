// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type possibleValueDiff struct {
	checkBase

	Want []string
	Got  []string

	Missed []string // value not exists in doc
	Spare  []string // value redundant in doc, not exists in code
}

func newPossibleValueDiff(checkBase checkBase, want []string, got []string, missed []string, spare []string) *possibleValueDiff {
	return &possibleValueDiff{
		checkBase: checkBase,
		Want:      want,
		Got:       got,
		Missed:    missed,
		Spare:     spare,
	}
}

func possibleValueStr(values []string) string {
	var codes []string
	for _, val := range values {
		codes = append(codes, util.ItalicCode(val))
	}
	return fmt.Sprintf("[%s]", strings.Join(codes, ", "))
}

func (p possibleValueDiff) String() string {
	var missInDoc, missInCode string
	if len(p.Missed) > 0 {
		missInDoc = fmt.Sprintf(" the following possible values are missing in the documentation: %s.", possibleValueStr(p.Missed))
	}
	if len(p.Spare) > 0 {
		missInCode = fmt.Sprintf(" the following possible values are missing in the schema: %v.", possibleValueStr(p.Spare))
	}
	return fmt.Sprintf(`%s:%s%s`,
		p.checkBase.Str(),
		missInDoc,
		missInCode,
	)
}

func (p possibleValueDiff) Fix(line string) (result string, err error) {
	if len(p.Want) == 0 {
		return line, nil
	}
	result = line
	// replace from field.EnumStart to field.EnumEnd
	var bs strings.Builder
	if len(p.Got) == 0 || (p.MDField().EnumStart == 0 && len(p.Missed) > 0) {
		// skip this kind of field. may submit in a separate run
		// find default index
		idx := strings.Index(line, "Defaults to")
		if idx < 0 {
			idx = strings.Index(line, "Changing this forces")
		}
		if idx > 0 {
			bs.WriteString(line[:idx])
		} else {
			bs.WriteString(line)
			bs.WriteByte(' ')
		}
		if len(p.Want) == 1 {
			bs.WriteString("The only possible value is ")
		} else {
			bs.WriteString("Possible values are ")
		}
		bs.WriteString(patchWantEnums(p.Want))
		bs.WriteByte('.')
		if idx > 0 {
			bs.WriteByte(' ')
			bs.WriteString(line[idx:])
		}
		result = bs.String()
	} else if len(p.Missed) > 0 {
		// only replace missed values
		bs.WriteString(line[:p.MDField().EnumStart])
		if len(p.Want) == 1 {
			bs.WriteString("The only possible value is ")
		} else {
			bs.WriteString("Possible values are ")
		}
		bs.WriteString(patchWantEnums(p.Want))
		if end := p.MDField().EnumEnd; end > 0 && end < len(line) {
			bs.WriteString(line[p.MDField().EnumEnd:])
		} else {
			f := p.MDField()
			log.Printf("warning enum end %s:L%d len %dvs%d; %s", path.Base(f.Path), f.Line, f.EnumEnd, len(line), line)
		}
		result = bs.String()
	}
	return result, nil
}

var _ Checker = (*possibleValueDiff)(nil)

func patchWantEnums(want []string) string {
	res := make([]string, len(want))
	for idx, val := range want {
		res[idx] = "`" + val + "`"
	}
	if len(res) == 1 {
		return res[0]
	}

	s := res[0]
	if len(res) >= 3 {
		s = strings.Join(res[:len(res)-1], ", ")
	}
	s += " and " + res[len(res)-1]
	return s
}

// check possible values
func checkPossibleValues(r *schema.Resource, md *model.ResourceDoc) (res []Checker) {
	schemModel := r.Schema.Schema
	_ = schemModel
	if md == nil {
		log.Printf("%s no match document exists", r.ResourceType)
		return
	}
	// loop over document model
	for name, field := range md.Args {
		partRes := diffField(r, field, []string{name})
		res = append(res, partRes...)
	}
	return
}

// xPath property name for parent nodes
func diffField(r *schema.Resource, mdField *model.Field, xPath []string) (res []Checker) {
	fullPath := strings.Join(xPath, ".")
	if isSkipProp(r.ResourceType, fullPath) {
		return
	}

	// if end property
	if mdField.Subs == nil {
		want := r.PossibleValues[fullPath]
		docVal := mdField.PossibleValues()
		if missed, spare := SliceDiff(want, docVal, true); len(missed)+len(spare) > 0 {
			if !mayExistsInDoc(mdField.Content, want) {
				base := newCheckBase(mdField.Line, fullPath, mdField)
				res = append(res, newPossibleValueDiff(base, want, docVal, missed, spare))
			}
		}
		return
	}
	// check if r has such path
	if !r.HasPathFor(xPath) {
		log.Printf("%s %s has no path [%s], there must be an error in markdwon", color.YellowString("[WARN]"), r.ResourceType, strings.Join(xPath, "."))
		return
	}
	for _, sub := range mdField.Subs {
		subRes := diffField(r, sub, append(xPath, sub.Name))
		res = append(res, subRes...)
	}
	return
}

func SliceDiff(want, got []string, caseInSensitive bool) (missed, spare []string) {
	// if `want` is nil then it may only write in doc, skip this
	if len(want) == 0 {
		return
	}
	// cross-check
	wantCpy, gotCpy := want, got
	if caseInSensitive {
		wantCpy = make([]string, len(want))
		gotCpy = make([]string, len(got))
		for idx := range want {
			wantCpy[idx] = strings.ToLower(want[idx])
		}
		for idx := range got {
			gotCpy[idx] = strings.ToLower(got[idx])
		}
	}
	wantMap := util.Slice2Map(wantCpy)
	gotMap := util.Slice2Map(gotCpy)

	for idx, k := range wantCpy {
		if _, ok := gotMap[k]; !ok {
			missed = append(missed, want[idx])
		}
	}

	for idx, k := range gotCpy {
		if _, ok := wantMap[k]; !ok {
			spare = append(spare, got[idx])
		}
	}

	return
}

// return true values exists in doc but may not with code quote
func mayExistsInDoc(docLine string, want []string) bool {
	for _, val := range want {
		if !strings.Contains(docLine, val) {
			return false
		}
	}
	return true
}
