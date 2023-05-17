package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var (
	ascii = `
___ _ _  _ _   _ ___  ____ _    ____ _  _ ____ ____ ____ 
 |  | |\ |  \_/  |__] |__| |    |__| |\ | |    |___ |__/ 
 |  | | \|   |   |__] |  | |___ |  | | \| |___ |___ |  \                                        
`
)

type Configure struct {
	Schema              string      `yaml:"schema"`
	Port                int         `yaml:"port"`
	SSLCertificate      string      `yaml:"ssl_certificate"`
	SSLCertificateKey   string      `yaml:"ssl_certificate_key"`
	HealthCheck         bool        `yaml:"tcp_health_check"`
	HealthCheckInterval uint        `yaml:"health_check_interval"`
	MaxAllowed          uint        `yaml:"max_allowed"`
	Location            []*Location `yaml:"location"`
}
type Location struct {
	Pattern     string   `yaml:"pattern"`
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

func ReadConfig(fileName string) (*Configure, error) {
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var config Configure
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &config, nil
}

func (c *Configure) Validation() error {
	if c.Schema != "http" && c.Schema != "https" {
		return fmt.Errorf("the Schema \"%s\" nod supported", c.Schema)
	}
	if len(c.Location) == 0 {
		return errors.New("the details of location cannot be null")
	}
	if c.Schema == "https" && (len(c.SSLCertificate) == 0 || len(c.SSLCertificateKey) == 0) {
		return errors.New("the https requires ssl_certificate and ssl_certificate_key")
	}
	if c.HealthCheckInterval < 1 {
		return errors.New("health_check_interval must be greater than 0")
	}

	return nil
}
func (c *Configure) Print() {
	fmt.Printf("Schema: %s\nPort: %d\nHealth Check: %v\nLocation:\n", c.Schema, c.Port, c.HealthCheck)
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tProxy Pass: %s\n\tMode: %s\n\n",
			l.Pattern, l.ProxyPass, l.BalanceMode)
	}
}
