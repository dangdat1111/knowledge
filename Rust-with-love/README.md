# RUST FOR RUSTACEANS — Jon Gjengset (Ghi chú học)

> Sách dành cho người **đã biết Rust cơ bản** (đã đọc *The Rust Programming Language*). Mục tiêu không phải dạy cú pháp, mà giải thích **cơ chế bên dưới**: bộ nhớ, kiểu, trait, async... hoạt động thực sự ra sao.

Bản PDF Early Access có 8 chương: **2-Foundations, 3-Types, 4-Designing Interfaces, 5-Error Handling, 6-Project Structure, 7-Testing, 8-Macros, 9-Async**.
(Thiếu Chương 1 và Chương 10-15: Unsafe Rust, Concurrency, FFI, no_std, Large Projects, Putting It Together.)

---

# CHƯƠNG 2 — FOUNDATIONS (Nền tảng)

## 2.1. Nói về bộ nhớ — Value, Variable, Pointer

Bản chất: Rust phân biệt 3 thứ mà các ngôn ngữ khác thường gộp chung:

- **Value (giá trị)**: = kiểu + một phần tử trong tập giá trị của kiểu đó. Ví dụ số `6` kiểu `u8`. Nó là "ý nghĩa", độc lập với việc nằm ở đâu.
- **Place (nơi chứa)**: một vị trí có thể chứa value (trên stack, heap...). **Variable** là một place có tên trên stack.
- **Pointer (con trỏ)**: là một value, nhưng value của nó là **địa chỉ** của một place.

```
let x = 42;        let y = 43;
let var1 = &x;     let var2 = &x;

   VALUE        PLACE (variable)      POINTER
  ┌──────┐      ┌────────────┐       ┌──────────────┐
  │  42  │◄─────│ x @0x1000  │◄──────│ var1 = 0x1000│
  └──────┘      └────────────┘   │   └──────────────┘
                                 └───┤ var2 = 0x1000│  (cùng trỏ x)
                                     └──────────────┘
```

Điểm "aha": `let s = "Hello world";` → **value thực sự** của `s` là một con trỏ tới ký tự đầu, **không phải** chuỗi. Chuỗi nằm nơi khác (static memory).

## 2.2. Hai mô hình mental về biến

**High-Level Model (mô hình luồng / flow)** — dùng để suy luận lifetime & borrow:
Biến là tên gắn với value. Mỗi lần truy cập, vẽ một "đường" (flow) từ lần dùng trước tới lần dùng này. Biến chỉ "tồn tại" khi đang giữ value hợp lệ.

```
Flow (dòng đời của value):
  x = 42 ●─────────►● đọc x ─────►● đọc x      (OK: các flow tương thích)
         \
          ●──► moved ✗ (sau move không vẽ được flow nữa)
```

Borrow checker thực chất kiểm tra: **không có 2 flow song song xung đột** (vd 2 flow mutable, hoặc 1 flow borrow khi không có flow nào sở hữu value).

**Low-Level Model (ô nhớ)** — dùng khi suy luận unsafe/raw pointer:
Biến = "value slot" (ô nhớ). Gán = đổ đầy ô. Đọc = compiler kiểm tra ô không rỗng (chưa moved/uninit).

> Bản chất: hai mô hình đều đúng, đều là đơn giản hóa. Hiểu cả hai giúp đọc code khó.

## 2.3. Các vùng bộ nhớ

```
        BỘ NHỚ CHƯƠNG TRÌNH
  ┌────────────────────────────┐
  │  STACK  (lớn dần xuống)     │  frame mỗi lần gọi hàm; tự thu hồi khi return
  │  ┌──────────┐              │
  │  │ frame baz│ ← top        │
  │  ├──────────┤              │
  │  │ frame bar│              │
  │  ├──────────┤              │
  │  │ frame foo│              │
  │  ├──────────┤              │
  │  │ frame main│             │
  │  └──────────┘              │
  │         ▲                  │
  │         │ khoảng trống     │
  │         ▼                  │
  │  HEAP  (cấp phát thủ công) │  Box::new → con trỏ; sống tới khi drop
  ├────────────────────────────┤
  │  STATIC  (mã nhị phân,      │  sống suốt chương trình; chuỗi, `static`
  │  hằng, biến static)         │  thường read-only
  └────────────────────────────┘
```

- **Stack**: frame chứa biến cục bộ + tham số. Return → frame biến mất → mọi reference tới nó phải có lifetime ≤ frame.
- **Heap**: sống độc lập với frame. `Box::new(v)` đặt v lên heap, trả con trỏ. Quên drop = **memory leak**. Cố ý leak: `Box::leak` → lấy `'static`.
- **Static**: `'static` lifetime = "sống tới khi chương trình tắt". Hay gặp ở **bound `T: 'static`** = "T tự cung tự cấp, không mượn gì non-static". Ví dụ `thread::spawn` cần `'static` vì thread mới có thể sống lâu hơn thread cũ.
- **`const` vs `static`**: `const` không có địa chỉ/ô nhớ — chỉ là "tên tiện lợi cho một value", được tính ở compile-time và chèn vào nơi dùng.

## 2.4. Ownership (Quyền sở hữu)

Bản chất: **mỗi value có đúng 1 owner**. Owner chịu trách nhiệm drop value. Move = chuyển quyền sở hữu.

```
let x1 = 42;            let y1 = Box::new(84);
{
   let z = (x1, y1);    // x1 COPY vào z (42 là Copy)
}                       // y1 MOVE vào z (Box không Copy)
// sau scope:
//   x1 vẫn dùng được  ✓   (bản sao)
//   y1 KHÔNG dùng được ✗   (đã move)
```

- **Copy trait**: value được sao chép bit thay vì move (i32, bool, f64...). Không thể Copy nếu type chứa thứ phải dealloc (như `Box`) — vì sẽ **double free**.

```
Nếu Box mà Copy được:
  box2 = box1   →   box1 ─┐
                          ├─► [heap 84]   ← cả hai cùng nghĩ mình sở hữu
                  box2 ─┘                   → drop 2 lần = THẢM HỌA
```

- **Drop order**: biến drop theo **thứ tự ngược** (vì biến sau có thể chứa reference tới biến trước — vd hashtable chứa ref tới string khai báo trước). Nhưng giá trị **lồng nhau** (tuple, array, struct) drop theo **thứ tự thuận** (Rust chưa cho self-reference trong 1 value nên không cần ngược).

## 2.5. Borrowing & Lifetimes (Mượn)

**Shared reference `&T`**: nhiều ref đọc cùng lúc, Copy được, **bất biến**. Compiler được phép giả định value **không đổi** khi shared ref còn sống → có thể đọc 1 lần và tái dùng.

**Mutable reference `&mut T`**: **độc quyền** (không ref nào khác tồn tại đồng thời). Compiler giả định không ai alias → tối ưu mạnh.

```
fn noalias(input: &i32, output: &mut i32) {
   if *input == 1 { *output = 2; }   // vì &mut độc quyền,
   if *input != 1 { *output = 3; }   // input và output KHÔNG thể trỏ cùng ô
}                                     // → compiler gộp thành if/else
```

`&mut` chỉ cho mutate **đúng ô** nó trỏ. Khi move value ra khỏi `&mut`, **phải để lại value khác** (không owner sẽ drop "khoảng trống"):

