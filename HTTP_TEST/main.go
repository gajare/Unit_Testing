package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	ID     int    `json:"id"`
	PostID int    `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

var posts []Post
var currentID int = 1

func main() {
	r := mux.NewRouter()

	// Initialize some mock data
	posts = []Post{
		{
			ID:     1,
			UserID: 1,
			Title:  "Sample Post",
			Body:   "This is a sample post",
		},
	}
	currentID = 2

	// Define routes
	r.HandleFunc("/posts/{id}", getPost).Methods("GET")
	r.HandleFunc("/posts", getPosts).Methods("GET")
	r.HandleFunc("/posts", createPost).Methods("POST")
	r.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", patchPost).Methods("PATCH")
	r.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")
	r.HandleFunc("/posts/{id}/comments", getComments).Methods("GET")

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for _, post := range posts {
		if post.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(post)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	var filteredPosts []Post

	if userId != "" {
		uid, _ := strconv.Atoi(userId)
		for _, post := range posts {
			if post.UserID == uid {
				filteredPosts = append(filteredPosts, post)
			}
		}
	} else {
		filteredPosts = posts
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredPosts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	newPost.ID = currentID
	currentID++
	posts = append(posts, newPost)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updatedPost Post
	json.NewDecoder(r.Body).Decode(&updatedPost)

	for i, post := range posts {
		if post.ID == id {
			updatedPost.ID = id
			posts[i] = updatedPost
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedPost)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func patchPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updates map[string]interface{}
	json.NewDecoder(r.Body).Decode(&updates)

	for i, post := range posts {
		if post.ID == id {
			if title, ok := updates["title"].(string); ok {
				posts[i].Title = title
			}
			if body, ok := updates["body"].(string); ok {
				posts[i].Body = body
			}
			if userId, ok := updates["userId"].(float64); ok {
				posts[i].UserID = int(userId)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(posts[i])
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for i, post := range posts {
		if post.ID == id {
			posts = append(posts[:i], posts[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func getComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["id"])

	comments := []Comment{
		{
			ID:     1,
			PostID: postID,
			Name:   "John Doe",
			Email:  "john@example.com",
			Body:   "Sample comment 1",
		},
		{
			ID:     2,
			PostID: postID,
			Name:   "Jane Smith",
			Email:  "jane@example.com",
			Body:   "Sample comment 2",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
