<?php
require './vendor/autoload.php';
// require './Restaurant/empty_param.php';
// require './Restaurant/RestaurantClient.php';
// // require './Restaurant/RestaurantRequest_timing_time.php';
// // require './Restaurant/RestaurantRequest_timing.php';
// require './Restaurant/RestaurantRequest.php';
// require './Restaurant/RestaurantResponse.php';
// require './GPBMetadata/Restaurant.php';

// $docker = getenv("ENV");
$client = new Restaurant\RestaurantClient(getenv("ENV")==="docker"?'host.docker.internal':"localhost".':4000',[
    'credentials' => Grpc\ChannelCredentials::createInsecure(),
]);

function decodeRestuarant($coded_restaurant){
    $restaurant = array();
    $restaurant["Name"] = $coded_restaurant->getName();
    $restaurant["Address"] = $coded_restaurant->getAddress();
    $restaurant["Id"] = $coded_restaurant->getId();
    $restaurant["Cusines"] = $coded_restaurant->getCusines();
    $restaurant["Rating"] = $coded_restaurant->getRating();
    $restaurant["Starttime"] = $coded_restaurant->getStarttime();
    $restaurant["Endtime"] = $coded_restaurant->getEndtime();
    $restaurant["Cft"] = $coded_restaurant->getCft();
    $restaurant["ImgUrl"] = $coded_restaurant->getImgUrl();
    return $restaurant;
}

function encodeRestaurant($restaurant){
    //create request variable
    $request = new Restaurant\RestaurantRequest();
    
    $request->setName(array_key_exists("Name",$restaurant)?$restaurant["Name"]:"");
    $request->setAddress(array_key_exists("Address",$restaurant)?$restaurant["Address"]:"");
    $request->setId(array_key_exists("Id",$restaurant)?$restaurant["Id"]:0);
    $request->setCusines(array_key_exists("Cusines",$restaurant)?$restaurant["Cusines"]:"");
    $request->setRating(array_key_exists("Rating",$restaurant)?$restaurant["Rating"]:0);
    $request->setStarttime(array_key_exists("Starttime",$restaurant)?$restaurant["Starttime"]:"");
    $request->setEndtime(array_key_exists("Endtime",$restaurant)?$restaurant["Endtime"]:"");
    $request->setCft(array_key_exists("Cft",$restaurant)?$restaurant["Cft"]:0);
    $request->setImgUrl(array_key_exists("ImgUrl",$restaurant)?$restaurant["ImgUrl"]:"");
    return $request;
}

function runCreateRestaurant($client){
    $requestBody = file_get_contents('php://input');    //get POST request body
    $requestBody = json_decode($requestBody,true);  //json object to php dictionary
    if(json_last_error() != JSON_ERROR_NONE){
        http_response_code(400);
        echo json_encode(array("Status"=>"Invalid json data"));
        return;
    }
    $requestBody["Id"] = 0; //since we are creating new restaurant and it is given by system
    $request = encodeRestaurant($requestBody);
    list($resp, $status) = $client->CreateRestaurant($request)->wait();
    if($resp->getStatus()=="201")
        http_response_code(201);
    else
        http_response_code(400);
    echo json_encode(array("Status" => $resp->getStatus()));
}

function runGetRestaurant($client){
    $empParam = new Restaurant\empty_param();
    $restaurants = $client->GetRestaurant($empParam)->responses();
    echo "Get request sent to go service";
    $decoded_restaurants = array();
    foreach($restaurants as $restaurant){
        $decoded_restaurant = decodeRestuarant($restaurant);
        array_push($decoded_restaurants,$decoded_restaurant);
    }
    http_response_code(200);
    echo json_encode($decoded_restaurants);
}

function runUpdateRestaurant($client){
    $requestBody = file_get_contents('php://input');
    $requestBody = json_decode($requestBody,true);  //json object to php dictionary
    if(json_last_error() != JSON_ERROR_NONE){
        http_response_code(400);
        echo json_encode(array("Status"=>"Invalid json data"));
        return;
    }
    $request = encodeRestaurant($requestBody);
    list($resp,$status) = $client->UpdateRestaurant($request)->wait();
    if($resp->getStatus()=="200"){
        http_response_code(200);
        echo json_encode(array("Status" => "Restaurant Updated"));
    }
    else{
        http_response_code(400);
        echo json_encode(array("Status"=>"Unable to proceed!"));
    }
}

function runDeleteRestaurant($client){
    $requestBody = file_get_contents('php://input');
    $requestBody = json_decode($requestBody,true);  //json object to php dictionary
    if(json_last_error() != JSON_ERROR_NONE){
        http_response_code(400);
        echo json_encode(array("Status"=>"Invalid json data"));
        return;
    }
    $request = encodeRestaurant($requestBody);
    list($resp,$status) = $client->DeleteRestaurant($request)->wait();
    if($resp->getStatus()=="200"){
        http_response_code(200);
        echo json_encode(array("Status" => "Restaurant Deleted"));
    }
    else{
        http_response_code(400);
        echo json_encode(array("Status"=>"Unable to Delete!"));
    }   
}

// //create request variable
// $request = new Restaurant\RestaurantRequest();
// $request->setId(10);
// $request->setName("ABCD");
// $request->setAddress("32/AB");
// $request->setCusines("South Indian, Chinese");
// $request->setRating(4.1);
// $request->setStarttime("08:00:30");
// $request->setEndtime("22:00:00");
// $request->setCft(800);
// $request->setImgUrl("https://b.zmtcdn.com/data/pictures/9/7319/e1b7673ed0aa2993b55b177409d5596c.jpg");

$request = $_SERVER['PATH_INFO'];
$method = $_SERVER['REQUEST_METHOD'];

if ($method=='POST' && $request=='/create'){
    runCreateRestaurant($client);
}
else if($method=='GET' && $request=='/read'){
    runGetRestaurant($client);
}
else if($method=='PUT' && $request=='/update'){
    runUpdateRestaurant($client,$request);
}
else if($method=="DELETE" && $request=='/delete'){
    runDeleteRestaurant($client,$request);
}
else{
    http_response_code(400);
    echo json_encode(array("error" => "404 bad request"));
}