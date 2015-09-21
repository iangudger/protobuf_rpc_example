package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"

	"github.com/iangudger/protobuf_rpc_example/proto/message"
)

func apiHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Println("Got api request from:", req.RemoteAddr)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	reqMsg := message.Message{}
	if err := proto.Unmarshal(buf.Bytes(), &reqMsg); err != nil {
		log.Println("Error unmarshaling:", err)
		writeRes(res, err.Error())
		return
	}

	if reqMsg.GetText() == "hello" {
		log.Println("got hello")
		writeRes(res, "hello back")
	} else {
		log.Println("not hello")
		writeRes(res, "not hello")
	}
}

func staticHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Println("Got static file request from:", req.RemoteAddr)
	res.Header().Set("cache-control", "public, no-transform")
	path := "./static" + ps.ByName("filepath")
	log.Println("Serving file:", path)
	http.ServeFile(res, req, path)
}

func writeRes(res http.ResponseWriter, msgtxt string) {
	msg := message.Message{Text: proto.String(msgtxt)}
	bin, err := proto.Marshal(&msg)
	if err != nil {
		log.Println("Error marshaling:", err)
		return
	}
	res.Write(bin)
}

func main() {
	router := httprouter.New()

	router.POST("/api", apiHandler)

	router.GET(
		"/*filepath",
		staticHandler,
	)

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", router))
}
