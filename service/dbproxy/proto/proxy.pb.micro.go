// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proxy.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for DBProxyService service

func NewDBProxyServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for DBProxyService service

type DBProxyService interface {
	// 请求执行sql语句
	ExecuteAction(ctx context.Context, in *ReqExec, opts ...client.CallOption) (*RespExec, error)
}

type dBProxyService struct {
	c    client.Client
	name string
}

func NewDBProxyService(name string, c client.Client) DBProxyService {
	return &dBProxyService{
		c:    c,
		name: name,
	}
}

func (c *dBProxyService) ExecuteAction(ctx context.Context, in *ReqExec, opts ...client.CallOption) (*RespExec, error) {
	req := c.c.NewRequest(c.name, "DBProxyService.ExecuteAction", in)
	out := new(RespExec)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DBProxyService service

type DBProxyServiceHandler interface {
	// 请求执行sql语句
	ExecuteAction(context.Context, *ReqExec, *RespExec) error
}

func RegisterDBProxyServiceHandler(s server.Server, hdlr DBProxyServiceHandler, opts ...server.HandlerOption) error {
	type dBProxyService interface {
		ExecuteAction(ctx context.Context, in *ReqExec, out *RespExec) error
	}
	type DBProxyService struct {
		dBProxyService
	}
	h := &dBProxyServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&DBProxyService{h}, opts...))
}

type dBProxyServiceHandler struct {
	DBProxyServiceHandler
}

func (h *dBProxyServiceHandler) ExecuteAction(ctx context.Context, in *ReqExec, out *RespExec) error {
	return h.DBProxyServiceHandler.ExecuteAction(ctx, in, out)
}
