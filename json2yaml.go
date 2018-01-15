package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	err = json.Unmarshal(input, &v)
	if err != nil {
		log.Fatal("Input parse error:", err)
	}

	outName := strings.TrimSuffix(inputName, ".json")
	outName = outName + ".yaml"
	out, err := os.Create(outName)
	if err != nil {
		log.Fatal("Output file error:", err)
	}
	defer out.Close()

	indent := NewIndenter("    ")
	PrintValue(out, indent, false, v)
}

func PrintValue(w io.Writer, i Indenter, byMap bool, v interface{}) {
	switch v := v.(type) {
	case nil:
		fmt.Fprint(w, "nil")
	case float64:
		fmt.Fprintf(w, "%f", v)
	case bool:
		fmt.Fprintf(w, "%v", v)
	case string:
		PrintString(w, i, v)
	case map[string]interface{}:
		PrintMap(w, i, byMap, v)
	case []interface{}:
		PrintSlice(w, i, byMap, v)
	default:
		fmt.Fprintf(w, "{{{UNKNOWN| %v <%t> }}}}", v, v)
	}
}

func PrintMap(w io.Writer, i Indenter, byMap bool, v map[string]interface{}) {
	if byMap {
		fmt.Fprint(w, "\n")
	}
	b := false
	for key, val := range v {
		if val == nil {
			continue
		}
		if b {
			fmt.Fprint(w, "\n")
		}
		fmt.Fprintf(w, `%s"%s": `, i.S, key)
		PrintValue(w, i.Increment(), true, val)
		b = true
	}
}

func PrintSlice(w io.Writer, i Indenter, byMap bool, v []interface{}) {
	if byMap {
		fmt.Fprint(w, "\n")
	}
	b := false
	for _, val := range v {
		if b {
			fmt.Fprint(w, "\n")
		}
		fmt.Fprintf(w, `%s - `, i.S)
		PrintValue(w, i.Increment(), true, val)
		b = true
	}
}

func PrintString(w io.Writer, i Indenter, v string) {
	if !strings.Contains(v, "\r\n") {
		fmt.Fprintf(w, "%q", v)
		return
	}

	vs := strings.Split(v, "\r\n")
	if len(vs) == 1 {
		fmt.Fprintf(w, "%q", v)
		return
	}

	if strings.HasSuffix(v, "\r\n") {
		fmt.Fprintln(w, "|")
	} else {
		fmt.Fprintln(w, "|+")
	}

	for _, vv := range vs {
		fmt.Fprintf(w, "%s%s\n", i.S, vv)
	}
}

type Indenter struct {
	S string
	I string
}

func NewIndenter(s string) Indenter {
	return Indenter{
		S: "",
		I: s,
	}
}

func (i Indenter) Increment() Indenter {
	return Indenter{
		S: i.S + i.I,
		I: i.I,
	}
}
