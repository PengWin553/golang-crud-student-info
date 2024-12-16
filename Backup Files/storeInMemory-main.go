package main

// Import necessary packages
import (
	"fmt" // Used for printing messages to the console
	"log" // Used for logging errors
	"os"
	"strconv" //to properly convert the ID from a string to an integer, which matches the struct's ID type

	"github.com/gofiber/fiber/v2" // The Fiber web framework for building APIs
	"github.com/joho/godotenv"
)

// Define the structure of a Student info
type Student struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

func main() {
	// Print a startup message to the console
	fmt.Println("Hello, World. Peng Win")

	// Create a new Fiber application instance
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	// A slice (dynamic array) to store Student items in memory
	students := []Student{}

	// GET ALL STUDENTS API
	app.Get("/api/students", func(c *fiber.Ctx) error {
		// Return all students as JSON
		return c.Status(200).JSON(students)
	})

	// CREATE STUDENT API
	app.Post("/api/students", func(c *fiber.Ctx) error {
		// Create a new instance of the Student struct
		student := &Student{}

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

		// Assign a unique ID to the new Student
		student.ID = len(students) + 1

		// Add the new Student item to the `students` slice
		students = append(students, *student)

		// Respond to the client with a 201 (Created) status and the new Student item in JSON format
		return c.Status(201).JSON(student)
	})

	// UPDATE STUDENT API
	app.Patch("/api/students/:id", func(c *fiber.Ctx) error {
		// Extract the "id" parameter from the URL path
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid student ID"})
		}

		// Create a new instance to parse the incoming update data
		updateStudent := &Student{}

		// Parse the JSON body from the incoming HTTP request
		if err := c.BodyParser(updateStudent); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Iterate over the `students` slice to find the Student with a matching ID
		for i, student := range students {
			if student.ID == id {
				// Perform partial update - only update fields that are provided
				if updateStudent.FirstName != "" {
					students[i].FirstName = updateStudent.FirstName
				}
				if updateStudent.LastName != "" {
					students[i].LastName = updateStudent.LastName
				}
				if updateStudent.PhoneNumber != "" {
					// Validate phone number if it's provided
					if len(updateStudent.PhoneNumber) < 10 || len(updateStudent.PhoneNumber) > 15 {
						return c.Status(400).JSON(fiber.Map{"error": "Invalid phone number format"})
					}
					students[i].PhoneNumber = updateStudent.PhoneNumber
				}
				if updateStudent.Email != "" {
					students[i].Email = updateStudent.Email
				}
				if updateStudent.Address != "" {
					students[i].Address = updateStudent.Address
				}

				// Respond to the client with the updated Student and a 200 (OK) status
				return c.Status(200).JSON(students[i])
			}
		}

		// If no Student with the matching ID is found:
		// - Respond with a 404 (Not Found) status code
		// - Include a JSON error message for the client
		return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
	})

	// DELETE STUDENT API
	app.Delete("/api/students/:id", func(c *fiber.Ctx) error {
		// Extract the "id" parameter from the URL path
		// c.Params("id") retrieves the value of the `:id` placeholder in the route
		id := c.Params("id")

		// Iterate over the `students` slice to find the Student with a matching ID
		for i, student := range students {
			// Use `fmt.Sprint` to convert the integer `student.ID` to a string for comparison
			if fmt.Sprint(student.ID) == id {
				// If the ID matches, remove the Student from the `students` slice
				// Use slicing to create a new slice excluding the matched item
				students = append(students[:i], students[i+1:]...)

				// Respond to the client with a 200 (OK) status and a success message
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		// If no Student with the matching ID is found:
		// - Respond with a 404 (Not Found) status code
		// - Include a JSON error message for the client
		return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
	})

	// Start the web server on port 4000
	// log.Fatal ensures the application stops if the server fails to start
	log.Fatal(app.Listen(":" + PORT))
}
