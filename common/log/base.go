package log

import (
	"errors"
	"fmt"
	"github.com/kataras/golog"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

func init() {
	SetConfig(NewDefaultConfig())
}

var (
	lock           = sync.Mutex{}
	loggers        = make(map[string]*Logger)
	frameIgnored   = regexp.MustCompile(`(?)(github.com/kataras/golog)|(online/common/log/log.go)`)
	ErrUnknowLevel = errors.New("unknown log level")
)

const (
	DebugLevel = golog.DebugLevel
	InfoLevel  = golog.InfoLevel
	WarnLevel  = golog.WarnLevel
	ErrorLevel = golog.ErrorLevel
	FatalLevel = golog.FatalLevel
	PanicLevel = golog.FatalLevel
	TraceLevel = golog.DebugLevel
)

type Logger struct {
	*golog.Logger
	name string
}

func formatter(l *golog.Log, name string) bool {
	if strings.HasSuffix(strings.ToLower(name), ".yak") {
		name = name[:len(name)-4]
		l.Message = fmt.Sprintf("[%v] %v", name, l.Message)
		return false
	}

	// 使用 runtime.Caller(n) 如果 n 是硬编码的，就可能会获取到错误的行号
	// 比如 logger.Debug 和 logger.Debugf 等调用栈是不一样的
	// 还比如 main.go 中是直接使用了 log.DefaultLogger 而不是 GetLogger
	// 这里需要遍历一下栈信息，找到非 golog 的第一个 frame，然后作为 caller
	file := "???"
	line := 0

	pc := make([]uintptr, 64)
	n := runtime.Callers(3, pc)
	if n != 0 {
		pc = pc[:n]
		frames := runtime.CallersFrames(pc)

		for {
			frame, more := frames.Next()
			if !frameIgnored.MatchString(frame.File) {
				file = frame.File
				line = frame.Line
				break
			}
			if !more {
				break
			}
		}
	}

	slices := strings.Split(file, "/")
	file = slices[len(slices)-1]

	l.Message = fmt.Sprintf("[%s:%s:%d] %s", name, file, line, l.Message)

	return false
}

// GetLogger 返回一个新的 logger, 生产环境每个新的 logger 将会把日志记录到单独的文件中，
// 因此，期望的是整个主进程直接使用 DefaultLogger, 每个插件使用新的 logger 方便 debug
// 本函数调用的时候 因为都是全局变量的位置，所以配置文件还没加载，会返回一个设置了默认配置的 logger
// 之后加载配置文件的时候，会重新设置每个 logger 的属性
func GetLogger(name string) *Logger {
	lock.Lock()
	defer lock.Unlock()
	logger, exists := loggers[name]
	if exists {
		return logger
	} else {
		logger = &Logger{
			Logger: golog.New(),
			name:   name,
		}
		logger.Handle(func(l *golog.Log) bool {
			return formatter(l, name)
		})
		logger.SetTimeFormat("2006-01-02 15:04:05 -0700")
		logger.SetLevel(GetConfig().Level)
		loggers[name] = logger
		return logger
	}
}

func CheckLogDir(dir string) error {
	if dir == "" {
		return nil
	} else {
		testFilepath := path.Join(dir, "test-log-dir.test")
		defer os.Remove(testFilepath)
		return ioutil.WriteFile(testFilepath, []byte("test log file"), 0640)
	}
}

var DefaultLogger = GetLogger("default")

// Print prints a log message without levels and colors.
func Print(v ...interface{}) {
	DefaultLogger.Print(v...)
}

// Printf formats according to a format specifier and writes to `Printer#Output` without levels and colors.
func Printf(format string, args ...interface{}) {
	DefaultLogger.Printf(format, args...)
}

// Println prints a log message without levels and colors.
// It adds a new line at the end, it overrides the `NewLine` option.
func Println(v ...interface{}) {
	DefaultLogger.Println(v...)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

// Warn will print when logger's Level is warn, info or debug.
func Warn(v ...interface{}) {
	DefaultLogger.Warn(v...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Info will print when logger's Level is info or debug.
func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

// Debug will print when logger's Level is debug.
func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Trace is named after Debug
var (
	Trace  = Debug
	Tracef = Debugf
)

func SetLevel(level golog.Level) {
	DefaultLogger.Level = level
}

func GetLevel() golog.Level {
	return DefaultLogger.Level
}

func SetOutput(w io.Writer) {
	DefaultLogger.SetOutput(w)
}

func ParseLevel(raw string) (golog.Level, error) {
	disable := golog.Levels[golog.DisableLevel]
	if disable.Name == raw {
		return golog.DisableLevel, nil
	}
	for _, s := range disable.AlternativeNames {
		if raw == s {
			return golog.DisableLevel, nil
		}
	}
	level := golog.ParseLevel(raw)
	if level == golog.DisableLevel {
		return level, ErrUnknowLevel
	}
	return level, nil
}

func Warningf(raw string, args ...interface{}) {
	Warnf(raw, args...)
}
