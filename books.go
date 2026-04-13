package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Book struct {
	ISBN          string   `json:"isbn"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Price         float64  `json:"price"`
	ImageURL      string   `json:"image_url"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	PublishedDate string   `json:"published_date"`
	Stock         int      `json:"stock"`
	Reviews       []string `json:"reviews"`
	Rating        float64  `json:"rating"`
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func getBooks() ([]Book, error) {
	booksSlice, err := os.ReadFile("./books.json")

	if err != nil {
		return nil, err
	}

	var books []Book
	err = json.Unmarshal(booksSlice, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func saveBooks(books []Book) error {
	booksSlice, err := json.Marshal(books)
	checkError(err)
	err = os.WriteFile("./books.json", booksSlice, 0644)
	return err
}

func handleGetBooks(cli *flag.FlagSet, all *bool, isbn *string) {
	cli.Parse(os.Args[2:])

	if !*all && *isbn == "" {
		fmt.Println("sub command --all or --ISBN needed")
		os.Exit(1)
	}
	books, err := getBooks()
	checkError(err)

	if *all {
		fmt.Printf("%-15s %-30s %-20s %-8s %-10s\n", "ISBN", "Title", "Author", "Price", "Rating")
		for _, book := range books {
			fmt.Printf("%-15s %-30s %-20s %-8.2f %-10.1f\n", book.ISBN, book.Title, book.Author, book.Price, book.Rating)
		}
	}

	if *isbn != "" {
		var foundBook bool
		fmt.Printf("ISBN \t title \t Author \t Price \t IMAGE URL \n")

		for _, book := range books {

			if *isbn == book.ISBN {
				foundBook = true
				fmt.Printf("%s \t %s \t %s \t %f \t %s \n", book.ISBN, book.Title, book.Author, book.Price, book.ImageURL)
			}
		}

		if !foundBook {
			fmt.Println("Book not found")
			os.Exit(1)
		}
	}
}

func handleAddBook(addCli *flag.FlagSet, isbn *string, title *string, author *string, price *float64, imageURL *string, description *string, category *string, publishedDate *string, stock *int, reviews *string, rating *float64) {
	addCli.Parse(os.Args[2:])

	if *isbn == "" || *title == "" || *author == "" || *price == 0 || *imageURL == "" || *description == "" || *category == "" || *publishedDate == "" || *stock == 0 || *reviews == "" || *rating == 0 {
		fmt.Println("Please Provide all required fileds")
		addCli.PrintDefaults()
		os.Exit(1)
	}

	books, err := getBooks()
	checkError(err)
	var foundBook bool

	for i, book := range books {
		if book.ISBN == *isbn {
			foundBook = true
			books[i].Title = *title
			books[i].Author = *author
			books[i].Price = *price
			books[i].ImageURL = *imageURL
			books[i].Description = *description
			books[i].Category = *category
			books[i].PublishedDate = *publishedDate
			books[i].Stock = *stock
			books[i].Reviews = strings.Split(*reviews, ",")
			books[i].Rating = *rating
		}
	}

	if !foundBook {
		newBook := Book{
			ISBN:          *isbn,
			Title:         *title,
			Author:        *author,
			Price:         *price,
			ImageURL:      *imageURL,
			Description:   *description,
			Category:      *category,
			PublishedDate: *publishedDate,
			Stock:         *stock,
			Reviews:       strings.Split(*reviews, ","),
			Rating:        *rating,
		}
		books = append(books, newBook)
	}
	err = saveBooks(books)
	checkError(err)
	fmt.Println("Book added Successfully")
}

func handleDeleteBook(deleteCli *flag.FlagSet, isbn *string) {

	deleteCli.Parse(os.Args[2:])

	if *isbn == "" {
		fmt.Println("Please Provide ISBN")
		deleteCli.PrintDefaults()
		os.Exit(1)
	}

	books, err := getBooks()
	var foundBook bool

	for i, book := range books {
		if book.ISBN == *isbn {
			books = append(books[:i], books[i+1:]...)
			foundBook = true
		}
	}

	if !foundBook {
		fmt.Println("Book not found")
		os.Exit(1)
	}

	err = saveBooks(books)
	checkError(err)
	fmt.Println("Book deleted successfully")

}

// handle Update Book
func handleUpdateBook(updateCli *flag.FlagSet, isbn *string, title *string, author *string, price *float64, imageURL *string, description *string, category *string, publishedDate *string, stock *int, reviews *string, rating *float64) {
	updateCli.Parse(os.Args[2:])
	if *isbn == "" {
		updateCli.PrintDefaults()
		os.Exit(1)
	}
	books, err := getBooks()
	var foundBook bool
	for i, book := range books {
		if book.ISBN == *isbn {
			foundBook = true
			books[i].Title = *title
			books[i].Author = *author
			books[i].Price = *price
			books[i].ImageURL = *imageURL
			books[i].Description = *description
			books[i].Category = *category
			books[i].PublishedDate = *publishedDate
			books[i].Stock = *stock
			books[i].Reviews = strings.Split(*reviews, ",")
			books[i].Rating = *rating
		}
	}
	if !foundBook {
		fmt.Println("Book not found")
		os.Exit(1)
	}
	err = saveBooks(books)
	checkError(err)
	fmt.Println("Book updated successfully")
}
