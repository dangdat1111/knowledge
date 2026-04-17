# Chương 01 - Your First Go Application

## Đối tượng và mục tiêu tài liệu
- Đối tượng: người mới bắt đầu với Go, đã biết một ngôn ngữ lập trình bất kỳ.
- Mục tiêu: hiểu vòng đời của một chương trình Go đơn giản, từ tạo mã nguồn đến chạy và kiểm tra kết quả.

## Mục tiêu học tập
Sau khi hoàn thành chương này, người học có thể:
- [x] Tạo và chạy ứng dụng Go đầu tiên bằng `go run`.
- [x] Biên dịch chương trình bằng `go build` và chạy file thực thi.
- [x] Giải thích vai trò của package `main` và hàm `main()`.
- [x] Đọc được luồng thực thi cơ bản của một chương trình Go.

## Yêu cầu trước khi học
- Kiến thức nền tảng cần có:
  - [x] Biết khái niệm biến, hàm, và đầu vào/đầu ra cơ bản.
- Công cụ và môi trường:
  - [x] Phiên bản Go: 1.20+ (khuyến nghị bản mới nhất ổn định)
  - [x] Trình soạn thảo/IDE: VS Code hoặc GoLand
  - [x] Hệ điều hành: macOS/Linux/Windows

## Tóm tắt nhanh (TL;DR)
1. Chương mở đầu giúp bạn dựng ứng dụng Go nhỏ để làm quen với cú pháp và công cụ.
2. Điểm cốt lõi: để chạy được chương trình độc lập, bạn cần package `main` và hàm `main()`.
3. `go run` phù hợp khi thử nhanh; `go build` phù hợp khi muốn tạo binary để phân phối/chạy lại.
4. Go khuyến khích mã nguồn rõ ràng, ít ma thuật, dễ đọc ngay từ ví dụ đầu tiên.
5. Khi gặp lỗi, đọc kỹ thông báo biên dịch là cách debug nhanh nhất ở giai đoạn đầu.

## Bản đồ khái niệm
- Chủ đề cốt lõi: cấu trúc chương trình Go tối thiểu.
- Liên hệ với chương trước: không có (chương nhập môn).
- Mở rộng sang chương sau: bối cảnh sử dụng Go và hệ sinh thái công cụ.
- Bối cảnh áp dụng thực tế: tạo service nội bộ nhỏ, script CLI, hoặc demo API ban đầu.

## Thuật ngữ và định nghĩa
| Thuật ngữ | Ý nghĩa trong chương | Vì sao quan trọng |
|---|---|---|
| `package main` | Package đặc biệt cho chương trình thực thi | Không có `main`, chương trình không chạy như app độc lập |
| `func main()` | Điểm vào (entrypoint) của ứng dụng | Xác định nơi bắt đầu thực thi |
| `go run` | Biên dịch tạm và chạy ngay | Tối ưu vòng lặp học và thử nghiệm nhanh |
| `go build` | Biên dịch ra file thực thi | Dùng cho chạy lặp lại và chuẩn bị phát hành |

## Ghi chú chi tiết

