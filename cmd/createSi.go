/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/commands"
	"github.com/web-of-things-open-source/tm-catalog-cli/internal/remotes"
)

func closeIndex(index bleve.Index) {
	var log = slog.Default()
	//log.Debug("pre  Close index")
	closeError := index.Close()
	log.Debug("post Close index", "error", closeError)
	// time.Sleep(10 * time.Second)
	// log.Debug("post Close index after sleep")
}

var catalogPath = "../catalog.bleve"

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
		searchResult, err, _ := listCmd.List(repoSpec, nil)
		//toc, err := spec.List(nil)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

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
		defer closeIndex(index)
		count := 0
		contents := searchResult.Entries
		for _, value := range contents {
			//			fmt.Printf("%s\t%s\n", value.Name, value.Mpn)

			// fn, err := commands.ParseFetchName(value.Name)
			// if err != nil {
			// 	log.Error(err.Error())
			// 	return //"", err
			// }

			for _, version := range value.Versions {
				//fn := &commands.FetchName{}
				//fqName := value.Name + ":" + version.Version.Model
				fqName := version.TMID
				//blID := fqName + "_" + version.Digest
				//blID = strings.ReplaceAll(blID, ":", "_")
				blID := fqName
				//fn, err := commands.ParseFetchName(fqName)
				if err != nil {
					log.Error(err.Error())
					return //"", err
				}
				id, thing, err, _ := commands.NewFetchCommand(rm).FetchByTMID(repoSpec, fqName, false)
				// thing_copy := make([]byte, len(thing))
				// copy(thing_copy[:], thing)
				_ = id
				if err != nil {
					fmt.Println(err.Error())
					//os.Exit(1)
					continue
				}
				saveSnapshot(thing, blID, "afterfetch-"+remoteName)

				// ask if Document is already indexed
				doc, _ := index.Document(blID)

				if doc != nil {
					deleteErr := index.Delete(blID)
					fmt.Printf("deleted exisiting document with id=%s first%v\n", blID, deleteErr)
				} else {
					fmt.Printf("new document with id=%s\n", blID)
				}

				var data any
				unmErr := json.Unmarshal(thing, &data)
				if unmErr != nil {
					fmt.Println(unmErr.Error())
					os.Exit(1)
				}
				saveSnapshot(thing, blID, "afterunmarshal-"+remoteName)

				idxErr := index.Index(blID, data)

				if idxErr != nil {
					fmt.Println(idxErr.Error())
					return
				}
				count++
				if count > 50000 {
					return
				}
			}
		}
	},
}

func saveSnapshot(thing []byte, blID string, sName string) {
	filename := "../snapshot-" + sName + "/" + blID
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = os.WriteFile(filename, thing, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type visitField func(parent any, data any, path string) (interface{}, error)

func RangeJSON(parent any, data any, path string, vf visitField) (any, error) {
	// if data == nil || strings.HasSuffix(path, ".forms") {
	// 	return nil, nil
	// }
	var err error
	//hideField := strings.HasSuffix(path, ".properties")
	// if data is a map, walk deeper in the fields of the map
	if aMap, isMap := data.(map[string]interface{}); isMap {

		for key, val := range aMap {
			var err2 error
			var val2 any
			if path == "" {
				val2, err2 = RangeJSON(aMap, val, key, vf)
			} else {
				val2, err2 = RangeJSON(aMap, val, path+"."+key, vf)
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
	// if data is a array, walk deeper in the each element of the array
	if aArr, isArr := data.([]any); isArr {

		j := 0
		for i := range aArr {
			val2, err2 := RangeJSON(nil, aArr[i], path+".["+strconv.Itoa(i)+"]", vf)
			err = ErrorCoalesce(err, err2)
			if val2 != nil {
				aArr[j] = val2
				j++
			}
		}
		return aArr[:j], err
	}
	// its a literal, so call the visitField function
	vf(parent, data, path)
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
