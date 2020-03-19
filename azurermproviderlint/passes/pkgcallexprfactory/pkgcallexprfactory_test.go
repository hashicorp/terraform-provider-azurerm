package pkgcallexprfactory

import (
	"log"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestValidateAnalyzer(t *testing.T) {
	fooanalyzer := BuildAnalyzer("foo", "bar")
	err := analysis.Validate([]*analysis.Analyzer{fooanalyzer})
	if err != nil {
		log.Fatal(err)
	}
}
