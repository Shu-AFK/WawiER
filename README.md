# JTL Wawi - Email Sender on Oversell

![WawiER Banner](assets/WawiER-Banner.png)

WawiER is a small Go project for processing orders from JTL Wawi.  
It checks if items are oversold and notifies customers via email.

**⚠️ Note:** WawiER requires **administrator privileges** to run because it needs to exclude its folder from Windows Defender. This is necessary for receiving workflow web requests without interference from the antivirus.

---

## Features

- Receives order notifications via JTL Wawi workflow web requests
- Fetches order data via REST API
- Checks stock levels for oversell
- Sends email notifications to inform customers if an oversell occurred
- Configurable logging (console, file, both, or none)
- Supports SMTP email servers for notifications

---

## Configuration

WawiER uses a `config.json` file for configuration. A full guide is available in [`CONFIG.md`](CONFIG.md).  
Example configuration:

```json
{
  "ApiBaseURL": "http://127.0.0.1:5883/api/eazybusiness/",
  "ApiVersion": "1.1",
  "ExcludedOrderIdStart": [
    "EXAMPLE"
  ],
  "SmtpHost": "smtp.example.com",
  "SmtpPort": "587",
  "SmtpUsername": "user@example.com",
  "SmtpPassword": "yourpassword",
  "SmtpSenderEmail": "no-reply@example.com",
  "LogMode": "both",
  "LogFile": "WawiER.log"
}
```

- `LogMode` options: `"none"`, `"console"`, `"file"`, `"both"`
- `LogFile` is used when `LogMode` is `"file"` or `"both"` (default: `WawiER.log`)

---

## Installation

1. Make sure [Go](https://golang.org/dl/) is installed.
2. Clone the repository:

```bash
git clone https://github.com/Shu-AFK/WawiER.git
cd WawiER
```

3. Build the project:

```bash
go build -o WawiER.exe github.com/Shu-AFK/WawiER
```

4. Place your `config.json` in the same folder or specify the path using the `-config` flag.

**⚠️ Reminder:** Run the executable as an **administrator** to allow Windows Defender exclusions.

---

## Usage

Run the compiled executable with administrator privileges:

```bash
./WawiER.exe
```

Optional flags:

- `-config <path>` → Path to a custom config file (defaults to `config.json`)
- `-log <mode>` → Override logging mode (`none`, `console`, `file`, `both`)
- `-logfile <path>` → Override log file path if using file or both mode

---

## Logging

Logs can be written to:

- Console only
- File only
- Both console and file
- Disabled (none)

By default, logs go to the console.

---

## Contributing

Contributions are welcome! Please submit pull requests or issues on [GitHub](https://github.com/Shu-AFK/WawiER).

---

## License

This project is licensed under the [Apache-2.0](LICENSE.md) License.