```
std::mem::take(s)     → lấy value ra, để lại Default::default()
std::mem::swap(a, b)  → tráo 2 &mut, không cần sở hữu
*s = new_value        → ghi đè (value cũ drop ngay)
```

## 2.6. Interior Mutability (Khả biến nội tại)

Bản chất: cho phép mutate **qua shared reference `&`** — "phá luật" nhưng an toàn nhờ cơ chế khác (atomic, runtime check).

```
Hai loại:
┌─────────────────────────────┬──────────────────────────────┐
│ Cho lấy &mut từ &            │ Chỉ cho thay value, không &mut│
│ Mutex, RefCell               │ Cell, atomic (AtomicUsize)    │
│ (dựa trên UnsafeCell)        │ (replace/get-set, không ref)  │
└─────────────────────────────┴──────────────────────────────┘
```

`UnsafeCell` là **cách duy nhất** hợp lệ để mutate qua `&`. `Cell` an toàn vì không bao giờ cho reference vào trong + không share giữa thread.

## 2.7. Lifetimes — chi tiết hơn

Bản chất: lifetime **không phải** là scope. Nó là **vùng code mà một reference phải còn hợp lệ**. Trung tâm = **borrow checker**.

**Lifetime không cần liên tục** — có thể có "lỗ hổng":

```
let mut x = Box::new(42);
let mut z = &x;          // 'a bắt đầu
for ... {
   println!("{}", z);    // dùng z (flow từ x)
   x = Box::new(i);      // x bị move → 'a "kết thúc" ở đây
   z = &x;               // 'a "khởi động lại" từ x mới
}
```

```
'a:  ●────dùng────● x moved   ●────dùng────●
     (hợp lệ)     ✗ (lỗ hổng)  (hợp lệ lại)
```

**Generic lifetimes** — khi type chứa reference:

```rust
struct StrSplit<'s, 'p> {
    delimiter: &'p str,   // ← cần 2 lifetime RIÊNG
    document:  &'s str,   //    vì giá trị trả về gắn với 's,
}                         //    không phải 'p (delimiter có thể là String tạm)
```

Nếu dùng chung 1 lifetime → buộc return gắn với cả document lẫn delimiter → không viết được hàm như `str_before` (delimiter là String cục bộ).

Lưu ý drop: nếu type **impl Drop**, drop sẽ "tính là một lần dùng" mọi lifetime → borrow checker kiểm tra ref còn hợp lệ lúc drop.

## 2.8. Lifetime Variance (Biến thiên)

Bản chất: quy định **subtype nào dùng thay được cho type nào**. `A` là subtype của `B` nếu A "ít nhất cũng hữu ích như" B.

```
'static  là subtype của  'a    (sống lâu hơn = hữu ích hơn)
Turtle   là subtype của  Animal (làm được nhiều hơn)
```

Ba loại biến thiên:

```
┌──────────────┬───────────────────────────────────────────────┐
│ COVARIANT    │ Dùng subtype thay được. &'a T covariant theo   │
│ (đồng biến)  │ 'a và T. → truyền &'static str cho &'a str OK  │
├──────────────┼───────────────────────────────────────────────┤
│ INVARIANT    │ Phải đúng y hệt. &mut T invariant theo T.      │
│ (bất biến)   │ Cell<T> invariant. → KHÔNG truyền             │
│              │ &mut Vec<&'static str> cho &mut Vec<&'a str>   │
├──────────────┼───────────────────────────────────────────────┤
│ CONTRAVARIANT│ Ngược lại. Tham số hàm: Fn(T) contravariant    │
│ (nghịch biến)│ theo T. fn(&'a str) hữu ích hơn fn(&'static)   │
└──────────────┴───────────────────────────────────────────────┘
```

Tại sao `&mut T` phải invariant? Vì nếu được nới lỏng, ta có thể nhét chuỗi sống-ngắn vào `Vec<&'static str>`, caller dùng tiếp tưởng là `'static` → dangling.

> Bản chất ứng dụng: **giữ type covariant nhất có thể**. Invariance làm type khó dùng.

---

# CHƯƠNG 3 — TYPES (Kiểu dữ liệu)

## 3.1. Kiểu trong bộ nhớ — Alignment & Layout

**Alignment (căn lề)**: value phải đặt ở địa chỉ là bội số của alignment. Mọi value tối thiểu byte-aligned (con trỏ trỏ byte, không trỏ bit). CPU 64-bit đọc theo khối 8 byte → đọc value "lệch" (misaligned) phải đọc 2 lần + ghép → chậm và sai khi concurrent.

```
i64 đặt ĐÚNG (8-aligned):        i64 đặt LỆCH (bắt đầu byte 4):
 byte: 0       8                  byte: 0   4       12
 ┌───────────┐                    ┌────┬───────┬────┐
 │   i64     │ ← 1 lần đọc        │////│  i64  │////│ ← 2 lần đọc + ghép
 └───────────┘                    └────┴───────┴────┘
```

Quy tắc: kiểu built-in align theo size (u8→1, u16→2, u32→4, u64→8). Kiểu phức tạp align theo field có alignment lớn nhất.

**Layout** — cách compiler sắp xếp field:

```rust
#[repr(C)]
struct Foo { tiny: bool, normal: u32, small: u8, long: u64, short: u16 }
```

```
repr(C) — GIỮ NGUYÊN THỨ TỰ field, chèn padding (26→32 byte):
 [tiny:1][pad:3][normal:4][small:1][pad:7][long:8][short:2][pad:6]
  ▲ bool   ▲để normal     ▲       ▲để long align    ▲để cả struct
           4-aligned                                 align 8 (bội alignment)

repr(Rust) mặc định — ĐƯỢC sắp xếp lại field (giảm dần size) → 16 byte:
 [long:8][normal:4][short:2][small:1][tiny:1]   (gần như không padding)

repr(packed) — KHÔNG padding → tiết kiệm RAM/mạng, nhưng misaligned → chậm/crash
repr(align(n)) — ÉP alignment lớn hơn (tránh false sharing giữa CPU cache line)
```

> Bản chất: Rust mặc định **không hứa** layout → để compiler tự do tối ưu. Cần layout xác định (FFI, raw pointer) thì dùng `repr(C)`.

Các kiểu phức hợp khác: **Tuple** ~ struct cùng thứ tự; **Array** = chuỗi liên tục, không padding; **Union** = layout chọn riêng mỗi variant, alignment = max; **Enum** = như union + 1 field ẩn lưu discriminant (variant nào).

## 3.2. DST & Fat Pointer

**DST (Dynamically Sized Type)**: type không biết size lúc compile — `[u8]` (slice), `dyn Trait` (trait object). Compiler cần size để sinh code → DST không xài trực tiếp được.

Giải pháp: đặt sau **fat pointer** = con trỏ thường + 1 word phụ:

```
&[u8] (slice):              &dyn Trait (trait object):
 ┌──────────┬─────────┐      ┌──────────┬──────────┐
 │ ptr →data│ len     │      │ ptr →data│ ptr →vtable│
 └──────────┴─────────┘      └──────────┴──────────┘
   2 × usize                   2 × usize → CHÍNH NÓ là Sized
```

`T: Sized` là bound ngầm khắp nơi; muốn nhận DST phải opt-out: `T: ?Sized`. `Box`, `Arc` chứa fat pointer được nên hỗ trợ `?Sized`.

## 3.3. Compilation & Dispatch (Static vs Dynamic)

Bản chất của generic: **monomorphization** — compiler "copy-paste" code generic cho từng type cụ thể.

