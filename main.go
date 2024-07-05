package main

import (
	"fyksz/cmd"
	"fyksz/k8s"
	"github.com/alecthomas/kong"
	"log"
)

type cli struct {
	Import     cmd.Import `cmd:"" usage:"Read one file or fill directory, and print out file content with --- separated (fyksz stream)"`
	Apply      cmd.Apply  `cmd:"" usage:"Execute a subprocess, standard input will be one fragment of the full input. If return 0, content will replace the original one."`
	Save       cmd.Save   `cmd:"" usage:"Save each stream element to a separated file"`
	Kubernetes k8s.K8s    `cmd:"" usage:"helpers to transform K8s resources files" aliases:"k8s"`
	Run        cmd.Run    `cmd:"" usage:"run steps from build.fyksz file"`
	Repeat     cmd.Repeat `cmd:"" usage:"repeat existing templates multiple times"`
}

func main() {
	ktx := kong.Parse(&cli{})
	err := ktx.Run()
	if err != nil {
		log.Fatalf("%++v", err)
	}
}
