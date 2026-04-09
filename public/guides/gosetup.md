# Setting up Golang

Setting up Go is really easy for all platforms and should only take two minutes.

---

## Downloading

Go to the [Go download page](https://go.dev/dl/) and download the appropriate file for your OS:

- Linux: .tar.gz
- Windows: .msi
- Mac: .pkg

---

## OS-Specific installation

### Linux

1.  Remove any previous Go installations:

        ```bash
        rm -rf /usr/local/go
        ```

2.  Extract the downloaded archive to /usr/local (make sure the filename/path is correct!):

        ```bash
        tar -C /usr/local -xzf go1.26.2.linux-amd64.tar.gz
        ```

3.  Add /usr/local/go/bin to the PATH environment variable.

        ```bash
        export PATH=$PATH:/usr/local/go/bin
        ```

    You might have to log out and back in for the PATH to take effect, or you can try running `source $HOME/.profile`

4.  Verify the installation:

        ```bash
        go version
        ```
---

### Windows

1.  Run the downloaded .msi installer and follow the prompts.

2.  After installation, open Command Prompt and verify the installation:

        ```cmd
        go version
        ```

---

### Mac

1.  Open the package file you downloaded and follow the prompts to install Go
2.  After installation, open Terminal and verify the installation:

        ```bash
        go version
        ```
