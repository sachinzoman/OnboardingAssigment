package main

import (
	"io"
	"fmt"
	"log"
	pb "restaurant_listing/restaurant"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	addr = ":8080"
)

func createRestaurant(client pb.RestaurantClient,restaurant *pb.RestaurantRequest){
	res,err := client.CreateRestaurant(context.Background(),restaurant)
	if err!=nil{
		log.Printf("Unable to proceed request: %v",err)
	}
	if res.Status == "200"{
		log.Printf("Succesfully created!")
	}
}

func getRestaurant(client pb.RestaurantClient, emp *pb.EmptyParam){
	stream, err := client.GetRestaurant(context.Background(), emp)
	if err!=nil{
		log.Printf("Unable to fetch data: %v", err)
	}
	for{
		restaurant,err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Fatalf("%v.GetRestaurant(_) = _, %v", client, err)
		}
		// log.Printf("Restaurant: %v", restaurant)
		log.Printf("Restaurants: %v",restaurant.Name)
		// for key,value:= range(restaurant){
		// 	log.Printf("%v: %v",key,value)
		// }
	}
}

func main(){
	//connection to grpc server
	conn,err := grpc.Dial(addr,grpc.WithInsecure())
	if err!=nil{
		fmt.Println("Unable to connect: %v",err)
	}

	defer conn.Close()

	client := pb.NewRestaurantClient(conn)

	restaurant := &pb.RestaurantRequest{
		Id :  0,
		Name : "Rajinder Da Dhaba",
		Rating : 4.1,
		Cusines : "North Indian, Rolls",
		Address : "AB-14B, Nauroji Nagar Marg, Opposite, Safdarjung Enclave, New Delhi, Delhi 110029",
		Times: &pb.RestaurantRequestTiming{
			Starttime: &pb.RestaurantRequestTimingTime{
				Hr: 17,
				Min: 24,
			},
			Endtime: &pb.RestaurantRequestTimingTime{
				Hr: 22,
				Min: 30,
			},
		},
		Cft : 800,
		ImgUrl: "https://b.zmtcdn.com/data/pictures/9/7319/e1b7673ed0aa2993b55b177409d5596c.jpg",
	}
	createRestaurant(client,restaurant)

	restaurant = &pb.RestaurantRequest{
		Id :  0,
		Name : "Prem Dhaba",
		Rating : 4,
		Cusines : "North Indian, Kebab",
		Address : "11139, East Park Road, Opposite JD Titler School, Karol Bagh, New Delhi",
		Times: &pb.RestaurantRequestTiming{
			Starttime: &pb.RestaurantRequestTimingTime{
				Hr: 12,
				Min: 0,
			},
			Endtime: &pb.RestaurantRequestTimingTime{
				Hr: 22,
				Min: 30,
			},
		},
		Cft : 650,
	}
	createRestaurant(client,restaurant)

	restaurant = &pb.RestaurantRequest{
		Id :  0,
		Name : "Karim's",
		Rating : 4,
		Cusines : "Mughlai, North Indian, Kebab, Rolls",
		Address : "16, Gali Kababian, Jama Masjid, New Delhi",
		Times: &pb.RestaurantRequestTiming{
			Starttime: &pb.RestaurantRequestTimingTime{
				Hr: 9,
				Min: 0,
			},
			Endtime: &pb.RestaurantRequestTimingTime{
				Hr: 0,
				Min: 30,
			},
		},
		Cft : 800,
		ImgUrl: "https://b.zmtcdn.com/data/pictures/6/18933156/e7256704bca7767f5475a7ec7cef9a94.jpg",
	}
	createRestaurant(client,restaurant)

	getRestaurant(client,&pb.EmptyParam{})
}