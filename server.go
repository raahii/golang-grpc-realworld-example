package main

import (
	"fmt"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/raahii/golang-grpc-realworld-example/db"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	w := zerolog.ConsoleWriter{Out: os.Stderr}
	l := zerolog.New(w).With().Timestamp().Caller().Logger()

	d, err := db.New()
	if err != nil {
		err = fmt.Errorf("failed to connect database: %w", err)
		l.Fatal().Err(err).Msg("failed to connect the database")
	}
	l.Info().Str("name", d.Dialect().GetName()).
		Str("database", d.Dialect().CurrentDatabase()).
		Msg("succeeded to connect to the database")

	err = db.AutoMigrate(d)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to migrate database")
	}

	us := store.NewUserStore(d)
	as := store.NewArticleStore(d)

	h := handler.New(&l, us, as)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		l.Panic().Err(fmt.Errorf("failed to listen: %w", err))
	}

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, h)
	pb.RegisterArticlesServer(s, h)
	l.Info().Str("port", port).Msg("starting server")
	if err := s.Serve(lis); err != nil {
		l.Panic().Err(fmt.Errorf("failed to serve: %w", err))
	}
}
