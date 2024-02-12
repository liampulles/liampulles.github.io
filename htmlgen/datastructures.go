package htmlgen

type Tree[T any] struct {
	item     T
	children []*Tree[T]
}

// Iterate the tree in breadth first order, i..e. level-by-level.
// If the func parameter returns true, then we stop iterating.
func (t *Tree[T]) IterateBreadthFirst(doWith func(node *Tree[T]) (stop bool)) {
	// We process the front of the queue, but add children to the back.
	queue := &DoubleLinkedList[*Tree[T]]{}
	queue.PushBack(t)

	for item, ok := queue.PopFront(); ok; item, ok = queue.PopFront() {
		if stop := doWith(item); stop {
			return
		}
		for _, child := range item.children {
			queue.PushBack(child)
		}
	}
}

// Can be used as a FIFO queue.
type DoubleLinkedList[T any] struct {
	first *dllNode[T]
	last  *dllNode[T]
}

type dllNode[T any] struct {
	prev *dllNode[T]
	item T
	next *dllNode[T]
}

func (n *dllNode[T]) pushBack(item T) *dllNode[T] {
	newNode := dllNode[T]{
		prev: n,
		item: item,
		next: n.next,
	}
	if n.next != nil {
		n.next.prev = &newNode
	}
	n.next = &newNode

	return &newNode
}

func (dll *DoubleLinkedList[T]) PushBack(item T) {
	// Add to the end

	// -> If last is nil, then first must also be nil, so set this single node for both.
	if dll.last == nil {
		dll.setRoot(item)
		return
	}

	// -> Else, push onto last
	dll.last = dll.last.pushBack(item)
}

func (dll *DoubleLinkedList[T]) PopFront() (found T, ok bool) {
	// Get from front

	// -> If front is nil, then last must also be nil, and there is nothing to return.
	if dll.first == nil {
		return
	}

	// -> Then its definitely the first item we need to return, we just to rewire first.
	found, ok = dll.first.item, true

	// -> Does the front have any successors? If not, then the list is empty - mark it so.
	if dll.first.next == nil {
		dll.first = nil
		dll.last = nil
		return
	}

	// -> Ok, then mark the successor as the front.
	dll.first = dll.first.next
	dll.first.prev = nil
	return
}

func (dll *DoubleLinkedList[T]) setRoot(item T) {
	newNode := dllNode[T]{item: item}
	dll.last = &newNode
	dll.first = &newNode
}
