package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	archiveBin = pctx.StaticRule("archiveBin", blueprint.RuleParams{
		Command:     "cd $workDir && zip $outputPath -j $toArchive",
		Description: "Archive $toArchive binary",
	}, "workDir", "toArchive", "outputPath")
)

type archiveModuleType struct {
	blueprint.SimpleName

	properties struct {
		ToArchive string
	}
}


func (am *archiveModuleType)GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)
	outputPath := path.Join(config.BaseOutputDir, "archive","bood_archive")

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("go binary archivation by module %s", name),
		Rule:        archiveBin,
		Outputs:     []string{outputPath},
		Inputs:      []string{path.Join(ctx.ModuleDir(), "out", "bin", am.properties.ToArchive)},
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"toArchive":  path.Join("out", "bin", am.properties.ToArchive),
		},
	})
}

func SimpleArchiveFactory() (blueprint.Module, []interface{}) {
	mType := &archiveModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
