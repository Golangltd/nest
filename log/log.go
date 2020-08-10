package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	tracker *zap.Logger
	gLogger *Logger
	eLogger *Logger
)

// levels
const (
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	fatalLevel = 4
)

const (
	printDebugLevel = "[debug  ] "
	printInfoLevel  = "[release] " //这里只是为了和leaf的日志格式保持一致
	printWarnLevel  = "[warn   ] "
	printErrorLevel = "[error  ] "
	printFatalLevel = "[fatal  ] "
)

type Logger struct {
	level      int32
	baseLogger *log.Logger
	baseFile   *os.File
	fileLogger *log.Logger
}

func (logger *Logger) init(strLevel string, flag int, dir string, isError bool) {
	logger.initLevel(strLevel, isError)
	logger.initBaseLogger(flag, isError)
	logger.initFileLogger(flag, dir, isError)
}

func (logger *Logger) initLevel(strLevel string, isError bool) {
	if isError {
		logger.level = warnLevel
	} else {
		logger.level = getLevel(strLevel)
	}
}

func (logger *Logger) initBaseLogger(flag int, isError bool) {
	var out io.Writer
	if isError {
		out = os.Stderr
	} else {
		out = os.Stdout
	}

	logger.baseLogger = log.New(out, "", flag)
}

func (logger *Logger) initFileLogger(flag int, dir string, isError bool) {
	if dir == "" {
		return
	}
	path := logger.getFilePath(dir, isError)
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	logger.fileLogger = log.New(f, "", flag)
}

func (logger *Logger) getFilePath(dir string, isError bool) string {
	dataStr := logger.getDateStr()

	var logType string
	if isError {
		logType = "error"
	} else {
		logType = "base"
	}

	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	return fmt.Sprintf("%s%s_%s.log", dir, dataStr, logType)
}

func (logger *Logger) getDateStr() string {
	t := time.Now()
	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
}

func getLevel(level string) int32 {
	switch strings.ToLower(level) {
	case "debug":
		return debugLevel
	case "info":
		return infoLevel
	case "warn":
		return warnLevel
	case "error":
		return errorLevel
	case "fatal":
		return fatalLevel
	default:
		return debugLevel
	}
}

func newLogger(strLevel string, flag int, dir string) (*Logger, *Logger) {
	logger := new(Logger)
	logger.init(strLevel, flag, dir, false)

	errLogger := new(Logger)
	errLogger.init(strLevel, flag, dir, true)

	return logger, errLogger
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		_ = logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int32, printLevel string, format string, a ...interface{}) {
	if level < atomic.LoadInt32(&logger.level) {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	_ = logger.baseLogger.Output(4, fmt.Sprintf(format, a...))

	if logger.fileLogger != nil {
		_ = logger.fileLogger.Output(4, fmt.Sprintf(format, a...))
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) SetLevel(level int32) {
	atomic.StoreInt32(&logger.level, level)
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.doPrintf(infoLevel, printInfoLevel, format, a...)
}

func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.doPrintf(warnLevel, printWarnLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Debug(format string, a ...interface{}) {
	gLogger.Debug(format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.Info(format, a...)
}

func Warn(format string, a ...interface{}) {
	gLogger.Warn(format, a...)
	eLogger.Warn(format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.Error(format, a...)
	eLogger.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.Fatal(format, a...)
	eLogger.Fatal(format, a...)
}

func Track(msg string, fields ...zap.Field) {
	tracker.Info(msg, fields...)
}

//这里是为了和python中track中用的格式保持一致
func luckyTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000000"))
}

func ReloadLogger(level string) {
	gLogger.SetLevel(getLevel(level))
}

//使用log之前必须初始化
func InitLogger(level string, dir string) {
	if gLogger != nil && tracker != nil {
		return
	}

	if dir != "" {
		mkdir(dir)
	}

	gLogger, eLogger = newLogger(level, log.Ldate|log.Lmicroseconds|log.Lshortfile, dir)

	//zap的logger
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "@timestamp"
	encoderCfg.EncodeTime = luckyTimeEncoder
	rawJSON := []byte(fmt.Sprintf(`{
	  "level": "%s",
	  "encoding": "json",
	  "outputPaths": ["%s/track.json"],
	  "encoderConfig": {
	    "levelEncoder": "uppercase"
	  },
      "disableCaller": true
	}`, level, dir))
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = encoderCfg
	tracker, _ = cfg.Build()
}

//关闭服务器之前调用，同步缓冲区
func CloseLogger() {
	_ = tracker.Sync()
}

func isDirExists(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func mkdir(dir string) {
	if isDirExists(dir) {
		return
	}
	_ = os.MkdirAll(dir, os.ModePerm)
}
