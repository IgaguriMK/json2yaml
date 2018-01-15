package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ghodss/yaml"
)

func main() {
	logf, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logf)
	log.SetFlags(log.Lshortfile)
	log.Println("---- launch ----")

	if len(os.Args) < 2 {
		log.Fatal("Too few arguments")
	}

	inputName := os.Args[1]
	input, err := ioutil.ReadFile(inputName)
	if err != nil {
		log.Fatal("Input file error:", err)
	}

	input = bytes.TrimPrefix(input, []byte("\xef\xbb\xbf"))

	var v interface{}
	err = yaml.Unmarshal(input, &v)
	if err != nil {
		log.Fatal("Input parse error:", err)
	}

	bs, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal("Output encoding error:", err)
	}

	outName := strings.TrimSuffix(inputName, ".yml")
	outName = strings.TrimSuffix(outName, ".yaml")
	outName = outName + ".json"
	out, err := os.Create(outName)
	if err != nil {
		log.Fatal("Output file error:", err)
	}
	defer out.Close()

	bs = bytes.Replace(bs, []byte(`\n`), []byte(`\r\n`), -1)

	out.Write(bs)
}
