package main

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"log"
	"os/exec"
	"path/filepath"
)

func csvToPdf(arcFile, vertFile string) {
	baseDir := filepath.Dir(arcFile)
	// Convert CSV files to DOT files
	cmd1 := fmt.Sprintf(".\\bin\\sed.exe -f .\\misc\\VERTICES.sed %s >> %s\\AttackGraph.dot", vertFile, baseDir)
	cmd2 := fmt.Sprintf(".\\bin\\sed.exe -f .\\misc\\ARCS.sed %s >> %s\\AttackGraph.dot", arcFile, baseDir)
	exec.Command("pwsh.exe", "-c", cmd1)
	exec.Command("pwsh.exe", "-c", cmd2)

	graph, err := graphviz.ParseFile(baseDir + "\\AttackGraph.dot")
	if err != nil {
		log.Fatalf("Error parsing DOT file: %v", err)
	}

	g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, fmt.Sprintf("%v\\AttackGraph.png", baseDir)); err != nil {
		log.Fatalf("Error rendering graph to file: %v", err)
	}
}
