package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
	"github.com/VictoriaMetrics/fastcache"
)

var (
	cacheSize = 1024 * 1024 * 32
)

type CacheService struct {
	cache *fastcache.Cache
	cacheExpiration int
	cacheChannel chan []byte
}

func NewCacheService() *CacheService {
	return &CacheService{
		cache: fastcache.New(cacheSize),
		cacheExpiration: 15,
		cacheChannel: make(chan []byte),
	}
}

func (s *CacheService) FindCache(namespace CacheNamespace, key string, entry interface{}) (bool, error) {
	buf2, has := s.cache.HasGet(nil, append([]byte(namespace), []byte(key)...))
	if has {
		decoder := gob.NewDecoder(bytes.NewReader(buf2))
		if err := decoder.Decode(entry); err != nil {
			return false, fmt.Errorf("failed to decode namespace %s, error: %w", namespace, err)
		}
		return true, nil
	}

	return false, nil
}

func (s *CacheService) UpsertCache(namespace CacheNamespace, key string, entry interface{}) error {
	var buf2 bytes.Buffer
	gob.Register(map[string]interface{}{})
	encoder := gob.NewEncoder(&buf2)
	if err := encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode for namespace: %s, error: %w", namespace, err)
	}

	k := append([]byte(namespace), []byte(key)...)
	s.cache.Set(k, buf2.Bytes())
	s.cacheChannel <- k

	return nil
}

func (s *CacheService) RunCacheHandler() {
	 for {
		 select {
		 	case data := <- s.cacheChannel:
				go func() {
					time.Sleep(time.Minute * time.Duration(s.cacheExpiration))

					s.cache.Del([]byte(data))
				}()
		 }
	 }
}