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

### Example Email Notification

![Example Email](assets/example_email/example_email.png)

This is an example of the email notification sent to customers if oversell occurs.

---

## Configuring JTL Wawi Workflow

To integrate WawiER with JTL Wawi, create a workflow that triggers when an order is created:

1. Open JTL Wawi and navigate to **Workflow Designer**.
2. Create a new workflow for **"Auftrag Erstellt"** (Order Created).
3. Add a **Webhook / HTTP Request** action with the following settings:

**URL:**
```
http://127.0.0.1:8080/api/neuerAuftrag
```

**Headers:**
```
Content-Type: application/json
Authentication: Bearer c4b55569-3d82-44a0-b9e1-79a06b79eaf1
```

**Body:**
```json
{
    "orderId": "{{ Vorgang.Stammdaten.Auftragsnummer }}"
}
```

4. Save and activate the workflow. WawiER will now receive new orders in real time.

> **Tip:** Make sure WawiER is running as administrator so it can receive the requests without interference from Windows Defender.

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
