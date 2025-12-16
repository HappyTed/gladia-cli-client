package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel uint8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type OutputTypeFlag uint

const (
	STD OutputTypeFlag = 1 << iota
	FILE
)

type ILogger interface {
	Debug(a ...any)
	Info(a ...any)
	Warn(a ...any)
	Error(a ...any)
	Fatal(a ...any)
	DebugF(format string, a ...any)
	InfoF(format string, a ...any)
	WarnF(format string, a ...any)
	ErrorF(format string, a ...any)
	FatalF(format string, a ...any)
	Close() error
}

type Logger struct {
	level           LogLevel
	output          io.Writer
	outputCloseFunc func() error
}

func DefaultLogger() ILogger {
	return &Logger{
		level:           DEBUG,
		output:          os.Stdout,
		outputCloseFunc: func() error { return nil },
	}
}

func NewLogger(level LogLevel, outputType OutputTypeFlag, outputPath string) ILogger {

	// default logger
	l := &Logger{
		level:           level,
		output:          os.Stdout,
		outputCloseFunc: func() error { return nil },
	}

	// set output flags
	var outs []io.Writer
	if outputType&STD != 0 {
		outs = append(outs, os.Stdout)
	}
	if outputType&FILE != 0 {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to init logger: failed open log file %s: %w", outputPath, err))
		}
		l.outputCloseFunc = file.Close
		outs = append(outs, file)
	}

	l.output = io.MultiWriter(outs...)
	log.SetOutput(l.output)

	// output format
	log.SetFlags(log.Ldate | log.Ltime)
	return l
}

func (l *Logger) Close() error {
	return l.outputCloseFunc()
}

func (l *Logger) Debug(a ...any) {
	if l.level == DEBUG {
		log.Println("[DEBUG]", fmt.Sprint(a...))
	}
}

func (l *Logger) DebugF(format string, a ...any) {
	if l.level == DEBUG {
		log.Printf("[DEBUG] "+format, a...)
	}
}

func (l *Logger) Info(a ...any) {
	if l.level <= INFO {
		log.Println("[INFO]", fmt.Sprint(a...))
	}
}

func (l *Logger) InfoF(format string, a ...any) {
	if l.level <= INFO {
		log.Printf("[INFO] "+format, a...)
	}
}

func (l *Logger) Warn(a ...any) {
	if l.level <= WARN {
		log.Println("[WARN]", fmt.Sprint(a...))
	}
}

func (l *Logger) WarnF(format string, a ...any) {
	if l.level <= WARN {
		log.Printf("[WARN] "+format, a...)
	}
}

func (l *Logger) Error(a ...any) {
	if l.level <= ERROR {
		log.Println("[ERROR]", fmt.Sprint(a...))
	}
}

func (l *Logger) ErrorF(format string, a ...any) {
	if l.level <= ERROR {
		log.Printf("[ERROR] "+format, a...)
	}
}

func (l *Logger) Fatal(a ...any) {
	if l.level <= FATAL {
		log.Println("[FATAL]", fmt.Sprint(a...))
	}
}

func (l *Logger) FatalF(format string, a ...any) {
	if l.level <= FATAL {
		log.Printf("[FATAL] "+format, a...)
	}
}

