# Sheet Sync

Use data from google sheets and template files to keep your content updated.

## Installation

You can install using the command `go get klaidliadon.dev/sheet-sync`

Obtain a new credential file from the [quickstart page](https://developers.google.com/sheets/api/quickstart/go).

1. Click on the **Enable Google Sheets API** button.
2. Choose a Project name.
3. Select **Web server** and add `http://localhost:8192` as callback.
4. Click on the **DOWNLOAD CLIENT CONFIGURATION** button.

Place the downloaded file in `~/.config/sheet-sync/credentials.json`.

## Usage

You can now run the `sheet-sync` command with the following arguments:

1. **Speadsheet ID**: the unique identifier of the source Google Sheet.
2. **Range**: the cell range to use as a source (e.g. `Sheet1!A1:B6`).
3. **Template**: the template file to use with the selected data.
4. **Output**: the destination file for the output.

