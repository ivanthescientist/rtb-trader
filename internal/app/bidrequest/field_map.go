package bidrequest

import (
	"fmt"
	"strings"
)

type FieldMap map[string]string

func (fm FieldMap) String() string {
	var builder = &strings.Builder{}

	_, _ = builder.WriteString("FieldMap {\n")

	for field, value := range fm {
		_, _ = fmt.Fprintf(builder, "\t %s = %s\n", field, value)
	}

	_, _ = builder.WriteString("}")

	return builder.String()
}
