## Introduction
gopwd is a (cross platform 😁) command-line password manager written in Golang. It provides a secure, convenient, and portable way to store and manage your passwords.

## Installation Guide

To install and use `gopwd`, follow the steps below:

### Step 1: Download the Binary

1. Go to the [GitHub releases page](https://github.com/pvwnthem/gopwd/releases) for `gopwd`.
2. Locate the release version you want to install. For example, v1.0.2.
3. Download the binary file suitable for your operating system. Choose the appropriate file based on your system architecture (e.g., 32-bit or 64-bit) and operating system (Linux, macOS, or Windows).

### Step 2: Extract the Binary (if necessary)
if you used a system package like a .deb or .rpm this is not necessary, the package manager will take care of this for you.
If you downloaded a compressed archive file (e.g., `.tar.gz` or `.zip`), extract its contents to a directory of your choice.

### Step 3: Add the Binary to Your System's Path

if you used a system package like a .deb or .rpm this is not necessary, the package manager will take care of this for you.
To use `gopwd` from anywhere on your system, you need to add the binary file to your system's executable path.

**Linux and macOS:**

1. Open a terminal.
2. Navigate to the directory where you extracted the binary file.
3. Use the following command to move the binary to a directory in your system's PATH (e.g., `/usr/local/bin`):

   ```bash
   sudo mv gopwd /usr/bin/
   ```

   Enter your password when prompted.

**Windows:**

1. Open File Explorer and navigate to the directory where you extracted the binary file.
2. Right-click on the `gopwd.exe` file and select "Cut" from the context menu.
3. Press `Win + X` and select "System" from the menu.
4. In the System window, click on "Advanced system settings" on the left.
5. Click on the "Environment Variables" button.
6. In the "System variables" section, scroll down and select the "Path" variable.
7. Click on the "Edit" button.
8. In the "Edit Environment Variable" window, click on the "New" button.
9. Paste the path to the directory where you extracted the `gopwd.exe` file.
10. Click "OK" on all windows to save the changes.

### Step 4: Verify the Installation

To verify that `gopwd` is installed correctly, open a new terminal or command prompt and run the following command:

```bash
gopwd help
```

If the installation was successful, you should see commands `gopwd` printed on the screen.

Congratulations! You have successfully installed `gopwd` on your system. You can now start using it to manage your passwords.

## Initializing a Vault
To initialize a vault, use the following command:

```
gopwd vault init [-p <path>] [-n <name>]
```

- `-p` (optional): Specify the path to the vault (default: `$HOME/.gopwd`).
- `-n` (optional): Specify the name of the vault (default: "vault").

**Example:**

```
gopwd vault init -p /home/'your username'/vaults -n main
```

This command initializes a vault at `/home/'your username'/vaults/main`.

### Additional Information on Initializing a Vault
If you use the default options, all other commands can be executed without specifying flags or a config file. They will automatically use the default vault location. However, if you have used a custom location, it is recommended to set up the path and name of your vault in a config file. Without a config file, you will need to specify the path and name flags for every command you run (`-p [path]` and `-n [name]`).

To set up a config file, you have two options:

1. Run the following command:

   ```
   gopwd config init -p "path to the directory where your vault is (not the full path)" -n "name of your vault (name of the actual vault folder)"
   ```

2. Run the following command to initialize the config file:

   ```
   gopwd config init
   ```

   Then, manually set the config options using:

   ```
   gopwd config set <field> <value>
   ```

   The `<field>` can be either "path" or "name".

You can change the config settings anytime using the above commands.

## Commands

### Inserting a Password

To insert a password into your vault for a specific service, use the following command:

```
gopwd insert <service> 
```

- `<service>`: Specify the service for which the password is used. The service name can be anything, including a website, username, or any other identifier. **Important:** The service name cannot contain spaces. Use a dash (-) instead. Special characters other than a dash or underscore are also not allowed.

You can nest the service names to organize your passwords. For example, if you have multiple GitHub passwords, you can set one as `github/personal` and another as `github/work`. Nesting is optional and is useful for organizational purposes. If you don't want to use this feature, simply use the service name as the service identifier.

### Generating a Password

To generate a password for a specific service and insert it into the vault, use the following command:

```
gopwd generate <service> <length>
```

- `<service>`: Specify the service for which the password is used.
- `<length>`: Specify the length of the generated password.
- `-m`, `--memorable` (optional): Generate a more memorable password which includes full words.
- `--no-symbols` (optional): Generate a password without symbols.

Example service names and their representation in the vault:

```
github
- vault
 | - github
    | - password

github/personal
- vault
 | - github
    | - personal
       | - password

github/work and github/personal in the same vault
- vault
 | - github
    | - personal
       | - password
    | - work
       | - password
```

**Note:** The password file is the encrypted version of your password, which will be automatically retrieved when you run the command `gopwd show <service>`.

### Showing a Password

To view the password for a specific service, use the following command:

```
gopwd show <service>
```

- `<service>`: Specify the service for which you want to view the password.
- `-c`, `--copy` (optional): Copy the password to your clipboard and don't show it on stdout. 
- `-l`, `--line` `<line>` (optional): Print or copy only a certain line of the password file. This is useful if you have metadata in your files that you dont want copied or shown (or the other way around). When this flag is not provided, the whole file is copied or printed.

### Removing a Password

To remove a password and its associated folder for a specific service, use the following command:

```
gopwd rm <service>
```

- `<service>`: Specify the service for which you want to remove the password.

### Editing a Password

To edit a password or add metadata such as an email or username, use the following command:

```
gopwd edit <service>
```

## Copy a Password to Another Service

To copy a password from one service to another, use the following command:

```
gopwd cp <service> <new service>
```

### Removing a Vault

To remove a vault, use the following command:

```
gopwd vault rm [-p <path>] [-n <name>]
```