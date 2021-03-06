package format

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/axialmarket/releaser/notes"
)

func ReadRelNotesVars(r io.Reader, vars *notes.RelNotesVars) error {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, vars)
}

func CompileRelNotesTemplate(vars notes.RelNotesVars, tmplFile io.Reader, w io.Writer) error {
	contents, err := ioutil.ReadAll(tmplFile)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("relnotes").Parse(fmt.Sprintf("%s", contents)))
	return tmpl.Execute(w, vars)
}
