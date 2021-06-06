package middleware

import (
	"encoding/json"
	"fmt"
	"go-crud-mongodb/models"
	mongoclient "go-crud-mongodb/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func CreateUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    var user models.User

    user.ID = primitive.NewObjectID()

    err := json.NewDecoder(r.Body).Decode(&user)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to decode the request body.  %v", err))
        return 
    }

    if user.Name == "" || user.Location == "" || user.Age == 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode("One or more values are empty")
        return
    }

    insertID, err := mongoclient.InsertUser(user)
    
    if err != nil {
        fmt.Printf("Unable to add user. %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to add user. %v", err))
        return
    }

    res := models.Response{
        ID:      insertID.Hex(),
        Message: "User created successfully",
    }

    json.NewEncoder(w).Encode(res)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    params := mux.Vars(r)
    
    objID, err := primitive.ObjectIDFromHex(params["id"])

    if err != nil {
        log.Printf("Invalid object id.  %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to get object ID.  %v", err))
        return
    }

    user, err := mongoclient.GetUser(objID)

    if err != nil {
        fmt.Printf("Unable to get user. %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to get user. %v", err))
        return
    }

    json.NewEncoder(w).Encode(user)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    users, err := mongoclient.GetAllUsers()

    if err != nil {
        fmt.Printf("Unable to get all users. %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to get all users. %v", err))
        return
    }

    json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    params := mux.Vars(r)

    objID, err := primitive.ObjectIDFromHex(params["id"])

    if err != nil {
        fmt.Printf("Unable to get object ID.  %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to get object ID.  %v", err))
        return
    }

    var user models.User

    err = json.NewDecoder(r.Body).Decode(&user)

    if err != nil {
        fmt.Printf("Unable to decode the request body.  %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to decode the request body.  %v", err))
        return
    }

    id, err := mongoclient.UpdateUser(objID, user)

    if err != nil {
        fmt.Printf("Unable to update user. %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to update user. %v", err))
        return
    }

    res := models.Response{
        ID:      id.Hex(),
        Message: "Updated Successfully",
    }

    json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    params := mux.Vars(r)

    objID, err := primitive.ObjectIDFromHex(params["id"])

    if err != nil {
        fmt.Printf("Unable to get object ID.  %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to get object ID.  %v", err))
        return
    }

    deletedRows, err := mongoclient.DeleteUser(objID)

    if err != nil {
        fmt.Printf("Unable to Delete User.  %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(fmt.Sprintf("Unable to Delete User.  %v", err))
        return
    }

    fmt.Println(deletedRows)
    res := models.Response{
        ID:      objID.Hex(),
        Message: "User Deleted Successfully",
    }

    json.NewEncoder(w).Encode(res)
}