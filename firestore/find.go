package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	// Use a service account
	sa := option.WithCredentialsFile("../firebase.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}

	var strs []string

	for i := range os.Args {
		if i > 0 {
			strs = append(strs, os.Args[i])
		}
	}

	iter := client.Collection("modelPerson").Where("FirstName", "==", strings.Join(strs, " ")).Documents(context.Background())

	i := 0
	var modelFirestore ModelFirestore

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		doc.DataTo(&modelFirestore)
		fmt.Println(modelFirestore)
		i++
	}

	fmt.Println(i)
}
