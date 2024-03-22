package middleware

import (
	em "emperror.dev/emperror"
	"emperror.dev/errors"
	ee "emperror.dev/errors"

	"github.com/Yu-Qi/GoAuth/api/response"
	"github.com/Yu-Qi/GoAuth/pkg/log"
	"github.com/gin-gonic/gin"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandlePanic that recovers from any panics and handles the error
func HandlePanic(c *gin.Context) {
	handleError := em.ErrorHandlerFunc(func(err error) {
		log.Error(err.Error())
		errTracer, ok := err.(stackTracer) // ok is false if errors doesn't implement stackTracer
		if ok {
			log.ErrorWithData("stack trace", errTracer.StackTrace())
		}
		response.AbortByAny(c, ee.WithStackDepth(err, 10))
	})
	defer em.HandleRecover(handleError)
	c.Next()
}
