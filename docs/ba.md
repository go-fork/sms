## **Business Analysis: Go Module sms**

### **1. Tổng quan**
#### **Mục tiêu**
Xây dựng một **Go module** dưới dạng **service component** (thư viện) để gửi **SMS** và **Voice Call** tích hợp vào các hệ thống Go. Module này giải quyết vấn đề kết nối với nhiều nhà cung cấp dịch vụ SMS khác nhau (Twilio, eSMS, SpeedSMS, v.v.) thông qua một API thống nhất.

Các tính năng chính:
- **Message Model** chuẩn hóa định dạng tin nhắn gửi đi
- **Template System** cho phép binding dữ liệu động vào nội dung tin nhắn
- **Provider Abstraction** giúp dễ dàng chuyển đổi và quản lý nhiều nhà cung cấp dịch vụ
- **Configuration Management** tự động xử lý việc đọc và xác thực cấu hình

Người dùng chỉ cần:
1. Cung cấp đường dẫn đến file cấu hình (`configFile`)
2. Khởi tạo module thông qua `sms.NewModule(configFile)`
3. Đăng ký các provider cần thiết qua `module.AddProvider(name, provider)`
4. Gọi các phương thức `SendSMS` hoặc `SendVoiceCall` trực tiếp từ module

#### **Giá trị mang lại**
- **Tính linh hoạt**: Người dùng tự chọn và đăng ký các provider cần thiết.
- **Tính tái sử dụng**: Tích hợp dễ dàng vào hệ thống Go (web, microservice, CLI).
- **Tính chuẩn hóa**: Message model (`From`, `To`, `By`) định nghĩa thông tin gửi.
- **Tính mở rộng**: Dễ dàng phát triển và tích hợp provider mới.
- **Hiệu suất**: Tối ưu gửi tin nhắn, retry logic xử lý lỗi.
- **Bảo mật**: Quản lý thông tin xác thực an toàn.
- **Đơn giản**: Giao diện module rõ ràng, trực tiếp, không qua lớp client trung gian.
- **Cấu hình dễ dàng**: Người dùng chỉ cung cấp `configFile`, module tự xử lý cấu hình.

#### **Đối tượng người dùng**
- Nhà phát triển Go cần gửi SMS/Voice Call cho khách hàng.
- Doanh nghiệp sử dụng SMS/Voice Call cho thông báo, tiếp thị, xác thực (bao gồm OTP), hoặc giao dịch (ngân hàng, thương mại điện tử, logistics).
- Hệ thống cần chuẩn hóa thông tin gửi và binding nội dung động.

#### **Trường hợp sử dụng**
- Gửi thông báo (giao hàng, cập nhật đơn hàng).
- Gửi mã xác thực (OTP cho đăng nhập, giao dịch).
- Gửi tin nhắn tiếp thị (khuyến mãi, sự kiện).
- Gửi lời thoại Voice Call (thông báo khẩn cấp, xác thực bằng giọng nói).

---

### **2. Yêu cầu chức năng**
#### **Yêu cầu chính**
1. **Gửi SMS**:
   - Gửi tin nhắn SMS tới khách hàng thông qua module interface.
   - Sử dụng **message model** để định nghĩa thông tin gửi (`From`, `To`, `By`).
   - Nội dung (message body) có thể được binding từ template hoặc truyền trực tiếp.
2. **Gửi Voice Call**:
   - Gửi nội dung qua cuộc gọi thoại tới khách hàng.
   - Sử dụng **message model** và nội dung lời thoại text (có thể binding từ template).
3. **Message model**:
   - Định nghĩa trong `model/message.go` với các trường:
     - `From`: Người gửi (số điện thoại, brandname, tùy provider).
     - `To`: Số điện thoại nhận tin nhắn/call.
     - `By`: Người khởi tạo (tên ứng dụng, ví dụ: "MyApp").
   - Nhúng trong `SendSMSRequest`, `SendVoiceRequest`.
   - Bao gồm tính năng template (binding dữ liệu):
     - Hàm `Render(string, data map[string]interface{}) string` trong `model/message.go`.
     - Input: Template (chuỗi) và dữ liệu (`map[string]interface{}` với key như `message`, `app_name`, `code`).
     - Output: Message body/lời thoại text (ví dụ: "Your code for MyApp is 123456").
     - Template mặc định:
       - Định nghĩa trong `model/message.go` hoặc cấu hình (`sms_template`, `voice_template`).
       - Ví dụ: SMS: "Your message is {message}"; Voice: "Your message is {message}".
     - Người dùng có thể override template qua `SendSMSRequest.Template` hoặc `SendVoiceRequest.Template`.
