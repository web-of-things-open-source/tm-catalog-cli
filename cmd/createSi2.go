/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/commands"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/remotes"
)

func checkFatal(err error, txt string) {
	if err != nil {
		slog.Default().Error(txt, "error", err)
		os.Exit(1)
	}
}

// createSiCmd represents the createSi command
var createSi2Cmd = &cobra.Command{
	Use:   "createSi2",
	Short: "Creates or updates a search index",
	Long:  `Creates or updates a search index for all entries in the "Table of Contents"`,
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := cmd.Flag("remote").Value.String()
		dirName := cmd.Flag("directory").Value.String()
		repoSpec, err := remotes.NewSpec(remoteName, dirName)
		checkFatal(err, fmt.Sprintf("could not initialize a remote instance for %s. check config", remoteName))

		rm := remotes.DefaultManager()
		listCmd := commands.NewListCommand(rm)
		searchResult, err, _ := listCmd.List(repoSpec, nil)
		checkFatal(err, "")

		// try to open index, if it not there create a fresh one
		index, err := bleve.Open(catalogPath)

		if err != nil {
			// open a new index
			indexMapping := bleve.NewIndexMapping()
			index, err = bleve.NewUsing(catalogPath, indexMapping, bleve.Config.DefaultIndexType, bleve.Config.DefaultKVStore, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		closeIndex(index)

		contents := searchResult.Entries
		for _, value := range contents {
			for _, version := range value.Versions {
				fqName := version.TMID
				// for now we use the relative filename as documentID
				blID := fqName
				id, thing, err, _ := commands.NewFetchCommand(rm).FetchByTMID(repoSpec, fqName, false)
				_ = id
				if err != nil {
					slog.Default().Error("cant Fetch TM", "error", err.Error())
					continue
				}
				var data any
				unmErr := json.Unmarshal(thing, &data)
				checkFatal(unmErr, "Unmarshal")

				// open index
				index, err := bleve.Open(catalogPath)
				checkFatal(err, "open Index")

				// ask if Document is already indexed
				doc, _ := index.Document(blID)
				if doc != nil {
					deleteErr := index.Delete(blID)
					fmt.Printf("deleted exisiting document with id=%s first%v\n", blID, deleteErr)
				} else {
					fmt.Printf("new document with id=%s\n", blID)
				}

				idxErr := index.Index(blID, data)
				checkFatal(idxErr, "index Document")
				err = index.Close()
				checkFatal(err, "close index")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(createSi2Cmd)
	createSi2Cmd.Flags().StringP("remote", "r", "", "name of the remote to pull from")
	createSi2Cmd.Flags().StringP("directory", "d", "", "TM repository directory to pull from")
}
