package gomodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

func TestSimpleBinFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			tested_binary {
			  name: "test-out",
			  srcs: ["test-src.go"],
              testPkg: "./...",
              testSrcs: ["**/*_test.go"],
			  pkg: ".",
	          vendorFirst: true
			}
		`),
		"test-src.go": nil,
		"test-src_test.go": nil,
	})

	ctx.RegisterModuleType("tested_binary", SimpleBinFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if !strings.Contains(text, "out/bin/test-out: ") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
		if !strings.Contains(text, " test-src_test.go") {
			t.Errorf("Generated ninja file does not have test file")
		}
		if !strings.Contains(text, " test-src.go") {
			t.Errorf("Generated ninja file does not have source dependency")
		}
		if !strings.Contains(text, "build vendor: g.gomodule.vendor | go.mod") {
			t.Errorf("Generated ninja file does not have vendor build rule")
		}
	}
}