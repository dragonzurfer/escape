package algo

import (
	"escape/dijkstra-master"
	"fmt"
	"log"
	"math"
)

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type City struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Location  Location `json:"location"`
	Continent string   `json:"contId"`
	Country   string   `json:"countryName"`
}

type Continent struct {
	Cities []City
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Returns distance toCity in kms
func (city *City) Distance(toCity *City) float64 {
	R := 6371.0 // Radius of the earth in km
	dLat := deg2rad(toCity.Location.Latitude - city.Location.Latitude)
	dLon := deg2rad(toCity.Location.Longitude - city.Location.Longitude)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(deg2rad(city.Location.Latitude))*math.Cos(deg2rad(toCity.Location.Latitude))*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c // Distance in km
	return d
}

func SetEdge(graph *dijkstra.Graph, cityNameToIndex *map[string]int, distance float64, source, dest City) {
	sourceIndex := (*cityNameToIndex)[source.ID]
	destIndex := (*cityNameToIndex)[dest.ID]
	graph.AddArc(sourceIndex, destIndex, int64(distance))
	graph.AddArc(sourceIndex, destIndex, int64(distance))
}

func createGraph(originCity City, terminalNodeIndex int, continentsExcludingOrigin []Continent, indexToCity *map[int]City, cityNameToIndex *map[string]int) *dijkstra.Graph {
	//create graph
	graph := dijkstra.NewGraph()
	for i := 0; i <= terminalNodeIndex; i += 1 {
		graph.AddVertex(i)
	}

	for i := 0; i+1 < len(continentsExcludingOrigin); i += 1 {
		j := i + 1
		for _, cityA := range continentsExcludingOrigin[i].Cities {
			for _, cityB := range continentsExcludingOrigin[j].Cities {

				distance := cityA.Distance(&cityB)
				SetEdge(graph, cityNameToIndex, distance, cityA, cityB)
			}
		}

	}

	for _, city := range continentsExcludingOrigin[0].Cities {
		distance := originCity.Distance(&city)
		SetEdge(graph, cityNameToIndex, distance, originCity, city)
	}

	for _, city := range continentsExcludingOrigin[4].Cities {
		distance := originCity.Distance(&city)
		destIndex := terminalNodeIndex
		sourceIndex := (*cityNameToIndex)[city.ID]
		graph.AddArc(sourceIndex, destIndex, int64(distance))
	}

	return graph
}

func GetContinentsExcludingOrigin(originCity City, continentToCities *map[string][]City) []Continent {
	continentsExcludingOrigin := []Continent{}
	for continent, cities := range *continentToCities {
		if continent != originCity.Continent {
			continentsExcludingOrigin = append(continentsExcludingOrigin, Continent{Cities: cities})

		}
	}
	return continentsExcludingOrigin
}

func setIndexCityMappers(indexToCity *map[int]City, cityNameToIndex *map[string]int, continentsExcludingOrigin *[]Continent) {
	currentIndex := 1
	for _, continent := range *continentsExcludingOrigin {
		for _, city := range continent.Cities {
			(*indexToCity)[currentIndex] = city
			(*cityNameToIndex)[city.ID] = currentIndex
			currentIndex += 1
		}
	}
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func GetContinentOrder(permutation []int, continentsExcludingOrigin *[]Continent) []Continent {
	continentOrder := []Continent{}
	for _, index := range permutation {
		continentOrder = append(continentOrder, (*continentsExcludingOrigin)[index])
	}
	return continentOrder
}

// Generate 5! possible continent visit order graphs
// Run Dijkstra
// Return path with minimum total distance
// Time Complexity: 5!*NlogN where N is the total number of cities
func PrintBestPath(originCity City, continentToCities *map[string][]City) {
	continentsExcludingOrigin := GetContinentsExcludingOrigin(originCity, continentToCities)

	indexToCity := make(map[int]City)
	cityNameToIndex := make(map[string]int)
	setIndexCityMappers(&indexToCity, &cityNameToIndex, &continentsExcludingOrigin)

	//Reserve 0 and largest index for source and terminal
	terminalNodeIndex := len(indexToCity)
	indexToCity[terminalNodeIndex] = originCity
	indexToCity[0] = originCity

	//Get all possible permutation of continent visit order
	continentIndexPermutations := permutations([]int{0, 1, 2, 3, 4})

	//Set cost as max = 10*Earth radius
	lowestCost := 65630
	var optimal dijkstra.BestPath

	// Build graph over each continent order and run djikstra
	for index, permutation := range continentIndexPermutations {
		fmt.Println("Running permutation ", index+1, "/", len(continentIndexPermutations))
		continentOrder := GetContinentOrder(permutation, &continentsExcludingOrigin)
		graph := createGraph(originCity, terminalNodeIndex, continentOrder, &indexToCity, &cityNameToIndex)
		bestCurrentPath, err := graph.Shortest(cityNameToIndex[originCity.ID], terminalNodeIndex)

		if len(bestCurrentPath.Path) < 6 {
			continue
		}
		if int64(lowestCost) > bestCurrentPath.Distance {
			lowestCost = int(bestCurrentPath.Distance)
			optimal = bestCurrentPath
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	// Print best path among all possible continent permutations
	for _, index := range optimal.Path {
		fmt.Printf("%s(%s) -> ", indexToCity[index].Name, indexToCity[index].Continent)
	}
	fmt.Printf("\nDistance Travelled: %d KMS\n", optimal.Distance)
}
