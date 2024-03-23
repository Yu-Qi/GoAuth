package middleware

import (
	"net/http"

	em "emperror.dev/emperror"
	"emperror.dev/errors"
	ee "emperror.dev/errors"

	"github.com/sirupsen/logrus"

	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/gin-gonic/gin"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandlePanic that recovers from any panics and handles the error
func HandlePanic(c *gin.Context) {
	handleError := em.ErrorHandlerFunc(func(err error) {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Panic occurred")
		errTracer, ok := err.(stackTracer) // ok is false if errors doesn't implement stackTracer
		if ok {
			logrus.WithFields(logrus.Fields{
				"error": errTracer.StackTrace(),
			}).Error("stack trace")
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"code":    code.InternalUnknownError,
			"message": ee.WithStackDepth(err, 10).Error(),
		})
	})
	defer em.HandleRecover(handleError)
	c.Next()
}
