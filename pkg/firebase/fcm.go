// pkg/firebase/fcm.go
package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type FCMClient interface {
	SendToDevice(token string, title, body string, data map[string]string) (string, error)
	SendToDevices(tokens []string, title, body string, data map[string]string) (int, error)
	SendToTopic(topic, title, body string, data map[string]string) (string, error)
}

type fcmClient struct {
	client *messaging.Client
	logger *zap.Logger
}

func NewFCMClient(credentialsFile string, logger *zap.Logger) (FCMClient, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &fcmClient{
		client: client,
		logger: logger,
	}, nil
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
		f.logger.Error("Error sending message to device", zap.Error(err), zap.String("token", token[:10]+"..."))
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
		f.logger.Error("Error sending message to devices", zap.Error(err), zap.Int("token_count", len(tokens)))
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
		f.logger.Error("Error sending message to topic", zap.Error(err), zap.String("topic", topic))
		return "", err
	}

	return response, nil
}
