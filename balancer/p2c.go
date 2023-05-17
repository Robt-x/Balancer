package balancer

import (
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

func init() {
	factories[P2CBalancer] = NewP2C
}

const Salt = "%#!"

type host struct {
	name string
	load uint64
}

type P2C struct {
	sync.RWMutex
	hosts   []*host
	rnd     *rand.Rand
	loadMap map[string]*host
}

func NewP2C(hosts []string) Balancer {
	p := &P2C{
		hosts:   []*host{},
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
		loadMap: make(map[string]*host),
	}
	for _, h := range hosts {
		p.Add(h)
	}
	return p
}
func (p *P2C) Balance(key string) (string, error) {
	p.RLock()
	defer p.RUnlock()
	h1, h2 := p.hash(key)
	if p.loadMap[h1].load <= p.loadMap[h2].load {
		return h1, nil
	}
	return h2, nil
}
func (p *P2C) hash(Key string) (string, string) {
	var n1, n2 string
	if len(Key) > 0 {
		saltKey := Key + Salt
		n1 = p.hosts[crc32.ChecksumIEEE([]byte(Key))%uint32(len(p.hosts))].name
		n2 = p.hosts[crc32.ChecksumIEEE([]byte(saltKey))%uint32(len(p.hosts))].name
		return n1, n2
	}
	return p.hosts[p.rnd.Intn(len(p.hosts))].name, p.hosts[p.rnd.Intn(len(p.hosts))].name

}
func (p *P2C) Add(hostName string) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[hostName]; ok {
		return
	}
	h := &host{
		name: hostName,
		load: 0,
	}
	p.hosts = append(p.hosts, h)
	p.loadMap[hostName] = h
}

func (p *P2C) Remove(hostName string) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[hostName]; !ok {
		return
	}
	delete(p.loadMap, hostName)
	for i, h := range p.hosts {
		if h.name == hostName {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			return
		}
	}
}
func (p *P2C) Inc(host string) {
	p.Lock()
	defer p.Unlock()
	if h, ok := p.loadMap[host]; ok {
		h.load++
	}
}

func (p *P2C) Done(host string) {
	p.Lock()
	defer p.Unlock()
	if h, ok := p.loadMap[host]; ok {
		if h.load > 0 {
			h.load++
		}
	}
}
