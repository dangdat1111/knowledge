# Chương 04 - Basic Types Values And Pointers

## Đối tượng và mục tiêu tài liệu
- Đối tượng: người học Go từ cơ bản đến trung cấp, muốn nắm chắc nền tảng kiểu dữ liệu, hằng số, biến và quy tắc ép kiểu.
- Mục tiêu: chuyển nội dung từ `main.go` thành ghi chú có cấu trúc để ôn tập nhanh và áp dụng vào bài toán thực tế.

## Mục tiêu học tập
Sau khi hoàn thành chương này, người học có thể:
- [x] Giải thích vì sao Go không tự động chuyển đổi kiểu dữ liệu.
- [x] Sử dụng `const`, `var`, `:=`, và `iota` đúng ngữ cảnh.
- [x] Nhận diện lỗi `mismatched types` và cách sửa.
- [x] Phân biệt suy luận kiểu giữa hằng số và biến.
- [x] Sử dụng `_` (blank identifier) để tránh lỗi biến không dùng.
- [x] Giải thích sự khác nhau giữa sao chép giá trị và dùng con trỏ.

## Yêu cầu trước khi học
- Kiến thức nền tảng cần có:
  - [x] Cú pháp Go cơ bản (`package`, `import`, `func main`).
  - [x] Biết sử dụng `fmt.Println`.
- Công cụ và môi trường:
  - [x] Phiên bản Go: 1.20+ (khuyến nghị 1.22+).
  - [x] Trình soạn thảo/IDE: VS Code hoặc GoLand.
  - [x] Hệ điều hành: macOS/Linux/Windows.

## Tóm tắt nhanh (TL;DR)
1. Chương này tập trung vào cách Go quản lý kiểu dữ liệu có tính chất chặt chẽ, đặc biệt với hằng số và biến.
2. Go không cho phép phép toán giữa các kiểu khác nhau nếu không tự chuyển đổi rõ ràng (ví dụ `int` với `float32`).
3. `const` dùng cho giá trị không đổi; `var` dùng cho giá trị có thể thay đổi.
4. Trình biên dịch có thể suy luận kiểu, nhưng điều này dễ gây nhầm lẫn giữa `float64` mặc định và `float32` khai báo thủ công.
5. `iota` giúp tạo dãy hằng số nguyên tăng dần nhanh gọn.
6. Cú pháp ngắn `:=` chỉ dùng trong hàm và có thể tái khai báo nếu ít nhất có 1 biến mới.
7. `_` giúp bỏ qua giá trị/biến tạm thời để tránh lỗi `declared but not used`.
8. Con trỏ (`*T`) cho phép tham chiếu cùng vùng nhớ; ví dụ cuối file in ra `101` vì `second` trỏ đến `first`.

## Bản đồ khái niệm
- Chủ đề cốt lõi: constants, variables, type inference, strict typing, `iota`, short declaration, blank identifier, pointers.
- Liên hệ với chương trước: cú pháp cơ bản và in dữ liệu ra màn hình.
- Mở rộng sang chương sau: pointers và cách kiểu dữ liệu ảnh hưởng đến bộ nhớ/hành vi hàm.
- Bối cảnh áp dụng thực tế: tính tổng tiền, quản lý tồn kho, tạo enum cho domain.

## Thuật ngữ và định nghĩa
| Thuật ngữ | Ý nghĩa trong chương | Vì sao quan trọng |
|---|---|---|
| `const` | Hằng số, không thay đổi giá trị sau khi khai báo | Tránh thay đổi ngoài ý muốn, dễ đọc code |
| `var` | Biến có thể thay đổi giá trị | Phù hợp dữ liệu runtime |
| Type inference | Trình biên dịch tự suy ra kiểu từ giá trị khởi tạo | Giảm độ dài code nhưng cần hiểu kiểu mặc định |
| `iota` | Tạo dãy hằng số nguyên tăng dần trong `const (...)` | Viết enum đơn giản, tránh gán thủ công |
| `:=` | Khai báo biến ngắn gọn trong hàm | Nhanh gọn, dễ viết code hằng ngày |
| `_` | Blank identifier để bỏ qua giá trị không cần dùng | Tránh lỗi compile `declared but not used` |
| Pointer (`*T`) | Biến lưu địa chỉ bộ nhớ của biến khác | Cho phép nhiều biến cùng tham chiếu một dữ liệu |

## Ghi chú chi tiết

### 1. Ý chính 1
- Giải thích: Go áp dụng quy tắc kiểu dữ liệu nghiêm ngặt. Không có chuyển đổi tự động giữa `int`, `float32`, `float64` trong phép toán.
- Mô hình tư duy: "Mọi toán hạng trong biểu thức nên cùng một kiểu, hoặc được chuyển đổi rõ ràng".
- Ví dụ tối thiểu:

