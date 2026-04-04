package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"

	"github.com/jedib0t/go-pretty/v6/table"
)

type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
)

type Output struct {
	Format OutputFormat
}

func New(format string) *Output {
	if format == "json" {
		return &Output{Format: FormatJSON}
	}
	return &Output{Format: FormatTable}
}

func (o *Output) Print(data interface{}) error {
	if o.Format == FormatJSON {
		return o.PrintJSON(data)
	}
	return o.PrintTable(data)
}

func (o *Output) PrintJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (o *Output) PrintTable(data interface{}) error {
	switch v := data.(type) {
	case []map[string]interface{}:
		return o.printSliceOfMaps(v)
	case map[string]interface{}:
		return o.printMap(v)
	case []interface{}:
		return o.printSliceInterface(v)
	case []string:
		for _, s := range v {
			fmt.Println(s)
		}
		return nil
	case string:
		fmt.Println(v)
		return nil
	default:
		slice := o.convertToSliceOfMaps(data)
		if slice != nil {
			return o.printSliceOfMaps(slice)
		}
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}
}

func (o *Output) convertToSliceOfMaps(data interface{}) []map[string]interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]map[string]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return nil
		}

		m := make(map[string]interface{})
		for j := 0; j < elem.NumField(); j++ {
			field := elem.Type().Field(j)
			tag := field.Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			m[tag] = elem.Field(j).Interface()
		}
		result[i] = m
	}

	return result
}

func (o *Output) printSliceInterface(data []interface{}) error {
	if len(data) == 0 {
		fmt.Println("No data")
		return nil
	}

	maps := make([]map[string]interface{}, 0, len(data))
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			maps = append(maps, m)
		}
	}

	if len(maps) > 0 {
		return o.printSliceOfMaps(maps)
	}

	return o.PrintJSON(data)
}

func (o *Output) printSliceOfMaps(data []map[string]interface{}) error {
	if len(data) == 0 {
		fmt.Println("No data")
		return nil
	}

	headers := make([]string, 0)
	for k := range data[0] {
		headers = append(headers, k)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	headerRow := table.Row{}
	for _, h := range headers {
		headerRow = append(headerRow, h)
	}
	t.AppendHeader(headerRow)

	for _, row := range data {
		values := make([]interface{}, len(headers))
		for i, h := range headers {
			values[i] = row[h]
		}
		t.AppendRow(values)
	}

	t.Render()
	return nil
}

func (o *Output) printMap(data map[string]interface{}) error {
	if len(data) == 0 {
		fmt.Println("No data")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")
	fmt.Fprintln(w, "---\t-----")
	for k, v := range data {
		fmt.Fprintf(w, "%s\t%v\n", k, v)
	}
	w.Flush()
	return nil
}

func (o *Output) PrintError(err error) {
	if o.Format == FormatJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]string{"error": err.Error()})
	} else {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	}
}

func (o *Output) PrintSuccess(msg string) {
	if o.Format == FormatJSON {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]string{"message": msg})
	} else {
		fmt.Println("✓", msg)
	}
}
