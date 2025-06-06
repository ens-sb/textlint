/*
See the NOTICE file distributed with this work for additional information
regarding copyright ownership.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package chunk

import (
	"log"
	"math/rand/v2"
	"os"
)

type Chunk struct {
	File  string
	Start int64
	End   int64
}

// GetChunks divides a file into a specified number of chunks of approximately equal size.
// It takes a file path and the number of chunks to create, and returns a slice of Chunk pointers.
// If the file size divided by numChunks is zero, it sets the chunk size to 1 byte.
// The last chunk may be larger than others to ensure all file content is covered.
// If there's an error accessing the file, the function will terminate the program.
func GetChunks(file string, numChunks int) []*Chunk {
	stats, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	size := stats.Size()
	chunkSize := size / int64(numChunks)
	if chunkSize == 0 {
		chunkSize = 1
	}

	chunks := make([]*Chunk, numChunks)

	for i := range numChunks {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numChunks-1 {
			end = size
		}
		chunks[i] = &Chunk{file, start, end}
	}

	// Randomize chunk order
	rand.Shuffle(len(chunks), func(i, j int) {
		chunks[i], chunks[j] = chunks[j], chunks[i]
	})

	return chunks
}
