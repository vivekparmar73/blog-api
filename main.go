package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BlogPost struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ImageURL  string    `json:"imageURL"`
}

var posts = make(map[int]BlogPost)
var idCounter = 1

func main() {
	http.HandleFunc("/posts", postsHandler)        // for GET and CREATE all posts
	http.HandleFunc("/posts/", postsDetailHandler) // for single post

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPosts(w, r)
	case http.MethodPost:
		createPost(w, r)
	default:
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func postsDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getPost(w, r, id)
	case http.MethodPut:
		updatePost(w, r, id)
	case http.MethodDelete:
		deletePost(w, r, id)
	default:
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}

}

func getPosts(w http.ResponseWriter, r *http.Request) {
	var allposts []BlogPost

	for _, p := range posts {
		allposts = append(allposts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allposts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post BlogPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post.ID = idCounter
	idCounter++
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	posts[post.ID] = post

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)

}

func getPost(w http.ResponseWriter, r *http.Request, id int) {
	post, ok := posts[id]

	if !ok {
		http.Error(w, "post not fount", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func updatePost(w http.ResponseWriter, r *http.Request, id int) {
	existingPost, ok := posts[id]
	if !ok {
		http.Error(w, "post not found!", http.StatusNotFound)
		return
	}

	var updatedPost BlogPost

	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingPost.Title = updatedPost.Title
	existingPost.Content = updatedPost.Content
	existingPost.Author = updatedPost.Author
	existingPost.ImageURL = updatedPost.ImageURL
	existingPost.UpdatedAt = time.Now()
	posts[id] = existingPost

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingPost)
}

func deletePost(w http.ResponseWriter, r *http.Request, id int) {
	if _, ok := posts[id]; !ok {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	delete(posts, id)
	w.WriteHeader(http.StatusNoContent)
}
