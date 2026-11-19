package data

import (
	"gorm.io/gen"
)

func gopsTable(g *gen.Generator) {
	g.ApplyBasic(g.GenerateModelAs("settings", "Settings"))
}
