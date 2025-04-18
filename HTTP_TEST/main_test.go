package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestMain(m *testing.M) {
	// Initialize mock data before tests
	posts = []Post{
		{ID: 1, UserID: 1, Title: "Test Post", Body: "Test Content"},
	}
	currentID = 2
	m.Run()
}

func TestGetPost(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{id}", getPost).Methods("GET")

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   Post
	}{
		{"Existing post", "/posts/1", http.StatusOK, posts[0]},
		{"Non-existent post", "/posts/999", http.StatusNotFound, Post{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var got Post
				json.NewDecoder(resp.Body).Decode(&got)
				if got != tt.wantBody {
					t.Errorf("got body %v, want %v", got, tt.wantBody)
				}
			}
		})
	}
}

func TestGetPosts(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts", getPosts).Methods("GET")

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantCount  int
	}{
		{"All posts", "/posts", http.StatusOK, 1},
		{"Filter by user", "/posts?userId=1", http.StatusOK, 1},
		{"No results", "/posts?userId=999", http.StatusOK, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			var posts []Post
			json.NewDecoder(resp.Body).Decode(&posts)
			if len(posts) != tt.wantCount {
				t.Errorf("got %d posts, want %d", len(posts), tt.wantCount)
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts", createPost).Methods("POST")

	t.Run("Create new post", func(t *testing.T) {
		body := `{"title":"New","body":"Content","userId":1}`
		req := httptest.NewRequest("POST", "/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusCreated)
		}

		var post Post
		json.NewDecoder(resp.Body).Decode(&post)
		if post.ID != 2 {
			t.Errorf("got ID %d, want 2", post.ID)
		}
	})
}

func TestUpdatePost(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{id}", updatePost).Methods("PUT")

	tests := []struct {
		name       string
		url        string
		body       string
		wantStatus int
	}{
		{"Update existing", "/posts/1", `{"id":1,"title":"Updated","body":"Updated","userId":1}`, http.StatusOK},
		{"Update non-existent", "/posts/999", `{"id":999,"title":"Updated","body":"Updated","userId":1}`, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestPatchPost(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{id}", patchPost).Methods("PATCH")

	tests := []struct {
		name       string
		url        string
		body       string
		wantStatus int
	}{
		{"Patch existing", "/posts/1", `{"title":"Patched"}`, http.StatusOK},
		{"Patch non-existent", "/posts/999", `{"title":"Patched"}`, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PATCH", tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"Delete existing", "/posts/1", http.StatusNoContent},
		{"Delete non-existent", "/posts/999", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestGetComments(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{id}/comments", getComments).Methods("GET")

	req := httptest.NewRequest("GET", "/posts/1/comments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var comments []Comment
	json.NewDecoder(resp.Body).Decode(&comments)
	if len(comments) != 2 {
		t.Errorf("got %d comments, want 2", len(comments))
	}
}

func TestSetupRouter(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/posts/1", http.StatusOK},
		{"POST", "/posts", http.StatusCreated},
		{"PUT", "/posts/1", http.StatusOK},
		{"PATCH", "/posts/1", http.StatusOK},
		{"DELETE", "/posts/1", http.StatusNoContent},
		{"GET", "/posts/1/comments", http.StatusOK},
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, ts.URL+tt.path, nil)
			if tt.method == "POST" || tt.method == "PUT" || tt.method == "PATCH" {
				req.Body = http.NoBody
				req.Header.Set("Content-Type", "application/json")
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func TestStartServer(t *testing.T) {
	router := setupRouter()

	// Create a test server with a random port
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Verify server is running
	resp, err := http.Get(testServer.URL + "/posts/1")
	if err != nil {
		t.Fatalf("Server not running: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
