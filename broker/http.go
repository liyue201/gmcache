package broker

import (
	"encoding/json"
	"github.com/codinl/binding"
	"github.com/codinl/go-logger"
	"github.com/liyue201/gmcache/broker/rpc"
	"github.com/liyue201/gmcache/proto"
	"github.com/liyue201/martini"
	"net/http"
)

const (
	StatusOK                  = 200
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

var StatusText = map[int]string{
	StatusOK:                  "OK",
	StatusNoContent:           "No Content",
	StatusBadRequest:          "Bad request",
	StatusInternalServerError: "Internal server error",
}

type SetReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	Ttl   uint64 `json:"ttl" binding:"required"`
}

type GetReq struct {
	Key string `json:"key" binding:"required"`
}

type DelReq struct {
	Key string `json:"key" binding:"required"`
}

func InitRouter(m *martini.ClassicMartini) {
	m.Post("(v1/set)", binding.Bind(SetReq{}), doSet)
	m.Post("(v1/get)", binding.Bind(GetReq{}), doGet)
	m.Post("(v1/del)", binding.Bind(DelReq{}), doDelete)
}

func doSet(resp http.ResponseWriter, req *http.Request, r SetReq) {
	if r.Key == "" || r.Value == "" || r.Ttl <= 0 {
		RespErr(resp, StatusBadRequest)
		return
	}
	conn, err := rpc.GetClientConn()
	if err != nil {
		logger.Errorf("doSet: GetClientConn error")
		RespErr(resp, StatusInternalServerError)
		return
	}
	arg := &proto.SetOptArg{Key: r.Key, Val: []byte(r.Value), Ttl: r.Ttl}
	client := proto.NewRpcServiceClient(conn)

	_, err = client.Set(rpc.NewRpcContext(arg.Key), arg)
	if err != nil {
		logger.Errorf("doSet:", err)
		return
	}
	RespJson(resp, StatusOK, nil)
}

func doGet(resp http.ResponseWriter, req *http.Request, r GetReq) {
	if r.Key == "" {
		RespErr(resp, StatusBadRequest)
		return
	}
	var err error
	conn, err := rpc.GetClientConn()
	if err != nil {
		logger.Errorf("doGet: GetClientConn error")
		RespErr(resp, StatusInternalServerError)
		return
	}
	arg := &proto.GetOptArg{Key: r.Key}
	client := proto.NewRpcServiceClient(conn)

	ret, err := client.Get(rpc.NewRpcContext(arg.Key), arg)
	if err != nil {
		logger.Errorf("doGet:", err)
		return
	}
	if ret.Code == proto.RCODE_FAILURE {
		RespJson(resp, StatusNoContent, nil)
		return
	}
	RespJson(resp, StatusOK, string(ret.Val))
}

func doDelete(resp http.ResponseWriter, req *http.Request, r DelReq) {
	if r.Key == "" {
		RespErr(resp, StatusBadRequest)
		return
	}
	var err error
	conn, err := rpc.GetClientConn()
	if err != nil {
		logger.Errorf("doDel: GetClientConn error")
		RespErr(resp, StatusInternalServerError)
		return
	}
	arg := &proto.DelOptArg{Key: r.Key}
	client := proto.NewRpcServiceClient(conn)

	ret, err := client.Delete(rpc.NewRpcContext(arg.Key), arg)
	if err != nil {
		logger.Errorf("doDel:", err)
		return
	}
	if ret.Code == proto.RCODE_FAILURE {
		RespJson(resp, StatusNoContent, nil)
		return
	}
	RespJson(resp, StatusOK, nil)
}

func RespErr(resp http.ResponseWriter, code int) {
	http.Error(resp, StatusText[code], code)
}

func RespJson(resp http.ResponseWriter, code int, data interface{}) {
	errByte := []byte(`{"code":500, "desc":"Internal server error", "data":null}`)

	result := struct {
		Code int         `json:"code"`
		Desc string      `json:"desc"`
		Data interface{} `json:"data"`
	}{
		Code: code,
		Desc: StatusText[code],
		Data: data,
	}

	b, err := json.Marshal(result)
	if err != nil {
		resp.Write(errByte)
		logger.Error(err)
		return
	}

	logger.Debugf("resp: %s", string(b))
	resp.Write(b)
}
