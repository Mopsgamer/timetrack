package main

import (
	"fmt"
	"os"
	"timet/manager"
	"timet/timet"

	"github.com/spf13/cobra"
)

const (
	argName         = `name`
	argNewName      = `new name`
	descrIntegerGr0 = `integer (greater than 0)`
	descrBelow      = `add once from below, not from above`
	descrGlobFirst  = `glob pattern string. the first record on top. see https://pkg.go.dev/github.com/gobwas/glob#Compile`
	descrGlobSet    = `glob pattern string. set of records. see https://pkg.go.dev/github.com/gobwas/glob#Compile`
	descrDate       = `date time string. example: '` + timet.DateFormat + `'.`
)

var DefaultManager = manager.New(timet.PathData)

var RootCmd = &cobra.Command{
	Use:     "timet",
	Short:   "A command line tool for time monitoring",
	Long:    `Timet is a CLI tool for time records.`,
	Version: timet.Version,
}

func ProcessAction(message string, err error) {
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", message)
}

func main() {
	if DefaultManager.DataLoadFromFile() != nil {
		DefaultManager.Reset()
	}
	var ListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "display record list",
		Run: func(cmd *cobra.Command, args []string) {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			msg, err := DefaultManager.List(name)
			ProcessAction(msg, err)
		},
	}
	RootCmd.AddCommand(ListCmd)

	var CreateCmd = &cobra.Command{
		Use:     `create [` + argName + `]`,
		Aliases: []string{"new", "add"},
		Short:   "create new record",
		Run: func(cmd *cobra.Command, args []string) {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			date, _ := cmd.Flags().GetString("date")
			below, _ := cmd.Flags().GetBool("below")
			msg, err := DefaultManager.Create(name, date, below)
			ProcessAction(msg, err)
		},
	}
	CreateCmd.PersistentFlags().String("date", "", descrDate)
	CreateCmd.PersistentFlags().Bool("below", false, descrBelow)
	RootCmd.AddCommand(CreateCmd)

	var RemoveCmd = &cobra.Command{
		Use:     "remove <name>",
		Short:   "delete a couple of records",
		Aliases: []string{"rm"},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			msg, err := DefaultManager.Remove(name)
			ProcessAction(msg, err)
		},
	}
	RootCmd.AddCommand(RemoveCmd)

	var ResetCmd = &cobra.Command{
		Use:   "reset",
		Short: "reset all records and settings",
		Run: func(cmd *cobra.Command, args []string) {
			msg, err := DefaultManager.Reset()
			ProcessAction(msg, err)
		},
	}
	RootCmd.AddCommand(ResetCmd)

	ProcessAction("", RootCmd.Execute())
}
