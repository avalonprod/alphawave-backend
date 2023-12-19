package v1

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initAiChatRoutes(api *gin.RouterGroup) {
	ai := api.Group("/ai")
	{
		// authenticated := ai.Group("/", h.userIdentity, h.setTeamSessionFromCookie)
		{
			ai.POST("/new-message", h.newMessage)
		}
	}
}

type streamContent struct {
	Content string `json:"content"`
}

func (h *HandlerV1) newMessage(c *gin.Context) {

	var input []types.Message

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	message, err := h.service.AiChatService.NewMessage(input)
	if err != nil {
		logger.Errorf("error gateway. err: %v", err)
		newResponse(c, http.StatusBadGateway, errors.New("error gateway").Error())

		return
	}
	chanStream := make(chan streamContent)
	go func() {
		defer close(chanStream)

		for {
			stream, err := message.Stream.Recv()
			if errors.Is(err, io.EOF) {
				c.Status(http.StatusOK)

				return
			}
			if err != nil {
				logger.Errorf("error stream data. err: %v", err)
				newResponse(c, http.StatusBadGateway, errors.New("error gateway").Error())

				return
			}
			if len(stream.Choices) != 0 {
				chanStream <- streamContent{Content: stream.Choices[0].Delta.Content}
			}
		}
	}()
	defer message.Stream.Close()
	// c.SSEvent("start", message.Role)
	c.Stream(func(w io.Writer) bool {

		if msg, ok := <-chanStream; ok {
			c.SSEvent("message", msg)
			return true
		}
		c.SSEvent("end", "")
		return false
	})
}
