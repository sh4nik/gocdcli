package client

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// ListPipes fetches a list of pipelines with group names
func ListPipes() {
	conf := getConfig()
	for _, pipelineGroup := range conf.PipelineGroups {
		for _, pipeline := range pipelineGroup.Pipelines {
			printPipe(pipelineGroup.Group, pipeline.Name)
		}
	}
}

// ComparePipes compares two piplelines and prints the output
func ComparePipes(p1Name string, p2Name string, diffOnly bool) {
	conf := getConfig()
	p1 := getPipeline(conf, p1Name)
	p2 := getPipeline(conf, p2Name)
	comp := compareEnvVars(p1.EnvVars, p2.EnvVars)
	printDiff("Pipeline", comp, diffOnly)
	for _, stage1 := range p1.Stages {
		for _, stage2 := range p2.Stages {
			if stage1.Name == stage2.Name {
				comp := compareEnvVars(stage1.EnvVars, stage2.EnvVars)
				printDiff(stage1.Name, comp, diffOnly)
				break
			}
		}
	}
}

const colorOff = "\x1b[0;1m"
const colorBlack = "\x1b[30;1m"
const colorGreen = "\x1b[32;1m"
const colorBlue = "\x1b[94;1m"
const colorRed = "\x1b[31;1m"
const backgroundColorBlue = "\x1b[104;1m"

type client struct {
	baseURL  string
	user     string
	password string
	client   *http.Client
}

type config struct {
	XMLName        xml.Name        `xml:"cruise"`
	PipelineGroups []pipelineGroup `xml:"pipelines"`
}

type pipelineGroup struct {
	XMLName   xml.Name   `xml:"pipelines"`
	Group     string     `xml:"group,attr"`
	Pipelines []pipeline `xml:"pipeline"`
}

type pipeline struct {
	XMLName xml.Name `xml:"pipeline"`
	Name    string   `xml:"name,attr"`
	EnvVars envVars  `xml:"environmentvariables"`
	Stages  []stage  `xml:"stage"`
}

type envVars struct {
	XMLName   xml.Name   `xml:"environmentvariables"`
	Variables []variable `xml:"variable"`
}

type stage struct {
	XMLName xml.Name `xml:"stage"`
	Name    string   `xml:"name,attr"`
	EnvVars envVars  `xml:"environmentvariables"`
}

type variable struct {
	XMLName xml.Name `xml:"variable"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value"`
}

type varComp struct {
	Name   string
	Value1 string
	Value2 string
}

func getConfig() config {
	client := newClient()
	conf, err := client.fetchConfig()
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func newClient() *client {
	return &client{
		baseURL:  viper.GetString("url"),
		user:     viper.GetString("username"),
		password: viper.GetString("password"),
		client:   &http.Client{},
	}
}

func (c *client) fetchConfig() (config, error) {

	var conf config

	req, err := http.NewRequest("GET", c.baseURL+"admin/config_xml", nil)
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.client.Do(req)

	if err != nil {
		return conf, errors.New("Could not reach server")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return conf, errors.New("Error reading server response")
	}

	if resp.StatusCode != http.StatusOK {
		return conf, errors.New("Server communication error")
	}

	xml.Unmarshal(body, &conf)

	return conf, nil
}

func getPipeline(conf config, name string) pipeline {
	var pipe pipeline
	for _, pipelineGroup := range conf.PipelineGroups {
		for _, pipe := range pipelineGroup.Pipelines {
			if name == pipe.Name {
				return pipe
			}
		}
	}
	return pipe
}

func compareEnvVars(vars1 envVars, vars2 envVars) []varComp {
	var diff []varComp
	for _, var1 := range vars1.Variables {
		found := false
		for _, var2 := range vars2.Variables {
			if var1.Name == var2.Name {
				diff = append(diff, varComp{var1.Name, var1.Value, var2.Value})
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, varComp{var1.Name, var1.Value, "MISSING"})
		}
	}
	for _, var2 := range vars2.Variables {
		found := false
		for _, var1 := range vars1.Variables {
			if var2.Name == var1.Name {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, varComp{var2.Name, "MISSING", var2.Value})
		}
	}
	return diff
}

func printPipe(group string, name string) {
	fmt.Println(colorBlue + group + ": " + colorOff + name)
}

func printDiff(header string, comp []varComp, diffOnly bool) {
	fmt.Println()
	fmt.Println(backgroundColorBlue + colorBlack + "       " + colorOff + "  " + header)
	fmt.Println()
	for _, variable := range comp {
		if variable.Value1 == variable.Value2 {
			if !diffOnly {
				fmt.Println(colorGreen + "[MATCH] " + colorBlue + variable.Name + ": " + colorOff + variable.Value1)
			}
		} else {
			fmt.Println(colorRed + "[DIFF]  " + colorBlue + variable.Name + ": " + colorOff + variable.Value1 + " -> " + variable.Value2)
		}
	}
}
