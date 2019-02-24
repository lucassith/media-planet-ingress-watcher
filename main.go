package main

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/lucassith/kube-watch/kubectl"

	"os"
	"os/signal"
	"syscall"

	logger "github.com/apsdehal/go-logger"
	"github.com/fsnotify/fsnotify"
)

var outputDirectory string
var log *logger.Logger
var watchDirectory string

func init() {
	log, _ = logger.New("main", 1, os.Stdout)
	args := os.Args
	if len(args) != 3 {
		log.Criticalf("You must specify arguments:\n\t1. Directory to listen for changes and new files.\n\t2. Output of yamls.")
		os.Exit(1)
	}
	os.MkdirAll(args[1], os.ModePerm)
	os.MkdirAll(args[2], os.ModePerm)
	watchDirectory = args[1]
	outputDirectory = args[2]
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Unable to start file watcher, error: %v\n", err)
		os.Exit(2)
	}
	defer watcher.Close()
	watcher.Add(watchDirectory)
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	log.Infof("Application ready.\n\nListening for changes in: %s\n", outputDirectory)
	for {
		select {
		case _ = <-gracefulStop:
			{
				os.Exit(0)
			}
		case newFile := <-watcher.Events:
			{
				if newFile.Op != fsnotify.Write {
					continue
				}
				contents, err := ioutil.ReadFile(newFile.Name)
				if err != nil {
					log.Errorf("Unable to read file, %v\n", err)
					continue
				}
				match, err := regexp.Match("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$", contents)
				if !match || err != nil {
					log.Errorf("[] Hostname regexp did not match %s\nerror: %v\n", contents, err)
					continue
				}
				filepath, err := makeNewYaml(string(contents[:]))
				log.Noticef("Created new ingress for %s in %s\n", contents, filepath)
				if err != nil {
					log.Errorf("Unable to create new yaml file, error: %v\n", err.Error())
				}
				cmdOutput, err := kubectl.ExecuteKubectl(filepath)
				log.Infof("Executed kubectl - output:\n%s\n error: %v\n", cmdOutput, err)
			}
		}
	}
}

func makeNewYaml(hostname string) (string, error) {

	filename := path.Join(outputDirectory, strings.Replace(hostname, ".", "-", -1)+"-ingress.yaml")
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = kubectl.MakeIngressFile(hostname, file)
	if err != nil {
		return "", err
	}
	return filename, nil
}
