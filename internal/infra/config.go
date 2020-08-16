package infra

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// LoadConf load env files according dir string passed as argument and set them into a map
func LoadConf(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var configs []map[string]string
	for _, f := range files {
		fn, _ := loadFile(dir + f.Name())
		configs = append(configs, fn)
	}

	return merge(configs...), nil
}

func merge(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func loadFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	m := make(map[string]string)

	for scanner.Scan() {
		t := scanner.Text()
		m[getKey(t)] = getValue(t)
	}

	return m, nil
}

func getRegex(txt, regex string) string {
	re := regexp.MustCompile(regex)
	return re.FindString(txt)
}

func getKey(l string) string {
	return getRegex(l, `^[^=]+`)
}

func getValue(l string) string {
	v := getRegex(l, `='[^']+`)
	return strings.ReplaceAll(v, `='`, "")
}
