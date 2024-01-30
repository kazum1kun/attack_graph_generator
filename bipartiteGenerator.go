package main

import (
	"fmt"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
)

type NodeType string

const (
	LEAF NodeType = "LEAF"
	AND           = "AND"
	OR            = "OR"
)

type CNode struct {
	Id   int
	Type NodeType
	Desc string
	Pred set.Set[int]
	ICap int
	OCap int
	Adj  set.Set[int]
}

func constructGraph(leaf, and, or, edge int, cycle, noCheck bool) {
	// Initialize CNodes for all nodes
	V := make([]*CNode, leaf+and+or+1)
	for i := 1; i <= leaf+and+or; i++ {
		var _type NodeType
		var iCap, oCap int
		var desc string
		if i <= leaf {
			_type = LEAF
			iCap = 0
			oCap = and
			desc = fmt.Sprintf("p%v", i)
		} else if i <= leaf+and {
			_type = AND
			iCap = leaf + or
			if noCheck {
				oCap = or
			} else {
				oCap = 1
			}
			desc = fmt.Sprintf("r%v", i-leaf)
		} else if i <= leaf+and+or-1 {
			_type = OR
			iCap = and
			oCap = (or - 1) * and
			desc = fmt.Sprintf("d%v", i-leaf-and)
		} else {
			// goal node
			_type = OR
			iCap = and
			oCap = 0
			desc = "goal"
		}
		V[i] = &CNode{
			Id:   i,
			Type: _type,
			Desc: desc,
			Pred: set.Set[int]{},
			ICap: iCap,
			OCap: oCap,
			Adj:  set.Set[int]{},
		}
	}

}
