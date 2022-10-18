package common

import "log"

func Check(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}

func GetOrDefault(index int, list []string) string {
	if len(list) >= index+1 {
		return list[index]
	}
	return ""
}
