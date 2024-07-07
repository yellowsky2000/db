package main

import (
	"database/sql"

	"github.com/yellowsky2000/db/handler"
	pb "github.com/yellowsky2000/db/proto"

	_ "github.com/jackc/pgx/v4/stdlib"
	admin "github.com/yellowsky2000/pkg/service/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config/client"
	"micro.dev/v4/service/logger"
)

var dbAddress = "postgresql://postgres:posgres@localhost:5432/db?sslmode=disable"

func main() {
	// create a service
	srv := service.New(
		service.Name("db"),
		service.Version("latest"),
	)

	// connect to the database
	cfg, err := client.NewConfig().Get("db.address")
	if err != nil {
		logger.Fatalf("error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}

	h := &handler.Db{}
	h.DBConn(sqlDB)

	// Register handler
	pb.RegisterDbHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
