package main

import "github.com/kataras/iris"
import "github.com/iris-contrib/middleware/cors"
import "fmt"

type Promotion struct {
    Id      string  `json:"id"`
}

type Statistics struct {
    Total   int     `json:"total"`
    Present int     `json:"present"`
    Absent  int     `json:"absent"`
    Late    int     `json:"late"`
    Waiting int     `json:"waiting"`
}

func main() {
	update := make(chan string)
    _promotion := Promotion{}
    _statistics := Statistics{}

    iris.Use(cors.Default())
	iris.Post("/update/statistics", func (ctx *iris.Context) {
        fmt.Println("ON POST")
        fmt.Println("PROMOTION: ", _promotion.Id)
        if err := ctx.ReadJSON(&_promotion); err != nil {
            ctx.JSON(400, map[string]string {"success": "false"})
        } else {
            if err := ctx.ReadJSON(&_statistics); err != nil {
                ctx.JSON(400, map[string]string {"success": "false"})
            } else {
                go func() {update <- "statistics"}()
        		ctx.JSON(200, map[string]string {"success": "true"})
            }
        }
	})

	iris.Config.Websocket.Endpoint = "/websocket"
	iris.Websocket.OnConnection(func (c iris.WebsocketConnection) {
        fmt.Println("ON CONNECTION")
        c.Join("statistics")

		continu := true
        go func() {
            for continu {
        		<-update
                fmt.Println("ON UPDATE")
                emit := fmt.Sprintf("statistics_%s", _promotion.Id)
                c.To("statistics").Emit(emit, _statistics)
                _promotion = Promotion{}
                _statistics = Statistics{}
        	}
    	}()

		c.OnDisconnect(func() {
			continu = false
		})
	})

	iris.Listen("127.0.0.1:8001")
}
