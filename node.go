package tree

type direction bool

const (
	left  = direction(false)
	right = direction(true)
	null  = -1
)

type node struct {
	_parent int32
	_left   int32
	_right  int32
	height  int8
}

func (n *node) parent() int {
	return int(n._parent)
}

func (n *node) left() int {
	return int(n._left)
}

func (n *node) right() int {
	return int(n._right)
}

func (n *node) set_parent(i int) {
	n._parent = int32(i)
}

func (n *node) set_left(i int) {
	n._left = int32(i)
}

func (n *node) set_right(i int) {
	n._right = int32(i)
}

func (n *node) link(i direction) int {
	if i == right {
		return int(n._right)
	} else {
		return int(n._left)
	}
}

func (n *node) set_link(i direction, ix int) {
	if i == right {
		n._right = int32(ix)
	} else {
		n._left = int32(ix)
	}
}
