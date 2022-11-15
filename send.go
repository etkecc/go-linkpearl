package linkpearl

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Send a message to the roomID and automatically try to encrypt it, if the destination room is encrypted
func (l *Linkpearl) Send(roomID id.RoomID, content interface{}) (id.EventID, error) {
	if !l.store.IsEncrypted(roomID) {
		return l.SendPlaintext(roomID, content)
	}

	encrypted, err := l.EncryptEvent(roomID, content)
	if err != nil {
		l.log.Error("cannot send encrypted message into %s: %v, sending plaintext...", roomID, err)
		return l.SendPlaintext(roomID, content)
	}

	// try to add missing relations
	if encrypted.RelatesTo == nil {
		originalContent, ok := content.(*event.Content)
		if ok {
			encrypted.RelatesTo = originalContent.AsMessage().RelatesTo
		}
	}

	return l.SendEncrypted(roomID, encrypted)
}

// SendFile to a matrix room
func (l *Linkpearl) SendFile(roomID id.RoomID, req *mautrix.ReqUploadMedia, msgtype event.MessageType, relation *event.RelatesTo) error {
	resp, err := l.GetClient().UploadMedia(*req)
	if err != nil {
		l.log.Error("cannot upload file %s: %v", req.FileName, err)
		return err
	}
	_, err = l.Send(roomID, &event.MessageEventContent{
		MsgType:   msgtype,
		Body:      req.FileName,
		URL:       resp.ContentURI.CUString(),
		RelatesTo: relation,
	})
	if err != nil {
		l.log.Error("cannot send uploaded file: %s: %v", req.FileName, err)
	}

	return err
}

// SendPlaintext sends plaintext event only
func (l *Linkpearl) SendPlaintext(roomID id.RoomID, content interface{}) (id.EventID, error) {
	resp, err := l.api.SendMessageEvent(roomID, event.EventMessage, content)
	if err != nil {
		return "", err
	}
	return resp.EventID, nil
}

// SendEncrypted sends encrypted event only
func (l *Linkpearl) SendEncrypted(roomID id.RoomID, content *event.EncryptedEventContent) (id.EventID, error) {
	resp, err := l.api.SendMessageEvent(roomID, event.EventEncrypted, content)
	if err != nil {
		return "", err
	}
	return resp.EventID, nil
}

// EncryptEvent before sending
func (l *Linkpearl) EncryptEvent(roomID id.RoomID, content interface{}) (*event.EncryptedEventContent, error) {
	encrypted, err := l.olm.EncryptMegolmEvent(roomID, event.EventMessage, content)
	if crypto.IsShareError(err) {
		err = l.olm.ShareGroupSession(roomID, l.store.GetRoomMembers(roomID))
		if err != nil {
			return nil, err
		}
		encrypted, err = l.olm.EncryptMegolmEvent(roomID, event.EventMessage, content)
	}

	return encrypted, err
}
