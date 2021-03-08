package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	matrixSize = 10
)

var (
	matrixA = [matrixSize][matrixSize] int {}
	matrixB = [matrixSize][matrixSize] int {}
	result = [matrixSize][matrixSize] int{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize] int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	for col := 0; col < matrixSize; col++ {
		for i := 0; i < matrixSize; i++ {
			result[row][col] += matrixA[row][i] * matrixB[i][col]
		}
	}
}

func main() {
	fmt.Println("Working...")
	start := time.Now()
	generateRandomMatrix(&matrixA)
	generateRandomMatrix(&matrixB)
	for row := 0; row < matrixSize; row++ {
		workOutRow(row)
		fmt.Printf("%v\n", result[row])
	}
	elapsed := time.Since(start)
	fmt.Printf("Done! Processing took %s", elapsed)
}


