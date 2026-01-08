package collab

type stringNode struct {
	str []byte
	leftSize int
	l 	*stringNode
	r	*stringNode
}

func (node stringNode) toByteArray() ([]byte, int) {
	if node.str != nil {
		return node.str, len(node.str)
	}
	var leftBytes []byte = nil
	var rightBytes []byte = nil
	var leftSize int
	var rightSize int

	if node.l != nil {
		leftBytes, leftSize = node.l.toByteArray()
	}
	if node.r != nil {
		rightBytes, rightSize = node.r.toByteArray()
	}
	
	size := leftSize + rightSize
	curr := 0
	totalBytes := make([]byte, size)

	for i := 0; i < leftSize; i++ {
		totalBytes[curr] = leftBytes[i]
		curr++
	}
	for i := 0; i < rightSize; i++ {
		totalBytes[curr] = rightBytes[i]
		curr++
	}
	return totalBytes, size
}

// implements the rope data structure
type FileText struct {
	head *stringNode
	nodeSize int
}

func CreateFileText(nodeSize int) FileText {
	text := FileText {
		head: nil,
		nodeSize: nodeSize,
	}

	// return pointer because struct is quite small
	return text
}

func (f FileText) ToByteArray() []byte {
	arr, _ := f.head.toByteArray()
	return arr
}


