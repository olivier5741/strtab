Golang library to manipulate a table of string with optional row and/or column header :

- append from `[]string...` or `maps[string]map[string]string`
- set row or column header 
- transpose 
- render without / with row / with column / with both header(s) (optional empty top left string)
- string sort according to row header (can transpose to sort on both headers)

Useful before marshalling to csv
