package spreadsheetmanager

import (
	"context"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// GoogleSpreadsheet represents a Google Sheets spreadsheet.
type GoogleSpreadsheet struct {
	SpreadsheetId  string
	CredentialFile string
	SheetService   *sheets.Service
	SheetContents  *sheets.ValueRange
}

// createClient creates a new Google Sheets client.
func createClient(credentialFile string) (*http.Client, error) {
	// Read the credentials file.
	ctx := context.Background()
	data, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		return nil, err
	}

	// Create a client.
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	client := conf.Client(ctx)

	return client, err
}

// createSheetService creates a new Google Sheets service.
func createSheetService(credentialFile string) (*sheets.Service, error) {
	// Create a client.
	client, err := createClient(credentialFile)
	if err != nil {
		return nil, err
	}

	// Create a sheets service.
	srv, err := sheets.New(client)

	return srv, err

}

// NewGoogleSpreadsheet creates a new Google Sheets spreadsheet.
func NewGoogleSpreadsheet(spreadsheetId, credentialFile string) *GoogleSpreadsheet {
	return &GoogleSpreadsheet{
		SpreadsheetId:  spreadsheetId,
		CredentialFile: credentialFile,
	}
}

// ReadData reads data from the Google Sheets spreadsheet.
func (g *GoogleSpreadsheet) ReadData(readRange string) (*sheets.ValueRange, error) {
	// Create a new Google Sheets service.
	srv, err := createSheetService(g.CredentialFile)
	if err != nil {
		return nil, err
	}
	g.SheetService = srv

	// Read data from the spreadsheet.
	resp, err := srv.Spreadsheets.Values.Get(g.SpreadsheetId, readRange).Do()
	g.SheetContents = resp

	return resp, err
}
