package repository

import (
	"CRUD/model"
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newMongoClient() *mongo.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	uri := os.Getenv("MONGO_URI")
	log.Println(uri)
	mongoTestClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln("Error while connecting to MongoDB", err)
	}

	log.Println("MongoDB connected successfully")

	err = mongoTestClient.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatalln("Ping failed", err)
	}

	log.Println("Ping successful")

	return mongoTestClient
}

func TestMongoOperation(t *testing.T) {
	mongoTestClient := newMongoClient()
	defer mongoTestClient.Disconnect(context.Background())

	emp1 := uuid.New().String()
	//emp2 := uuid.New().String()

	coll := mongoTestClient.Database("companydb").Collection("employee_test")

	empRepo := EmployeeRepo{MongoCollection: coll}

	t.Run("InsertEmployee", func(t *testing.T) {
		emp := model.Employee{
			Name:       "Tony Stark",
			Department: "Engineering",
			EmployeeID: emp1,
		}

		result, err := empRepo.InsertEmployee(&emp)

		if err != nil {
			t.Errorf("Error while inserting employee: %s", err)
		}

		t.Log("Inserted employee with ID: ", result)
	})

	t.Run("Get employee 1", func(t *testing.T) {
		result, err := empRepo.FindEmployeeById(emp1)

		if err != nil {
			t.Fatal("Get operation failed", err)
		}

		t.Log("Employee 1: ", result)
	})

	t.Run("Get all employees", func(t *testing.T) {
		results, err := empRepo.FindAllEmployees()

		if err != nil {
			t.Fatal("Get all operation failed", err)
		}

		t.Log("All employees: ", results)
	})

	t.Run("Update employee 1", func(t *testing.T) {
		emp := model.Employee{
			Name:       "Ngo Nam",
			Department: "IT",
			EmployeeID: emp1,
		}

		result, err := empRepo.UpdateEmployeeById(emp1, &emp)

		if err != nil {
			t.Fatal("Update operation failed", err)
		}

		t.Log("Updated employee 1: ", result)
	})

	t.Run("Get employee 1 after update", func(t *testing.T) {
		result, err := empRepo.FindEmployeeById(emp1)

		if err != nil {
			t.Fatal("Get operation failed", err)
		}

		t.Log("Employee 1: ", result)
	})

	t.Run("Delete employee 1", func(t *testing.T) {
		result, err := empRepo.DeleteEmployeeById(emp1)

		if err != nil {
			t.Fatal("Delete operation failed", err)
		}

		t.Log("Deleted employee 1: ", result)
	})

	t.Run("Delete all employees", func(t *testing.T) {
		result, err := empRepo.DeleteAllEmployees()

		if err != nil {
			t.Fatal("Delete all operation failed", err)
		}

		t.Log("Deleted all employees: ", result)
	})
}
