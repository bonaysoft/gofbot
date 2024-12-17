package messenger

import (
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

const stringContainNewline = `aaa
bbb
ccc
`

func TestDefaultManager_BuildReply(t *testing.T) {
	fmt.Println(strconv.Quote(stringContainNewline)[1 : len(stringContainNewline)-1])
	tmpl, err := template.New("").Parse(stringContainNewline)
	assert.NoError(t, err)

	tmpl.Execute(nil, nil)
}
