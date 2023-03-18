package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/opengoats/cmdb/apps/book"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// func TestBookQuery(t *testing.T) {
// 	should := assert.New(t)

// 	conf := client.NewDefaultConfig()
// 	// 设置GRPC服务地址
// 	conf.SetAddress("127.0.0.1:8050")
// 	// 携带认证信息
// 	// conf.SetClientCredentials("secret_id", "secret_key")
// 	c, err := client.NewClient(conf)
// 	if should.NoError(err) {
// 		resp, err := c.Book().QueryBook(
// 			context.Background(),
// 			// book.NewQueryBookRequest(),
// 		)
// 		should.NoError(err)
// 		fmt.Println(resp.Items)
// 	}
// }

func TestCreateBook(t *testing.T) {
	should := assert.New(t)
	conn, err := grpc.Dial("localhost:18050", grpc.WithInsecure())
	defer conn.Close()
	if should.NoError(err) {
		client := book.NewServiceClient(conn)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
		req := &book.CreateBookRequest{BookName: "水浒呼呼1111111111", Author: "张三"}
		reply, err := client.CreateBook(ctx, req)
		if should.NoError(err) {
			fmt.Println(reply)
		}
	}
}

func TestQueryBook(t *testing.T) {
	should := assert.New(t)
	conn, err := grpc.Dial("localhost:18050", grpc.WithInsecure())
	defer conn.Close()
	if should.NoError(err) {
		client := book.NewServiceClient(conn)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
		req := book.NewQueryBookRequest()
		reply, err := client.QueryBook(ctx, req)
		if should.NoError(err) {
			fmt.Println(reply)
		}
	}
}

func TestDescribeBook(t *testing.T) {
	should := assert.New(t)
	conn, err := grpc.Dial("localhost:18050", grpc.WithInsecure())
	defer conn.Close()
	if should.NoError(err) {
		client := book.NewServiceClient(conn)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
		req := &book.DescribeBookRequest{Id: "cg4covivlkjhml7adi30"}
		reply, err := client.DescribeBook(ctx, req)
		if should.NoError(err) {
			fmt.Println(reply)
		}
	}
}

func TestUpdateBook(t *testing.T) {
	should := assert.New(t)
	conn, err := grpc.Dial("localhost:18050", grpc.WithInsecure())
	defer conn.Close()
	if should.NoError(err) {
		client := book.NewServiceClient(conn)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
		req := &book.UpdateBookRequest{Id: "cg4covivlkjhml7adi30", Data: &book.CreateBookRequest{BookName: "小瘪三"}, UpdateMode: 1}
		reply, err := client.UpdateBook(ctx, req)
		if should.NoError(err) {
			fmt.Println(reply)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	should := assert.New(t)
	conn, err := grpc.Dial("localhost:18050", grpc.WithInsecure())
	defer conn.Close()
	if should.NoError(err) {
		client := book.NewServiceClient(conn)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
		req := &book.DeleteBookRequest{Id: "cg4crqavlkjhml7adi40"}
		reply, err := client.DeleteBook(ctx, req)
		if should.NoError(err) {
			fmt.Println(reply)
		}
	}
}
