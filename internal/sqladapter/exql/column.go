package exql

import (
	"fmt"
	"strings"
)

type columnT struct {
	Name  string
	Alias string
}

// Column represents a SQL column.
type Column struct {
	Name interface{}
	hash hash
}

// ColumnWithName creates and returns a Column with the given name.
func ColumnWithName(name string) *Column {
	return &Column{Name: name}
}

// Hash returns a unique identifier for the struct.
func (c *Column) Hash() string {
	return c.hash.Hash(c)
}

// Compile transforms the ColumnValue into an equivalent SQL representation.
func (c *Column) Compile(layout *Template) (compiled string) {
	if z, ok := layout.Read(c); ok {
		return z
	}

	switch value := c.Name.(type) {
	case string:
		input := trimString(value)

		chunks := separateByAS(input)

		if len(chunks) == 1 {
			chunks = separateBySpace(input)
		}

		name := chunks[0]

		nameChunks := strings.SplitN(name, layout.ColumnSeparator, 2)

		for i := range nameChunks {
			nameChunks[i] = trimString(nameChunks[i])
			nameChunks[i] = mustParse(layout.IdentifierQuote, Raw{Value: nameChunks[i]})
		}

		name = strings.Join(nameChunks, layout.ColumnSeparator)

		var alias string

		if len(chunks) > 1 {
			alias = trimString(chunks[1])
			alias = mustParse(layout.IdentifierQuote, Raw{Value: alias})
		}

		compiled = mustParse(layout.ColumnAliasLayout, columnT{name, alias})
	case Raw:
		compiled = value.String()
	default:
		compiled = fmt.Sprintf("%v", c.Name)
	}

	layout.Write(c, compiled)

	return
}
