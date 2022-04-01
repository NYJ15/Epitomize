package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/pilinux/gorest/controller"
	"github.com/pilinux/gorest/database"
	"github.com/pilinux/gorest/database/model"
)

var posts model.Post
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	Userid   uint
	jwt.StandardClaims
}

func Authentification(w http.ResponseWriter, r *http.Request) uint {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return http.StatusUnauthorized
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	tknStr := reqToken
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return http.StatusUnauthorized
		}
		w.WriteHeader(http.StatusBadRequest)
		return http.StatusBadRequest
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return http.StatusUnauthorized
	}
	fmt.Println("userid", claims.Userid)
	return claims.Userid
}
func AllPostTest(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	fmt.Println(userid)
	posts := controller.GetPosts(userid, false)
	result := model.GetAllPost{
		Posts: posts,
	}
	json.NewEncoder(w).Encode(result)
}
func AllPost(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	fmt.Println(userid)
	posts := controller.GetPosts(userid, true)
	result := model.GetAllPost{
		Posts: posts,
	}
	json.NewEncoder(w).Encode(result)
}

func TopTagsTest(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	topTag, responseType := controller.GetTopTags(false)
	if responseType == http.StatusOK {
		result := make(map[string][]string)
		result["TagList"] = topTag
		json.NewEncoder(w).Encode(result)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
func AllTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topTag, responseType := controller.GetAllTags(true)
	if responseType == http.StatusOK {
		result := make(map[string][]string)
		result["TagList"] = topTag
		json.NewEncoder(w).Encode(result)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
func TopTags(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	topTag, responseType := controller.GetTopTags(true)
	if responseType == http.StatusOK {
		result := make(map[string][]string)
		result["TagList"] = topTag
		json.NewEncoder(w).Encode(result)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
func CreateNewPostTest(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	responseType := controller.CreatePost(post, 1, false)
	json.NewEncoder(w).Encode(http.StatusText(responseType))
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		errorResponse := controller.CreateUser(user)
		if errorResponse.HTTPCode == http.StatusOK {
			json.NewEncoder(w).Encode(errorResponse.Message)
			return
		}
		http.Error(w, errorResponse.Message, errorResponse.HTTPCode)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	user, responseType := controller.GetUser(userid)

	if responseType == http.StatusOK {
		json.NewEncoder(w).Encode(user)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func SearchUserPost(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	var search model.Search
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responseType := controller.SearchPost(search)
		result := model.GetAllSearchPost{
			Posts: responseType,
		}
		json.NewEncoder(w).Encode(result)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var login model.Login
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responseType := controller.Login(login)
		json.NewEncoder(w).Encode(responseType)
	}
}

func GetUserFeed(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	posts, responseType := controller.GetUserFeed(userid)
	if responseType == http.StatusOK {
		json.NewEncoder(w).Encode(posts)
		return
	}
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func GetUserRecommendations(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	posts, responseType := controller.GetUserRecommendations(userid)
	if responseType == http.StatusOK {
		json.NewEncoder(w).Encode(posts)
		return
	}
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func UserList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("userid")
	userid := Authentification(w, r)
	fmt.Println(userid)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	responseType := controller.UserList(userid)
	result := model.UserListResponses{
		Users: responseType,
	}

	json.NewEncoder(w).Encode(result)
}

func CreateNewPost(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	fmt.Println(userid)
	fmt.Println(userid)
	var post model.Post
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responseType := controller.CreatePost(post, userid, true)
		json.NewEncoder(w).Encode(http.StatusText(responseType))
	}
}

func GetPostTest(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	post, responseType := controller.GetPost(1, false)
	if responseType == http.StatusOK {
		json.NewEncoder(w).Encode(post)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	params := mux.Vars(r)
	id := params["id"]
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, responseType := controller.GetPost(postId, true)
	if responseType == http.StatusOK {
		json.NewEncoder(w).Encode(post)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func EditPostTest(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	var post model.Post
	err, responseType := controller.EditPost(1, post, code, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	} else {
		json.NewEncoder(w).Encode(http.StatusText(responseType))
	}
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	var post model.Post
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		postId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err, responseType := controller.EditPost(postId, post, code, true)
		// result := http.Response{
		// 	StatusCode: responseType,
		// 	Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		// }
		// if err != nil {
		// 	json.NewEncoder(w).Encode(err.Error())
		// } else {
		// 	json.NewEncoder(w).Encode(result)
		// }
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			json.NewEncoder(w).Encode(http.StatusText(responseType))
		}
	}
}

func DeletePostTest(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	params := mux.Vars(r)
	id := params["id"]
	controller.DeletePost(id, false)
}
func FollowUser(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	params := mux.Vars(r)
	id := params["userid"]
	Val, _ := strconv.ParseUint(id, 10, 64)
	responseType := controller.FollowUser(uint(Val), userid)
	json.NewEncoder(w).Encode(http.StatusText(responseType))
}
func UnFollowUser(w http.ResponseWriter, r *http.Request) {
	userid := Authentification(w, r)
	if userid == http.StatusUnauthorized || userid == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(userid)), int(userid))
		return
	}
	params := mux.Vars(r)
	id := params["userid"]
	Val, _ := strconv.ParseUint(id, 10, 64)
	responseType := controller.UnFollowUser(uint(Val), userid)
	json.NewEncoder(w).Encode(http.StatusText(responseType))
}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	code := Authentification(w, r)
	if code == http.StatusUnauthorized || code == http.StatusBadRequest {
		http.Error(w, http.StatusText(int(code)), int(code))
		return
	}
	params := mux.Vars(r)
	id := params["id"]
	controller.DeletePost(id, true)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HomePage")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
		log.Println("Executing middleware again")
	})
}

func HandleRequests() {
	myRouter := mux.NewRouter()
	// myRouter.Use(accessControlMiddleware)
	myRouter.HandleFunc("/", HomePage)
	myRouter.HandleFunc("/post", AllPost).Methods("GET")
	myRouter.HandleFunc("/alltags", AllTags).Methods("GET")
	myRouter.HandleFunc("/topTags", TopTags).Methods("GET")
	myRouter.HandleFunc("/post/{id}", GetPost).Methods("GET")
	myRouter.HandleFunc("/post", CreateNewPost).Methods("POST")
	myRouter.HandleFunc("/post/{id}", EditPost).Methods("PUT")
	myRouter.HandleFunc("/deleteposts/{id}", DeletePost).Methods("DELETE")
	myRouter.HandleFunc("/login", LoginUser).Methods("POST")
	myRouter.HandleFunc("/userlist", UserList).Methods("GET")
	myRouter.HandleFunc("/user", CreateNewUser).Methods("POST")
	myRouter.HandleFunc("/user", GetUser).Methods("GET")
	myRouter.HandleFunc("/follow/{userid}", FollowUser).Methods("GET")
	myRouter.HandleFunc("/unfollow/{userid}", UnFollowUser).Methods("GET")
	myRouter.HandleFunc("/user/feed", GetUserFeed).Methods("GET")
	myRouter.HandleFunc("/user/recommended", GetUserRecommendations).Methods("GET")
	myRouter.HandleFunc("/search", SearchUserPost).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", CorsMiddleware(myRouter)))
}

func main() {
	if err := database.InitDB().Error; err != nil {
		fmt.Println(err)
		return
	}
	HandleRequests()
}