```
STATIC DISPATCH (impl Trait / <T>):           DYNAMIC DISPATCH (dyn Trait):
─────────────────────────────             ─────────────────────────────
fn contains(p: impl Pattern)              fn contains(p: &dyn Pattern)

Mỗi type → 1 bản copy, địa chỉ hàm        1 bản code, caller đưa vtable
biết lúc compile.                          tra địa chỉ hàm lúc runtime.

  String::contains_u8()                    p ──► [vtable]
  String::contains_char()                          ├─ is_contained_in → 0x...
  String::contains_str()                           ├─ layout/size/align
  (code lặp, cache instruction kém,                └─ drop
   compile lâu, NHƯNG nhanh + inline)      (1 code, compile nhanh,
                                            mất tối ưu inline, +1 lần tra vtable)
```

**Trait object** = (con trỏ tới data) + (con trỏ tới vtable). Mọi vtable đều chứa `drop` → **mọi `dyn Trait` cũng là `dyn Drop`**.

**Object-safety** (điều kiện biến thành trait object): method không được generic, không trả `Self`, không có static method (không deref `self`). Vd `Clone` (trả `Self`) → không object-safe. Có thể dùng `where Self: Sized` để loại method khỏi kiểm tra object-safety.

> Quy tắc: **library dùng static dispatch** (để user tự chọn), **binary dùng dynamic dispatch** (code gọn, compile nhanh).

## 3.4. Generic Traits — hai cách

```
trait Foo<T> { }            ← GENERIC PARAMETER: nhiều impl cho 1 type
trait Foo { type Bar; }     ← ASSOCIATED TYPE: chỉ 1 impl cho 1 type
```

> Quy tắc ngón tay cái: **dùng associated type khi chỉ có 1 impl hợp lý**, dùng generic param khi cần nhiều impl. Associated type dễ dùng hơn (không phải lặp lại bound, không phải disambiguate `FromIterator::<u32>::from_iter`).

## 3.5. Coherence & Orphan Rule (Tính nhất quán)

Bản chất: với mỗi (type, method), chỉ có **đúng 1** lựa chọn impl. Nếu không, vd bạn impl `Display for bool` và std cũng có → compiler không biết chọn cái nào.

**Orphan Rule**: chỉ được impl `Trait for Type` nếu **trait HOẶC type là local** trong crate của bạn.

```
✓ impl MyTrait for bool      (trait local)
✓ impl Display for MyType    (type local)
✗ impl Display for bool      (cả hai đều foreign → CẤM)
```

Ngoại lệ:
- **Blanket impl** `impl<T> MyTrait for T`: chỉ crate định nghĩa trait được làm.
- **Fundamental types** (`&`, `&mut`, `Box`) đánh dấu `#[fundamental]`: coi như "trong suốt" → `impl IntoIterator for &MyType` được phép.
- **Covered impl**: param generic được "bọc" trong type khác — `impl From<MyType> for Vec<i32>` OK.

## 3.6. Trait Bounds (Ràng buộc)

Bound không nhất thiết dạng `T: Trait`. Có thể là biểu thức tùy ý:

```rust
where String: Clone,                    // luôn đúng, không có type local
where io::Error: From<MyError<T>>,      // generic ở vế phải
where HashMap<T, usize, S>: FromIterator // ràng buộc cả type phức hợp
where for<'a> &'a Self: IntoIterator,   // HIGHER-RANKED: "với MỌI lifetime 'a"
```

`for<'a>` = higher-ranked trait bound = "reference này impl trait với **bất kỳ** lifetime nào". Hay dùng với `Fn`. Tham chiếu associated type: `I::Item` hoặc `<Type as Trait>::AssocType` để disambiguate.

## 3.7. Marker Traits & Marker Types

```
MARKER TRAIT (không method): Send, Sync, Copy, Sized, Unpin
  → chỉ "đánh dấu một tính chất" của type. Send = "an toàn gửi qua thread".
  → cho phép viết bound thể hiện yêu cầu ngữ nghĩa mà code không biểu đạt trực tiếp.
  → hầu hết là AUTO-TRAIT (compiler tự impl trừ khi type chứa thứ không impl).

MARKER TYPE (unit struct, không data): struct Unauthenticated;
  → đánh dấu TRẠNG THÁI, khiến API không thể dùng sai.
```

```
SshConnection<Unauthenticated>  ──connect()──►  SshConnection<Authenticated>
        │                                               │
   chỉ có connect()                          chỉ ở đây mới có run_command()
```

## 3.8. Existential Types (Kiểu tồn tại)

`impl Trait` ở vị trí return + `async fn` = existential type: "**tồn tại một** type cụ thể thỏa trait, nhưng tôi không nói là gì".

```
fn make() -> impl Iterator<Item=u32>   // caller chỉ biết "nó là Iterator"
                                        // type thật bị ẩn (closure/map/filter...)
```

> Bản chất: **zero-cost type erasure** — ẩn type cụ thể (như iterator/future) khỏi public API → đổi implementation sau mà không phá vỡ code dùng. Khác với `impl Trait` ở **argument** (chỉ là cú pháp tắt cho generic param ẩn danh). Lưu ý: auto-trait (Send/Sync) vẫn được lan truyền qua `impl Trait`.

---

# CHƯƠNG 4 — DESIGNING INTERFACES (Thiết kế giao diện)

Bốn nguyên tắc: **Unsurprising, Flexible, Obvious, Constrained**.

```
       INTERFACE TỐT
  ┌────────────────────────┐
  │ UNSURPRISING  không bất ngờ (tên & trait quen thuộc)
  │ FLEXIBLE      ít ràng buộc, nhiều hứa hẹn linh hoạt
  │ OBVIOUS       dễ hiểu đúng, khó dùng sai
  │ CONSTRAINED   thận trọng khi thay đổi (tránh breaking)
  └────────────────────────┘
```

Tham khảo: Rust API Guidelines, RFC 1105 (API evolution), `cargo clippy`.

## 4.1. Unsurprising (Không gây bất ngờ) — Principle of Least Surprise

- **Naming**: tên giống nhau → hành vi giống nhau. `iter` nhận `&self`, `into_inner` nhận `self`, `SomethingError` impl `Error`.
- **Common traits**: eagerly impl các trait tiêu chuẩn (vì coherence cấm user tự impl foreign trait cho type bạn):

```
NÊN impl: Debug (gần như luôn), Send+Sync (+Unpin), Clone, Default,
          PartialEq, PartialOrd, Hash, Eq, Ord, serde::Serialize/Deserialize
CẨN THẬN với Copy: user không kỳ vọng; bỏ Copy sau này là BREAKING change.
```

- **Ergonomic trait impl**: nên cung cấp blanket impl cho `&T`, `&mut T`, `Box<T>` để `fn foo<T: Trait>(t: T)` gọi được với `&Bar`. Cân nhắc impl `IntoIterator` cho `&MyType` để `for` loop hoạt động.
- **Wrapper Types**: `Deref`/`AsRef` cho cảm giác "kế thừa".

```
t: T  với  T: Deref<Target=U>
   t.method()  →  nếu T không có method, tự deref sang U
```

`Borrow` khác `Deref`/`AsRef`: hẹp hơn, dùng khi type "tương đương" nhau (HashSet cho tra bằng `&str` hoặc `&String`), yêu cầu Hash/Eq/Ord giống hệt.

