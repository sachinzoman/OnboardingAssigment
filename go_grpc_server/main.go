package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//protobuf file
	pb "restaurant_listing/go_grpc_server/protos"

	//for database
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	addr = "localhost:8080"
)

type server struct {
	db *sql.DB
}

func (s *server) connect(ctx context.Context) (*sql.Conn, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return conn, err
}

func (s *server) CreateRestaurant(ctx context.Context, req *pb.RestaurantRequest) (*pb.RestaurantResponse, error) {

	// log.Printf("Create Message Triggered: %v", time.Now())
	c, err := s.connect(ctx)
	if err != nil {
		log.Printf("Unable to connect to MySQL: %v", err)
		return nil, err
	}
	defer c.Close()
	// starttime := strconv.Itoa(int(req.Times.Starttime.Hr)) + ":" + strconv.Itoa(int(req.Times.Starttime.Min)) + ":00"
	// endtime := strconv.Itoa(int(req.Times.Endtime.Hr)) + ":" + strconv.Itoa(int(req.Times.Endtime.Min)) + ":00"
	query := "insert into restaurant(`restaurantName`, `rating`, `cusines`, `address`, `startime`, `endtime`, `cft`, `img_url`) values(?,?,?,?,?,?,?,?)"
	res, err := c.ExecContext(ctx, query, req.Name, req.Rating, req.Cusines, req.Address, req.Starttime, req.Endtime, req.Cft, req.ImgUrl)
	if err != nil {
		log.Panic(err)
		return &pb.RestaurantResponse{Status: "400"}, err
	}
	id, err := res.LastInsertId()

	log.Printf("Restaurnat creted with id: %v", id)
	return &pb.RestaurantResponse{Status: "201"}, nil
}

func (s *server) GetRestaurant(emp *pb.EmptyParam, stream pb.Restaurant_GetRestaurantServer) error {

	log.Println("Get requested from client triggred")
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
			Id:        restaurantId,
			Name:      restaurantName,
			Rating:    rating,
			Cusines:   cusines,
			Address:   address,
			Starttime: startime,
			Endtime:   endtime,
			Cft:       cft,
			ImgUrl:    img_url,
		}
		if err := stream.Send(result); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil

}

func (s *server) UpdateRestaurant(ctx context.Context, req *pb.RestaurantRequest) (*pb.RestaurantResponse, error) {

	//create connection to database
	conn, err := s.connect(ctx)
	if err != nil {
		log.Fatal("Could not connect to database: %v", err)
		return nil, err
	}

	//insert restaurant entry
	query := "update restaurant set restaurantName=?, rating=?, cusines=?, address=?, startime=?, endtime=?, cft=?, img_url=? where restaurantId=?"
	res, err := conn.ExecContext(ctx, query, req.Name, req.Rating, req.Cusines, req.Address, req.Starttime, req.Endtime, req.Cft, req.ImgUrl, req.Id)

	rows, err := res.RowsAffected()

	if err != nil {
		log.Fatal("Could not proceed restuarant update: %v", err)
		return &pb.RestaurantResponse{Status: "400"}, nil
	}

	if rows < 1 {
		return &pb.RestaurantResponse{Status: "400"}, nil
	}

	log.Print("Item inserted succesfully!, %v", res)

	return &pb.RestaurantResponse{Status: "200"}, nil
}

func (s *server) DeleteRestaurant(ctx context.Context, req *pb.RestaurantRequest) (*pb.RestaurantResponse, error) {

	//create connection to database
	conn, err := s.connect(ctx)
	if err != nil {
		log.Fatal("Could not connect to database: %v", err)
		return nil, err
	}

	//delete given entry
	res, err := conn.ExecContext(ctx, "delete from restaurant where restaurantId= ?", req.Id)
	if err != nil {
		log.Fatal("Could not proceed restuarant update: %v", err)
		return &pb.RestaurantResponse{Status: "400"}, nil
	}

	rows, err := res.RowsAffected()

	if rows < 1 {
		return &pb.RestaurantResponse{Status: "400"}, nil
	}

	log.Print("Item deleted succesfully!, %v", res)

	return &pb.RestaurantResponse{Status: "200"}, nil
}

func main() {

	//setup mysql database
	db, err := sql.Open("mysql", "jarvis:jarvis123@tcp(0.0.0.0)/restaurantDB")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err)
	}

	res, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		log.Fatal(err)
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
			break
		} else {
			log.Println(restaurantName)
		}
	}

	// defer the close till after the main function has finished
	defer db.Close()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Server unable to listen: %v", err)
	} else {
		fmt.Println("==>Server is up and running!")
	}
	s := grpc.NewServer()
	pb.RegisterRestaurantServer(s, &server{db: db})
	s.Serve(lis)
}
