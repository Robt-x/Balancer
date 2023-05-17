package balancer

import (
	"errors"
	"log"
)

var (
	NoHostError                = errors.New("no hosts")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}
type Factory func([]string) Balancer

var factories = make(map[string]Factory)

func Build(algorithm string, hosts []string) (Balancer, error) {
	factory, ok := factories[algorithm]
	if !ok {
		log.Fatal(AlgorithmNotSupportedError)
		return nil, AlgorithmNotSupportedError
	}
	return factory(hosts), nil
}
