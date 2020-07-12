package main

import (
	"os"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/mattn/go-colorable"
	"github.com/nlowe/euchre/cmd"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	petname.NonDeterministicMode()

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
		FullTimestamp:   true,
	})

	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
