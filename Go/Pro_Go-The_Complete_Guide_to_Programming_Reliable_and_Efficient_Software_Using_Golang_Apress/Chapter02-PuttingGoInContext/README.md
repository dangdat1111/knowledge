# Chương 02 - Putting Go In Context

## Đối tượng và mục tiêu tài liệu
- Đối tượng: người đã chạy được chương trình Go đầu tiên và muốn hiểu vì sao Go phù hợp cho backend/system.
- Mục tiêu: đặt Go vào đúng bối cảnh kỹ thuật, nắm ưu/nhược điểm và tiêu chí chọn Go cho dự án thực tế.

## Mục tiêu học tập
Sau khi hoàn thành chương này, người học có thể:
- [x] Trình bày ngắn gọn triết lý thiết kế của Go.
- [x] So sánh Go với ngôn ngữ động và ngôn ngữ hệ thống khác ở mức thực dụng.
- [x] Chọn được tình huống nên/không nên dùng Go.
- [x] Xác định được các trụ cột giúp Go vận hành tốt trong môi trường production.

## Yêu cầu trước khi học
- Kiến thức nền tảng cần có:
  - [x] Đã nắm cấu trúc app Go cơ bản (Chapter 01).
  - [x] Hiểu khái niệm compile, runtime, concurrency ở mức tổng quan.
- Công cụ và môi trường:
  - [x] Phiên bản Go: 1.20+
  - [x] Trình soạn thảo/IDE: VS Code/GoLand
  - [x] Hệ điều hành: macOS/Linux/Windows

## Tóm tắt nhanh (TL;DR)
1. Go được thiết kế cho năng suất đội ngũ và độ tin cậy khi vận hành.
2. Điểm mạnh nổi bật: biên dịch nhanh, binary gọn, chuẩn thư viện mạnh, mô hình concurrency rõ ràng.
3. Go ưu tiên sự đơn giản, nhất quán hơn là quá nhiều tính năng ngôn ngữ.
4. Trong backend, microservices, tooling, xử lý mạng, Go thường là lựa chọn rất tốt.
5. Với tác vụ cần generic metaprogramming phức tạp hoặc UI desktop dày đặc, có thể cân nhắc ngôn ngữ khác.

## Bản đồ khái niệm
- Chủ đề cốt lõi: vị trí của Go trong hệ sinh thái ngôn ngữ lập trình hiện đại.
- Liên hệ với chương trước: từ “chạy được chương trình” sang “biết khi nào nên dùng Go”.
- Mở rộng sang chương sau: dùng Go toolchain để hiện thực hóa ưu điểm thực tế.
- Bối cảnh áp dụng thực tế: backend API, job xử lý dữ liệu, CLI nội bộ, hạ tầng cloud-native.

## Thuật ngữ và định nghĩa
| Thuật ngữ | Ý nghĩa trong chương | Vì sao quan trọng |
|---|---|---|
| Simplicity | Triết lý thiết kế đơn giản, ít cơ chế khó đoán | Giảm chi phí đọc code và onboarding |
| Concurrency | Chạy nhiều tác vụ đồng thời với goroutine/channel | Phù hợp ứng dụng I/O và mạng |
| Toolchain | Bộ công cụ như `go build`, `go test`, `go fmt` | Chuẩn hóa quy trình phát triển |
| Deployability | Khả năng đóng gói và triển khai nhanh | Tác động trực tiếp đến vận hành production |

## Ghi chú chi tiết

### 1. Triết lý của Go: đơn giản để mở rộng đội ngũ
- Giải thích: Go tránh quá nhiều “ma thuật” ngôn ngữ, ưu tiên code dễ đọc và dễ review.
- Mô hình tư duy: tối ưu cho “nhiều người cùng bảo trì lâu dài”, không chỉ cho tác giả ban đầu.
- Ví dụ tối thiểu:

```go
func calculateTotal(price, qty float64) float64 {
    return price * qty
}
```

- Khi nào nên dùng: hệ thống cần bàn giao, mở rộng nhóm nhanh.
- Khi nào không nên dùng: dự án đòi hỏi DSL/ngôn ngữ cực kỳ biểu đạt theo hướng meta.

### 2. Go trong backend và hạ tầng
- Giải thích: Go phù hợp dịch vụ mạng nhờ runtime ổn định, thư viện chuẩn mạnh, concurrency tốt.
- Mô hình tư duy: một binary, cấu hình rõ ràng, observability tốt, dễ triển khai.
- Ví dụ tối thiểu:

```go
// Pseudo: khởi tạo HTTP server đơn giản
// mux := http.NewServeMux()
// http.ListenAndServe(":8080", mux)
```

- Ghi chú về hiệu năng/độ dễ đọc: hiệu năng đủ cao cho đa số hệ thống nghiệp vụ, đồng thời giữ codebase dễ bảo trì.

### 3. Khi nào không chọn Go
- Giải thích: không có ngôn ngữ nào tối ưu cho mọi bài toán.
- Ví dụ tối thiểu:

```text
Bài toán phân tích dữ liệu ad-hoc, notebook-first -> thường hợp Python hơn.
```

