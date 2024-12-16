package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the structure of a Student info
type Student struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	PhoneNumber string             `json:"phoneNumber"`
	Email       string             `json:"email"`
	Address     string             `json:"address"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello, world!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB Atlas")

	collection = client.Database("golang_crud_student_info_db").Collection("students")

	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // Your frontend URL
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	app.Get("/api/students", getStudents)
	app.Post("/api/students", createStudent)
	app.Patch("/api/students/:id", updateStudent)
	app.Delete("/api/students/:id", deleteStudent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

// FETCH STUDENTS
func getStudents(c *fiber.Ctx) error {
	var students []Student

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var student Student
		if err := cursor.Decode(&student); err != nil {
			return err
		}

		students = append(students, student)
	}

	return c.JSON(students)
}

// CREATE STUDENT
func createStudent(c *fiber.Ctx) error {
	student := new(Student)

	if err := c.BodyParser(student); err != nil {
		return err
	}

	// Parse the JSON body from the incoming HTTP request into the Student struct
	if err := c.BodyParser(student); err != nil {
		// If parsing fails, return a 400 (Bad Request) error with a helpful message
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate all fields are not empty and handle edge cases
	if student.FirstName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Student First Name is required"})
	}

	if student.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Student Last Name is required"})
	}

	if student.PhoneNumber == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Student Phone Number is required"})
	} else if len(student.PhoneNumber) < 10 || len(student.PhoneNumber) > 15 {
		// Ensure phone number length is valid (10-15 characters for flexibility)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid phone number format"})
	}

	if student.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Student Email is required"})
	}

	if student.Address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Student Address is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), student)
	if err != nil {
		return err
	}

	student.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(student)
}

// UPDATE STUDENT
func updateStudent(c *fiber.Ctx) error {
	// Extract the student ID from the URL parameter
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	// Create a new instance to parse the incoming update data
	updateStudent := &Student{}
	// Parse the JSON body from the incoming HTTP request
	if err := c.BodyParser(updateStudent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Prepare the update document
	update := bson.M{}
	unset := bson.M{}

	// Conditionally add fields to update or remove old fields
	if updateStudent.FirstName != "" {
		update["firstName"] = updateStudent.FirstName
		unset["firstname"] = "" // Remove the old field
	}
	if updateStudent.LastName != "" {
		update["lastName"] = updateStudent.LastName
		unset["lastname"] = "" // Remove the old field
	}
	if updateStudent.PhoneNumber != "" {
		// Validate phone number if it's provided
		if len(updateStudent.PhoneNumber) < 10 || len(updateStudent.PhoneNumber) > 15 {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid phone number format"})
		}
		update["phoneNumber"] = updateStudent.PhoneNumber
		unset["phonenumber"] = "" // Remove the old field
	}
	if updateStudent.Email != "" {
		update["email"] = updateStudent.Email
	}
	if updateStudent.Address != "" {
		update["address"] = updateStudent.Address
	}

	// If no fields to update, return an error
	if len(update) == 0 && len(unset) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No update fields provided"})
	}

	// Perform the update
	filter := bson.M{"_id": objectId}
	updateResult, err := collection.UpdateOne(
		context.Background(),
		filter,
		bson.M{
			"$set":   update, // Add or update new fields
			"$unset": unset,  // Remove old fields
		},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update student"})
	}

	// Check if a document was actually updated
	if updateResult.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
	}

	// Retrieve the updated student to return to the client
	var updatedStudent Student
	err = collection.FindOne(context.Background(), filter).Decode(&updatedStudent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve updated student"})
	}

	// Respond with the updated student
	return c.Status(200).JSON(updatedStudent)
}

func deleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
