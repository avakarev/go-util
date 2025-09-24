package gormutil

import (
	"reflect"

	"gorm.io/gorm/schema"
)

const (
	// HookAfterCreate is event name for AfterCreate hook
	HookAfterCreate = "AfterCreate"
	// HookAfterUpdate is event name for AfterUpdate hook
	HookAfterUpdate = "AfterUpdate"
	// HookAfterDelete is event name for AfterDelete hook
	HookAfterDelete = "AfterDelete"
)

// HookEvent represents event that triggered hook
type HookEvent string

// String returns event's string representation
func (e HookEvent) String() string {
	return string(e)
}

// IsAfterCreate checks whether event is "AfterCreate"
func (e HookEvent) IsAfterCreate() bool {
	return e.String() == HookAfterCreate
}

// IsAfterUpdate checks whether event is "AfterUpdate"
func (e HookEvent) IsAfterUpdate() bool {
	return e.String() == HookAfterUpdate
}

// IsAfterDelete checks whether event is "AfterDelete"
func (e HookEvent) IsAfterDelete() bool {
	return e.String() == HookAfterDelete
}

// Hook defines hook event
type Hook struct {
	Table string
	Model any
	Event HookEvent
}

// HookHandlerFunc is a subscription's callback
type HookHandlerFunc func(hook *Hook)

// HookSubscription defines hook's subscription
type HookSubscription struct {
	handler HookHandlerFunc
	table   string
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

func (hb *HookBus) publish(model any, event HookEvent) {
	hb.publishChan <- &Hook{
		Table: tableName(model),
		Model: model,
		Event: event,
	}
}

func (hb *HookBus) subscribe(model any, fn HookHandlerFunc) {
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
				if sub.table == hook.Table {
					go sub.handler(hook)
				}
			}
		}
	}
}

func tableName(v any) string {
	value := reflect.ValueOf(v)
	modelType := reflect.Indirect(value).Type()
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	namer := schema.NamingStrategy{}
	return namer.TableName(modelType.Name())
}