Hộp lưu ý: type "trong suốt" (như `Arc`) nên **tránh inherent method** — dùng `Arc::clone(&x)` thay vì `x.clone()` để không nhập nhằng "clone Arc hay clone inner?".

## 4.2. Flexible (Linh hoạt) — Contract = Requirements + Promises

Bản chất: mỗi API là một **hợp đồng**. Requirement (ràng buộc đầu vào) và Promise (đảm bảo đầu ra).

```
fn frobnicate1(s: String)        -> String           ← chặt: đòi sở hữu, hứa trả String
fn frobnicate2(s: &str)          -> Cow<'_, str>      ← lỏng hơn: ref vào, Cow ra
fn frobnicate3(s: impl AsRef<str>) -> impl AsRef<str> ← lỏng nhất

Nguyên tắc: ÍT ràng buộc + ÍT hứa cụ thể = dễ tiến hóa.
Thêm ràng buộc / bỏ hứa = BREAKING (major version).
```

- **Generic arguments**: làm tham số generic nếu user có thể muốn truyền type khác. Có thể thay generic bằng `&dyn Trait` để giảm code, nhưng **đừng ép dynamic dispatch lên user** (họ không opt-out được). Đổi từ concrete → generic không hẳn backward-compatible (type inference của user có thể vỡ).
- **Object safety** là một phần hợp đồng: ưu tiên giữ trait object-safe.
- **Borrowed vs Owned**: cần sở hữu (gọi method nhận `self`, gửi qua thread) → lưu owned, để **caller** cấp owned. Không cần → nhận reference. Không chắc → dùng `Cow`.
- **Fallible/Blocking destructors**: `Drop` không báo lỗi được, không async được. → cung cấp thêm **explicit destructor** (method nhận `self`, trả `Result`/dùng `async fn`). Ba cách xử lý xung đột với `Drop`:

```
1. Bọc field trong Option → Option::take trong cả 2 destructor
2. Mỗi field "takeable" → std::mem::take (None/Vec rỗng/...)
3. ManuallyDrop (deref tới inner, nhưng take là UNSAFE)
```

## 4.3. Obvious (Hiển nhiên) — dễ hiểu đúng, khó dùng sai

- **Documentation**: ghi rõ panic, lỗi, điều kiện unsafe. End-to-end example cấp module. Dùng `#[doc(hidden)]`, `#[doc(cfg(..))]`, `#[doc(alias="...")]`, intra-doc links.
- **Type System Guidance**:

```
SEMANTIC TYPING: thay 3 bool bằng 3 enum 2-variant
  fn(DryRun::Yes, Overwrite::No)  ← không thể nhầm thứ tự
                                     (compiler báo lỗi nếu sai chỗ)

ZERO-SIZED TYPE đánh dấu trạng thái (Rocket<Stage>):
  struct Grounded; struct Launched;
  struct Rocket<Stage = Grounded> { stage: PhantomData<Stage> }

  Rocket<Grounded> ──launch()──► Rocket<Launched>
       │                              │
   không có launch                accelerate()/decelerate()
   của Launched                   (chỉ tồn tại sau khi phóng)

  → "làm cho trạng thái bất hợp lệ KHÔNG THỂ biểu diễn được"
```

- `#[must_use]`: cảnh báo nếu user bỏ qua giá trị trả về (như `Result`).

## 4.4. Constrained (Bị ràng buộc) — quản lý breaking change

```
THAY ĐỔI                                        BREAKING?
─────────────────────────────────────────────────────────
Đổi tên/xóa public type                          CÓ
Thêm private field vào struct (mất constructor)  CÓ
tuple struct → struct thường (pattern match vỡ)  CÓ
Thêm method có default vào trait                  KHÔNG (trừ khi mất object-safe)
Thêm method KHÔNG default vào trait               CÓ (phá impl)
impl foreign trait cho existing type              CÓ (gọi nhập nhằng)
impl trait cho NEW type                           KHÔNG
```

Công cụ:
- `#[non_exhaustive]` — cấm user dùng constructor ngầm/pattern match đầy đủ → bạn được thêm field/variant sau.
- `pub(crate)`, `pub(in path)` — giảm bề mặt public.
- **Sealed traits** — trait chỉ dùng được, không impl được bởi crate khác → bạn thêm method thoải mái:

```rust
pub trait CanUseCannotImplement: sealed::Sealed { }
mod sealed {
    pub trait Sealed {}
    impl<T> Sealed for T where T: TraitBounds {}  // chỉ type bạn cho phép
}
```

- **Hidden contracts**: re-export type foreign → version của foreign crate trở thành một phần hợp đồng của bạn. **Auto-traits** (Send/Sync) tự động → nếu private field mất Send thì type public cũng mất Send = breaking ngầm.

```
THE SEMVER TRICK (David Tolnay): khi type T giữ nguyên qua major bump,
release bản 1.x mới re-export T từ 2.0 → cả hai major chia sẻ MỘT type T.

Test bắt auto-trait breaking:
  fn is_normal<T: Sized + Send + Sync + Unpin>() {}
  #[test] fn normal_types() { is_normal::<MyType>(); }  // không chạy, chỉ compile
```

---

# CHƯƠNG 5 — ERROR HANDLING (Xử lý lỗi)

## 5.1. Biểu diễn lỗi — Enumeration vs Erasure

```
ENUMERATION (liệt kê)              ERASURE (xóa kiểu / opaque)
──────────────────────            ───────────────────────────
pub enum CopyError {              Box<dyn Error + Send + Sync + 'static>
   In(io::Error),                 hoặc struct opaque với field private
   Out(io::Error),
}
Dùng khi: caller CẦN phân biệt    Dùng khi: caller KHÔNG thể làm gì
nguyên nhân để xử lý khác nhau    khác nhau dù biết nguyên nhân
(vd lỗi input vs output)          (vd image decoder lỗi size vs nén)
```

Type lỗi tốt nên:
1. impl `std::error::Error` (đặc biệt `Error::source` → in backtrace tới gốc).
2. impl `Display` (1 dòng, chữ thường, không dấu chấm cuối) + `Debug` (chi tiết: port, request id...).
3. impl `Send + Sync` (dùng được đa luồng, ghép với `io::Error`).
4. là `'static` (dễ propagate, dễ type-erase, cho phép **downcasting**).

```
DOWNCASTING — lấy lại type cụ thể từ dyn Error:
  let e: &dyn Error = ...;
  if let Some(io) = e.downcast_ref::<io::Error>() { ... }
  → CHỈ chạy khi 'static. Cơ chế: TypeId (std::any::TypeId)
    Error::type_id (qua vtable) so với TypeId của type đích.
```

Lưu ý: type-erased errors **compose** đẹp — hàm trả `Box<dyn Error>` thì dùng `?` với mọi loại lỗi đều "just work". Cộng đồng đồng thuận: lỗi nên hiếm → đặt sau con trỏ (`Box`/`Arc`) để không phình `Result` ở happy path.

## 5.2. Special Error Cases (Lỗi đặc biệt)

```
Result<T, ()>   vs   Option<T>
  Err(())        =  "thất bại, cần xử lý"   (#[must_use], () KHÔNG impl Error)
  None           =  "không có gì để trả"    (không phải lỗi)
  → ĐỪNG "đơn giản hóa" Result<T,()> thành Option<T>: khác ngữ nghĩa.
  → Tốt hơn: tự định nghĩa unit struct impl Error.

never type (!) — value không bao giờ tạo được:
  fn loop_forever() -> Result<(), !>   // Err không bao giờ xảy ra
  → compiler bỏ luôn code panic cho Result<T, !>, pattern match không cần liệt kê !

std::thread::Result<T> = Result<T, Box<dyn Any + Send + 'static>>
  → error là dyn Any (chỉ biết "là gì đó"), vì nó = "đối số truyền cho panic!"
```