4. **Hỗ trợ nhiều provider**:
   - Tích hợp:
     - **Việt Nam**: eSMS, SpeedSMS, Stringee, Fibo, VoiceFPT, 1SS.
     - **Quốc tế**: Twilio, Dexatel, Plivo, Exotel, Gupshup, Msg91, 2Factor.
   - Chọn provider qua `default_provider` trong cấu hình.
5. **Cấu hình linh hoạt**:
   - Người dùng cung cấp `configFile` (đường dẫn đến file cấu hình, ví dụ: `/path/to/config.yaml`).
   - Module sử dụng Viper (`github.com/spf13/viper`) để load cấu hình từ `configFile`.
   - **Cấu hình chung** (top level): `default_provider`, `http_timeout`, `retry_attempts`, `retry_delay`, `sms_template`, `voice_template`.
   - **Cấu hình provider** (key `providers`): API key, token, endpoint.
   - Mỗi adapter tự load/validate cấu hình từ `providers.<provider>` trong `config.go`.
6. **Tính mở rộng**:
   - Adapter pattern, mỗi adapter có `config.go`.
   - Dễ thêm provider mới.
7. **Hiệu suất và bảo mật**:
   - Retry logic với exponential backoff.
   - Quản lý thông tin xác thực an toàn.
   - Gửi tin nhắn/call trong 5-10 giây.

#### **Yêu cầu phi chức năng**
- **Hiệu suất**: Gửi SMS/Voice Call trong 5-10 giây.
- **Tính tương thích**: Go 1.18+.
- **Tính bảo trì**: Code rõ ràng, có unit test, tài liệu.
- **Kích thước**: Module nhẹ, phụ thuộc tối thiểu (Viper, Resty).
- **Không triển khai API**: Chỉ là thư viện.

---

### **3. Cấu trúc thư mục**
Cấu trúc thư mục gọn gàng, mỗi adapter là một Go module riêng biệt:

```
sms/ (module chính: github.com/go-fork/sms)
├── sms.go                    # Module container, điểm vào chính cho thư viện, quản lý providers
├── config/
│   ├── config.go             # Xử lý cấu hình chung từ configFile
│   ├── validate.go           # Validate cấu hình chung
│   └── example.yaml          # File cấu hình mẫu (YAML)
├── model/
│   ├── provider.go           # Interface Provider (SendSMS, SendVoiceCall)
│   ├── message.go            # Message model (From, To, By) và hàm Render
│   ├── request.go            # Structs cho request (SendSMSRequest, SendVoiceRequest)
│   └── response.go           # Structs cho response (SendSMSResponse, SendVoiceResponse)
├── client/
│   └── client.go             # Core logic, cung cấp client.NewClient và hàm Send*
├── retry/
│   └── retry.go              # Retry logic với exponential backoff
├── tests/
│   ├── client_test.go        # Unit test cho client
│   ├── sms_test.go           # Unit test cho sms.go
│   ├── config_test.go        # Unit test cho cấu hình
│   └── message_test.go       # Unit test cho message và template
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md     # Template cho báo cáo lỗi
│   │   └── feature_request.md # Template cho yêu cầu tính năng
│   └── PULL_REQUEST_TEMPLATE.md # Template cho Pull Request
├── examples/
│   ├── simple/
│   │   └── main.go           # Ví dụ cơ bản về cách sử dụng
│   └── template/
│       └── main.go           # Ví dụ sử dụng template
├── LICENSE                   # File license (MIT hoặc Apache-2.0)
├── CHANGELOG.md              # Lịch sử thay đổi của module
├── CONTRIBUTING.md           # Hướng dẫn đóng góp vào dự án
├── go.mod                    # Go module definition (bao gồm Viper)
├── go.sum                    # Go module dependencies
├── .gitignore                # Danh sách file bỏ qua trong Git
└── README.md                 # Hướng dẫn cài đặt, sử dụng, tích hợp

adapters/ (mỗi thư mục là một Go module riêng)
├── twilio/ (module: github.com/go-fork/sms/adapters/twilio)
│   ├── config.go             # Load và validate cấu hình Twilio từ Viper
│   ├── adapter.go            # Logic gọi API Twilio
│   ├── adapter_test.go       # Unit test cho adapter Twilio
│   ├── LICENSE               # File license (giống module chính)
│   ├── README.md             # Hướng dẫn sử dụng adapter Twilio
│   ├── .gitignore            # File gitignore cho module adapter
│   ├── go.mod                # Go module định nghĩa cho adapter Twilio
│   └── go.sum                # Go module dependencies cho Twilio
├── esms/ (module: github.com/go-fork/sms/adapters/esms)
│   ├── config.go             # Load và validate cấu hình eSMS
│   ├── adapter.go            # Logic gọi API eSMS
│   ├── adapter_test.go       # Unit test cho adapter eSMS
│   ├── LICENSE               # File license (giống module chính)
│   ├── README.md             # Hướng dẫn sử dụng adapter eSMS
│   ├── .gitignore            # File gitignore cho module adapter
│   ├── go.mod                # Go module định nghĩa cho adapter eSMS
│   └── go.sum                # Go module dependencies cho eSMS
├── speedsms/ (module: github.com/go-fork/sms/adapters/speedsms)
│   ├── config.go             # Load và validate cấu hình SpeedSMS
│   ├── adapter.go            # Logic gọi API SpeedSMS
│   ├── adapter_test.go       # Unit test cho adapter SpeedSMS
│   ├── LICENSE               # File license (giống module chính)
│   ├── README.md             # Hướng dẫn sử dụng adapter SpeedSMS
│   ├── .gitignore            # File gitignore cho module adapter
│   ├── go.mod                # Go module định nghĩa cho adapter SpeedSMS
│   └── go.sum                # Go module dependencies cho SpeedSMS
```

