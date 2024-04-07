package main

import (
	"encoding/csv"
	"math/rand/v2"
	"os"
	"sort"
)

const SAMPLE_SIZE = 100

func main() {
	reposFilepath := "data/repos.csv"

	randomSample(reposFilepath, SAMPLE_SIZE, "data/sample.csv")
}

func randomSample(reposFilepath string, sampleSize int, outputFile string) error {
	pcg := rand.NewPCG(123, 420)

	file, err := os.Open(reposFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	allRecords, err := reader.ReadAll()
	if err != nil {
		return err
	}

	header := allRecords[0]
	allRecords = allRecords[1:]

	if sampleSize > len(allRecords) {
		sampleSize = len(allRecords)
	}

	indices := *getIndices(sampleSize, len(allRecords), pcg)

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write header
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write sampled rows
	for _, index := range indices {
		if err := writer.Write(allRecords[index]); err != nil {
			return err
		}
	}

	return nil
}

func getIndices(sampleSize int, populationSize int, seed *rand.PCG) *[]int {
	random := rand.New(seed)

	indices := make([]int, sampleSize)
	generated := map[int]bool{}

	count := 0
	for count < sampleSize {
		index := random.IntN(populationSize - 1)
		if !generated[index] {
			indices[count] = index
			generated[index] = true
			count++
		}
	}

	sort.Slice(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})

	return &indices
}
