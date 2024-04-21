package lazyskiplist

func (this *SkipList) findCeiling(query interface{}) (node *Node) {
	if query == nil {
		return nil
	}
	pred := this.head
	for lv := MAX_LEVEL - 1; lv >= 0; lv-- {
		curr := pred.next[lv]
		for curr != this.tail && this.comparator(query, curr.key) > 0 {
			pred = curr
			curr = pred.next[lv]
		}

		if curr != this.tail && this.comparator(query, curr.key) == 0 {
			return curr
		}
	}
	return pred.next[0]
}

func (this *SkipList) findFloor(query interface{}) (node *Node) {
	if query == nil {
		return nil
	}
	pred := this.head
	for lv := MAX_LEVEL - 1; lv >= 0; lv-- {
		curr := pred.next[lv]
		for curr != this.tail && this.comparator(query, curr.key) > 0 {
			pred = curr
			curr = pred.next[lv]
		}

		if curr != this.tail && this.comparator(query, curr.key) == 0 {
			return curr
		}
	}
	return pred
}

func (this *SkipList) Ceiling(query interface{}) (key, value interface{}, found bool) {
	if node := this.findCeiling(query); node != nil {
		return node.key, node.value, true
	}
	return nil, nil, false
}

func (this *SkipList) Floor(query interface{}) (key, value interface{}, found bool) {
	if node := this.findFloor(query); node != nil {
		return node.key, node.value, true
	}
	return nil, nil, false
}

func (list *SkipList) Begin(query interface{}) *Iterator {
	if node := list.findCeiling(query); node != nil {
		return &Iterator{list: list, node: node}
	}
	return &Iterator{list: list, node: list.head.next[0]}
}

func (list *SkipList) End(query interface{}) *Iterator {
	if node := list.findFloor(query); node != nil {
		return &Iterator{list: list, node: node}
	}
	return &Iterator{list: list, node: list.tail.prev}
}
