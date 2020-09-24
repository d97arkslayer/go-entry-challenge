package Database

import (
	"github.com/d97arkslayer/go-entry-challenge/Types"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
	"log"
	"os"
)



/**
 * GetClientDatabase
 * Use to get DGraph client connection
 */
func GetDgraphClient() (*dgo.Dgraph, Types.CancelFunc) {
	conn, err := grpc.Dial(os.Getenv("DATABASE_HOST"), grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}