## 5.3. Handling Errors — toán tử `?`

```
?  =  "unwrap hoặc return sớm"  +  chuyển kiểu qua trait From
  fn f() -> Result<T, E> {
     let x = fallible()?;   // nếu Err(X), tự From::<X> → E rồi return Err(E)
  }
```

```
FROM & INTO: tại sao cần CẢ HAI?
  - Lịch sử: coherence rules cũ buộc phải có cả hai.
  - Ergonomic: fn(impl Into<Foo>) vs fn<T>(...) where Foo: From<T>
  - QUY TẮC: impl FROM, dùng INTO trong bound (Into có blanket impl từ From).

Try trait (chưa ổn định): ? là cú pháp tắt cho Try.
  Áp dụng được cho Result, Option, Poll<Result<T,E>>...
  Try blocks: phạm vi hóa ? để KHÔNG return cả hàm (chạy cleanup):
    let r = try { ...code dùng ?... };  thing.cleanup();  r
```

---

# CHƯƠNG 6 — PROJECT STRUCTURE (Cấu trúc dự án)

## 6.1. Features (Cờ tính năng)

Bản chất: feature = **build flag** bật chức năng tùy chọn. Nguyên tắc tối thượng: features phải **ADDITIVE** (chỉ thêm, không xóa/đổi) — vì Cargo lấy **hợp** (union) features khi nhiều crate cùng phụ thuộc.

```
crate A ─┐ cần feature X của C ─┐
         ├─► crate C  ──────────┤ Cargo bật C MỘT lần với X ∪ Y
crate B ─┘ cần feature Y của C ─┘
  → nếu X, Y loại trừ nhau (mutually exclusive) → build vỡ!
```

```toml
[features]
derive = ["syn"]          # bật optional dependency
default = ["derive"]      # feature mặc định
[dependencies]
syn = { version="1", optional=true }   # optional dep tự thành feature cùng tên
```

Trong code dùng **conditional compilation**: `#[cfg(feature="x")]` (chọn item để compile) và `cfg!(feature="x")` (biểu thức runtime). Ưu tiên `#[cfg]` vì additive. CI nên test mọi tổ hợp feature (tool: `cargo-hack`).

## 6.2. Workspaces (Không gian làm việc)

Bản chất: crate = **một đơn vị compile** (compiler coi cả crate như 1 file lớn). Crate lớn → đổi 1 dòng phải compile lại cả crate. Workspace chia crate lớn thành nhiều **subcrate** phụ thuộc lẫn nhau:

```
my-workspace/
├── Cargo.toml          [workspace] members = ["foo", "bar/one", "bar/two"]
├── Cargo.lock          (CHUNG cho cả workspace)
├── foo/Cargo.toml
├── bar/one/Cargo.toml  → deps: foo = { path="../../foo" }
└── bar/two/Cargo.toml  → deps: one = { path="../one" }

Đổi bar/two → CHỈ bar/two compile lại (foo, one không đổi).
cargo check/test ở root → chạy cho TẤT CẢ members.
```

Hộp lưu ý: dùng **path dependency chỉ khi phụ thuộc thay đổi chưa publish**; nếu không dùng `version` để Cargo không coi 2 phiên bản same-name là 2 crate khác nhau.

## 6.3. Project Configuration

```
[patch]   — thay nguồn dependency tạm thời (test bản sửa lỗi):
  [patch.crates-io]
  regex = { path = "/home/jon/regex" }
  serde = { git = "...", branch = "faster" }
  (KHÔNG đưa vào crate publish — chỉ áp dụng crate gốc)

[profile] — tùy chọn compile:
  opt-level (0..3, "s"=size)  codegen-units (nhiều=compile nhanh, tối ưu kém)
  lto (link-time optimization: tối ưu xuyên compilation unit)
  debug, debug-assertions, overflow-checks

[profile.*.panic]:
  unwind (mặc định) — gỡ stack, drop value, có cleanup
  abort — thoát ngay, không cleanup (embedded, panic toàn cục)
```

```
UNWINDING khi panic:
  baz panic → bar → foo → main → exit
  (mỗi frame drop value bình thường → cleanup tài nguyên)
  Thread khác vẫn chạy; chỉ chết khi thread main thoát.
```

Profile overrides: `[profile.dev.package.serde]` opt-level=3 — tăng tối ưu cho 1 dep nặng ở debug mode.

**CRATE vs PACKAGE**: crate = cây module bắt đầu từ 1 file gốc (`lib.rs`/`main.rs`). Package = tập hợp crate + metadata (mô tả bởi 1 `Cargo.toml`), có thể chứa nhiều binary, test, workspace member.

## 6.4. Conditional Compilation & Versioning

```
#[cfg(condition)] options:
  feature = "x"        (tính năng)
  target_os="windows"/unix, windows, target_family
  test, doc, doctest   (ngữ cảnh)
  miri (tool)          target_arch="x86", target_feature="avx"  target_env="gnu/msvc/musl"

[target.'cfg(windows)'.dependencies]  → dep theo nền tảng
  winrt = "0.7"
```

CI audit dep: `cargo-deny`, `cargo-audit` (phát hiện dep trùng major, lỗ hổng bảo mật, license cấm).

```
VERSIONING — Semantic Versioning (breaking=major, add=minor, fix=patch):

MSRV (Minimum Supported Rust Version): tăng minor khi đổi MSRV
  → user kẹt Rust cũ dùng version="2, <2.7" để pin.

Minimal Dependency Versions: ghi version SỚM NHẤT thỏa mãn, không phải mới nhất.
  hugs="1.7.3" có thể vỡ nếu crate khác cần hugs<1.6. Dùng -Z minimal-versions test.

Changelogs: giữ thủ công (Keep a Changelog), đừng dump git log.

Unreleased versions: sau release 2.0.3 → đặt ngay 2.0.4-alpha.1.
  thay đổi additive → 2.1.0-alpha.1; breaking → 3.0.0-alpha.1.
```

---

# CHƯƠNG 7 — TESTING (Kiểm thử)

## 7.1. Cơ chế test của Rust

```
cargo test --lib  →  rustc --test  →  2 hiệu ứng:
  1. bật cfg(test)         (conditional compilation cho code test)
  2. sinh TEST HARNESS     (hàm main tự sinh, gọi mọi #[test])

UNIT test (#[test] trong src)     INTEGRATION test (thư mục tests/)
  → cùng crate, thấy private        → mỗi file = 1 crate riêng, chỉ thấy public API
  → cfg(test) = true                → chạy với crate KHÔNG có cfg(test)

harness = false  → tự viết fn main làm test runner (cho benchmark, fuzzer)
```

Tham số test harness (qua `cargo test -- ...`): `--nocapture`, `--test-threads=1`, `--skip`, `--ignored`, `--list`.

## 7.2. `#[cfg(test)]` — code chỉ tồn tại khi test

