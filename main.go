package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/Axili39/yyt/godict"
	"gopkg.in/yaml.v2"
)

// yyt features:
// + Merging multiple yaml files
// + extract data node from path (output string, yaml, json)
// + apply transformation base on template

// CommandLine Options
// yyt [-h] [-v] [-version] [-node PATH] [-exec TEMPLATE] [-execfile FILENAME] [-override] [-prio last/first] FILES
// -h show usage
// -v verbose mode
// -version show version
// -node PATH : select a node describe by PATH
// -exec TEMPLATE : execute template on selected node. By default selected node is root document
// -execfile FILENAME : execute template on selected node. By default, selected node is root document
// -override : (NYI) Not Yet Implemented
// -prio : replace priority order (NYI)
// -out	FILENAME	: output result to filename
// FILES : list of files. If empty, take stdin stream
//
// output:
// if -exec or -execfile options are empty, yyt marshall selected node to Yaml encoding format to output.

func runTemplate(node interface{}, tmpl string, output io.Writer) {
	generator := template.Must(template.New("").Parse(tmpl))
	err := generator.Execute(output, node)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template", err)
		os.Exit(1)
	}
}

func loadFiles(files []string) godict.Dict {
	var dict godict.Dict
	var err error

	if len(files) == 0 {
		// TODO load stdin
	} else {
		dict, err = godict.LoadFromYamlFiles(files)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading files :", err)
			os.Exit(1)
		}
	}
	return dict
}

func selectNode(nodePath string, dict godict.Dict) interface{} {
	// by default output node to yaml
	var node interface{} = dict

	if nodePath != "" {
		var err error
		node, err = dict.ExtractFromXPath(nodePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if node == nil {
		// TODO show Error
		os.Exit(1)
	}
	return node
}

func processTemplate(node interface{}, exec string, execFile string, output io.Writer) {
	if exec != "" {
		runTemplate(node, exec, output)
		os.Exit(0)
	}

	if execFile != "" {
		//TODO
		os.Exit(0)
	}
}

func main() {
	exec := flag.String("exec", "", "execute template on node")
	execFile := flag.String("execfile", "", "execute template on node")
	nodePath := flag.String("node", "", "path to node to do action")
	out := flag.String("out", "", "output file")
	showVersion := flag.Bool("version", false, "Show version")
	outformat := flag.String("out-format", "yaml", "Output format")

	flag.Parse()
	files := flag.Args()

	// Schow version
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Load Files
	dict := loadFiles(files)

	// Select Node
	node := selectNode(*nodePath, dict)

	// Set output
	output := os.Stdout
	if *out != "" {
		output, _ = os.Create(*out)
	}

	// Process Template if necessary
	processTemplate(node, *exec, *execFile, output)

	// Marshal to yaml if necessary
	var data []byte
	switch *outformat {
	case "yaml":
		data, _ = yaml.Marshal(node)
		_, err := output.WriteString("---\n")
		if err != nil {
			fmt.Fprintln(os.Stderr, "i/o error")
			os.Exit(1)
		}
	case "json":
		data, _ = json.Marshal(godict.Y2JConvert(node))
	}

	_, err := output.Write(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "i/o error")
		os.Exit(1)
	}
}
