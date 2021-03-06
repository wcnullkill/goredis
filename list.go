package goredis

type listNode struct {
	prev *listNode
	next *listNode
	// string
	value *Sdshdr
}

type list struct {
	head *listNode
	tail *listNode
	len  int
}

// redis listCreate
func Newlist() *list {
	return &list{}
}

// redis listDup
func (l *list) Dup() *list {

	return nil
}

// redis listFirst
func (l *list) First() (*listNode, bool) {
	if l.len > 0 {
		return l.head, true
	}
	return nil, false
}

// redis listLast
func (l *list) Last() (*listNode, bool) {
	if l.len > 0 {
		return l.tail, true
	}
	return nil, false
}

// redis listPrevNode
func (l *list) Prev(node *listNode) (*listNode, bool) {
	if node != nil {
		return node.prev, true
	}
	return nil, false
}

// redis listNextNode
func (l *list) Next(node *listNode) (*listNode, bool) {
	if node != nil {
		return node.next, true
	}
	return nil, false
}

// redis listAddNodeHead
func (l *list) AddHead(node *listNode) bool {
	if node == nil {
		return false
	}
	if l.len == 0 {
		l.head = node
		l.tail = node
	} else {
		n := l.head
		l.head = node
		node.next = n
		n.prev = node
		node.prev = nil
	}
	l.len++
	return true
}

// redis listAddNodeTail
func (l *list) AddTail(node *listNode) bool {
	if node == nil {
		return false
	}
	if l.len == 0 {
		l.head = node
		l.tail = node
	} else {
		n := l.tail
		l.tail = node
		node.prev = n
		n.next = node
		node.next = nil
	}
	l.len++
	return true
}

// 将node插入old之前
// redis listInsertNode
func (l *list) InsertBefore(old *listNode, node *listNode) bool {
	if old == nil || node == nil {
		return false
	}

	// 如果old是head
	if old.isHead() {
		l.head = node
		node.prev = nil
	} else {
		old.prev.next = node
		node.prev = old.prev
	}
	node.next = old
	old.prev = node
	l.len++
	return true
}

// 将node插入old之后
// redis listInsertNode
func (l *list) InsertAfter(old *listNode, node *listNode) bool {
	if old == nil || node == nil {
		return false
	}
	// 如果old是tail
	if old.isTail() {
		l.tail = node
		node.next = nil
	} else {
		old.next.prev = node
		node.next = old.next
	}
	old.next = node
	node.prev = old
	l.len++
	return true

}

// redis listSearchKey
func (l *list) SearchKey(val *Sdshdr) (*listNode, bool) {
	node, exist := l.First()
	for exist {
		v := node.value
		if v.Compare(val) {
			return node, true
		}
		node, exist = l.Next(node)
	}
	return nil, false
}

// 查找第n个元素，从1开始
// redis listIndex
func (l *list) Index(n int) (*listNode, bool) {
	if n > l.len {
		return nil, false
	}
	node := l.head
	if n == 1 {
		return node, true
	}
	for i := 2; i < n; i++ {
		node, _ = l.Next(node)
	}
	return node, true
}

// redis listDel
func (l *list) Del(node *listNode) bool {
	if node == nil {
		return false
	}
	// if node is head and tail
	if node.prev == nil && node.next == nil {
		l.head = nil
		l.tail = nil
		l.len = 0
		node.free()
		return true
	}
	if node.prev == nil {
		l.head = node.next
		node.next.prev = nil
		node.free()
		l.len--
		return true
	}
	if node.next == nil {
		l.tail = node.prev
		node.prev.next = nil
		node.free()
		l.len--
		return true
	}
	node.next.prev = node.prev
	node.prev.next = node.next
	node.free()
	l.len--
	return true
}

// 将末尾元素弹出，插入表头
// redis listRotate
func (l *list) Rotate() bool {
	node, exist := l.Last()
	if !exist {
		return false
	}
	ok := l.Del(node)
	if !ok {
		return false
	}
	ok = l.AddHead(node)
	return ok
}

func (node *listNode) free() {
	node.next = nil
	node.prev = nil
}

func (node *listNode) isHead() bool {
	return node.prev == nil
}
func (node *listNode) isTail() bool {
	return node.next == nil
}