### 1. Cấu trúc chương trình tối thiểu
- Giải thích: một chương trình Go có thể chạy được cần `package main` và `func main()`.
- Mô hình tư duy: nghĩ `main()` như cánh cửa chính của ứng dụng.
- Ví dụ tối thiểu:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go")
}
```

- Khi nào nên dùng: mọi ứng dụng CLI/service độc lập đều bắt đầu theo cấu trúc này.
- Khi nào không nên dùng: khi viết thư viện dùng lại thì package không nên là `main`.

### 2. Chạy nhanh bằng `go run`
- Giải thích: `go run` phù hợp khi bạn muốn thử code ngay.
- Mô hình tư duy: như nút "Run" để học nhanh, không cần giữ binary.
- Ví dụ tối thiểu:

```bash
go run .
```

- Ghi chú về hiệu năng/độ dễ đọc: tiện cho phát triển ban đầu, nhưng không thay cho quy trình build phát hành.

### 3. Biên dịch bằng `go build`
- Giải thích: `go build` tạo file thực thi để bạn chạy nhiều lần mà không cần biên dịch lại ngay.
- Ví dụ tối thiểu:

```bash
go build -o app
./app
```

## Phân tích mã theo từng bước
Sử dụng ví dụ in lời chào:
1. Đầu vào và đầu ra mong muốn: không có đầu vào, đầu ra là chuỗi `Hello, Go`.
2. Kiểu dữ liệu/cấu trúc dữ liệu chính: chuỗi (string literal).
3. Luồng xử lý: vào `main()` -> gọi `fmt.Println` -> in ra terminal.
4. Chiến lược xử lý lỗi: ví dụ cơ bản chưa có nhánh lỗi; lỗi chủ yếu đến từ biên dịch.
5. Kiểm tra kết quả cuối cùng: terminal hiển thị đúng nội dung mong muốn.

## Lỗi thường gặp và mẹo debug
| Lỗi | Triệu chứng | Nguyên nhân gốc | Cách khắc phục |
|---|---|---|---|
| Thiếu `package main` | Không chạy được như ứng dụng | Sai loại package | Đổi về `package main` |
| Sai tên hàm `main` | Build lỗi hoặc không có entrypoint | Không đúng chữ ký chuẩn | Dùng đúng `func main()` |
| Import không dùng | Biên dịch báo lỗi import | Go không cho import thừa | Xóa import dư hoặc dùng nó |

## Checklist thực hành tốt
- [x] Tên biến/hàm rõ ràng, nhất quán.
- [x] Hàm nhỏ, làm một việc.
- [x] Xử lý lỗi tường minh khi bắt đầu có I/O.
- [x] Có chạy thử sau mỗi thay đổi nhỏ.
- [x] Mã nguồn được format chuẩn (`gofmt`).

## So sánh và đánh đổi
- Cách tiếp cận A: `go run`
  - Ưu điểm: nhanh, tiện thử nghiệm.
  - Nhược điểm: không tạo artifact rõ ràng để phân phối.
- Cách tiếp cận B: `go build`
  - Ưu điểm: tạo binary ổn định, tiện chạy lại.
  - Nhược điểm: thêm một bước so với chạy nhanh.
- Hướng dẫn chọn giải pháp: học và debug nhanh dùng `go run`; demo/đóng gói dùng `go build`.

## Bài tập luyện tập
### Cơ bản
1. Sửa chương trình để in ra tên của bạn thay vì `Hello, Go`.
2. Thêm dòng in thứ hai mô tả mục tiêu học Go của bạn.

### Trung cấp
1. Tách logic in lời chào thành hàm `printGreeting(name string)` và gọi từ `main()`.
2. Truyền tên từ biến môi trường, nếu không có thì dùng giá trị mặc định.

### Nâng cao
1. Viết CLI nhỏ nhận tên qua cờ lệnh `-name`.
2. Đóng gói chương trình thành binary và thử chạy trên môi trường khác.

## Ý tưởng mini project
- Tên dự án: `hello-go-cli`
- Mục tiêu: tạo chương trình chào người dùng theo tên và thời điểm trong ngày.
- Yêu cầu: hỗ trợ tham số dòng lệnh, giá trị mặc định, thông báo lỗi thân thiện.
- Mở rộng tùy chọn: thêm đa ngôn ngữ (vi/en) qua flag.

## Câu hỏi phỏng vấn
1. Vì sao ứng dụng Go cần `package main` và `func main()`?
2. Khi nào nên dùng `go run` thay vì `go build`?
3. Go xử lý import không sử dụng như thế nào và vì sao?

## Tự đánh giá sau khi học
- Phần mình đã hiểu chắc: cấu trúc app Go tối thiểu, cách chạy/build cơ bản.
- Phần còn mơ hồ: khác biệt sâu hơn giữa module/package ở project lớn.
- Câu hỏi cần quay lại: cách tổ chức thư mục chuẩn khi ứng dụng mở rộng.

## Nhật ký thực hành
| Ngày | Bài thực hành | Thời gian | Kết quả | Hành động tiếp theo |
|---|---|---|---|---|
| 2026-04-17 | Tạo app Hello Go và chạy bằng `go run` | 20 phút | Thành công | Tìm hiểu thêm `go build` và module |

## Tài liệu tham khảo
- Trang sách: Chapter 01 của Pro Go.
- Tài liệu Go chính thức: https://go.dev/doc/
- Bài viết/video bổ sung: A Tour of Go (https://go.dev/tour/)

## Lịch sử cập nhật
| Ngày | Cập nhật | Ghi chú |
|---|---|---|
| 2026-04-17 | Tạo nội dung mẫu ban đầu | Dùng làm chuẩn điền cho các chương khác |
