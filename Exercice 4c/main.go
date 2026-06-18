package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type result struct {
	workerID int
	number   int
	sum      int
}

func sumDivisors(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			sum += i
		}
	}
	return sum
}

func generateNumbers(numJobs int, jobs chan<- int) {
	defer close(jobs)

	for number := 1; number <= numJobs; number++ {
		jobs <- number * 100000
	}
}

func worker(id int, jobs <-chan int, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		results <- result{
			workerID: id,
			number:   job,
			sum:      sumDivisors(job),
		}
	}
}

func runWorkerPool(numWorkers int, numJobs int) []result {
	jobs := make(chan int)
	results := make(chan result, numJobs)
	var wg sync.WaitGroup

	go generateNumbers(numJobs, jobs)

	for id := 1; id <= numWorkers; id++ {
		wg.Add(1)
		go worker(id, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	collectedResults := make([]result, 0, numJobs)
	for currentResult := range results {
		collectedResults = append(collectedResults, currentResult)
	}

	return collectedResults
}

func printResults(results []result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].number < results[j].number
	})

	for _, currentResult := range results {
		fmt.Printf(
			"Worker %d -> nombre %d : somme des diviseurs = %d\n",
			currentResult.workerID,
			currentResult.number,
			currentResult.sum,
		)
	}
}

func main() {
	const numWorkers = 4
	const numJobs = 100

	startTime := time.Now()
	results := runWorkerPool(numWorkers, numJobs)
	duration := time.Since(startTime)

	fmt.Printf("Worker pool avec %d workers et %d taches\n", numWorkers, numJobs)
	fmt.Println()
	printResults(results)
	fmt.Println()
	fmt.Printf("Temps d'execution total : %s\n", duration)
}
