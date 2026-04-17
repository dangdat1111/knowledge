# Trình soạn thảo (Neovim)

Tài liệu này dùng chung cho bộ ghi chú sách, với cấu hình Neovim từ repo của bạn:

- Repo: [https://github.com/dangdat1111/nvim](https://github.com/dangdat1111/nvim)

## Mục đích
- Chuẩn hóa môi trường soạn thảo khi học Go.
- Giảm khác biệt cấu hình giữa các máy.
- Dễ onboarding cho người học mới.

## Thiết lập nhanh
1. Cài Neovim (khuyến nghị bản ổn định mới).
2. Clone repo cấu hình:

```bash
git clone https://github.com/dangdat1111/nvim ~/.config/nvim
```

3. Mở Neovim để plugin tự cài lần đầu:

```bash
nvim
```

## Gợi ý plugin/tính năng cần cho Go
- LSP cho Go (gopls)
- Formatter (`gofmt`)
- Linter (`golangci-lint` nếu dùng)
- Tree-sitter cho syntax highlight
- Telescope để tìm file/ký hiệu nhanh

## Quy ước khi viet tai lieu
- Mỗi chapter cập nhật trong `ChapterXX-.../README.md`.
- Ưu tiên ví dụ chạy được và có đầu vào/đầu ra rõ ràng.
- Luôn cập nhật `Nhật ký thực hành` sau mỗi buổi học.
