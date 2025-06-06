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
package cmd

import (
	"fmt"
	"log"
	"runtime"

	"github.com/ens-sb/textlint/chunk"
	"github.com/ens-sb/textlint/lnc"
	"github.com/spf13/cobra"
)

var lncCmd = &cobra.Command{
	Use:   "lnc",
	Short: "Count the number of newline, tab and null characters in a file",
	Long:  `Count the number of newline, tab and null characters in a file.`,
	Run: func(cmd *cobra.Command, args []string) {
		threads, err := cmd.Flags().GetInt("threads")
		checkErr(err)
		runtime.GOMAXPROCS(threads)

		batches, err := cmd.Flags().GetInt("batches")
		checkErr(err)

		files := cmd.Flags().Args()
		if len(files) != 1 {
			log.Fatal("Please provide a single file to check")
		}
		chunks := chunk.GetChunks(files[0], batches)
		ctrl := make(chan *lnc.Counts, batches)

		for _, chunk := range chunks {
			go lnc.CountChars(chunk, ctrl)
		}

		globalStats := &lnc.Counts{}
		for range batches {
			st := <-ctrl
			globalStats.Newlines += st.Newlines
			globalStats.Nulls += st.Nulls
			globalStats.Tabs += st.Tabs
		}

		fmt.Println("Newlines\tNulls\tTabs")
		fmt.Printf("%d\t%d\t%d\n", globalStats.Newlines, globalStats.Nulls, globalStats.Tabs)
	},
}

func init() {
	rootCmd.AddCommand(lncCmd)

	// Here you will define your flags and configuration settings.
	lncCmd.Flags().IntP("threads", "j", 4, "Number of cores to use")
	lncCmd.Flags().IntP("batches", "b", 64, "Number of batches to use")

}
