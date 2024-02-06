package main

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/rwtodd/Go.Sed/sed"
	"log"
	"os"
	"path/filepath"
)

func csvToPdf(arcFile, vertFile string) {
	baseDir := filepath.Dir(arcFile)

	// Convert CSV files to DOT files
	// Create output DOT file
	outFile, err := os.Create(baseDir + "/AttackGraph.dot")
	if err != nil {
		log.Fatalf("Error creating DOT file: %v", err)
	}
	_, _ = outFile.WriteString("digraph G {\n")

	// Convert vertices
	vertSedFile, err := os.Open("./misc/VERTICES_no_metric.sed")
	if err != nil {
		log.Fatalln("Error opening vertices sed definition files")
	}
	engine, err := sed.New(vertSedFile)
	if err != nil {
		log.Fatalln("Error parsing provided vertices sed files")
	}
	vertString, _ := os.ReadFile(vertFile)
	runString, err := engine.RunString(string(vertString))
	if err != nil {
		log.Fatalln("Error parsing CSV to .dot files")
	}
	_, _ = outFile.WriteString(runString)
	if err != nil {
		log.Fatalf("Error converting CSV to DOT file: %v", err)
	}

	// Convert edges
	arcSedFile, err := os.Open("./misc/ARCS_noLabel.sed")
	if err != nil {
		log.Fatalln("Error opening arcs sed definition files")
	}
	engine, err = sed.New(arcSedFile)
	if err != nil {
		log.Fatalln("Error parsing provided arcs sed files")
	}
	arcString, _ := os.ReadFile(arcFile)
	runString, err = engine.RunString(string(arcString))
	if err != nil {
		log.Fatalln("Error parsing CSV to .dot files")
	}
	_, _ = outFile.WriteString(runString)

	_, _ = outFile.WriteString("}")
	_ = outFile.Sync()
	_ = outFile.Close()

	graph, err := graphviz.ParseFile(baseDir + "/AttackGraph.dot")
	if err != nil {
		log.Fatalf("Error parng DOT file: %v", err)
	}

	g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, fmt.Sprintf("%v/AttackGraph.png", baseDir)); err != nil {
		log.Fatalf("Error rendering graph to file: %v", err)
	}
}
