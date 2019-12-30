package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "strconv"
	"log"

	"carsMongoDB/models"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note: username and password are provided by Atlas account. Free online instance of MongoDB
const connectionString = "mongodb+srv://<Username>:<Password>@cluster0-kwb5x.mongodb.net/test?retryWrites=true&w=majority"

// Database Name
const databaseName = "carsDB"

// Collection name
const collectionName = "cars"

// collection object/instance
var carsDB *mongo.Collection

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("Error with connect - ", err)
	}

	// Ping the server
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error with ping - ", err)
	}

	carsDB = client.Database(databaseName).Collection(collectionName)

	fmt.Println("Collection instance created!")
}

// ReadAllCars will return all cars
func ReadAllCars(w http.ResponseWriter, r *http.Request) {
	// Set the way we will serve data between frontend and backend
	w.Header().Set("Content-Type", "application/json")
	// Allow cross origin connections making the routes accessible for everyone
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := readAllCars()
	json.NewEncoder(w).Encode(payload)
}
func readAllCars() []primitive.M {
	// Find all cars from database
	cursor, err := carsDB.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Error in db.Find() - ", err)
	}

	var cars []primitive.M
	for cursor.Next(context.Background()) {
		var car bson.M
		errorDecode := cursor.Decode(&car)
		if errorDecode != nil {
			log.Fatal("Error in decode - ", errorDecode)
		}
		cars = append(cars, car)
	}

	errorCursor := cursor.Err()
	if errorCursor != nil {
		log.Fatal("Error in cursor - ", errorCursor)
	}
	// Close the cursor when the function is done executing
	defer cursor.Close(context.Background())
	return cars
}

// ReadCar will return one car
func ReadCar(w http.ResponseWriter, r *http.Request) {
	// Set the way we will serve data between frontend and backend
	w.Header().Set("Content-Type", "application/json")
	// Allow cross origin connections making the routes accessible for everyone
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Get all params from url
	params := mux.Vars(r)
	// Retreive id from params and convert it from string into primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	payload := readCar(id)
	json.NewEncoder(w).Encode(payload)
}
func readCar(id primitive.ObjectID) models.Car {
	fmt.Println(id)
	var car models.Car
	// Find a car from database
	errorFindOne := carsDB.FindOne(context.Background(), bson.M{"_id": id}).Decode(&car)
	if errorFindOne != nil {
		log.Fatal("Error in decode - ", errorFindOne)
	}
	return car
}

// CreateCar will add a new car to the slice
func CreateCar(w http.ResponseWriter, r *http.Request) {
	// Set the way we will serve data between frontend and backend
	w.Header().Set("Content-Type", "application/json")
	// Allow cross origin connections making the routes accessible for everyone
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow the server to perform post operation
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	// Allow the content type that is specified by client to be processed on server
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Declare an empty car
	var car models.Car
	// Take the car json from the client and decode it into car struct
	_ = json.NewDecoder(r.Body).Decode(&car)
	payload := createCar(car)
	json.NewEncoder(w).Encode(payload)
}
func createCar(car models.Car) models.Car {
	insertResult, errorInsert := carsDB.InsertOne(context.Background(), car)
	if errorInsert != nil {
		log.Fatal(errorInsert)
	}
	// Add an ID field to the newly created car
	// MongoDB returns special interface that needs to be mapped to our ID type which is primitive.ObjectID
	car.ID = insertResult.InsertedID.(primitive.ObjectID)
	return car
}

// UpdateCar will modify a car chosen by id
func UpdateCar(w http.ResponseWriter, r *http.Request) {
	// Set the way we will serve data between frontend and backend
	w.Header().Set("Content-Type", "application/json")
	// Allow cross origin connections making the routes accessible for everyone
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow the server to perform post operation
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	// Allow the content type that is specified by client to be processed on server
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Get all params from url
	params := mux.Vars(r)
	// Declare an empty car
	var car models.Car
	// Take the car json from the client and decode it into car struct
	_ = json.NewDecoder(r.Body).Decode(&car)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	payload := updateCar(id, car)
	json.NewEncoder(w).Encode(payload)
}
func updateCar(id primitive.ObjectID, car models.Car) models.Car {
	// Update a value from database, ignore the result from action
	_, errorUpdateOne := carsDB.UpdateOne(
		context.Background(),
		models.Car{ID: id},
		bson.M{
			"$set": bson.M{
				"model": car.Model,
				"price": car.Price,
			},
		},
	)
	// Check for errors
	if errorUpdateOne != nil {
		log.Fatal("Error with UpdateOne - ", errorUpdateOne)
	}
	// Build up a car with id from params and fields from the client
	car.ID = id
	return car
}

// DeleteCar will remove a car from db by its id
func DeleteCar(w http.ResponseWriter, r *http.Request) {
	// Set the way we will serve data between frontend and backend
	w.Header().Set("Content-Type", "application/json")
	// Allow cross origin connections making the routes accessible for everyone
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow the server to perform post operation
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	// Allow the content type that is specified by client to be processed on server
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Get all params from url
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	payload := deleteCar(id)
	json.NewEncoder(w).Encode(payload)
}
func deleteCar(id primitive.ObjectID) string {
	_, errorDelete := carsDB.DeleteOne(context.Background(), bson.M{"_id": id})
	if errorDelete != nil {
		log.Fatal(errorDelete)
	}
	return "Removed a car"
}
