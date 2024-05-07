package component

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

// saves file to github
func SaveFileToGithub(usernameGhp, emailGhp, repoGhp, path string, r *http.Request) (string, error) {
	// ambil file dari form
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("error 1: %s", err)
	}
	defer file.Close()

	// generate random file name
	randomFileName, err := generateRandomFileName(handler.Filename)
	if err != nil {
		return "", fmt.Errorf("error 2: %s", err)
	}

	// cek apakah file sudah ada
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error 5: %s", err)
	}

	// get access token
	access_token := os.Getenv("GITHUB_TOKEN")
	if access_token == "" {
		return "", fmt.Errorf("error access token: %s", err)
	}

	// save file to github
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access_token},
	)
	tc := oauth2.NewClient(r.Context(), ts)
	client := github.NewClient(tc)

	_, _, err = client.Repositories.CreateFile(r.Context(), usernameGhp, repoGhp, path+"/"+randomFileName, &github.RepositoryContentFileOptions{
		Message:   github.String("Add new file"),
		Content:   fileContent,
		Committer: &github.CommitAuthor{Name: github.String(usernameGhp), Email: github.String(emailGhp)},
	})
	if err != nil {
		return "", fmt.Errorf("error 6: %s", err)
	}

	// get file url
	fileUrl := "https://" + usernameGhp + ".github.io/" + repoGhp + "/" + path + "/" + randomFileName

	return fileUrl, nil
}

func generateRandomFileName(originalFilename string) (string, error) {
	// generate random file name
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomFileName := fmt.Sprintf("%x%s", randomBytes, filepath.Ext(originalFilename))
	return randomFileName, nil
}
