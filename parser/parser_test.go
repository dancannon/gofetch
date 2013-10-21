package parser

import (
	"testing"
)

const str = `
    <div id="header">
        <h1>Header</h1>
    </div>
    <div id="content">
        <p>Test</p>
    </div>
    <div id="footer">
        Footer
    </div>
`

func TestRequest(t *testing.T) {
	NewParser().Parse(str)
}
