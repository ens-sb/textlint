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
package null

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/ens-sb/textlint/chunk"
)

// CheckChunk processes a chunk of text from a file for validation.
// It reads the file specified in the chunk from the start to end position,
// validates the content, and signals completion through the control channel.
// The function handles file reading in buffered chunks to manage memory efficiently.
//
// Parameters:
//   - chunk: A pointer to a Chunk structure containing file path and byte range information.
//   - ctrl: A boolean channel used to signal completion of the chunk processing.
//
// The function will log a fatal error if it cannot open the file or seek to the
// specified position. Upon completion, it attempts to send a true value on the
// control channel, logging a message if the channel is not ready to receive.
func CheckChunk(chunk *chunk.Chunk, ctrl chan bool) {
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

	for {
		c, err := br.Read(buff)
		count += c

		ValidateBuffer(buff[:c])

		if count >= int(chunk.End-chunk.Start) {
			break
		}

		if err == io.EOF {
			break
		}
	}

	select {
	case ctrl <- true:
	default:
		log.Println("Control channel is closed or not ready to receive")
	}
}

// ValidateBuffer checks the input buffer for null bytes.
// It iterates through each byte in the buffer and triggers a fatal error
// if a null byte (0x00) is encountered.
func ValidateBuffer(buff []byte) {
	for _, char := range buff {
		if char == 0 {
			log.Fatal("Null byte found in file!")
		}
	}
}
