package forjava

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	normal = iota
	object
	list
)

func toUpper(s string, first bool) string {
	rst := strings.Split(s, "_")
	buf := strings.Builder{}
	for i, v := range rst {
		if i == 0 && first {
			buf.WriteString(v)
			continue
		}
		if 'a' <= v[0] && v[0] <= 'z' {
			buf.WriteByte(v[0] - 'a' + 'A')
			buf.WriteString(v[1:])
		} else {
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func deduceNormal(tab, k string, v interface{}) (string, string) {
	switch v.(type) {
	case string, byte, []byte:
		return fmt.Sprintf("%sprivate String %s;\n", tab, string(k)), "String"
	case bool:
		return fmt.Sprintf("%sprivate Boolean %s;\n", tab, string(k)), "Boolean"
	case float64, float32:
		return fmt.Sprintf("%sprivate Double %s;\n", tab, string(k)), "Double"
	}
	return fmt.Sprintf("//unknown %s %v \n", string(k), v), "Object"
}

//return content and type and identify
func deduceItf(tab, name string, v interface{}) (content, rtype string, id int) {
	switch v.(type) {
	case []interface{}:
		content, rtype = deduceArr(tab, name, v.([]interface{}))
		id = list
	case map[string]interface{}:
		content, rtype = deduceMap(v.(map[string]interface{}), tab, name)
		id = object
	default:
		content, rtype = deduceNormal(tab, name, v)
		id = normal
	}
	return
}

func deduceArr(tab, name string, v []interface{}) (string, string) {
	r := strings.Builder{}
	val, itemType, id := deduceItf(tab, name, v[0])
	listType := "List<" + itemType + ">"
	//r.WriteString(fmt.Sprintf("%sprivate %s %s;\n", tab, listType, name))
	if id != normal{
		r.WriteString(val)
	}
	return r.String(), listType
}

func deduceMap(data map[string]interface{}, tab, name string) (string, string) {
	r := strings.Builder{}
	mapType := toUpper(name, false)
	vals := make([]string, 0)

	r.WriteString(tab + "@Getter @Setter\n")
	r.WriteString(fmt.Sprintf("%sprivate static class %s{\n", tab, mapType))
	for k, v := range data {
		val, rtype , id := deduceItf(tab+"\t", k, v)
		r.WriteString(fmt.Sprintf("%sprivate %s %s;\n", tab+"\t", rtype, k))
		if id != normal{
			vals = append(vals, val)
		}
	}
	for _, e := range vals {
		r.WriteString(e)
	}
	r.WriteString(tab + "}\n")
	return r.String(), mapType
}

func Json2java(str []byte, name string) (string, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal(str, &data)
	if err != nil {
		fmt.Println("err", err.Error())
		return "", err
	}
	if name == "" {
		name = "AutoGen"
	}
	content, _ ,_:= deduceItf("", name, data)
	return content, nil
}
