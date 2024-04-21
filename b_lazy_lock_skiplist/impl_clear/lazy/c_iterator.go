package lazyskiplist

type Iterator struct {
	list *SkipList
	node *Node
}

func (it *Iterator) Next() bool {
	if finger := it.node.next[0]; finger != nil {
		it.node = finger
		return true
	}
	return false
}

func (it *Iterator) Prev() bool {
	if finger := it.node.prev; finger != nil {
		it.node = finger
		return true
	}
	return false
}

func (it *Iterator) Present() bool {
	return it.node != it.list.head && it.node != it.list.tail
}

func (it *Iterator) IsMarked() bool {
	return it.node.marked
}

func (it *Iterator) CompareTo(key interface{}) int {
	return it.list.comparator(it.node.key, key)
}

func (it *Iterator) Key() interface{} {
	return it.node.key
}

func (it *Iterator) Value() interface{} {
	return it.node.value
}
