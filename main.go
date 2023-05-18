package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Profile struct {
	ID             string `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Headline       string `json:"headline"`
	PublicProfileUrl string `json:"publicProfileUrl"`
	Positions      []struct {
		Title      string `json:"title"`
		StartDate  struct {
			Month int `json:"month"`
			Year  int `json:"year"`
		} `json:"startDate"`
		EndDate struct {
			Month int `json:"month"`
			Year  int `json:"year"`
		} `json:"endDate"`
	} `json:"positions"`
}

func main() {
	// Read the configuration file using Viper
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}
	client_id := viper.GetString("client_id")
	client_secret := viper.GetString("client_secret")
	hashtag := viper.GetString("hashtag")
	experience_years := viper.GetInt("experience_years")
	job_title := viper.GetString("job_title")
	spreadsheet_id := viper.GetString("spreadsheet_id")
	sheet_name := viper.GetString("sheet_name")

	// Set up the LinkedIn provider
	linkedinProvider := linkedin.New(client_id, client_secret, "http://localhost:3000/auth/linkedin/callback", linkedin.ScopeEmail, linkedin.ScopeBasicProfile, linkedin.ScopeReadWriteWShare)
	goth.UseProviders(linkedinProvider)

	// Configure the OAuth2 client
	conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint:     linkedinProvider.Endpoint(),
		Scopes:       []string{linkedin.ScopeEmail, linkedin.ScopeBasicProfile, linkedin.ScopeReadWriteWShare},
		RedirectURL:  "http://localhost:3000/auth/linkedin/callback",
	}

	// Obtain an access token using the OAuth2 client
	token, err := conf.PasswordCredentialsToken(context.Background(), "", url.Values{
		"grant_type": {"client_credentials"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	accessToken := token.AccessToken

	profiles, err := searchProfiles(accessToken, hashtag)
	if err != nil {
		fmt.Println(err)
		return
	}

	filteredProfiles := filterProfiles(profiles, job_title, experience_years)

	err = addProfilesToGoogleSheets(filteredProfiles, accessToken, spreadsheet_id, sheet_name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Profiles added to Google Sheets.")
}

func searchProfiles(accessToken string, hashtag string) ([]Profile, error) {
	searchParams := url.Values{}
	searchParams.Set("q", "hashtag:"+hashtag)
	searchParams.Set("count", "10")
	searchParams.Set("start", "0")
	searchParams.Set("fields", "id,firstName,lastName,headline,publicProfileUrl,positions:(title,timePeriod)")
	searchUrl := "https://api.linkedin.com/v2/search?q=" + searchParams.Get("q") + "&count=" + searchParams.Get("count") + "&start=" + searchParams.Get("start") + "&fields=" + searchParams.Get("fields")
	searchHeaders := map[string]string{"Authorization": "Bearer " + accessToken}
	searchResponse, err := doRequest("GET", searchUrl, searchHeaders, nil)
	if err != nil {
		return nil, err
	}
	var profiles struct {
		Elements []Profile `json:"elements"`
	}
	err = json.Unmarshal([]byte(searchResponse), &profiles)
	if err != nil {
		return nil, err
	}
	return profiles.Elements, nil
}

func filterProfiles(profiles []Profile, jobTitle string, experienceYears int) []Profile {
	var filteredProfiles []Profile
	for _, profile := range profiles {
		hasJobTitle := false
		hasExperience := false
		for _, position := range profile.Positions {
			if strings.Contains(strings.ToLower(position.Title), strings.ToLower(jobTitle)) {
				hasJobTitle = true
			}
			startDate := position.StartDate.Year*12 + position.StartDate.Month
			endDate := position.EndDate.Year*12 + position.EndDate.Month
			duration := (endDate - startDate) / 12
			if duration >= experienceYears {
				hasExperience = true
			}
		}
		if hasJobTitle && hasExperience {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}
	return filteredProfiles
}

func addProfilesToGoogleSheets(profiles []Profile, accessToken string, spreadsheetId string, sheetName string) error {
	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		return err
	}
	sheetsService, err := sheets.NewService(ctx, creds)
	if err != nil {
		return err
	}
	sheetRange := sheetName + "!A1:E"
	var valueRange sheets.ValueRange
	var rows [][]interface{}
	for _, profile := range profiles {
		row := []interface{}{profile.ID, profile.FirstName, profile.LastName, profile.Headline, profile.PublicProfileUrl}
		rows = append(rows, row)
	}
	valueRange.Values = rows
	_, err = sheetsService.Spreadsheets.Values.Append(spreadsheetId, sheetRange, &valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}
	return nil
}
