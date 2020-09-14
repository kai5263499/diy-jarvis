package diyjarvis

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func NewMqttComms(clientID string, mqttBroker string) (*MqttComms, error) {

	opts := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(clientID)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		logrus.WithError(token.Error()).Fatal("mqtt new client")
		return nil, token.Error()
	}

	mc := &MqttComms{
		mqttClient: mqttClient,
		baseChan:   make(chan pb.Base, 100),
	}

	if token := mc.mqttClient.Subscribe(domain.RequestTopic, 0, mc.messageHandler); token.Wait() && token.Error() != nil {
		logrus.WithError(token.Error()).Fatal("mqtt subscribe")
		return nil, token.Error()
	}

	return mc, nil
}

type MqttComms struct {
	mqttClient mqtt.Client
	baseChan   chan pb.Base
}

func (m *MqttComms) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var req pb.Base
	if err := proto.Unmarshal(msg.Payload(), &req); err != nil {
		logrus.WithError(err).Errorf("unable to unmarshal request")
		return
	}

	logrus.Debugf("got message with %d bytes from %s with type=%#v", len(msg.Payload()), msg.Topic(), req.Type)

	m.baseChan <- req
}

func (m *MqttComms) RequestChan() chan pb.Base {
	return m.baseChan
}

func (m *MqttComms) SendRequest(req pb.Base) error {
	pubBytes, err := proto.Marshal(&req)
	if err != nil {
		logrus.WithError(err).Errorf("unable to marshal proto")
		return err
	}

	if token := m.mqttClient.Publish(domain.RequestTopic, 0, false, pubBytes); token.Error() != nil {
		logrus.WithError(token.Error()).Errorf("error sending request to marshal proto")
		return token.Error()
	}
	logrus.Debugf("published %d bytes to %s", len(pubBytes), domain.RequestTopic)

	return nil
}

func (m *MqttComms) Close() {
	m.mqttClient.Disconnect(100)
}

func (m *MqttComms) MQTTClient() mqtt.Client {
	return m.mqttClient
}
