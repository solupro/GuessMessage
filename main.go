package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/olahol/melody"
	"github.com/labstack/gommon/log"
	"strconv"
)

func main()  {
	r := gin.Default()
	m := melody.New()
	log.SetLevel(log.DEBUG)
	user := make(map[int64]*melody.Session)

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.POST("/push", func(c *gin.Context) {
		strId := c.PostForm("userId")
		message := c.PostForm("message")

		userId, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": 0,
				"code": -3,
				"msg": err.Error(),
			})

			return
		}

		msg := "没有找到用户"
		if s, ok := user[userId]; ok {
			log.Debug("send userId:" + strconv.FormatInt(userId, 10) + " message:" + message)
			m.BroadcastMultiple([]byte(message), []*melody.Session{s})
			msg = "成功发送"
		}

		c.JSON(http.StatusOK, gin.H{
			"result": 1,
			"code": 0,
			"msg": msg,
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		userId := _getUserIdFromSession(s)
		if 0 != userId {
			log.Debug("add:" + strconv.FormatInt(userId, 10))
			user[userId] = s
		}
	})

	m.HandleClose(func(se *melody.Session, i int, s string) error {
		userId := _getUserIdFromSession(se)
		if 0 != userId {
			log.Debug("remove:" + strconv.FormatInt(userId, 10))
			delete(user, userId)
		}

		return nil
	})

	r.Run(":5000")
}

func _getUserIdFromSession(s *melody.Session) int64 {
	s.Request.ParseForm()
	strId := s.Request.Form.Get("user_id")

	userId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		log.Warn(err.Error())
		return 0
	}

	return userId
}