## Phân tích mã theo từng bước
Ví dụ tình huống chọn ngôn ngữ cho một service nội bộ:
1. Đầu vào: yêu cầu API ổn định, tải trung bình, team cần onboarding nhanh.
2. Kiểu dữ liệu/cấu trúc dữ liệu chính: request/response JSON, cấu hình env.
3. Luồng xử lý: nhận request -> xử lý nghiệp vụ -> trả response.
4. Chiến lược xử lý lỗi: trả mã lỗi rõ ràng, log có ngữ cảnh.
5. Kiểm tra kết quả cuối cùng: dịch vụ chạy ổn định, thời gian phản hồi và vận hành đạt kỳ vọng.

## Lỗi thường gặp và mẹo debug
| Lỗi | Triệu chứng | Nguyên nhân gốc | Cách khắc phục |
|---|---|---|---|
| Chọn Go theo trào lưu | Dự án khó đạt năng suất | Không đánh giá yêu cầu thực | Lập checklist chọn công nghệ theo tiêu chí |
| Tối ưu hiệu năng quá sớm | Mất thời gian, code khó đọc | Thiếu số liệu đo đạc | Profile trước, tối ưu sau |
| Bỏ qua chuẩn công cụ Go | Chất lượng code không đồng đều | Thiếu quy ước nhóm | Chuẩn hóa `go fmt`, `go test`, CI |

## Checklist thực hành tốt
- [x] Chọn ngôn ngữ theo bài toán, không theo xu hướng.
- [x] Ưu tiên code dễ đọc, dễ review.
- [x] Thiết lập quy trình toolchain ngay từ đầu.
- [x] Theo dõi metric vận hành sớm (log, latency, lỗi).
- [x] Đánh giá đánh đổi trước khi cam kết kiến trúc.

## So sánh và đánh đổi
- Cách tiếp cận A: Go cho backend service
  - Ưu điểm: binary độc lập, concurrency tốt, vận hành gọn.
  - Nhược điểm: ít linh hoạt cho kiểu lập trình quá “động”.
- Cách tiếp cận B: Ngôn ngữ động cho backend
  - Ưu điểm: tốc độ thử nghiệm ý tưởng nhanh ở giai đoạn đầu.
  - Nhược điểm: cần thêm kỷ luật hiệu năng và vận hành khi hệ thống lớn dần.
- Hướng dẫn chọn giải pháp: nếu ưu tiên ổn định dài hạn và hiệu năng vừa-cao, Go thường là lựa chọn cân bằng.

## Bài tập luyện tập
### Cơ bản
1. Liệt kê 3 lý do vì sao Go phù hợp cho backend.
2. Liệt kê 2 trường hợp bạn sẽ không chọn Go.

### Trung cấp
1. Viết tài liệu 1 trang “Why Go” cho dự án giả định của bạn.
2. So sánh Go với 1 ngôn ngữ bạn đang dùng theo 5 tiêu chí kỹ thuật.

### Nâng cao
1. Đề xuất kiến trúc service nhỏ bằng Go và mô tả quyết định kỹ thuật.
2. Xây dựng tiêu chí đánh giá công nghệ để team dùng thống nhất.

## Ý tưởng mini project
- Tên dự án: `tech-choice-playbook`
- Mục tiêu: tạo tài liệu mẫu giúp team quyết định khi nào dùng Go.
- Yêu cầu: có bảng tiêu chí, ví dụ use-case, khuyến nghị rõ ràng.
- Mở rộng tùy chọn: thêm template chấm điểm theo trọng số.

## Câu hỏi phỏng vấn
1. Go giải quyết vấn đề gì tốt hơn so với một số ngôn ngữ backend khác?
2. Bạn sẽ giải thích triết lý “simplicity over cleverness” của Go như thế nào?
3. Hãy nêu một tình huống không nên chọn Go và lý do.

## Tự đánh giá sau khi học
- Phần mình đã hiểu chắc: vị trí của Go trong backend và vận hành.
- Phần còn mơ hồ: biên giới giữa “đủ tốt” và “tối ưu sâu” trong từng loại hệ thống.
- Câu hỏi cần quay lại: khi scale rất lớn, cần bổ sung chiến lược nào ngoài ngôn ngữ.

## Nhật ký thực hành
| Ngày | Bài thực hành | Thời gian | Kết quả | Hành động tiếp theo |
|---|---|---|---|---|
| 2026-04-17 | Viết ma trận so sánh Go với ngôn ngữ hiện tại | 35 phút | Hoàn thành bản nháp | Chuẩn hóa thành template dùng chung cho team |

## Tài liệu tham khảo
- Trang sách: Chapter 02 của Pro Go.
- Tài liệu Go chính thức: https://go.dev/doc/
- Bài viết/video bổ sung: Go Proverbs, Effective Go

## Lịch sử cập nhật
| Ngày | Cập nhật | Ghi chú |
|---|---|---|
| 2026-04-17 | Tạo nội dung mẫu ban đầu | Dùng làm chuẩn điền cho các chương khác |
