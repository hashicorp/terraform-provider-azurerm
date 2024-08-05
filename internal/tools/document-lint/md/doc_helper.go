// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package md

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	docRDir             string
	resourceFilePathMap map[string]string
	file2Reosurce       = map[string]string{}
	once                sync.Once
)

// MDPathFor return full path of markdown file of resource
func MDPathFor(resourceType string) string {
	// find source
	fullPath := path.Join(ResourceDir(), fmt.Sprintf("%s.html.markdown", strings.TrimPrefix(resourceType, "azurerm_")))
	// check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return getMappingPath(resourceType)
	}
	return fullPath
}

func getMappingPath(resourceName string) (res string) {
	if resourceFilePathMap == nil {
		once.Do(func() {
			tmpMap := map[string]string{}

			dir, err := os.ReadDir(ResourceDir())
			_ = err
			for _, en := range dir {
				if en.IsDir() {
					continue
				}
				fullPath := path.Join(ResourceDir(), en.Name())
				name := fileResource(fullPath)
				tmpMap[name] = fullPath
				if _, ok := file2Reosurce[fullPath]; !ok {
					file2Reosurce[fullPath] = name
				}
			}
			resourceFilePathMap = tmpMap
		})
	}
	return resourceFilePathMap[resourceName]
}

var titleReg = regexp.MustCompile(`\npage_title:[^\n]*(azurerm_[a-zA-Z0-9_]+)"?`)

func fileResource(path string) string {
	// match content
	content, _ := os.ReadFile(path)
	// if content match pattern
	if subs := titleReg.FindStringSubmatch(string(content)); len(subs) > 1 {
		return subs[1]
	}
	return ""
}

func docDir() string {
	file, _ := util.FuncFileLine(utils.Int32)
	return path.Join(path.Dir(path.Dir(file)), "website", "docs")
}

func ResourceDir() string {
	if docRDir == "" {
		docRDir = path.Join(docDir(), "r")
	}
	return docRDir
}
