tree provides balanced tree structure which acts as index
for sort.Interface and could be used for searching like sort.Search

```go
    data := sort.IntSlice{}
    tree := &tree.Tree{}
    for i:=0; i<N; i++ {
        data = append(data, rand.Intn(1<<30))
        tree.Insert(data)
    }

    v := rand.Intn(1<<30)
    ix := tree.Search(func(i int) bool {
        return data[i] >= v
    })

    fmt.Println("Min:", data[tree.Min()], " Max:", data[tree.Max()])
    for ix := tree.Next(-1); ix < tree.Len(); ix = tree.Next(ix) {
        fmt.Printf("%d ", data[ix])
    }
```


```
Benchmark_TreeInsert10   	  300000	      4770 ns/op
Benchmark_TreeInsert100  	   30000	     47859 ns/op
Benchmark_TreeInsert1000 	    2000	    568310 ns/op
Benchmark_TreeInsert10000	     200	   8360731 ns/op
Benchmark_TreeInsert30000	      50	  32095120 ns/op
Benchmark_TreeSearch10   	 1000000	      1301 ns/op
Benchmark_TreeSearch100  	  100000	     18330 ns/op
Benchmark_TreeSearch1000 	    5000	    261579 ns/op
Benchmark_TreeSearch10000	     500	   3741459 ns/op
Benchmark_TreeSearch30000	     100	  13748600 ns/op
Benchmark_SortInsert10   	  300000	      3629 ns/op
Benchmark_SortInsert100  	   30000	     42644 ns/op
Benchmark_SortInsert1000 	    1000	   1270424 ns/op
Benchmark_SortInsert10000	      10	 121639685 ns/op
Benchmark_SortInsert30000	       1	1225768439 ns/op
Benchmark_SortSearch10   	 1000000	      1088 ns/op
Benchmark_SortSearch100  	  100000	     15561 ns/op
Benchmark_SortSearch1000 	   10000	    224691 ns/op
Benchmark_SortSearch10000	     500	   3050891 ns/op
Benchmark_SortSearch30000	     100	  10371959 ns/op
ok  	_/home/yura/Projects/go-tree	35.386s
```
