// tree provides balanced tree structure which acts as index
// for sort.Interface and could be used for searching like sort.Search
//
//     data := sort.IntSlice{}
//     tree := &tree.Tree{}
//     for i:=0; i<N; i++ {
//         data = append(data, rand.Intn(1<<30))
//         tree.Insert(data)
//     }
//
//     v := rand.Intn(1<<30)
//     ix := tree.Search(func(i int) bool {
//         return data[i] >= v
//     })
//
//     fmt.Println("Min:", data[tree.Min()], " Max:", data[tree.Max()])
//     for ix := tree.Next(-1); ix < tree.Len(); ix = tree.Next(ix) {
//         fmt.Printf("%d ", data[ix])
//     }
package tree

import "sort"

// Tree provides balanced tree structure
// which keeps order of is external sort.Interface
type Tree struct {
	root, min, max int
	nodes          []node
}

// Len returns number of indexed elements
func (t *Tree) Len() int {
	return len(t.nodes)
}

// Min returns index of minimum element
// panics if called on empty tree
func (t *Tree) Min() int {
	if len(t.nodes) == 0 {
		panic("Tree.Min should not be called on empty tree")
	}
	return t.min
}

// Max returns index of maximum element
// panics if called on empty tree
func (t *Tree) Max() int {
	if len(t.nodes) == 0 {
		panic("Tree.Max should not be called on empty tree")
	}
	return t.max
}

// Search returns first index for which predicate is true
// returns Len() if no element satisfies predicate
func (t *Tree) Search(pred func(i int) bool) int {
	if len(t.nodes) == 0 {
		return 0
	}
	now := t.root
	last_true := len(t.nodes)
	for {
		node := &t.nodes[now]
		if pred(now) {
			last_true = now
			if node._left == null {
				return now
			}
			now = int(node._left)
		} else if node._right == null {
			return last_true
		} else {
			now = int(node._right)
		}
	}
}

// SearchLast returns last index for which predicate is true
// returns -1 if no element satisfies predicate
func (t *Tree) SearchLast(pred func(i int) bool) int {
	if len(t.nodes) == 0 {
		return 0
	}
	now := t.root
	last_true := -1
	for {
		node := &t.nodes[now]
		if pred(now) {
			last_true = now
			if node._right == null {
				return now
			}
			now = int(node._right)
		} else if node._left == null {
			return last_true
		} else {
			now = int(node._left)
		}
	}
}

// Next returns index of next in-order element.
// if argument is -1, then return index of minimal element.
// returns t.Len() on finish.
func (t *Tree) Next(i int) int {
	if i > len(t.nodes) {
		panic("Tree index overflow")
	}
	if i == len(t.nodes) || i == t.max {
		return len(t.nodes)
	}
	if i == -1 {
		return t.min
	}
	node := &t.nodes[i]
	if node._right != null {
		i = int(node._right)
		for t.nodes[i]._left != null {
			i = int(t.nodes[i]._left)
		}
		return i
	}
	for node._parent != null {
		pix := int(node._parent)
		dir := t.dir(i, pix)
		if dir == left {
			return pix
		}
		i, node = pix, &t.nodes[pix]
	}
	panic("tree broken")
}

// Prev returns index of previos in-order element.
// if argument is Tree.Len(), then return index of maximal element.
// returns -1 on finish.
func (t *Tree) Prev(i int) int {
	if i > len(t.nodes) {
		panic("Tree index overflow")
	}
	if i == -1 || i == t.min {
		return -1
	}
	if i == len(t.nodes) {
		return t.max
	}
	node := &t.nodes[i]
	if node._left != null {
		i = int(node._left)
		for t.nodes[i]._right != null {
			i = int(t.nodes[i]._right)
		}
		return i
	}
	for node._parent != null {
		pix := int(node._parent)
		dir := t.dir(i, pix)
		if dir == right {
			return pix
		}
		i, node = pix, &t.nodes[pix]
	}
	panic("tree broken")
}

