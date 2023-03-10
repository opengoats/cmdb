package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/opengoats/cmdb/apps/book"
	"github.com/opengoats/goat/app"
	"github.com/opengoats/goat/http/response"
	"github.com/opengoats/goat/logger"
	"github.com/opengoats/goat/logger/zap"
)

var (
	h = &handler{}
)

type handler struct {
	service book.ServiceServer
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(book.AppName)
	h.service = app.GetGrpcApp(book.AppName).(book.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return book.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{"books"}

	ws.Route(ws.POST("").To(h.CreateBook).
		Doc("create a book").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(book.CreateBookRequest{}).
		Writes(response.NewMessage(book.Book{})))

	ws.Route(ws.GET("/").To(h.QueryBook).
		Doc("get all books").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata("action", "list").
		Reads(book.CreateBookRequest{}).
		Writes(response.NewMessage(book.BookSet{})).
		Returns(200, "OK", book.BookSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeBook).
		Doc("get a book").
		Param(ws.PathParameter("id", "identifier of the book").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(response.NewMessage(book.Book{})).
		Returns(200, "OK", response.NewMessage(book.Book{})).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{id}").To(h.PutBook).
		Doc("update a book").
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(book.CreateBookRequest{}))

	ws.Route(ws.PATCH("/{id}").To(h.PatchBook).
		Doc("patch a book").
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(book.CreateBookRequest{}))

	ws.Route(ws.DELETE("/{id}").To(h.DeleteBook).
		Doc("delete a book").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")))
}

func init() {
	app.RegistryRESTfulApp(h)
}
