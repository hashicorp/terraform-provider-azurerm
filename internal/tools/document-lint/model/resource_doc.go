// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type FieldType int

const (
	FieldTypeAttr FieldType = iota
	FieldTypeBlock
)

type Properties map[string]*Field

// Merge p1 from attribute
func (p Properties) Merge(p1 Properties) {
	if p == nil {
		return
	}
	for k, v := range p1 {
		// if exists already, it maybe should be in attr
		if f, ok := p[k]; ok {
			f.SameNameAttr = v
		} else {
			p[k] = v
		}
	}
}

func (p Properties) AddField(f *Field) {
	if p == nil {
		log.Printf("try add field `%s` @L%d to nil property", f.Name, f.Line)
		return
	}
	p[f.Name] = f
}

func (p Properties) hasCircularRef() string {
	for name, f := range p {
		if f.Typ == FieldTypeBlock {
			if f.Subs.hasCircularRefFor(name) {
				return name
			}
		}
	}
	return ""
}

func (p Properties) hasCircularRefFor(name string) bool {
	if p == nil {
		return false
	}

	if f, ok := p[name]; ok {
		if f.Typ == FieldTypeBlock {
			return true
		}
		return f.Subs.hasCircularRefFor(name)
	}
	return false
}

const (
	Default  = 0
	Optional = 1 << iota
	Required
	Computed // if in attribute area, should not use
)

type Field struct {
	Name      string
	Path      string              // like xpath a.b.c
	Required  int                 `json:",omitempty"`
	Line      int                 `json:",omitempty"`
	Typ       FieldType           `json:",omitempty"` // is a simple attr field or a TypeList property
	Pos       PosType             // in Argument or in Attribute
	Default   string              `json:",omitempty"` // default value as string, empty for no default value
	ForceNew  bool                `json:",omitempty"` // if contains force new string
	Content   string              `json:",omitempty"` // origin doc line
	EnumStart int                 `json:",omitempty"`
	EnumEnd   int                 `json:",omitempty"`
	Enums     map[string]struct{} `json:",omitempty"`
	Subs      Properties          `json:",omitempty"`
	Skip      bool                `json:",omitempty"` // if skip this field
	FormatErr string              // the field has a format err causes parse error

	BlockTypeName string `json:",omitempty"` // block type name may not equal block name
	SameNameAttr  *Field // same name field exists in both arguments and attributes

	enumsInOrder []string // keep enums in order
	GuessEnums   []string // guess all code block as possible values
}

func NewField(line string) *Field {
	return &Field{
		Content: line,
		// Required: false,
	}
}

func (f *Field) SetGuessEnums(values []string) {
	// remove repeat values and ` mark
	hys := make(map[string]struct{}, len(values))
	var res []string
	for _, val := range values {
		val = strings.Trim(val, "`\"'")
		if _, ok := hys[val]; !ok {
			hys[val] = struct{}{}
			res = append(res, val)
		}
		hys[val] = struct{}{}
	}
	f.GuessEnums = res
}

func (f *Field) AddEnum(val ...string) {
	if f.Enums == nil {
		f.Enums = map[string]struct{}{}
	}
	for _, v := range val {
		if _, ok := f.Enums[v]; !ok {
			f.Enums[v] = struct{}{}
			f.enumsInOrder = append(f.enumsInOrder, v)
		}
	}
}

func (f *Field) PossibleValues() (res []string) {
	if len(f.enumsInOrder) > 0 {
		return f.enumsInOrder
	}
	// try fetch all code block as code
	return f.GuessEnums
}

func (f *Field) AddSubField(sub *Field) {
	if f.Subs == nil {
		f.Subs = map[string]*Field{}
	}
	f.Subs[sub.Name] = sub
}

func (f *Field) MatchName(name string) *Field {
	if f.Name == name {
		return f
	}
	for _, f1 := range f.Subs {
		if res := f1.MatchName(name); res != nil {
			return res
		}
	}
	return nil
}

func (f *Field) AllSubBlock(name string, needBlock bool) (res []*Field) {
	// name **or** guessed block type name equals name
	if f.Typ == FieldTypeBlock && f.BlockTypeName == name {
		res = append(res, f)
		return
	}
	if !needBlock {
		// if not need block type
		if f.BlockTypeName == "" && f.Name == name {
			res = append(res, f)
			return
		}
	}
	for _, f1 := range f.Subs {
		res = append(res, f1.AllSubBlock(name, needBlock)...)
	}
	return
}

