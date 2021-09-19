# CrosstownBus: a Golang Integration Event Bus

A Go (golang) event and command bus that integrates with message brokers.
Brokers supported so far:
* Redis
* RabbitMQ

The objective of this project is to abstract the communication between projects and message brokers, facilitating the creation of event-based services and micro-services.

## Instalation
Run the command below to add crosstownbus as a library in your project:

```shell
go get github.com/paoloposso/crosstownbus
```

## Usage Example

To create a bus to communicate using RabbitMQ, call the respective method, after installing and importing the crosstown module.

```shell
bus, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
```

After that, use the `bus` object (sample name) to do perform the operations your deserve in your project, like publishing messages / events and / or subscribe.

### Subscribing to Events

When subscribing to an event, it's necessary to pass a handler as parameter, as shown bellow. It's also necessary to inform the event type. Use the type of the event you are creating to enforce the same type between publisher and subscribers 

```shell
bus.Subscribe(reflect.TypeOf(UserCreated{}), UserCreatedHandler{})
```

If you subscribe to a same event by using the same event type (like UserCreated in this example) the result will be multiple handlers receiving the same message.
When using RabbitMQ we are doing it by creating a `fanout exchange` for each `event` and a `queue`, connected to this exchange, for each `handler`.

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