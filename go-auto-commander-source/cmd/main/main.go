package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/lutfailham96/go-auto-commander/internal/notificator"
	"github.com/lutfailham96/go-auto-commander/internal/spreadsheetmanager"

	"google.golang.org/api/sheets/v4"
)

// sendDiscordMessage sends a message to a Discord user.
func sendDiscordMessage(tokenBot string, recepientId string, message string) error {
	discord, err := notificator.NewDiscord(tokenBot, recepientId)
	err = discord.SendMessage(message)

	return err
}

// executeCommand executes a command and returns the output.
func executeOSCommand(command string) (string, error) {
	cmd := exec.Command(command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), err
}

// proceedSheetData proceeds the sheet data.
func proceedSheetData(tokenBot string, googleSpreadsheet *spreadsheetmanager.GoogleSpreadsheet) {
	if len(googleSpreadsheet.SheetContents.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		layout := "02/01/2006 15:04:05" // layout to parse the timestamp from the sheet
		currentDate := time.Now()

		// Iterate through the rows in the sheet.
		for i, row := range googleSpreadsheet.SheetContents.Values {
			// Skip the first row
			if i == 0 {
				continue
			}

			// Parse the timestamp from the sheet.
			timestamp, err := time.Parse(layout, row[0].(string))
			if err != nil {
				log.Fatalf("Unable to parse timestamp: %v", err)
			}

			// Compare the date part of the timestamp with the current date.
			if timestamp.Year() == currentDate.Year() && timestamp.Month() == currentDate.Month() && timestamp.Day() == currentDate.Day() {
				fullName := row[2].(string)
				projectName := row[3].(string)
				environment := row[4].(string)
				forwarderType := row[5].(string)
				discordId := row[8].(string)

				// Skip the row if the status is "DONE".
				if (len(row) > 9 && row[9].(string) == "DONE") || len(row) > 9 && row[9].(string) == "FAILED" { // assuming status is in column F
					fmt.Println("SKIP: ", fullName, projectName, environment, forwarderType, discordId)
					continue
				}

				// Execute the command.
				osCommand := fmt.Sprintf("./scripts/%s/%s/%s/generate.sh", projectName, environment, forwarderType)
				cmdOutput, err := executeOSCommand(osCommand)
				if err != nil {
					log.Println("execute-command-failed: ", err)
					sendDiscordMessage(tokenBot, discordId, fmt.Sprintf("Hello %s, permintaan akses PF anda ke %s %s %s tidak dapat diproses, harap untuk menghubungi PIC terkait", fullName, forwarderType, projectName, environment))
					updateRange := fmt.Sprintf("Form Responses 1!J%d", i+1) // assuming status is in column F
					valueRange := &sheets.ValueRange{
						Values: [][]interface{}{{"FAILED"}},
					}
					googleSpreadsheet.SheetService.Spreadsheets.Values.Update(googleSpreadsheet.SpreadsheetId, updateRange, valueRange).ValueInputOption("USER_ENTERED").Do()
					continue
				}
				log.Println("execute-command: ", osCommand)
				cmdOutput = fmt.Sprintf("Hello %s, berikut terkait akses PF %s %s %s\n%s", fullName, forwarderType, projectName, environment, cmdOutput)

				// Send message to discord
				err = sendDiscordMessage(tokenBot, discordId, cmdOutput)
				if err != nil {
					log.Println("discord-send-message-failed: ", err)
					continue
				}
				log.Println("discord-send-message: to ", fullName, " (", discordId, ")")

				// Update the status in the sheet to "DONE".
				updateRange := fmt.Sprintf("Form Responses 1!J%d", i+1) // assuming status is in column F
				valueRange := &sheets.ValueRange{
					Values: [][]interface{}{{"DONE"}},
				}
				_, err = googleSpreadsheet.SheetService.Spreadsheets.Values.Update(googleSpreadsheet.SpreadsheetId, updateRange, valueRange).ValueInputOption("USER_ENTERED").Do()
				if err != nil {
					log.Println("update-sheet-failed: ", err)
				}

				log.Println("PROCEED: ", fullName, projectName, environment, forwarderType, discordId)
			}
		}
	}
}

func main() {
	spreadsheetId := "YOUR-SPREADSHEET-ID"
	readRange := "Form Responses 1"
	gs := spreadsheetmanager.NewGoogleSpreadsheet(spreadsheetId, "./YOUR-SPREADSHEET-CREDENTIALS.json")
	_, err := gs.ReadData(readRange)
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	tokenBot := "YOUR-DISCORD-TOKEN-BOT"
	proceedSheetData(tokenBot, gs)
}
