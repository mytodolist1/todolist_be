package component

import (
	"time"
)

func AddReminder(username, phonenumber, title string, remindAt time.Time) error {
	duration := time.Until(remindAt)
	
	time.Sleep(duration)

	message := `Hai ` + username + `\n\nTugas kamu ` + title + ` telah melewati batas waktu, Segera selesaikan tugas kamu!`

	SendWhatsAppConfirmation(message, phonenumber)

	return nil
}
