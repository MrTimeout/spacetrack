/*
Copyright © 2022 MrTimeout estonoesmiputocorreo@gmail.com

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
	"time"

	"github.com/MrTimeout/spacetrack/client"
	"github.com/MrTimeout/spacetrack/data"
	"github.com/MrTimeout/spacetrack/model"
	"github.com/MrTimeout/spacetrack/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	limit      model.Limit
	format     model.Format
	orderBy    model.OrderBy
	predicates []string
	dryRun     bool
)

var gpCmd = &cobra.Command{
	Use:   "gp",
	Short: "Command which refers to the Request class GP or General Perturbations",
	Long: `Command which refers to the Request class GP or General Perturbations. We can fetch all the data from the satellite catalog and filter it.
We can limit, order and sort asceding or descending by any field present in the response. We can also format the response to 4 different ones and filter
response by a lot of fields.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, predicate := range predicates {
			if !model.IsOperandValid(predicate) {
				if help := model.OperandHelp(predicate); help != "" {
					cmd.Println(help)
				}
				return fmt.Errorf("trying to parse predicate argument: %s", predicate)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		defer func() {
			utils.Logger.Info("Time consumed from start of the application", zap.Duration("duration", time.Since(now)))
		}()

		p, err := model.ToPredicates(predicates)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}

		query := client.SpaceRequest{
			Limit:           limit,
			OrderBy:         orderBy,
			Format:          format,
			ShowEmptyResult: true,
			Predicates:      p,
		}.BuildQuery()

		if dryRun {
			cmd.Println("URL: ", query)
		} else {
			rsp, err := client.GetSpaceClientInstance().FetchData(query)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if err = data.Persist(workDir, rsp); err != nil {
				cmd.PrintErrln(err)
			}
		}
	},
	Example: `
	spacetrack gp --dry-run --format json --limit 10 --orderby norad_cat_id --sort asc

	spacetrack gp --format json --limit 10 --orderby norad_cat_id --sort desc --filter "decay_date<>null-val" --filter "epoch<now-30"

	spacetrack gp --format xml --orderby norad_cat_id --sort asc --filter "decay_date<>null-val" --log-level info --work-dir /tmp/my/spacetrack

	spacetrack gp --format xml --log-level debug --log-file /var/log/spacetrack.log --work-dir /tmp/my/spacetrack
	`,
}

func init() {
	gpCmd.Flags().Var(&format, "format", "Formatting output of the response. Possible values are html, json, csv, xml")
	gpCmd.Flags().Var(&orderBy.Sort, "sort", "Sort response Ascending or Descending. By default, it is asc")
	gpCmd.Flags().IntVar(&limit.Max, "limit", -1, "Limitting output to a restrictive number of results")
	gpCmd.Flags().IntVar(&limit.Skip, "skip", -1, "Skipping first n elements")
	gpCmd.Flags().StringVar(&orderBy.By, "orderby", "norad_cat_id", "Order results by specified field, which is present on the response. Default value is norad_cat_id. It is used in conjuction with sort, which default value is asc")
	gpCmd.Flags().StringArrayVar(&predicates, "filter", []string{}, "Filter response by all the fields allowed in the response. Default is none")
	gpCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Just build the path and prompt it to the console")

	utils.CheckErr(gpCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return model.FormatValues, cobra.ShellCompDirectiveDefault
	}))

	utils.CheckErr(gpCmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return model.SortValues, cobra.ShellCompDirectiveDefault
	}))

	utils.CheckErr(gpCmd.RegisterFlagCompletionFunc("orderby", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return model.ByPossibleValues, cobra.ShellCompDirectiveDefault
	}))

	utils.CheckErr(gpCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return model.PredicatePossibleValues, cobra.ShellCompDirectiveDefault
	}))

	rootCmd.AddCommand(gpCmd)
}
