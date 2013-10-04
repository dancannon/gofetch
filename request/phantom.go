package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

type PhantomRequest struct{}

const scriptTmpl string = `
    var page = require('webpage').create();

    // Setup PhantomJS
    page.settings.userAgent = '{{.userAgent}}';

    // Make request
    page.open('{{.url}}', function(status) {
        if (status !== 'success') {
            console.log(JSON.stringify({
                "error": "Unable to load page"
            }));
            phantom.exit();
        } else {
            window.setInterval(function() {
                console.log(JSON.stringify(page));
                phantom.exit();
            }, {{.waitTime}});
        }
    })
`

func (r *PhantomRequest) Send(url string) (string, error) {
	f, err := ioutil.TempFile("", "gofetch")
	if err != nil {
		return "", err
	}

	defer f.Close()
	defer os.Remove(f.Name())

	// Write script to temp file
	tmpl, err := template.New("phantom_script").Parse(scriptTmpl)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(f, map[string]string{
		"url":       url,
		"userAgent": "GoFetch",
		"waitTime":  "3000",
	})
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	// Execute script and return the result
	cmd := exec.Command("phantomjs", f.Name())
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	// Decode and return response
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		return "", err
	}

	if err, ok := response["error"]; ok {
		return "", fmt.Errorf(err.(string))
	}

	if _, ok := response["content"]; !ok {
		return "", fmt.Errorf("Unexpected response")
	}

	return response["content"].(string), nil
}
