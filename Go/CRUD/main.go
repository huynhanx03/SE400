package main

import (
	"CRUD/usecase"
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
)

var mongoClient *mongo.Client

func init() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println(".env file loaded successfully")
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is empty")
	}

	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln("Error while connecting to MongoDB", err)
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatalln("Ping failed", err)
	}

	log.Println("MongoDB connected successfully")
}

func main() {
	defer mongoClient.Disconnect(context.Background())

	r := mux.NewRouter()

	coll := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))
	empService := usecase.EmployeeService{MongoCollection: coll}

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	r.HandleFunc("/employee", empService.CreateEmployee).Methods(http.MethodPost)
	r.HandleFunc("/employee/{id}", empService.GetEmployeeById).Methods(http.MethodGet)
	r.HandleFunc("/employee", empService.GetAllEmployee).Methods(http.MethodGet)
	r.HandleFunc("/employee/{id}", empService.UpdateEmployeeById).Methods(http.MethodPut)
	r.HandleFunc("/employee/{id}", empService.DeleteEmployeeById).Methods(http.MethodDelete)
	r.HandleFunc("/employee", empService.DeleteAllEmployee).Methods(http.MethodDelete)

	log.Print("Server is running on port 4444")
	http.ListenAndServe(":4444", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ngo Nam dz"))
}
