# achelous
Sendmail replacement, which queues and uploads emails to a web service. We designed it for IoT devices and containers, which might send emails for administrative purposes.

**Naming:** The Achelous (Greek: Αχελώος, Ancient Greek: Ἀχελῷος Akhelôios), also Acheloos, is a river in western Greece. [Wikipedia](https://en.wikipedia.org/wiki/Achelous_River)

The project is split into two programs:

- `spring` is responsible for mail submission and implementes the `sendmail` command line interface. The `sendmail` binaries are provided by symlinks pointing to `spring`.
- `upstream` is responsible for mail upload and success reporting. It is implemented as a daemon process and can be managed via `systemd` or `rc.d` scripts.

## Features

- Sendmail compatible command line interface, designed as a drop-in replacement. Provides the virtual `mail-transport-agent` package.
- Compatible with ISC `cron` and GNU `mail`.
- A hierarchical queue structure, which allows embedding of container queues in sub-queues on the host system.
- URL and header variables of target web service are freely configurable.
- Reports successfully processed queues to the web service. Allows the implementation of a health check.

## Installation

Achelous conflicts with all packages providing the virtual `mail-transport-agent` package, e.g. `exim4-daemon-light`. To be save, please remove them first.

```sh
# add repository
echo "deb [arch=amd64,armhf] https://coreprocess.github.io/achelous-dist/deb $(lsb_release -s -c) main" | sudo tee /etc/apt/sources.list.d/achelous.list

# add key of maintainer
sudo apt-key adv --keyserver hkp://keys.gnupg.net --recv C7192F40A7C34E5A25339476D1E482C66415ACC5

# update package index
sudo apt-get update

# install achelous
sudo apt-get install achelous
```

See [here](https://github.com/coreprocess/achelous-dist) for more.

Create and adjust the configuration files `/etc/achelous/spring.json` and `/etc/achelous/upstream.json`.

## Configuration

**Note:** Usually specific reloading procedures are not required. `spring` reads its configuration file on every call. `upstream` reads its configuration file on every queue run.

The following examples are provided for your convenience. You will find formal specifications in the [specs](./specs/) folder.

Example `/etc/achelous/spring.json`

```json
{
  "AbortOnParseErrors": false,
  "DefaultQueue": "",
  "PrettyJSON": false,
  "TriggerQueueRun": true,
  "GenerateDebugMail": {
    "OnInvalidParameters": true,
    "OnParsingErrors": true,
    "OnOtherErrors": true,
    "Message": {
      "Sender": {
        "Name": "Achelous Spring",
        "Email": ""
      },
      "Receiver": {
        "Name": "Devops",
        "Email": ""
      },
      "Subject": "ACHELOUS SPRING DEBUG MESSAGE",
      "Body": "Activity: %[1]s\nReference: %[2]s\nError: %[3]v\nData: %+[4]v"
    }
  }
}
```

Example `/etc/achelous/upstream.json`

```json
{
  "PauseBetweenRuns": {
    "PreviousRunOK": "5m",
    "PreviousRunWithErrors": "30s"
  },
  "Target": {
    "Upload": {
      "URL": "https://myservice/api/email/upload",
      "Header": {
        "Authorization": [ "Bearer example" ]
      }
    },
    "Report": {
      "URL": "https://myservice/api/email/report",
      "Header": {
        "Authorization": [ "Bearer example" ]
      }
    },
    "RetriesPerRun": {
      "Attempts": 3,
      "PauseBetweenAttempts": "5s"
    }
  }
}
```

**Note:** You might want to increase `PauseBetweenRuns.PreviousRunOK` significantly in some IoT scenarios, e.g. `1h` or more.

## Usage

Submit emails exactly as before, e.g.

```sh

# EXAMPLE 1
sendmail -i "receiver@mail.com" <<EOF
Subject: Lorem subject...
From: sender@mail.com

Lorem body...
EOF

# EXAMPLE 2
# requires the GNU mail command (apt-get install mailutils)
echo "Lorem body..." | mail -r "sender@mail.com" -s "Lorem subject..." "receiver@mail.com"

```

## Upload and Report Protocol

Uploading and reporting is done via a `POST` request. The URL as well as the headers to be sent are fully configurable via the configuration file `/etc/achelous/upstream.json`.

We provide the following example of the request body for your convenience. You will find formal specifications in the [specs](./specs/) folder.

```json
{
  "id": "01C8EDWZM57VXWY2VJNBNXHA01",
  "timestamp": "2018-03-13T01:59:26.213363817+01:00",
  "participants": {
    "from": {
      "name": "Sender User",
      "email": "sender@mail.com"
    },
    "to": [
      {
        "name": "Receiver User",
        "email": "receiver@mail.com"
      }
    ]
  },
  "subject": "Lorem subject...",
  "body": {
    "text": "Lorem body...\n",
    "html": ""
  },
  "attachments": []
}
```

The reporting request is performed on successful queue runs only. The request body is empty.

## Building

**Note:** Pre-built packages are available [here](https://github.com/coreprocess/achelous/releases). Have a look and skip this section, if you find the appropriate `deb` file for your architecture.

Make sure the following tools are available in your `PATH`:

- `make`
- `gcc`
- `go`
- `glide`
- `dpkg-deb` (for packaging)

Run `make` to build the binaries from source. The binaries will be placed in the `.build/bin` directoriy.

Run `make dist` to build the package file. The `deb` file will be placed in the `.build/dist` directory.

Install achelous with `sudo dpkg -i achelous_*.deb`.
