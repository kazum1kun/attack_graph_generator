package generator

import (
	"fmt"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
	hset "github.com/ugurcsen/gods-generic/sets/linkedhashset"
	"math/rand"
)
import sll "github.com/ugurcsen/gods-generic/lists/singlylinkedlist"

// ConstructGraphAlt A memory-hungry version of the generator
func ConstructGraphAlt(leaf, and, or, edge int, cycleOk, relaxed bool, rnd *rand.Rand) *[]*CNode {
	// We first construct 4 lists:
	// inRequiredOr is the list of OR that needs an inbound edge
	inRequiredOr := sll.New[int]()
	// inRequiredAnd is the list of AND that needs an inbound edge
	inRequiredAnd := sll.New[int]()
	// outRequiredOr is the list of OR/LEAF that needs an outbound edge
	outRequiredOr := sll.New[int]()
	// outRequiredAnd is the list of AND that needs an outbound edge
	outRequiredAnd := sll.New[int]()

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
			inRequiredOr.Add(i)
		} else if i <= or {
			_type = OR
			iCap = and
			oCap = and
			desc = fmt.Sprintf("d%v", i)
			inRequiredOr.Add(i)
			outRequiredOr.Add(i)
		} else if i <= or+leaf {
			_type = LEAF
			iCap = 0
			oCap = and
			desc = fmt.Sprintf("p%v", i-or)
			outRequiredOr.Add(i)
		} else {
			_type = AND
			iCap = leaf + or
			if relaxed {
				oCap = or
			} else {
				oCap = 1
			}
			desc = fmt.Sprintf("r%v", i-leaf-or)
			inRequiredAnd.Add(i)
			outRequiredAnd.Add(i)
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

	// Process the lists (i.e. make connections)
	// Required AND -> OR edges
	// Shuffle both "outRequiredAnd" and "inRequiredOr" lists, then we match them index-wise (similar to Python `zip`)
	// Note that, because rand.Perm only makes a sequence with range [0, n), we need to offset the numbers
	outRequiredAndPerm := rnd.Perm(and)
	inRequiredOrPerm := rnd.Perm(or)

	var andNodeId, orNodeId int
	for i := 0; i < min(or, and); i++ {
		andNodeId = outRequiredAndPerm[i] + andPadding
		orNodeId = inRequiredOrPerm[i] + 1
		addEdge(V[andNodeId], V[orNodeId], cycleOk, &V)
		edge -= 1
	}
	// We run out of AND first, use `and` as our starting index for remaining OR nodes
	// Then connect a random AND to the OR
	if and < or {
		for i := and; i < or; i++ {
			orNodeId = inRequiredOrPerm[i] + 1
		reroll1:
			andNodeId = rnd.Intn(and) + andPadding
			if !addEdge(V[andNodeId], V[orNodeId], cycleOk, &V) {
				goto reroll1
			}
			edge -= 1
		}
	} else if or < and {
		// The opposite situation
		for i := or; i < and; i++ {
			andNodeId = outRequiredAndPerm[i] + andPadding
		reroll2:
			orNodeId = rnd.Intn(or) + 1
			if !addEdge(V[andNodeId], V[orNodeId], cycleOk, &V) {
				goto reroll2
			}
			edge -= 1
		}
	}

	// Required LEAF/OR -> AND edges
	outRequiredOrPerm := rnd.Perm(or + leaf - 1)
	inRequiredAndPerm := rnd.Perm(and)

	for i := 0; i < min(or+leaf-1, and); i++ {
		orNodeId = outRequiredOrPerm[i] + 2
		andNodeId = inRequiredAndPerm[i] + andPadding
		addEdge(V[orNodeId], V[andNodeId], cycleOk, &V)
		edge -= 1
	}
	// We run out of OR/LEAF first, use `or+leaf-1` as our starting index for remaining AND nodes
	if or+leaf-1 < and {
		for i := or + leaf - 1; i < and; i++ {
			andNodeId = inRequiredAndPerm[i] + andPadding
		reroll3:
			orNodeId = rnd.Intn(or+leaf-1) + 2
			if !addEdge(V[orNodeId], V[andNodeId], cycleOk, &V) {
				goto reroll3
			}
			edge -= 1
		}
		// The opposite situation
	} else if and < or+leaf-1 {
		for i := and; i < or+leaf-1; i++ {
			orNodeId = outRequiredOrPerm[i] + 2
		reroll4:
			andNodeId = rnd.Intn(and) + andPadding
			if !addEdge(V[orNodeId], V[andNodeId], cycleOk, &V) {
				goto reroll4
			}
			edge -= 1
		}
	}

	// First find all nodes with available OCap
	availSet := hset.New[int]()
	for i := 1; i <= total; i++ {
		if V[i].OCap > 0 {
			availSet.Add(i)
		}
	}

	for edge > 0 {
		src := rnd.Intn(availSet.Size())
		var targets *hset.Set[int]
		if V[src].Type == AND {
			targets = hset.New[int](makeRange(0, or-1)...)
		} else {
			targets = hset.New[int](makeRange(andPadding, total)...)
		}
		targets.Remove(V[src].Adj.Values()...)

		dstIdx := rnd.Intn(targets.Size())
		it := targets.Iterator()
		for dstIdx > 0 {
			it.Next()
			dstIdx -= 1
		}
		dst := it.Value()

		if !addEdge(V[src], V[dst], cycleOk, &V) {
			continue
		} else {
			if V[src].OCap < 1 {
				availSet.Remove(src)
			}
		}
	}

	return &V
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
