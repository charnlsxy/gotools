package forjava

import (
	"encoding/json"
	"fmt"
	"strings"
)

func printNormal(tab, k string, v interface{}) (string, string) {
	switch v.(type) {
	case string, byte, []byte:
		return fmt.Sprintf("%sprivate String %s;\n", tab, string(k)), "String"
	case bool:
		return fmt.Sprintf("%sprivate boolean %s;\n", tab, string(k)), "Boolean"
	case float64, float32:
		return fmt.Sprintf("%sprivate double %s;\n", tab, string(k)), "Double"
	}
	return fmt.Sprintf("//unknown %s %v \n", string(k), v), "Object"
}

func printItf(tab,k string, v interface{}) (string, string) {
	switch v.(type) {
	case string, byte, []byte, bool, float64, float32:
		return printNormal(tab,k, v)
	case []interface{}:
		return printArr(tab,k, v.([]interface{}))
	case map[string]interface{}:
		return printMap(v.(map[string]interface{}),tab)
	default:
		return fmt.Sprintf("//unknown %s %v \n", string(k), v), "Object"
	}
}

func printArr(tab,k string, v []interface{}) (string, string) {
	r := strings.Builder{}
	val, t := printItf(tab,k, v[0])
	t = "List<" + t + ">"
	r.WriteString(fmt.Sprintf("%sprivate %s %s;\n",tab, t, k))
	r.WriteString(val)
	return r.String(), t
}

var varCnt = 0
var tmpName = "VarClass_"

func getNewName() string {
	varCnt++
	return fmt.Sprintf("%s%d", tmpName, varCnt)
}
func printMap(data map[string]interface{}, tab string) (string, string) {
	r := strings.Builder{}
	name := getNewName()
	r.WriteString(tab+"@Getter\n")
	r.WriteString(tab+"@Setter\n")
	r.WriteString(fmt.Sprintf("%sstatic class %s{\n", tab, name))
	for k, v := range data {
		val, t := printItf(tab+"\t",k, v)
		if len(t) > len(tmpName) && t[:len(tmpName)] == tmpName {
			r.WriteString(fmt.Sprintf("%sprivate %s %s;\n", tab+"\t", t, k))
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
		val, t := printItf("\t",k, v)
		if len(t) > len(tmpName) && t[:len(tmpName)] == tmpName {
			r.WriteString(fmt.Sprintf("\tprivate %s %s;\n", t, k))
		}
		r.WriteString(val)
	}
	r.WriteString("}\n")
	return r.String()
}
func Json2java(str []byte,name string) (string, error) {
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
