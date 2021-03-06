package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"users-admin/app/db"
	"users-admin/app/logger"
	"users-admin/app/model"

	"strconv"

	"github.com/gorilla/mux"
)

// Get all users.
func GetAllUsers(dbSource *sql.DB, w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("find all users end point")
	dao := &db.Dao{}
	users := dao.GetAllUsers(dbSource)
	respondJSON(w, http.StatusOK, users)
}

func FindUserById(dbSource *sql.DB, w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("find by user id end point")
	vars := mux.Vars(r)
	userid, err := strconv.ParseInt(vars["userid"], 10, 8)
	if err != nil {
		panic(err)
	}
	dao := &db.Dao{}
	user, found := dao.FindUserById(dbSource, userid)

	if found == true {
		logger.Info.Printf("user name found %s", user.Username)
		respondJSON(w, http.StatusOK, user)
	} else {
		respondJSON(w, http.StatusNoContent, nil)
	}
}

func CreateUser(dbSource *sql.DB, w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Create user end point")
	dao := &db.Dao{}
	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := dao.SaveUser(dbSource, &user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, u)
	defer r.Body.Close()
}

func ModifyUser(dbSource *sql.DB, w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Modify user end point")
	vars := mux.Vars(r)
	userid, errConv := strconv.ParseInt(vars["userid"], 10, 8)
	if errConv != nil {
		panic(errConv)
	}
	dao := &db.Dao{}
	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := dao.UpdateUser(dbSource, &user, userid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, u)
	defer r.Body.Close()
}

func DeleteUser(dbSource *sql.DB, w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Delete user end point")
	vars := mux.Vars(r)
	userid, errConv := strconv.ParseInt(vars["userid"], 10, 8)
	if errConv != nil {
		panic(errConv)
	}
	dao := &db.Dao{}
	errDelete := dao.DeleteUser(dbSource, userid)
	if errDelete != nil {
		respondError(w, http.StatusInternalServerError, errDelete.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func SearchUsersByFilters(dbSource *sql.DB, w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Search users end point")
	dao := &db.Dao{}
	usrSearchReq := model.UserSearchRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usrSearchReq); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := dao.GetAllUsersByFilters(dbSource, &usrSearchReq)
	userResp := model.UserSearchPaginator{}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userResp.TotalElements = len(users)
	userResp.Page = usrSearchReq.Page
	userResp.Users = users
	respondJSON(w, http.StatusOK, userResp)
	defer r.Body.Close()
}
