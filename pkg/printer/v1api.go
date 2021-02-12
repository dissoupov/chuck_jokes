package printer

import (
	"fmt"
	"io"
	"time"

	v1 "github.com/dissoupov/chuck_jokes/api/v1"
	"github.com/olekukonko/tablewriter"
)

var (
	headerForClusterMember = []string{"ID", "Name", "URLs"}
	headerForNodeStatus    = []string{"Node", "Name", "Host", "Port", "Version", "Uptime", "Lease", "TTL"}
)

// PrintServerStatus will print v1.ServerStatus
func PrintServerStatus(w io.Writer, r *v1.ServerStatus) {
	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Append([]string{"Host", r.HostName})
	table.Append([]string{"Port", r.Port})
	table.Append([]string{"StartedAt", r.StartedAt.Format(time.RFC3339)})
	table.Append([]string{"Uptime", (r.Uptime / time.Second * time.Second).String()})
	table.Append([]string{"Version", r.Version})
	table.Render()
	fmt.Fprintln(w)
}

// PrintList prints a list of values
func PrintList(w io.Writer, header string, list []string) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{header})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, r := range list {
		table.Append([]string{r})
	}

	table.Render()
	fmt.Fprintln(w)
}
