package main

import (
	"fmt"
	"log"
	"os"
)

func graphToCsv(V *[]*CNode, outDir string) {
	// Assume the outDir is good as it was tested in the main
	vertFile, err := os.Create(fmt.Sprintf("%s/VERTICES.CSV", outDir))
	arcFile, err := os.Create(fmt.Sprintf("%s/ARCS.CSV", outDir))

	defer vertFile.Close()
	defer arcFile.Close()
	if err != nil {
		log.Panicf("Error creating file in %s\n", outDir)
	}
	for _, node := range *V {
		var endNum int
		if node.Type == LEAF {
			endNum = 1
		} else {
			endNum = 0
		}
		_, err := vertFile.WriteString(fmt.Sprintf("%d,\"%s\",\"%s\",%d",
			node.Id, node.Desc, node.Type, endNum))
		if err != nil {
			log.Panicln("Error writing to VERTEX.CSV file")
		}

		for _, neighbor := range node.Adj.Values() {
			_, err := arcFile.WriteString(fmt.Sprintf("%d,%d,-1", neighbor, node.Id))
			if err != nil {
				log.Panicln("Error writing to ARCS.CSV file")
			}
		}
	}
	vertFile.Sync()
	arcFile.Sync()
}
