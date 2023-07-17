package logic

import (
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Run(resources ...string) {
	targets := map[string]*untypedResource{}
	for _, res := range resources {
		targets[res] = nil
	}

	allRes := loadAllUntypedResources()
	for _, res := range allRes {
		if _, ok := targets[res.name]; ok {
			targets[res.name] = res
		} else if _, ok = targets[strings.TrimPrefix(res.name, "azurerm_")]; ok {
			targets[res.name] = res
		}
	}

	for _, name := range resources {
		target := targets[name]
		if target == nil {
			log.Printf("no such resource: %s", name)
			continue
		}

		g, err := newGenerator(target)
		if err != nil {
			log.Printf("[Error] build generator for %s: %v", name, err)
		}
		meta := g.buildMeta()

		code, err := meta.codeGen()
		if err != nil {
			log.Printf("[Error] code gen %s err: %v", name, err)
			continue
		}

		if code2, err := format.Source(code); err != nil {
			log.Printf("[Error] format code: %v", err)
		} else {
			code = code2
		}

		fileName := filepath.Base(target.path)
		path := filepath.Join(filepath.Dir(target.path), fileName[:len(fileName)-len(filepath.Ext(fileName))]+"_typed.go")
		if err := os.WriteFile(path, code, 0644); err != nil {
			log.Printf("[Error] write file %s err: %v", name, err)
		}
		goimports(path)
	}
}
