package tree

type direction bool
type index int32

const MaxSize = (1 << 30) - 1

const (
	left  = direction(false)
	right = direction(true)
	null  = -1
)

type node struct {
	_parent index
	_left   index
	_right  index
	height  int8
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
		n._right = index(ix)
	} else {
		n._left = index(ix)
	}
}
