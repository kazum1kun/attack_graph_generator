package main

import (
	"fmt"
	"github.com/kazum1kun/attack_graph_generator/hashset/generator"
	"github.com/kazum1kun/attack_graph_generator/hashset/utils"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	app := &cli.App{
		Name:      "AGG",
		Usage:     "Mulval-compatible attack graph generator",
		UsageText: "agg command [command options]",
		Version:   "0.2.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "number of OR, PF and AND nodes (in this order)",
				Aliases:  []string{"n"},
				Required: true,
				Category: "GRAPH",
			},
			&cli.IntFlag{
				Name:     "edge",
				Usage:    "number of edges",
				Aliases:  []string{"e"},
				Required: true,
				Category: "GRAPH",
				Action: func(ctx *cli.Context, i int) error {
					if i <= 0 {
						return fmt.Errorf("flag edge must be positive")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:        "cycle",
				Usage:       "whether cycles are permitted (does not guarantee them)",
				Aliases:     []string{"c"},
				Value:       false,
				DefaultText: "false",
				Category:    "GENERATION",
			},
			&cli.Int64Flag{
				Name:        "seed",
				Usage:       "random seed",
				Aliases:     []string{"s"},
				Value:       time.Now().UnixNano(),
				DefaultText: "current Unix epoch in nanoseconds",
				Category:    "GENERATION",
			},
			&cli.StringFlag{
				Name:        "outdir",
				Usage:       "output to `DIR`",
				Value:       ".",
				Aliases:     []string{"o"},
				DefaultText: "current directory",
				Category:    "OUTPUT",
				Action: func(ctx *cli.Context, s string) error {
					_, err := os.Stat("s")
					if os.IsNotExist(err) {
						return fmt.Errorf("%v does not exist", s)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:        "graph",
				Usage:       "generate a graphical rendition of the attack graph",
				Aliases:     []string{"g"},
				Value:       false,
				DefaultText: "false",
				Category:    "OUTPUT",
			},
			&cli.BoolFlag{
				Name:        "relaxed",
				Usage:       "relax the constraint so that AND node can have multiple outgoing edges",
				Value:       false,
				DefaultText: "false",
				Category:    "GENERATION",
			},
			&cli.BoolFlag{
				Name:        "altgen",
				Usage:       "alternative generation method that used a lot more RAM, but converge quicker",
				Value:       false,
				Aliases:     []string{"alt"},
				DefaultText: "false",
				Category:    "GENERATION",
			},
			&cli.StringFlag{
				Name:        "vertsed",
				Usage:       "sed `FILE` to be used to process VERTICES.CSV",
				Value:       "./misc/VERTICES_no_metric.sed",
				Aliases:     []string{"vs"},
				DefaultText: "./misc/VERTICES_no_metric.sed",
				Category:    "OUTPUT",
			},
			&cli.StringFlag{
				Name:        "arcsed",
				Usage:       "sed `FILE` to be used to process ARCS.CSV",
				Value:       "./misc/ARCS_noLabel.sed",
				Aliases:     []string{"as"},
				DefaultText: "./misc/ARCS_noLabel.sed",
				Category:    "OUTPUT",
			},
			&cli.BoolFlag{
				Name:        "randw",
				Usage:       "add random weights to the edges",
				Value:       false,
				Aliases:     []string{"rw"},
				DefaultText: "false",
				Category:    "OUTPUT",
			},
		},
		Action: func(ctx *cli.Context) error {
			node := ctx.String("node")
			tokens := strings.Split(node, " ")
			hasError := false
			if len(tokens) != 3 {
				hasError = true
			}
			or, err := strconv.Atoi(tokens[0])
			if err != nil {
				hasError = true
			}
			leaf, err := strconv.Atoi(tokens[1])
			if err != nil {
				hasError = true
			}
			and, err := strconv.Atoi(tokens[2])
			if err != nil {
				hasError = true
			}
			if min(min(or, leaf), and) <= 0 {
				hasError = true
			}
			if hasError {
				return fmt.Errorf("flag node should contain exactly 3 positive integers")
			}

			edgeNum := ctx.Int("edge")

			minEdge := (leaf + 2*and + or) / 2
			var maxEdge int
			if !ctx.Bool("relaxed") {
				maxEdge = leaf*and + or*and
				if or > and {
					return fmt.Errorf("number of OR node cannot be greater than number of AND node")
				}
			} else {
				maxEdge = (2*leaf + or + 4) * and / 2
			}
			if edgeNum < minEdge || edgeNum > maxEdge {
				return fmt.Errorf("flag edge out of bound for current input, valid range [%v-%v], current %v",
					minEdge, maxEdge, edgeNum)
			}

			var cycleOk bool
			if ctx.Bool("cycle") {
				cycleOk = true
			} else {
				cycleOk = false
			}
			seed := ctx.Int64("seed")
			outDir := ctx.String("outdir")
			generateGraph := ctx.Bool("graph")
			relaxed := ctx.Bool("relaxed")
			randW := ctx.Bool("randw")
			alt := ctx.Bool("altgen")

			rnd := rand.New(rand.NewSource(seed))

			var V *[]*generator.CNode
			if alt {
				V = generator.ConstructGraphAlt(leaf, and, or, edgeNum, cycleOk, relaxed, rnd)
			} else {
				V = generator.ConstructGraph(leaf, and, or, edgeNum, cycleOk, relaxed, rnd)
			}

			if randW {
				utils.SetGaussianParams(math.Log2(float64(edgeNum)), 0.2, 0.01)
			}

			arcSed := ctx.String("arcsed")
			vertSed := ctx.String("vertsed")

			if randW {
				utils.GraphToCsv(V, outDir, rnd)
			} else {
				utils.GraphToCsv(V, outDir, nil)
			}

			if generateGraph {
				utils.CsvToPdf(fmt.Sprintf("%s/ARCS.CSV", outDir), fmt.Sprintf("%s/VERTICES.CSV", outDir),
					arcSed, vertSed)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
