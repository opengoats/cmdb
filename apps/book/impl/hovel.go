package impl

import (
	"context"

	"github.com/opengoats/cmdb/apps/book"
	"github.com/opengoats/goat/exception"
	"github.com/opengoats/goat/pb/request"
)

func (s *service) save(ctx context.Context, ins *book.Book) (*book.Book, error) {
	// 开启事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, exception.NewInternalServerError("start tx err %s", err)
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
		return nil, exception.NewInternalServerError("insert table book err %s", err)
	}
	defer rstmt.Close()

	_, err = rstmt.ExecContext(ctx, ins.Id, ins.Status, ins.CreateAt, ins.CreateBy, ins.Data.BookName, ins.Data.Author)
	if err != nil {
		s.log.Named("CreateBook").Error(err)
		return nil, exception.NewInternalServerError("insert table book err %s", err)
	}

	return ins, nil
}

func (s *service) query(ctx context.Context, req *book.QueryBookRequest) (*book.BookSet, error) {
	// 数据库插入参数
	args := []interface{}{req.Keywords, req.Keywords, req.Page.ComputeOffset(), uint(req.Page.PageSize)}

	// query stmt, 构建一个Prepare语句
	s.log.Named("QueryBook").Debugf("sql: %s; %v", queryBook, args)
	stmt, err := s.db.PrepareContext(ctx, queryBook)
	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, exception.NewInternalServerError("query table book err %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, exception.NewInternalServerError("query table book err %s", err)
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
			return nil, exception.NewInternalServerError("query table book err %s", err)
		}
		set.Add(ins)
	}

	// total统计
	set.Total = int64(len(set.Items))

	return set, nil

}

func (s *service) describe(ctx context.Context, req *book.DescribeBookRequest) (*book.Book, error) {
	args := []interface{}{req.Id}

	// query stmt, 构建一个Prepare语句
	s.log.Named("DescribeBook").Debugf("sql: %s; %v", describeBook, args)
	stmt, err := s.db.PrepareContext(ctx, describeBook)
	if err != nil {
		return nil, exception.NewInternalServerError("describe book err %s", err)
	}
	defer stmt.Close()

	// 取出数据，赋值结构体
	ins := book.NewDefaultBook()

	err = stmt.QueryRowContext(ctx, args...).Scan(&ins.Id, &ins.Status, &ins.CreateAt, &ins.CreateBy, &ins.UpdateAt, &ins.UpdateBy,
		&ins.DeleteAt, &ins.DeleteBy, &ins.Data.BookName, &ins.Data.Author)

	if err != nil {
		s.log.Named("QueryBook").Error(err)
		return nil, exception.NewInternalServerError("describe book err %s", err)
	}
	return ins, nil
}

func (s *service) update(ctx context.Context, req *book.UpdateBookRequest, ins *book.Book) (*book.Book, error) {
	// 根据更新模式进行数据库操作
	switch req.UpdateMode {
	case request.UpdateMode_PUT:
		ins.Update(req)
	case request.UpdateMode_PATCH:
		if err := ins.Patch(req); err != nil {
			s.log.Named("UpdateBook").Error(err)
			return nil, exception.NewInternalServerError("update book err %s", err)
		}
	}

	// 校验更新后数据合法性
	if err := ins.Validate(); err != nil {
		s.log.Named("UpdateBook").Error(err)
		return nil, exception.NewInternalServerError("update book err %s", err)
	}

	// 更新数据库
	// 开启一个事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, exception.NewInternalServerError("update book err %s", err)
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
		return nil, exception.NewInternalServerError("update book err %s", err)
	}
	defer bookStmt.Close()

	return ins, nil
}

func (s *service) delete(ctx context.Context, req *book.DeleteBookRequest, ins *book.Book) (*book.Book, error) {
	// 开启一个事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, exception.NewInternalServerError("delete book err %s", err)
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
		return nil, exception.NewInternalServerError("delete book err %s", err)
	}
	defer bookStmt.Close()

	_, err = bookStmt.ExecContext(ctx, req.Id)
	if err != nil {
		return nil, exception.NewInternalServerError("delete book err %s", err)
	}

	return ins, nil
}
