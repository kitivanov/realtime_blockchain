package collector

import (
	"encoding/json"
	"log"
)

func HandleEvents(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println("marshal error:", err)
		return
	}

	log.Println(string(b))
}
