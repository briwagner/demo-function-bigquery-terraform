package bigquery_etl

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
)

// GCSEvent provides access to the bucket event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// Global API client
var storageClient *storage.Client

// BigQuery identifiers.
var projectID string
var dataSet string
var tableID string

func init() {
	// Declare a separate err variable to avoid shadowing the client variables.
	var err error

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
}

// LoadData is the function entrypoint.
func LoadData(ctx context.Context, e GCSEvent) error {
	// Set environment variables.
	projectID = os.Getenv("PROJECT_ID")
	dataSet = os.Getenv("DATASET")
	tableID = os.Getenv("TABLE_ID")

	// Build uri for new file uploaded.
	uri := fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name)
	log.Printf("Received a new file %s", uri)

	// Read file extension to process accordingly.
	fileType := getFiletype(e.Name)
	if fileType == "" {
		return fmt.Errorf("Cannot determine filetype for: %s", e.Name)
	}
	log.Printf("Detected file type: %s", fileType)

	// Process data transfer to BigQuery.
	gcsRef := bigquery.NewGCSReference(uri)
	// Allow for skipping columns/fields that are not exported.
	gcsRef.IgnoreUnknownValues = true

	// Tell BigQuery how to process data, depending on type.
	switch fileType {
	case ".csv":
		gcsRef.SourceFormat = bigquery.CSV
		// Exclude title row from import.
		gcsRef.SkipLeadingRows = 1

	case ".json":
		gcsRef.SourceFormat = bigquery.JSON
	}

	// Use struct to define data schema.
	schema, err := bigquery.InferSchema(Person{})
	if err != nil {
		return err
	}
	gcsRef.Schema = schema

	// Initiate data loader.
	client, err := bigquery.NewClient(ctx, projectID)
	loader := client.Dataset(dataSet).Table(tableID).LoaderFrom(gcsRef)
	// Overwrite data in table.
	loader.WriteDisposition = bigquery.WriteTruncate

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		// Errors from processing the file are captured in status variable.
		return fmt.Errorf("File transfer completed with error: %v", status.Err())
	}

	return nil
}

func getFiletype(filename string) string {
	r, err := regexp.Compile(`\.[a-z]+$`)
	if err != nil {
		log.Printf("Failed to compile regex: %e", err)
		return ""
	}
	return r.FindString(filename)
}
