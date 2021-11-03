package controllers

import (
	"container/heap"
	"context"
	"math"
	"strconv"
	"strings"

	"googlemaps.github.io/maps"
)

type Costs struct {
	G int
	H int
}

type Node struct {
	Name   string
	Parent *Node
	Cost   int
	F      int
	Index  int
}

type Graph struct {
	Edges  map[string][]string
	Costs  map[string]Costs
	Cities map[string]string
	Nodes  *[]Node
}

func (graph *Graph) Initialise(initial string, goal string) {

	graph.Edges = make(map[string][]string)
	graph.Costs = make(map[string]Costs)
	graph.Cities = make(map[string]string)

	// Add Cities nad placeIDs
	{
		graph.AddCity("Blagoevgrad", "ChIJYw9_jwD2qhQRjfVAhLbcfEE")
		graph.AddCity("Burgas", "ChIJkZ38-WaSpkAR8E2_aRKgAAQ")
		graph.AddCity("Dobrich", "ChIJlXyLx4cWpUAR1BJ3D0ZVzrU")
		graph.AddCity("Gabrovo", "ChIJRYeV1uIPqUAREFC_aRKgAAQ")
		graph.AddCity("Haskovo", "ChIJd40tNtNQrRQR_LM3iebZYPI")
		graph.AddCity("Kardzhali", "ChIJh61qMnFurRQRIFK_aRKgAAQ")
		graph.AddCity("Kyustendil", "ChIJjzLFk2acqhQRZCSV9oAzXsw")
		graph.AddCity("Lovech", "ChIJU1G3GQreq0ARvTNimQqJT34")
		graph.AddCity("Montana", "ChIJ6cLcvaQ0q0AR1GVaXsGBXlI")
		graph.AddCity("Pazardzhik", "ChIJEXMb8QSwrBQRq2OZCeyXtdw")
		graph.AddCity("Pernik", "ChIJz0npQUnLqhQRgFW_aRKgAAQ")
		graph.AddCity("Pleven", "ChIJuwC0NInzq0ARIFa_aRKgAAQ")
		graph.AddCity("Plovdiv", "ChIJPXZIogjRrBQRoDgTb_rRcGQ")
		graph.AddCity("Razgrad", "ChIJk2qc9pwOr0ARh2DS1gUde-c")
		graph.AddCity("Ruse", "ChIJhaPif7RgrkARAmgoFlLbJ74")
		graph.AddCity("Shumen", "ChIJa8aV82eKpUARG-bOD1NZ2k4")
		graph.AddCity("Silistra", "ChIJUXS_U9fgr0ARkFm_aRKgAAQ")
		graph.AddCity("Sliven", "ChIJyTn6b8YgpkAR4Fm_aRKgAAQ")
		graph.AddCity("Smolyan", "ChIJu2HT0IpZrBQRgFq_aRKgAAQ")
		graph.AddCity("Sofia", "ChIJ9Xsxy4KGqkARYF6_aRKgAAQ")
		graph.AddCity("Stara Zagora", "ChIJNTnoMAtpqEARsZ-HYhTtunQ")
		graph.AddCity("Targovishte", "ChIJE2sLgW9Dr0ARuB90qa6IiVc")
		graph.AddCity("Varna", "ChIJodfzqotTpEARfIulcRyUJ1c")
		graph.AddCity("Veliko Tarnovo", "ChIJ3ZGUPUshqUARwF2_aRKgAAQ")
		graph.AddCity("Vidin", "ChIJJZAAoux5U0cRt9h3-OI6hR8")
		graph.AddCity("Vratsa", "ChIJp6svIl0Wq0ARQF2_aRKgAAQ")
		graph.AddCity("Yambol", "ChIJtRC3Vik0pkARaluYYwp8F_Q")
	}
	// Add edges
	{
		// Blagoevgrad
		graph.AddEdge("Blagoevgrad", "Kyustendil")
		graph.AddEdge("Blagoevgrad", "Pernik")
		graph.AddEdge("Blagoevgrad", "Sofia")
		graph.AddEdge("Blagoevgrad", "Pazardzhik")
		graph.AddEdge("Blagoevgrad", "Smolyan")
		// Burgas
		graph.AddEdge("Burgas", "Yambol")
		graph.AddEdge("Burgas", "Sliven")
		graph.AddEdge("Burgas", "Shumen")
		graph.AddEdge("Burgas", "Varna")
		// Dobrich
		graph.AddEdge("Dobrich", "Silistra")
		graph.AddEdge("Dobrich", "Shumen")
		graph.AddEdge("Dobrich", "Varna")
		// Gabrovo
		graph.AddEdge("Gabrovo", "Lovech")
		graph.AddEdge("Gabrovo", "Veliko Tarnovo")
		graph.AddEdge("Gabrovo", "Targovishte")
		// Haskovo
		graph.AddEdge("Haskovo", "Kardzhali")
		graph.AddEdge("Haskovo", "Yambol")
		graph.AddEdge("Haskovo", "Stara Zagora")
		// Kardzhali
		graph.AddEdge("Kardzhali", "Smolyan")
		graph.AddEdge("Kardzhali", "Haskovo")
		graph.AddEdge("Kardzhali", "Plovdiv")
		// Kyustendil
		graph.AddEdge("Kyustendil", "Blagoevgrad")
		graph.AddEdge("Kyustendil", "Pernik")
		// Lovech
		graph.AddEdge("Lovech", "Sofia")
		graph.AddEdge("Lovech", "Vratsa")
		graph.AddEdge("Lovech", "Pleven")
		graph.AddEdge("Lovech", "Veliko Tarnovo")
		graph.AddEdge("Lovech", "Gabrovo")
		// Montana
		graph.AddEdge("Montana", "Pernik")
		graph.AddEdge("Montana", "Sofia")
		graph.AddEdge("Montana", "Vratsa")
		graph.AddEdge("Montana", "Vidin")
		// Pazardzhik
		graph.AddEdge("Pazardzhik", "Sofia")
		graph.AddEdge("Pazardzhik", "Blagoevgrad")
		graph.AddEdge("Pazardzhik", "Plovdiv")
		graph.AddEdge("Pazardzhik", "Smolyan")
		// Pernik
		graph.AddEdge("Pernik", "Sofia")
		graph.AddEdge("Pernik", "Blagoevgrad")
		graph.AddEdge("Pernik", "Kyustendil")
		graph.AddEdge("Pernik", "Montana")
		graph.AddEdge("Pernik", "Vidin")
		// Pleven
		graph.AddEdge("Pleven", "Vratsa")
		graph.AddEdge("Pleven", "Lovech")
		graph.AddEdge("Pleven", "Veliko Tarnovo")
		graph.AddEdge("Pleven", "Ruse")
		// Plovdiv
		graph.AddEdge("Plovdiv", "Pazardzhik")
		graph.AddEdge("Plovdiv", "Smolyan")
		graph.AddEdge("Plovdiv", "Kardzhali")
		graph.AddEdge("Plovdiv", "Stara Zagora")
		// Razgrad
		graph.AddEdge("Razgrad", "Ruse")
		graph.AddEdge("Razgrad", "Silistra")
		graph.AddEdge("Razgrad", "Veliko Tarnovo")
		graph.AddEdge("Razgrad", "Targovishte")
		graph.AddEdge("Razgrad", "Shumen")
		// Ruse
		graph.AddEdge("Ruse", "Pleven")
		graph.AddEdge("Ruse", "Razgrad")
		graph.AddEdge("Ruse", "Silistra")
		graph.AddEdge("Ruse", "Veliko Tarnovo")
		// Shumen
		graph.AddEdge("Shumen", "Targovishte")
		graph.AddEdge("Shumen", "Razgrad")
		graph.AddEdge("Shumen", "Silistra")
		graph.AddEdge("Shumen", "Dobrich")
		graph.AddEdge("Shumen", "Varna")
		graph.AddEdge("Shumen", "Burgas")
		// Silistra
		graph.AddEdge("Silistra", "Ruse")
		graph.AddEdge("Silistra", "Razgrad")
		graph.AddEdge("Silistra", "Shumen")
		graph.AddEdge("Silistra", "Dobrich")
		// Sliven
		graph.AddEdge("Sliven", "Stara Zagora")
		graph.AddEdge("Sliven", "Yambol")
		graph.AddEdge("Sliven", "Burgas")
		// Smolyan
		graph.AddEdge("Smolyan", "Blagoevgrad")
		graph.AddEdge("Smolyan", "Pazardzhik")
		graph.AddEdge("Smolyan", "Plovdiv")
		graph.AddEdge("Smolyan", "Kardzhali")
		// Sofia
		graph.AddEdge("Sofia", "Vratsa")
		graph.AddEdge("Sofia", "Lovech")
		graph.AddEdge("Sofia", "Pazardzhik")
		graph.AddEdge("Sofia", "Blagoevgrad")
		graph.AddEdge("Sofia", "Pernik")
		graph.AddEdge("Sofia", "Montana")
		// Stara Zagora
		graph.AddEdge("Stara Zagora", "Plovdiv")
		graph.AddEdge("Stara Zagora", "Sliven")
		graph.AddEdge("Stara Zagora", "Haskovo")
		// Targovishte
		graph.AddEdge("Targovishte", "Veliko Tarnovo")
		graph.AddEdge("Targovishte", "Gabrovo")
		graph.AddEdge("Targovishte", "Razgrad")
		graph.AddEdge("Targovishte", "Shumen")
		// Varna
		graph.AddEdge("Varna", "Dobrich")
		graph.AddEdge("Varna", "Shumen")
		graph.AddEdge("Varna", "Burgas")
		// Veliko Tarnovo
		graph.AddEdge("Veliko Tarnovo", "Razgrad")
		graph.AddEdge("Veliko Tarnovo", "Gabrovo")
		graph.AddEdge("Veliko Tarnovo", "Targovishte")
		graph.AddEdge("Veliko Tarnovo", "Lovech")
		graph.AddEdge("Veliko Tarnovo", "Pleven")
		// Vidin
		graph.AddEdge("Vidin", "Montana")
		graph.AddEdge("Vidin", "Pernik")
		// Vratsa
		graph.AddEdge("Vratsa", "Montana")
		graph.AddEdge("Vratsa", "Sofia")
		graph.AddEdge("Vratsa", "Pleven")
		graph.AddEdge("Vratsa", "Lovech")
		// Yambol
		graph.AddEdge("Yambol", "Sliven")
		graph.AddEdge("Yambol", "Haskovo")
		graph.AddEdge("Yambol", "Burgas")
	}
	// Add costs
	graph.Costs = graph.findCosts(goal)
}
func (graph *Graph) AddCity(node string, placeID string) {

	graph.Cities[node] = placeID
}
func (graph *Graph) AddEdge(node1 string, node2 string) {

	graph.Edges[node1] = append(graph.Edges[node1], node2)
}
func (graph *Graph) findCosts(goal string) map[string]Costs {

	costs := make(map[string]Costs)
	for k, v := range graph.Cities {

		if k != goal {
			// Google API code to get the costs for V
			c, err := maps.NewClient(maps.WithAPIKey("AIzaSyDcZ1Dks3_6CnPF11dfvwEJtk60JiBaQBc"))
			checkError(err)

			request := &maps.DistanceMatrixRequest{Mode: "ModeDriving"}
			request.Origins = append(request.Origins, "place_id:"+v)
			request.Destinations = append(request.Destinations, "place_id:"+graph.Cities[goal])
			response, err := c.DistanceMatrix(context.Background(), request)
			checkError(err)
			distance := response.Rows[0].Elements[0].Distance
			duration := response.Rows[0].Elements[0].Duration
			distString := strings.Split(distance.HumanReadable, " ")[0]
			dist, err := strconv.ParseFloat(distString, 64)
			checkError(err)
			distInt := int(math.Round(dist))
			costs[k] = Costs{H: distInt, G: int(duration.Minutes())}
		}
	}

	return costs
}

