package api

import (
	"backend/app/form"
	"backend/app/repository"
	"backend/db"
	err2 "backend/utils/err"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ApplyEntryAPI(app *gin.RouterGroup, resource *db.Resource) {
	EntryEntity := repository.NewEntryEntity(resource)

	entryRoute := app.Group("")
	entryRoute.GET("/indexing", indexing(EntryEntity))

	entryRoute.POST("/kurs", createEntry(EntryEntity))
	entryRoute.PUT("/kurs", updateEntry(EntryEntity))
	entryRoute.DELETE("/kurs/:date", deleteEntry(EntryEntity))

	entryRoute.GET("/kurs/startdate=:startdate/enddate=:enddate", getEntryBetweenDates(EntryEntity))
	entryRoute.GET("/kurs/startdate=:startdate/enddate=:enddate/:symbol", getEntryListOneSymbol(EntryEntity))

	entryRoute.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func indexing(entryEntity repository.IEntry) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		entry, code, err := entryEntity.Indexing()
		response := map[string]interface{}{
			"entry": entry,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func createEntry(entryEntity repository.IEntry) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		entryRequest := form.Entry{}
		if err := ctx.Bind(&entryRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		entry, code, err := entryEntity.CreateOneEntry(entryRequest)
		response := map[string]interface{}{
			"entry": entry,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func getEntryBetweenDates(entryEntity repository.IEntry) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		startDate := ctx.Param("startdate")
		endDate := ctx.Param("enddate")

		if endDate == "" {
			endDate = time.Now().Format("2006-01-02")
		}
		list, code, err := entryEntity.GetEntryListByDate(startDate, endDate)
		response := map[string]interface{}{
			"entry": list,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func getEntryListOneSymbol(entryEntity repository.IEntry) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		startDate := ctx.Param("startdate")
		endDate := ctx.Param("enddate")
		symbol := ctx.Param("symbol")
		list, code, err := entryEntity.GetEntryListBySymbol(startDate, endDate, symbol)
		response := map[string]interface{}{
			"entry": list,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func updateEntry(entryEntity repository.IEntry) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		entryRequest := form.Entry{}
		if err := ctx.Bind(&entryRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		entry, code, err := entryEntity.UpdateEntry(entryRequest.Date, entryRequest)
		response := map[string]interface{}{
			"entry": entry,
			"err":   err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func deleteEntry(entryEntity repository.IEntry) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		date := ctx.Param("date")
		entry, code, err := entryEntity.DeleteEntry(date)
		response := map[string]interface{}{
			"entry": entry,
			"err":   err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