// Insert adds in-order element of sort.Interface at index Tree.Len()
// You should not put element with key equal to existed key.
func (t *Tree) Insert(data sort.Interface) {
	ix := len(t.nodes)
	t.nodes = append(t.nodes, node{null, null, null, 1})
	if ix == 0 {
		t.root, t.min, t.max = 0, 0, 0
		return
	}
	var dir direction
	cur := t.root
	curnode := &t.nodes[cur]
	for {
		dir = direction(data.Less(int(cur), ix))

		if curnode.link(dir) == null {
			break
		}
		cur = curnode.link(dir)
		curnode = &t.nodes[cur]
	}
	node := &t.nodes[ix]
	node._parent = int32(cur)
	curnode.set_link(dir, ix)
	if dir == right {
		if cur == t.max {
			t.max = ix
		}
	} else if cur == t.min {
		t.min = ix
	}
	t.balance(cur)
}

// InsertBefore adds new element at specified position.
// It trust you and doesn't check insertion position
func (t *Tree) InsertBefore(data sort.Interface, cur int) {
	ix := len(t.nodes)
	dir := left
	t.nodes = append(t.nodes, node{null, null, null, 1})
	if ix == 0 {
		if cur != 0 {
			panic("InsertBefore on empty tree accepts only 0")
		}
		t.root, t.min, t.max = 0, 0, 0
		return
	}
	var curnode *node
	if cur == ix {
		dir, cur = right, t.max
		curnode = &t.nodes[cur]
	} else {
		curnode = &t.nodes[cur]
		if curnode._left != null {
			dir = right
			cur = t.Prev(cur)
			curnode = &t.nodes[cur]
		}
	}
	node := &t.nodes[ix]
	node._parent = int32(cur)
	curnode.set_link(dir, ix)
	if dir == right {
		if cur == t.max {
			t.max = ix
		}
	} else if cur == t.min {
		t.min = ix
	}
	t.balance(cur)
}

// Delete removes element from a tree and return index of next in-order element
func (t *Tree) Delete(data sort.Interface, ix int) int {
	if ix < 0 || ix >= len(t.nodes) {
		panic("Tree.Delete out of range")
	}
	if len(t.nodes) == 0 {
		t.nodes = t.nodes[:0]
		return 0
	}
	node := &t.nodes[ix]
	next := t.Next(ix)
	if node._left != null && node._right != null {
		data.Swap(ix, next)
		next, ix, node = ix, next, &t.nodes[next]
		/* at this moment order is temporary broken,
		   but it will be restored after complete */
	}
	return t.del(data, node, ix, next)
}

// DeleteAndPrev removes element from a tree and return index of next in-order element
func (t *Tree) DeleteAndPrev(data sort.Interface, ix int) int {
	if ix < 0 || ix >= len(t.nodes) {
		panic("Tree.Delete out of range")
	}
	if len(t.nodes) == 0 {
		t.nodes = t.nodes[:0]
		return 0
	}
	node := &t.nodes[ix]
	prev := t.Prev(ix)
	if node._left != null && node._right != null {
		data.Swap(ix, prev)
		prev, ix, node = ix, prev, &t.nodes[prev]
		/* at this moment order is temporary broken,
		   but it will be restored after complete */
	}
	return t.del(data, node, ix, prev)
}

