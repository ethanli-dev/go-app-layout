/*
Copyright © 2025 lixw
*/
package main

import (
	"os"
	"runtime/debug"

	"github.com/ethanli-dev/go-app-layout/cmd"
)

// @title Go App Layout API
// @version 1.0
// @description This is a sample server celler server.
func main() {
	setCrashOutput()
	cmd.Execute()
}

func setCrashOutput() {
	crashFile, _ := os.Create("crash.log")
	_ = debug.SetCrashOutput(crashFile, debug.CrashOptions{})
}
