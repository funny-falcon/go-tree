package tree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

var _ = fmt.Println

func check(t *testing.T, data sort.Interface, tree *Tree, ix int) int8 {
	node := &tree.nodes[ix]
	l, r := int(node._left), int(node._right)
	var lh, rh int8
	if l != null {
		if data.Less(ix, l) {
			t.Fatalf("%d < %d", ix, l)
		}
		lh = check(t, data, tree, l)
	} else {
		lh = 0
	}
	if r != null {
		if data.Less(r, ix) {
			t.Fatalf("%d < %d", r, ix)
		}
		rh = check(t, data, tree, r)
	} else {
		rh = 0
	}
	bal := lh - rh
	if bal < -1 || bal > 1 || tree.bal(ix) != bal {
		t.Fatalf("height fails: %d [%d, %d]",
			ix, lh, rh)
	}
	return max_i8(lh, rh) + 1
}

func check_iter(t *testing.T, data sort.Interface, tree *Tree) {
	if tree.Min() != tree.Next(-1) {
		t.Fatalf("min or next is wrong")
	}
	cnt := 0
	lesser := tree.Min()
	for ix := tree.Next(lesser); ix < tree.Len(); lesser, ix = ix, tree.Next(ix) {
		if data.Less(ix, lesser) {
			t.Fatalf("%d < %d", ix, lesser)
		}
		cnt++
	}
	if cnt != tree.Len()-1 {
		t.Fatalf("Iteration: %d < %d", cnt+1, tree.Len())
	}
}

func test_insert(t *testing.T, data sort.IntSlice) *Tree {
	tree := Tree{}
	for range data {
		tree.Insert(data)
		check(t, data, &tree, tree.root)
		check_iter(t, data, &tree)
	}
	return &tree
}

func Test_Insert(t *testing.T) {
	data := sort.IntSlice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	tree := test_insert(t, data)
	if tree.Min() != 0 {
		t.Fatalf("Min is not 0")
	}
	if tree.Max() != len(data)-1 {
		t.Fatalf("Max is not %d", len(data)-1)
	}
	for k := 0; k < 1000; k++ {
		tree = &Tree{}
		data = sort.IntSlice{}
		for i := 0; i < 100; i++ {
			v := rand.Intn(1000)
			ix := tree.Search(func(i int) bool {
				return data[i] >= v
			})
			if ix < len(data) && data[ix] == v {
				for j := 0; j < len(data); j++ {
					if data[j] == v && j != ix {
						t.Fatalf("missed duplicate")
					}
				}
				continue
			}
			data = append(data, v)
			if k&1 == 0 {
				tree.Insert(data)
			} else {
				tree.InsertBefore(ix)
			}
			check(t, data, tree, tree.root)
			check_iter(t, data, tree)
		}
		test_insert(t, data)
	}
}

func test_delete(t *testing.T, data sort.IntSlice, tree *Tree, maxn int) {
	for tree.Len() > 0 {
		v := 1 + rand.Intn(maxn)
		ix := tree.Search(func(i int) bool {
			return data[i] >= v
		})
		if ix == tree.Len() || data[ix] != v {
			for j := 0; j < tree.Len(); j++ {
				if data[j] == v {
					t.Fatalf("search failed")
				}
			}
			continue
		}
		nxt := tree.Next(ix)
		var nextv int
		if nxt != tree.Len() {
			nextv = data[nxt]
		}
		nxt = tree.Delete(data, ix)
		if data[tree.Len()] != v {
			t.Fatalf("Delete don't place value at last position")
		}
		if nxt < tree.Len() && data[nxt] != nextv {
			t.Fatalf("Delete returns wrong next")
		}
	}
}

func Test_Delete(t *testing.T) {
	data := sort.IntSlice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	tree := Tree{}
	for range data {
		tree.Insert(data)
	}
	test_delete(t, data, &tree, 13)
	data = sort.IntSlice{}
	f := map[int]struct{}{}
	tree = Tree{}
	for i := 100; i > 0; i-- {
		v := rand.Intn(1000)
		if _, ok := f[v]; ok {
			continue
		}
		data = append(data, v)
	}
	test_delete(t, data, &tree, 1000)
}

