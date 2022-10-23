package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ObjExists(m map[string]interface{}, key string) string {
	if _, ok := m[key]; !ok {
		return key
	} else {
		return ObjExists(m, key+"-Copy")
	}
}

func PrettyPrint(m map[string]interface{}) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func GetCleanedStatement(s string) []string {
	if strings.HasPrefix(s, "#") {
		return []string{""}
	}
	var ss []string
	for _, v := range strings.SplitAfterN(strings.TrimSpace(strings.ReplaceAll(s, "\t", "    ")), " ", 2) {
		if v != "" {
			ss = append(ss, Dequote(strings.Trim(v, " ")))
		}
	}
	return ss
}

func Dequote(s string) string {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return s[1 : len(s)-1]
	} else {
		return s
	}
}

func GenerateMap(path string) map[string]interface{} {
	m := make(map[string]interface{})
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("No File Found at path : " + path)
		os.Exit(-1)
	}
	a := strings.Split(string(content), "\n")
	var curr_obj string = ""
	obj_map := make(map[string]string)
	for k := range a {
		s := GetCleanedStatement(a[k])
		if len(s) == 2 && s[1] == "{" {
			curr_obj = ObjExists(m, s[0])
		} else if curr_obj != "" {
			if len(s) == 2 {
				obj_map[s[0]] = s[1]
			} else if len(s) == 1 && s[0] == "}" {
				m[curr_obj] = obj_map
				obj_map = make(map[string]string)
				curr_obj = ""
			}
		} else if len(s) == 2 {
			m[s[0]] = s[1]
		}
	}
	return m
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File Provided!")
		os.Exit(-1)
	} else {
		PrettyPrint(GenerateMap(os.Args[1]))
	}
}
