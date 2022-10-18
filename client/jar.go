package shellgame

import (
	"golang.org/x/net/publicsuffix"
	"net/http/cookiejar"
	"sync"
)

var lock = &sync.Mutex{}

var jar *cookiejar.Jar

func getJar() (*cookiejar.Jar, error) {
	if jar == nil {
		lock.Lock()
		defer lock.Unlock()
		if jar == nil {
			var err error
			jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
			if err != nil {
				return nil, err
			}
		}
	}
	return jar, nil
}
