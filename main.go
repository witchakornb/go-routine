package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// โครงสร้างสำหรับเก็บผลลัพธ์จาก API แต่ละตัว
// อาจจะเก็บข้อมูลที่ parse แล้ว หรือ เก็บ error ที่เกิดขึ้น
type APIResult struct {
	URL     string
	Body    []byte
	Error   error
	Latency time.Duration // เก็บเวลาที่ใช้ในการดึงข้อมูล (optional)
}

// ฟังก์ชันสำหรับดึงข้อมูลจาก API เดียว
// รับ URL, WaitGroup สำหรับจัดการ goroutine, และ channel สำหรับส่งผลลัพธ์กลับ
func fetchAPI(url string, wg *sync.WaitGroup, resultsChan chan<- APIResult) {
	// defer wg.Done() จะถูกเรียกเมื่อฟังก์ชันนี้ทำงานเสร็จสิ้น
	// เพื่อบอก WaitGroup ว่า goroutine นี้ทำงานเสร็จแล้ว
	defer wg.Done()

	start := time.Now() // เริ่มจับเวลา

	// สร้าง HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		resultsChan <- APIResult{URL: url, Error: fmt.Errorf("error creating request: %w", err), Latency: time.Since(start)}
		return
	}

	// ส่ง request
	client := &http.Client{Timeout: 10 * time.Second} // ตั้ง timeout ป้องกันการรอคอยนานเกินไป
	resp, err := client.Do(req)
	if err != nil {
		resultsChan <- APIResult{URL: url, Error: fmt.Errorf("error sending request: %w", err), Latency: time.Since(start)}
		return
	}
	// defer resp.Body.Close() สำคัญมาก เพื่อคืนทรัพยากรเมื่อสิ้นสุดการทำงาน
	defer resp.Body.Close()

	// ตรวจสอบ Status Code
	if resp.StatusCode != http.StatusOK {
		resultsChan <- APIResult{URL: url, Error: fmt.Errorf("unexpected status code: %d", resp.StatusCode), Latency: time.Since(start)}
		return
	}

	// อ่านข้อมูลจาก response body
	body, err := io.ReadAll(resp.Body)
	latency := time.Since(start) // หยุดจับเวลา
	if err != nil {
		resultsChan <- APIResult{URL: url, Error: fmt.Errorf("error reading response body: %w", err), Latency: latency}
		return
	}

	// ส่งผลลัพธ์ (ข้อมูลที่ได้) กลับไปที่ channel
	resultsChan <- APIResult{URL: url, Body: body, Latency: latency}
}

func main() {
	// --- กำหนดค่าเริ่มต้น ---
	// URL ของ API ที่ต้องการดึง (ใช้ API ตัวอย่าง)
	apiURL1 := "https://httpbin.org/get?source=api1"  // API ตัวอย่างที่ trả về JSON เกี่ยวกับ request ที่ส่งไป
	apiURL2 := "https://httpbin.org/delay/1" // API ตัวอย่างที่จะหน่วงเวลา 1 วินาทีก่อนตอบกลับ

	// สร้าง WaitGroup เพื่อรอให้ goroutine ทั้งหมดทำงานเสร็จ
	var wg sync.WaitGroup

	// สร้าง Channel เพื่อรับผลลัพธ์จาก goroutine ต่างๆ
	// กำหนด buffer size เท่ากับจำนวน goroutine ที่จะสร้าง เพื่อไม่ให้ goroutine บล็อกตอนส่งข้อมูล
	resultsChan := make(chan APIResult, 2)

	// --- เริ่มการทำงานพร้อมกัน ---
	fmt.Println("เริ่มต้นดึงข้อมูลจาก API พร้อมกัน...")

	// เพิ่ม counter ใน WaitGroup เท่ากับจำนวน goroutine ที่จะรัน
	wg.Add(2)

	// รัน goroutine ที่ 1 เพื่อดึงข้อมูลจาก apiURL1
	go fetchAPI(apiURL1, &wg, resultsChan)

	// รัน goroutine ที่ 2 เพื่อดึงข้อมูลจาก apiURL2
	go fetchAPI(apiURL2, &wg, resultsChan)

	// --- รอและปิด Channel ---
	// สร้าง goroutine แยกต่างหากเพื่อรอให้ wg.Wait() เสร็จสิ้น แล้วจึงปิด Channel
	// ทำแบบนี้เพื่อป้องกัน deadlock กรณีที่ main goroutine รออ่านจาก channel ที่ไม่มีใครส่งมาแล้ว
	go func() {
		wg.Wait()      // รอจนกว่า counter ของ WaitGroup จะเป็น 0 (goroutine ทั้ง 2 ตัวเรียก Done())
		close(resultsChan) // ปิด Channel หลังจาก goroutine ทั้งหมดทำงานเสร็จ
		fmt.Println("Goroutines ทั้งหมดทำงานเสร็จสิ้น, ปิด channel.")
	}()

	// --- ประมวลผลผลลัพธ์ ---
	fmt.Println("รอรับผลลัพธ์จาก API...")

	// วนลูปเพื่อรับผลลัพธ์จาก Channel จนกว่า Channel จะถูกปิด
	for result := range resultsChan {
		fmt.Printf("\nได้รับผลลัพธ์จาก: %s (ใช้เวลา: %v)\n", result.URL, result.Latency)
		if result.Error != nil {
			// ถ้ามี error เกิดขึ้น
			fmt.Printf("เกิดข้อผิดพลาด: %v\n", result.Error)
		} else {
			// ถ้าสำเร็จ พิมพ์ข้อมูลที่ได้ (ตัวอย่างนี้พิมพ์แค่ความยาว)
			// ในการใช้งานจริง อาจจะทำการ unmarshal JSON หรือประมวลผลอื่นๆ
			fmt.Printf("ข้อมูลที่ได้รับ (ขนาด %d bytes): %s\n", len(result.Body), string(result.Body))
			// หมายเหตุ: การแปลง []byte เป็น string โดยตรงอาจจะไม่เหมาะกับข้อมูลขนาดใหญ่มาก
		}
	}

	fmt.Println("\nประมวลผลผลลัพธ์ทั้งหมดเรียบร้อย")
}