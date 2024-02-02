package main

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func csvToPdf(arcFile, vertFile string) {
	baseDir := filepath.Dir(arcFile)

	// Convert CSV files to DOT files
	outFile, err := os.Create(baseDir + "\\AttackGraph.dot")
	if err != nil {
		log.Fatalf("Error creating DOT file: %v", err)
	}
	_, _ = outFile.WriteString("digraph G {")
	_ = outFile.Sync()
	_ = outFile.Close()

	c1 := fmt.Sprintf(".\\bin\\sed.exe -f .\\misc\\VERTICES_no_metric.sed %s >> %s\\AttackGraph.dot", vertFile, baseDir)
	c2 := fmt.Sprintf(".\\bin\\sed.exe -f .\\misc\\ARCS_noLabel.sed %s >> %s\\AttackGraph.dot", arcFile, baseDir)
	cmd1 := exec.Command("pwsh.exe", "-c", c1)
	cmd2 := exec.Command("pwsh.exe", "-c", c2)
	err = cmd1.Run()
	err = cmd2.Run()
	if err != nil {
		log.Fatalf("Error converting CSV to DOT file: %v", err)
	}

	outFile, err = os.OpenFile(baseDir+"\\AttackGraph.dot", os.O_APPEND, 0644)
	_, _ = outFile.WriteString("}")
	_ = outFile.Sync()
	_ = outFile.Close()

	graph, err := graphviz.ParseFile(baseDir + "\\AttackGraph.dot")
	if err != nil {
		log.Fatalf("Error parsing DOT file: %v", err)
	}

	g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, fmt.Sprintf("%v\\AttackGraph.png", baseDir)); err != nil {
		log.Fatalf("Error rendering graph to file: %v", err)
	}
}
