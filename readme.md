# CrosstownBus: a Golang Integration Event Bus

A Go (golang) event bus that integrates with message brokers.
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

To create a bus to communicate with Redis, call the respective method, after installing and importing the crosstown module.

```shell
bus, err := crosstownbus.CreateRedisBus(reflect.TypeOf(UserCreated{}), "localhost:6379", "")
```

After that, use the `bus` object (sample name) to do perform the operations your deserve in your project, like publishing messages / events and / or subscribe.

### Subscribe

When subscribing to an event, it's necessary to pass a handler as parameter, as shown bellow.

```shell
bus.Subscribe(UserCreatedHandler{})
```

This Handler must implement the interface IntegrationEventHandler, with a method `handler`.
Any message received by the broker will be passed to this handler as an array of bytes `[]byte`.
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