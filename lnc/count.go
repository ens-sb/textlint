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
package lnc

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/ens-sb/textlint/chunk"
)

type Counts struct {
	Newlines int
	Nulls    int
	Tabs     int
}

// CountChars processes a chunk of text from a file and counts various character statistics.
// It reads the specified chunk, updates statistics, and sends the results through a control channel.
//
// Parameters:
//   - chunk: A pointer to a chunk.Chunk that defines the file path and the range of bytes to process.
//   - ctrl: A channel to send the resulting character counts.
//
// The function will log.Fatal in the following cases:
//   - If chunk.End is less than chunk.Start
//   - If the file cannot be opened
//   - If seeking to the chunk's start position fails
//
// When processing is complete, the function attempts to send the statistics through the control channel.
// If the channel is not ready to receive, it logs a message and continues.
func CountChars(chunk *chunk.Chunk, ctrl chan *Counts) {
	if chunk.End < chunk.Start {
		log.Fatal("chunk.End is less than chunk.Start")
		return
	}
	fh, err := os.Open(chunk.File)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = fh.Seek(chunk.Start, 0)
	if err != nil {
		log.Fatal(err)
	}
	br := bufio.NewReaderSize(fh, 1024*1024)
	buff := make([]byte, chunk.End-chunk.Start)

	count := 0

	stats := &Counts{}

	for {
		c, err := br.Read(buff)
		count += c

		UpdateStats(stats, buff[:c])

		if count >= int(chunk.End-chunk.Start) {
			break
		}

		if err == io.EOF {
			break
		}
	}

	ctrl <- stats
}

// UpdateStats calculates and updates statistics from a buffer of bytes.
// It counts newlines and null bytes found in the input buffer and
// updates the provided Counts structure accordingly.
//
// Parameters:
//   - stats: A pointer to a Counts structure to be updated
//   - buff: A byte slice to be analyzed
//
// This function does not return any value as it modifies the stats in place.
func UpdateStats(stats *Counts, buff []byte) {
	for _, b := range buff {
		switch b {
		case '\n':
			stats.Newlines++
		case '\x00':
			stats.Nulls++
		case '\t':
			stats.Tabs++
		}
	}
}
