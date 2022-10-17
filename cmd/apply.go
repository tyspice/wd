package cmd

import (
	"bytes"
	"log"
	"os"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		apply(path)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func apply(path string) {
	patch, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	files, _, err := gitdiff.Parse(patch)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		file, err := os.OpenFile(f.OldName, os.O_RDWR, f.OldMode)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var output bytes.Buffer
		if err := gitdiff.Apply(&output, file, f); err != nil {
			log.Fatal(err)
		}

		file.Truncate(0)
		_, err = file.Write(output.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
}