func Test_InitSorted(t *testing.T) {
	data := sort.IntSlice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	tree := Tree{}
	tree.InitSorted(data.Len())
	check(t, data, &tree, tree.root)
	check_iter(t, data, &tree)

	data = sort.IntSlice{}
	f := map[int]struct{}{}
	tree = Tree{}
	for i := 100; i > 0; i-- {
		v := rand.Intn(1000)
		if _, ok := f[v]; ok {
			continue
		}
		data = append(data, v)
	}
	sort.Sort(data)
	tree.InitSorted(data.Len())
	check(t, data, &tree, tree.root)
	check_iter(t, data, &tree)
}

type tstruct struct {
	I  int
	Ix int
}

type tslice []tstruct

func (b tslice) Len() int           { return len(b) }
func (b tslice) Less(i, j int) bool { return b[i].I < b[j].I }
func (b tslice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b tslice) CheckSorted(t *testing.T) {
	sz := b.Len()
	for i := 1; i < sz; i++ {
		if b.Less(i, i-1) {
			t.Fatalf("b[%d].I=%v > b[%d].I=%v",
				i-1, b[i-1].I, i, b[i].I)
		} else if !b.Less(i-1, i) && b[i-1].Ix > b[i].Ix {
			t.Logf("b[%d].I=%v == b[%d].I=%v",
				i-1, b[i-1].I, i, b[i].I)
			t.Fatalf("b[%d].Ix=%v > b[%d].Ix=%v",
				i-1, b[i-1].Ix, i, b[i].Ix)
		}
	}
}

func trand(sz int) tslice {
	res := make(tslice, sz)
	for i := 0; i < sz; i++ {
		res[i].I = rand.Intn(sz / 4)
		res[i].Ix = i
	}
	return res
}

func Test_StableSort(t *testing.T) {
	sizes := [...]int{10, 50, 150}
	for _, sz := range sizes {
		res := trand(sz)
		StableSort(res)
		res.CheckSorted(t)
	}
}

type bigstruct struct {
	I  int
	Sl [2][]int
}

type benchslice []bigstruct

func (b benchslice) Len() int           { return len(b) }
func (b benchslice) Less(i, j int) bool { return b[i].I < b[j].I }
func (b benchslice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func random_tree(data *benchslice, tree *Tree, n int) {
	for j := 0; j < n; j++ {
		v := rand.Intn(1 << 30)
		*data = append(*data, bigstruct{I: v})
		tree.Insert(*data)
	}
}

func random_slice(data *benchslice, n int) {
	for j := 0; j < n; j++ {
		v := rand.Intn(1 << 30)
		ix := sort.Search(len(*data), func(i int) bool {
			return (*data)[i].I >= v
		})
		*data = append(*data, bigstruct{})
		copy((*data)[ix+1:], (*data)[ix:])
		(*data)[ix] = bigstruct{I: v}
	}
}

func benchmark_TreeInsert(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		data := benchslice{}
		tree := Tree{}
		random_tree(&data, &tree, n)
	}
}

func benchmark_TreeSearch(b *testing.B, n int) {
	data := benchslice{}
	tree := Tree{}
	random_tree(&data, &tree, n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			v := rand.Intn(1 << 30)
			ix := tree.Search(func(i int) bool {
				return data[i].I >= v
			})
			if ix < tree.Len() {
				if data[ix].I < v {
					b.Fatalf("search failed")
				}
				prev := tree.Prev(ix)
				if prev > -1 && data[prev].I >= v {
					b.Fatalf("search failed")
				}
			} else {
				if data[tree.Max()].I >= v {
					b.Fatalf("search failed")
				}
			}
		}
	}
}

func benchmark_TreeSort(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		data := benchslice{}
		tree := Tree{}
		random_tree(&data, &tree, n)
		tree.LeaveSorted(data)
		if !sort.IsSorted(data) {
			panic("not sorted")
		}
	}
}

func benchmark_SortInsert(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		data := benchslice{}
		random_slice(&data, n)
	}
}

func benchmark_SortSearch(b *testing.B, n int) {
	data := benchslice{}
	random_slice(&data, n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			v := rand.Intn(1 << 30)
			ix := sort.Search(len(data), func(i int) bool {
				return data[i].I >= v
			})
			if ix < len(data) {
				if data[ix].I < v {
					b.Fatalf("search failed")
				}
				if ix > 0 && data[ix-1].I >= v {
					b.Fatalf("search failed")
				}
			} else {
				if data[len(data)-1].I >= v {
					b.Fatalf("search failed")
				}
			}
		}
	}
}

