package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"google.golang.org/api/sheets/v4"
)

var (
	credPath  string
	tokenPath string
)

var cfg struct {
	Spreadsheet string
	Range       string
	Template    string
	Output      string
}

func init() {
	if len(os.Args) != 5 {
		fmt.Println(os.Args[0], "<speadsheet_id>", "<range>", "<template>", "<output>")
		os.Exit(1)
	}
	cfg.Spreadsheet = os.Args[1]
	cfg.Range = os.Args[2]
	cfg.Template = os.Args[3]
	cfg.Output = os.Args[4]

	u, err := user.Current()
	if err != nil {
		log.Fatalf("Cannot get current user: %v", err)
	}
	baseDir := filepath.Join(u.HomeDir, ".config", "sheet-sync")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("Cannot create configuration directory: %v", err)
	}
	credPath = filepath.Join(baseDir, "credentials.json")
	tokenPath = filepath.Join(baseDir, "token.json")
}

func main() {
	client, err := NewClient(credPath, tokenPath)
	if err != nil {
		log.Fatalf("Cannot get client: %v", err)
	}

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(cfg.Spreadsheet, cfg.Range).Do()
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := template.ParseFiles(cfg.Template)
	if err != nil {
		log.Fatal(err)
	}
	dst, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()
	if err := tpl.Execute(dst, resp.Values); err != nil {
		log.Fatal(err)
	}
}
