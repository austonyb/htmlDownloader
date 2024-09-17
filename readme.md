## HTML Downloader
### Build Instructions:

1. **Build the Executable**:
    Navigate to your project directory and run the following command to build the Go file into an executable:
    ```sh
    go build -o html_downloader main.go
    ```

2. **Run the Executable**: You can run the executable directly from the terminal:
    ```sh
    ./html_downloader -url <URL> -name <Directory Name> [--keep-js]
    ```

3. **Install Globally (Optional)**: To make the executable accessible from anywhere in your terminal, move it to a directory included in your system's PATH:
    ```sh
    sudo mv html_downloader /usr/local/bin/
    ```

Now you can run `html_downloader` from any directory in your terminal.
```sh
Replace `<URL>` with the URL of the HTML page you want to download, `<Directory Name>` with the name of the directory where you want to save the files, and optionally include `--keep-js` if you want to retain the `<script>` tags.
```