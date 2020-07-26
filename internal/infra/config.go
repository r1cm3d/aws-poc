package infra

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func loadConf() map[string]string {
	envDir := "../../scripts/env/"
	files, err := ioutil.ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}

	var configs []map[string]string
	for _, f := range files {
		configs = append(configs, loadFile(envDir+f.Name()))
	}

	return merge(configs...)
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

func loadFile(path string) map[string]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	m := make(map[string]string)

	for scanner.Scan() {
		t := scanner.Text()
		m[getKey(t)] = getValue(t)
	}

	return m
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