```
TẠI SAO cần code chỉ-khi-test?
┌──────────────────────────────────────────────────────────┐
│ MOCKING: tạo type/hàm "giả" để kiểm soát (network, data)   │
│   → thường qua generics hoặc cfg(test)/feature. crate mockall│
│                                                            │
│ TEST-ONLY API: lộ trạng thái nội bộ cho test mà không cho   │
│   production (vd HashMap → RawTable::buckets()):           │
│     impl RawTable { #[cfg(test)] pub(crate) fn buckets() }  │
│                                                            │
│ BOOKKEEPING: đếm thêm chỉ khi test (vd BufWriter đếm số     │
│   lần gọi write thực) — không tốn phí ở production         │
└──────────────────────────────────────────────────────────┘
```

## 7.3. Doctests

Code trong `///` doc comment được chạy như test. Chạy như **integration test** (chỉ thấy public API). Tự bọc `fn main`. Tính năng:

```
/// ```
/// # let hidden = 0;        ← dòng có # ẩn khỏi doc nhưng VẪN chạy
/// assert!(frobnify(x).is_ok());
/// ```
Attribute: should_panic, ignore, no_run, compile_fail
  ```compile_fail   ← test rằng code KHÔNG compile (vd type không Send)
  # struct MyNonSendType(Rc<()>);
  fn is_send<T: Send>() {}
  is_send::<MyNonSendType>();   // phải FAIL compile
  ```
```

(Test lỗi compile chi tiết hơn: crate `trybuild`.)

## 7.4. Công cụ test bổ sung

```
LINTING: clippy có lint "correctness" (bắt bug gần chắc chắn):
  a=b; b=a (không swap), mem::forget(&ref), for x in y.next()
  rustc: #![warn(rust_2018_idioms)], missing_docs, missing_debug_implementations

TEST GENERATION:
  FUZZING — sinh input ngẫu nhiên, tìm crash. libfuzzer + crate arbitrary
    (Arbitrary trait biến byte ngẫu nhiên → type Rust). cargo-fuzz.
  PROPERTY-BASED — mô tả tính chất, framework sinh input kiểm tra. proptest.
    (so kết quả impl thật vs impl naive đơn giản)
    "fuzzing = property testing với property = 'không crash'"
```

```
TEST AUGMENTATION (bắt bug phi xác định):
  MIRI — interpreter chạy MIR, phát hiện UB: đọc uninit, dùng sau drop,
         truy cập ngoài biên, tạo 2 &mut tới cùng ô.
    error: Undefined Behavior: trying to reborrow ...
  LOOM — chạy test với MỌI thứ tự interleaving của thread đồng bộ
         (Mutex A trước B, rồi B trước A...) → bắt race condition.
```

## 7.5. Performance Testing

```
3 cạm bẫy:

1. PERFORMANCE VARIANCE: kết quả dao động (CPU clock, paging, nhiệt độ...)
   → chạy NHIỀU lần, xem PHÂN PHỐI (p95) không phải 1 số.
   crate: hdrhistogram, criterion (thống kê + null hypothesis testing)

2. COMPILER OPTIMIZATION: compiler xóa code "vô dụng" → benchmark sai:
   for i in 0..4 { vs.push(i); }  ← bị tối ưu mất hoàn toàn!
   → std::hint::black_box: "giả định arg được dùng theo cách hợp lệ"
     black_box(vs.as_ptr()); vs.push(i); black_box(vs.as_ptr());
     (dùng as_ptr không &vs — vì black_box giả định thao tác HỢP LỆ)

3. I/O OVERHEAD: vô tình đo nhầm:
   for i in 0..1_000_000 { println!("{}", i); my_function(); }
   ← thực ra đo thời gian PRINT 1 triệu số, không phải my_function!
   → vòng lặp benchmark chỉ chứa ĐÚNG code cần đo.
```

---

# CHƯƠNG 8 — MACROS (Vĩ lệnh)

Bản chất: macro = công cụ để **compiler viết code thay bạn**. Bạn đưa "công thức" sinh code, compiler thay mỗi lần gọi macro bằng kết quả. Khác C/C++ `#define` (thay text bừa bãi) — macro Rust theo luật rõ ràng, **kháng lạm dụng**.

```
2 loại:  DECLARATIVE (macro_rules!)    PROCEDURAL (crate riêng, chạy code)
```

## 8.1. Declarative Macros (`macro_rules!`)

Bản chất: "tìm-và-thay-thế có cấu trúc, được compiler hỗ trợ". Bạn **khai báo** output trông thế nào với input thế nào — compiler tự lo việc parse.

```rust
macro_rules! test_battery {
    ($($t:ty as $name:ident),*) => {     // ← MATCHER (mẫu)
        $(                                // ← lặp cho mỗi cặp
            mod $name {
                #[test] fn frobnified()   { test_inner::<$t>(1, true) }
                #[test] fn unfrobnified() { test_inner::<$t>(1, false) }
            }
        )*
    }
}
test_battery! { u8 as u8_tests, i128 as i128_tests };
```

```
CẤU TRÚC: macro_rules! tên {
   (matcher 1) => { transcriber 1 };   ← thử matcher từ trên xuống
   (matcher 2) => { transcriber 2 };
}

MATCHER: token tree khớp input. Fragment types:
  $x:ident (định danh)  $e:expr (biểu thức)  $t:ty (kiểu)  $t:tt (1 token tree)
  $($k:expr => $v:expr),+   ← lặp 1+ lần, phân cách dấu phẩy

TRANSCRIBER: metavariable ($key) được thay bằng phần input đã khớp.
  $(map.insert($key, $value);)+   ← lặp, mỗi lần thay đúng cặp
```

```
HYGIENE (vệ sinh): macro_rules! HYGIENIC với biến — biến trong macro
ở "vũ trụ riêng", không đụng biến ở call site:
  macro_rules! let_foo { ($x:expr) => { let foo = $x; } }
  let foo = 1;
  let_foo!(2);          // sinh `let foo = 2` nhưng ở vũ trụ khác
  assert_eq!(foo, 1);   // ✓ foo gốc = 1 (KHÔNG bị shadow)

  → NHƯNG hygiene CHỈ áp dụng cho biến. Type/module/function CHIA SẺ
    namespace với call site (cố ý — để macro sinh được type/impl).
  → Để chia sẻ biến CÓ CHỦ Ý: nhận $i:ident làm tham số rồi mới gán.
```

Lưu ý quan trọng: macro_rules! chỉ tồn tại **sau khi khai báo** trong source (textual scoping). `mod foo` (khai báo macro) phải đứng **trước** `mod bar` (dùng macro) trong `lib.rs`! Ngoại lệ: `#[macro_export]` nâng macro lên gốc crate. Trong macro tái dùng được nên dùng full path `::core::option::Option`, tránh `::std` để chạy cả `no_std`; dùng `$crate` để trỏ crate định nghĩa macro.

## 8.2. Macros hoạt động thế nào — Token & AST

```
NGUỒN  ──lexer──►  TOKENS  ──parser──►  AST (cây cú pháp)
"(value + 4)"      ( value + 4 )         biểu thức cộng

Macro được đánh giá GIỮA token → AST:
  - compiler parse tới chỗ gọi macro, HOÃN parse phần trong macro
  - đánh giá macro trên chuỗi token → sinh token → parse → chèn AST vào

Input macro = chuỗi TOKEN TREE (không cần là Rust hợp lệ, nhưng phải parse được):
  for <- x    ← hợp lệ TRONG macro (compiler parse được)
  for {       ← KHÔNG (thiếu ngoặc đóng)

Output declarative macro = LUÔN là Rust hợp lệ (expr/statement/item/type/match arm)
  → không thể sinh code Rust sai → kháng lạm dụng.
```

