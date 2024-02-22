package generator

import set "github.com/ugurcsen/gods-generic/sets/hashset"

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
