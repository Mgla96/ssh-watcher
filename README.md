# ssh-watcher

[![Run Tests](https://github.com/Mgla96/ssh-watcher/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/Mgla96/ssh-watcher/actions/workflows/main.yml)

**SSH Watcher** monitors SSH logs and sends alerts to Slack for quick incident response. It is open for extension to send alerts to other notification services beyond Slack.

This project is still a work in progress.

## Installation

### Download Binary


`AMD64` is the only architecture currently supported, although I'm open to building binaries for other architectures in the future. The binary to run ssh-watcher can be found under the "Assets" section of a release. To download the binary programmatically, you can run the following commands:


1. Download binary, make it executable, move binary to desired executable path location. Update `v0.1.0-alpha-2` with the version of ssh-watcher that you would like to install.

    ```bash
    curl -L -o ssh-watcher-linux https://github.com/Mgla96/ssh-watcher/releases/download/v0.1.0-alpha-2/ssh-watcher-linux-amd64 && \
    chmod +x ssh-watcher-linux && \
    sudo mv ssh-watcher-linux /usr/local/bin/ssh-watcher
    ```

### Prepare the Systemd Service File

1. Copy the `ssh-watcher.service` file and populate the environment variables. Save this file in `/etc/systemd/system/`

2. Update the `ExecStart` field in `ssh-watcher.service` file to point to the location of the ssh-watcher binary

### Start the Service

1. **Reload systemd**: to make systemd aware of the new service file.

    ```bash
    sudo systemctl daemon-reload
    ```

2. **Enable the Service**: to ensure `ssh-watcher` starts automatically at boot.

    ```bash
    sudo systemctl enable ssh-watcher.service
    ```

3. **Start the Service**: Start service immediately without rebooting.

    ```bash
    sudo systemctl start ssh-watcher.service
    ```

4. **Check the Service Status**: Verify that the service is active and running.

    ```bash
    sudo systemctl status ssh-watcher.service
    ```