func benchmark_SortSort(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		data := benchslice{}
		for j := 0; j < n; j++ {
			data = append(data, bigstruct{I: rand.Intn(1 << 30)})
		}
		sort.Sort(data)
		if !sort.IsSorted(data) {
			panic("not sorted")
		}
	}
}

func benchmark_SortStable(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		data := benchslice{}
		for j := 0; j < n; j++ {
			data = append(data, bigstruct{I: rand.Intn(1 << 30)})
		}
		sort.Stable(data)
		if !sort.IsSorted(data) {
			panic("not sorted")
		}
	}
}

func Benchmark_TreeInsert10(b *testing.B)    { benchmark_TreeInsert(b, 10) }
func Benchmark_TreeInsert100(b *testing.B)   { benchmark_TreeInsert(b, 100) }
func Benchmark_TreeInsert1000(b *testing.B)  { benchmark_TreeInsert(b, 1000) }
func Benchmark_TreeInsert10000(b *testing.B) { benchmark_TreeInsert(b, 10000) }
func Benchmark_TreeInsert30000(b *testing.B) { benchmark_TreeInsert(b, 30000) }
func Benchmark_TreeSearch10(b *testing.B)    { benchmark_TreeSearch(b, 10) }
func Benchmark_TreeSearch100(b *testing.B)   { benchmark_TreeSearch(b, 100) }
func Benchmark_TreeSearch1000(b *testing.B)  { benchmark_TreeSearch(b, 1000) }
func Benchmark_TreeSearch10000(b *testing.B) { benchmark_TreeSearch(b, 10000) }
func Benchmark_TreeSearch30000(b *testing.B) { benchmark_TreeSearch(b, 30000) }
func Benchmark_TreeSort10(b *testing.B)      { benchmark_TreeSort(b, 10) }
func Benchmark_TreeSort100(b *testing.B)     { benchmark_TreeSort(b, 100) }
func Benchmark_TreeSort1000(b *testing.B)    { benchmark_TreeSort(b, 1000) }
func Benchmark_TreeSort10000(b *testing.B)   { benchmark_TreeSort(b, 10000) }
func Benchmark_TreeSort30000(b *testing.B)   { benchmark_TreeSort(b, 30000) }
func Benchmark_SortInsert10(b *testing.B)    { benchmark_SortInsert(b, 10) }
func Benchmark_SortInsert100(b *testing.B)   { benchmark_SortInsert(b, 100) }
func Benchmark_SortInsert1000(b *testing.B)  { benchmark_SortInsert(b, 1000) }
func Benchmark_SortInsert10000(b *testing.B) { benchmark_SortInsert(b, 10000) }
func Benchmark_SortInsert30000(b *testing.B) { benchmark_SortInsert(b, 30000) }
func Benchmark_SortSearch10(b *testing.B)    { benchmark_SortSearch(b, 10) }
func Benchmark_SortSearch100(b *testing.B)   { benchmark_SortSearch(b, 100) }
func Benchmark_SortSearch1000(b *testing.B)  { benchmark_SortSearch(b, 1000) }
func Benchmark_SortSearch10000(b *testing.B) { benchmark_SortSearch(b, 10000) }
func Benchmark_SortSearch30000(b *testing.B) { benchmark_SortSearch(b, 30000) }
func Benchmark_SortSort10(b *testing.B)      { benchmark_SortSort(b, 10) }
func Benchmark_SortSort100(b *testing.B)     { benchmark_SortSort(b, 100) }
func Benchmark_SortSort1000(b *testing.B)    { benchmark_SortSort(b, 1000) }
func Benchmark_SortSort10000(b *testing.B)   { benchmark_SortSort(b, 10000) }
func Benchmark_SortSort30000(b *testing.B)   { benchmark_SortSort(b, 30000) }
func Benchmark_SortStable10(b *testing.B)    { benchmark_SortStable(b, 10) }
func Benchmark_SortStable100(b *testing.B)   { benchmark_SortStable(b, 100) }
func Benchmark_SortStable1000(b *testing.B)  { benchmark_SortStable(b, 1000) }
func Benchmark_SortStable10000(b *testing.B) { benchmark_SortStable(b, 10000) }
func Benchmark_SortStable30000(b *testing.B) { benchmark_SortStable(b, 30000) }
