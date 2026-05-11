# Catapult

A simple, lightweight, speedy bootstrapper to skip downloading the Meowcoin blockchain from slow peers.

## Why I created Catapult

Meowcoin's blockchain has been growing for years and is almost 4GB. The problem is that the number of high-speed peers available to sync from is relatively small, which can cause a "Sync Stall." In some cases, it can take several days just to reach 100%.

If we want Meowcoin to grow, we need a way to get people on the network now, not next week. Most people simply aren't going to wait that long just to open a crypto wallet. Catapult exists to remove that barrier, providing a high-speed "shortcut" so new users can start participating in the ecosystem immediately.

## 📋 Prerequisites

- Ensure you have at least 10GB of free space on your device.
- You may need to run as administrator as Catapult interacts with your AppData folder.
- Make sure your wallet is fully closed.

## 🛠️ How to Use Catapult

### **Step 1: Download**

Go to the [Releases](https://github.com/Dem0-cs/Catapult/releases) page and download the binary for your operating system:

- **Windows:** `catapult-windows-x86_64.exe`

### **Step 2: Preparation**

**CRITICAL:** Ensure your Meowcoin Core wallet is **completely closed**.
Catapult needs to delete old blockchain data to make room for the fresh bootstrap. If the wallet is open, your OS will most likely lock those files and the process will fail.

### **Step 3: Run the Tool**

#### **🪟 Windows**

1. Double-click `catapult-windows-x86_64.exe`.
2. If a "Windows protected your PC" message appears, click **More Info** and then **Run Anyway**.
3. Follow the instructions within Catapult.

#### **🐧 Linux**

1. Move to the directory where you downloaded the binary `cd /path/to/catapult`
2. Give Catapult permissions `chmod +x catapult-(Version)-linux-x86_64`
3. Run Catapult using `./catapult-(Version)-linux-x86_64`

## 💖 Donating

If you want to donate to this project feel free to send $MEWC to this address
MATQBDBMdHCf6cjuH5q8zWAj8zNgNoqWYN
