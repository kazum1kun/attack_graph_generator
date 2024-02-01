package main

import (
	"fmt"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
	"math/rand"
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

func constructGraph(leaf, and, or, edge int, cycle, noCheck bool, seed int64) {
	rnd := rand.New(rand.NewSource(seed))

	// Initialize CNodes for all nodes
	// Initialize CNodes for all nodes
	V := make([]*CNode, leaf+and+or+1)
	for i := 1; i <= leaf+and+or; i++ {
		var _type NodeType
		var iCap, oCap int
		var desc string
		if i == 1 {
			// goal node
			_type = OR
			iCap = and
			oCap = 0
			desc = "goal"
		} else if i <= leaf {
			_type = LEAF
			iCap = 0
			oCap = and
			desc = fmt.Sprintf("p%v", i-1)
		} else if i <= leaf+or {
			_type = OR
			iCap = and
			oCap = (or - 1) * and
			desc = fmt.Sprintf("d%v", i-leaf-1)
		} else {
			_type = AND
			iCap = leaf + or
			if noCheck {
				oCap = or
			} else {
				oCap = 1
			}
			desc = fmt.Sprintf("r%v", i-leaf-and-1)
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
		V[i].Pred.Add(i)
	}

	// Stage 1: priority matches
	// Nodes are connected such that a minimal, valid attack graph is generated that uses all the nodes
	// To do so, we create a random permutation of numbers for PF, AND, and OR nodes
	// Then we simply match them from top to bottom
	// First match inbound edges from PF/OR to AND
	priorityInSrc := rnd.Perm(leaf + or - 1) // skip goal node
	priorityInDst := rnd.Perm(and)

	// Preconditions guarantee that the number of AND is geq than then number of OR
	// However this may not hold in case the check is turned off
	// Match priority src to priority dst
	minIn := min(and, leaf+or-1)
	for i := 0; i < minIn; i++ {
		// add an edge from src[i] to dst[i]
		srcId := priorityInSrc[i] + 2
		dstId := priorityInDst[i] + 1 + leaf + or
		addEdge(V[srcId], V[dstId], cycle)
		edge -= 1
	}
	// Match the remaining
	if minIn == and {
		for i := len(priorityInDst); i < len(priorityInSrc); i++ {
			srcId := priorityInSrc[i] + 2
			dstId := rnd.Intn(and) + 1 + leaf + or
			addEdge(V[srcId], V[dstId], cycle)
			edge -= 1
		}
	} else {
		for i := len(priorityInSrc); i < len(priorityInDst); i++ {
			srcId := rnd.Intn(leaf+or-1) + 2
			dstId := priorityInDst[i] + 1 + leaf + or
			addEdge(V[srcId], V[dstId], cycle)
			edge -= 1
		}
	}

	// Now match outbound edges from AND to OR
	priorityOutSrc := rnd.Perm(and)
	priorityOutDst := rnd.Perm(or)

	minOut := min(and, or)
	for i := 0; i < minOut; i++ {
		srcId := priorityOutSrc[i] + 1 + leaf + or
		dstId := priorityOutDst[i] + 1
		addEdge(V[srcId], V[dstId], cycle)
		edge -= 1
	}
	// Match the remaining
	if minOut == and {
		for i := len(priorityOutSrc); i < len(priorityOutDst); i++ {
			srcId := rnd.Intn(and) + leaf + or + 1
			dstId := priorityOutDst[i] + 1
			addEdge(V[srcId], V[dstId], cycle)
			edge -= 1
		}
	} else {
		for i := len(priorityOutDst); i < len(priorityOutSrc); i++ {
			srcId := priorityOutSrc[i] + leaf + or + 1
			dstId := rnd.Intn(or) + 1
			addEdge(V[srcId], V[dstId], cycle)
			edge -= 1
		}
	}

	// Stage 2: match the remaining edge quota randomly
	for ; edge > 0; edge-- {

	}

}

func addEdge(src, dst *CNode, cycle bool) bool {
	// check for cycles (a cycle is found when the dst's predecessors are a subset of src's predecessor)
	if !cycle && src.Pred.Contains(dst.Pred.Values()...) {
		return false
	}
	src.Adj.Add(dst.Id)
	dst.Pred.Union(&src.Pred)
	src.OCap -= 1
	dst.ICap -= 1
	return true
}
