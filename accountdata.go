package linkpearl

import (
	"reflect"
	"strings"

	"maunium.net/go/mautrix/id"
)

func (l *Linkpearl) getAccountData(name string, output interface{}) error {
	l.log.Debug("retrieving account data %s", name)
	err := l.GetClient().GetAccountData(name, output)
	if err != nil && strings.Contains(err.Error(), "M_NOT_FOUND") {
		return err
	}

	l.log.Debug("storing account data %s to the cache", name)
	l.acc.Add(name, output)
	return nil
}

// GetAccountData returns the user's account data of this type (from cache and API)
func (l *Linkpearl) GetAccountData(name string, output interface{}) error {
	defer func() {
		rerr := recover()
		if rerr == nil {
			return
		}
		l.log.Error("failed to set cached account data %s using reflect: %v", name, rerr)
		err := l.getAccountData(name, output)
		if err != nil {
			l.log.Error("failed to retrieve account data %s: %v", name, err)
		}
	}()

	cached, ok := l.acc.Get(name)
	if ok {
		rv := reflect.ValueOf(output)
		rv.Elem().Set(reflect.ValueOf(cached).Elem())
		l.log.Debug("retrieved account data %s from cache", name)
		return nil
	}

	return l.getAccountData(name, output)
}

// SetAccountData sets the user's account data of this type (to cache and API)
func (l *Linkpearl) SetAccountData(name string, data interface{}) error {
	l.log.Debug("storing account data %s to the cache", name)
	l.acc.Add(name, data)

	return l.GetClient().SetAccountData(name, data)
}

func (l *Linkpearl) getRoomAccountData(roomID id.RoomID, name string, output interface{}) error {
	l.log.Debug("retrieving room %s account data %s", roomID, name)
	err := l.GetClient().GetRoomAccountData(roomID, name, output)
	if err != nil && strings.Contains(err.Error(), "M_NOT_FOUND") {
		return err
	}

	l.log.Debug("storing room %s account data %s to the cache", roomID, name)
	l.acc.Add(roomID.String()+name, output)
	return nil
}

// GetRoomAccountData returns the rooms's account data of this type (from cache and API)
func (l *Linkpearl) GetRoomAccountData(roomID id.RoomID, name string, output interface{}) error {
	defer func() {
		rerr := recover()
		if rerr == nil {
			return
		}
		l.log.Error("failed to set cached room %s account data %s using reflect: %v", roomID, name, rerr)
		err := l.getAccountData(name, output)
		if err != nil {
			l.log.Error("failed to retrieve room %s account data %s: %v", roomID, name, err)
		}
	}()

	key := roomID.String() + name
	cached, ok := l.acc.Get(key)
	if ok {
		rv := reflect.ValueOf(output)
		rv.Elem().Set(reflect.ValueOf(cached).Elem())
		l.log.Debug("retrieved room %s account data %s from cache", roomID, name)
		return nil
	}

	return l.getRoomAccountData(roomID, name, output)
}

// SetRoomAccountData sets the rooms's account data of this type (to cache and API)
func (l *Linkpearl) SetRoomAccountData(roomID id.RoomID, name string, data interface{}) error {
	l.log.Debug("storing room %s account data %s to the cache", roomID, name)
	l.acc.Add(roomID.String()+name, data)

	return l.GetClient().SetRoomAccountData(roomID, name, data)
}
