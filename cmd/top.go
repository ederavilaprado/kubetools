// Copyright © 2016 Eder Ávila Prado <eder.prado@luizalabs.com>
//

package cmd

import "github.com/spf13/cobra"

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "A brief description of your command",
}

func init() {
	RootCmd.AddCommand(topCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// topCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// topCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
