// Package strtab provide useful manipulation on a table of string
// with optional headers for rows and colum.
// Use-case : before marhsalling to csv
package strtab

import (
	"fmt"
	"sort"
)

var (
	errRowSize = fmt.Errorf(
		"The size of the row doesn't match the size of content row")
	errColSize = fmt.Errorf(
		"The size of the col doesn't match the size of content col")
)

// T the table with optional row and column headers
type T struct {
	rowHeader []string
	colHeader []string
	content   [][]string
}

func (t T) Len() int { return len(t.content) }
func (t *T) Swap(i, j int) {
	t.content[i], t.content[j] = t.content[j], t.content[i]
	t.rowHeader[i], t.rowHeader[j] = t.rowHeader[j], t.rowHeader[i]
}
func (t T) Less(i, j int) bool { return t.rowHeader[i] < t.rowHeader[j] }

func (t T) String() (s string) {
	for _, l := range t.GetContentWithHeaders(true) {
		for _, v := range l {
			s += v + ","
		}
		s += "\n"
	}
	return
}

// NewEmpytT creates empty table
func NewEmptyT() *T {
	t := T{make([]string, 0), make([]string, 0), make([][]string, 0)}
	return &t
}

// NewT creates a table from r rows and h as column headers
// The first argument of each row will be the corresponding row header
func NewT(h []string, r ...[]string) *T {
	t := NewEmptyT()
	t.SetHeader(h)
	t.Append(r...)
	//sort.Sort(t)

	return t
}

// Transpose perform a matrix transformation of content
// and switch headers of the table
func (t *T) Transpose() *T {

	t.rowHeader, t.colHeader = t.colHeader, t.rowHeader

	if len(t.content) == 0 {
		return t
	}

	rows := len(t.content[0])
	cols := len(t.content)

	trans := make([][]string, rows)

	for i := 0; i < rows; i++ {
		trans[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			trans[i][j] = t.content[j][i]
		}
	}

	t.content = trans
	sort.Sort(t)
	return t
}

// SetHeader sets the column header
func (t *T) SetHeader(c []string) error {
	if len(t.content) != 0 && len(t.content) != len(c) {

		fmt.Println(len(t.content))
		return errRowSize
	}
	h := make([]string, len(c))
	copy(h, c)
	t.colHeader = h
	return nil
}

// NewTfromMap creates table from multi map
// and set headers according to map keys
func NewTfromMap(vs map[string]map[string]string) *T {
	h := make(map[string]int, 0)
	var c [][]string
	n := 0
	for i, r := range vs {
		newRow := make([]string, n+1)
		newRow[0] = i
		for j, v := range r {
			if e, ok := h[j]; ok {
				newRow[e+1] = v
			} else {
				for oldKey, cOld := range c {
					c[oldKey] = append(cOld, "")
				}
				newRow = append(newRow, v)
				h[j] = n
				n++
			}
		}
		c = append(c, newRow)
	}

	hSlice := make([]string, len(h))
	for k, v := range h {
		hSlice[v] = k
	}
	t := NewT(hSlice, c...)
	return t
}

// Append appends rows to the table and take the first argument
// of each of them as the corresponding row header
func (t *T) Append(vs ...[]string) error {
	var h []string
	c := make([][]string, len(vs))
	for i, v := range vs {
		if len(t.content) != 0 && len(t.content[0]) != len(v)-1 {
			return errRowSize
		}
		h = append(h, v[0])
		c[i] = v[1:]
	}

	t.content = append(t.content, c...)
	t.rowHeader = append(t.rowHeader, h...)
	if len(t.content) != 0 && len(t.colHeader) < len(t.content[0]) {
		t.colHeader = append(t.colHeader,
			make([]string, len(t.content[0])-len(t.colHeader))...)
	}
	return nil
}

// Sort perform a string sort on the row header
// and move the content rows accordingly
func (t *T) Sort() *T {
	sort.Sort(t)
	return t
}

// GetContent outputs the content of the table
func (t T) GetContent() [][]string {
	var out [][]string
	for _, r := range t.content {
		newRow := make([]string, len(r))
		copy(newRow, r)
		out = append(out, newRow)
	}
	return out
}

// GetContentWithColHeader outputs the content  of the table
// with only the column header
func (t T) GetContentWithColHeader() [][]string {
	var out [][]string
	out = append(out, t.colHeader)
	for _, c := range t.GetContent() {
		out = append(out, c)
	}
	return out
}

func prepend(base [][]string, add []string) [][]string {
	var out [][]string
	for i, l := range base {
		out = append(out, append([]string{add[i]}, l...))
	}
	return out
}

// GetContentWithRowHeader outputs the content of the table
// with only the row header
func (t T) GetContentWithRowHeader() [][]string {
	return prepend(t.content, t.rowHeader)
}

// GetContentWithHeaders outputs the content of the table with both headers
// If 'ok' is set to through add an empty element at the top left corner to make the
// column header fit the content
func (t T) GetContentWithHeaders(ok bool) [][]string {
	var out [][]string
	var colHead []string
	if ok {
		colHead = append(colHead, "")
	}
	colHead = append(colHead, t.colHeader...)
	// should put a 'if'
	if len(t.content) > 0 {
		out = append(out, colHead[:len(t.content[0])+1])
	}
	inter := prepend(t.content, t.rowHeader)

	for _, it := range inter {
		out = append(out, it)
	}

	return out
}
