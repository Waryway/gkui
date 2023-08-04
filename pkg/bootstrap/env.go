package bootstrap

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
	"strconv"
)

const (
	YML = "yaml"
	DEF = "default"
	REQ = "required"
)

type BootStrap interface {
}

type Config[boot BootStrap] struct {
	BootStrap boot
	BootCh    chan Config[boot]
	ErrCh     chan error
	file      *[]byte
	location  string
}

// Init will read the settings.yml file into the bootstrap config.
// It will override and Environment Variable settings.
// Also initializes the channels
func (c *Config[boot]) Init() *Config[boot] {
	c.location = "settings.yaml"
	c.BootCh = make(chan Config[boot])
	c.ErrCh = make(chan error)

	return c
}

// Load will persist the current applications bootstrap settings to the settings.yml file.
func (c *Config[boot]) Load() *Config[boot] {
	go func(conf *Config[boot]) {
		c.ReadSettingFile().LoadYaml().LoadEnv()
		c.BootCh <- *c
	}(c)

	select {
	case err := <-c.ErrCh:
		fmt.Println(err)
	case res := <-c.BootCh:
		return &res
	}

	return c
}

// Save will persist the current applications bootstrap settings to the settings.yml file.
func (c *Config[boot]) Save() {
	if d, err := yaml.Marshal(c.BootStrap); err != nil {
		c.ErrCh <- err
	} else {
		if f, err := os.Create(c.location); err != nil {
			c.ErrCh <- err
		} else {
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					c.ErrCh <- err
				}
			}(f)
			if _, err = f.Write(d); err != nil {
				c.ErrCh <- err
			}
		}
	}
}

// LoadEnv checks for any configured BootStrap values from the environment.
// Effectively can only override top level yaml configuration settings.
// Adapted from https://github.com/Vish511/GoENV/blob/master/envconfigloader.go
func (c *Config[boot]) LoadEnv() *Config[boot] {
	v := reflect.ValueOf(c.BootStrap)
	var failed []string
	failedCount := 0

	for i := 0; i < v.NumField(); i++ {
		tagName := v.Type().Field(i).Tag.Get(YML)
		tagDefault := v.Type().Field(i).Tag.Get(DEF)

		if tagName == "" || tagName == "-" {
			continue
		}

		envVal, Defaulted := loadFromEnv(tagName, tagDefault, c.ErrCh)

		if Defaulted && envVal != "" {
			log.Println("Loaded: " + tagName + "=" + envVal + " with default")
		} else if envVal != "" {
			log.Println("Loaded: " + tagName + "=" + envVal)
		}

		if on, err := strconv.ParseBool(v.Type().Field(i).Tag.Get(REQ)); err == nil && on && envVal == "" {
			failedCount++
			failed = append(failed, tagName)
			log.Println("Missing: " + tagName + "=" + "")
		}

		if reflect.ValueOf(&c.BootStrap).Elem().Field(i).String() == "" && envVal != "" {
			reflect.ValueOf(&c.BootStrap).Elem().Field(i).SetString(envVal)
		} else {
			log.Println("Loaded: " + tagName + "=" + reflect.ValueOf(&c.BootStrap).Elem().Field(i).String())
		}

		if failedCount > 0 {
			c.ErrCh <- errors.New("Failed to load" + strconv.Itoa(failedCount) + " vars")
		}
	}

	return c
}

func (c *Config[boot]) LoadYaml() *Config[boot] {
	if err := yaml.Unmarshal(*c.file, &c.BootStrap); err != nil {
		c.ErrCh <- err
	}

	return c
}

func (c *Config[boot]) ReadSettingFile() *Config[boot] {
	if file, err := os.ReadFile(c.location); err != nil {
		log.Println("Failed to read" + c.location)
		c.ErrCh <- err
	} else {
		c.file = &file
	}

	return c
}

// loadFromEnv Given a key lookup the envvar with a default fallback. return the value and true if the default was used.
// k is the envvar key.
// def is the default value to use if the envvar is empty.
func loadFromEnv(k string, def string, errorCh chan error) (string, bool) {
	if v := os.Getenv(k); v == "" {
		if err := os.Setenv(k, def); err != nil {
			errorCh <- err
		} else {
			v = def
		}
		return v, true
	} else {
		return v, false
	}
}
