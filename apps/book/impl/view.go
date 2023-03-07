package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/opengoats/cmdb/apps/book"
	"github.com/opengoats/cmdb/common/pb/base"
)

var (
	validate = validator.New()
)

func (s *service) CreateBook(ctx context.Context, req *book.CreateBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := validate.Struct(req); err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, status.Errorf(codes.InvalidArgument, "Request parameter error")
	}
	// book结构体赋值
	ins := &book.Book{
		Base: &base.Base{
			Id:       xid.New().String(),
			Status:   1,
			CreateAt: time.Now().UnixMicro(),
			CreateBy: "",
		},
		Data: req,
	}

	// 开启事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, status.Errorf(codes.Internal, "Creation failed. Please contact the administrator")
	}

	// 通过Defer处理事务提交方式
	// 1. 无报错，则Commit 事务
	// 2. 有报错, 则Rollback 事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Named("CreateBook").Error("rollback error ", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				s.log.Named("CreateBook").Error("commit error ", err)
			}
		}
	}()

	// 插入book表
	rstmt, err := tx.PrepareContext(ctx, insertBook)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, status.Errorf(codes.Internal, "Creation failed. Please contact the administrator")
	}
	defer rstmt.Close()

	_, err = rstmt.ExecContext(ctx,
		ins.Base.Id, ins.Base.Status, ins.Base.CreateAt, ins.Base.CreateBy, ins.Data.Name, ins.Data.Author,
	)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, status.Errorf(codes.Internal, "Creation failed. Please contact the administrator")
	}

	return ins, nil
}

func (s *service) QueryBook(ctx context.Context, req *book.QueryBookRequest) (*book.BookSet, error) {

	// 请求体校验
	if err := validate.Struct(req); err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, status.Errorf(codes.InvalidArgument, "Request parameter error")
	}

	// 参数
	offSet := int64((req.Page.PageNumber - 1) * req.Page.PageSize)
	args := []interface{}{req.Name, req.Author, offSet, req.Page.PageSize}

	s.log.Named("QueryBook").Infof("query sql: %s; %v", queryBook, args)

	// query stmt, 构建一个Prepare语句
	stmt, err := s.db.PrepareContext(ctx, queryBook)
	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, status.Errorf(codes.Internal, "Query failed. Please contact the administrator")
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, status.Errorf(codes.Internal, "Query failed. Please contact the administrator")
	}
	defer rows.Close()

	// 结构体赋值
	set := &book.BookSet{Items: []*book.Book{}}
	for rows.Next() {
		// 每扫描一行,就需要读取出来
		ins := &book.Book{Base: &base.Base{}, Data: &book.CreateBookRequest{}}
		if err := rows.Scan(&ins.Base.Id, &ins.Base.Status, &ins.Base.CreateAt, &ins.Base.CreateBy,
			&ins.Base.UpdateAt, &ins.Base.UpdateBy, &ins.Base.DeleteAt, &ins.Base.DeleteBy,
			&ins.Data.Name, &ins.Data.Author); err != nil {
			s.log.Named("QueryBook").Error(err)
			return nil, status.Errorf(codes.Internal, "Query failed. Please contact the administrator")
		}
		set.Items = append(set.Items, ins)
	}
	fmt.Printf("set: %+v\n", set)
	// total统计
	set.Total = int64(len(set.Items))

	return set, nil
}

// func (s *service) DescribeBook(ctx context.Context, req *book.DescribeBookRequest) (
// 	*book.Book, error) {
// 	query := sqlbuilder.NewQuery(queryBookSQL)
// 	querySQL, args := query.Where("id = ?", req.Id).BuildQuery()
// 	s.log.Debugf("sql: %s", querySQL)

// 	queryStmt, err := s.db.Prepare(querySQL)
// 	if err != nil {
// 		return nil, exception.NewInternalServerError("prepare query book error, %s", err.Error())
// 	}
// 	defer queryStmt.Close()

// 	ins := book.NewDefaultBook()
// 	err = queryStmt.QueryRow(args...).Scan(
// 		&ins.Id, &ins.CreateAt, &ins.Data.CreateBy, &ins.UpdateAt, &ins.UpdateBy,
// 		&ins.Data.Name, &ins.Data.Author,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, exception.NewNotFound("%s not found", req.Id)
// 		}
// 		return nil, exception.NewInternalServerError("describe book error, %s", err.Error())
// 	}

// 	return ins, nil
// }

// func (s *service) UpdateBook(ctx context.Context, req *book.UpdateBookRequest) (
// 	*book.Book, error) {
// 	ins, err := s.DescribeBook(ctx, book.NewDescribeBookRequest(req.Id))
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch req.UpdateMode {
// 	case request.UpdateMode_PUT:
// 		ins.Update(req)
// 	case request.UpdateMode_PATCH:
// 		err := ins.Patch(req)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// 校验更新后数据合法性
// 	if err := ins.Data.Validate(); err != nil {
// 		return nil, err
// 	}

// 	if err := s.updateBook(ctx, ins); err != nil {
// 		return nil, err
// 	}

// 	return ins, nil
// }

// func (s *service) DeleteBook(ctx context.Context, req *book.DeleteBookRequest) (
// 	*book.Book, error) {
// 	ins, err := s.DescribeBook(ctx, book.NewDescribeBookRequest(req.Id))
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := s.deleteBook(ctx, ins); err != nil {
// 		return nil, err
// 	}

// 	return ins, nil
// }
