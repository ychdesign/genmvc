package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/fatih/camelcase"
)

func GetGoMod() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	f := path.Join(pwd, "go.mod")
	if !fileExists(f) {
		return "", errors.New(f + " not exists")
	}

	file, err := os.Open(path.Join(pwd, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), "module") {
			return strings.TrimSpace(strings.TrimLeft(fileScanner.Text(), "module")), nil
		} else {
			return "", nil
		}
	}
	return "", nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lowerCaseFirst(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}

func upperCaseFirst(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}

func isPublicName(name string) bool {
	for _, c := range name {
		return unicode.IsUpper(c)
	}
	return false
}

func quote(tag string) string {
	return "`" + tag + "`"
}

func addGormTag(fieldName string) string {
	splitted := camelcase.Split(fieldName)
	var lowerSplitted []string
	for _, s := range splitted {
		lowerSplitted = append(lowerSplitted, strings.ToLower(s))
	}

	name := strings.Join(lowerSplitted, "_")
	return fmt.Sprintf("gorm:\"column:%s\"", name)
}

func addJsonTag(fieldName string) string {
	splitted := camelcase.Split(fieldName)
	var lowerSplitted []string
	for _, s := range splitted {
		lowerSplitted = append(lowerSplitted, strings.ToLower(s))
	}

	name := strings.Join(lowerSplitted, "_")
	return fmt.Sprintf("json:\"%s\"", name)
}
