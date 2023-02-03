package service

import (
	"context"
	"crypto/tls"
	"ehang.io/nps/api/protos/account"
	"github.com/beego/beego"
	"github.com/beego/beego/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

var (
	conn *grpc.ClientConn
	once sync.Once
)

func getConn() *grpc.ClientConn {
	once.Do(func() {
		var credential credentials.TransportCredentials
		if isTls, err := beego.AppConfig.Bool("GRPC_CLIENT_ACCOUNT_TLS"); err == nil && isTls {
			credential = credentials.NewTLS(&tls.Config{})
		} else {
			credential = insecure.NewCredentials()
		}

		client, err := grpc.Dial(beego.AppConfig.String("GRPC_CLIENT_ACCOUNT_URL"), grpc.WithTransportCredentials(credential))
		if err != nil {
			logs.Error("grpc connect account: %v", err)
		}
		conn = client
	})
	return conn
}

func GetSignInUrl() string {
	c := account.NewOauthClient(getConn())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := c.GetSignInUrl(ctx, &account.GetSignInUrlRequest{
		ApplicationSecret: beego.AppConfig.String("GRPC_CLIENT_ACCOUNT_SECRET"),
	})
	if err != nil {
		log.Printf("grpc get error GetSignInUrl: %v", err)
	}

	return r.GetSignInUrl()
}

func GetUserInfoByToken(token string) (*account.GetUserInfoByTokenResponse, error) {
	c := account.NewInfoClient(getConn())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := c.GetUserInfoByToken(ctx, &account.GetUserInfoByTokenRequest{
		ApplicationSecret: beego.AppConfig.String("GRPC_CLIENT_ACCOUNT_SECRET"),
		Token:             token,
	})
	if err != nil {
		log.Printf("grpc get error GetUserInfoByEncodeUserId: %v", err)
	}

	return r, err
}
