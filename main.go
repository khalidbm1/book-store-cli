package main

import (
	"flag"
	"fmt"
	"os"
)

/*
Main function to handle get, add, update, delete commands.
Expected commands are as follows:
./bookstore get --all
./bookstore get --ISBN 5
./bookstore add --ISBN 6 --title test-book --author name --price 200 --image_url http://name.com/test.png --description "desc" --category fiction --published_date 2020-01-01 --stock 10 --reviews "good,great" --rating 4.5
./bookstore update --ISBN 6 --title test-book-1 --author name --price 2001 --image_url http://image-name.com/test.png1
./bookstore delete --ISBN 6
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand: get, add, update, delete")
		os.Exit(1)
	}

	getCli := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCli.Bool("all", false, "List all the books")
	getISBN := getCli.String("ISBN", "", "Get book by ISBN")

	addCli := flag.NewFlagSet("add", flag.ExitOnError)
	addISBN := addCli.String("ISBN", "", "Book ISBN")
	addTitle := addCli.String("title", "", "Book title")
	addAuthor := addCli.String("author", "", "Book author")
	addPrice := addCli.Float64("price", 0, "Book price")
	addImageUrl := addCli.String("image_url", "", "Book image URL")
	addDescription := addCli.String("description", "", "Book description")
	addCategory := addCli.String("category", "", "Book category")
	addPublishedDate := addCli.String("published_date", "", "Book published date")
	addStock := addCli.Int("stock", 0, "Book stock")
	addReviews := addCli.String("reviews", "", "Book reviews (comma-separated)")
	addRating := addCli.Float64("rating", 0, "Book rating")

	updateCli := flag.NewFlagSet("update", flag.ExitOnError)
	updateISBN := updateCli.String("ISBN", "", "Book ISBN")
	updateTitle := updateCli.String("title", "", "Book title")
	updateAuthor := updateCli.String("author", "", "Book author")
	updatePrice := updateCli.Float64("price", 0, "Book price")
	updateImageUrl := updateCli.String("image_url", "", "Book image URL")
	updateDescription := updateCli.String("description", "", "Book description")
	updateCategory := updateCli.String("category", "", "Book category")
	updatePublishedDate := updateCli.String("published_date", "", "Book published date")
	updateStock := updateCli.Int("stock", 0, "Book stock")
	updateReviews := updateCli.String("reviews", "", "Book reviews (comma-separated)")
	updateRating := updateCli.Float64("rating", 0, "Book rating")

	deleteCli := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteISBN := deleteCli.String("ISBN", "", "Delete book by ISBN")

	switch os.Args[1] {
	case "get":
		getCli.Parse(os.Args[2:])
		if !*getAll && *getISBN == "" {
			getCli.PrintDefaults()
			os.Exit(1)
		}
		handleGetBooks(getCli, getAll, getISBN)
	case "add":
		addCli.Parse(os.Args[2:])
		handleAddBook(addCli, addISBN, addTitle, addAuthor, addPrice, addImageUrl, addDescription, addCategory, addPublishedDate, addStock, addReviews, addRating)
	case "update":
		updateCli.Parse(os.Args[2:])
		if *updateISBN == "" {
			updateCli.PrintDefaults()
			os.Exit(1)
		}
		handleUpdateBook(updateCli, updateISBN, updateTitle, updateAuthor, updatePrice, updateImageUrl, updateDescription, updateCategory, updatePublishedDate, updateStock, updateReviews, updateRating)
	case "delete":
		deleteCli.Parse(os.Args[2:])
		if *deleteISBN == "" {
			deleteCli.PrintDefaults()
			os.Exit(1)
		}
		handleDeleteBook(deleteCli, deleteISBN)
	default:
		fmt.Println("Invalid command")
		os.Exit(1)
	}
}
