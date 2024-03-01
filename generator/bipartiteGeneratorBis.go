package generator

import (
	"fmt"
	dll "github.com/kazum1kun/attack_graph_generator/doublylinkedlist"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
	"math/rand"
)

// ConstructGraphAlt A memory-hungry version of the generator
func ConstructGraphAlt(leaf, and, or, edge int, cycleOk, relaxed bool, rnd *rand.Rand) *[]*CNode {
	total := or + leaf + and
	andPadding := 1 + leaf + or

	// Initialize CNodes for all nodes
	V := make([]*CNode, total+1)
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
		} else if i <= or {
			_type = OR
			iCap = and
			oCap = and
			desc = fmt.Sprintf("d%v", i)
		} else if i <= or+leaf {
			_type = LEAF
			iCap = 0
			oCap = and
			desc = fmt.Sprintf("p%v", i-or)
		} else {
			_type = AND
			iCap = leaf + or
			if relaxed {
				oCap = or
			} else {
				oCap = 1
			}
			desc = fmt.Sprintf("r%v", i-leaf-or)
		}
		V[i] = &CNode{
			Id:   i,
			Type: _type,
			Desc: desc,
			Pred: *set.New[int](),
			ICap: iCap,
			OCap: oCap,
			Adj:  *set.New[int](),
		}
		V[i].Pred.Add(i)
	}

	// Enumerate ALL possible combinations of edges... this will use a lot of aux space
	andToOr := dll.New[Edge]()
	orToAnd := dll.New[Edge]()

	orTargets := makeRange(1, or-1)
	andTargets := makeRange(andPadding, total)

	// Skip the goal node as it cannot be a source
	for i := 2; i <= total; i++ {
		if V[i].Type == AND {
			andToOr.Add(generateEdges(i, orTargets)...)
		} else {
			orToAnd.Add(generateEdges(i, andTargets)...)
		}
	}

	// Randomly pick edges from the universe of all edges
	// Here we use the ratio of AND vs OR + PF to determine the probability each is drawn
	andRatio := float64(and / (or + leaf))
	var target *dll.List[Edge]
	for edge > 0 {
		if rnd.Float64() > andRatio {
			target = orToAnd
		} else {
			target = andToOr
		}
	reroll:
		idx := rnd.Intn(target.Size())
		if !attemptAdd(target, idx, !relaxed, cycleOk, &V) {
			goto reroll
		} else {
			edge -= 1
		}
	}

	return &V
}

type Edge struct {
	Src int
	Dst int
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func generateEdges(src int, dst []int) []Edge {
	result := make([]Edge, len(dst))
	for i, v := range dst {
		result[i] = Edge{src, v}
	}
	return result
}

func attemptAdd(list *dll.List[Edge], target int, deleteAll, cycleOk bool, V *[]*CNode) bool {
	it := list.Iterator()
	var result Edge

	for ; it.Index() != target; it.Next() {
	}
	result = it.Value()

	if !addEdge((*V)[result.Src], (*V)[result.Dst], cycleOk, V) {
		return false
	}

	if deleteAll && (*V)[result.Src].Type == AND {
		var start, end int
		// seek to the starting block
		for ; it.Value().Src == start; it.Prev() {
		}
		it.Next()
		start = it.Index()
		for ; it.Value().Src == start; it.Next() {
		}
		it.Prev()
		end = it.Index()
		list.RangeRemove(start, end)
	} else {
		list.Remove(target)
	}

	return true
}
