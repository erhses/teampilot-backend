package integrations

import (
	"context"
	"fmt"
	"log"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

var (
	fcmClient *messaging.Client
	fonce     sync.Once
)

func initFirebaseApp() error {
	var initError error
	fonce.Do(func() {
		opt := option.WithCredentialsFile("serviceAccountKey.json")

		config := &firebase.Config{
			ProjectID: "realestate-b3bda",
		}

		app, err := firebase.NewApp(context.Background(), config, opt)
		if err != nil {
			initError = fmt.Errorf("error initializing firebase app: %v", err)
			return
		}

		client, err := app.Messaging(context.Background())
		if err != nil {
			initError = fmt.Errorf("error initializing messaging client: %v", err)
			return
		}

		fcmClient = client
	})
	return initError
}

type FNotificationRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func FSendNotifToAll(c *fiber.Ctx) error {
	if err := initFirebaseApp(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var req NotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := SendNotification(req.Title, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"response": response,
	})
}

func FSendNotifWithID(title, content string, deviceToken string) (string, error) {
	if err := initFirebaseApp(); err != nil {
		return "", err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  content,
		},
		Token: deviceToken,
	}
	response, err := fcmClient.Send(context.Background(), message)
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	return response, nil
}

func FSendNotification(title, content string) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  content,
		},
		Topic: "all",
	}

	response, err := fcmClient.Send(context.Background(), message)
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	return response, nil
}

func SendMulticastNotification(title, content string, deviceTokens []string) (*messaging.BatchResponse, error) {
	if err := initFirebaseApp(); err != nil {
		return nil, err
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  content,
		},
		Tokens: deviceTokens,
	}

	br, err := fcmClient.SendEachForMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	return br, nil
}
func SendToTopics(title, content string, topic string) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  content,
		},
		Topic: topic,
	}

	response, err := fcmClient.Send(context.Background(), message)
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}
		

	return response, nil
}

func SubscribeToTopic(deviceToken, topic string) error {
	if err := initFirebaseApp(); err != nil {
		return err
	}

	_, err := fcmClient.SubscribeToTopic(context.Background(), []string{deviceToken}, topic)
	if err != nil {
		return fmt.Errorf("error subscribing to topic: %v", err)
	}

	return nil
}

func UnsubscribeFromTopic(deviceToken, topic string) error {
	if err := initFirebaseApp(); err != nil {
		return err
	}

	_, err := fcmClient.UnsubscribeFromTopic(context.Background(), []string{deviceToken}, topic)
	if err != nil {
		return fmt.Errorf("error unsubscribing from topic: %v", err)
	}

	return nil
}
