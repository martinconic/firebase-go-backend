package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
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
		fmt.Println("Missing parameter, provide one!")
		return
	}

	var strs []string
	for i := range os.Args {
		if i > 0 {
			strs = append(strs, os.Args[i])
		}
	}

	query := client.Collection("modelPerson").Where("FirstName", "==", strings.Join(strs, " "))
	iter := query.Documents(context.Background())

	var modelfirestore ModelFirestore

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		doc.DataTo(&modelfirestore)

		names := strings.SplitN(modelfirestore.FirstName, " ", 2)

		_, err = doc.Ref.Update(context.Background(), []firestore.Update{{Path: "FirstName", Value: names[0]},
			{Path: "LastName", Value: names[1]}})

		if err != nil {
			log.Fatalln(err)
		}
	}
}
