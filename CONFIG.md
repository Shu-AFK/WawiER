# WawiER Configuration Guide

The `config.json` file defines the runtime behavior for the WawiER application. Below is the full explanation of each field, their types, defaults, and validation rules.

---

## JSON Example

```json
{
  "ApiBaseURL": "http://127.0.0.1:5883/api/eazybusiness/",
  "ApiVersion": "1.1",
  "ExcludedOrderIdStart": [
    "EXAMPLE"
  ],
  "SmtpHost": "[ SMTP Host ]",
  "SmtpPort": "587",
  "SmtpUsername": "[ SMTP Username ]",
  "SmtpPassword": "[ SMTP Password ]",
  "SmtpSenderEmail": "[ Sender Email ]",
  "LogMode": "both",
  "LogFile": "WawiER.log"
}
```

---

## Configuration Fields

### API Settings

| Field        | Type   | Required | Description                                                                  |
|--------------|--------|----------|------------------------------------------------------------------------------|
| `ApiBaseURL` | string | Yes      | Base URL of the API server. Must include protocol (`http://` or `https://`). |
| `ApiVersion` | string | Yes      | Version of the API to use (e.g., `"1.1"`).                                   |


---

### Exclusions

| Field                  | Type             | Required | Description                                                                             |
|------------------------|------------------|----------|-----------------------------------------------------------------------------------------|
| `ExcludedOrderIdStart` | array of strings | No       | Optional list of order ID prefixes to exclude from processing. Can be empty or omitted. |

---

### SMTP Settings

| Field             | Type   | Required | Description                                   |
|-------------------|--------|----------|-----------------------------------------------|
| `SmtpHost`        | string | Yes      | Hostname of the SMTP server.                  |
| `SmtpPort`        | string | Yes      | SMTP port number as a string (e.g., `"587"`). |
| `SmtpUsername`    | string | Yes      | Username for SMTP authentication.             |
| `SmtpPassword`    | string | Yes      | Password for SMTP authentication.             |
| `SmtpSenderEmail` | string | Yes      | Email address used as the sender.             |


---

### Logging Settings

| Field     | Type   | Required | Default        | Description                                                                                               |
|-----------|--------|----------|----------------|-----------------------------------------------------------------------------------------------------------|
| `LogMode` | string | No       | `"console"`    | Determines where logs are written. Valid values: `"none"`, `"console"`, `"file"`, `"both"`.               |
| `LogFile` | string | No       | `"WawiER.log"` | Path to the log file if `LogMode` is `"file"` or `"both"`. Relative paths are resolved to absolute paths. |

**Behavior:**

- `"none"` → disables logging.
- `"console"` → logs only to console (default).
- `"file"` → logs only to the specified log file.
- `"both"` → logs to both console and file.

---

### Validation Rules

1. All required fields must be set; otherwise, WawiER will fail to start.
2. `LogMode` defaults to `"console"` if empty.
3. `LogFile` defaults to `"WawiER.log"` if missing and `LogMode` is `"file"` or `"both"`.
4. `ExcludedOrderIdStart` is optional.
5. `SmtpPort` must be a valid numeric string.

---

### Notes

- File paths in `LogFile` are converted to absolute paths.
- If the log file directory does not exist, the application may fail unless the directory is created manually or programmatically.
- Sensitive fields like `SmtpPassword` must be protected. Do not commit credentials to public repositories.
