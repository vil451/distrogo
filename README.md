# Distrogo

Distrogo is a command-line tool for managing Docker containers easily. It provides a set of commands to create, start,
enter, and list Docker containers with a user-friendly interface.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Features

- Create and manage Docker containers with ease.
- Enter running containers with an interactive shell.
- List containers with filtering options.
- Simple command-line interface for quick access.

## Installation

To install Distrogo, you need to have Go installed on your machine. Follow the instructions below:

1. Clone the repository:

   ```bash
   git clone https://github.com/vil451/distrogo.git
   ```
2. Navigate to the project directory:
   ```bash
   cd distrogo
   ```
3. Build the project:
    ```bash
    go build
    ```

4. Move the binary to a directory in your PATH (optional):

   ```bash
   mv distrogo /usr/local/bin/
   ```

### Usage

To use Distrogo, simply run the distrogo command followed by the desired sub-command. For example:

```bash
distrogo <sub-command> [options]
```

#### Commands

Create a Container
Create a new container from a specified image:

```bash
distrogo create <image> [options]
```

#### Start a Container

Start a previously created container:

```bash
distrogo start <container_name>
```

#### Enter a Container

Access a running container's shell interactively:

```bash
distrogo enter <container_name>
```

#### List Containers

List all containers with options to filter by name and status:

```bash
distrogo list [options]
```

#### Examples

Create a new container:

```bash
distrogo create ubuntu
```

#### Start a container:

```bash
distrogo start my_container
```

#### Enter a container:

```bash
distrogo enter my_container
```

#### List all containers:

```bash
distrogo list
```