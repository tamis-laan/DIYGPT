package wiki

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

func readBytes(reader *bufio.Reader, delimiter string, remove bool) ([]byte, error) {
	// Create a buffer
	var buf bytes.Buffer

	delBytes := []byte(delimiter)
	delLen := len(delBytes)
	lastBytes := make([]byte, delLen) // Keep track of the last few bytes for delimiter comparison

	for {

		// Read ReadBytes
		b, err := reader.ReadByte()

		// Block on EOF unill new messages are available
		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Handle error
		if err != nil && err != io.EOF {
			return buf.Bytes(), err
		}

		// Write the byte to the buffer
		buf.WriteByte(b)
		// Update the lastBytes queue
		copy(lastBytes, append(lastBytes[1:], b))

		// Check if the last few bytes match our delimiter
		if bytes.Equal(lastBytes, delBytes) {
			// If remove is true, remove the delimiter characters from the buffer before returning
			if remove {
				buf.Truncate(buf.Len() - delLen)
			}
			return buf.Bytes(), nil
		}
	}
}

type WikiClient struct {
	response *http.Response
	reader   *bufio.Reader
}

func NewWikiClient(url string) (*WikiClient, error) {
	// Create a new HTTP client.
	client := http.Client{}

	// Send GET request to the SSE endpoint.
	req, err := http.NewRequest("GET", url, nil)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Set header for SSE
	req.Header.Set("Accept", "text/event-stream")

	// Make the request.
	response, err := client.Do(req)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Create buffered io reader
	reader := bufio.NewReader(response.Body)

	// Read until double line break
	message, err := readBytes(reader, "\n\n", true)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Handle error
	if string(message) != ":ok" {
		return nil, err
	}

	// Create
	wikiClient := WikiClient{response, reader}

	// No errors
	return &wikiClient, err
}

func (connection *WikiClient) Close() {
	connection.response.Body.Close()
}

func (connection *WikiClient) Read() ([]string, error) {

	// Read until double line break
	raw, err := readBytes(connection.reader, "\n\n", true)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Cast Cast raw bytes to string
	txt := string(raw)

	// Split the text into it's message components
	message := strings.Split(string(txt), "\n")

	// Skip anything that is not really a message
	if message[0] != "event: message" {
		return nil, nil
	}

	// Remove id: from string
	message[1] = strings.TrimPrefix(message[1], "id:")

	// Remove data: from string
	message[2] = strings.TrimPrefix(message[2], "data:")

	// Return message
	return message[1:], nil
}
