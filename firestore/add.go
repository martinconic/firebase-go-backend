package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func readJSONmodel(name string) []modelJson {
	jsonFile, err := os.Open(name)
	if err != nil {
		log.Fatalf("Connot open '%s': %s\n ", name, err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("Cannot read JSON data: ", err.Error())
	}

	var modeljson []modelJson
	json.Unmarshal(byteValue, &modeljson)

	fmt.Println(modeljson)
	return modeljson
}

func saveJsonToFirestore(modelfirestore []ModelFirestore, client *firestore.Client) {
	batch := client.Batch()

	for i := range modelfirestore {

		orderRef := client.Collection("modelPerson").NewDoc()
		batch.Set(orderRef, modelfirestore[i])
	}

	_, err := batch.Commit(context.Background())
	if err != nil {
		log.Fatalln("Error commiting batch: ", err.Error())
	}
}

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
		fmt.Println("Missing paramete(s)!")
		return
	}

	rows := readJSONmodel(os.Args[1])

	var modelsfirestore []ModelFirestore

	for i := range rows {
		modelfirestore := ModelFirestore{
			rows[i].ID,
			rows[i].FirstName,
			rows[i].LastName,
			rows[i].Description,
		}

		modelsfirestore = append(modelsfirestore, modelfirestore)

		if i == 450 || i == (len(rows)-1) {
			saveJsonToFirestore(modelsfirestore, client)
			modelsfirestore = []ModelFirestore{}
		}
	}
}
