package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/ryanuber/columnize"
)

type Testsuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Properties Properties `xml:"properties"`
	Testcases  []Testcase `xml:"testcase"`
	Name       string     `xml:"name,attr"`
	Text       string     `xml:",chardata"`
	TestCount  int        `xml:"tests,attr"`
	Failures   int        `xml:"failures,attr"`
	Errors     int        `xml:"errors,attr"`
	Skipped    int        `xml:"skipped,attr"`
	Time       int        `xml:"time,attr"`
	Timestamp  string     `xml:"timestamp,attr"`
	Hostname   string     `xml:"hostname,attr"`
	SystemOut  string     `xml:"system-out"`
	SystemErr  string     `xml:"system-err"`
}

type Properties struct {
	Text     string     `xml:",chardata"`
	Property []Property `xml:"property"`
}

type Property struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Testcase struct {
	XMLName   xml.Name `xml:"testcase"`
	Text      string   `xml:",chardata"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      int      `xml:"time,attr"`
	Failure   Failure  `xml:"failure"`
	Skipped   string   `xml:"skipped"`
}

type Failure struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

var Version string

func report(filename string, showSkipped bool) {
	data, _ := ioutil.ReadFile(filename)
	testsuite := &Testsuite{}

	_ = xml.Unmarshal([]byte(data), &testsuite)

	success := testsuite.TestCount - testsuite.Failures - testsuite.Errors - testsuite.Skipped

	if success == testsuite.TestCount {
		return
	}

	output := []string{
		"Test Count | Failures | Errors | Skipped",
		"---------- | -------- | ------ | -------",
		fmt.Sprintf("%v | %v | %v | %v", testsuite.TestCount, testsuite.Failures, testsuite.Errors, testsuite.Skipped),
	}
	result := columnize.SimpleFormat(output)

	fmt.Println("-------------------------------------")
	fmt.Println(fmt.Sprintf("Test %s run @ %s", testsuite.Name, testsuite.Timestamp))
	fmt.Println(fmt.Sprintf("Hostname: %s", testsuite.Hostname))
	fmt.Println("-------------------------------------")
	fmt.Println(result)

	fmt.Println("\nFailures")
	for _, testcase := range testsuite.Testcases {
		if testcase.Failure.Text == "" {
			continue
		}

		fmt.Println("-------------------------------------")
		fmt.Println(fmt.Sprintf("- name: %s", testcase.Name))
		fmt.Println(fmt.Sprintf("  time: %vsec", testcase.Time))
		fmt.Println("  output:")
		fmt.Println(testcase.Failure.Text)
	}

	if !showSkipped {
		return
	}

	fmt.Println("\nSkipped Tests")
	fmt.Println("-------------------------------------")
	for _, testcase := range testsuite.Testcases {
		if testcase.Skipped == "" {
			continue
		}

		fmt.Println(fmt.Sprintf("- name: %s", testcase.Name))
		fmt.Println(fmt.Sprintf("  reason: %s", testcase.Skipped))
	}

	fmt.Println("\n")
}

func main() {
	parser := argparse.NewParser("xunit-reader", "Prints provided string to stdout")
	version := parser.Flag("v", "version", &argparse.Options{Help: "show version"})
	showSkipped := parser.Flag("s", "show-skipped", &argparse.Options{Help: "show skipped tests"})
	f := parser.String("p", "pattern", &argparse.Options{Required: true, Help: "pattern referencing files to process"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", parser.Usage(err))
		os.Exit(1)
		return
	}

	if *version {
		fmt.Printf("xunit-reader %v\n", Version)
		os.Exit(0)
		return
	}

	matches, _ := filepath.Glob(*f)
	for _, match := range matches {
		report(match, *showSkipped)
	}
}
