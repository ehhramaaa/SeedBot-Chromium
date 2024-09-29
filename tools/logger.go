package tools

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrintLogo() {
	levelColor := color.New(color.FgCyan)
	levelColor.Println(`
  /$$$$$$                            /$$ /$$$$$$$              /$$    
 /$$__  $$                          | $$| $$__  $$            | $$    
| $$  \__/  /$$$$$$   /$$$$$$   /$$$$$$$| $$  \ $$  /$$$$$$  /$$$$$$  
|  $$$$$$  /$$__  $$ /$$__  $$ /$$__  $$| $$$$$$$  /$$__  $$|_  $$_/  
 \____  $$| $$$$$$$$| $$$$$$$$| $$  | $$| $$__  $$| $$  \ $$  | $$    
 /$$  \ $$| $$_____/| $$_____/| $$  | $$| $$  \ $$| $$  | $$  | $$ /$$
|  $$$$$$/|  $$$$$$$|  $$$$$$$|  $$$$$$$| $$$$$$$/|  $$$$$$/  |  $$$$/
 \______/  \_______/ \_______/ \_______/|_______/  \______/    \___/  
`)

	levelColor.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")

	levelColor = color.New(color.FgRed)
	levelColor.Println("[!] All risks are your responsibility. This tool is intended for educational purposes and to make your life easier.....")
}

func Logger(level, message string) {
	level = strings.ToLower(level)
	var levelColor *color.Color

	switch level {
	case "info":
		levelColor = color.New(color.FgWhite)
	case "error":
		levelColor = color.New(color.FgRed)
	case "success":
		levelColor = color.New(color.FgGreen)
	case "warning":
		levelColor = color.New(color.FgYellow)
	default:
		levelColor = color.New(color.FgWhite)
	}

	if level == "input" {
		levelColor.Printf("[+] %s", message)
	} else {
		levelColor.Println(fmt.Sprintf("[*] %s", message))
	}
}
