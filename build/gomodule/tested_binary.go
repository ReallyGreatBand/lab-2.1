package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/ReallyGreatBand/lab-2.1/build/gomodule")

	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	goTest = pctx.StaticRule("binaryTest", blueprint.RuleParams{
		Command: "cd $workDir && go test -v $testPkg > $output",
		Description: "test go command $testPkg",
	}, "workDir", "testPkg", "output")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")
)

// goBinaryModuleType implements the simplest Go binary build without running tests for the target Go package.
type goBinaryModuleType struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg string
		// Go package name to test
		TestPkg string
		// Go test excludes
		TestSrcs []string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		// If to call vendor command.
		VendorFirst bool

		DefaultBuild bool
		// Example of how to specify dependencies.
		Deps []string
	}
}

func (gb *goBinaryModuleType) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return gb.properties.Deps
}

func (gb *goBinaryModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	output := path.Join(config.BaseOutputDir, name, "bood_test")

	var inputs []string
	var testInputs []string
	inputErrors := false
	for _, src := range gb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gb.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErrors = true
		}
	}
	for _, src := range gb.properties.TestSrcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			testInputs = append(testInputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErrors = true
		}
	}
	if inputErrors {
		return
	}

	if gb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	build := gb.properties.DefaultBuild

	if build {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Build %s as Go binary", name),
			Rule:        goBuild,
			Outputs:     []string{outputPath},
			Implicits:   inputs,
			Optional:    true,
			Args: map[string]string{
				"outputPath": outputPath,
				"workDir":    ctx.ModuleDir(),
				"pkg":        gb.properties.Pkg,
			},
		})
	}

	testInputs = append(testInputs, inputs...)
	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Test %s as Go binary", name),
		Rule:        goTest,
		Outputs:     []string{output},
		Implicits:   testInputs,
		Args: map[string]string{
			"output":     output,
			"workDir":    ctx.ModuleDir(),
			"testPkg":    gb.properties.TestPkg,
		},
	})
}

// SimpleBinFactory is a factory for go binary module type which supports Go command packages with running tests.
func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &goBinaryModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}