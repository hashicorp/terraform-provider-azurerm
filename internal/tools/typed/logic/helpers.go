package logic

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func camelCase(name string) string {
	bs := strings.Builder{}
	for idx, c := range name {
		if c == '_' {
			continue
		}
		if bs.Len() == 0 || name[idx-1] == '_' {
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
		}
		bs.WriteByte(byte(c))
	}
	if bs.Len() == 0 {
		return "_"
	}
	return bs.String()
}

func modelName(key string) string {
	return camelCase(key) + "Model"
}

func valueType(s *pluginsdk.Schema, name string) string {
	t := s.Type
	switch t {
	case pluginsdk.TypeBool:
		return "bool"
	case pluginsdk.TypeInt:
		return "int"
	case pluginsdk.TypeString:
		return "string"
	case pluginsdk.TypeFloat:
		return "float64"
	case pluginsdk.TypeList:
		switch el := s.Elem.(type) {
		case *pluginsdk.Schema:
			return fmt.Sprintf("[]%s", valueType(el, name))
		case *pluginsdk.Resource:
			return fmt.Sprintf("[]%s", modelName(name))
		}
	case pluginsdk.TypeSet:
		switch el := s.Elem.(type) {
		case *pluginsdk.Schema:
			return fmt.Sprintf("[]%s", valueType(el, name))
		case *pluginsdk.Resource:
			return fmt.Sprintf("[]%s", modelName(name))
		}
	case pluginsdk.TypeMap:
		if name == "tag" {
			return "map[string]string"
		}
		return "map[string]string"
	}
	return "undef"
}

func StringContainsAny(name string, needle ...string) bool {
	for _, n := range needle {
		if strings.Contains(name, n) {
			return true
		}
	}
	return false
}

func safeRun(fn func()) {
	defer func() {
		_ = recover()
	}()
	fn()
}

func goimports(file string) {
	var out, stdErr bytes.Buffer
	cmd := exec.Command("goimports", "-w", file)
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	if err := cmd.Run(); err != nil {
		log.Printf("[Error] run goimports: %v, stdout: %s, stderr: %s", err, out.String(), stdErr.String())
	}
}
