package io

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func OpenFile(file string) ([]byte, error) {
	data, err := os.Open(file)
	var buf bytes.Buffer
	tee := io.TeeReader(data, &buf)
	bytes, _ := ioutil.ReadAll(tee)
	data.Close()
	return bytes, err
}

func ReadAll(data []byte) *bytes.Buffer {
	//byteValue, _ := ioutil.ReadAll(reader)
	//if err := json.Unmarshal(byteValue, &someStruct); err != nil {
	//	fmt.Println("There was an error:", err)
	//}
	return bytes.NewBuffer(data)
}

//func ReadAll(payload *os.File) *bytes.Buffer {
//	byteValue, _ := ioutil.ReadAll(payload)
//	//if err := json.Unmarshal(byteValue, &someStruct); err != nil {
//	//	fmt.Println("There was an error:", err)
//	//}
//	return bytes.NewBuffer(byteValue)
//}
