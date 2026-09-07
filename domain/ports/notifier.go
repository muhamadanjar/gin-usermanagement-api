package ports

// Notifier is the port for pushing notifications to devices/topics.
// Implementations: pkg/firebase FCM client (structural typing, no adapter
// needed).
type Notifier interface {
	SendToDevice(token, title, body string, data map[string]string) (string, error)
	SendToDevices(tokens []string, title, body string, data map[string]string) (int, error)
	SendToTopic(topic, title, body string, data map[string]string) (string, error)
}
