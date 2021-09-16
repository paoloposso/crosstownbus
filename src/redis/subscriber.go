package redis

import eventbus "github.com/paoloposso/crosstown-bus/src/event-bus"

type RedisSubscriber struct{}

func CreateRedisSubscriber() eventbus.Subscriber {
	return RedisSubscriber{}
}

func (subs RedisSubscriber) SubscribeEvent(config struct{}, eventName string, eventHandler eventbus.IntegrationEventHandler) {

}
