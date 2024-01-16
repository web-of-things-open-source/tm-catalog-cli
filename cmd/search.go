/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"strings"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/blevesearch/bleve/v2/search/query"

	// _ "github.com/blevesearch/bleve/v2/search/highlight"
	// _ "github.com/blevesearch/bleve/v2/search/highlight/highlighter/ansi"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search",
	Long:  `search.`,
	RunE:  executeSearch,
}

var limit, skip, repeat int
var explain, highlight, fields bool
var qtype, qfield, sortby string
var idx bleve.Index

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "Repeat the query this many times.")
	searchCmd.Flags().IntVarP(&limit, "number limit", "n", 10, "Limit number of results returned.")
	searchCmd.Flags().IntVarP(&skip, "skip", "s", 0, "Skip the first N results.")
	searchCmd.Flags().BoolVarP(&explain, "explain", "x", false, "Explain the result scoring.")
	searchCmd.Flags().BoolVar(&highlight, "highlight", true, "Highlight matching text in results.")
	searchCmd.Flags().BoolVar(&fields, "fields", false, "Load stored fields.")
	searchCmd.Flags().StringVarP(&qtype, "type", "t", "query_string", "Type of query to run.")
	searchCmd.Flags().StringVarP(&qfield, "field", "f", "", "Restrict query to field, not applicable to query_string queries.")
	searchCmd.Flags().StringVarP(&sortby, "sort-by", "b", "", "Sort by field.")
}

func executeSearch(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must specify query")
	}
	var errOpen error
	idx, errOpen = bleve.Open("../catalog.bleve")
	if errOpen != nil {
		return fmt.Errorf("error opening bleve index: %v", errOpen)
	}
	defer idx.Close()
	query := buildQuery(args)
	for i := 0; i < repeat; i++ {
		req := bleve.NewSearchRequestOptions(query, limit, skip, explain)
		if highlight {
			req.Highlight = bleve.NewHighlightWithStyle("ansi")
		}
		if fields {
			req.Fields = []string{"*"}
		}
		if sortby != "" {
			if strings.Contains(sortby, ",") {
				req.SortBy(strings.Split(sortby, ","))
			} else {
				req.SortBy([]string{sortby})
			}
		}
		res, err := idx.Search(req)
		if err != nil {
			return fmt.Errorf("error running query: %v", err)
		}
		fmt.Println(res)
	}
	return nil

}

func buildQuery(args []string) query.Query {
	var q query.Query
	switch qtype {
	case "prefix":
		pquery := bleve.NewPrefixQuery(strings.Join(args, " "))
		if qfield != "" {
			pquery.SetField(qfield)
		}
		q = pquery
	case "term":
		pquery := bleve.NewTermQuery(strings.Join(args, " "))
		if qfield != "" {
			pquery.SetField(qfield)
		}
		q = pquery
	case "wild":
		pquery := bleve.NewWildcardQuery(strings.Join(args, " "))
		if qfield != "" {
			pquery.SetField(qfield)
		}
		q = pquery
	default:
		// build a search with the provided parameters
		queryString := strings.Join(args, " ")
		q = bleve.NewQueryStringQuery(queryString)
	}
	return q
}
