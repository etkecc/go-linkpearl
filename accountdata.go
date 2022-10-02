package linkpearl

import (
	"strings"
)

// GetAccountData returns the user's account data of this type (from cache and API)
func (l *Linkpearl) GetAccountData(name string, output interface{}) (interface{}, error) {
	if l.acc.Contains(name) {
		l.log.Debug("account data cache contains %s", name)
		cached, _ := l.acc.Get(name)
		return cached, nil
	}

	l.log.Debug("retrieving account data %s", name)
	err := l.GetClient().GetAccountData(name, output)
	if err != nil && strings.Contains(err.Error(), "M_NOT_FOUND") {
		return nil, err
	}

	l.log.Debug("storing account data %s to the cache", name)
	l.acc.Add(name, output)
	return output, nil
}

// SetAccountData sets the user's account data of this type (to cache and API)
func (l *Linkpearl) SetAccountData(name string, data interface{}) error {
	l.log.Debug("stroing account data %s to the cache", name)
	l.acc.Add(name, data)

	return l.GetClient().SetAccountData(name, data)
}
