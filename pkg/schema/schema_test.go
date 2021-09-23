package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana/pkg/schema/rtinstance"
	"golang.org/x/tools/txtar"
)

var CasesDir = filepath.Join("testdata", "trimapplydefaults")

type Case struct {
	Name    string
	CUE     string
	Full    string
	Trimmed string
}

func TestGenerate(t *testing.T) {
	cases, err := loadCases(CasesDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		t.Run(c.Name+" apply default value", func(t *testing.T) {
			scmInstance, err := rtinstance.Rt.Compile(c.Name+".cue", c.CUE)
			if err != nil {
				t.Fatal(err)
			}
			inputResource := Resource{Value: c.Trimmed}
			out, err := ApplyDefaults(inputResource, scmInstance.Value())
			if err != nil {
				t.Fatal(err)
			}
			b := []byte(out.Value.(string))
			b, _ = JsonRemarshal(b)

			if s := cmp.Diff(c.Full, string(b)); s != "" {
				t.Fatal(s)
			}
		})
	}

	for _, c := range cases {
		t.Run(c.Name+" trim default value", func(t *testing.T) {
			scmInstance, err := rtinstance.Rt.Compile(c.Name+".cue", c.CUE)
			if err != nil {
				t.Fatal(err)
			}
			inputResource := Resource{Value: c.Full}
			out, err := TrimDefaults(inputResource, scmInstance.Value())
			if err != nil {
				t.Fatal(err)
			}
			b := []byte(out.Value.(string))
			b, _ = JsonRemarshal(b)

			if s := cmp.Diff(c.Trimmed, string(b)); s != "" {
				t.Fatal(s)
			}
		})
	}
}

func JsonRemarshal(bytes []byte) ([]byte, error) {
	var ifce interface{}
	err := json.Unmarshal(bytes, &ifce)
	if err != nil {
		return []byte{}, err
	}
	output, err := json.Marshal(ifce)
	outputstring := string(output)
	if err != nil {
		return []byte{}, err
	}
	outputstring = strings.Replace(outputstring, "\\u003c", "<", -1)
	outputstring = strings.Replace(outputstring, "\\u003e", ">", -1)
	outputstring = strings.Replace(outputstring, "\\u0026", "&", -1)
	return []byte(outputstring), nil
}

func loadCases(dir string) ([]Case, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var cases []Case
	for _, fi := range files {
		file := filepath.Join(dir, fi.Name())
		a, err := txtar.ParseFile(file)
		if err != nil {
			return nil, err
		}

		if len(a.Files) != 3 {
			return nil, fmt.Errorf("Malformed test case '%s': Must contain exactly three files (CUE, Full and Trimed), but has %d", file, len(a.Files))
		}

		fullBuffer := new(bytes.Buffer)
		fullJson := a.Files[1].Data
		if err := json.Compact(fullBuffer, fullJson); err != nil {
			return nil, err
		}

		trimBuffer := new(bytes.Buffer)
		trimedJson := a.Files[2].Data
		if err := json.Compact(trimBuffer, trimedJson); err != nil {
			return nil, err
		}

		cases = append(cases, Case{
			Name:    fi.Name(),
			CUE:     string(a.Files[0].Data),
			Full:    fullBuffer.String(),
			Trimmed: trimBuffer.String(),
		})
	}
	return cases, nil
}
