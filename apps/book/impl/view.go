package impl

import (
	"context"

	"github.com/go-playground/validator/v10"

	"github.com/opengoats/cmdb/apps/book"
	"github.com/opengoats/goat/pb/request"
)

var (
	validate = validator.New()
)

func (s *service) CreateBook(ctx context.Context, req *book.CreateBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, err
	}
	// book结构体赋值
	ins := book.NewBook()
	ins.Data = req

	// 开启事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, err
	}

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
	s.log.Named("CreateBook").Debugf("sql: %s", insertBook)
	rstmt, err := tx.PrepareContext(ctx, insertBook)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, err
	}
	defer rstmt.Close()

	_, err = rstmt.ExecContext(ctx, ins.Id, ins.Status, ins.CreateAt, ins.CreateBy, ins.Data.BookName, ins.Data.Author)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, err
	}

	return ins, nil
}

func (s *service) QueryBook(ctx context.Context, req *book.QueryBookRequest) (*book.BookSet, error) {

	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, err
	}

	// 数据库插入参数
	args := []interface{}{req.BookName, req.Author, req.Page.ComputeOffset(), uint(req.Page.PageSize)}

	// query stmt, 构建一个Prepare语句
	s.log.Named("QueryBook").Debugf("sql: %s; %v", queryBook, args)
	stmt, err := s.db.PrepareContext(ctx, queryBook)
	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, err
	}
	defer rows.Close()

	// 结构体赋值
	set := book.NewBookSet()
	for rows.Next() {
		ins := book.NewDefaultBook()
		err := rows.Scan(&ins.Id, &ins.Status, &ins.CreateAt, &ins.CreateBy, &ins.UpdateAt, &ins.UpdateBy,
			&ins.DeleteAt, &ins.DeleteBy, &ins.Data.BookName, &ins.Data.Author)

		if err != nil {
			s.log.Named("QueryBook").Error(err)
			return nil, err
		}
		set.Add(ins)
	}

	// total统计
	set.Total = int64(len(set.Items))

	return set, nil
}

func (s *service) DescribeBook(ctx context.Context, req *book.DescribeBookRequest) (*book.Book, error) {

	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("DescribeBook").Error(err)
		return nil, err
	}

	args := []interface{}{req.Id}

	// query stmt, 构建一个Prepare语句
	s.log.Named("DescribeBook").Debugf("sql: %s; %v", describeBook, args)
	stmt, err := s.db.PrepareContext(ctx, describeBook)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// 取出数据，赋值结构体
	ins := book.NewDefaultBook()

	err = stmt.QueryRowContext(ctx, args...).Scan(&ins.Id, &ins.Status, &ins.CreateAt, &ins.CreateBy, &ins.UpdateAt, &ins.UpdateBy,
		&ins.DeleteAt, &ins.DeleteBy, &ins.Data.BookName, &ins.Data.Author)

	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, err
	}
	return ins, nil
}

func (s *service) UpdateBook(ctx context.Context, req *book.UpdateBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("UpdateBook").Error(err)
		return nil, err
	}

	// 验证更新id,查询不到直接返回
	ins, err := s.DescribeBook(ctx, &book.DescribeBookRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	// 根据更新模式进行数据库操作
	switch req.UpdateMode {
	case request.UpdateMode_PUT:
		ins.Update(req)
	case request.UpdateMode_PATCH:
		if err := ins.Patch(req); err != nil {
			s.log.Named("UpdateBook").Error(err)
			return nil, err
		}
	}

	// 校验更新后数据合法性
	if err := ins.Validate(); err != nil {
		s.log.Named("UpdateBook").Error(err)
		return nil, err
	}

	// 更新数据库
	// 开启一个事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 通过Defer处理事务提交方式
	// 1. 无报错，则Commit 事务
	// 2. 有报错，则Rollback 事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Error("rollback error, %s", err.Error())
			}
		} else {
			if err := tx.Commit(); err != nil {
				s.log.Error("commit error, %s", err.Error())
			}
		}
	}()

	s.log.Named("UpdateBook").Debugf("sql: %s", updateBook)
	bookStmt, err := tx.PrepareContext(ctx, updateBook)
	_, err = bookStmt.ExecContext(ctx, ins.UpdateAt, ins.UpdateBy, ins.Data.BookName, ins.Data.Author, ins.Id)
	if err != nil {
		return nil, err
	}
	defer bookStmt.Close()

	return ins, nil
}

func (s *service) DeleteBook(ctx context.Context, req *book.DeleteBookRequest) (*book.Book, error) {
	// 请求体校验
	if err := req.Validate(); err != nil {
		s.log.Named("DeleteBook").Error(err)
		return nil, err
	}

	// 验证更新id,查询不到直接返回
	ins, err := s.DescribeBook(ctx, &book.DescribeBookRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	// 开启一个事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 通过Defer处理事务提交方式
	// 1. 无报错，则Commit 事务
	// 2. 有报错，则Rollback 事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Error("rollback error, %s", err.Error())
			}
		} else {
			if err := tx.Commit(); err != nil {
				s.log.Error("commit error, %s", err.Error())
			}
		}
	}()

	s.log.Named("DeleteBook").Debugf("sql: %s", deleteBook)
	bookStmt, err := tx.PrepareContext(ctx, deleteBook)
	if err != nil {
		return nil, err
	}
	defer bookStmt.Close()

	_, err = bookStmt.ExecContext(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return ins, nil
}
