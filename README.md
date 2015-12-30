Eddystone Scanner in Golang
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/EddystoneScanner/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/EddystoneScanner?status.svg)](https://godoc.org/github.com/kkdai/EddystoneScanner)  [![Build Status](https://travis-ci.org/kkdai/react-diff.svg?branch=master)](https://travis-ci.org/kkdai/EddystoneScanner)

![](images/eddystone.jpg)


Features
---------------


- Scan [Eddystone™](https://github.com/google/eddystone) beacon in multiple environment (see detail requirement in [paypal/gatt](https://github.com/paypal/gatt) support in future)
- Print out detail beacon data from eddystone



Install
---------------
`go get github.com/kkdai/EddystoneScanner`


Run it
---------------

```bash

$./EddystoneScanner

State: PoweredOn
Scanning for eddystone beacon...
=====================================
Eddystone Beacon Found!.....
It is UID beancon, parse data..
Beacon Name: abeacon_4CA2
TxPower: -59
UID: FDA50693-A4E24FB1-AFCFC6EB-07647825
-------------------------------------
Eddystone Beacon Found!.....
It is URL beacon, parse data..
Beacon Name: abeacon_4CA2
TxPower: -59
URL: http://google.com
-------------------------------------
```

Eddystone Beacon simulator
---------------

Using following eddystone beacon simulator if you don't have real one in hand.

- [Beacon simulator for iBeacon and Eddystone](https://github.com/kkdai/beacon)
- [Node.js Eddystone simulator](https://github.com/don/node-eddystone-beacon)


Inspired
---------------

- [https://github.com/paypal/gatt](https://github.com/paypal/gatt)
- [https://github.com/suapapa/go_eddystone](https://github.com/suapapa/go_eddystone)


Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.

