package gormutil

import (
	"reflect"

	"gorm.io/gorm/schema"
)

const (
	// HookAfterCreate is event name for AfterCretate hook
	HookAfterCreate = "AfterCreate"
	// HookAfterUpdate is event name for AfterUpdate hook
	HookAfterUpdate = "AfterUpdate"
	// HookAfterDelete is event name for AfterDelete hook
	HookAfterDelete = "AfterDelete"
)

// HookEvent represents event that triggered hook
type HookEvent string

// IsAfterCreate checks whether event is "AfterCreate"
func (e HookEvent) IsAfterCreate() bool {
	return string(e) == HookAfterCreate
}

// IsAfterUpdate checks whether event is "AfterUpdate"
func (e HookEvent) IsAfterUpdate() bool {
	return string(e) == HookAfterUpdate
}

// IsAfterDelete checks whether event is "AfterDelete"
func (e HookEvent) IsAfterDelete() bool {
	return string(e) == HookAfterDelete
}

// String returns event's string representation
func (e HookEvent) String() string {
	return string(e)
}

// HookHandlerFunc is a subscription's callback
type HookHandlerFunc func(model interface{}, event HookEvent)

// HookSubscription defines hook's subscription
type HookSubscription struct {
	handler HookHandlerFunc
	table   string
}

// Hook defines hook event
type Hook struct {
	table string
	model interface{}
	event HookEvent
}

// HookBus maintains the set of subscriptions and broadcast any incoming hooks
type HookBus struct {
	// subscriptions defines list of subscriptions
	subscriptions map[*HookSubscription]struct{}

	// subscribeChan adds given subscription to the list
	subscribeChan chan *HookSubscription

	// publishChan broadcasts given hook to its subscriptions
	publishChan chan *Hook
}

func (hb *HookBus) publish(model interface{}, event HookEvent) {
	hb.publishChan <- &Hook{
		table: tableName(model),
		model: model,
		event: event,
	}
}

func (hb *HookBus) subscribe(model interface{}, fn HookHandlerFunc) {
	hb.subscribeChan <- &HookSubscription{
		table:   tableName(model),
		handler: fn,
	}
}

func newHookBus() *HookBus {
	return &HookBus{
		subscriptions: make(map[*HookSubscription]struct{}),
		subscribeChan: make(chan *HookSubscription),
		publishChan:   make(chan *Hook),
	}
}

func (hb *HookBus) run() {
	for {
		select {
		case sub := <-hb.subscribeChan:
			hb.subscriptions[sub] = struct{}{}
		case hook := <-hb.publishChan:
			for sub := range hb.subscriptions {
				if sub.table == hook.table {
					go sub.handler(hook.model, hook.event)
				}
			}
		}
	}
}

func tableName(v interface{}) string {
	value := reflect.ValueOf(v)
	modelType := reflect.Indirect(value).Type()
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	namer := schema.NamingStrategy{}
	return namer.TableName(modelType.Name())
}