```go
const price float32 = 275.00
const tax float32 = 27.50
const quantity int = 2

// Lỗi: invalid operation (mismatched types int and float32)
// fmt.Println("Total:", quantity*(price+tax))

fmt.Println("Total:", float32(quantity)*(price+tax))
```

- Khi nào nên dùng: khi cần code an toàn kiểu và hạn chế lỗi ngầm.
- Khi nào không nên dùng: không có "không nên"; đây là quy tắc ngôn ngữ, cần thích nghi đúng cách.

### 2. Ý chính 2
- Giải thích: `const` và `var` khác nhau về khả năng thay đổi giá trị; với `var`, kiểu có thể suy luận từ giá trị khởi tạo.
- Mô hình tư duy: `const` cho giá trị bất biến; `var` cho trạng thái thay đổi theo thời gian.
- Ví dụ tối thiểu:

```go
const price, tax float32 = 275.00, 27.50
var quantity, inStock = 2, true

fmt.Println("Total:", float32(quantity)*(price+tax))
fmt.Println("In Stock:", inStock)
```

- Lưu ý quan trọng từ `main.go`: nếu viết `var price = 275.00` thì `price` thành `float64`; cộng với `tax float32` sẽ báo lỗi.

### 3. Ý chính 3
- Giải thích: `iota` và cú pháp `:=` giúp code gọn, nhưng cần dùng đúng phạm vi.
- Ví dụ tối thiểu:

```go
const (
    Watersports = iota // 0
    Soccer             // 1
    Chess              // 2
)

price, tax, inStock := 275.00, 27.50, true
price2, tax := 200.00, 25.00 // hợp lệ vì có biến mới price2

fmt.Println(Watersports, Soccer, Chess)
fmt.Println(price, tax, inStock, price2)
```

### 4. Ý chính 4
- Giải thích: blank identifier (`_`) dùng để nhận giá trị nhưng cố ý không sử dụng.
- Ví dụ tối thiểu:

```go
price, tax, inStock, _ := 275.00, 27.50, true, true
var _ = "Alice"

fmt.Println(price, tax, inStock)
```

- Mẹo: dùng `_` khi tạm thời chưa dùng biến, hoặc khi chỉ cần một phần kết quả trả về.

### 5. Ý chính 5
- Giải thích: gán trực tiếp tạo bản sao giá trị; con trỏ giúp tham chiếu cùng ô nhớ.
- Ví dụ đang chạy trong `main.go`:

```go
first := 100
second := &first
first++

fmt.Println(*second) // 101
```

- Kết luận: `second` trỏ đến `first`, nên khi `first` tăng thì dereference `*second` cũng thấy giá trị mới.

## Phân tích mã theo từng bước
Sử dụng đúng ví dụ đang được bật trong `main.go`:
1. Đầu vào và đầu ra mong muốn
   - Đầu vào: `first := 100`.
   - Đầu ra: in `101`.
2. Kiểu dữ liệu/cấu trúc dữ liệu chính
   - `first`: kiểu `int`.
   - `second`: kiểu `*int` (con trỏ tới `int`).
3. Luồng xử lý hoặc thuật toán
   - Khai báo `first`.
   - Gán `second := &first` để lấy địa chỉ bộ nhớ của `first`.
   - Tăng `first` lên 1.
   - In `*second` để đọc giá trị tại địa chỉ mà `second` đang trỏ tới.
4. Chiến lược xử lý lỗi
   - Chủ yếu dựa vào lỗi compile-time (`mismatched types`, `redeclared`, `declared but not used`).
5. Kiểm tra kết quả cuối cùng
   - Chạy `go run main.go`.
   - Kỳ vọng đầu ra là `101`.

## Lỗi thường gặp và mẹo debug
| Lỗi | Triệu chứng | Nguyên nhân gốc | Cách khắc phục |
|---|---|---|---|
| `mismatched types int and float32` | Build thất bại khi nhân/chia/cộng/trừ | Trộn `int` với `float32` trong cùng biểu thức | Ép kiểu rõ ràng, ví dụ `float32(quantity)` |
| `mismatched types float64 and float32` | Build thất bại khi cộng `price + tax` | Biến suy luận thành `float64`, biến kia khai báo `float32` | Đồng nhất kiểu ngay từ đầu hoặc ép kiểu |
| `redeclared in this block` với `:=` | Báo lỗi tái khai báo biến cùng tên | Dùng `:=` nhưng không có biến mới nào | Dùng `=` để gán lại, hoặc thêm ít nhất 1 biến mới |
| `declared and not used` | Build thất bại dù logic có vẻ đúng | Có biến/hằng được khai báo nhưng không dùng | Xóa biến, dùng biến, hoặc tạm dùng `_` |

