# JTL Wawi - Email sender on oversell

![WawiER Banner](assets/WawiER-Banner.png)

A small Go project for processing orders.  
It checks if items are oversold and notifies customers via email.

## Features
- Gets notified on order via JTL Wawi workflow web request
- Fetch order data via REST API
- Check stock levels 
- Send email notifications informing the customer if oversell happened