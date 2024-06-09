# ibeacon-identity

This repository contains a Golang service that identifies and authorizes iBeacon identifiers broadcasted by iOS or Android apps. The identifiers are validated based on a signature calculated using CRC32 checksum with a secret.

The service listens for iBeacon advertisements and extracts the UUID, Major, and Minor from the identifier. It then calculates a signature using the UUID, Major, and Minor values along with a secret. The signature is then compared with the signature received in the identifier. If they match, the identifier is considered valid.
