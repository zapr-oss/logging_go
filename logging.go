package logging

import (
	"github.com/zapr-oss/logging_go/hook"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type LogConfig struct {
	Path            string `json:"path"`
	Level           string `json:"level"`
	MaxSizeInMb     int    `json:"maxSizeInMb"`
	MaxBackups      int    `json:"maxBackups"`
	MaxAgeInDays    int    `json:"maxAgeInDays"`
	GZipCompress    bool   `json:"gZipCompress"`
	IDEF            bool   `json:"isDifferentErrorFile"`
	Formatter       string `json:"formatter"`
	Env             string `json:"env"`
	ShouldSetCaller bool   `json:"shouldSetCaller"`
}

func GetLogger(name string) *log.Logger {

	logConf, err := loadConfiguration()

	if err != nil {
		log.Warn("Error loading log config file")
	}

	logger := log.New()
	logger.Level = log.DebugLevel

	if logConf != nil {

		// Setting formatter accordingly, using config. Default is text so we are not explicitly setting it.
		if strings.ToLower(logConf.Formatter) == "json" {
			logger.SetFormatter(new(log.JSONFormatter))
		} else if logConf.Formatter != "text" {
			log.Println("Formatter Unknown: Setting to default formatter `text`")
		}

		logger.SetReportCaller(logConf.ShouldSetCaller)

		if strings.ToLower(logConf.Env) == "dev" || strings.ToLower(logConf.Env) == "local" {
			log.Println("Loglevel is default set to DEBUG for environments `dev` and `local`")
			logger.Level = log.DebugLevel
		} else {
			if strings.EqualFold(logConf.Level, "trace") {
				logger.Level = log.TraceLevel
			} else if strings.EqualFold(logConf.Level, "debug") {
				logger.Level = log.DebugLevel
			} else if strings.EqualFold(logConf.Level, "info") {
				logger.Level = log.InfoLevel
			} else if strings.EqualFold(logConf.Level, "warn") {
				logger.Level = log.WarnLevel
			} else if strings.EqualFold(logConf.Level, "error") {
				logger.Level = log.ErrorLevel
			} else if strings.EqualFold(logConf.Level, "fatal") {
				logger.Level = log.FatalLevel
			} else if strings.EqualFold(logConf.Level, "panic") {
				logger.Level = log.PanicLevel
			}
			setupLogs(logger, name, logConf)
		}

		log.Println("LogLevel set to: ", logger.Level.String())
	}

	return logger
}

func getExecutableFileDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex)
}

func loadConfiguration() (*LogConfig, error) {

	logConf := &LogConfig{}

	fileDir := getExecutableFileDirectory()

	filePathLink := fileDir + "/log_config.json"

	log.Println("Checking File Path: ", filePathLink)

	configFile, err := os.Open(filePathLink)

	if err != nil {
		filePathLink = fileDir + "/resources/log_config.json"

		configFile, err = os.Open(filePathLink)

		if err != nil {
			log.Printf("Didn't find file at: %v OR %v", fileDir+"/log_config.json", filePathLink)
			return nil, err
		}
	}

	log.Println("Log Config found at File Path: ", filePathLink)

	defer configFile.Close()

	logFileBytes, err := ioutil.ReadAll(configFile)

	if err != nil {
		log.Fatal("Error reading log file")
	}

	err = json.Unmarshal(logFileBytes, logConf);
	return logConf, err
}

// setupLogs adds hooks to send logs to different destinations depending on level
func setupLogs(logger *log.Logger, fileName string, logConf *LogConfig) {

	if !logConf.IDEF {
		logPath := logConf.Path + fileName + ".log"
		logger.SetOutput(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    logConf.MaxSizeInMb, // megabytes
			MaxBackups: logConf.MaxBackups,
			MaxAge:     logConf.MaxAgeInDays, //days
			Compress:   logConf.GZipCompress, // disabled by default
		})

		log.Println("Redirecting all LEVEL logs to file: ", logPath)
		return
	}

	logger.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	errorLogPath := logConf.Path + fileName + ".err"
	logger.AddHook(&hook.WriterHook{ // Send logs with level higher than warning to stderr
		Writer: &lumberjack.Logger{
			Filename:   errorLogPath,
			MaxSize:    logConf.MaxSizeInMb, // megabytes
			MaxBackups: logConf.MaxBackups,
			MaxAge:     logConf.MaxAgeInDays, //days
			Compress:   logConf.GZipCompress, // disabled by default
		},
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.Println("Redirecting all WARN,ERROR,FATAL,PANIC logs to file: ", errorLogPath)

	infoLogPath := logConf.Path + fileName + ".log"
	logger.AddHook(&hook.WriterHook{ // Send info and debug logs to stdout
		Writer: &lumberjack.Logger{
			Filename:   infoLogPath,
			MaxSize:    logConf.MaxSizeInMb, // megabytes
			MaxBackups: logConf.MaxBackups,
			MaxAge:     logConf.MaxAgeInDays, //days
			Compress:   logConf.GZipCompress, // disabled by default
		},
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
		},
	})

	log.Println("Redirecting all DEBUG,INFO logs to file: ", infoLogPath)

}