#### **Giải thích**
- **Module chính (github.com/go-fork/sms)**: Chứa logic cốt lõi, mô hình message và client.
  - `sms.go` làm container, quản lý các providers thông qua `AddProvider` và `SwitchProvider`.
- **Các modules adapter (github.com/go-fork/sms/adapters/[name])**: Mỗi adapter là một Go module riêng biệt.
  - Mỗi adapter đều phụ thuộc vào module chính để truy cập vào interface Provider và models.
  - Tách riêng giúp quản lý phiên bản độc lập cho từng adapter.
  - Người dùng chỉ cần import các adapter cần thiết.

---

### **4. Development Guidelines**
#### **4.1 Môi trường phát triển**
- **Go version**:
  - Go 1.18+ (yêu cầu tối thiểu)
  - Go 1.20+ (khuyến nghị)

- **Thư viện chính**:
  | Thư viện | Phiên bản tối thiểu | Mục đích |
  |----------|---------------------|----------|
  | github.com/spf13/viper | v1.15.0 | Quản lý cấu hình |
  | github.com/go-resty/resty/v2 | v2.7.0 | HTTP client |
  | github.com/stretchr/testify | v1.8.0 | Testing framework |
  | golang.org/x/time | latest | Rate limiting & backoff |

- **Công cụ phát triển**:
  - VS Code với Go extension
  - GoLand IDE
  - golangci-lint (v1.50.0+)

#### **4.2 Code Rules Guidelines**
##### **Cấu trúc mã nguồn**
- **Package Structure**: Mỗi package nên tập trung vào một chức năng cụ thể.
- **File Organization**: 
  - Mỗi file không quá 500 dòng
  - Mỗi function không quá 50 dòng
  - Mỗi package có một file `doc.go` mô tả chức năng

##### **Dependency Injection**
- **Constructor Pattern**: Sử dụng constructor để khởi tạo struct với dependencies.
```go
// Sử dụng pattern này
func NewClient(config *Config) (*Client, error) {
    // Validate config
    if err := config.Validate(); err != nil {
        return nil, err
    }
    return &Client{
        config: config,
        providers: make(map[string]Provider),
    }, nil
}
```
- **Avoid Global State**: Không sử dụng biến global, singletons không có lý do chính đáng.
- **Testability**: Thiết kế mã nguồn để dễ dàng mock dependencies trong tests.

##### **Error Handling**
- **Error Types**: Sử dụng error types cụ thể cho từng loại lỗi.
```go
type ProviderNotFoundError struct {
    Name string
}
func (e ProviderNotFoundError) Error() string {
    return fmt.Sprintf("provider %s not found", e.Name)
}
```
- **Context**: Luôn truyền context.Context trong các API public.
- **Wrapping Errors**: Sử dụng `fmt.Errorf("failed to do X: %w", err)` để wrap errors.

