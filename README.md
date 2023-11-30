<div align="center">
    <a href="https://hermes.fractum.nl">
        <img src="public/logo.png" alt="Logo" width="80" height="80">
    </a>
    <h1>Heimdall</h1>
    <p><i>Docker container state monitor with notifications for Discord</i></p>
</div>

<!-- TOC -->
<details>
    <summary>Table of Contents</summary>
    <ol>
        <li>
          <a href="#about-heimdall">About Heimdall</a>
        </li>
        <li>
          <a href="#features">Features</a>
        </li>
        <li>
            <a href="#usage">Usage</a>
            <ul>
                <li><a href="#docker-container">Docker container</a></li>
                <ul>
                    <li><a href="#using-docker-desktop-1">Using Docker Desktop</a></li>
                </ul>
                <li><a href="#standalone-application">Standalone application</a></li>
                <ul>
                    <li><a href="#using-docker-desktop-2">Using Docker Desktop</a></li> 
                </ul>
            </ul>
        </li>
        <li>
          <a href="#technologies">Technologies</a>
          <ul>
            <li><a href="#language">Language</a></li>
            <li><a href="#deployed-to">Deployed to</a></li>
            <li><a href="#ci--cd">CI/CD</a></li>
            <li><a href="#released-using">Released using</a></li>
            <li><a href="#logo-created-using">Logo created using</a></li>
          </ul>
        </li>
        <li>
          <a href="#screenshots">Screenshots</a>
        </li>
      </ol>
</details>
<!-- TOC -->

## About Heimdall
Heimdall is a monitoring application for your Docker containers. It sends you notifications through a webhook whenever the state of a container changes.
It does this by using the Docker socket to listen for events. Heimdall also provides the option to receive periodic notifications about the state of your containers, where it sends you an overview of every container's status.

## Features
- [x] Easy monitoring for Docker containers
- [x] Receive notifications through Discord webhooks
- [x] Receive periodic notifications about the state of your containers
- [ ] Status API
- [ ] Web UI
- [ ] Bugs (hopefully)

## Usage
Heimdall can be used in a couple different ways:
1. As a Docker container
2. As a standalone application

### Docker container
The easiest way to use Heimdall is by running it as a Docker container. You can do this by running the following command:
```bash
docker run -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e HEIMDALL_WEBHOOK_URL=<your-webhook-url> \
    --name heimdall \
    drfractum/heimdall:latest
```

#### Using Docker Desktop
In some cases, Docker Desktop puts the Docker socket in a different location than the standard Docker installation.
You can change the above command to the following to point to the right location:
```diff
docker run -d \
-    -v /var/run/docker.sock:/var/run/docker.sock \
+    -v ~/.docker/desktop/docker.sock:/var/run/docker.sock \
    -e HEIMDALL_WEBHOOK_URL=<your-webhook-url> \
    --name heimdall \
    drfractum/heimdall:latest
```

### Standalone application
You can also run Heimdall as a standalone application. This is useful if you want to run it on a server or your local machine.
To run Heimdall as a standalone application, you can download the latest release from the [releases page](https://github.com/RobinHeidenis/heimdall/releases).

Be sure to download the right binary for your operating system.

| Platform                                                                                                 | Binary                  |
|----------------------------------------------------------------------------------------------------------|-------------------------|
| ![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)       | heimdall_Linux_x86_64   |
| ![Raspberry Pi](https://img.shields.io/badge/-RaspberryPi-C51A4A?style=for-the-badge&logo=Raspberry-Pi)  | heimdall_Linux_armv7    |
| ![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white) | heimdall_Windows_x86_64 |
| ![macOS](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0)   | heimdall_Darwin_x86_64  |

Other binaries are available in case you have a different architecture or operating system.
Other binaries are available in case you have a different architecture or operating system.


After downloading the release, you can run it by executing the following command:
```bash
./heimdall --webhook-url=<your-webhook-url>
```

On Windows you can run it by executing the following command:
```bash
heimdall.exe --webhook-url=<your-webhook-url>
```

#### Using Docker Desktop
In some cases if you're using Docker Desktop, you might run into an error where Heimdall can't connect to the Docker socket.
By default, Docker Desktop puts the Docker socket in a different location than the standard Docker installation.
If Heimdall doesn't automatically detect your environment and the right location for the Docker socket, you'll need to create a symlink to the docker.sock file. You can do this by running the following command:
```bash
sudo ln -s ~/.docker/desktop/docker.sock /var/run/docker.sock
```

## Customisation
Heimdall can be customised by using the following environment variables:

| Long flag                 | Short flag | Environment Variable             | Default   | Required                                                                                         | Explanation                                                                                      |
|---------------------------|------------|----------------------------------|-----------|--------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------|
| `--periodic-notification` | `-n`       | `HEIMDALL_PERIODIC_NOTIFICATION` | `false`   | No                                                                                               | Enable periodic notifications                                                                    |
| `--notification-interval` | `-i`       | `HEIMDALL_NOTIFICATION_INTERVAL` | `60`      | Only if periodic notifications are enabled and you want a different value than default           | How often (in minutes) periodic notifications should be sent                                     |
| `--all-containers`        | `-a`       | `HEIMDALL_ALL_CONTAINERS`        | `false`   | Only if periodic notifications are enabled and you want periodic notifications on all containers | Enable periodic notification reporting on all containers, including stopped ones                 |
| `--retry`                 | `-r`       | `HEIMDALL_RETRY`                 | `10`      | No                                                                                               | How long Heimdall should sleep before retrying in case the Docker event stream ends unexpectedly |
| `--provider`              | `-p`       | `HEIMDALL_PROVIDER`              | `discord` | No                                                                                               | What notification provider should be used. Possible values: `discord`                            |
| `--webhook-url`           | `-w`       | `HEIMDALL_WEBHOOK_URL`           | -         | Yes                                                                                              | What URL Heimdall should use to send notifications to                                            |
| `--debug`                 | `-d`       | `HEIMDALL_DEBUG`                 | `false`   | No                                                                                               | Enable extra debug messages                                                                      |


## Technologies
Heimdall was created using Go. The CI/CD pipeline is handled by GitHub Actions and the Docker image is hosted on Docker Hub.

The program is mostly based on the [Docker SDK for Go](https://pkg.go.dev/github.com/docker/docker/client#section-readme)

It uses the following technologies:
### Language
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)

### Deployed to
[![Docker](https://img.shields.io/badge/docker-%231d63ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/r/drfractum/heimdall)

### CI/CD
[![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)](https://github.com/RobinHeidenis/heimdall/actions)

### Released using
[![GoReleaser](https://img.shields.io/badge/goreleaser-%23000.svg?style=for-the-badge&logo=data:image/png;base64,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)](https://goreleaser.com)

### Logo created using
[![Bing Image Creator](https://img.shields.io/badge/bing%20image%20creator-%230078D4.svg?style=for-the-badge&logo=microsoftbing&logoColor=white)](https://www.bing.com/images/create)


## Screenshots
### Terminal output
![logo.png](public/terminal-output.png)

### Discord notifications
#### Container started
![container-started.png](public/container-started.png)

#### Container stopped
![container-stopped.png](public/container-stopped.png)

#### Container healthy
![container-healthy.png](public/container-healthy.png)

#### Container unhealthy
![container-unhealthy.png](public/container-unhealthy.png)

#### Container errored
![container-errored.png](public/container-errored.png)

### Periodic notification
#### Running containers
![container-list.png](public/container-list.png)

#### All containers
![container-list-all.png](public/container-list-all.png)
