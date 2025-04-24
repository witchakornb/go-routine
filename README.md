# Go Routine Example

This project demonstrates how to use Go routines and channels to fetch data from multiple APIs concurrently. It includes an example of using `sync.WaitGroup` to manage goroutines and a buffered channel to collect results.

## Features
- Fetch data from multiple APIs concurrently.
- Handle errors gracefully for each API call.
- Measure the latency of each API request.
- Use `sync.WaitGroup` to synchronize goroutines.
- Use a buffered channel to collect results without blocking.

## How It Works

1. **API Fetching**:
   - The `fetchAPI` function is responsible for making HTTP GET requests to a given URL.
   - It measures the time taken for the request and handles errors such as request creation, response status, and reading the response body.

2. **Concurrency**:
   - Two API URLs are fetched concurrently using goroutines.
   - A `sync.WaitGroup` is used to wait for all goroutines to complete.

3. **Channel for Results**:
   - A buffered channel is used to collect results from each goroutine.
   - The channel is closed once all goroutines finish their work.

4. **Result Processing**:
   - Results are read from the channel and processed to display the URL, latency, and any errors or data received.

## Code Overview

### `fetchAPI` Function
This function takes a URL, a `WaitGroup`, and a results channel as arguments. It performs the following steps:
- Creates an HTTP GET request.
- Sends the request and measures the latency.
- Handles errors and unexpected status codes.
- Reads the response body and sends the result to the channel.

### `main` Function
The `main` function demonstrates the following:
- Initializes two API URLs.
- Creates a `WaitGroup` and a buffered channel.
- Starts two goroutines to fetch data from the APIs.
- Waits for all goroutines to complete and closes the channel.
- Processes the results from the channel.

## How to Run

1. **Prerequisites**:
   - Install [Go](https://golang.org/dl/).

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/witchakornb/go-routine.git
   cd go-routine
   ```

3. **Run the Program**:
   ```bash
   go run main.go
   ```

4. **Expected Output**:
   - The program will fetch data from two APIs concurrently and display the results, including latency and any errors.

## Example Output
```
go run main.go
เริ่มต้นดึงข้อมูลจาก API พร้อมกัน...
รอรับผลลัพธ์จาก API...

ได้รับผลลัพธ์จาก: https://httpbin.org/get?source=api1 (ใช้เวลา: 7.2033787s)
ข้อมูลที่ได้รับ (ขนาด 309 bytes): {
  "args": {
    "source": "api1"
  }, 
  "headers": {
    "Accept-Encoding": "gzip",
    "Host": "httpbin.org",
    "User-Agent": "Go-http-client/2.0",
    "X-Amzn-Trace-Id": "Root=1-680a7788-042afb6847ea76ec43e8dd8d"
  },
  "origin": "202.28.118.119",
  "url": "https://httpbin.org/get?source=api1"
}

Goroutines ทั้งหมดทำงานเสร็จสิ้น, ปิด channel.

ได้รับผลลัพธ์จาก: https://httpbin.org/delay/1 (ใช้เวลา: 10.0007335s)
เกิดข้อผิดพลาด: error sending request: Get "https://httpbin.org/delay/1": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

ประมวลผลผลลัพธ์ทั้งหมดเรียบร้อย
```

## Author
This project was created by Witchakorn Boonprakom. You can find more projects and contributions on [GitHub](https://github.com/witchakornb).

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.