// pkg/firebase/fcm.go
package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMClient interface {
	SendToDevice(token string, title, body string, data map[string]string) (string, error)
	SendToDevices(tokens []string, title, body string, data map[string]string) (int, error)
	SendToTopic(topic, title, body string, data map[string]string) (string, error)
}

type fcmClient struct {
	client *messaging.Client
}

func NewFCMClient(credentialsFile string) (FCMClient, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &fcmClient{client: client}, nil
}

func (f *fcmClient) SendToDevice(token string, title, body string, data map[string]string) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Token: token,
	}

	response, err := f.client.Send(context.Background(), message)
	if err != nil {
		log.Printf("Error sending message to device: %v\n", err)
		return "", err
	}

	return response, nil
}

func (f *fcmClient) SendToDevices(tokens []string, title, body string, data map[string]string) (int, error) {
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
	}

	response, err := f.client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Printf("Error sending message to devices: %v\n", err)
		return 0, err
	}

	return response.SuccessCount, nil
}

func (f *fcmClient) SendToTopic(topic, title, body string, data map[string]string) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Topic: topic,
	}

	response, err := f.client.Send(context.Background(), message)
	if err != nil {
		log.Printf("Error sending message to topic: %v\n", err)
		return "", err
	}

	return response, nil
}