## 8.3. Procedural Macros (Macro thủ tục)

Bản chất: parser + sinh code, bạn viết "keo" ở giữa. Bạn viết **CÁCH** sinh code (chạy chương trình), không phải code cuối. Compiler đưa chuỗi token, bạn trả chuỗi token.

```
3 loại:
┌────────────────┬──────────────────────────────────────────────┐
│ FUNCTION-LIKE  │ giống macro_rules! → THAY thế. Không hygienic.  │
│   foo!(...)    │ Dùng khi: declarative quá phức tạp / cần tính   │
│                │ toán compile-time (phf, hex-literal).          │
├────────────────┼──────────────────────────────────────────────┤
│ ATTRIBUTE      │ #[test] → THAY item, nhận 2 input (attr+item). │
│   #[route(...)]│ Dùng: test gen, framework (rocket #[get]),      │
│                │ middleware (tracing), type transform (pin_project)│
├────────────────┼──────────────────────────────────────────────┤
│ DERIVE         │ #[derive(X)] → THÊM (append) vào item.         │
│                │ Dùng: tự động hóa impl trait HIỂN NHIÊN         │
│                │ (Debug, Clone, serde). Có helper attr #[serde(skip)]│
└────────────────┴──────────────────────────────────────────────┘
```

```
CHI PHÍ: procedural macro làm CHẬM compile vì:
  1. kéo theo dep nặng (syn parse Rust mất hàng chục giây)
  2. dễ sinh nhiều code mà không nhận ra → compiler vẫn phải parse/compile
  → giảm bằng: tắt feature thừa, compile macro ở debug mode.
```

## 8.4. TokenStream & Span

```
TokenStream = chuỗi TokenTree (token đơn, hoặc nhóm trong ()/{}/[])
  parse: dùng crate syn → AST dễ duyệt
  sinh: dựng thủ công + Extend, hoặc "...".parse::<TokenStream>()

SPAN — phần "ma thuật": mỗi token mang span = vị trí gốc trong source.
  → lỗi compiler trong code macro sinh ra được TRỎ về đúng chỗ user viết:
     name_as_debug!(u31)  → compiler trỏ "u31" dù lỗi nằm trong code sinh ra
  → compile_error! + span = sinh thông báo lỗi đẹp trỏ đúng phần input.

Span cũng điều khiển HYGIENE của proc macro:
  Span::call_site()  → định danh resolve ở call site (không cách ly)
  Span::mixed_site() → resolve ở định nghĩa macro (hygienic như macro_rules!)
```

> Quy tắc: code đổi theo **type** → dùng generics; lặp lại thuần túy → dùng macro. Generics ergonomic hơn macro.

---

# CHƯƠNG 9 — ASYNCHRONOUS PROGRAMMING (Lập trình bất đồng bộ)

## 9.1. Vì sao cần async? — CPU dành phần lớn thời gian CHỜ

```
CPU thực tế:  ███ chạy  ░░░░░░░░░░░░░░ CHỜ (network, disk, RAM, chuột) ░░░░░░░
              vài lệnh   "eons" trôi qua giữa các sự kiện → utilization vài %
```

```
SYNCHRONOUS / BLOCKING:                 ASYNCHRONOUS / NONBLOCKING:
─────────────────────                   ──────────────────────────
chờ op trước xong mới chạy op sau        op có thể "chưa xong, quay lại sau"
1 thread = 1 việc tại 1 thời điểm        1 thread = nhiều việc xen kẽ

Giải pháp cũ: MULTITHREADING            Giải pháp: trả về Poll
  + concurrency (OS chọn thread chạy)     enum Poll<T> { Ready(T), Pending }
  - tốn: nhiều thread, đổi thread đắt,    poll_xxx: "làm được bao nhiêu làm,
    buộc dùng Arc/Mutex thay Rc/Cell        rồi return, NHỚ chỗ dừng để tiếp"
```

```
CONCURRENCY vs PARALLELISM:
  concurrency (xen kẽ):   _-_-_-_   (nhiều việc, có thể 1 core)
  parallelism (song song): =====    (nhiều việc thật sự cùng lúc, nhiều core)
```

## 9.2. Future trait — polling chuẩn hóa

```rust
// Phiên bản đơn giản hóa:
trait Future {
    type Output;
    fn poll(&mut self) -> Poll<Self::Output>;
}
```

Future = "value chưa có ngay". `poll` → `Ready(T)` (xong, gọi là **resolve**) hoặc `Pending`. Quy tắc: **không poll lại sau khi Ready** (future được quyền panic; future an toàn poll lại gọi là *fused*).

## 9.3. Vì sao cần async/await? — viết tay Future quá khổ

```
async fn forward<T>(rx: Receiver<T>, tx: Sender<T>) {   ← NGẮN, dễ đọc
    while let Some(t) = rx.next().await {
        tx.send(t).await;
    }
}
```

Nếu viết tay thì phải dựng **state machine** (máy trạng thái):

```
enum Forward<T> {                       ← mỗi await = 1 trạng thái
    WaitingForReceive(ReceiveFuture, Option<Sender>),
    WaitingForSend(SendFuture, Option<Receiver>),
}
impl Future for Forward {
    fn poll(&mut self) -> Poll<()> {
        match self {                    ← khôi phục chỗ dừng lần trước
            WaitingForReceive(recv, tx) => { ...poll recv... }
            WaitingForSend(send, rx)    => { ...poll send... }
        }
    }
}
```

```
async/await tự sinh state machine này cho bạn.
  await/yield  ≈  "return nhưng giữ lại trạng thái bên lề"
```

(Receiver giống "async Iterator" → std đang tiến tới trait có `poll_next`, gọi là **stream**.)

## 9.4. Generators — cơ chế bên dưới async/await

Bản chất: generator = đoạn code có thể **dừng giữa chừng (yield)** rồi **tiếp tục** đúng chỗ. async/await desugar thành generator:

```
generator fn forward(...) {
    loop {
        let mut f = rx.next();
        let r = if let Poll::Ready(r) = f.poll() { r } else { yield };  ← chưa xong → yield
        ...
    }
}
```

Compiler **chèn ngầm** code lưu biến cục bộ vào "cấu trúc dữ liệu của generator" thay vì stack → khi resume, đọc lại từ đó.

```
THE SIZE OF GENERATORS: cấu trúc generator phải chứa state tại MỌI điểm yield.
  async fn chứa [u8; 8192] → generator chứa 8KiB. Future lồng nhau → generator
  PHÌNH TO → tốn memcpy khi move future. Quá to → Box::pin (đẩy lên heap).
```

## 9.5. Pin & Unpin — vấn đề self-referential

Bản chất vấn đề: generator giữ **reference tới chính biến cục bộ của nó** (self-referential). Nếu generator bị **move**, các reference đó trỏ địa chỉ cũ → **dangling**.

```
TRƯỚC khi move:                     SAU khi move:
 generator @0x1000                   generator @0x5000
 ┌─────────────┐                     ┌─────────────┐
 │ rx @0x1004  │◄──┐                 │ rx @0x5004  │     ┌── ref vẫn trỏ
 │ future {    │   │                 │ future {    │     │   0x1004 (CŨ!)
 │   ref→0x1004│───┘                 │   ref→0x1004│─────┘   → DANGLING ✗
 │ }           │                     │ }           │
 └─────────────┘                     └─────────────┘
```

