package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"io"
	"log"
	"math"
	"os"
	"time"
)

func processFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	streamSize := math.Round(float64((fileStat.Size() / 3.0)) + 0.5)

	fmt.Printf("file size: %d\n", fileStat.Size())

	buffer := make([]byte, int(streamSize))
	uuid := guuid.NewString()

	for i := 0; i < 3; i++ {
		_, err := file.Read(buffer)

		if err == io.EOF {
			log.Fatalln(err)
		}

		fileChunk := newFileChunk(buffer, i, uuid, "")

		// request
		go sendBytes(fileChunk, "file.png")
	}
}

type fileChunk struct {
	uuid   string
	data   []byte
	index  int
	meta   string
	asJson []byte
}

func newFileChunk(data []byte, index int, uuid string, meta string) *fileChunk {
	chunk := fileChunk{
		uuid:  uuid,
		data:  data,
		index: index,
		meta:  meta,
	}

	fileChunkMap := map[string]string{
		"data":  string(chunk.data),
		"index": string(index),
		"meta":  meta,
	}

	jsonData, err := json.Marshal(fileChunkMap)
	if err != nil {
		log.Fatalln(err)
	}

	chunk.asJson = jsonData

	return &chunk
}

func sendBytes(fileChunk *fileChunk, fileName string) {
	fmt.Printf("starting to send %d chunk \n", fileChunk.index)

	url := fmt.Sprintf("localhost:3001/%s", fileChunk.uuid)

	//START response mocking
	resp, err := postToMock(url, bytes.NewBuffer(fileChunk.asJson), fileChunk)
	if err != nil {
		log.Fatalln(err)
	}

	type ResponseMock struct {
		Data   string
		Status string
	}

	var respMock ResponseMock

	err = json.Unmarshal(resp, &respMock)

	fmt.Printf("%s\n", respMock.Data)
	// END response mocking
}

func postToMock(url string, buffer *bytes.Buffer, chunk *fileChunk) (resp []byte, err error) {
	time.Sleep(3 * time.Second)

	data := fmt.Sprintf("completed for part: %d, for uuid: %s", chunk.index, chunk.uuid)

	rawResponse := map[string]string{
		"status": "200",
		"data":   string(data),
	}

	jsonResponse, err := json.Marshal(rawResponse)
	if err != nil {
		log.Fatalln(err)
	}

	response := bytes.NewBuffer(jsonResponse).Bytes()

	return response, nil
}
