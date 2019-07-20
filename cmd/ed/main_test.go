package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var fSimple = `a
b
c`

type edTest struct {
	infile     string
	commands   []string
	wantOutput []string
	wantFile   []string
}

var tests = []edTest{
	{infile: fSimple,
		commands:   []string{"2,3p", "1,1p"},
		wantOutput: []string{"b", "c", "a"}},
	{infile: fSimple,
		commands: []string{"2,2d", "w"},
		wantFile: []string{"a", "c"}},
}

func runEd(infileContents string, commands []string) (out []string, fileContents string, err error) {

	// Create a temporary file to interact with
	f, err := ioutil.TempFile("", "edtest")
	if err != nil {
		return
	}
	err = f.Close()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(f.Name(), []byte(infileContents), 0644)
	if err != nil {
		return
	}
	defer os.Remove(f.Name())

	// Add "edit file" and "quit" commands to the input
	input := strings.NewReader("e " + f.Name() + "\n" + strings.Join(commands, "\n") + "\nq\n")
	var builder strings.Builder
	process(input, &builder, nil)
	out = strings.Split(builder.String(), "\n")

	// Remove first&last empty lines (from opening the file and closing the editor)
	out = out[1 : len(out)-1]
	fileContentsBytes, err := ioutil.ReadFile(f.Name())
	fileContents = string(fileContentsBytes)

	return out, fileContents, err
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMain(t *testing.T) {

	for _, test := range tests {
		gotOut, gotFile, err := runEd(test.infile, test.commands)
		if err != nil {
			t.Fatal(err)
		}
		if test.wantOutput != nil {
			if !stringSliceEqual(gotOut, test.wantOutput) {
				t.Errorf("output: got %q, wanted %q\n", gotOut, test.wantOutput)
			}
		}
		if test.wantFile != nil {
			if strings.Join(test.wantFile, "\n") != gotFile {
				t.Errorf("file: got %s, wanted %s\n", gotFile, test.wantFile)
			}

		}
	}
}
