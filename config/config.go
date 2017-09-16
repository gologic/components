package config

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

type ConfigInterface interface {
	Get(key string) string
	Has(key string) bool
	All() map[string]string
	GetWithFallback(key string, fallback string) string
}

type config struct {
	values map[string]string
}

func Load(filename string) (*config, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	line, isPrefix, err := r.ReadLine()
	keyValue := regexp.MustCompile(`^[A-Za-z0-9_-]+\s*=.*$`)
	values := make(map[string]string)

	for err == nil && !isPrefix {

		s := string(line)

		if keyValue.MatchString(s) {
			delimIndex := strings.Index(s, "=")
			key := strings.Trim(s[:delimIndex], " ")
			value := strings.Trim(s[delimIndex+1:], " ")
			values[key] = value
		}

		line, isPrefix, err = r.ReadLine()
	}

	if isPrefix {
		return nil, errors.New("config file line exceeded read buffer size")
	}

	return &config{values}, nil
}

func (c config) Get(key string) string {
	return c.values[key]
}

func (c config) Has(key string) bool {
	return c.values[key] != ""
}

func (c config) GetWithFallback(key string, fallback string) string {
	if c.values[key] != "" {
		return c.values[key]
	}
	return fallback
}

func (c config) All() map[string]string {
	return c.values
}
