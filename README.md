## Overview

This is a simple go logger that can be used to log logs of different package into different file. 

The simplest way to use this package is to import and then get logger by name. Include log_config.json file in either default **_project-directory_** or **_project-directory/resources_**
```go
package main

import (
	"github.com/zapr-oss/logging_go"
	"github.com/sirupsen/logrus"
	)

var log = logging.GetLogger("handler")

func main() {
        log.WithFields(logrus.Fields{
            "name":       "John",
            "id":         10}).
            Errorln("error getting data")
}
```

This will create a file name handler.log in the directory given in log_config.json. 
If you want to use the same file for logging all logs, log object can be passed to all functions and packages.

```json
{
  "path": "/home/ubuntu/content-inventory-scripts/ott_content_merger/content_merger/logs/",
  "level": "info",  //getting log level
  "maxSizeInMb": 50,
  "maxBackups": 5,
  "maxAgeInDays": 30,
  "gZipCompress": true, // gzip compress the backups
  "isDifferentErrorFile": true, // will log error logs to file with .err extention
  "formatter": "text", //setting formatter. (text/json)
  "env": "prod", //this is used to set local/dev env to always use debug logs.
  "shouldSetCaller": true //used to print the line at which logging happened
}
```