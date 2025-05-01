package integrations

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/tbalthazar/onesignal-go"
)

// Global OneSignal client instance
var (
	oneSignalClient *onesignal.Client
	once            sync.Once
)

// Initialize OneSignal client once
func initOneSignalClient() {
	once.Do(func() {
		oneSignalClient = onesignal.NewClient(nil)
		oneSignalClient.AppKey = os.Getenv("ONESIGNAL_API")
	})
}

type NotificationRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func SendNotifToAll(c *fiber.Ctx) error {
	initOneSignalClient()

	var req NotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createRes, err := SendNotification(req.Title, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"response": createRes,
	})
}

func SendNotifWithID(title, content string, notificationToken string) (*onesignal.NotificationCreateResponse, *http.Response, error) {
	initOneSignalClient()
	nc := &onesignal.NotificationRequest{
		AppID:            os.Getenv("ONESIGNAL_APP"),
		Contents:         map[string]string{"en": content},
		Headings:         map[string]string{"en": title},
		IncludePlayerIDs: []string{notificationToken},
	}
	return oneSignalClient.Notifications.Create(nc)
}

func SendNotification(title, content string) (*onesignal.NotificationCreateResponse, error) {
	notificationReq := &onesignal.NotificationRequest{
		AppID:            os.Getenv("ONESIGNAL_APP"),
		Headings:         map[string]string{"en": title},
		Contents:         map[string]string{"en": content},
		IncludedSegments: []string{"All"},
	}

	createRes, _, err := oneSignalClient.Notifications.Create(notificationReq)
	if err != nil {
		fmt.Println("Error creating notification:", err)
		return nil, err
	}

	return createRes, nil
}
