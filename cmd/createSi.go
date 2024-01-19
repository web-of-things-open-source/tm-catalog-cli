/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/commands"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/remotes"
)

// createSiCmd represents the createSi command
var createSiCmd = &cobra.Command{
	Use:   "createSi",
	Short: "Creates or updates a search index",
	Long:  `Creates or updates a search index for all entries in the "Table of Contents"`,
	Run: func(cmd *cobra.Command, args []string) {
		var log = slog.Default()
		remoteName := cmd.Flag("remote").Value.String()
		dirName := cmd.Flag("directory").Value.String()
		repoSpec, err := remotes.NewSpec(remoteName, dirName)
		rm := remotes.DefaultManager()
		//		remote, err := rm.Get(remoteName)
		if err != nil {
			// TODO: error seems specific to remotes.Get()
			log.Error(fmt.Sprintf("could not initialize a remote instance for %s. check config", remoteName), "error", err)
			os.Exit(1)
		}
		listCmd := commands.NewListCommand(rm)
		searchResult, err := listCmd.List(repoSpec, nil)
		_ = searchResult
		//toc, err := spec.List(nil)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		// open a new index
		mapping := bleve.NewIndexMapping()
		index, err := bleve.New("../catalog.bleve", mapping)
		if err != nil {
			fmt.Println(err)
			return
		}
		contents := searchResult.Entries
		for _, value := range contents {
			fmt.Printf("%s\t%s\n", value.Name, value.Mpn)
			fmt.Println(string(value.Name))
			// fn, err := commands.ParseFetchName(value.Name)
			// if err != nil {
			// 	log.Error(err.Error())
			// 	return //"", err
			// }

			for _, version := range value.Versions {
				//fn := &commands.FetchName{}
				fqName := value.Name + ":" + version.Version.Model
				fn, err := commands.ParseFetchName(fqName)
				if err != nil {
					log.Error(err.Error())
					return //"", err
				}
				id, thing, err := commands.NewFetchCommand(rm).FetchByName(repoSpec, fn)
				//thing, err := commands.(fn, remote)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				fmt.Println(id)
				fmt.Println(string(thing))
				var data any
				json.Unmarshal(thing, &data)

				vf := func(data interface{}, path string) (interface{}, error) {
					fmt.Printf("%s: %v(%T)\n", path, data, data)
					return data, nil
				}

				RangeJSON(data, "", vf)
				index.Index(fqName, data)

			}

		}
	},
}

type visitField func(data interface{}, path string) (interface{}, error)

func RangeJSON(data interface{}, path string, vf visitField) (interface{}, error) {
	if data == nil || strings.HasSuffix(path, ".forms") {
		return nil, nil
	}
	var err error
	//hideField := strings.HasSuffix(path, ".properties")
	if aMap, isMap := data.(map[string]interface{}); isMap {

		for key, val := range aMap {
			var err2 error
			var val2 any
			if path == "" {
				val2, err2 = RangeJSON(val, key, vf)
			} else {
				val2, err2 = RangeJSON(val, path+"."+key, vf)
			}
			err = ErrorCoalesce(err, err2)
			if val2 == nil {
				delete(aMap, key)
			} else {
				aMap[key] = val2
			}
		}
		return aMap, err
	}

	if aArr, isArr := data.([]interface{}); isArr {

		j := 0
		for i := range aArr {
			val2, err2 := RangeJSON(aArr[i], path+".["+strconv.Itoa(i)+"]", vf)
			err = ErrorCoalesce(err, err2)
			if val2 != nil {
				aArr[j] = val2
				j++
			}
		}
		return aArr[:j], err
	}
	vf(data, path)
	return data, nil
}

func ErrorCoalesce(searchIn ...error) error {
	for _, err := range searchIn {
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(createSiCmd)
	createSiCmd.Flags().StringP("remote", "r", "", "name of the remote to pull from")
	createSiCmd.Flags().StringP("directory", "d", "", "TM repository directory to pull from")
}
