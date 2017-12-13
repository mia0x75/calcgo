package calcgo

// OptimizedAST holds an optimized ast.
// For all integer and decimal the value was already interpreted.
// All operations are already interpreted, if both child nodes could already be
// interpreted.
type OptimizedAST struct {
	Node *OptimizedNode
}

// OptimizedNode holds an optimized node
type OptimizedNode struct {
	Type        NodeType       `json:"type"`
	Value       float64        `json:"value"`
	OldValue    string         `json:"old_value"`
	IsOptimized bool           `json:"is_optimized"`
	LeftChild   *OptimizedNode `json:"left"`
	RightChild  *OptimizedNode `json:"right"`
}

// Optimize optimizes an ast.
// Interprets all integer and decimal nodes.
// Interprets all operations, if their child nodes can already be interpreted
func Optimize(ast *AST) (*OptimizedAST, error) {
	if ast == nil {
		return nil, nil
	}

	optimizedNode, err := optimizeNode(ast.Node)
	if err != nil {
		return nil, err
	}

	return &OptimizedAST{Node: optimizedNode}, nil
}

func optimizeNode(node *Node) (*OptimizedNode, error) {
	var result float64
	var err error

	switch node.Type {
	case NInteger:
		result, err = interpretInteger(node)
	case NDecimal:
		result, err = interpretDecimal(node)
	case NVariable:
		return &OptimizedNode{
			Type:        NVariable,
			Value:       0,
			OldValue:    node.Value,
			IsOptimized: false,
			LeftChild:   nil,
			RightChild:  nil,
		}, nil
	case NAddition, NSubtraction, NMultiplication, NDivision:
		return optimizeOperator(node)
	default:
		return nil, ErrorInvalidNodeType
	}

	if err != nil {
		return nil, err
	}

	return &OptimizedNode{
		Type:        NDecimal,
		Value:       result,
		OldValue:    "",
		IsOptimized: true,
		LeftChild:   nil,
		RightChild:  nil,
	}, nil
}

func optimizeOperator(node *Node) (*OptimizedNode, error) {
	left, right, err := getOptimizedNodeChilds(node)
	if err != nil {
		return nil, err
	}

	if !left.IsOptimized || !right.IsOptimized {
		return &OptimizedNode{
			Type:        node.Type,
			Value:       0,
			OldValue:    "",
			IsOptimized: false,
			LeftChild:   left,
			RightChild:  right,
		}, nil
	}

	var result float64

	switch node.Type {
	case NAddition:
		result = left.Value + right.Value
	case NSubtraction:
		result = left.Value - right.Value
	case NMultiplication:
		result = left.Value * right.Value
	case NDivision:
		if right.Value == 0 {
			return nil, ErrorDivisionByZero
		}
		result = left.Value / right.Value
	}

	return &OptimizedNode{
		Type:        NDecimal,
		Value:       result,
		OldValue:    "",
		IsOptimized: true,
		LeftChild:   nil,
		RightChild:  nil,
	}, nil
}

func getOptimizedNodeChilds(node *Node) (*OptimizedNode, *OptimizedNode, error) {
	if node.LeftChild == nil {
		return nil, nil, ErrorMissingLeftChild
	}
	if node.RightChild == nil {
		return nil, nil, ErrorMissingRightChild
	}

	left, err := optimizeNode(node.LeftChild)
	if err != nil {
		return nil, nil, err
	}
	right, err := optimizeNode(node.RightChild)
	if err != nil {
		return nil, nil, err
	}

	return left, right, nil
}