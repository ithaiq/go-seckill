package core

import "github.com/gin-gonic/gin"

type IMid interface {
	OnRequest(context *gin.Context) error
}
