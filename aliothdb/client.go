package aliothdb

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2_api "github.com/influxdata/influxdb-client-go/v2/api"
)

func GetClient(uri string, token string) (client influxdb2.Client) {
	fmt.Println("Acquiring InfluxDB client...")
	client = influxdb2.NewClient("http://localhost:8086", token)
	return
}

func GetWriteAPI(client influxdb2.Client, organization string, bucket string) (writeAPI influxdb2_api.WriteAPIBlocking) {
	writeAPI = client.WriteAPIBlocking(organization, bucket)
	return
}

func WriteTemperature(writeAPI influxdb2_api.WriteAPIBlocking, avg float64, max float64) {
	// Create point using full params constructor
	const table = "kevin_testd"
	currentTime := time.Now()

	p := influxdb2.NewPoint(table,
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": avg, "max": max},
		currentTime)
	err := writeAPI.WritePoint(context.Background(), p)
	if err == nil {
		fmt.Printf("Writing avg=%f, max=%f, time=%s to table %s\n", avg, max, currentTime.String(), table)
	} else {
		fmt.Printf("Write failed:\n\t%s\n", err.Error())
	}

	/*

		// Create point using fluent style
		p = influxdb2.NewPointWithMeasurement("kevin_testd").
			AddTag("unit", "temperature").
			AddField("avg", avg).
			AddField("max", max).
			SetTime(time.Date(2021, time.August, 2, 12, 0, 0, 0, time.UTC))
		writeAPI.WritePoint(context.Background(), p)

		// Or write directly line protocol
		line := fmt.Sprintf("kevin_testd,unit=temperature avg=%f,max=%f", avg, max)
		writeAPI.WriteRecord(context.Background(), line)

	*/

}