## Checklist thực hành tốt
- [x] Tên biến rõ ràng, nhất quán (`price`, `tax`, `quantity`, `inStock`).
- [x] Đồng nhất kiểu dữ liệu trong biểu thức tính toán.
- [x] Ưu tiên bắt lỗi ở compile-time.
- [x] Sử dụng `const` cho giá trị cố định.
- [x] Chỉ dùng `:=` trong hàm và khi cần khai báo mới.
- [x] Dùng `_` hợp lý để xử lý giá trị không cần dùng.
- [x] Phân biệt rõ sao chép giá trị và tham chiếu qua con trỏ.

## So sánh và đánh đổi
- Cách tiếp cận A: khai báo rõ ràng kiểu (`var price float32 = 275.00`)
  - Ưu điểm: dễ kiểm soát, tránh sai kiểu ngầm.
  - Nhược điểm: dài dòng hơn.
- Cách tiếp cận B: để trình biên dịch suy luận (`price := 275.00`)
  - Ưu điểm: gọn, dễ viết nhanh.
  - Nhược điểm: dễ vô tình dùng `float64` khi hệ thống cần `float32`.
- Hướng dẫn chọn giải pháp:
  - Bài toán nghiệp vụ/dữ liệu nhạy cảm kiểu -> ưu tiên khai báo rõ ràng.
  - Prototype nhanh/nội bộ -> có thể dùng suy luận, nhưng phải kiểm tra kiểu.

## Bài tập luyện tập
### Cơ bản
1. Viết chương trình tính tổng tiền từ `price`, `tax`, `quantity` mà không bị lỗi kiểu.
2. Tạo enum 5 môn thể thao bằng `iota` và in ra giá trị.
3. Viết ví dụ dùng `_` để bỏ qua một giá trị trả về.

### Trung cấp
1. Refactor một đoạn code đang dùng `float64` sang `float32` đồng nhất.
2. Viết ví dụ thể hiện sự khác nhau giữa `const` và `var` trong 1 flow bán hàng.
3. Viết đoạn code chứng minh sự khác nhau giữa `second := first` và `second := &first`.

### Nâng cao
1. Thiết kế package nhỏ dùng `iota` để định nghĩa trạng thái đơn hàng.
2. Viết test xác minh các hàm tính toán không còn lỗi trộn kiểu dữ liệu.

## Ý tưởng mini project
- Tên dự án: Price Calculator CLI.
- Mục tiêu: tính tổng đơn hàng và kiểm tra tồn kho với kiểu dữ liệu an toàn.
- Yêu cầu: nhập giá, thuế, số lượng; in tổng tiền; xử lý chuyển đổi kiểu rõ ràng.
- Mở rộng tùy chọn: thêm enum trạng thái đơn hàng bằng `iota`.

## Câu hỏi phỏng vấn
1. Tại sao Go không tự động chuyển đổi giữa `int` và `float32`?
2. Khi nào nên dùng `const` thay vì `var`?
3. Quy tắc tái khai báo biến với `:=` là gì?
4. Blank identifier (`_`) giải quyết vấn đề gì trong Go?
5. Vì sao `fmt.Println(*second)` in giá trị mới sau khi `first++`?

## Tự đánh giá sau khi học
- Phần mình đã hiểu chắc: strict typing, `const`/`var`, `:=`, `iota`.
- Phần còn mơ hồ: lựa chọn `float32` hay `float64` theo bài toán cụ thể.
- Câu hỏi cần quay lại: khi nào nên ưu tiên độ chính xác số học thay vì tối ưu bộ nhớ?

## Nhật ký thực hành
| Ngày | Bài thực hành | Thời gian | Kết quả | Hành động tiếp theo |
|---|---|---|---|---|
| 2026-04-26 | Ôn tập constants, variables, iota, `:=`, `_`, pointers | 60p | Hiểu được khác biệt copy giá trị và con trỏ, chạy ví dụ ra `101` | Viết thêm bài tập về pointer và test |

## Tài liệu tham khảo
- Trang sách: Chương 04 - Basic Types Values And Pointers.
- Tài liệu Go chính thức: https://go.dev/ref/spec
- Bài viết/video bổ sung: Go by Example (constants, variables, iota).

## Lịch sử cập nhật
| Ngày | Cập nhật | Ghi chú |
|---|---|---|
| 2026-04-26 | Điền đầy đủ nội dung từ `main.go` vào khung README | Bản đầu tiên theo template |
| 2026-04-26 | Bổ sung blank identifier và pointers theo mã mới trong `main.go` | Hoàn thiện phiên bản tiếng Việt có dấu |
