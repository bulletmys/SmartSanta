package algorithm

type node struct {
	id            int
	sccID         int
	visited       bool
	edges, rEdges []uint64
}

type scc struct {
	root *node
	len  int
}

type graph struct {
	nodes map[uint64]*node
	scc   map[int]*scc
	path  []uint64
}

func newGraph() *graph {
	var g graph
	g.nodes = make(map[uint64]*node)
	return &g
}

func newNode() *node {
	var n node
	n.sccID = -1
	return &n
}

func (g *graph) addEdge(t, h uint64) bool {
	if _, ok := g.nodes[t]; !ok {
		return false
	}
	if _, ok := g.nodes[h]; !ok {
		return false
	}
	g.nodes[t].edges = append(g.nodes[t].edges, h)
	g.nodes[h].rEdges = append(g.nodes[h].rEdges, t)

	return true
}

func (g *graph) addNode(label uint64) bool {
	if _, ok := g.nodes[label]; !ok {
		n := newNode()
		g.nodes[label] = n
		return true
	}
	return false
}

func (g *graph) resetVisited() {
	for _, n := range g.nodes {
		n.visited = false
	}
}

func (g *graph) createFinishingOrder() []*node {
	g.resetVisited()
	t := make([]*node, 0, len(g.nodes))
	for _, v := range g.nodes {
		if v.visited == false {
			dfsAssignFinishingNumber(v, g, &t)
		}
	}
	return t
}

func (g *graph) generateScc() {
	p := 1
	fo := g.createFinishingOrder()
	g.resetVisited()
	for i := len(fo) - 1; i >= 0; i-- {
		n := fo[i]
		if n.visited == false {
			s := p
			p++

			g.scc[s] = &scc{root: n}

			dfsMarkScc(n, g, s)
		}
	}
}

func dfsAssignFinishingNumber(n *node, g *graph, t *[]*node) {
	n.visited = true
	for _, neighbor := range n.rEdges {
		if g.nodes[neighbor].visited == false {
			dfsAssignFinishingNumber(g.nodes[neighbor], g, t)
		}
	}
	*t = append(*t, n)
}

func dfsMarkScc(n *node, g *graph, s int) {
	n.visited = true

	// ?? Если нода считалась другой компонентой связности ??
	if n.sccID != 0 && n.sccID != s {
		g.scc[n.sccID].len--
	}

	n.sccID = s
	for _, neighbor := range n.edges {
		if g.nodes[neighbor].visited == false {
			dfsMarkScc(g.nodes[neighbor], g, s)
		}
	}
}

// Нахождение Гамильтонова цикла в графе, состоящем
// из вершин одной компоненты связанности
func dfsFindPaths(first, n *node, g *graph, sccLen int) bool {
	n.visited = true
	for _, neighbor := range n.edges {
		if g.nodes[neighbor].sccID != n.sccID {
			continue
		}
		if g.nodes[neighbor].visited == false {

			g.path = append(g.path, neighbor)

			res := dfsFindPaths(first, g.nodes[neighbor], g, sccLen)
			if res {
				return true
			}
		}
		if n.id == first.id && len(g.path) == sccLen {
			return true
		}
	}

	g.path = g.path[:len(g.path)-1]
	n.visited = false

	return false
}

func pathToPairs(path []uint64, pairs map[uint64]uint64) {
	for i := 0; i < len(path)-1; i++ {
		pairs[path[i]] = path[i+1]
	}

	pairs[path[len(path)-1]] = path[0]
}

func calculatePairs(g *graph) map[uint64]uint64 {
	pairs := make(map[uint64]uint64)

	for _, v := range g.scc {
		g.path = make([]uint64, 0)

		res := dfsFindPaths(v.root, v.root, g, v.len)
		if !res {
			return nil
		}

		pathToPairs(g.path, pairs)
	}

	return pairs
}

func CountPreferences(nodes map[uint64][]uint64) map[uint64]uint64 {
	g := newGraph()

	for k, v := range nodes {
		g.addNode(k)
		for _, elem := range v {
			g.addNode(elem)
			if ok := g.addEdge(k, elem); !ok {
				return nil
			}
		}
	}

	g.generateScc()

	return calculatePairs(g)
}