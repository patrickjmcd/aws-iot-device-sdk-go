# AWS IoT SDK for Go lang
The aws-iot-device-sdk-go package allows developers to write Go lang applications which access the AWS IoT Platform via MQTT.
## Install
`go get "github.com/joshuarose/aws-iot-device-sdk-go"`
## Example
```
package main
import (
    "github.com/joshuarose/aws-iot-device-sdk-go/device"
    "fmt"
)

func main() {
    thing, err := device.NewThingFromFiles(
        device.KeyPair{
            PrivateKeyPath: "path/to/private/key",
            CertificatePath: "path/to/certificate",
            CACertificatePath: "path/to/rootCA",
        },
        "xxxxxxxxxx.iot.us-east-1.amazonaws.com",  // AWS IoT endpoint
        device.ThingName("thing_name"),
    )
    if err != nil {
        panic(err)
    }

    s, err := thing.GetThingShadow()
    if err != nil {
        panic(err)
    }
    fmt.Println(s)

    shadowChan, _, err := thing.SubscribeForThingShadowChanges()
    if err != nil {
        panic(err)
    }

    for {
        select {
        case s, ok := <- shadowChan:
            if !ok {
                panic("failed to read from shadow channel")
            }
            fmt.Println(s)
        }
    }
}
```
## Contributing
```
- To test create a certificates directory in device/
- Add device cert as certificates/cert.pem
- Add device private key to certificates/private.key
- Add Amazon root from https://www.amazontrust.com/repository/AmazonRootCA1.pem and add to certificates/root.ca.pem
```
## Reference
```
// NewThingFromFiles returns a new instance of Thing
func NewThingFromFiles(keyPair KeyPair, thingName ThingName, region Region) (*Thing, error)
```
```
// NewThingFromStrings returns a new instance of Thing
func NewThingFromStrings(cert string, key string, awsEndpoint string, thingName ThingName) (*Thing, error
```
```
// GetThingShadow gets the current thing shadow
func (t *Thing) GetThingShadow() (Shadow, error)
```
```
// UpdateThingShadow publish a message with new thing shadow
func (t *Thing) UpdateThingShadow(payload Shadow) error
```
```
// SubscribeForThingShadowChanges returns the channel with the shadow updates
func (t *Thing) SubscribeForThingShadowChanges() (chan Shadow, error) 
```
```
// SubscribeForThingShadowChanges returns the channel with the shadow updates
func (t *Thing) SubscribeForThingShadowChanges() (chan Shadow, error) 
```
```
// ListenForJobs is a helper function that subscribes to the topic responsible for notifying on IoT Core Jobs
func (t *Thing) ListenForJobs() (chan Payload, error)
```
```
// Publish to a custom topic
func (t *Thing) PublishToCustomTopic(payload Payload, topic string) error
```
```
// SubscribeForCustomTopic subscribes for the custom topic and returns the channel with the topic messages.
func (t *Thing) SubscribeForCustomTopic(topic string) (chan Payload, error)
```
