package main

import (
	"net"
	"fmt"
	"google.golang.org/grpc"
	"golang.org/x/net/context"

	//protobuf file
	pb "restaurant_listing/restaurant"

	//for database
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const(
	addr = ":8080"
)

type server struct{
	savedRestaurants []*pb.RestaurantRequest
}

func (s *server) CreateRestaurant(ctx context.Context, in *pb.RestaurantRequest) (*pb.RestaurantResponse,error){
	s.savedRestaurants = append(s.savedRestaurants,in)
	return &pb.RestaurantResponse{Status:"200"},nil
}

func (s *server) GetRestaurant(emp *pb.EmptyParam, stream pb.Restaurant_GetRestaurantServer) error{
	for _,restaurant := range(s.savedRestaurants){
		if err:=stream.Send(restaurant); err!=nil{
			return err
		}
	}
	return nil
}

func main(){

	//setup mysql database
	db,err := sql.Open("mysql","jarvis:jarvis123@tcp(localhost)/restaurant")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	defer db.Close()

	lis,err := net.Listen("tcp",addr)
	if err!=nil{
		fmt.Println("Server unable to listen: %v",err)
	}
	s := grpc.NewServer()
	pb.RegisterRestaurantServer(s,&server{})
	s.Serve(lis)
}