func TestLogger() {

	// debug
	lvl := DEBUG
	debug := NewLogger(lvl, STD|FILE, "logs/debug.log")

	defer debug.Close()
	// not format
	fmt.Println()
	fmt.Println("log level:", lvl)
	fmt.Println()
	debug.Debug("THIS IS _DEBUG_ FROM - ", lvl)
	debug.Info("THIS IS _INFO_ FROM - ", lvl)
	debug.Warn("THIS IS _WARN_ FROM - ", lvl)
	debug.Error("THIS IS _ERROR_ FROM - ", lvl)
	debug.Fatal("THIS IS _FATAL_ FROM - ", lvl)
	// format string
	fmt.Println()
	fmt.Println("format string")
	fmt.Println()
	debug.DebugF("THIS IS _DEBUG_ FROM - %d", lvl)
	debug.InfoF("THIS IS _INFO_ FROM - %d", lvl)
	debug.WarnF("THIS IS _WARN_ FROM - %d", lvl)
	debug.ErrorF("THIS IS _ERROR_ FROM - %d", lvl)
	debug.FatalF("THIS IS _FATAL_ FROM - %d", lvl)
	fmt.Println("---")

	// info
	lvl = INFO
	info := NewLogger(lvl, STD|FILE, "logs/info.log")

	defer info.Close()
	// not format
	fmt.Println()
	fmt.Println("log level:", lvl)
	fmt.Println()
	info.Debug("THIS IS _DEBUG_ FROM - ", lvl)
	info.Info("THIS IS _INFO_ FROM - ", lvl)
	info.Warn("THIS IS _WARN_ FROM - ", lvl)
	info.Error("THIS IS _ERROR_ FROM - ", lvl)
	info.Fatal("THIS IS _FATAL_ FROM - ", lvl)
	// format string
	fmt.Println()
	fmt.Println("format string")
	fmt.Println()
	info.DebugF("THIS IS _DEBUG_ FROM - %d", lvl)
	info.InfoF("THIS IS _INFO_ FROM - %d", lvl)
	info.WarnF("THIS IS _WARN_ FROM - %d", lvl)
	info.ErrorF("THIS IS _ERROR_ FROM - %d", lvl)
	info.FatalF("THIS IS _FATAL_ FROM - %d", lvl)
	fmt.Println("---")

	// warn
	lvl = WARN
	warn := NewLogger(lvl, STD|FILE, "logs/warn.log")

	defer warn.Close()
	// not format
	fmt.Println()
	fmt.Println("log level:", lvl)
	fmt.Println()
	warn.Debug("THIS IS _DEBUG_ FROM - ", lvl)
	warn.Info("THIS IS _INFO_ FROM - ", lvl)
	warn.Warn("THIS IS _WARN_ FROM - ", lvl)
	warn.Error("THIS IS _ERROR_ FROM - ", lvl)
	warn.Fatal("THIS IS _FATAL_ FROM - ", lvl)
	// format string
	fmt.Println()
	fmt.Println("format string")
	fmt.Println()
	warn.DebugF("THIS IS _DEBUG_ FROM - %d", lvl)
	warn.InfoF("THIS IS _INFO_ FROM - %d", lvl)
	warn.WarnF("THIS IS _WARN_ FROM - %d", lvl)
	warn.ErrorF("THIS IS _ERROR_ FROM - %d", lvl)
	warn.FatalF("THIS IS _FATAL_ FROM - %d", lvl)
	fmt.Println("---")

	// error
	lvl = ERROR
	erlog := NewLogger(lvl, STD|FILE, "logs/error.log")

	defer erlog.Close()
	// not format
	fmt.Println()
	fmt.Println("log level:", lvl)
	fmt.Println()
	erlog.Debug("THIS IS _DEBUG_ FROM - ", lvl)
	erlog.Info("THIS IS _INFO_ FROM - ", lvl)
	erlog.Warn("THIS IS _WARN_ FROM - ", lvl)
	erlog.Error("THIS IS _ERROR_ FROM - ", lvl)
	erlog.Fatal("THIS IS _FATAL_ FROM - ", lvl)
	// format string
	fmt.Println()
	fmt.Println("format string")
	fmt.Println()
	erlog.DebugF("THIS IS _DEBUG_ FROM - %d", lvl)
	erlog.InfoF("THIS IS _INFO_ FROM - %d", lvl)
	erlog.WarnF("THIS IS _WARN_ FROM - %d", lvl)
	erlog.ErrorF("THIS IS _ERROR_ FROM - %d", lvl)
	erlog.FatalF("THIS IS _FATAL_ FROM - %d", lvl)
	fmt.Println("---")

	// fatal
	lvl = FATAL
	fatal := NewLogger(lvl, STD|FILE, "logs/fatal.log")

	defer fatal.Close()
	// not format
	fmt.Println()
	fmt.Println("log level:", lvl)
	fmt.Println()
	fatal.Debug("THIS IS _DEBUG_ FROM - ", lvl)
	fatal.Info("THIS IS _INFO_ FROM - ", lvl)
	fatal.Warn("THIS IS _WARN_ FROM - ", lvl)
	fatal.Error("THIS IS _ERROR_ FROM - ", lvl)
	fatal.Fatal("THIS IS _FATAL_ FROM - ", lvl)
	// format string
	fmt.Println()
	fmt.Println("format string")
	fmt.Println()
	fatal.DebugF("THIS IS _DEBUG_ FROM - %d", lvl)
	fatal.InfoF("THIS IS _INFO_ FROM - %d", lvl)
	fatal.WarnF("THIS IS _WARN_ FROM - %d", lvl)
	fatal.ErrorF("THIS IS _ERROR_ FROM - %d", lvl)
	fatal.FatalF("THIS IS _FATAL_ FROM - %d", lvl)
	fmt.Println("---")
}
