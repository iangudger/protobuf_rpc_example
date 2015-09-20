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
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	reqMsg := message.Message{}
	if err := proto.Unmarshal(buf.Bytes(), &reqMsg); err != nil {
		log.Println("unmarshaling:", err)
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
	res.Header().Set("cache-control", "public, no-transform")
	path := "./static" + ps.ByName("filepath")
	log.Println("Serving file:", path)
	http.ServeFile(res, req, path)
}

func writeRes(res http.ResponseWriter, msgtxt string) {
	msg := message.Message{Text: proto.String(msgtxt)}
	bin, err := proto.Marshal(&msg)
	if err != nil {
		log.Println("marshaling:", err)
		return
	}
	log.Printf("Sending %d bytes\n", len(bin))
	log.Println("Sending message: ", bin)
	log.Println("Sending message: ", proto.MarshalTextString(&msg))
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