package main

import (
	"fmt"
	al "github.com/ugurcsen/gods-generic/lists/arraylist"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
	"log"
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

func constructGraph(leaf, and, or, edge int, cycleOk, relaxed bool, seed int64) *[]*CNode {
	rnd := rand.New(rand.NewSource(seed))
	// We first construct 4 lists:
	// inRequiredOr is the list of OR that needs an inbound edge
	inRequiredOr := al.New[int]()
	// inRequiredAnd is the list of AND that needs an inbound edge
	inRequiredAnd := al.New[int]()
	// outRequiredOr is the list of OR/LEAF that needs an outbound edge
	outRequiredOr := al.New[int]()
	// outRequiredAnd is the list of AND that needs an outbound edge
	outRequiredAnd := al.New[int]()

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
			inRequiredOr.Add(i)
		} else if i <= or {
			_type = OR
			iCap = and
			oCap = (or - 1) * and
			desc = fmt.Sprintf("d%v", i-1)
			inRequiredOr.Add(i)
			outRequiredOr.Add(i)
		} else if i <= or+leaf {
			_type = LEAF
			iCap = 0
			oCap = and
			desc = fmt.Sprintf("p%v", i-leaf-1)
			outRequiredOr.Add(i)
		} else {
			_type = AND
			iCap = leaf + or
			if relaxed {
				oCap = or
			} else {
				oCap = 1
			}
			desc = fmt.Sprintf("r%v", i-leaf-or-1)
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

	total := or + leaf + and
	andPadding := 1 + leaf + or

	// Process the lists (i.e. make connections)
	// Required AND -> OR edges
	it := outRequiredAnd.Iterator()
	for it.Next() {
		andNodeId := it.Value()
		var orNodeId int
		if inRequiredOr.Size() > 0 {
			orNodeIdx := rnd.Intn(inRequiredOr.Size())
			orNodeId, _ = inRequiredOr.Get(orNodeIdx)
			inRequiredOr.Remove(orNodeIdx)
		} else {
			orNodeId = rnd.Intn(or) + 1
		}

		// Impossible for cycles to form at this stage, no cycle check necessary
		addEdge(V[andNodeId], V[orNodeId], cycleOk)
		edge -= 1
	}
	// Process residual OR nodes
	it = inRequiredOr.Iterator()
	for it.Next() {
		orNodeId := it.Value()
		andNodeId := rnd.Intn(and) + andPadding

		addEdge(V[andNodeId], V[orNodeId], cycleOk)
		edge -= 1
	}

	// We are done with required AND -> OR edges now.
	// Required LEAF/OR -> AND edges
	it = outRequiredOr.Iterator()
	for it.Next() {
		orNodeId := it.Value()
		var andNodeId, andNodeIdx int

	pickAnother:
		if inRequiredAnd.Size() > 0 {
			andNodeIdx = rnd.Intn(inRequiredAnd.Size())
			andNodeId, _ = inRequiredAnd.Get(andNodeIdx)
		} else {
			andNodeId = rnd.Intn(and) + andPadding
		}

		// Cycle is possible at this time
		if !addEdge(V[orNodeId], V[andNodeId], cycleOk) {
			goto pickAnother
		} else {
			edge -= 1
			inRequiredAnd.Remove(andNodeIdx)
		}
	}

	// Process residual AND nodes
	it = inRequiredAnd.Iterator()
	for it.Next() {
		andNodeId := it.Value()
		// Skip the goal node
	pickAnother2:
		orNodeId := rnd.Intn(or+leaf-1) + 2

		if !addEdge(V[orNodeId], V[andNodeId], cycleOk) {
			goto pickAnother2
		} else {
			edge -= 1
		}
	}

	// All required edges are satisfied, proceed to generate random edges for remaining edge quota
	attempts := 1
	for edge > 0 {
		if attempts > 100 {
			log.Println("WARN: edge generation seems to be stuck, consider reducing the edge parameter!")
		}
		// Exclude the goal node
		src := rnd.Intn(total-1) + 2
		// Determine the correct dst type for the given src
		var dst int
		if V[src].Type == OR || V[src].Type == LEAF {
			dst = rnd.Intn(and) + andPadding
		} else {
			dst = rnd.Intn(or) + 1
		}
		if addEdge(V[src], V[dst], cycleOk) {
			// Decrement edge count only if the add was successful; otherwise retry adding
			edge--
		} else {
			attempts += 1
		}
	}

	return &V
}

func addEdge(src, dst *CNode, cycleOk bool) bool {
	// check for cycles (a cycle is found when the dst predecessors are a subset of src predecessor)
	if !cycleOk && src.Pred.Contains(dst.Pred.Values()...) {
		return false
	}
	// do not allow for duplicate edges in any case
	if src.Adj.Contains(dst.Id) {
		return false
	}

	src.Adj.Add(dst.Id)
	dst.Pred.Union(&src.Pred)
	src.OCap -= 1
	dst.ICap -= 1
	return true
}