func (graph *Graph) aStarSearch(intial string, goal string) []string {

	start := Node{Name: intial, Parent: nil, Cost: 0, F: 0}
	var path []string
	path = graph.search(start, goal)

	return path
}

func (graph *Graph) search(node Node, goal string) []string {

	priorityQueue := make(PriorityQueue, 0, 1000)
	heap.Init(&priorityQueue)
	heap.Push(&priorityQueue, &node)
	reached := make(map[string]*Node)

	for priorityQueue.Len() != 0 {

		node := heap.Pop(&PriorityQueue{}).(*Node)
		if node.Name == goal {

			nodes := make([]string, 0)
			i := 0
			for _, v := range reached {
				nodes[i] = v.Name
				i++
			}
			return nodes
		}
		children := graph.getChildren(node)
		for _, v := range children {
			if _, ok := reached[v.Name]; !ok {

				reached[v.Name] = v
				heap.Push(&priorityQueue, v)
			} else {
				if v.F < reached[v.Name].F {
					reached[v.Name] = v
					heap.Push(&priorityQueue, v)
				}
			}
		}
	}

	return nil
}

func (graph *Graph) getChildren(node *Node) []*Node {

	var children []*Node
	for _, v := range graph.Edges[node.Name] {

		c, err := maps.NewClient(maps.WithAPIKey("AIzaSyDcZ1Dks3_6CnPF11dfvwEJtk60JiBaQBc"))
		checkError(err)

		request := &maps.DistanceMatrixRequest{Mode: "ModeDriving"}
		request.Origins = append(request.Origins, "place_id:"+graph.Cities[node.Name])
		request.Destinations = append(request.Destinations, "place_id:"+graph.Cities[v])
		response, err := c.DistanceMatrix(context.Background(), request)
		checkError(err)
		duration := response.Rows[0].Elements[0].Duration
		cost := int(duration.Minutes())
		cost = node.Cost + cost

		n := Node{Name: v, Parent: node, Cost: cost, F: cost + graph.Costs[v].H}
		children = append(children, &n)
	}

	return children
}

// Priority Queue Definition
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Pop() interface{} {

	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}
func (pq *PriorityQueue) Push(node interface{}) {

	n := len(*pq)
	city := node.(*Node)
	city.Index = n
	*pq = append(*pq, city)
}
