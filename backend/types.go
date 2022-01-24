package main

import (
	"fmt"
	"io"
)

type SSE struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Data  string `json:"data"`
	Retry int    `json:"retry"`
}

func (sse *SSE) String() string {
	return fmt.Sprintf(`event: %s
data: %s
id: %s
retry: %d
`, sse.Event, sse.Data, sse.ID, sse.Retry)
}

func (sse *SSE) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintln(w, sse.String())
	if err != nil {
		return int64(n), err
	}
	return int64(n), nil
}

type SSEs []SSE

func (sses *SSEs) WriteTo(w io.Writer) (int64, error) {
	var total int64
	for _, sse := range *sses {
		n, err := sse.WriteTo(w)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}
