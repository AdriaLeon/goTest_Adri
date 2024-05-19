package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle requests for generating JSON data
	http.HandleFunc("/generate-json", HandleGenerateJSONAndCallPythonScript)

	// Handle requests to the root route

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	// Handle upload of files requests
	http.HandleFunc("/upload", HandleUploadFile)
}

type Instrument struct {
	Name      string `json:"name"`
	Intervals []bool `json:"intervals"`
}

type CarnaticData struct {
	Instruments []Instrument `json:"instruments"`
}

func HandleGenerateJSONAndCallPythonScript(w http.ResponseWriter, r *http.Request) {
	// Generate JSON file
	fileName := "input.json"
	err := GenerateAndWriteJSON(fileName)
	if err != nil {
		http.Error(w, "Error generating JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Call Python script
	cmd := exec.Command("python3", "lib/script.py", fileName)
	err = cmd.Run() // Run the command without capturing output
	if err != nil {
		http.Error(w, "Error calling Python script: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally, you can send a success status code
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `<img id="output-image" src="./output.png?t=%d" alt="Intervals">`, time.Now().Unix())
}

func GenerateAndWriteJSON(fileName string) error {
	// Seed the random number generator
	rand.NewSource(time.Now().UnixNano())

	// Helper function to generate random boolean intervals
	generateRandomIntervals := func(length int) []bool {
		intervals := make([]bool, length)
		for i := range intervals {
			intervals[i] = rand.Intn(2) == 1
		}
		return intervals
	}

	// Define the directory for storing data files
	directory := "data"

	// Ensure the directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Define the file path for the JSON file
	filePath := filepath.Join(directory, fileName)

	// Create the output file in the specified directory
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Define the Carnatic instruments with random boolean intervals
	data := CarnaticData{
		Instruments: []Instrument{
			{
				Name:      "Veena",
				Intervals: generateRandomIntervals(10),
			},
			{
				Name:      "Mridangam",
				Intervals: generateRandomIntervals(10),
			},
			{
				Name:      "Flute",
				Intervals: generateRandomIntervals(10),
			},
			{
				Name:      "Violin",
				Intervals: generateRandomIntervals(10),
			},
			{
				Name:      "Kanjira",
				Intervals: generateRandomIntervals(10),
			},
			{
				Name:      "Ghatam",
				Intervals: generateRandomIntervals(10),
			},
		},
	}

	// Convert the data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func HandleGenerateJSON(w http.ResponseWriter, r *http.Request) {
	fileName := "input.json"
	err := GenerateAndWriteJSON(fileName)
	if err != nil {
		http.Error(w, "Error generating JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON file successfully created"))
}

func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Max size of 10MB
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extracts the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error extracting file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Saves file on test_uploads TODO: Tendría que crear un carpeta uploads para guardarlo?
	fileName := filepath.Join(".", "uploads", "test_uploads")
	outFile, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Error creating file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copies the file on the outFile
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error copying file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Sends the answer to the client
	fmt.Fprintf(w, "File uploaded successfully: %s", "test_uploads")
}
