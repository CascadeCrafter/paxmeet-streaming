package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

// Config structure defining structure of config.yaml
type Config struct {
	IngressURL string `yaml:"ingress_url"`
}

// ReadConfig reads YAML configuration from file
func ReadConfig() (*Config, error) {
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// main function
func main() {
	config, err := ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	target, err := url.Parse(config.IngressURL)
	if err != nil {
		fmt.Println("Error parsing ingress URL:", err)
		os.Exit(1)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(rw, req)
	})

	fmt.Println("Starting proxy server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
