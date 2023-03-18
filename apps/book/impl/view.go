package impl

import (
	"context"

	"github.com/opengoats/cmdb/apps/book"
	"github.com/opengoats/goat/exception"
)

func (s *service) CreateBook(ctx context.Context, req *book.CreateBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, exception.NewBadRequest("validate create book error, %s", err)
	}
	// book结构体赋值
	ins := book.NewBook()
	ins.Data = req
	// 存入数据库
	return s.save(ctx, ins)
}

func (s *service) QueryBook(ctx context.Context, req *book.QueryBookRequest) (*book.BookSet, error) {

	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, exception.NewBadRequest("validate query book error, %s", err)
	}

	// 查询默认值填充
	if req.Keywords == "" {
		req.Keywords = "%"
	} else {
		req.Keywords = req.Keywords + "%"
	}

	// 数据库查询
	return s.query(ctx, req)
}

func (s *service) DescribeBook(ctx context.Context, req *book.DescribeBookRequest) (*book.Book, error) {

	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("DescribeBook").Error(err)
		return nil, exception.NewBadRequest("validate describe book error, %s", err)
	}

	// 数据库查询
	return s.describe(ctx, req)
}

func (s *service) UpdateBook(ctx context.Context, req *book.UpdateBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("UpdateBook").Error(err)
		return nil, exception.NewBadRequest("validate update book error, %s", err)
	}

	// 验证更新id,查询不到直接返回
	ins, err := s.DescribeBook(ctx, &book.DescribeBookRequest{Id: req.Id})
	if err != nil {
		return nil, exception.NewBadRequest("id not exist, %s", err)
	}
	// 数据库更新
	return s.update(ctx, req, ins)
}

func (s *service) DeleteBook(ctx context.Context, req *book.DeleteBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("DeleteBook").Error(err)
		return nil, exception.NewBadRequest("validate delete book error, %s", err)
	}

	// 验证更新id,查询不到直接返回
	ins, err := s.DescribeBook(ctx, &book.DescribeBookRequest{Id: req.Id})
	if err != nil {
		return nil, exception.NewBadRequest("id not exist, %s", err)
	}
	// 数据库操作
	return s.delete(ctx, req, ins)
}
