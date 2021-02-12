package printer_test

import (
	"bytes"
	"testing"

	v1 "github.com/dissoupov/chuck_jokes/api/v1"
	"github.com/dissoupov/chuck_jokes/pkg/printer"
	"github.com/stretchr/testify/assert"
)

const projFolder = "../../"

func Test_PrintServerStatus(t *testing.T) {
	w := bytes.NewBuffer([]byte{})
	printer.PrintServerStatus(w, &v1.ServerStatus{
		HostName: "host2",
		Port:     "123",
		Version:  "v1",
	})

	out := string(w.Bytes())
	assert.Equal(t,
		"  Host      | host2                 \n"+
			"  Port      | 123                   \n"+
			"  StartedAt | 0001-01-01T00:00:00Z  \n"+
			"  Uptime    | 0s                    \n"+
			"  Version   | v1                    \n\n",
		out)
}

func TestPrintList(t *testing.T) {
	w := bytes.NewBuffer([]byte{})

	printer.PrintList(w, "key", []string{"value"})

	out := string(w.Bytes())
	assert.Contains(t, out, "   KEY   \n---------\n  value  \n\n")
}
