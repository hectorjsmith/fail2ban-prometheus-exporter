package textfile

import (
	"log"
	"net/http"
)

func (c *Collector) WriteTextFileMetrics(w http.ResponseWriter, r *http.Request) {
	if !c.enabled {
		return
	}

	for _, f := range c.fileMap {
		_, err := w.Write(f.fileContents)
		if err != nil {
			c.appendErrorForPath(f.fileName)
			log.Printf("error writing file contents to response writer '%s': %v", f.fileName, err)
		}
	}
}
