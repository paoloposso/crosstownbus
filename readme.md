# CrosstownBus: a Golang Integration Event and Command Bus

A lightweight Go (golang) event bus that integrates with message brokers.
Brokers supported so far:
* RabbitMQ

The objective of this project is to abstract the communication between projects and message brokers, facilitating the creation of event-based services and micro-services.

## Installation
Run the command below to add crosstownbus as a library in your project:

```shell
go get github.com/paoloposso/crosstownbus
```

## Using CrosstownBus as an Event Bus

To create a bus to communicate using RabbitMQ, call the respective method, after installing and importing the crosstown module.

```shell
bus, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
```

After that, use the `bus` object (sample name) to perform the operations your deserve in your project, like publishing messages / events and / or subscribe.

### Subscribing to Events

When subscribing to an event, it's necessary to pass a handler as parameter, as shown bellow. 
* It's also necessary to inform the event type. 
* Use the type of the event you are creating to enforce the same type between publisher and subscribers 

```shell
bus.Subscribe(reflect.TypeOf(UserCreated{}), UserCreatedHandler{}, nil)
```

If you subscribe to a same event by using the same event type (like UserCreated in this example) the result will be multiple handlers receiving the same message.
When using RabbitMQ we are doing it by creating a `fanout exchange` for each `event` and a `queue`, connected to this exchange, for each `handler`. Note that, in this case, if you have multiple handlers with the same event pointing to the same event, the result is that a single queue will be created and, therefore, the messages will be balanced over the handlers.

This Handler must implement the interface IntegrationEventHandler, with a method `handler`.
Any message received by the broker will be passed to the handlers as an array of bytes `[]byte`.
The message will then be able to be handled as wanted inside of your project.
```shell
type IntegrationEventHandler interface {
	Handle(event []byte)
}
```

An example of a handler is shown next:
```shell
type UserCreatedHandler struct{}

func (handler UserCreatedHandler) Handle(event []byte) {
	var user *UserCreated
	json.Unmarshal(event, &user)
	fmt.Println(user.Name, "received:", time.Now())
	time.Sleep(5 * time.Second)
}
```

### Publishing Events
The same bus object can be used to publish events as well. 
* Use the `Publish` method to do it.
* Passing the event object as parameter.

```shell
bus.Publish(UserCreated{Name: "test", Id: 55})
```

### Retry Options: Retry Queues on RabbitMQ
If you need to retry processing a message when you have an error, you can add the following options when subscribing to an event:
* MaxRetryTimes: Number of times that the message will be moved to the retry queue and requeued to be consumed again.
* RetrySeconds: Time to wait before the message is requeued.

Example:
```
err = bus.Subscribe(reflect.TypeOf(
	eventsamples.UserCreated{}), 
	eventsamples.UserCreatedSendMailHandler{}, 
	&eventbus.ResilienceOptions{RetrySeconds: 5, MaxRetryTimes: 3},
)
```

In this case, when your handler returns an error, crosstownbus framework will verify the MaxRetryTimes and, if the max number was not reached, will send the message to the retry queue. When the time to wait (RetrySeconds) is reached, the message will be requeued and the consumer will receive the message again.