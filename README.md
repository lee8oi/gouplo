Gouplo is a simple & easy-to-use Go based fileserver with multi-file upload capability.

# Features
* jQuery/Ajax for dynamic form submissions without page reload.
* Automatically uses HTTPS for better security.
* HTTP Basic Authentication for simple login capability.
* Multiple file handling for batch upload jobs (no progress bar).
* Easily configurable via json config file.

# Usage

##### Generate SSL certificates.
On Linux/Unix systems the simplest way to obtain your SSL certificate & key file is by using Certbot
[https://certbot.eff.org/](https://certbot.eff.org/). You'll need the fullchain.pem and the privkey.pem files.

##### Build the server binary.
```Bash
$ go build
```

##### Edit server configuration.
Rename `config.json.example` file to `config.json` and edit the configuration values to suit your needs.

##### Start the server.
Note: If -config is omitted gouplo assumes config.json is in current directory.
```Bash
$ gouplo -config "path/to/config.json"
```

##### Test the server.
Visit your web address in a web browser to test the server and verify your site is using a secure connection.

##### Modify the code.
Once you have successfully tested the server you are ready to edit the code to better suit your needs. Good luck and have fun!
