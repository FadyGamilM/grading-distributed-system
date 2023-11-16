package logger

import (
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
)

var log *stdlog.Logger

// this tpye will implement the io.Reader interface by implementing the Write([]byte)(int, error) method
type logFile string

// we will persist our logs into a file
// so the file in golang has a method called write which receives a slice of bytes and returns an int and error, so i will follow the same convention of file.write method
func (lfp logFile) Write(data []byte) (int, error) {
	// first we open the file/create it if it doesn't exist
	file, err := os.OpenFile(
		string(lfp),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)
	// handle errors
	if err != nil {
		return 0, err
	}

	// at the end clean the resources up
	defer file.Close()

	// write to the file and return the result
	return file.Write(data)
}

// initialize the log variable
func Run(destination string) {
	log = stdlog.New(logFile(destination), "", stdlog.LstdFlags)
}

// register the routers
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		msg, err := ioutil.ReadAll(r.Body)
		// if there is any error or the content of the log message is empty , we should return bad request response
		if err != nil || len(msg) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// we process the request by logging it to our log-file (current database)
		Log(string(msg))
	})
}

func Log(msg string) {
	log.Printf(" ➜ %v\n", msg)
}
