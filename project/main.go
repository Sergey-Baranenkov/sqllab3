package main

import (
	"context"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"log"
)

var (
	Router = router.New()
	PostgresConn *pgx.Conn
	onTrueDatabase = false
)

func main (){
	if Conn, err := pgx.Connect(context.Background(),"host=localhost user=me password=12345 dbname=postgres"); err!= nil{
		log.Fatal("Невозможно подключиться к бд")
	}else{
		PostgresConn = Conn
		if _, err = PostgresConn.Exec(context.Background(), dropCreateFunctions); err!= nil{
			log.Fatal("Невозможно создать функции createdb dropdb в мастер бд")
		}
	}

	Router.GET("/", fasthttp.FSHandler("./frontend/index.html",1))
	Router.GET("/static/*filepath", fasthttp.FSHandler("./frontend/static", 1))

	Router.GET("/initdb", InitDb)
	Router.GET("/dropdb",DropDb)
	Router.GET("/inittables",CheckTrueDb(InitTables))
	Router.GET("/truncateusers",CheckTrueDb(TruncateUsers))
	Router.GET("/truncatecomments",CheckTrueDb(TruncateComments))
	Router.GET("/getallusers",CheckTrueDb(GetAllUsers))
	Router.GET("/getallcomments",CheckTrueDb(GetAllComments))
	Router.GET("/selectcommentbymessage",CheckTrueDb(SelectCommentsByMessage))
	Router.GET("/selectusersbynickname",CheckTrueDb(SelectUsersByNickname))
	Router.GET("/insertuser",CheckTrueDb(InsertUser))
	Router.GET("/insertcomment",CheckTrueDb(InsertComment))
	Router.GET("/updatecommenttext",CheckTrueDb(UpdateCommentText))
	Router.GET("/updateusernickname",CheckTrueDb(UpdateUserNickname))

	Router.GET("/deletecommentsbymessage",CheckTrueDb(DeleteCommentsByMessage))
	Router.GET("/deleteparticularcomment",CheckTrueDb(DeleteParticularComment))

	Router.GET("/deleteusersbynickname",CheckTrueDb(DeleteUsersByNickname))
	Router.GET("/deleteparticularuser",CheckTrueDb(DeleteParticularUser))


	fmt.Println("Listen and serve at port 9000")
	if err := fasthttp.ListenAndServe(":9000", Router.Handler); err != nil{
		log.Fatal("Невозможно запустить сервер")
	}
}