##### **Naming Conventions**
- **Package Names**: Ngắn gọn, đơn số (ví dụ: `model`, không phải `models`).
- **Variable Names**: CamelCase, tránh viết tắt không phổ biến.
- **Interface Names**: Tên interface cho hành động nên kết thúc bằng "er" (ví dụ: `Provider`).
- **Constant Groups**: Sử dụng `const` groups khi liên quan đến nhau.

##### **Comments & Documentation**
- **Package Comments**: Mỗi package nên có comment mô tả ở đầu file `doc.go`.
- **Function Comments**: Các public functions/methods nên có comment theo chuẩn GoDoc.
- **Example Code**: Thêm examples cho các API quan trọng.

##### **Testing**
- **Unit Tests**: Coverage tối thiểu 80% cho code mới.
- **Table-Driven Tests**: Sử dụng table-driven tests khi testing nhiều cases.
- **Mock Data**: Sử dụng mock data cho tất cả external calls.
- **Test Helpers**: Tách logic test helper vào các functions riêng.

##### **Performance**
- **Pooling**: Sử dụng connection pooling cho HTTP clients.
- **Caching**: Cache kết quả khi thích hợp (ví dụ: provider instances).
- **Memory Allocation**: Tránh memory allocation không cần thiết, đặc biệt trong hot paths.

#### **4.3 Code Review Guidelines**
- **Pull Requests**: Mỗi PR không quá 500 dòng code thay đổi.
- **Review Checklist**:
  - Code đúng theo style guide
  - Tests đầy đủ và pass
  - Documentation đã cập nhật
  - Error handling đầy đủ
  - Không có hardcoded values
  - Không có security vulnerabilities

