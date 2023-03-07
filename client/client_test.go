package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/opengoats/cmdb/apps/book"
	"github.com/opengoats/cmdb/common/pb/page"
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
		req := &book.CreateBookRequest{Name: "水浒", Author: "赵六"}
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
		req := &book.QueryBookRequest{Page: &page.PageRequest{PageSize: 10, PageNumber: 1}, Name: "%", Author: "%"}
		reply, err := client.QueryBook(ctx, req)
		if should.NoError(err) {
			fmt.Println(reply)
		}
	}
}
