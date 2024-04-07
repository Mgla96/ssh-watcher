# ssh-watcher

[![Run Tests](https://github.com/Mgla96/ssh-watcher/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/Mgla96/ssh-watcher/actions/workflows/main.yml)

**SSH Watcher** monitors SSH logs and sends alerts to Slack for quick incident response. It is open for extension to send alerts to other notification services beyond Slack.

This project is still a work in progress.

## Installation

### Download Binary


Currently, `ssh-watcher` supports only the `AMD64` architecture, although binaries for other architectures can be built in the future. The binary to run ssh-watcher can be found under the `Assets` section of a release. To programmatically download the binary, follow the instructions below:

Downloading latest version of ssh-watcher. To download the latest version, you need `curl` and `jq` installed on your system.

```bash
curl -L $(curl -s https://api.github.com/repos/Mgla96/ssh-watcher/releases/latest | \
jq -r '.assets[] | select(.name == "ssh-watcher-linux-amd64") | .browser_download_url') \
-o ssh-watcher-linux && \
chmod +x ssh-watcher-linux && \
sudo mv ssh-watcher-linux /usr/local/bin/ssh-watcher
```

Alternatively, you can install a specific release with the following command and replacing `v0.1.0` with the version of
`ssh-watcher` you would like to install.

```bash
curl -L -o ssh-watcher-linux https://github.com/Mgla96/ssh-watcher/releases/download/v0.1.0/ssh-watcher-linux-amd64 && \
chmod +x ssh-watcher-linux && \
sudo mv ssh-watcher-linux /usr/local/bin/ssh-watcher
```

### Prepare the Systemd Service File

1. Copy the `ssh-watcher.service` file and populate the environment variables. Save this file in `/etc/systemd/system/`.

2. Update the `ExecStart` field in `ssh-watcher.service` file to point to the location of the ssh-watcher binary.

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
