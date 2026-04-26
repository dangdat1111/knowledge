package main

import (
	"fmt"
	//"math/rand"
)

func main() {

	//fmt.Println(rand.Int())

	//fmt.Println("Hello, Go")
	//fmt.Println(20+20)
	//fmt.Println(20+30)

	//const price float32 = 275.00
	//const tax float32 = 27.50
	//fmt.Println(price+tax)

	//Go có những quy tắc nghiêm ngặt về kiểu dữ liệu và không thực hiện chuyển đổi kiểu tự động,
	//điều này có thể làm phức tạp các tác vụ lập trình thông thường, như minh họa trong Danh sách 4-8.
	//const quantity int = 2
	//fmt.Println("Total:", quantity * (price + tax))
	// Kiểu dữ liệu của hằng số mới là int,
	//đây là lựa chọn phù hợp cho một đại lượng chỉ có thể biểu thị một số nguyên sản phẩm, chẳng hạn.
	//Hằng số này được sử dụng trong biểu thức được truyền cho hàm fmt.Println để tính tổng giá.
	//Nhưng trình biên dịch báo lỗi sau khi biên dịch mã:
	//.\main.go:12:26: invalid operation: quantity * (price + tax) (mismatched types int and float32)
	// => không khớp kiểu dữ liệu -> Golang không cho phép điều đó

	// Understanding IOTA
	//Từ khóa iota có thể được sử dụng để tạo một chuỗi các hằng số nguyên không có kiểu dữ liệu liên tiếp mà không cần gán giá trị riêng lẻ cho chúng.
	//const (
	//	Watersports = iota // 0
	//	Soccer 			   // 1
	//	Chess			   // 2
	//)
	//Mẫu này tạo ra một chuỗi các hằng số, mỗi hằng số được gán một giá trị số nguyên, bắt đầu từ số không.
	//fmt.Println(Watersports, Soccer, Chess)

	// Định nghĩa nhiều hằng số bằng một câu lệnh duy nhất
	//const price, tax float32 = 275.00, 27.50
	//const quantity, inStock = 2, true
	//fmt.Println("Total: ", quantity * (price + tax))
	//fmt.Println("In Stock: ", inStock)

	// Using Variables
	// Biến được định nghĩa bằng từ khoá "var" và không giống như hằng số "const", giá trị của biến có thể thay đổi
	//var price float32 = 275.00
	//price = 300
	//fmt.Println(price) // 300

	//Omitting the Variables's Data Type ( Bỏ qua kiểu dữ liệu của biến ! )
	// Trình biên dịch của Go có thể suy ra kiểu dữ liệu của biến dựa trên giá trị ban đầu ==> cho phép bỏ qua kiểu dữ liệu
	//var price = 275.00 // float64
	//var price2 = price
	//fmt.Println(price2)
	// Việc bỏ qua kiểu dữ liệu không có tác dụng giống nhau đối với 1 biến "var" và 1 hằng số "const"
	// => Trình biên dịch sẽ không cho phép trộn lẫn các kiểu dữ liệu khác nhau
	//var price = 275.00 // Go tự suy luận dấu phẩy động (275.00) là kiểu dữ liệu float64
	//var tax float32 = 27.50
	//fmt.Println(price + tax) // -> invalid operation: price + tax (mismatched types float64 and float32)

	//Omitting the Variables's Data Type Value Assignment (Bỏ qua việc gán giá trị cho biến)
	// các biến có thể định nghĩa mà không cần giá trị ban đầu
	//var price float32
	//price = 300
	//fmt.Println(price)

	// Define Multiple Variables with a Single Statement (Định nghĩa nhiều biến bằng 1 câu lệnh duy nhất)
	//var price, tax = 275.00,27.50
	//var price, tax float32
	//price = 275.00
	//tax = 27.50
	//fmt.Println(price, tax)

	//Using the Short Variables Declaration Syntax (Sử dụng cú pháp khai báo biến ngắn gọn)
	//price,tax, inStock := 275.00, 27.50, true
	//Cú pháp khai báo biến ngắn gọn chỉ có thể sử dụng bên trong các hàm ( func )

	//Using the Short Variables Syntax to Redefine Variables
	//->Sử dụng cú pháp biến rút gọn để định nghĩa lại biến
	// Lỗi: tax redeclared in this block
	//price, tax, inStock := 275.00, 27.50, true
	//var price2, tax = 200.00, 25.00
	//==> Việc định nghĩa lại 1 biến được cho phép nếu sử dụng cú pháp khai báo ngắn gọn (:=)
	// Cách làm đúng:
	//price, tax , inStock := 275.00, 27.50, true
	//price2, tax := 200.00, 25.00

	//Using the Blank Identifier
	//Sử dụng mã định danh trống
	//price, tax, inStock, _ := 275.00, 27.50, true, true
	//var _ = "Alice"
	//Sử dụng "_" để tránh lỗi declared but not used ( Khai báo nhưng không sử dụng )

	//Understanding Pointers
	//first := 100 //-> first = 101
	//second:= first //-> second = 100
	//first++
	//đoạn code trên tạo ra 2 biến -> Go sẽ copy giá trị hiện tại của biến first sau khi tạo biến second => 2 biến này độc lập với nhau
	//==> mỗi biến tham chiếu tới 1 bộ nhớ riêng biệt

	//Defining a Pointers
 	// Con trỏ là 1 biến có giá trị là 1 địa chỉ bộ nhớ
	//first := 100
	//var second *int = &first
	// *int: pointer type
	//first ++
	first := 100
	second := &first
	first++
	fmt.Println(*second)



}
