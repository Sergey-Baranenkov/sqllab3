package main

import (
	"bytes"
	"context"
	json2 "encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"log"
)

var null = json2.RawMessage("null")
var emptyArray = json2.RawMessage("[]")
func UpdateCommentText (ctx *fasthttp.RequestCtx){
	commentId:= string(ctx.QueryArgs().Peek("commentId"))
	messageText:= string(ctx.QueryArgs().Peek("messageText"))
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select update_comment_text($1,$2)", commentId, messageText).Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func SelectUsersByNickname (ctx *fasthttp.RequestCtx){
	nickname:= string(ctx.QueryArgs().Peek("nickname"))
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select select_users_by_nickname($1)", nickname).Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func DeleteUsersByNickname (ctx *fasthttp.RequestCtx){
	nickname:= string(ctx.QueryArgs().Peek("nickname"))
	if _, err := PostgresConn.Exec(context.Background(), "select delete_users_by_nickname($1)", nickname); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}

func SelectCommentsByMessage (ctx *fasthttp.RequestCtx){
	messageText:= string(ctx.QueryArgs().Peek("messageText"))
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select select_comments_by_message($1)", messageText).Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}


func DeleteCommentsByMessage (ctx *fasthttp.RequestCtx){
	messageText:= string(ctx.QueryArgs().Peek("messageText"))
	if _, err := PostgresConn.Exec(context.Background(), "select delete_comments_by_message($1)", messageText); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}


func DeleteParticularUser (ctx *fasthttp.RequestCtx){
	userId:= string(ctx.QueryArgs().Peek("userId"))
	if _, err := PostgresConn.Exec(context.Background(), "select delete_particular_user($1)", userId); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}


func DeleteParticularComment (ctx *fasthttp.RequestCtx){
	commentId:= string(ctx.QueryArgs().Peek("commentId"))
	if _, err := PostgresConn.Exec(context.Background(), "select delete_particular_comment($1)", commentId); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}

func UpdateUserNickname (ctx *fasthttp.RequestCtx){
	userId:= string(ctx.QueryArgs().Peek("userId"))
	newNickname := string(ctx.QueryArgs().Peek("newNickname"))
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select update_user_nickname($1,$2)", userId, newNickname).Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func InsertComment (ctx *fasthttp.RequestCtx){
	authId:= string(ctx.QueryArgs().Peek("authId"))
	message := string(ctx.QueryArgs().Peek("messageText"))
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select insert_comment($1,$2)", authId, message).Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func InsertUser (ctx *fasthttp.RequestCtx){
	nickname:= string(ctx.QueryArgs().Peek("nickname"))
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select insert_user($1)", nickname).Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func TruncateComments (ctx *fasthttp.RequestCtx){
	if _, err := PostgresConn.Exec(context.Background(), "select truncate_comments()"); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}

func TruncateUsers (ctx *fasthttp.RequestCtx){
	if _, err := PostgresConn.Exec(context.Background(), "select truncate_users()"); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}

func GetAllComments (ctx *fasthttp.RequestCtx){
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select get_all_comments()").Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func GetAllUsers (ctx *fasthttp.RequestCtx){
	var json json2.RawMessage
	if err := PostgresConn.QueryRow(context.Background(), "select get_all_users()").Scan(&json); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	if bytes.Equal(json, null){
		json = emptyArray
	}
	_, _ = ctx.WriteString(string(json))
}

func InitTables(ctx *fasthttp.RequestCtx){
	if _, err := PostgresConn.Exec(context.Background(), "select create_tables()"); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}

func InitDb(ctx *fasthttp.RequestCtx){
	if _, err := PostgresConn.Exec(context.Background(), "select create_database()"); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	Conn, err := pgx.Connect(context.Background(),"host=localhost user=me password=12345 dbname=lab4_db")

	if err!= nil {
		log.Fatal("Невозможно подключиться к бд")
	}

	err = PostgresConn.Close(context.Background())
	if err != nil {
		log.Fatal("Невозможно закрыть старое подключение")
	}
	PostgresConn = Conn
	onTrueDatabase = true

	if _, err = PostgresConn.Exec(context.Background(), dropCreateFunctions); err!= nil{
		log.Fatal("Невозможно создать функции createdb dropdb в lab4_db")
	}

	if _, err := PostgresConn.Exec(context.Background(), sqlFunctions); err != nil{
		ctx.Error(err.Error(),400)
		log.Fatal("Невозможно создать функции")
	}

	ctx.SetStatusCode(200)
}

func DropDb(ctx *fasthttp.RequestCtx){
	if Conn, err := pgx.Connect(context.Background(),"host=localhost user=me password=12345 dbname=postgres"); err!= nil{
		ctx.Error(err.Error(),400)
		log.Fatal("Невозможно подключиться к бд")
	}else{
		_ = PostgresConn.Close(context.Background())
		PostgresConn = Conn
		onTrueDatabase = false
	}

	if _, err := PostgresConn.Exec(context.Background(), "select drop_database()"); err != nil{
		ctx.Error(err.Error(),400)
		return
	}
	ctx.SetStatusCode(200)
}

func CheckTrueDb(next fasthttp.RequestHandler) fasthttp.RequestHandler{
	return func(ctx *fasthttp.RequestCtx) {
		if onTrueDatabase {
			next(ctx)
		}else{
			ctx.Error("бд не создана", 400)
		}
	}
}