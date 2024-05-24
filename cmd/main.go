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

// AudioData represents the structure of the audio data
type AudioData struct {
	ID          int       `json:"id"`
	FileName    string    `json:"file_name"`
	SR          int       `json:"sr"`
	ArrayLength int       `json:"arrayLength"`
	AudioArray  []float64 `json:"audio_array"`
	IsVoice     []bool    `json:"is_voice"`
	IsViolin    []bool    `json:"is_violin"`
	IsMridangam []bool    `json:"is_mridangam"`
	IsGhatam    []bool    `json:"is_ghatam"`
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

	response := fmt.Sprintf(`
		<div id="output-images" class="image-container">
			<a href="./output.html">
				<img id="output-image" src="./output.png?t=%d" alt="Intervals">
			</a>
			<img id="sound-wave" src="./sound_wave.png?t=%d" alt="sound_wave">
		</div>`, time.Now().Unix(), time.Now().Unix())

	fmt.Fprint(w, response)
}

// generateRandomBoolArray generates an array of random boolean values
func generateRandomBoolArray(length int) []bool {
	rand.NewSource(time.Now().UnixNano())
	boolArray := make([]bool, length)
	for i := range boolArray {
		boolArray[i] = rand.Intn(2) == 1
	}
	return boolArray
}

// generateRandomFloatArray generates an array of random float64 values
func generateRandomFloatArray(length int) []float64 {
	rand.NewSource(time.Now().UnixNano())
	floatArray := make([]float64, length)
	for i := range floatArray {
		floatArray[i] = rand.Float64()
	}
	return floatArray
}

// GenerateAndWriteJSON generates the audio data with random values and writes it to a JSON file
func GenerateAndWriteJSON(fileName string) error {
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

	// Define the length of the arrays
	arrayLength := 1000

	// Create an instance of AudioData with random values
	audioData := AudioData{
		ID:          0,
		FileName:    "performance1mixed.wav",
		SR:          44100,
		ArrayLength: arrayLength,
		AudioArray:  generateRandomFloatArray(arrayLength),
		IsVoice:     generateRandomBoolArray(arrayLength),
		IsViolin:    generateRandomBoolArray(arrayLength),
		IsMridangam: generateRandomBoolArray(arrayLength),
		IsGhatam:    generateRandomBoolArray(arrayLength),
	}

	// Convert the AudioData instance to JSON
	jsonData, err := json.MarshalIndent(audioData, "", "  ")
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
