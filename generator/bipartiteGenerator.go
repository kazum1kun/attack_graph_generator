package generator

import (
	"fmt"
	wr "github.com/mroth/weightedrand/v2"
	al "github.com/ugurcsen/gods-generic/lists/arraylist"
	queue "github.com/ugurcsen/gods-generic/queues/linkedlistqueue"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
	"math"
	"math/rand"
)

//goland:noinspection GoDeprecation
func ConstructGraph(leaf, and, or, edge int, cycleOk, relaxed bool, rnd *rand.Rand) *[]*CNode {
	// We first construct 4 lists:
	// inRequiredOr is the list of OR that needs an inbound edge
	inRequiredOr := al.New[int]()
	// inRequiredAnd is the list of AND that needs an inbound edge
	inRequiredAnd := al.New[int]()
	// outRequiredOr is the list of OR/LEAF that needs an outbound edge
	outRequiredOr := al.New[int]()
	// outRequiredAnd is the list of AND that needs an outbound edge
	outRequiredAnd := al.New[int]()

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
			oCap = (or - 1) * and
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

	// Initialize a weighted random draw
	outChooser, inOrChooser, inAndChooser := initRand(&V, and, or, andPadding)

	success := 0.0
	// Randomly draw from the weighted out/in distributions
	for edge > 0 {
		// Re-populate the chooser with the latest info
		if success > math.Sqrt(float64(total)) {
			outChooser, inOrChooser, inAndChooser = initRand(&V, and, or, andPadding)
			success = 0.0
		}

		src := outChooser.PickSource(rnd)
		if V[src].Type == AND {
		reroll5:
			dst := inOrChooser.PickSource(rnd)
			if addEdge(V[src], V[dst], cycleOk, &V) {
				success += 1
				edge -= 1
			} else {
				goto reroll5
			}
		} else {
		reroll6:
			dst := inAndChooser.PickSource(rnd)
			if addEdge(V[src], V[dst], cycleOk, &V) {
				success += 1
				edge -= 1
			} else {
				goto reroll6
			}
		}
	}

	return &V
}

func addEdge(src, dst *CNode, cycleOk bool, V *[]*CNode) bool {
	// check for cycles (a cycle is found when the dst predecessors are a subset of src predecessor)
	if !cycleOk && src.Pred.Contains(dst.Pred.Values()...) {
		return false
	}
	// do not allow for duplicate edges in any case
	if src.Adj.Contains(dst.Id) {
		return false
	}

	src.Adj.Add(dst.Id)

	// Update the Pred info in the subtree rooted at dst
	// Disable this step for cycles (since they will lead to infinite loops)
	if !cycleOk {
		q := queue.New[*CNode]()
		q.Enqueue(dst)
		for !q.Empty() {
			node, _ := q.Dequeue()
			node.Pred = *node.Pred.Union(&src.Pred)
			for _, adjNode := range node.Adj.Values() {
				q.Enqueue((*V)[adjNode])
			}
		}
	}

	src.OCap -= 1
	dst.ICap -= 1
	return true
}

func initRand(V *[]*CNode, and, or, andPadding int) (*wr.Chooser[int, int], *wr.Chooser[int, int], *wr.Chooser[int, int]) {
	// Construct weighted avg arrays for outgoing and incoming edges
	outgoingTemp := make([]wr.Choice[int, int], len(*V))
	incomingOrTemp := make([]wr.Choice[int, int], or)
	incomingAndTemp := make([]wr.Choice[int, int], and)
	for i := 1; i < len(*V); i++ {
		if (*V)[i].OCap > 0 {
			outgoingTemp[i-1] = wr.NewChoice(i, (*V)[i].OCap)
		}
		if (*V)[i].ICap > 0 {
			if (*V)[i].Type == OR {
				incomingOrTemp[i-1] = wr.NewChoice(i, (*V)[i].ICap)
			} else {
				incomingAndTemp[i-andPadding] = wr.NewChoice(i, (*V)[i].ICap)
			}
		}
	}
	outgoingChooser, _ := wr.NewChooser(outgoingTemp...)
	incomingOrChooser, _ := wr.NewChooser(incomingOrTemp...)
	incomingAndChooser, _ := wr.NewChooser(incomingAndTemp...)

	return outgoingChooser, incomingOrChooser, incomingAndChooser
}
