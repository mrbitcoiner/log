package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type Log struct {
	init bool
	fd io.Writer
	level byte
	fNameMode byte
}

type CfgFunc func(*Log) error

const (
	LOGFATAL byte = 1
	LOGERR byte = 2
	LOGWARN byte = 3
	LOGINFO byte = 4
	LOGDEBUG byte = 5
	LOGTRACE byte = 6
)

const (
	FILENAME_MODE_NONE byte = 0
	FILENAME_MODE_FILEONLY byte = 1
	FILENAME_MODE_FILEPATH byte = 2
)

var (
	ErrInvalidLogLevel = errors.New("log: invalid log level")
	ErrLogFatal = errors.New("log: error log FATAL")
	ErrLog = errors.New("log error")
)

// NewLog create a new logger instance
func NewLog(cfgFuncs ...CfgFunc) (*Log, error) {
	log := defaultLog();
	err := error(nil)
	for _, f := range cfgFuncs {
		err = f(log)
		if err != nil { return nil, errors.Join(ErrLog, err) }
	}
	return log, nil
}

// Fatal prints the fatal log and panics
func (l *Log) Fatal(params ...interface{}) {
	l.log(LOGFATAL, params...)
	panic(ErrLogFatal)
}

// Fatalf prints the fatal log and panics
func (l *Log) Fatalf(fmtstr string, params ...interface{}) {
	l.logf(LOGFATAL, fmtstr, params...)
	panic(ErrLogFatal)
}

// Err prints the err log
func (l *Log) Err(params ...interface{}) {
	l.log(LOGERR, params...)
}

// Errf prints the err log
func (l *Log) Errf(fmtstr string, params ...interface{}) {
	l.logf(LOGERR, fmtstr, params...)
}

// Warn prints the warn log
func (l *Log) Warn(params ...interface{}) {
	l.log(LOGWARN, params...)
}

// Warnf prints the warn log
func (l *Log) Warnf(fmtstr string, params ...interface{}) {
	l.logf(LOGWARN, fmtstr, params...)
}

// Info prints the info log
func (l *Log) Info(params ...interface{}) {
	l.log(LOGINFO, params...)
}

// Infof prints the info log
func (l *Log) Infof(fmtstr string, params ...interface{}) {
	l.logf(LOGINFO, fmtstr, params...)
}

// Debug prints the debug log
func (l *Log) Debug(params ...interface{}) {
	l.log(LOGDEBUG, params...)
}

// Debugf prints the debug log
func (l *Log) Debugf(fmtstr string, params ...interface{}) {
	l.logf(LOGDEBUG, fmtstr, params...)
}

// Trace prints the trace log
func (l *Log) Trace(params ...interface{}) {
	l.log(LOGTRACE, params...)
}

// Tracef prints the trace log
func (l *Log) Tracef(fmtstr string, params ...interface{}) {
	l.logf(LOGTRACE, fmtstr, params...)
}

func (l *Log) log(lvl byte, params ...interface{}) {
	toLog := l.toLog(lvl)
	if !toLog { return }
	msg := fmt.Sprintln(params...)
	l.logWrite(lvl, msg)
}

func (l *Log) logf(lvl byte, fmtstr string, params ...interface{}) {
	toLog := l.toLog(lvl)
	if !toLog { return }
	msg := fmt.Sprintf(fmtstr, params...)
	l.logWrite(lvl, msg)
}

func (l *Log) logWrite(lvl byte, msg string) {
	log := ""
	switch l.fNameMode {
	case FILENAME_MODE_NONE:
		log = fmt.Sprintf(
			"%s |%s| %s",
			time.Now().Format(time.DateTime),
			lvlStr(lvl),
			msg,
		)
	case FILENAME_MODE_FILEONLY:
		_, file, line, ok := runtime.Caller(3)
		slashOffset := strings.LastIndex(file, "/")
		file = string([]byte(file)[slashOffset+1:])
		if !ok {
			fmt.Fprintln(os.Stderr, "log: err getting caller")
		}
		log = fmt.Sprintf(
			"%s |%s| %s:%d: %s",
			time.Now().Format(time.DateTime),
			lvlStr(lvl),
			file,
			line,
			msg,
		)
	case FILENAME_MODE_FILEPATH:
		_, file, line, ok := runtime.Caller(3)
		if !ok {
			fmt.Fprintln(os.Stderr, "log: err getting caller")
		}
		log = fmt.Sprintf(
			"%s |%s| %s:%d: %s",
			time.Now().Format(time.DateTime),
			lvlStr(lvl),
			file,
			line,
			msg,
		)
	default:
		log = fmt.Sprintf(
			"%s |%s| %s",
			time.Now().Format(time.DateTime),
			lvlStr(lvl),
			msg,
		)
	}
	n, err := l.fd.Write([]byte(log))
	if err != nil {
		fmt.Fprintln(os.Stderr, "log: err writing log")
	}
	if n != len(log) {
		fmt.Fprintln(os.Stderr, "log: err unexpected bytes written:", len(log), n)
	}
}

func (l *Log) toLog(lvl byte) bool {
	return l.level >= lvl
}

func isValid(lvl byte) bool {
	return lvl >= LOGFATAL && lvl <= LOGTRACE
}

func defaultLog() *Log {
	return &Log{
		fd: os.Stderr,
		level: LOGINFO,
		fNameMode: FILENAME_MODE_FILEONLY,
	}
}

func lvlStr(lvl byte) string {
	switch lvl {
	case LOGFATAL:
		return "FATAL"
	case LOGERR:
		return "ERR"
	case LOGWARN:
		return "WARN"
	case LOGINFO:
		return "INFO"
	case LOGDEBUG:
		return "DEBUG"
	case LOGTRACE:
		return "TRACE"
	}
	return ""
}

func WithWriter(writer io.Writer) CfgFunc {
	return func(l *Log) error {
		l.fd = writer
		return nil
	}
}

func WithLevel(level byte) CfgFunc {
	return func(l *Log) error {
		valid := isValid(level)
		if !valid { return ErrInvalidLogLevel }
		l.level = level
		return nil
	}
}

func WithFileName(l *Log) error {
	l.fNameMode = FILENAME_MODE_FILEONLY
	return nil
}

func WithFilePath(l *Log) error {
	l.fNameMode = FILENAME_MODE_FILEPATH
	return nil
}