Giải pháp: **`Pin<P>`** (ngăn move) + **`Unpin`** (marker: "type này move an toàn dù đã pin").

```rust
struct Pin<P> { pointer: P }   // P là POINTER type (Box<T>, &mut T...), không phải T
impl<P> Pin<P> where P: Deref {
    pub unsafe fn new_unchecked(pointer: P) -> Self;   // unsafe: hứa target không move
}
```

```
Future trait THẬT (gần đúng):
  fn poll(self: Pin<&mut Self>, cx: &mut Context) -> Poll<Self::Output>

Pin giữ POINTER (không phải T trực tiếp) vì: nếu Pin giữ T, move Pin = move T.
  Pin<Box<T>>, Pin<&mut T> → move Pin chỉ move con trỏ, T đứng yên. ✓

new_unchecked UNSAFE vì compiler không kiểm tra được pointer có thật sự
  giữ target bất động (caller có thể Pin::new(&mut foo), gọi method, drop Pin,
  rồi move foo → self-ref vỡ).
```

```
UNPIN — "key to safe pinning":
  impl Unpin for T = "T move an toàn kể cả khi đã pin" (không xài đảm bảo của Pin)
  → auto-trait (như Send/Sync). Generator là !Unpin.
  → với T: Unpin, cung cấp API AN TOÀN:
      Pin::new (safe), DerefMut (safe), Pin::into_inner
  Box LUÔN Unpin (move Box<T> không move T) → Box::pin AN TOÀN.

OBTAIN PIN:
  Unpin   → Pin::new(&mut future)               (an toàn)
  !Unpin  → Box::pin (heap, an toàn, +1 alloc)  HOẶC
            pin_mut! macro (stack, cần unsafe, shadowing để cấm dùng value cũ)
```

```rust
// Desugar <expr>.await thực sự:
match expr {
    mut pinned => loop {
        match unsafe { Pin::new_unchecked(&mut pinned) }.poll() {
            Poll::Ready(r) => break r,
            Poll::Pending  => yield,   // ← an toàn vì pinned nằm trong generator đã pin
        }
    }
}
```

## 9.6. Executor, Waker, Task — ai gọi poll?

```
EXECUTOR: vòng lặp poll các future. KHÔNG busy-loop (phí CPU) → NGỦ khi
  không future nào tiến triển được, chỉ DẬY khi có việc.

WAKER: cơ chế đánh thức executor. Executor đưa Waker vào mỗi poll (qua Context).
  Future trait THẬT:
    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output>

  Waker.wake() → "future này RUNNABLE" (đáng poll lại), không hẳn "wake".
  Cơ chế: RawWakerVTable (vtable thủ công như dynamic dispatch chương 3).
```

```
THE POLL CONTRACT: nếu poll trả Pending, future PHẢI đảm bảo có gì đó
  gọi wake(Waker) khi nó sẵn sàng tiến triển.

LEAF FUTURE (lá, không chứa future con — TCP, channel):
  - sự kiện NỘI BỘ (channel): lưu Waker vào channel; sender gọi wake khi gửi.
  - sự kiện NGOÀI (TCP, timer): executor gom mọi event source → 1 lệnh blocking
    tới OS (epoll trên Linux); OS báo có event → executor gọi wake.

REACTOR = phần executor mà leaf future đăng ký event source.
```

```
CÂY FUTURE (asynchronous program):
                  TASK (root future, điểm tiếp xúc với executor)
                   │  executor gọi poll(task)
        ┌──────────┴──────────┐
     future                future (select/join = "subexecutor")
        │                     │
     future                future
        │                     │
   LEAF (TCP)            LEAF (channel)  ← chỗ tương tác với Waker

TASK = root future. Executor chỉ poll task; từ đó mỗi future tự quyết
  poll future con nào (theo từng "cạnh" của cây).
  select = chờ future ĐẦU TIÊN xong; join = chờ TẤT CẢ.
```

## 9.7. Spawn — tạo task mới

```
async fn server(socket: TcpListener) -> Result<()> {
    while let Some(stream) = socket.accept().await? {
        spawn(handle_client(stream));   // ← biến future thành TASK riêng
    }                                   //    chạy nền, đa luồng nếu executor cho phép
}
```

```
TIẾN HÓA xử lý kết nối:
1. await tuần tự:  handle_client(stream).await  → 1 lần 1 kết nối (KHÔNG concurrent)
2. manual executor: giữ Vec<future>, poll thủ công → concurrent nhưng busy-loop, lộn xộn
3. spawn:          mỗi client = 1 task → executor tự multiplex, đa luồng nếu Send

spawn(future) ≈ spawn thread, NHƯNG:
  - task vẫn phụ thuộc executor poll (executor chết → task chết)
```

```
BLOCKING IN ASYNC CODE (cạm bẫy lớn):
  Gọi code đồng bộ (std::sync::sleep, vòng lặp chặt) trong async
  → chiếm thread executor → task khác KHÔNG chạy được → delay lớn.
  Quy tắc: KHÔNG future nào nên chạy quá ~1ms mà không trả Poll::Pending.
  Giải pháp: chuyển sang async, hoặc chạy trên thread riêng + channel.
  (Một số executor đa luồng dùng work-stealing để giảm thiểu.)
```

---

# TỔNG KẾT TOÀN SÁCH — Sợi chỉ xuyên suốt

```
                    "BẢN CHẤT" CỦA RUST QUA 8 CHƯƠNG
  ┌────────────────────────────────────────────────────────────────┐
  │ Ch.2 Foundations  → MỌI THỨ là value/place/pointer; ownership &   │
  │                     borrow checker = kiểm tra "flow" tương thích  │
  │ Ch.3 Types        → type = cách diễn giải bit; generic=copy-paste │
  │                     (mono) vs dyn=vtable; coherence giữ 1 lựa chọn│
  │ Ch.4 Interfaces   → API = hợp đồng; mã hóa luật vào TYPE để       │
  │                     "trạng thái sai không biểu diễn được"         │
  │ Ch.5 Errors       → enumerate (phân biệt) vs erase (opaque);      │
  │                     ? = unwrap-or-return + From                   │
  │ Ch.6 Structure    → crate=đơn vị compile; feature phải additive;  │
  │                     SemVer + MSRV + minimal versions              │
  │ Ch.7 Testing      → cfg(test) lộ nội bộ; Miri/Loom bắt UB/race;   │
  │                     black_box chống compiler tối ưu benchmark     │
  │ Ch.8 Macros       → compiler viết code hộ; declarative=hygienic   │
  │                     (token→AST), procedural=chạy code (TokenStream)│
  │ Ch.9 Async        → Future=Poll; async/await→generator→state      │
  │                     machine; Pin chống move self-ref; Waker đánh   │
  │                     thức executor; spawn=task                     │
  └────────────────────────────────────────────────────────────────┘
```

**Triết lý chung của Jon Gjengset xuyên suốt:** Rust không có "ma thuật" — mọi tính năng cao cấp (lifetime, trait object, async, macro) đều quy về **các cơ chế cụ thể, có thể giải thích được**. Hiểu cơ chế bên dưới (bộ nhớ, monomorphization, state machine, vtable) thì những lỗi compiler khó hiểu trở nên hiển nhiên, và bạn viết được code vừa an toàn vừa hiệu quả.
