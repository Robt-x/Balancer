package balancer

import "github.com/lafikl/consistent"

func init() {
	factories[ConsistentHashBalancer] = NewConsistent
}

type ConsistentHash struct {
	ch *consistent.Consistent
}

func NewConsistent(hosts []string) Balancer {
	c := &ConsistentHash{
		ch: consistent.New(),
	}
	for _, h := range hosts {
		c.Add(h)
	}
	return c
}
func (c *ConsistentHash) Balance(key string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NoHostError
	}
	return c.ch.Get(key)
}

func (c *ConsistentHash) Add(hostName string) {
	c.ch.Add(hostName)
}

func (c *ConsistentHash) Remove(host string) {
	c.ch.Remove(host)
}
func (c *ConsistentHash) Inc(host string) {
	c.ch.Inc(host)
}

func (c *ConsistentHash) Done(host string) {
	c.ch.Done(host)
}
