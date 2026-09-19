package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"practice/api/models"
	"practice/api/shared/commonfun"

	"github.com/gin-gonic/gin"
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

func SaveActivitiesInDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		incomingRequestTime := time.Now()
		contentType := c.GetHeader("Content-Type")

		var requestBody string
		if strings.Contains(contentType, "multipart/form-data") {
			requestBody = "MULTIPART_FORM_DATA"
		} else if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body again
		}

		// ---------- RESPONSE WRITER ----------
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		c.Next()
		outgoingRequestTime := time.Now()

		logData := models.ActivityLog{
			Method:              c.Request.Method,
			Path:                c.FullPath(),
			Query:               c.Request.URL.RawQuery,
			IP:                  c.ClientIP(),
			UserAgent:           c.Request.UserAgent(),
			Authorization:       c.GetHeader("Authorization"),
			Origin:              c.GetHeader("Origin"),
			Referer:             c.GetHeader("Referer"),
			Headers:             c.Request.Header,
			RequestBody:         requestBody,
			StatusCode:          c.Writer.Status(),
			ResponseBody:        blw.body.String(),
			LatencySeconds:      outgoingRequestTime.Sub(incomingRequestTime).Seconds(),
			IncomingRequestTime: incomingRequestTime,
			OutgoingRequestTime: outgoingRequestTime,
		}
		go commonfun.SaveActivity(logData) // save into db asynchronously
	}
}
