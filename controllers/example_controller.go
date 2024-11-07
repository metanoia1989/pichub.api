package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/models"
	"pichub.api/repository"
)

func GetData(ctx *gin.Context) {
	var example []*models.Example
	repository.Get(&example)
	ctx.JSON(http.StatusOK, &example)

}
