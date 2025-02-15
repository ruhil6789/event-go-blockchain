package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ruhil6789/event-sky/database"
	"github.com/ruhil6789/event-sky/pkg/model"
	response "github.com/ruhil6789/event-sky/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection= database.GetCollection(database.DB,"users")
var validate = validator.New()




// Ctx represents the Context which hold the HTTP request and response. It has methods for the request query string, parameters, body, HTTP headers and so on.
func CreateUser(c *fiber.Ctx)  error{
	ctx,cancel:=context.WithTimeout(context.Background(),10 *time.Second)
	var user model.User
	defer cancel()

	if err := c.BodyParser(&user); err!= nil {
   
		return c.Status(http.StatusBadRequest).JSON(response.ApiResponse{Message: "error",Data: &fiber.Map{"data":err.Error()},Status: http.StatusBadRequest})

	}

	    //use the validator library to validate required fields
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(response.ApiResponse{ Message: "error", Data: &fiber.Map{"data": validationErr.Error()},Status: http.StatusBadRequest})
    }

	 newUser := model.User{
        Id:       primitive.NewObjectID(),
        Name:     user.Name,
        Location: user.Location,
        Title:    user.Title,
    }

    result, err := userCollection.InsertOne(ctx, newUser)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(response.ApiResponse{ Message: "error", Data: &fiber.Map{"data": err.Error()},Status: http.StatusInternalServerError})
    }

	 return c.Status(http.StatusCreated).JSON(response.ApiResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    var user model.User
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(response.ApiResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    return c.Status(http.StatusOK).JSON(response.ApiResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func EditAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    var user model.User
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    //validate the request body
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(response.ApiResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(response.ApiResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

    update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

    result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(response.ApiResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }
    //get updated user details
    var updatedUser model.User
    if result.MatchedCount == 1 {
        err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(response.ApiResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }
    }

    return c.Status(http.StatusOK).JSON(response.ApiResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}