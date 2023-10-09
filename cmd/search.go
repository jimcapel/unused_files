/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type file struct {
	path string
	name string
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search a directory for any unused files",
	Long:  `Recursively searches a directory for any files that are not imported by any other files. Specifically for javascript/typescript projects.`,
	Run: func(cmd *cobra.Command, args []string) {

		performanceFlagSet, _ := cmd.Flags().GetBool("performance")
		performanceFlagSetShorthand, _ := cmd.Flags().GetBool("p")

		start := time.Now()

		var files []file
		var usedFiles []file

		fileExtensionsToInclude := [4]string{".ts", ".tsx", ".js", ".jsx"}

		//	find a good value to use
		maxCapacity := 1000000

		root := "./"

		if len(args) >= 1 && args[0] != "" {
			root = args[0]
		}

		fmt.Printf("Search with root @ %s\n", root)
		//	first, we walk the tree and append all file paths & names (without the extension, as extension isn't needed for js import) to a slice.
		//	we also exclude any files with extensions not matching .ts/.tsx/.js/.jsx
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err == nil {
				if !d.Type().IsDir() && slices.Contains(fileExtensionsToInclude[:], filepath.Ext(d.Name())) {
					files = append(files, file{path: path, name: d.Name()[:len(d.Name())-len(filepath.Ext(d.Name()))]})
				}

			} else {
				println(err.Error())
			}
			return nil
		})

		/*
			We iterate over the files, this time opening each file and checking if any of the other file names are present

			If a fileName is present, we append the files path & name to usedFiles.
		*/

		for i := 0; i < len(files); i++ {
			file, err := os.Open(fmt.Sprintf("./%s", files[i].path))
			fmt.Print("\033[G\033[K")
			fmt.Printf("Searching: %s", files[i].name)

			if err == nil {
				// optionally, resize scanner's capacity for lines over 64K, needs to be adjusted for optimum capacity
				scanner := bufio.NewScanner(file)
				buf := make([]byte, maxCapacity)
				scanner.Buffer(buf, maxCapacity)

				for scanner.Scan() {

					words := strings.Split(scanner.Text(), " ")

					/*	if a line contains the words const/export/type, we can be certain that we have passed the imports
						and therefore can stop looking in this particular file
					*/
					if slices.Contains(words, "const") || slices.Contains(words, "export") || slices.Contains(words, "type") || slices.Contains(words, "enum") {
						break
					}

					//	Search for the filename
					for j := 0; j < len(files); j++ {
						//	We don't want to search a file for it's own name, or if the file has already been found
						if j == i || slices.Contains(usedFiles, files[j]) {
							continue
						}
						if strings.Contains(scanner.Text(), files[j].name) {
							usedFiles = append(usedFiles, files[j])
						}
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}

			} else {
				println(err.Error())
			}

			file.Close()
		}

		fmt.Printf("File name, Path\n")

		// we now compare all file names against used file names, to find unused file names
		for i := 0; i < len(files); i++ {

			if !slices.Contains(usedFiles, files[i]) {
				fmt.Printf("%s, %s\n", files[i].name, files[i].path)
			}

		}

		if performanceFlagSet || performanceFlagSetShorthand {
			elapsed := time.Since(start)
			log.Printf("Binomial took %s", elapsed)

		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().Bool("performance", false, "Used to report the run-time of the search command")

	searchCmd.PersistentFlags().Bool("p", false, "shorthand for performance flag, which is used to report the run-time of the search command")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
