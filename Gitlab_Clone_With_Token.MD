
# Hướng dẫn Clone Mã Nguồn GitLab sử dụng Access Token (Git CLI)

## 🔐 1. Tạo Project Access Token trên GitLab

**Thực hiện với quyền Owner hoặc cá nhân tự tạo token:**

1. Vào GitLab tại đường dẫn project:  
   `https://<gitlab-instance>/developmentvn/<project>`
2. Chọn **Settings** > **Access Tokens**
3. Điền các thông tin:
   - **Name:** Tên tùy chọn (ví dụ: `clone-token`)
   - **Expires at:** Ngày hết hạn (tuỳ chọn)
   - **Scopes:** Chọn `read_repository`
4. Nhấn **Create project access token**
5. Sao chép và lưu lại token (chỉ hiển thị một lần duy nhất)

---

## 💻 2. Clone project bằng dòng lệnh Git

### Cấu trúc cú pháp:

```
git -c http.extraHeader="Authorization: Bearer glpat-abc123xyz789" clone https://<gitlab-domain>/<group>/<project>.git
```

### Ví dụ:

- Token: `glpat-abc123xyz789`
- GitLab domain: `gitlab.company.com`
- Project: `examples/example-project`

```bash
git -c http.extraHeader="Authorization: Bearer glpat-abc123xyz789" clone https://gitlab.company.com/examples/example-project.git
```

> **Không cần nhập username/password**

---

## 📁 3. Sau khi clone thành công

```bash
cd example-project
code .   # mở VS Code tại thư mục
```

> Nếu chưa có `code` trong terminal: Mở VS Code > Command Palette > `Shell Command: Install 'code' command in PATH`

---

## 📎 Tips bảo mật

- Không commit hoặc chia sẻ token công khai
- Có thể thu hồi token bất kỳ lúc nào tại: `Settings > Access Tokens`
- Xoá lịch sử shell chứa token:

```bash
history -d $(history | grep "git clone" | tail -1 | awk '{print $1}')
```

---

**Tài liệu này dành cho thành viên dự án sử dụng GitLab Self-managed và VS Code hoặc terminal.**
