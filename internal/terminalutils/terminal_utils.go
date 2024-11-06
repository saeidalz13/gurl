package terminalutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// Reset
	FormatReset = "\033[0m"

	// Regular Colors
	RegularBlack  = "\033[0;30m" // Black
	RegularRed    = "\033[0;31m" // Red
	RegularGreen  = "\033[0;32m" // Green
	RegularYellow = "\033[0;33m" // Yellow
	RegularBlue   = "\033[0;34m" // Blue
	RegularPurple = "\033[0;35m" // Purple
	RegularCyan   = "\033[0;36m" // Cyan
	RegularWhite  = "\033[0;37m" // White

	// Bold
	BoldBlack  = "\033[1;30m" // Bold Black
	BoldRed    = "\033[1;31m" // Bold Red
	BoldGreen  = "\033[1;32m" // Bold Green
	BoldYellow = "\033[1;33m" // Bold Yellow
	BoldBlue   = "\033[1;34m" // Bold Blue
	BoldPurple = "\033[1;35m" // Bold Purple
	BoldCyan   = "\033[1;36m" // Bold Cyan
	BoldWhite  = "\033[1;37m" // Bold White

	// Underline
	UnderlineBlack  = "\033[4;30m" // Underline Black
	UnderlineRed    = "\033[4;31m" // Underline Red
	UnderlineGreen  = "\033[4;32m" // Underline Green
	UnderlineYellow = "\033[4;33m" // Underline Yellow
	UnderlineBlue   = "\033[4;34m" // Underline Blue
	UnderlinePurple = "\033[4;35m" // Underline Purple
	UnderlineCyan   = "\033[4;36m" // Underline Cyan
	UnderlineWhite  = "\033[4;37m" // Underline White

	// Background
	BackgroundBlack  = "\033[40m" // Background Black
	BackgroundRed    = "\033[41m" // Background Red
	BackgroundGreen  = "\033[42m" // Background Green
	BackgroundYellow = "\033[43m" // Background Yellow
	BackgroundBlue   = "\033[44m" // Background Blue
	BackgroundPurple = "\033[45m" // Background Purple
	BackgroundCyan   = "\033[46m" // Background Cyan
	BackgroundWhite  = "\033[47m" // Background White

	// High Intensity
	IntenseBlack  = "\033[0;90m" // High Intensity Black
	IntenseRed    = "\033[0;91m" // High Intensity Red
	IntenseGreen  = "\033[0;92m" // High Intensity Green
	IntenseYellow = "\033[0;93m" // High Intensity Yellow
	IntenseBlue   = "\033[0;94m" // High Intensity Blue
	IntensePurple = "\033[0;95m" // High Intensity Purple
	IntenseCyan   = "\033[0;96m" // High Intensity Cyan
	IntenseWhite  = "\033[0;97m" // High Intensity White

	// Bold High Intensity
	BoldIntenseBlack  = "\033[1;90m" // Bold High Intensity Black
	BoldIntenseRed    = "\033[1;91m" // Bold High Intensity Red
	BoldIntenseGreen  = "\033[1;92m" // Bold High Intensity Green
	BoldIntenseYellow = "\033[1;93m" // Bold High Intensity Yellow
	BoldIntenseBlue   = "\033[1;94m" // Bold High Intensity Blue
	BoldIntensePurple = "\033[1;95m" // Bold High Intensity Purple
	BoldIntenseCyan   = "\033[1;96m" // Bold High Intensity Cyan
	BoldIntenseWhite  = "\033[1;97m" // Bold High Intensity White

	// High Intensity Backgrounds
	BackgroundIntenseBlack  = "\033[0;100m" // High Intensity Background Black
	BackgroundIntenseRed    = "\033[0;101m" // High Intensity Background Red
	BackgroundIntenseGreen  = "\033[0;102m" // High Intensity Background Green
	BackgroundIntenseYellow = "\033[0;103m" // High Intensity Background Yellow
	BackgroundIntenseBlue   = "\033[0;104m" // High Intensity Background Blue
	BackgroundIntensePurple = "\033[0;105m" // High Intensity Background Purple
	BackgroundIntenseCyan   = "\033[0;106m" // High Intensity Background Cyan
	BackgroundIntenseWhite  = "\033[0;107m" // High Intensity Background White
)

func PrintWsError(errMsg string) {
	fmt.Printf("%s[ERROR]:%s %s", BoldRed, FormatReset, errMsg)
}

func PrintWsServerMsg(msg string) {
	fmt.Printf("%s[SERVER]:%s %s\n", BoldCyan, FormatReset, msg)
}

func PrintWsClientMsg(msg string) {
	fmt.Printf("%s[CLIENT]:%s %s\n", BoldGreen, FormatReset, msg)
}

func PrintHTTPClientInfo(ip, httpRequest string) {
	fmt.Printf("%s\n[To Server] >>%s\n", BoldWhite, FormatReset)

	fmt.Printf("%s\nServer Details%s\n", BoldPurple, FormatReset)
	fmt.Println("---------------------")
	fmt.Printf("%sServer IP:%s %s\n", RegularPurple, FormatReset, ip)
	// Other details
	fmt.Print("\n")

	fmt.Printf("%sRequest%s\n", BoldGreen, FormatReset)
	fmt.Println("---------------------")
	fmt.Print(httpRequest)

	// If the request didn't have body, it wouldn't
	// end with new line character.
	if httpRequest[len(httpRequest)-1] != '\n' {
		fmt.Print("\n\n")
	}
	fmt.Printf("%s[From Server] <<%s\n", BoldWhite, FormatReset)
}

func PrintWebSocketClientInfo(ip, wsRequest string) {
	fmt.Printf("%s\n[To Server] >>%s\n", BoldWhite, FormatReset)

	fmt.Printf("%s\nDetails%s\n", BoldPurple, FormatReset)
	fmt.Println("---------------------")
	fmt.Printf("%sServer IP:%s %s\n", RegularPurple, FormatReset, ip)
	// Other details
	fmt.Print("\n")

	fmt.Printf("%sRequest%s\n", BoldGreen, FormatReset)
	fmt.Println("---------------------")
	fmt.Print(wsRequest)

	fmt.Printf("%s[From Server] <<%s\n\n", BoldWhite, FormatReset)
}

func PrintAppWarning(msg string) {
	fmt.Printf("%s[WARNING]:%s %s\n", BoldYellow, FormatReset, msg)
}

func PrintAppError(msg string) {
	fmt.Printf("%s[ERROR]:%s %s\n", BoldRed, FormatReset, msg)
}

func GetWsInputFromStdin() []byte {
	// If we use fmt.Scanln(), then it only reads
	// the characters until the space. bufio lets
	// us consider all the characters until the
	// delimiter we decide. We choose '\n' that
	// shows the end of the input.
	reader := bufio.NewReader(os.Stdin)

	for {
		rawInput, err := reader.ReadString('\n')
		if err != nil {
			PrintWsError(err.Error())
			continue
		}

		// Trim spaces and newlines from the input
		rawInput = strings.TrimSpace(rawInput)

		// Check if the input is empty or contains only spaces
		if len(rawInput) == 0 {
			PrintWsError("empty input!")
			continue
		}

		// Remove spaces within the input
		modInput := make([]byte, 0, len(rawInput))
		for _, b := range []byte(rawInput) {
			if b != ' ' {
				modInput = append(modInput, b)
			}
		}

		return modInput
	}
}
