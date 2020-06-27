package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/Axili39/yyt/godict"
	"gopkg.in/yaml.v2"
)

// yyt features:
// + Merging 2 yaml files
// + extract data node from path (output string, yaml, json)
// + apply transformation base on template

// Multiples file in command lines
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func runTemplate(node interface{}, tmpl string, output io.Writer) {
	generator := template.Must(template.New("").Parse(tmpl))
	err := generator.Execute(output, node)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template", err)
		os.Exit(1)
	}
}

// Commands :
//
// - execute template
//		-t "template string"
//		-tf "template file"
// - merge yaml
// 		-o "output filename"
// - getting node
//		-n path to node

func main() {
	var files arrayFlags
	flag.Var(&files, "f", "yaml file")
	template := flag.String("t", "", "execute template on node")
	templateFile := flag.String("tf", "", "execute template on node")
	nodePath := flag.String("n", "", "path to node to do action")
	out := flag.String("o", "", "output file")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	dict, _ := godict.LoadFromYamlFiles(files)
	// by default output node to yaml
	var node interface{} = dict
	var err error
	if *nodePath != "" {
		if *verbose {
			fmt.Println("Selecting NodePath : ", *nodePath)
		}
		node, err = dict.ExtractFromXPath(*nodePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if node == nil {
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Processing Node:", node)
	}
	output := os.Stdout
	if *out != "" {
		output, _ = os.Create(*out)
	}

	if *template != "" {
		runTemplate(node, *template, output)
		os.Exit(0)
	}

	if *templateFile != "" {
		//TODO
		os.Exit(0)
	}

	data, _ := yaml.Marshal(node)
	output.WriteString("---\n")
	output.Write(data)

}