type Timeout struct {
	Line  int   // line number, if line == 0 means no such timeout line in document
	Value int64 // timeout value in seconds
}

// Timeouts parse timeout value as second
type Timeouts struct {
	Create Timeout // nil if missed in document
	Update Timeout
	Read   Timeout
	Delete Timeout
}

type Import struct {
	ResourceType string
	ResourceID   string
}

func (p Properties) FindField(name string) *Field {
	for _, f := range p {
		if res := f.MatchName(name); res != nil {
			return res
		}
	}
	return nil
}

func (p Properties) FindAllSubBlock(name string) (res []*Field) {
	for _, f := range p {
		res = append(res, f.AllSubBlock(name, true)...)
	}
	if len(res) == 0 {
		// find in not block property, if no such block property found
		for _, f := range p {
			res = append(res, f.AllSubBlock(name, false)...)
		}
	}
	return res
}

type PossibleValue struct {
	Valeus []string
	Field  *Field
}

func NewPossibleValue(values []string, f *Field) PossibleValue {
	return PossibleValue{
		Valeus: values,
		Field:  f,
	}
}

type ResourceDoc struct {
	ResourceName string
	Args         Properties
	Attr         Properties
	ExampleHCL   string
	Timeouts     *Timeouts // nil if no timeouts part in document
	Import       Import

	Blocks map[string]Properties // two pass get all blocks

	PossibleValues map[string]PossibleValue // save a.b.c possible values here
}

func (r *ResourceDoc) SetTimeout(lineNum int, line string) {
	if r.Timeouts == nil {
		r.Timeouts = &Timeouts{}
	}
	start, end := util.TimeoutValueIdx(line)
	if end <= start {
		return
	}

	text := line[start:end]
	values := strings.Split(strings.TrimSpace(text), " ")
	if len(values) < 2 {
		log.Printf("timeout value invalid: %s", line)
		return
	}

	// text is xxx hours or xxx minutes
	num, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		log.Printf("timeout number parse for %s err: %v", line, err)
		return
	}
	if strings.HasPrefix(values[1], "hour") {
		num *= 60 // all minutes
	}
	num *= 60 // set to seconds
	item := Timeout{
		Line:  lineNum,
		Value: num,
	}
	switch {
	case strings.Contains(line, "`create`"):
		r.Timeouts.Create = item
	case strings.Contains(line, "`update`"):
		r.Timeouts.Update = item
	case strings.Contains(line, "`read`"):
		r.Timeouts.Read = item
	case strings.Contains(line, "`delete`"):
		r.Timeouts.Delete = item
	default:
		log.Printf("line has not timeout type: %s", line)
	}
}

func (r *ResourceDoc) AllProp() Properties {
	res := Properties{}
	for k, v := range r.Args {
		res[k] = v
	}
	for k, v := range r.Attr {
		// merge back to args if exists in both arg and attr
		if arg, ok := res[k]; ok {
			if arg.Subs != nil {
				for k2, v2 := range v.Subs {
					if _, ok := arg.Subs[k2]; ok {
						arg.Subs[k2] = v2
					}
				}
			}
		} else {
			res[k] = v
		}
	}
	return res
}

func (r *ResourceDoc) CurProp(pos PosType) Properties {
	if pos == PosArgs {
		return r.Args
	} else if pos == PosAttr {
		return r.Attr
	}
	return nil
}

// TuneSubBlocks loop blocks to link empty block properties to block by name
func (r *ResourceDoc) TuneSubBlocks() (fixNames []string) {
	var partial func(f *Field)
	partial = func(f *Field) {
		if f.Typ == FieldTypeBlock {
			if f.Subs == nil {
				fixNames = append(fixNames, f.BlockTypeName)
				f.Subs = r.Blocks[f.BlockTypeName]
			}
			for _, f2 := range f.Subs {
				partial(f2)
			}
		}
	}

	for _, f := range r.Args {
		partial(f)
	}
	for _, f := range r.Attr {
		partial(f)
	}
	return
}

func (r *ResourceDoc) HasCircularRef() string {
	if name := r.Args.hasCircularRef(); name != "" {
		return name
	}
	if name := r.Attr.hasCircularRef(); name != "" {
		return name
	}
	return ""
}

func NewResourceDoc() *ResourceDoc {
	return &ResourceDoc{
		Args:           Properties{},
		Attr:           Properties{},
		PossibleValues: map[string]PossibleValue{},
		Blocks:         map[string]Properties{},
	}
}
