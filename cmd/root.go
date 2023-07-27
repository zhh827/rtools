/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rtools",
	Short: "rocketMQ 命令测试工具",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// var nameServer string
	// consumerCmd.PersistentFlags().StringVarP(&nameServer, "name", "n", "", " name server addr")
	// var groupName string
	// consumerCmd.PersistentFlags().StringVar(&groupName, "group", "tester", "consumer group name")
	// var TopicName string
	// consumerCmd.PersistentFlags().StringVar(&TopicName, "topic", "test-topic", "topic name")
	// consumerCmd.MarkFlagRequired("name")

}
