package cachetypes

import (
	"sync"
	"math/rand"
)

type RssCache struct {
	sync.Mutex
	Data []string //Writing array
}

//Choose random writing
func (rssCache *RssCache) ChooseRandom() string {
	return rssCache.Data[rand.Intn(len(rssCache.Data) - 1)]
}