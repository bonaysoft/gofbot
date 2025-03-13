package messenger

import (
	"bytes"
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

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]any{})
}

func Test_flattenMap(t *testing.T) {
	type args struct {
		m        map[string]any
		maxDepth int
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "empty map",
			args: args{m: map[string]any{}, maxDepth: 0},
			want: map[string]string{},
		},
		{
			name: "flat map",
			args: args{m: map[string]any{"a": 1, "b": "2"}, maxDepth: 0},
			want: map[string]string{"a": "1", "b": "2"},
		},
		{
			name: "nested map",
			args: args{m: map[string]any{"a": map[string]any{"b": map[string]any{"c": 1, "d": true}}}, maxDepth: 3},
			want: map[string]string{"a.b.c": "1", "a.b.d": "true"},
		},
		{
			name: "nested map with array value",
			args: args{m: map[string]any{"a": map[string]any{"b": []any{1, 2, 3}}}, maxDepth: 3},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, flattenMap(tt.args.m, tt.args.maxDepth), "flattenMap(%v, %v)", tt.args.m, tt.args.maxDepth)
		})
	}
}