#### **4.4 Git Workflow**
- **Branch Naming**: `feature/tên-tính-năng`, `fix/tên-lỗi`, `docs/tên-tài-liệu`
- **Commit Messages**: Tuân theo format [Conventional Commits](https://www.conventionalcommits.org/)
- **Versioning**: Tuân theo [Semantic Versioning](https://semver.org/)

---

### **5. Các bước phát triển**
1. **Thiết kế cấu trúc module**:
   - Xây dựng interface `Provider` (`model/provider.go`).
   - Định nghĩa `client/client.go` với các phương thức cơ bản.
   - Tạo `model/message.go` cho message model và hàm `Render`.
   - Xây dựng `sms.go` làm container module với quản lý providers thông qua `AddProvider` và `SwitchProvider`.
2. **Tích hợp Viper cho cấu hình**:
   - Xây dựng `config/config.go` và `config/validate.go` để load/validate cấu hình từ `configFile`.
   - Tạo file mẫu `config/example.yaml`.
3. **Triển khai adapters dưới dạng Go modules riêng**:
   - Tạo các Go modules riêng cho từng adapter:
     - `github.com/go-fork/sms/adapters/twilio`
     - `github.com/go-fork/sms/adapters/esms`
     - `github.com/go-fork/sms/adapters/speedsms`
   - Mỗi adapter có `config.go` và file cài đặt chính (vd: `adapter.go`).
   - Mỗi adapter tự quản lý go.mod và dependencies riêng.
4. **Thêm retry logic**:
   - Xây dựng `retry/retry.go` với exponential backoff.
5. **Viết unit test**:
   - Test client, sms.go, cấu hình, và message (bao gồm template) trong `tests/`.
   - Mỗi adapter module có test riêng.
6. **Tối ưu và tài liệu**:
   - Tối ưu hiệu suất (HTTP client, timeout).
   - Viết README, ví dụ tích hợp, và tài liệu template.
7. **Publish các module**:
   - Đưa lên GitHub repository `github.com/go-fork/sms` (module chính).
   - Đưa lên GitHub các modules adapter `github.com/go-fork/sms/adapters/[name]`.
   - Thiết lập issue tracker và PR template.

---

### **6. Hướng dẫn sử dụng**
#### **Tính năng chính**
- Gửi SMS/Voice Call tới khách hàng với **message model** (`From`, `To`, `By`).
- Hỗ trợ binding dữ liệu template thành chuỗi (tính năng nhỏ trong Message).
- Hỗ trợ nhiều provider (Twilio, eSMS, SpeedSMS).
- Retry tự động khi request thất bại.
- Cấu hình qua `configFile`, module tự xử lý Viper.

#### **Cấu hình mẫu**
File `config.yaml` (đường dẫn đến file cấu hình):
```yaml
# Cấu hình chung
default_provider: twilio
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms
sms_template: "Your message from {app_name}: {message}"
voice_template: "Your message from {app_name} is {message}"

# Cấu hình provider
providers:
  twilio:
    account_sid: ACxxxxxxxxxxxxxxxx
    auth_token: xxxxxxxxxxxxxxxxxxx
    from_number: +1234567890
  esms:
    api_key: xxxxxxxxxxxxxxxx
    secret: xxxxxxxxxxxxxxxxx
  speedsms:
    token: xxxxxxxxxxxxxxxx
```

#### **Cách sử dụng**
1. **Khởi tạo module**:
   - Người dùng chỉ cần cung cấp `configFile`, module tự động xử lý mọi thứ.
   ```go
   import (
       "github.com/go-fork/sms"
   )
   configFile := "./config.yaml" // Đường dẫn đến file cấu hình
   module, err := sms.NewModule(configFile)
   if err != nil {
       panic(err)
   }
   ```

2. **Đăng ký và chuyển đổi providers (nếu cần)**:
   ```go
   import (
       "github.com/go-fork/sms"
       "github.com/go-fork/sms/adapters/twilio"
       "github.com/go-fork/sms/adapters/custom"
   )
   
   // Khởi tạo module với file cấu hình
   configFile := "./config.yaml"
   module, err := sms.NewModule(configFile)
   if err != nil {
       panic(err)
   }
   
   // Khởi tạo provider với cùng file cấu hình
   twilioProvider, err := twilio.NewProvider(configFile)
   if err != nil {
       panic(err)
   }
   
   // Thêm provider vào module (Name() được gọi tự động bên trong)
   err = module.AddProvider(twilioProvider)
   if err != nil {
       panic(err)
   }
   
   // Chuyển sang sử dụng provider khác theo tên
   // "twilio" là giá trị trả về bởi twilioProvider.Name()
   err = module.SwitchProvider("twilio")
   if err != nil {
       panic(err)
   }
   ```

3. **Gửi SMS với message model**:
   ```go
   import (
       "context"
       "fmt"
       "github.com/go-fork/sms/model"
   )

   ctx := context.Background()
   req := model.SendSMSRequest{
       Message: model.Message{
           From: "+1234567890",
           To:   "+84912345678",
           By:   "MyApp",
       },
       Data: map[string]interface{}{
           "app_name": "MyApp",
           "message":  "Your OTP is 123456",
       },
       Template: "", // Nếu rỗng, dùng sms_template từ config.yaml
   }
   resp, err := module.SendSMS(ctx, req)
   if err != nil {
       panic(err)
   }
   fmt.Printf("SMS sent: %+v\n", resp)
   ```

#### **Message và Template**
- **Message model** (`model/message.go`):
  - Định nghĩa `From`, `To`, `By` để chuẩn hóa thông tin gửi.
  - Bao gồm hàm `Render(string, data map[string]interface{}) string` để binding dữ liệu template.
  - Nhúng trong `SendSMSRequest`, `SendVoiceRequest`.
- **Template** (tính năng nhỏ trong Message):
  - Hàm `Render` binding `template data` thành chuỗi:
    - Input: Template (chuỗi) và dữ liệu.
    - Output: Message body/lời thoại text.
  - Template mặc định trong `model/message.go`:
    - SMS: "Your message is {message}".
    - Voice: "Your message is {message}".
  - Tùy chỉnh:
    - Qua file cấu hình (`sms_template`, `voice_template`) trong `configFile`.
    - Qua `SendSMSRequest.Template` hoặc `SendVoiceRequest.Template`.

------

### **7. Kế hoạch triển khai**
#### **Triển khai module**
- **Publish lên GitHub**:
  - Module chính: `github.com/go-fork/sms`
  - Modules adapter: `github.com/go-fork/sms/adapters/[name]`
  - README với hướng dẫn cài đặt, cung cấp `configFile`, sử dụng.
  - File `config/example.yaml` với template mẫu.
- **Versioning**:
  - Semantic versioning (v1.0.0) cho mỗi module.
  - Module chính và các adapter có thể có phiên bản độc lập.
  - Tag release trên GitHub.
- **CI/CD**:
  - GitHub Actions để chạy unit test, lint code cho mỗi module.
  - Auto-publish khi merge vào main.

#### **Triển khai trong hệ thống**
- **Cài đặt**:
  ```go
  go get github.com/go-fork/sms
  // Chỉ cài đặt các adapter cần thiết
  go get github.com/go-fork/sms/adapters/twilio
  go get github.com/go-fork/sms/adapters/esms
  ```
- **Cấu hình**: Chuẩn bị `configFile` (ví dụ: `config.yaml`).
- **Import trong code**:
  ```go
  import (
      "github.com/go-fork/sms"
      _ "github.com/go-fork/sms/adapters/twilio" // Chỉ import adapter cần sử dụng
      _ "github.com/go-fork/sms/adapters/esms"   // Import sẽ tự đăng ký adapter
  )
  ```
- **Testing**:
  - Unit test cục bộ với mock API.
  - Test tích hợp với sandbox API.
- **Monitoring**:
  - Log `MessageID`, `Status`, `Error`.
  - Theo dõi latency và tỷ lệ thành công.

------

### **8. Rủi ro và giải pháp**
#### **Rủi ro**
1. **Khác biệt API**:
   - Provider có format request/response khác nhau.
   - **Giải pháp**: Chuẩn hóa qua `model/message.go`, `model/request.go`, `model/response.go`.
2. **Template lỗi**:
   - Sai cú pháp hoặc thiếu key trong `Data`.
   - **Giải pháp**: Validate template và `Data` trong `model/message.go`, cung cấp mặc định.
3. **Message model không hợp lệ**:
   - `From`, `To`, `By` không đúng định dạng của provider.
   - **Giải pháp**: Validate trong adapter và tài liệu rõ ràng.
4. **Lỗi cấu hình**:
   - File cấu hình trong `configFile` sai hoặc thiếu.
   - **Giải pháp**: Validate trong `config/validate.go` và `adapter/<provider>/config.go`.
5. **Chi phí provider**:
   - Twilio có chi phí cao.
   - **Giải pháp**: Tài liệu chi phí, ưu tiên eSMS, SpeedSMS.

#### **Giải pháp dự phòng**
- **Fallback provider**: Thử provider khác nếu provider chính thất bại.
- **Mock testing**: Dùng `httptest` để test không cần API thực.
- **Tài liệu**: FAQ và troubleshooting trong README.

------

### **9. Tài liệu và hỗ trợ**
#### **Tài liệu**
- **README.md**:
  - Hướng dẫn cài đặt, cung cấp `configFile`, sử dụng.
  - Ví dụ code với message model và template binding.
  - Danh sách provider, chi phí, placeholders hỗ trợ.
- **config/example.yaml**: Mẫu cấu hình với template.
- **Code comments**: Mô tả hàm, struct, message model, và template.

#### **Hỗ trợ**
- **GitHub Issues**: Báo lỗi, yêu cầu tính năng.
- **Community**: Khuyến khích PR để thêm provider.
- **Email**: Email hỗ trợ cho vấn đề phức tạp.

------

### **10. Kế hoạch mở rộng**
- **Thêm provider**: Plivo, Gupshup, Stringee.
- **Tính năng**:
  - Fallback provider.
  - Async API (channel/callback).
  - Custom logger.
  - Template đa ngôn ngữ.
- **Hiệu suất**:
  - Connection pooling cho HTTP client.
  - Batch sending SMS/Voice Call.
- **CI/CD**:
  - Auto-test với sandbox API.
  - Auto-generate tài liệu.

------

### **11. Kết luận**

Go module `sms` là một **service component** mạnh mẽ, linh hoạt để gửi **SMS** và **Voice Call** tới khách hàng. **Message model** (`From`, `To`, `By`) chuẩn hóa thông tin gửi, với tính năng nhỏ trong Message để binding `template data` thành chuỗi (`Render`). Module tự quản lý Viper, chỉ yêu cầu người dùng cung cấp đường dẫn đến file cấu hình (`configFile`). Thiết kế adapter pattern, cách sử dụng đơn giản (`sms.NewModule` + hàm `Send*`), tài liệu chi tiết, và hỗ trợ cộng đồng giúp module dễ tích hợp và bảo trì.