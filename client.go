package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewClient(ctx context.Context, credentials, token string) (*http.Client, error) {
	b, err := ioutil.ReadFile(credentials)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}

	var t *oauth2.Token
	if _, err := os.Stat(token); err == nil {
		f, err := os.Open(token)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		t = new(oauth2.Token)
		if err = json.NewDecoder(f).Decode(t); err != nil {
			return nil, err
		}
		if t.Expiry.After(time.Now()) {
			return config.Client(ctx, t), nil
		}
	}

	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	_ = browser.OpenURL(url)
	fmt.Printf("Access required: to enable access for your account use the sign-in page opened in your browser or the following link:\n%s\n", url)
	ch := make(chan string)
	go func(ch chan<- string) {
		defer close(ch)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			ch <- r.URL.Query().Get("code")
			fmt.Fprint(w, "You can close this tab.")
		})
		if err := http.ListenAndServe("localhost:8192", nil); err != nil {
			log.Fatalln("Cannot start server:", err)
		}
	}(ch)
	if t, err = config.Exchange(ctx, <-ch); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(token, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(t); err != nil {
		return nil, err
	}
	return config.Client(ctx, t), nil
}
