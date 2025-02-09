package thing_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/thing"
	"github.com/stretchr/testify/assert"
)

var thingName = ""
var testEndpoint = ""

func TestMain(m *testing.M) {
	var ok bool

	thingName, ok = os.LookupEnv("AWS_IOT_THING_NAME")
	if !ok {
		panic("AWS_IOT_THING_NAME environment variable must be defined")
	}

	testEndpoint, ok = os.LookupEnv("AWS_MQTT_ENDPOINT")
	if !ok {
		panic("AWS_MQTT_ENDPOINT environment variable must be defined")
	}

	code := m.Run()
	os.Exit(code)
}

var keyPair = models.KeyPair{
	CertificatePath:   "./certificates/cert.pem",
	PrivateKeyPath:    "./certificates/private.key",
	CACertificatePath: "./certificates/root.ca.pem",
}

type shadowStruct struct {
	State struct {
		Reported struct {
			Value int64 `json:"value"`
		} `json:"reported"`
	} `json:"state"`
}

func TestNewThingFromFiles(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")

	defer th.Disconnect()
}

func TestNewThingFromStrings(t *testing.T) {
	cert, err := ioutil.ReadFile(keyPair.CertificatePath)
	key, err := ioutil.ReadFile(keyPair.PrivateKeyPath)
	th, err := thing.NewThingFromStrings(string(cert), string(key), testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")

	defer th.Disconnect()
}

func TestThingShadow(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	thingShadowChan, _, err := th.SubscribeForThingShadowChanges()
	assert.NoError(t, err, "received thing shadow subscription channel without error")

	data := time.Now().Unix()

	shadow := fmt.Sprintf(`{"state": {"reported": {"value": %d}}}`, data)

	err = th.UpdateThingShadow(thing.Shadow(shadow))
	assert.NoError(t, err, "thing shadow updated without error")

	updatedShadow, ok := <-thingShadowChan
	assert.True(t, ok, "the reading updated shadow channel was successful")

	unmarshaledUpdatedShadow := &shadowStruct{}

	err = json.Unmarshal(updatedShadow, unmarshaledUpdatedShadow)
	assert.NoError(t, err, "thing shadow payload unmarshaled without error")

	assert.Equal(t, data, unmarshaledUpdatedShadow.State.Reported.Value, "thing shadow update has consistent data")

	gottenShadow, err := th.GetThingShadow()
	assert.NoError(t, err, "retrieved thing shadow without error")

	unmarshaledGottenShadow := &shadowStruct{}

	err = json.Unmarshal(gottenShadow, unmarshaledGottenShadow)
	assert.NoError(t, err, "retrieved thing shadow unmarshaling  without error")

	assert.Equal(t, data, unmarshaledGottenShadow.State.Reported.Value, "retrieved thing shadow has consistent data")
}

func TestThing_UpdateThingShadowShouldFail(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	_, thingShadowErrChan, err := th.SubscribeForThingShadowChanges()
	assert.NoError(t, err, "received thing shadow subscription channel without error")

	err = th.UpdateThingShadow(thing.Shadow("invalid JSON"))
	assert.NoError(t, err, "thing shadow updated without error")

	_, ok := <-thingShadowErrChan
	assert.True(t, ok, "the update shadow error has been handled successfully")
}

func TestThing_UpdateThingShadowDocument(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	shadowChan, err := th.SubscribeForCustomTopic(fmt.Sprintf("$aws/things/%s/shadow/update/documents", thingName))
	assert.NoError(t, err, "received thing shadow subscription channel without error")

	shadowDocument := fmt.Sprintf(`{"state": {"reported": {"value": %d}}}`, time.Now().UTC().Unix())

	err = th.UpdateThingShadowDocument(thing.Shadow(shadowDocument))
	assert.NoError(t, err, "thing shadow document updated without error")

	remoteShadow, ok := <-shadowChan
	assert.True(t, ok, "the update shadow document has been handled successfully")

	assert.Equal(t, thing.Payload(shadowDocument), remoteShadow)
}

func TestThing_DeleteThingShadow(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	err = th.DeleteThingShadow()
	assert.NoError(t, err, "thing shadow deleted without error")
}

func TestThing_ListenForJobs(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	_, err = th.ListenForJobs()
	assert.NoError(t, err, "subscribed to iot core jobs with no errors")
}

func TestThing_UnsubscribeFromJobs(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	err = th.UnsubscribeFromJobs()
	assert.NoError(t, err, "Unsubscribed to iot core jobs with no errors")
}

func TestThing_CustomTopic(t *testing.T) {
	th, err := thing.NewThingFromFiles(keyPair, testEndpoint, thingName)
	assert.NoError(t, err, "thing instance created without error")
	assert.NotNil(t, th, "thing instance is not nil")
	defer th.Disconnect()

	customTopic := fmt.Sprintf("$aws/things/%s/fancy", thingName)

	customChan, err := th.SubscribeForCustomTopic(customTopic)
	assert.NoError(t, err, "received thing shadow custom topic subscription channel without error")

	customPayload := thing.Payload(`{"state":{"reported":{"yo":true}}}`)

	err = th.PublishToCustomTopic(customPayload, customTopic)
	assert.NoError(t, err, "thing shadow published to custom topic updated without error")

	remotePayload, ok := <-customChan
	assert.True(t, ok, "the shadow in custom topic has been handled successfully")

	assert.Equal(t, customPayload, remotePayload)
}
