package tools

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

const stage string = "dev"

func PrintLogo() {
	levelColor := color.New(color.FgCyan)
	levelColor.Println(`
 /$$$$$$$$ /$$     /$$                                 /$$$$$$$                               
| $$_____/| $$    | $$                                | $$__  $$                              
| $$     /$$$$$$  | $$$$$$$   /$$$$$$   /$$$$$$       | $$  \ $$  /$$$$$$   /$$$$$$   /$$$$$$ 
| $$$$$ |_  $$_/  | $$__  $$ /$$__  $$ /$$__  $$      | $$  | $$ /$$__  $$ /$$__  $$ /$$__  $$
| $$__/   | $$    | $$  \ $$| $$$$$$$$| $$  \__/      | $$  | $$| $$  \__/| $$  \ $$| $$  \ $$
| $$      | $$ /$$| $$  | $$| $$_____/| $$            | $$  | $$| $$      | $$  | $$| $$  | $$
| $$$$$$$$|  $$$$/| $$  | $$|  $$$$$$$| $$            | $$$$$$$/| $$      |  $$$$$$/| $$$$$$$/
|________/ \___/  |__/  |__/ \_______/|__/            |_______/ |__/       \______/ | $$____/ 
                                                                                    | $$      
                                                                                    | $$      
                                                                                    |__/      
`)

	levelColor.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")

	levelColor = color.New(color.FgRed)
	levelColor.Println("[!] All risks are your responsibility. This tool is intended for educational purposes and to make your life easier.....")
}

func Logger(level, message string) {
	message = strings.ReplaceAll(message, "\n", "")
	message = strings.ReplaceAll(message, "\r", "")

	level = strings.ToLower(level)
	var levelColor *color.Color

	switch level {
	case "info":
		levelColor = color.New(color.FgWhite)
		levelColor.Println(fmt.Sprintf("📢  %s", message))
	case "error":
		if stage == "dev" {
			_, file, line, ok := runtime.Caller(1)
			if ok {
				levelColor := color.New(color.FgRed).SprintFunc()

				fmt.Println(levelColor(fmt.Sprintf("☠️  %s", message)))
				fmt.Println(levelColor(fmt.Sprintf("☠️  Error Path: %s, Line %d", file, line)))
			} else {
				fmt.Println("Tidak dapat mendapatkan informasi file dan baris kode.")
			}
		} else {
			levelColor = color.New(color.FgRed)
			levelColor.Println(fmt.Sprintf("☠️  %s", message))
		}
	case "success":
		levelColor = color.New(color.FgGreen)
		levelColor.Println(fmt.Sprintf("✅  %s", message))
	case "warning":
		levelColor = color.New(color.FgYellow)
		levelColor.Println(fmt.Sprintf("⚠️  %s", message))
	case "input":
		levelColor = color.New(color.FgCyan)
		levelColor.Printf("⌨️  %s", message)
	default:
		levelColor = color.New(color.FgWhite)
		levelColor.Println(fmt.Sprintf("[%s] %s", level, message))
	}
}
