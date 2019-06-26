package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//protobuf file
	pb "restaurant_listing/restaurant"

	//for database
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	addr = ":8080"
)

type server struct {
	db *sql.DB
	// savedRestaurants []*pb.RestaurantRequest
}

func (s *server) connect(ctx context.Context) (*sql.Conn, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return conn, err
}

func (s *server) CreateRestaurant(ctx context.Context, req *pb.RestaurantRequest) (*pb.RestaurantResponse, error) {

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	starttime := strconv.Itoa(int(req.Times.Starttime.Hr)) + ":" + strconv.Itoa(int(req.Times.Starttime.Min)) + ":00"
	endtime := strconv.Itoa(int(req.Times.Endtime.Hr)) + ":" + strconv.Itoa(int(req.Times.Endtime.Min)) + ":00"
	res, err := c.ExecContext(ctx, "insert into restaurant(`restaurantName`, `rating`, `cusines`, `address`, `startime`, `endtime`, `cft`, `img_url`) values(?,?,?,?,?,?,?,?)", req.Name, req.Rating, req.Cusines, req.Address, starttime, endtime, req.Cft, req.ImgUrl)
	if err != nil {
		return &pb.RestaurantResponse{Status: "400"}, err
	}
	id, err := res.LastInsertId()
	log.Printf("Restaurnat creted with id: %v", id)

	// s.savedRestaurants = append(s.savedRestaurants, req)

	return &pb.RestaurantResponse{Status: "200"}, nil

}

func (s *server) GetRestaurant(emp *pb.EmptyParam, stream pb.Restaurant_GetRestaurantServer) error {

	//get a connection from connection pool
	c, err := s.connect(stream.Context())

	res, err := c.QueryContext(stream.Context(), "select * from restaurant")
	if err != nil {
		return nil
	}

	var restaurantId int32
	var restaurantName string
	var rating float32
	var cusines string
	var address string
	var startime string
	var endtime string
	var cft float32
	var img_url string

	for res.Next() {
		err := res.Scan(&restaurantId, &restaurantName, &rating, &cusines, &address, &startime, &endtime, &cft, &img_url)
		if err != nil {
			log.Fatal(err)
			return err
		}
		result := &pb.RestaurantRequest{
			Id:      restaurantId,
			Name:    restaurantName,
			Rating:  rating,
			Cusines: cusines,
			Address: address,
			Times: &pb.RestaurantRequestTiming{
				Starttime: &pb.RestaurantRequestTimingTime{
					Hr:  8,
					Min: 30,
				},
				Endtime: &pb.RestaurantRequestTimingTime{
					Hr:  22,
					Min: 0,
				},
			},
			Cft:    cft,
			ImgUrl: img_url,
		}
		if err := stream.Send(result); err != nil {
			return err
		}
	}
	return nil

}

func main() {

	//setup mysql database
	db, err := sql.Open("mysql", "jarvis:jarvis123@tcp(localhost)/restaurantDB")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	defer db.Close()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Server unable to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRestaurantServer(s, &server{db: db})
	s.Serve(lis)
}
