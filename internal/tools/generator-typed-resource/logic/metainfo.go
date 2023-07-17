package logic

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"strings"
	"time"
)

type ModelField struct {
	Name string
	Type string
	Tag  string
}

type Model struct {
	Name   string
	Fields []ModelField
}

type SchemaField struct {
	Name    string
	Content string
}

func NewSchemaFieldFromAST(name string, content string) SchemaField {
	return SchemaField{
		Name:    name,
		Content: content,
	}
}

type Function struct {
	Name        string
	varModel    string
	Timout      time.Duration
	Ast         *ast.FuncDecl
	BodyContent string
	transformed bool // has been transformed
}

func (f Function) VarModel() string {
	if f.varModel == "" {
		return "model"
	}
	return f.varModel
}

func (f Function) Signature() string {
	return strings.TrimPrefix(types.ExprString(f.Ast.Type), "func")
}

type StateUpdate struct {
	SchemaVersion string
}

type MetaInfo struct {
	Name          string
	Package       string
	ModelName     string // root model name
	Recv          string
	IDValidator   string
	ResourceType  string
	SchemaVersion string
	StateUpgrade  string // upgrade string
	HasUpdate     bool
	Imports       []*ast.ImportSpec
	Models        []Model
	Arguments     []SchemaField
	Attributes    []SchemaField
	Functions     []*Function
	CreateFunc    *Function
	UpdateFunc    *Function
	ReadFunc      *Function
	DeleteFunc    *Function
}

func (m *MetaInfo) ImportsExpr() (res []string) {
	for _, im := range m.Imports {
		res = append(res, strings.TrimSpace(fmt.Sprintf("%s %s %s", identExpr(im.Name), im.Path.Value, im.Comment.Text())))
	}
	return
}

func (f Function) GenTimeout() string {
	if f.Timout == 0 {
		return "undef"
	}

	if f.Timout%time.Hour == 0 {
		return fmt.Sprintf("%d * time.Hour", f.Timout/time.Hour)
	}

	if f.Timout%time.Minute == 0 {
		return fmt.Sprintf("%d * time.Minute", f.Timout/time.Minute)
	}

	return fmt.Sprintf("%d * time.Second", f.Timout)
}

func (m *MetaInfo) findFieldByTag(modelName, tag string) *ModelField {
	if tag == "" {
		return nil
	}
	if modelName == "" {
		modelName = m.ModelName
	}
	keys := strings.SplitN(tag, ".", 2)
	for _, model := range m.Models {
		if model.Name == modelName {
			for _, f := range model.Fields {
				if f.Tag == keys[0] {
					if len(keys) >= 2 {
						return m.findFieldByTag(f.Type, keys[1])
					}
					return &f
				}
			}
		}
	}
	return nil
}

func (m *MetaInfo) topModel() *Model {
	for _, model := range m.Models {
		if model.Name == m.ModelName {
			return &model
		}
	}
	return nil
}

func (m *MetaInfo) topField(name string) *ModelField {
	top := m.topModel()
	if top == nil {
		log.Printf("[Error] no top model for %s", m.ModelName)
		return nil
	}
	for _, f := range top.Fields {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

func (m *MetaInfo) findModelByKeys(keys []string) *Model {
	for _, model := range m.Models {
		mkeys := make(map[string]bool, len(model.Fields))
		for _, f := range model.Fields {
			mkeys[f.Tag] = true
		}
		count := len(keys)
		for _, k := range keys {
			if mkeys[k] {
				count--
			}
		}
		if count == 0 {
			return &model
		}
	}
	return nil
}
