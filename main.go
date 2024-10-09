package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

// Item object containing name, weight and value
type Item struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Value  int     `json:"value"`
}

// Calculating total value and total weight of given solution
func computeEnergy(solution []int, items []Item) (totalValue int, totalWeight float64) {
	for i, included := range solution {
		// If items is included in knapsack
		if included == 1 {
			// Sum values and weights
			totalValue += items[i].Value
			totalWeight += items[i].Weight
		}
	}
	return
}

// Generating random solution array
func randomSolution(items []Item, rnd *rand.Rand) []int {
	// Initializing solution slice with same length as items array
	solution := make([]int, len(items))
	// Loop through solution slice
	for i := range solution {
		// Set random value: 0 or 1
		solution[i] = rnd.Intn(2)
	}
	return solution
}

// Generating the closest candidate solution array
func generateCandidate(solution []int, rnd *rand.Rand) []int {
	// Initializing candidate slice with same length as solution slice
	candidate := make([]int, len(solution))
	// Copy value from solution slice to candidate slice
	copy(candidate, solution)
	// Taking the random index of current solution
	index := rnd.Intn(len(solution))
	// Inverting index value
	candidate[index] = 1 - candidate[index]
	return candidate
}

// Returning 1 if candidate is better for sure
// Returning random float number from 0 to 1 if candidate might be better
func candidateIsBetter(curValue, candidateValue int, temp float64) float64 {
	if candidateValue > curValue {
		return 1.0
	}

	// Returning the base-e exponential of energy variation.
	return math.Exp(float64(candidateValue-curValue) / temp)
}

// Reading items from JSON file
func readItemsFromJSON(filename string) ([]Item, error) {
	// Opening the file
	file, err := os.Open(filename)
	// Checking if file exists
	if err != nil {
		return nil, err
	}
	// Closing file in the end of main function, even if error will occur
	defer file.Close()

	// Reading file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Deserializing JSON to array of items
	var items []Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Simulated Annealing algorithm
func simulatedAnnealing(items []Item, maxWeight, maxTemp, minTemp, coolingRate float64) ([]int, int) {
	// Randomizing seed for random
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)

	// Generating initial random solution
	curSolution := randomSolution(items, rnd)
	curValue, curWeight := computeEnergy(curSolution, items)

	// If weight of initial random solution exceeds maxWeight, trying to find a better solution
	for curWeight > maxWeight {
		curSolution = randomSolution(items, rnd)
		curValue, curWeight = computeEnergy(curSolution, items)
	}

	bestSolution := make([]int, len(curSolution))
	copy(bestSolution, curSolution)
	bestValue := curValue
	temp := maxTemp

	// Main simulated annealing loop
	for temp > minTemp {
		// Generating candidate solution and calculating it's weight and value
		candidateSolution := generateCandidate(curSolution, rnd)
		candidateValue, candidateWeight := computeEnergy(candidateSolution, items)

		// Skipping if weight of candidate solution is higher than max weight allowed
		if candidateWeight <= maxWeight {
			// Taking candidate solution if it's better or might be better
			if candidateIsBetter(curValue, candidateValue, temp) > rnd.Float64() {
				curSolution = candidateSolution
				curValue = candidateValue
			}

			// Updating best solution
			if candidateValue > bestValue {
				bestSolution = make([]int, len(candidateSolution))
				copy(bestSolution, candidateSolution)
				bestValue = candidateValue
			}

			// Cooling down the temperature
			temp *= coolingRate
		}
	}

	return bestSolution, bestValue
}

// Print list of items included in knapsack
func showKnapsack(solution []int, items []Item) {
	fmt.Println("List of items included in knapsack:")
	count := 0
	for i, included := range solution {
		if included == 1 {
			count++
			fmt.Printf(" - %s (Weight: %f, Value: %d)\n", items[i].Name, items[i].Weight, items[i].Value)
		}
	}

	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("Total items included: %d\n", count)
}

func main() {
	// Reading items from JSON file
	items, err := readItemsFromJSON("item_set_small.json")
	if err != nil {
		log.Fatalf("Error while reading the file: %v", err)
	}

	// Algorithm params
	maxWeight := 5.0
	maxTemp := 1000.0
	minTemp := 0.1
	coolingRate := 0.9

	// Record script start time
	start := time.Now()

	// Run simulated annealing algorithm
	bestSolution, bestValue := simulatedAnnealing(items, maxWeight, maxTemp, minTemp, coolingRate)
	fmt.Printf("Best solution: %v\n", bestSolution)
	showKnapsack(bestSolution, items)
	fmt.Printf("Total value: %d\n", bestValue)

	// Script execution time calculation
	duration := time.Since(start)
	fmt.Printf("Execution time: %v\n", duration)
	fmt.Printf("-------------------------------------------------------------")
}
