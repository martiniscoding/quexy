package store
import (
	"context"
	"fmt"
	"time"
	"os"
	//pgx pool is a driver to connect postgress to go lang 
	"github.com/jackc/pgx/v5/pgxpool"
)
type Store struct {
	pool *pgxpool.Pool
}
func New(ctx context.Context)(*Store,error){
	
	//short function pgxpool.New  would have connected tod db 
	// but it will not allow modificatons 
	//bhai docker me pg chalayag to url ye banti hai 
	//postgres://USER:PASSWORD@HOST:PORT/DATABASE?sslmode=disable
	url :=os.Getenv("DATABASE_URL")
	//parse xonfig  url ko todega 
	cfg,err:=pgxpool.ParseConfig(url)
	if err!=nil{
		fmt.Println("Error parsing config:", err)
		return nil,fmt.Errorf("Error parsing config: %v", err)
	}
	cfg.MaxConns = 10;
	cfg.MinConns = 2;
	cfg.MaxConnLifetime = 30 * time.Minute;
	pool,err := pgxpool.NewWithConfig(ctx, cfg)
	if err!=nil{
		fmt.Println("Error creating connection pool:", err)
		return nil,fmt.Errorf("Error creating connection pool: %v", err)
	}
	s:=Store{pool}
	fmt.Println("connecting to db ", s)
	return &s,nil
	
}

func (s *Store) Close() {
	s.pool.Close()
}