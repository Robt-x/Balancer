package balancer

import "hash/crc32"

type IPHash struct {
	BaseBalancer
}

func init() {
	factories[IPHashBalancer] = NewIPHash
}
func NewIPHash(hosts []string) Balancer {
	return &IPHash{BaseBalancer{hosts: hosts}}
}

func (i *IPHash) Balance(key string) (string, error) {
	i.RLock()
	defer i.RUnlock()
	if len(i.hosts) == 0 {
		return "", NoHostError
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(i.hosts))
	return i.hosts[value], nil
}
