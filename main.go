/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"rtools/cmd"

	"github.com/apache/rocketmq-client-go/v2/rlog"
)

func main() {
	rlog.SetLogLevel("error")
	cmd.Execute()
}
