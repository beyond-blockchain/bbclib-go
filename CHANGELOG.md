CHANGELOG
====

## v1.6.0
* change programming interfaces
  - BBcTransaction and its child classes provide utility methods to create objects.

## v1.5.1
* remove dependency on libbcsig
  - all key manipulation and sign/verify functions become golang native
* unsupport ECC NIST SECP-256 curve

## v1.5.0
* add BBcAssetRaw and BBcAssetHash classes
* the version in BBcTransaction header is 2

## v1.4.4
* bug fix

## v1.4.3
* Add utility to include signature (SignAndAdd function)
  * Note that the function does not work correctly for a transaction with BBcReference

## v1.4.2
* Add key import/export functions in keypair.go

## v1.4.1
* Add installation script (prepare_bbclib.sh)

## v1.4
* ID length configuration support (same as py-bbclib v1.4.1)
* External public key support (same as py-bbclib v1.4.1)
  * BBcSignature having 0-length public key indicates that the public key for verification is given externally.

## v1.3
* not released

## v1.2
* Golang implementation of bbclib.py in BBc-1 version 1.2
  - Cloned from quvox/bbclib-go 

