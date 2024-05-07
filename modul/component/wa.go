package component

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func SendWhatsAppConfirmation(message, phonenumber string) error {
	url := "https://api.wa.my.id/api/send/message/text"

	jsonStr := []byte(`{
        "to": "` + phonenumber + `",
        "isgroup": false,
        "messages": "` + message + `"
    }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", os.Getenv("WA_TOKEN"))
	// req.Header.Set("Token", "v4.public.eyJleHAiOiIyMDI0LTAyLTE5VDIxOjA3OjM2WiIsImlhdCI6IjIwMjQtMDEtMjBUMjE6MDc6MzZaIiwiaWQiOiI2MjgyMzE3MTUwNjgxIiwibmJmIjoiMjAyNC0wMS0yMFQyMTowNzozNloiff1YQuHHPwSzGpisAMb9rTLP58-jCqtByzePJACBLghprkq2HXtTSbVTShc49m3GIVkU42VSl8uSGme8c4vXnQc")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	return nil
}
