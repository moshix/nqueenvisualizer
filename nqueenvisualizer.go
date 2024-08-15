package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Step struct {
    Board [][]int `json:"board"`
}

func solveNQueens(n int, steps *[]Step) {
    board := make([][]int, n)
    for i := range board {
        board[i] = make([]int, n)
    }
    solve(0, n, board, steps)
}

func solve(row, n int, board [][]int, steps *[]Step) bool {
    if row == n {
        *steps = append(*steps, Step{copyBoard(board)})
        return true
    }

    for col := 0; col < n; col++ {
        if isSafe(row, col, n, board) {
            board[row][col] = 1
            *steps = append(*steps, Step{copyBoard(board)})

            if solve(row+1, n, board, steps) {
                return true
            }

            board[row][col] = 0
            *steps = append(*steps, Step{copyBoard(board)})
        }
    }
    return false
}

func isSafe(row, col, n int, board [][]int) bool {
    for i := 0; i < row; i++ {
        if board[i][col] == 1 {
            return false
        }
    }

    for i, j := row, col; i >= 0 && j >= 0; i, j = i-1, j-1 {
        if board[i][j] == 1 {
            return false
        }
    }

    for i, j := row, col; i >= 0 && j < n; i, j = i-1, j+1 {
        if board[i][j] == 1 {
            return false
        }
    }

    return true
}

func copyBoard(board [][]int) [][]int {
    n := len(board)
    copy := make([][]int, n)
    for i := range board {
        copy[i] = make([]int, n)
        for j := range board[i] {
            copy[i][j] = board[i][j]
        }
    }
    return copy
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)

    nStr := r.URL.Query().Get("n")
    n, err := strconv.Atoi(nStr)
    if err != nil || n < 1 {
        http.Error(w, "Invalid 'n' parameter", http.StatusBadRequest)
        return
    }

    steps := []Step{}
    solveNQueens(n, &steps)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(steps)
}

func main() {
    pwd, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/solve", handler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Serving file: %s", r.URL.Path)
        if r.URL.Path == "/" {
            http.ServeFile(w, r, filepath.Join(pwd, "index.html"))
        } else {
            http.ServeFile(w, r, filepath.Join(pwd, r.URL.Path[1:]))
        }
    })

    log.Println("Listening on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
