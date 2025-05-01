package constants

import "fmt"

type NotificationMessage struct {
    Title   string
    Content string
}

type PropertyNotifications struct {
    NewLeaseApplication NotificationMessage
    LeaseApproved      NotificationMessage
    LeaseRejected      NotificationMessage
}

var Property = PropertyNotifications{
    NewLeaseApplication: NotificationMessage{
        Title:   "шинэ түрээсийн хүсэлт!",
        Content: "You have a new lease application for your property: %s",
    },
    LeaseApproved: NotificationMessage{
        Title:   "Түрээсийн хүсэлт зөвшөөрөгдсөн",
        Content: "The lease application for %s has been approved",
    },
    LeaseRejected: NotificationMessage{
        Title:   "Түрээсийн хүсэлт rejected",
        Content: "The lease application for %s has been rejected",
    },
}

func (n NotificationMessage) Format(params ...interface{}) (string, string) {
    return n.Title, fmt.Sprintf(n.Content, params...)
}