package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"services/pkg/schema"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Payload struct {
	Sentence string `json:"sentence"`
}

// Read configuration
func init() {
	// Define command-line flags
	pflag.String("kafka.bootstrap", "localhost:9092", "Kafka bootstrap servers")
	pflag.String("kafka.groupid", "", "The kafka client id")
	pflag.String("topic.input", "", "The input topic to read from")
	pflag.String("topic.output", "", "The output topic to read from")
	pflag.String("model.host", "localhost", "Model host url")
	pflag.Int("model.port", 8000, "Model port number")

	// Parse command-line flags
	pflag.Parse()

	// Bind command-line flags to Viper
	viper.BindPFlags(pflag.CommandLine)

	// Replace . with _ in environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Make configurable from environment variables
	viper.AutomaticEnv()

}

func main() {

	// Create consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.bootstrap"),
		"group.id":          viper.GetString("kafka.groupid"),
		"auto.offset.reset": "smallest",
	})

	log.Println("bootstrap", viper.GetString("kafka.bootstrap"))
	log.Println("group id", viper.GetString("kafka.groupid"))
	log.Println("topic", viper.GetString("topic.input"))

	// Err check
	if err != nil {
		log.Fatal(err)
	}

	// Close consumer
	defer c.Close()

	// Log to console
	log.Printf("%+v\n", c)

	// Get input topic
	input_topic := viper.GetString("topic.input")

	// Get output topic
	// output_topic := viper.GetString("topic.output")

	// Subscribe to test topic
	err = c.Subscribe(input_topic, nil)

	// Err check
	if err != nil {
		log.Fatal(err)
	}

	// Construct model endpoint url
	model_enpoint := fmt.Sprintf("http://%s:%d/run", viper.GetString("model.host"), viper.GetInt("model.port"))

	// Create http client
	client := &http.Client{}

	// Go into consumption loop
	for {

		// Wait for a message
		msg, err := c.ReadMessage(-1)

		// Err check
		if err != nil {
			log.Fatal(err)
		}

		// Define content
		content := schema.WikiPageCreateLocal{}

		// Unpack msg into content
		err = json.Unmarshal(msg.Value, &content)

		// Error check
		if err != nil {
			log.Fatal(err)
		}

		// Create the payload
		payload := Payload{
			Sentence: content.Comment,
		}

		// Turn payload into json bytes
		payloadBytes, err := json.Marshal(payload)

		// Err check
		if err != nil {
			log.Fatal(err)
		}

		// Body reader
		body := bytes.NewReader(payloadBytes)

		// Create the request
		req, err := http.NewRequest("POST", model_enpoint, body)

		// Set context header to json type
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		resp, err := client.Do(req)

		// Err check
		if err != nil {
			log.Fatal(err)
		}

		// Read the response body
		responseBody, err := io.ReadAll(resp.Body)

		// Error check
		if err != nil {
			log.Fatal(err)
		}

		// Process message
		log.Println(string(payloadBytes), "\n", string(responseBody))

		// Commit message offset
		_, err = c.CommitMessage(msg)

		// Err check
		if err != nil {
			log.Fatal(err)

			// Close body
			resp.Body.Close()
		}
	}

}