func (t *Tree) del(data sort.Interface, node *node, ix, next int) int {
	pix := int(node._parent)
	if pix == null {
		if node._left == null {
			rix := int(node._right)
			t.root = rix
			if rix != null {
				t.nodes[rix]._parent = null
			}
			if t.min == ix {
				t.min = rix
			}
		} else {
			lix := int(node._left)
			t.root = lix
			if lix != null {
				t.nodes[lix]._parent = null
			}
			if t.max == ix {
				t.max = lix
			}
		}
		pix = t.root
	} else {
		pdir := t.dir(ix, pix)
		parent := &t.nodes[pix]
		if node._left == null {
			rix := int(node._right)
			parent.set_link(pdir, rix)
			if rix != null {
				t.nodes[rix]._parent = int32(pix)
				if t.min == ix {
					t.min = rix
				}
			} else {
				if t.max == ix {
					t.max = pix
				}
				if t.min == ix {
					t.min = pix
				}
			}
		} else {
			lix := int(node._left)
			parent.set_link(pdir, lix)
			if lix != null {
				t.nodes[lix]._parent = int32(pix)
				if t.max == ix {
					t.max = lix
				}
			} else if t.max == ix {
				t.max = pix
			}
		}
	}
	if ix != len(t.nodes)-1 {
		jx := len(t.nodes) - 1
		data.Swap(ix, jx)
		inode, jnode := &t.nodes[ix], &t.nodes[jx]
		*inode, *jnode = *jnode, *inode
		t.fixlinks(inode, jx, ix)
		if next == jx {
			next = ix
		}
		if pix == jx {
			pix = ix
		}
	}
	t.nodes = t.nodes[:len(t.nodes)-1]
	t.balance(pix)
	return next
}

func (t *Tree) fixlinks(inode *node, i, j int) {
	if inode._parent != null {
		parent := &t.nodes[inode._parent]
		dir := direction(int(parent._right) == i)
		parent.set_link(dir, j)
	} else if t.root == i {
		t.root = j
	} else {
		panic("tree broken")
	}
	if inode._left != null {
		lnode := &t.nodes[inode._left]
		if int(lnode._parent) != i {
			panic("tree broken")
		}
		lnode._parent = int32(j)
	} else if t.min == i {
		t.min = j
	}
	if inode._right != null {
		rnode := &t.nodes[inode._right]
		if int(rnode._parent) != i {
			panic("tree broken")
		}
		rnode._parent = int32(j)
	} else if t.max == i {
		t.max = j
	}
}

func (t *Tree) balance(cur int) {
	for cur != null {
		node := &t.nodes[cur]
		lh, rh := t.height(node._left), t.height(node._right)
		var dir direction
		if lh < rh-1 {
			dir = right
		} else if lh-1 > rh {
			dir = left
		} else {
			node.height = max_i8(lh, rh) + 1
			cur = int(node._parent)
			continue
		}
		chld := node.link(dir)
		chnode := &t.nodes[chld]
		hs := [2]int8{
			t.height(int32(chnode.link(!dir))),
			t.height(int32(chnode.link(dir)))}
		if hs[1]-hs[0] < 0 {
			/* rotate child */
			t.rotate(chld, !dir)
		}
		t.rotate(cur, dir)
		cur = int(node._parent)
	}
}

// rotate tree node to direction
func (t *Tree) rotate(ix int, dir direction) {
	node := &t.nodes[ix]
	p := node._parent
	ch := node.link(dir)
	if ch == null {
		panic("wrong rotation direction")
	}
	chnode := &t.nodes[ch]
	node.set_link(dir, chnode.link(!dir))
	if node.link(dir) != null {
		t.nodes[node.link(dir)]._parent = int32(ix)
	}
	chnode.set_link(!dir, ix)
	node._parent = int32(ch)
	chnode._parent = int32(p)
	t.fixheight(node)
	t.fixheight(chnode)
	if p != null {
		pnode := &t.nodes[p]
		pdir := direction(int(pnode._right) == ix)
		pnode.set_link(pdir, ch)
		t.fixheight(pnode)
	} else {
		t.root = ch
	}
}

func (t *Tree) fixheight(n *node) {
	lh, rh := t.height(n._left), t.height(n._right)
	n.height = max_i8(lh, rh) + 1
}

func (t *Tree) height(ix int32) int8 {
	if ix == null {
		return 0
	}
	return t.nodes[ix].height
}

func (t *Tree) dir(i, ipar int) direction {
	parent := &t.nodes[ipar]
	if int(parent._left) == i {
		return left
	} else if int(parent._right) == i {
		return right
	} else {
		panic("tree broken")
	}
}

func (t *Tree) bal(ix int) int8 {
	node := &t.nodes[ix]
	lh := t.height(node._left)
	rh := t.height(node._right)
	return lh - rh
}

func max_i8(i, j int8) int8 {
	if i < j {
		return j
	}
	return i
}
