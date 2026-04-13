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
./bookstore add --ISBN 6 --title test-book --author name --price 200 --image_url [http://name.com/test.png](http://name.com/test.png)
./bookstore update --ISBN 6 --title test-book-1 --author name --price 2001 --image_url [http://image-name.com/test.png1](http://image-name.com/test.png1)
./bookstore delete --ISBN 6
*/
func main() {
	/*
		get books --all or --ISBN
		./bookstore get --all
		./bookstore get --ISBN 5
	*/
	getCli := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCli.Bool("all", false, "List all the books")
	getISBN := getCli.String("ISBN", "", "Get book by ISBN")

	/*
		add a book with ISBN ,title, author, price, image_url
		./bookstore add --ISBN 6 --title test-book --author name --price 200 --image_url [http://name.com/test.png](http://name.com/test.png)
	*/

	addCli := flag.NewFlagSet("add", flag.ExitOnError)
	addISBN := addCli.String("ISBN", "", "Book ISBN")
	addTitle := addCli.String("title", "", "Book title")
	addAuthor := addCli.String("author", "", "Book author")
	addPrice := addCli.Float64("price", 0, "Book price")
	addImageUrl := addCli.String("image_url", "", "Book image URL")

	/*
		update a book with ISBN ,title, author, price, image_url
		./bookstore update --ISBN 6 --title test-book-1 --author name --price 2001 --image_url [http://name.com/test.png1](http://name.com/test.png1)
	*/

	updateCli := flag.NewFlagSet("update", flag.ExitOnError)
	updateISBN := updateCli.String("ISBN", "", "Book ISBN")
	updateTitle := updateCli.String("title", "", "Book title")
	updateAuthor := updateCli.String("author", "", "Book author")
	updatePrice := updateCli.Float64("price", 0, "Book price")
	updateImageUrl := updateCli.String("image_url", "", "Book image URL")

	/*
		delete a book by --ISBN
		./bookstore delete --ISBN 6
	*/
	deleteCli := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteISBN := deleteCli.String("ISBN", "", "Delete book by ISBN")

	// Parse command line flags
	getCli.Parse(os.Args[2:])
	addCli.Parse(os.Args[2:])
	updateCli.Parse(os.Args[2:])
	deleteCli.Parse(os.Args[2:])

	// Handle commands
	switch os.Args[1] {
	case "get":
		if *getAll {
			handleGetBooks(getCli, getAll, nil)
		} else if *getISBN != "" {
			handleGetBooks(getCli, nil, getISBN)
		} else {
			getCli.PrintDefaults()
			os.Exit(1)
		}
	case "add":
		if *addISBN == "" || *addTitle == "" || *addAuthor == "" || *addPrice == 0 || *addImageUrl == "" {
			addCli.PrintDefaults()
			os.Exit(1)
		}
		handleAddBook(addCli, addISBN, addTitle, addAuthor, addPrice, addImageUrl, nil, nil, nil, nil, nil, nil)
	case "update":
		if *updateISBN == "" {
			updateCli.PrintDefaults()
			os.Exit(1)
		}
		handleUpdateBook(updateCli, updateISBN, updateTitle, updateAuthor, updatePrice, updateImageUrl, nil, nil, nil, nil, nil, nil)
	case "delete":
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
