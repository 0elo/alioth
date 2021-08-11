package main

import (
	"alioth/aliothdb"
	"context"
	"fmt"
	"math/rand"
	"os"

	"rsc.io/quote"
)

func main() {
	const organization = "kevin"
	const bucket = "kevin"
	influxDBToken := os.Getenv("INFLUXDB_TOKEN")
	// const token = "PMfIKn1sIt-A7el5Z9fF1qHe_j4nThoxC-i2KPkXZmk874taNba9FF4OjEtyHMxqXF1QwFi5f-jYcqndGpaXoA=="
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := aliothdb.GetClient("http://localhost:8086", influxDBToken)
	// Use blocking write client for writes to desired bucket
	writeAPI := aliothdb.GetWriteAPI(client, organization, bucket)
	avg := rand.Float64()
	max := avg

	for max < avg {
		max = rand.Float64()
	}
	aliothdb.WriteTemperature(writeAPI, avg, max)

	// Get query client
	queryAPI := client.QueryAPI("kevin")
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"kevin")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "kevin_testd")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	// Ensures background processes finishes
	defer client.Close()

	fmt.Println(quote.Go())
}
