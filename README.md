# LinkedIn Profile Search and Filtering Tool

This is a command-line tool that searches for LinkedIn profiles based on a hashtag and filters the results based on job title and experience. The filtered profiles are then added to a new sheet in Google Sheets.
Prerequisites
Go 1.16 or later
A LinkedIn account
A Google Cloud Platform project with the Google Sheets API enabled
A Google service account with access to the Google Sheets API
Getting Started
Clone this repository:
```
git clone https://github.com/yourusername/linkedin-profile-search.git
```

Install the dependencies:
```
go mod download
```

Create a configuration file named config.yaml in the root directory of the project. Here's an example configuration file:
yaml
```
client_id: <your LinkedIn app client ID>
client_secret: <your LinkedIn app client secret>
redirect_uri: <your LinkedIn app redirect URI>
hashtag: <the hashtag to search for>
experience_years: <the minimum number of years of experience>
job_title: <the job title to filter by>
spreadsheet_id: <the ID of the Google Sheets spreadsheet to add the filtered profiles to>
sheet_name: <the name of the sheet in the Google Sheets spreadsheet to add the filtered profiles to>
```

Replace the placeholders in the configuration file with your own values:
client_id, client_secret, and redirect_uri: These are the credentials for your LinkedIn app. You can create a new app in the LinkedIn Developer Console.
hashtag: This is the hashtag to search for on LinkedIn.
experience_years: This is the minimum number of years of experience that a profile must have to be included in the results.
job_title: This is the job title to filter the results by.
spreadsheet_id and sheet_name: These are the ID and name, respectively, of the Google Sheets spreadsheet to add the filtered profiles to.
Run the tool:
```
go run main.go
```

Follow the instructions to authenticate with LinkedIn and authorize the tool to access your profile data.
Wait for the tool to search for profiles, filter the results, and add the filtered profiles to the specified sheet in Google Sheets. The tool will print a message when it's done.
Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request.
License
This project is licensed under the MIT License. See the LICENSE file for details.
