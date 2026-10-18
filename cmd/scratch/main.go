package main
import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/martiniscoding/quexy/internal/store"
)

func main(){
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	s,err:=store.New(ctx)
	if err!=nil{
		log.Fatal(err)
	}
	//s data type  se jo bhi functionn store folder me implemented hai wo use kar sakte hai ham s pr 
	//Close, *Store ka method hai (receiver), isliye s.Close() likh sakte hain agar parameter hota to Close(s) 
	//likhna padta 
	defer s.Close()

	fmt.Println("connection to db ", s)
} 