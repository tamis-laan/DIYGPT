package main

import (
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Read configuration
func init() {
	// Define command-line flags
	pflag.String("kafka.bootstrap", "localhost:9092", "Kafka bootstrap servers")
	pflag.String("kafka.groupid", "wiki consumer", "The kafka client id")
	pflag.String("topic", "wiki.test", "The topic to publish to")

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

	// Err check
	if err != nil {
		log.Fatal(err)
	}

	// Close consumer
	defer c.Close()

	// Log to console
	log.Printf("%+v\n", c)

	// Define the topic
	topic := viper.GetString("topic")

	// Subscribe to test topic
	err = c.Subscribe(topic, nil)

	// Err check
	if err != nil {
		log.Fatal(err)
	}

	// Go into consumption loop
	for {
		// Wait for a message
		msg, err := c.ReadMessage(-1)

		// Err check
		if err != nil {
			log.Fatal(err)
		}

		// Process message
		log.Println(msg.TopicPartition, msg.Key, string(msg.Value))

		// Commit message offset
		_, err = c.CommitMessage(msg)

		// Err check
		if err != nil {
			log.Fatal(err)
		}
	}

}
