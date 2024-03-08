package generator

import (
	"fmt"
	wr "github.com/mroth/weightedrand/v2"
	queue "github.com/ugurcsen/gods-generic/queues/linkedlistqueue"
	set "github.com/ugurcsen/gods-generic/sets/hashset"
	"math/rand"
)

//goland:noinspection GoDeprecation
func ConstructGraph(leaf, and, or, edge int, cycleOk, relaxed bool, rnd *rand.Rand) *[]*CNode {
	total := or + leaf + and
	leafPadding := 1 + or
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

	andTotalOCap := 0
	orTotalOCap := (or - 1) * and
	if relaxed {
		andTotalOCap = and * or
	} else {
		andTotalOCap = and
	}
	leafTotalOCap := leaf * and

	// Generate random edges
	for edge > 0 {
		// Pick from a random starting point
		// The probability that a point gets picked should be proportional to its available outgoing edges
		// (to avoid duplicate edges)
		// but since it's expensive to have a constantly-updating weighted random algorithms, we settle for just
		// tracking the categories and hope the RNG doesn't fail us...
		chooser, _ := wr.NewChooser(
			wr.NewChoice(OR, orTotalOCap),
			wr.NewChoice(AND, andTotalOCap),
			wr.NewChoice(LEAF, leafTotalOCap),
		)
		srcType := chooser.PickSource(rnd)

		var src, dst int
		// Pick a starting node
		if srcType == OR {
			// Skip index 0 and the goal node
			src = rnd.Intn(or) + 2
		} else if srcType == LEAF {
			src = rnd.Intn(leaf) + leafPadding
		} else {
			src = rnd.Intn(and) + andPadding
		}
		if V[src].OCap == 0 {
			continue
		}

		// Pick a ending node
		if srcType == AND {
			dst = rnd.Intn(or) + 1
		} else {
			dst = rnd.Intn(and) + andPadding
		}
		if V[dst].OCap == 0 {
			continue
		}

		// Update the numbers as necessary
		if addEdge(V[src], V[dst], cycleOk, &V) {
			edge -= 1
			if srcType == OR {
				orTotalOCap -= 1
			} else if srcType == LEAF {
				leafTotalOCap -= 1
			} else {
				andTotalOCap -= 1
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
