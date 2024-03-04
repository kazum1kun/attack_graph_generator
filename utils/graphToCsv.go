package utils

import (
	"fmt"
	"github.com/kazum1kun/attack_graph_generator/hashset/generator"
	"log"
	"math/rand"
	"os"
)

func GraphToCsv(V *[]*generator.CNode, outDir string, rnd *rand.Rand) {
	// Assume the outDir is good as it was tested in the main
	vertFile, err := os.Create(fmt.Sprintf("%s/VERTICES.CSV", outDir))
	arcFile, err := os.Create(fmt.Sprintf("%s/ARCS.CSV", outDir))

	defer func(vertFile *os.File) {
		err := vertFile.Close()
		if err != nil {
			log.Println("Error closing vertex file!")
		}
	}(vertFile)
	defer func(arcFile *os.File) {
		err := arcFile.Close()
		if err != nil {
			log.Println("Error closing arc file!")
		}
	}(arcFile)

	// Initialize randNorm if

	if err != nil {
		log.Panicf("Error creating file in %s\n", outDir)
	}
	for _, node := range *V {
		if node == nil {
			continue
		}
		var endNum int
		if node.Type == generator.LEAF {
			endNum = 1
		} else {
			endNum = 0
		}
		_, err := vertFile.WriteString(fmt.Sprintf("%d,\"%s\",\"%s\",%d\n",
			node.Id, node.Desc, node.Type, endNum))
		if err != nil {
			log.Panicln("Error writing to VERTEX.CSV file")
		}

		edgeProb := -1.0
		for _, neighbor := range node.Adj.Values() {
			if rnd != nil {
				edgeProb = RandNorm(rnd)
			}
			_, err := arcFile.WriteString(fmt.Sprintf("%d,%d,%f\n", neighbor, node.Id, edgeProb))
			if err != nil {
				log.Panicln("Error writing to ARCS.CSV file")
			}
		}
	}

	err = vertFile.Sync()
	err = arcFile.Sync()
	if err != nil {
		log.Fatalln("Error flushing ARCS.CSV file")
	}

}
