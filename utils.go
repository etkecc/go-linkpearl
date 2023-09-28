package linkpearl

import (
	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// EventRelatesTo uses evt as source for EventParent() and RelatesTo()
func EventRelatesTo(evt *event.Event) *event.RelatesTo {
	return RelatesTo(EventParent(evt.ID, evt.Content.AsMessage()))
}

// RelatesTo returns relation object of a matrix event (either threads with reply-to fallback or plain reply-to)
func RelatesTo(parentID id.EventID, noThreads ...bool) *event.RelatesTo {
	if parentID == "" {
		return nil
	}

	var nothreads bool
	if len(noThreads) > 0 {
		nothreads = noThreads[0]
	}

	if nothreads {
		return &event.RelatesTo{
			InReplyTo: &event.InReplyTo{
				EventID: parentID,
			},
		}
	}

	return &event.RelatesTo{
		Type:    event.RelThread,
		EventID: parentID,
		InReplyTo: &event.InReplyTo{
			EventID: parentID,
		},
		IsFallingBack: true,
	}
}

// EventParent returns parent event ID (either from thread or from reply-to relation)
func EventParent(currentID id.EventID, content *event.MessageEventContent) id.EventID {
	if content == nil {
		return currentID
	}

	relation := content.OptionalGetRelatesTo()
	if relation == nil {
		return currentID
	}

	threadParent := relation.GetThreadParent()
	if threadParent != "" {
		return threadParent
	}

	replyParent := relation.GetReplyTo()
	if replyParent != "" {
		return replyParent
	}

	return currentID
}

// EventContains checks if raw event content contains specified field with specified values
func EventContains[T comparable](evt *event.Event, field string, value T) bool {
	if evt.Content.Raw == nil {
		return false
	}
	if EventField[T](&evt.Content, field) != value {
		return false
	}

	return true
}

// EventField returns field value from raw event content
func EventField[T comparable](content *event.Content, field string) T {
	var zero T
	raw, ok := content.Raw[field]
	if !ok {
		return zero
	}

	if raw == nil {
		return zero
	}

	v, ok := raw.(T)
	if !ok {
		return zero
	}

	return v
}

func ParseContent(evt *event.Event, eventType event.Type, log *zerolog.Logger) {
	if evt.Content.Parsed != nil {
		return
	}
	perr := evt.Content.ParseRaw(eventType)
	if perr != nil {
		log.Error().Err(perr).Msg("cannot parse event content")
	}
}

// UnwrapError tries to unwrap a error into something meaningful, like mautrix.HTTPError or mautrix.RespError
func UnwrapError(err error) error {
	switch err.(type) {
	case nil:
		return nil
	case mautrix.HTTPError:
		return unwrapHTTPError(err)
	default:
		return err
	}
}

func unwrapHTTPError(err error) error {
	httperr, ok := err.(mautrix.HTTPError)
	if !ok {
		return err
	}

	uwerr := httperr.Unwrap()
	if uwerr != nil {
		return uwerr
	}

	return httperr
}
