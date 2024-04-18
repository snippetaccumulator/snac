package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type VSTColor string

const (
	ResetColor             VSTColor = "\033[0m"
	BlackForeground        VSTColor = "\033[30m"
	RedForeground          VSTColor = "\033[31m"
	GreenForeground        VSTColor = "\033[32m"
	YellowForeground       VSTColor = "\033[33m"
	BlueForeground         VSTColor = "\033[34m"
	PurpleForeground       VSTColor = "\033[35m"
	CyanForeground         VSTColor = "\033[36m"
	WhiteForeground        VSTColor = "\033[37m"
	GreyForeground         VSTColor = "\033[90m"
	BrightRedForeground    VSTColor = "\033[91m"
	BrightGreenForeground  VSTColor = "\033[92m"
	BrightYellowForeground VSTColor = "\033[93m"
	BrightBlueForeground   VSTColor = "\033[94m"
	BrightPurpleForeground VSTColor = "\033[95m"
	BrightCyanForeground   VSTColor = "\033[96m"
	BrightWhiteForeground  VSTColor = "\033[97m"

	BlackBackground        VSTColor = "\033[40m"
	RedBackground          VSTColor = "\033[41m"
	GreenBackground        VSTColor = "\033[42m"
	YellowBackground       VSTColor = "\033[43m"
	BlueBackground         VSTColor = "\033[44m"
	PurpleBackground       VSTColor = "\033[45m"
	CyanBackground         VSTColor = "\033[46m"
	WhiteBackground        VSTColor = "\033[47m"
	GreyBackground         VSTColor = "\033[100m"
	BrightRedBackground    VSTColor = "\033[101m"
	BrightGreenBackground  VSTColor = "\033[102m"
	BrightYellowBackground VSTColor = "\033[103m"
	BrightBlueBackground   VSTColor = "\033[104m"
	BrightPurpleBackground VSTColor = "\033[105m"
	BrightCyanBackground   VSTColor = "\033[106m"
	BrightWhiteBackground  VSTColor = "\033[107m"
)

type LogLevel int

const (
	ERROR LogLevel = iota
	WARN
	SUCCESS
	INFO
	DEBUG
)

var appliedLevel LogLevel

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case SUCCESS:
		return "SUCCESS"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l LogLevel) ForegroundColor() VSTColor {
	switch l {
	case DEBUG:
		return GreyForeground
	case INFO:
		return BrightBlueForeground
	case WARN:
		return BlackForeground
	case SUCCESS:
		return BrightGreenForeground
	case ERROR:
		return BlackForeground
	default:
		return BrightPurpleForeground
	}
}

func (l LogLevel) BackgroundColor() VSTColor {
	switch l {
	case DEBUG:
		return BlackBackground
	case INFO:
		return BlackBackground
	case WARN:
		return BrightYellowBackground
	case SUCCESS:
		return BlackBackground
	case ERROR:
		return BrightRedBackground
	default:
		return BlackBackground
	}
}

func log(level LogLevel, format string, a ...any) {
	fmt.Printf("%s%s%-7s [%s]%s: %s\n", level.BackgroundColor(), level.ForegroundColor(), level, time.Now().Format("2006-01-02 15:04:05"), ResetColor, fmt.Sprintf(format, a...))
}

func FromString(level string) LogLevel {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "SUCCESS":
		return SUCCESS
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

func SetLevel(level LogLevel) {
	appliedLevel = level
}

func Debug(format string, a ...any) {
	if appliedLevel >= DEBUG {
		log(DEBUG, format, a...)
	}
}

func Info(format string, a ...any) {
	log(INFO, format, a...)
}

func Warn(format string, a ...any) {
	log(WARN, format, a...)
}

func Success(format string, a ...any) {
	log(SUCCESS, format, a...)
}

func Err(exit bool, err error) {
	if err == nil {
		return
	}
	Error(exit, err.Error())
}

func Error(exit bool, format string, a ...any) {
	log(ERROR, format, a...)
	if exit {
		os.Exit(1)
	}
}
