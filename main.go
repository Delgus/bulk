package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/lib/pq"
)

type Book struct {
	Id     int
	Title  string
	Author string
	Price  float64
}

func main() {
	// configuration
	var batchesCount, batchSize, startID int
	flag.IntVar(&batchesCount, "batch", 50, "count of batches")
	flag.IntVar(&batchSize, "size", 100, "size of batch")
	flag.IntVar(&startID, "startID", 1, "start id")
	flag.Parse()
	booksCount := batchesCount * batchSize

	// open connection with postgres
	db, err := sql.Open("postgres", "user=postgres password=123456 dbname=postgres host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// generate batches
	batches := make([][]Book, batchesCount)
	numBatch := 0
	for i := 0; i < booksCount; i++ {
		bookID := startID + i
		if i%batchSize == 0 && i != 0 {
			numBatch++
		}
		batches[numBatch] = append(batches[numBatch], Book{
			Id:     bookID,
			Title:  fmt.Sprintf("title %d", bookID),
			Author: fmt.Sprintf("author %d", bookID),
			Price:  0.7,
		})
	}

	// bulk insert in many threads
	// TODO: add the ability to limit goroutines
	var wg sync.WaitGroup
	wg.Add(batchesCount)
	for i := 0; i < batchesCount; i++ {
		go func(i int) { bulkInsert(db, batches[i]); wg.Done() }(i)
	}
	wg.Wait()
}

// example from https://godoc.org/github.com/lib/pq#hdr-Bulk_imports
func bulkInsert(db *sql.DB, books []Book) {
	txn, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("books", "id", "title", "author", "price"))
	if err != nil {
		log.Println(err)
	}

	for _, book := range books {
		_, err = stmt.Exec(book.Id, book.Title, book.Author, book.Price)
		if err != nil {
			log.Println(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Println(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Println(err)
	}
}
