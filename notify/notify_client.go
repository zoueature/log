package notify

type Notifier interface {
	SendMarkdown(title, content string, atMobile ...string) error
}
