package handlers

import (
	"fmt"
	"net/http"

	"counter/internal/counter"

	"github.com/gin-gonic/gin"
)

var counterInstance = counter.NewObservableCount()

func HomePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func IncrementCounter(ctx *gin.Context) {
	counterInstance.Inc()
	ctx.JSON(http.StatusOK, gin.H{"response": "ok"})
}

func SSEHandler(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	clientGone := ctx.Writer.CloseNotify()

	observerCh, currentCount := counterInstance.Subscribe()
	defer counterInstance.Unsubscribe(observerCh)

	ctx.SSEvent("message", currentCount)
	ctx.Writer.Flush()

	for {
		select {
		case <-clientGone:
			fmt.Println("Client disconnected")
			return
		case newCount := <-observerCh:
			ctx.SSEvent("message", newCount)
			ctx.Writer.Flush()
			fmt.Println(newCount)
		}
	}
}
