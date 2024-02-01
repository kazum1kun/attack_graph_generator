package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntSliceFlag{
				Name:     "node",
				Usage:    "number of PF, AND, and OR nodes",
				Aliases:  []string{"n"},
				Required: true,
				Action: func(ctx *cli.Context, input []int) error {
					haveNeg := false
					for _, v := range input {
						if v <= 0 {
							haveNeg = false
						}
					}
					if len(input) != 3 || haveNeg {
						return fmt.Errorf("flag node should contain exactly 3 positive integers")
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:     "edge",
				Usage:    "number of edges",
				Aliases:  []string{"e"},
				Required: true,
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
			},
			&cli.Int64Flag{
				Name:        "seed",
				Usage:       "random seed",
				Aliases:     []string{"s"},
				Value:       time.Now().UnixNano(),
				DefaultText: "current Unix epoch in seconds",
			},
			&cli.StringFlag{
				Name:        "outdir",
				Usage:       "output to `DIR`",
				Value:       ".",
				Aliases:     []string{"o"},
				DefaultText: "current directory",
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
				Usage:       "generate a pdf rendition of the attack graph",
				Aliases:     []string{"g"},
				Value:       false,
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "graph",
				Usage:       "generate a pdf rendition of the attack graph",
				Aliases:     []string{"g"},
				Value:       false,
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "nocheck",
				Usage:       "disable basic input validation checks (may result in invalid graphs)",
				Value:       false,
				DefaultText: "false",
			},
		},
		Action: func(ctx *cli.Context) error {
			// PF AND OR
			nodeNum := ctx.IntSlice("node")
			edgeNum := ctx.Int("edge")

			if !ctx.Bool("nocheck") {
				// Sanity check
				minEdge := (nodeNum[0] + 2*nodeNum[1] + nodeNum[2]) / 2
				maxEdge := (nodeNum[0] + nodeNum[2]) * nodeNum[1]
				if edgeNum < minEdge || edgeNum > maxEdge {
					return fmt.Errorf("flag edge out of bound for current input, valid range [%v-%v], current %v",
						minEdge, maxEdge, edgeNum)
				}
				if nodeNum[2] > nodeNum[1] {
					return fmt.Errorf("number of OR node cannot be greater than number of AND node")
				}
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
