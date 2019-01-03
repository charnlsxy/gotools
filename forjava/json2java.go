package forjava

import (
	"encoding/json"
	"fmt"
	"strings"
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

func printNormal(tab, k string, v interface{}) (string, string) {
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

func printItf(tab, name string, v interface{}) (string, string) {
	switch v.(type) {
	case string, byte, []byte, bool, float64, float32:
		return printNormal(tab, name, v)
	case []interface{}:
		return printArr(tab, name, v.([]interface{}))
	case map[string]interface{}:
		return printMap(v.(map[string]interface{}), tab, name)
	default:
		return fmt.Sprintf("//unknown %s %v \n", string(name), v), "Object"
	}
}

func printArr(tab, name string, v []interface{}) (string, string) {
	r := strings.Builder{}
	val, t := printItf(tab, name, v[0])
	t = "List<" + toUpper(t, false) + ">"
	r.WriteString(fmt.Sprintf("%sprivate %s %s;\n", tab, t, name))
	r.WriteString(val)
	return r.String(), t
}

var varCnt = 0
var tmpName = "VarClass_"

func getNewName() string {
	varCnt++
	return fmt.Sprintf("%s%d", tmpName, varCnt)
}
func printMap(data map[string]interface{}, tab, name string) (string, string) {
	r := strings.Builder{}
	tp := toUpper(name, false)
	r.WriteString(fmt.Sprintf("%sprivate %s %s\n",tab, tp, name))
	r.WriteString(tab + "@Getter\n")
	r.WriteString(tab + "@Setter\n")
	r.WriteString(fmt.Sprintf("%sstatic class %s{\n", tab, tp))
	for k, v := range data {
		val, t := printItf(tab+"\t", k, v)
		if len(t) > len(tmpName) && t[:len(tmpName)] == tmpName {
			r.WriteString(fmt.Sprintf("%sprivate %s %s;\n", tab+"\t", toUpper(t, false), k))
		}
		r.WriteString(val)
	}
	r.WriteString(tab + "}\n")
	return r.String(), name
}
func build(name string, data map[string]interface{}) string {
	r := strings.Builder{}
	r.WriteString(fmt.Sprintf("@Data\npublic class %s{\n", name))
	for k, v := range data {
		val, t := printItf("\t", k, v)
		if len(t) > len(tmpName) && t[:len(tmpName)] == tmpName {
			r.WriteString(fmt.Sprintf("\tprivate %s %s;\n", t, k))
		}
		r.WriteString(val)
	}
	r.WriteString("}\n")
	return r.String()
}
func Json2java(str []byte, name string) (string, error) {
	varCnt = 0
	data := make(map[string]interface{})
	err := json.Unmarshal(str, &data)
	if err != nil {
		fmt.Println("err", err.Error())
		return "", err
	}
	if name == "" {
		name = getNewName()
	}
	return build(name, data), nil
}
