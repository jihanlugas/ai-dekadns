package middleware

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	logger "ai-dekadns/helper"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LoggerToFile() gin.HandlerFunc {
	// log file
	fileName := path.Join("logging", "lintasarta.log")
	// write file
	//src, err := os.OpenFile("lintasarta.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	src, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}
	// instantiation
	logger := logrus.New()
	// Set output
	logger.Out = src
	// Set log level
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		// start time
		startTime := time.Now()
		// Processing request
		c.Next()
		// Stop time
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqUri := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()
		// Log format
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

func LoggerToElastic(client *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		requestBody, _ := io.ReadAll(c.Request.Body)
		reader := io.NopCloser(bytes.NewBuffer(requestBody))
		c.Request.Body = reader

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// start time
		startTime := time.Now()

		// Processing request
		c.Next()

		// Response Body
		responseBody := bodyLogWriter.body.String()
		// Stop time
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqUri := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()

		go logger.SendLoggingToElastic(client, logger.ElasticLog{
			Request: logger.ElasticRequestLog{
				RequestTime: startTime,
				Method:      reqMethod,
				Uri:         reqUri,
				Proto:       c.Request.Proto,
				UserAgent:   c.Request.UserAgent(),
				Referer:     c.Request.Referer(),
				PostData:    string(requestBody),
				ClientIP:    clientIP,
				LatencyTime: latencyTime,
			},
			Response: logger.ElasticResponseLog{
				ResponseTime: endTime,
				StatusCode:   statusCode,
				ResponseBody: responseBody,
			},
		})
	}
}
