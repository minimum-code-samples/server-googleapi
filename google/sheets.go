package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// ReadSheetsTitles extracts the titles of the sheets from the supplied array.
func ReadSheetsTitles(ss []*sheets.Sheet) []string {
	titles := make([]string, len(ss))
	for i, s := range ss {
		fmt.Printf("- Title: %s\n", s.Properties.Title)
		fmt.Printf("- ID: %v\n", s.Properties.SheetId)
		titles[i] = s.Properties.Title
	}
	return titles
}

// FetchSpreadsheetSheets gets the data pertaining to the sheets in a particular spreadsheet.
func FetchSpreadsheetSheets(ctx context.Context, spreadsheetID string, credentials []byte, token *oauth2.Token) ([]*sheets.Sheet, error) {
	srv, err := makeSheetsService(ctx, credentials, token)
	if err != nil {
		return nil, err
	}
	call := srv.Spreadsheets.Get(spreadsheetID)
	if err != nil {
		return nil, err
	}
	resp, err := call.Fields("sheets.properties").Do()
	if err != nil {
		return nil, err
	}
	return resp.Sheets, nil
}

// FetchSpreadsheetValues retrieves all values of a single spreadsheet from Google Sheets.
func FetchSpreadsheetValues(ctx context.Context, spreadsheetID, sheetName string, credentials []byte, token *oauth2.Token) ([][]interface{}, error) {
	srv, err := makeSheetsService(ctx, credentials, token)
	if err != nil {
		return nil, err
	}
	rg := "'" + sheetName + "'!B2:M"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rg).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func makeSheetsService(ctx context.Context, credentials []byte, token *oauth2.Token) (*sheets.Service, error) {
	cfg, err := MakeConfig(credentials, ScopesWithSheets())
	if err != nil {
		return nil, err
	}
	srv, err := sheets.NewService(ctx, option.WithTokenSource(cfg.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
