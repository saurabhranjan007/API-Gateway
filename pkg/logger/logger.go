package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logging *zap.Logger

// InitLogger initializes the logger
func InitLogger() {
	log.Println("Initializing logger!!")

	rootPath, err := GetRootDirectoryPath()
	if err != nil {
		log.Fatalf("Error getting root directory path: %s", err.Error())
	} else {
		log.Printf("Root directory path: %s\n", rootPath)
	}

	// Check for log directory (if not, create)
	logDirectory := fmt.Sprintf("%s/.logs", rootPath)
	info, err := EnsureDirectory(logDirectory)
	if err != nil {
		log.Fatalf("Error creating log directory: %s", err.Error())
	} else {
		log.Printf("Log directory ensured: %s\n", logDirectory)
	}

	// Create log file based on current date (if exists, append)
	logPath := fmt.Sprintf("%s/log_%s.log", logDirectory, time.Now().Format("2006-01-02"))
	if info != nil {
		logPath = fmt.Sprintf("%s/log_%s.log", logDirectory, info.ModTime().Format("2006-01-02"))
	}

	if err := EnsureLogFile(logPath); err != nil {
		log.Fatalf("Error creating log file: %s", err.Error())
	} else {
		log.Printf("Log file ensured: %s\n", logPath)
	}

	// Custom log configuration
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    100, // MB
			MaxBackups: 3,
			MaxAge:     365, // days
			Compress:   false,
			LocalTime:  true,
		}),
		zapcore.DebugLevel,
	)

	// Zap logger core setup
	Logging = zap.New(core, zap.AddCaller())
}

// SyncLogger flushes any buffered log entries
func SyncLogger() {
	if Logging != nil {
		Logging.Sync()
	}
}

func LogInfo(source, activity, debugString string, object ...interface{}) {
	if Logging == nil {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	Logging.Info(debugString,
		zap.String("Source", source),
		zap.Any("Object", object),
		zap.String("Activity", activity),
		zap.String("Caller", caller),
	)
}

func LogError(source string, activity string, object interface{}, err error) {
	if Logging == nil {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	Logging.Error("Error",
		zap.String("Source", source),
		zap.Any("Object", object),
		zap.String("Activity", activity),
		zap.String("Caller", caller),
		zap.Error(err),
	)
}

func LogFatal(source string, activity string, object interface{}, err error) {
	if Logging == nil {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	Logging.Fatal("Fatal",
		zap.String("Source", source),
		zap.Any("Object", object),
		zap.String("Activity", activity),
		zap.String("Caller", caller),
		zap.Error(err),
	)
}

func LogWarning(source string, activity string, message string, object interface{}) {
	if Logging == nil {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	Logging.Warn("Warning",
		zap.String("Source", source),
		zap.Any("Object", object),
		zap.String("Activity", activity),
		zap.String("Caller", caller),
		zap.String("Message", message),
	)
}

// GetModuleDirectoryPath returns the directory of the current module
func GetModuleDirectoryPath() (string, error) {

	executablePath, err := os.Executable()
	if err != nil {
		log.Printf("Error getting executable path: %v\n", err)
		return "", err
	}

	return filepath.Dir(executablePath), nil
}

// GetRootDirectoryPath returns the root directory of the project
func GetRootDirectoryPath() (string, error) {

	rootPath, err := filepath.Abs(".")
	if err != nil {
		log.Printf("Error getting root directory path: %v\n", err)
		return "", err
	}

	return rootPath, nil
}

// EnsureDirectory checks if a directory exists, and creates it if it does not
func EnsureDirectory(directoryPath string) (os.FileInfo, error) {

	info, err := os.Stat(directoryPath)

	if os.IsNotExist(err) {
		log.Printf("Directory does not exist, creating: %s\n", directoryPath)
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v\n", err)
			return nil, err
		}
		info, err = os.Stat(directoryPath)
	}

	if err != nil {
		log.Printf("Error stating directory: %v\n", err)
	}

	return info, err
}

// EnsureLogFile checks if a log file exists, and creates it if it does not
func EnsureLogFile(path string) error {

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		log.Printf("Log file does not exist, creating: %s\n", path)
		file, err := os.Create(path)
		if err != nil {
			log.Printf("Error creating log file: %v\n", err)
			return fmt.Errorf("failed to create log file: %v", err)
		}
		file.Close()
	} else if err != nil {
		log.Printf("Error stating log file: %v\n", err)
	}

	return nil
}
