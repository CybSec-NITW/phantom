# phantom ctf platform

## DESCRIPTION

CTF platfom based on Go language

## HOW TO USE

### Windows platform

```bash
git clone https://github.com/CybSec-NITW/phantom.git
cd phantom
go get -u github.com/beego/bee
bee run
```

If it prompts SQLite installation failed, please install[tdm-gcc](http://tdm-gcc.tdragon.net/download) Or other gcc environment.

### Linux platform

```bash
apt-get update
apt-get install gcc
git clone https://github.com/CybSec-NITW/phantom.git
cd phantom
go get -u github.com/beego/bee
bee run
```

## TODO

- [x] Installation page
- [x] User Management
- [x] Problem management
- [x] Contest page
- [ ] flag dynamic anti-cheat
- [ ] Dynamic delivery of containers
- [ ]Beautiful front-end page under container dynamics
- [x] Initialization page (configuration database when first loaded, etc.)
